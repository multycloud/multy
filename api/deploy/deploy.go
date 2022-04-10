package deploy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/multycloud/multy/api/converter"
	"github.com/multycloud/multy/api/errors"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/configpb"
	"github.com/multycloud/multy/api/proto/credspb"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/encoder"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	rg "github.com/multycloud/multy/resources/resource_group"
	"github.com/multycloud/multy/resources/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	tfFile  = "main.tf"
	tfDir   = ".terraform"
	tfState = "terraform.tfstate"
)

var (
	AwsCredsNotSetErr   = status.Error(codes.InvalidArgument, "aws credentials are required but not set")
	AzureCredsNotSetErr = status.Error(codes.InvalidArgument, "azure credentials are required but not set")
)

type EncodedResources struct {
	HclString         string
	affectedResources []string
}

func Encode(credentials *credspb.CloudCredentials, c *configpb.Config, prev *configpb.Resource, curr *configpb.Resource) (EncodedResources, error) {
	result := EncodedResources{}
	decodedResources, err := GetResources(c, prev)
	if err != nil {
		return result, err
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
		if r.GetCloud() == commonpb.CloudProvider_AWS && (credentials.GetAwsCreds().GetAccessKey() == "" || credentials.GetAwsCreds().GetSecretKey() == "") {
			return result, AwsCredsNotSetErr
		}
		if r.GetCloud() == commonpb.CloudProvider_AZURE && (credentials.GetAzureCreds().GetSubscriptionId() == "" ||
			credentials.GetAzureCreds().GetClientId() == "" ||
			credentials.GetAzureCreds().GetTenantId() == "" ||
			credentials.GetAzureCreds().GetClientSecret() == "") {
			return result, AzureCredsNotSetErr
		}
	}

	result.affectedResources = updateMultyResourceGroups(decodedResources, encoded, c, prev, curr)

	return result, nil
}

func updateMultyResourceGroups(decodedResources *encoder.DecodedResources, encoded encoder.EncodedResources, c *configpb.Config, prev *configpb.Resource, curr *configpb.Resource) []string {
	var result []string
	if prev != nil {
		result = append(result, prev.DeployedResourceGroup.DeployedResource...)
	}

	resourcespbById := map[string]*configpb.Resource{}
	for _, resource := range c.Resources {
		resourcespbById[resource.ResourceId] = resource
	}
	deployedResourcesByGroupId := map[string]*configpb.DeployedResourceGroup{}
	for groupId, group := range decodedResources.Resources.GetMultyResourceGroups() {
		for _, resource := range group.Resources {
			for _, deployedResource := range encoded.DeployedResources[resource] {
				if _, ok := deployedResourcesByGroupId[groupId]; !ok {
					deployedResourcesByGroupId[groupId] = &configpb.DeployedResourceGroup{GroupId: groupId}
				}
				deployedResourcesByGroupId[groupId].DeployedResource = append(deployedResourcesByGroupId[groupId].DeployedResource, deployedResource)
			}
			if _, ok := resourcespbById[resource.GetResourceId()]; ok {
				resourcespbById[resource.GetResourceId()].DeployedResourceGroup = deployedResourcesByGroupId[groupId]
			} else {
				log.Printf("[DEBUG] not adding %s to a group, it doesn't exist in the state\n", resource.GetResourceId())
			}
		}
	}

	if curr != nil {
		result = append(result, resourcespbById[curr.ResourceId].DeployedResourceGroup.DeployedResource...)
	}

	return result
}

func GetResources(c *configpb.Config, prev *configpb.Resource) (*encoder.DecodedResources, error) {
	res := resources.NewResources()

	for _, r := range c.Resources {
		conv, err := converter.GetConverter(r.ResourceArgs.ResourceArgs.MessageName())
		if err != nil {
			return nil, err
		}
		err = addMultyResourceNew(r, res, conv)
		if err != nil {
			return nil, err
		}
	}

	provider, err := getExistingProvider(prev)
	if err != nil {
		return nil, err
	}

	// TODO add s3 state file backend
	decodedResources := encoder.DecodedResources{
		Resources: res,
		Providers: provider,
	}

	return &decodedResources, nil
}

