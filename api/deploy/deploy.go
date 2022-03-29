package deploy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/multycloud/multy/api/converter"
	"github.com/multycloud/multy/api/errors"
	"github.com/multycloud/multy/api/proto/common"
	"github.com/multycloud/multy/api/proto/config"
	"github.com/multycloud/multy/api/proto/creds"
	"github.com/multycloud/multy/api/proto/resources"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/decoder"
	"github.com/multycloud/multy/encoder"
	common_resources "github.com/multycloud/multy/resources"
	cloud_providers "github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	rg "github.com/multycloud/multy/resources/resource_group"
	"github.com/multycloud/multy/resources/types"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func Translate(credentials *creds.CloudCredentials, c *config.Config, prev *config.Resource, curr *config.Resource) (string, error) {
	// TODO: get rid of this translation layer and instead use protos directly
	translated := map[string]common_resources.CloudSpecificResource{}
	for _, r := range c.Resources {
		// TODO: move this to Converters
		resourceMessage := r.ResourceArgs.ResourceArgs
		if resourceMessage.MessageIs(&resources.VirtualNetworkArgs{}) {
			err := addMultyResource(r, translated, &converter.VnConverter{})
			if err != nil {
				return "", err
			}
		} else if resourceMessage.MessageIs(&resources.SubnetArgs{}) {
			err := addMultyResource(r, translated, &converter.SubnetConverter{})
			if err != nil {
				return "", err
			}
		} else if resourceMessage.MessageIs(&resources.NetworkInterfaceArgs{}) {
			err := addMultyResource(r, translated, &converter.NetworkInterfaceConverter{})
			if err != nil {
				return "", err
			}
		} else if resourceMessage.MessageIs(&resources.RouteTableArgs{}) {
			err := addMultyResource(r, translated, &converter.RouteTableConverter{})
			if err != nil {
				return "", err
			}
		} else if resourceMessage.MessageIs(&resources.RouteTableAssociationArgs{}) {
			err := addMultyResource(r, translated, &converter.RouteTableAssociationConverter{})
			if err != nil {
				return "", err
			}
		} else if resourceMessage.MessageIs(&resources.NetworkSecurityGroupArgs{}) {
			err := addMultyResource(r, translated, &converter.NetworkSecurityGroupConverter{})
			if err != nil {
				return "", err
			}
		} else if resourceMessage.MessageIs(&resources.DatabaseArgs{}) {
			err := addMultyResource(r, translated, &converter.DatabaseConverter{})
			if err != nil {
				return "", err
			}
		} else if resourceMessage.MessageIs(&resources.ObjectStorageArgs{}) {
			err := addMultyResource(r, translated, &converter.ObjectStorageConverter{})
			if err != nil {
				return "", err
			}
		} else if resourceMessage.MessageIs(&resources.ObjectStorageObjectArgs{}) {
			err := addMultyResource(r, translated, &converter.ObjectStorageObjectConverter{})
			if err != nil {
				return "", err
			}
		} else if resourceMessage.MessageIs(&resources.PublicIpArgs{}) {
			err := addMultyResource(r, translated, &converter.PublicIpConverter{})
			if err != nil {
				return "", err
			}
		} else if resourceMessage.MessageIs(&resources.KubernetesClusterArgs{}) {
			err := addMultyResource(r, translated, &converter.KubernetesClusterConverter{})
			if err != nil {
				return "", err
			}
		} else if resourceMessage.MessageIs(&resources.KubernetesNodePoolArgs{}) {
			err := addMultyResource(r, translated, &converter.KubernetesNodePoolConverter{})
			if err != nil {
				return "", err
			}
		} else if resourceMessage.MessageIs(&resources.LambdaArgs{}) {
			err := addMultyResource(r, translated, &converter.LambdaConverter{})
			if err != nil {
				return "", err
			}
		} else if resourceMessage.MessageIs(&resources.VaultArgs{}) {
			err := addMultyResource(r, translated, &converter.VaultConverter{})
			if err != nil {
				return "", err
			}
		} else if resourceMessage.MessageIs(&resources.VaultAccessPolicyArgs{}) {
			err := addMultyResource(r, translated, &converter.VaultAccessPolicyConverter{})
			if err != nil {
				return "", err
			}
		} else if resourceMessage.MessageIs(&resources.VaultSecretArgs{}) {
			err := addMultyResource(r, translated, &converter.VaultSecretConverter{})
			if err != nil {
				return "", err
			}
		} else if resourceMessage.MessageIs(&resources.VirtualMachineArgs{}) {
			err := addMultyResource(r, translated, &converter.VirtualMachineConverter{})
			if err != nil {
				return "", err
			}
		} else {
			return "", fmt.Errorf("unknown resource type %s", resourceMessage.MessageName())
		}
	}

	provider, err := getExistingProvider(prev)
	if err != nil {
		return "", err
	}

	// TODO add s3 state file backend
	decodedResources := decoder.DecodedResources{
		Resources: translated,
		GlobalConfig: decoder.DecodedGlobalConfig{
			Location:      "ireland",
			Clouds:        cloud_providers.GetAllCloudProviders(),
			DefaultRgName: rg.GetDefaultResourceGroupId(),
		},
		Providers: provider,
	}

	hclOutput, errs, err := encoder.Encode(&decodedResources, credentials)
	if len(errs) > 0 {
		return hclOutput, errors.ValidationErrors(errs)
	}
	if err != nil {
		return hclOutput, err
	}

	for _, r := range translated {
		if string(r.Cloud) == "aws" && (credentials.GetAwsCreds().GetAccessKey() == "" || credentials.GetAwsCreds().GetSecretKey() == "") {
			return hclOutput, AwsCredsNotSetErr
		}
		if string(r.Cloud) == "azure" && (credentials.GetAzureCreds().GetSubscriptionId() == "" ||
			credentials.GetAzureCreds().GetClientId() == "" ||
			credentials.GetAzureCreds().GetTenantId() == "" ||
			credentials.GetAzureCreds().GetClientSecret() == "") {
			return hclOutput, AzureCredsNotSetErr
		}
	}

	return hclOutput, nil
}

