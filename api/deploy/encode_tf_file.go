package deploy

import (
	"github.com/multycloud/multy/api/errors"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/credspb"
	"github.com/multycloud/multy/encoder"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/types"
	"github.com/multycloud/multy/util"
)

func EncodeTfFile(credentials *credspb.CloudCredentials, c *resources.MultyConfig, prev resources.Resource, curr resources.Resource) (EncodedResources, error) {
	result := EncodedResources{}

	provider, err := getExistingProvider(prev, credentials)
	if err != nil {
		return result, err
	}

	decodedResources := &encoder.DecodedResources{
		Resources: c.Resources,
		Providers: provider,
	}

	encoded, err := encoder.Encode(decodedResources, credentials)
	if err != nil {
		return result, err
	}
	if len(encoded.ValidationErrs) > 0 {
		return result, errors.ValidationErrors(encoded.ValidationErrs)
	}

	result.HclString = encoded.HclString
	for _, r := range decodedResources.Resources.ResourceMap {
		if r.GetCloud() == commonpb.CloudProvider_AWS && !hasValidAwsCreds(credentials) {
			return result, AwsCredsNotSetErr
		}
		if r.GetCloud() == commonpb.CloudProvider_AZURE && !hasValidAzureCreds(credentials) {
			return result, AzureCredsNotSetErr
		}
		if r.GetCloud() == commonpb.CloudProvider_GCP && !hasValidGcpCreds(credentials) {
			return result, GcpCredsNotSetErr
		}
	}

	affectedResources := map[string]struct{}{}
	if prev != nil {
		for _, dr := range c.GetAffectedResources(prev.GetResourceId()) {
			affectedResources[dr] = struct{}{}
		}
	}

	c.UpdateMultyResourceGroups()
	c.UpdateDeployedResourceList(encoded.DeployedResources)

	if curr != nil {
		for _, dr := range c.GetAffectedResources(curr.GetResourceId()) {
			affectedResources[dr] = struct{}{}
		}
	}

	result.affectedResources = util.SortedKeys(affectedResources)
	return result, nil
}

func hasValidAzureCreds(credentials *credspb.CloudCredentials) bool {
	return credentials.GetAzureCreds().GetSubscriptionId() != "" &&
		credentials.GetAzureCreds().GetClientId() != "" &&
		credentials.GetAzureCreds().GetTenantId() != "" &&
		credentials.GetAzureCreds().GetClientSecret() != ""
}

func hasValidAwsCreds(credentials *credspb.CloudCredentials) bool {
	return credentials.GetAwsCreds().GetAccessKey() != "" && credentials.GetAwsCreds().GetSecretKey() != ""
}

func hasValidGcpCreds(credentials *credspb.CloudCredentials) bool {
	return credentials.GetGcpCreds().GetCredentials() != ""
}

func getExistingProvider(r resources.Resource, creds *credspb.CloudCredentials) (map[commonpb.CloudProvider]map[string]*types.Provider, error) {
	providers := map[commonpb.CloudProvider]map[string]*types.Provider{}
	if r != nil {
		location := r.GetCloudSpecificLocation()
		providers[r.GetCloud()] = map[string]*types.Provider{
			location: {
				Cloud:        r.GetCloud(),
				Location:     location,
				NumResources: 1,
			},
		}
	}

	// Here we use a default location so that if there are lingering resources in the state we don't throw an error.
	// It doesn't work perfectly tho -- AWS resources will be removed by terraform from the state if they don't exist
	// in our config and will no longer be tracked.
	defaultAzureLocation := common.LOCATION[commonpb.Location_EU_WEST_1][commonpb.CloudProvider_AZURE]
	defaultAwsLocation := common.LOCATION[commonpb.Location_EU_WEST_1][commonpb.CloudProvider_AWS]

	if hasValidAwsCreds(creds) && providers[commonpb.CloudProvider_AZURE] == nil {
		providers[commonpb.CloudProvider_AZURE] = map[string]*types.Provider{
			defaultAzureLocation: {
				Cloud:        commonpb.CloudProvider_AZURE,
				Location:     defaultAzureLocation,
				NumResources: 1,
			},
		}
	}
	if hasValidAzureCreds(creds) && providers[commonpb.CloudProvider_AWS] == nil {
		providers[commonpb.CloudProvider_AWS] = map[string]*types.Provider{
			defaultAwsLocation: {
				Cloud:        commonpb.CloudProvider_AWS,
				Location:     defaultAwsLocation,
				NumResources: 1,
			},
		}
	}

	return providers, nil
}