type tfOutput struct {
	Level      string `json:"@level"`
	Message    string `json:"@message"`
	Diagnostic struct {
		Summary string `json:"summary"`
		Detail  string `json:"detail"`
	} `json:"diagnostic"`
}

func Deploy(ctx context.Context, c *configpb.Config, prev *configpb.Resource, curr *configpb.Resource) (*output.TfState, error) {
	tmpDir := filepath.Join(os.TempDir(), "multy", c.UserId)
	encoded, err := EncodeAndStoreTfFile(ctx, c, prev, curr)
	if err != nil {
		return nil, err
	}

	err = MaybeInit(ctx, c.UserId)
	if err != nil {
		return nil, err
	}

	var targetArgs []string

	log.Println("[INFO] Running apply for targets:")
	for _, id := range encoded.affectedResources {
		log.Printf("[INFO] %s", id)
		targetArgs = append(targetArgs, "-target="+id)
	}

	start := time.Now()
	defer func() {
		log.Printf("[DEBUG] apply finished in %s", time.Since(start))
	}()

	// TODO: only deploy targets given in the args
	outputJson := new(bytes.Buffer)
	cmd := exec.CommandContext(ctx, "terraform", append([]string{"-chdir=" + tmpDir, "apply", "-refresh=false", "-auto-approve", "--json"}, targetArgs...)...)
	cmd.Stdout = outputJson
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		outputs, parseErr := parseTfOutputs(outputJson)
		if parseErr != nil {
			return nil, errors.InternalServerErrorWithMessage("error deploying resources", parseErr)
		}
		if parseErr := getFirstError(outputs); parseErr != nil {
			return nil, errors.DeployError(parseErr)
		}
		return nil, errors.InternalServerErrorWithMessage("error deploying resources", err)
	}

	state, err := GetState(ctx, c.UserId)
	if err != nil {
		return state, errors.InternalServerErrorWithMessage("error parsing state", err)
	}

	return state, nil
}

func EncodeAndStoreTfFile(ctx context.Context, c *configpb.Config, prev *configpb.Resource, curr *configpb.Resource) (EncodedResources, error) {
	credentials, err := util.ExtractCloudCredentials(ctx)
	if err != nil {
		return EncodedResources{}, err
	}
	encoded, err := Encode(credentials, c, prev, curr)
	if err != nil {
		return encoded, err
	}

	// TODO: move this to a proper place
	hclOutput := RequiredProviders + encoded.HclString

	tmpDir := filepath.Join(os.TempDir(), "multy", c.UserId)
	err = os.WriteFile(filepath.Join(tmpDir, tfFile), []byte(hclOutput), os.ModePerm&0664)
	return encoded, err
}

func getFirstError(outputs []tfOutput) error {
	for _, o := range outputs {
		if o.Level == "error" {
			return fmt.Errorf(o.Diagnostic.Summary)
		}
	}
	return nil
}

func parseTfOutputs(outputJson *bytes.Buffer) ([]tfOutput, error) {
	var out []tfOutput
	line, err := outputJson.ReadString('\n')
	for ; err == nil; line, err = outputJson.ReadString('\n') {
		elem := tfOutput{}
		err = json.Unmarshal([]byte(line), &elem)
		if err != nil {
			return nil, err
		}
		out = append(out, elem)
	}

	if err == io.EOF {
		return out, nil
	}

	return nil, err
}

func MaybeInit(ctx context.Context, userId string) error {
	tmpDir := filepath.Join(os.TempDir(), "multy", userId)
	_, err := os.Stat(filepath.Join(tmpDir, tfDir))
	if os.IsNotExist(err) {
		start := time.Now()
		defer func() {
			log.Printf("[DEBUG] init finished in %s", time.Since(start))
		}()

		cmd := exec.CommandContext(ctx, "terraform", "-chdir="+tmpDir, "init")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}

func GetState(ctx context.Context, userId string) (*output.TfState, error) {
	tmpDir := filepath.Join(os.TempDir(), "multy", userId)
	state := output.TfState{}
	stateJson := new(bytes.Buffer)
	cmd := exec.CommandContext(ctx, "terraform", "-chdir="+tmpDir, "show", "-json")
	cmd.Stdout = stateJson
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return &state, err
	}

	err = json.Unmarshal(stateJson.Bytes(), &state)
	if err != nil {
		return nil, err
	}
	return &state, err
}