func Deploy(ctx context.Context, c *config.Config, prev *config.Resource, curr *config.Resource) (*output.TfState, error) {
	credentials, err := util.ExtractCloudCredentials(ctx)
	if err != nil {
		return nil, err
	}
	hclOutput, err := Translate(credentials, c, prev, curr)
	if err != nil {
		return nil, err
	}

	// TODO: move this to a proper place
	hclOutput = RequiredProviders + hclOutput

	fmt.Println(hclOutput)

	tmpDir := filepath.Join(os.TempDir(), "multy", c.UserId)
	err = os.WriteFile(filepath.Join(tmpDir, tfFile), []byte(hclOutput), os.ModePerm&0664)
	if err != nil {
		return nil, errors.InternalServerErrorWithMessage("error storing configuration", err)
	}

	err = MaybeInit(c.UserId)
	if err != nil {
		return nil, err
	}

	fmt.Println("Running tf apply")
	startApply := time.Now()

	// TODO: only deploy targets given in the args
	// TODO: parse errors and send them to user
	cmd := exec.Command("terraform", "-chdir="+tmpDir, "apply", "-auto-approve")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return nil, errors.InternalServerErrorWithMessage("error deploying resources", err)
	}
	log.Printf("tf apply ended in %s", time.Since(startApply))

	state, err := GetState(c.UserId)
	if err != nil {
		return state, errors.InternalServerErrorWithMessage("error parsing state", err)
	}

	return state, nil
}

func MaybeInit(userId string) error {
	tmpDir := filepath.Join(os.TempDir(), "multy", userId)
	_, err := os.Stat(filepath.Join(tmpDir, tfDir))
	if os.IsNotExist(err) {
		fmt.Println("running tf init")
		startInit := time.Now()

		cmd := exec.Command("terraform", "-chdir="+tmpDir, "init")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			return err
		}
		log.Printf("tf init ended in %s", time.Since(startInit))
	} else if err != nil {
		return err
	}
	return nil
}

func GetState(userId string) (*output.TfState, error) {
	tmpDir := filepath.Join(os.TempDir(), "multy", userId)
	state := output.TfState{}
	stateJson := new(bytes.Buffer)
	cmd := exec.Command("terraform", "-chdir="+tmpDir, "show", "-json")
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

type hasCommonArgs interface {
	GetCommonParameters() *common.ResourceCommonArgs
}

func addMultyResource(r *config.Resource, translated map[string]common_resources.CloudSpecificResource, c converter.MultyResourceConverter) error {
	m, err := r.ResourceArgs.ResourceArgs.UnmarshalNew()
	if err != nil {
		return err
	}
	// TODO: refactor this
	if commonArgs, ok := m.(hasCommonArgs); ok {
		if commonArgs.GetCommonParameters().ResourceGroupId == "" {
			rgId, _ := decoder.GetRgId(rg.GetDefaultResourceGroupId(), nil, &hcl.EvalContext{
				Variables: map[string]cty.Value{},
				Functions: map[string]function.Function{},
			}, c.GetResourceType())
			if rgId != "" {
				commonArgs.GetCommonParameters().ResourceGroupId = rgId
				resourceGroup := common_resources.CloudSpecificResource{
					ImplicitlyCreated: true,
					Cloud:             cloud_providers.CloudProvider(strings.ToLower(commonArgs.GetCommonParameters().CloudProvider.String())),
					Resource: &rg.Type{
						ResourceId: rgId,
						Name:       rgId,
						Location:   strings.ToLower(commonArgs.GetCommonParameters().Location.String()),
					},
				}
				translated[resourceGroup.GetResourceId()] = resourceGroup
			}
		}
	}

	translatedResource, err := c.ConvertToMultyResource(r.ResourceId, m, translated)
	if err != nil {
		return err
	}
	translated[translatedResource.Resource.GetResourceId()] = translatedResource

	return nil
}

type WithCommonParams interface {
	GetCommonParameters() *common.ResourceCommonArgs
}

func getExistingProvider(r *config.Resource) (map[cloud_providers.CloudProvider]map[string]*types.Provider, error) {
	if r != nil {
		args := r.GetResourceArgs().GetResourceArgs()
		m, err := args.UnmarshalNew()
		if err != nil {
			return nil, err
		}
		if wcp, ok := m.(WithCommonParams); ok {
			cloud := cloud_providers.CloudProvider(strings.ToLower(wcp.GetCommonParameters().CloudProvider.String()))
			location, err := cloud_providers.GetCloudLocation(strings.ToLower(wcp.GetCommonParameters().Location.String()), cloud)
			if err != nil {
				return nil, err
			}
			return map[cloud_providers.CloudProvider]map[string]*types.Provider{
				cloud: {
					location: &types.Provider{
						Cloud:             cloud,
						Location:          location,
						IsDefaultProvider: false,
						NumResources:      1,
					},
				},
			}, nil

		}

	}

	return map[cloud_providers.CloudProvider]map[string]*types.Provider{}, nil
}