func RefreshState(ctx context.Context, userId string) error {
	start := time.Now()
	defer func() {
		log.Printf("[DEBUG] refresh finished in %s", time.Since(start))
	}()
	tmpDir := filepath.Join(os.TempDir(), "multy", userId)
	outputJson := new(bytes.Buffer)
	cmd := exec.CommandContext(ctx, "terraform", "-chdir="+tmpDir, "refresh", "-json")
	cmd.Stdout = outputJson
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		outputs, parseErr := parseTfOutputs(outputJson)
		if parseErr != nil {
			return errors.InternalServerErrorWithMessage("error querying resources", parseErr)
		}
		if parseErr := getFirstError(outputs); parseErr != nil {
			return errors.InternalServerErrorWithMessage("error querying resources", parseErr)
		}
		return errors.InternalServerErrorWithMessage("error querying resources", err)
	}
	return err
}

type hasCommonArgs interface {
	GetCommonParameters() *commonpb.ResourceCommonArgs
}

func addMultyResourceNew(r *configpb.Resource, res resources.Resources, metadata *converter.ResourceMetadata) error {
	m, err := r.ResourceArgs.ResourceArgs.UnmarshalNew()
	if err != nil {
		return err
	}
	// TODO: refactor this
	var resourceGroup *rg.Type
	if commonArgs, ok := m.(hasCommonArgs); ok {
		if commonArgs.GetCommonParameters().ResourceGroupId == "" {
			rgId := rg.GetDefaultResourceGroupIdString(metadata.AbbreviatedName)
			if rgId != "" {
				commonArgs.GetCommonParameters().ResourceGroupId = rgId
				if commonArgs.GetCommonParameters().CloudProvider == commonpb.CloudProvider_AZURE {
					resourceGroup =
						&rg.Type{
							ResourceId: rgId,
							Name:       rgId,
							Location:   strings.ToLower(commonArgs.GetCommonParameters().Location.String()),
							Cloud:      commonArgs.GetCommonParameters().CloudProvider,
						}

					res.ResourceMap[resourceGroup.GetResourceId()] = resourceGroup
				}
			}
		}
	}

	translatedResource, err := metadata.InitFunc(r.ResourceId, m, res)
	if err != nil {
		return err
	}
	res.ResourceMap[translatedResource.GetResourceId()] = translatedResource
	if resourceGroup != nil {
		res.AddDependency(translatedResource.GetResourceId(), resourceGroup.GetResourceId())
	}
	return nil
}

type WithCommonParams interface {
	GetCommonParameters() *commonpb.ResourceCommonArgs
}

func getExistingProvider(r *configpb.Resource) (map[commonpb.CloudProvider]map[string]*types.Provider, error) {
	if r != nil {
		args := r.GetResourceArgs().GetResourceArgs()
		m, err := args.UnmarshalNew()
		if err != nil {
			return nil, err
		}
		if wcp, ok := m.(WithCommonParams); ok {
			location, err := common.GetCloudLocation(strings.ToLower(wcp.GetCommonParameters().Location.String()), wcp.GetCommonParameters().CloudProvider)
			if err != nil {
				return nil, err
			}
			return map[commonpb.CloudProvider]map[string]*types.Provider{
				wcp.GetCommonParameters().CloudProvider: {
					location: &types.Provider{
						Cloud:             wcp.GetCommonParameters().CloudProvider,
						Location:          location,
						IsDefaultProvider: false,
						NumResources:      1,
					},
				},
			}, nil

		}

	}

	return map[commonpb.CloudProvider]map[string]*types.Provider{}, nil
}
