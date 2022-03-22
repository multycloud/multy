package deploy

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/multycloud/multy/api/converter"
	"github.com/multycloud/multy/api/proto/common"
	"github.com/multycloud/multy/api/proto/config"
	"github.com/multycloud/multy/api/proto/resources"
	"github.com/multycloud/multy/decoder"
	"github.com/multycloud/multy/encoder"
	common_resources "github.com/multycloud/multy/resources"
	cloud_providers "github.com/multycloud/multy/resources/common"
	rg "github.com/multycloud/multy/resources/resource_group"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"google.golang.org/protobuf/proto"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	tfFile  = "main.tf"
	tfState = "terraform.tfstate"
)

func Translate(c *config.Config) (string, error) {
	// TODO: get rid of this translation layer and instead use protos directly
	translated := map[string]common_resources.CloudSpecificResource{}
	for _, r := range c.Resources {
		if len(r.ResourceArgs.ResourceArgs) == 0 {
			continue
		}

		// TODO: move this to Converters
		resourceMessage := r.ResourceArgs.ResourceArgs[0]
		if resourceMessage.MessageIs(&resources.CloudSpecificVirtualNetworkArgs{}) {
			err := addMultyResource(r, translated, &converter.VnConverter{})
			if err != nil {
				return "", err
			}
		} else if resourceMessage.MessageIs(&resources.CloudSpecificSubnetArgs{}) {
			err := addMultyResource(r, translated, &converter.SubnetConverter{})
			if err != nil {
				return "", err
			}
		} else if resourceMessage.MessageIs(&resources.CloudSpecificNetworkInterfaceArgs{}) {
			err := addMultyResource(r, translated, &converter.NetworkInterfaceConverter{})
			if err != nil {
				return "", err
			}
		} else if resourceMessage.MessageIs(&resources.CloudSpecificRouteTableArgs{}) {
			err := addMultyResource(r, translated, &converter.RouteTableConverter{})
			if err != nil {
				return "", err
			}
		} else if resourceMessage.MessageIs(&resources.CloudSpecificRouteTableAssociationArgs{}) {
			err := addMultyResource(r, translated, &converter.RouteTableAssociationConverter{})
			if err != nil {
				return "", err
			}
		} else if resourceMessage.MessageIs(&resources.CloudSpecificNetworkSecurityGroupArgs{}) {
			err := addMultyResource(r, translated, &converter.NetworkSecurityGroupConverter{})
			if err != nil {
				return "", err
			}
		} else if resourceMessage.MessageIs(&resources.CloudSpecificDatabaseArgs{}) {
			err := addMultyResource(r, translated, &converter.DatabaseConverter{})
			if err != nil {
				return "", err
			}
		} else if resourceMessage.MessageIs(&resources.CloudSpecificObjectStorageArgs{}) {
			err := addMultyResource(r, translated, &converter.ObjectStorageConverter{})
			if err != nil {
				return "", err
			}
		} else if resourceMessage.MessageIs(&resources.CloudSpecificObjectStorageObjectArgs{}) {
			err := addMultyResource(r, translated, &converter.ObjectStorageObjectConverter{})
			if err != nil {
				return "", err
			}
		} else if resourceMessage.MessageIs(&resources.CloudSpecificPublicIpArgs{}) {
			err := addMultyResource(r, translated, &converter.PublicIpConverter{})
			if err != nil {
				return "", err
			}
		} else if resourceMessage.MessageIs(&resources.CloudSpecificKubernetesClusterArgs{}) {
			err := addMultyResource(r, translated, &converter.KubernetesClusterConverter{})
			if err != nil {
				return "", err
			}
		} else if resourceMessage.MessageIs(&resources.CloudSpecificKubernetesNodePoolArgs{}) {
			err := addMultyResource(r, translated, &converter.KubernetesNodePoolConverter{})
			if err != nil {
				return "", err
			}
		} else if resourceMessage.MessageIs(&resources.CloudSpecificLambdaArgs{}) {
			err := addMultyResource(r, translated, &converter.LambdaConverter{})
			if err != nil {
				return "", err
			}
		} else if resourceMessage.MessageIs(&resources.CloudSpecificVaultArgs{}) {
			err := addMultyResource(r, translated, &converter.VaultConverter{})
			if err != nil {
				return "", err
			}
		} else if resourceMessage.MessageIs(&resources.CloudSpecificVaultAccessPolicyArgs{}) {
			err := addMultyResource(r, translated, &converter.VaultAccessPolicyConverter{})
			if err != nil {
				return "", err
			}
		} else if resourceMessage.MessageIs(&resources.CloudSpecificVaultSecretArgs{}) {
			err := addMultyResource(r, translated, &converter.VaultSecretConverter{})
			if err != nil {
				return "", err
			}
		} else if resourceMessage.MessageIs(&resources.CloudSpecificVirtualMachineArgs{}) {
			err := addMultyResource(r, translated, &converter.VirtualMachineConverter{})
			if err != nil {
				return "", err
			}
		} else {
			return "", fmt.Errorf("unknown resource type %s", resourceMessage.MessageName())
		}
	}

	// TODO add s3 state file backend
	decodedResources := decoder.DecodedResources{
		Resources: translated,
		GlobalConfig: decoder.DecodedGlobalConfig{
			Location:      "ireland",
			Clouds:        cloud_providers.GetAllCloudProviders(),
			DefaultRgName: rg.GetDefaultResourceGroupId(),
		},
	}

	hclOutput := encoder.Encode(&decodedResources)

	return hclOutput, nil
}

func Deploy(c *config.Config, resourceId string) error {

	hclOutput, err := Translate(c)
	if err != nil {
		return err
	}

	fmt.Println(hclOutput)

	tmpDir := filepath.Join(os.TempDir(), "multy", c.UserId)
	err = os.WriteFile(filepath.Join(tmpDir, tfFile), []byte(hclOutput), os.ModePerm&0664)
	if err != nil {
		return fmt.Errorf("error creating output file: %s", err.Error())
	}

	fmt.Println("running tf init")

	cmd := exec.Command("terraform", "-chdir="+tmpDir, "init")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}

	fmt.Println("Running tf apply")

	// TODO: only deploy targets given in the args
	// TODO: parse errors and send them to user
	cmd = exec.Command("terraform", "-chdir="+tmpDir, "apply", "-auto-approve")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}

	// TODO: store this in S3
	_, err = os.ReadFile(filepath.Join(tmpDir, tfState))
	//if err == nil {
	//	fmt.Println(string(stateBytes))
	//}

	return nil
}

type hasCommonArgs interface {
	GetCommonParameters() *common.CloudSpecificResourceCommonArgs
}

func addMultyResource(r *config.Resource, translated map[string]common_resources.CloudSpecificResource, c converter.MultyResourceConverter) error {
	var allResources []proto.Message
	for _, args := range r.ResourceArgs.ResourceArgs {
		m := c.NewArg()
		err := args.UnmarshalTo(m)
		if err != nil {
			return err
		}
		allResources = append(allResources, m)
		// TODO: refactor this
		if commonArgs, ok := m.(hasCommonArgs); ok {
			if commonArgs.GetCommonParameters().ResourceGroupId == "" {
				rgId, _ := decoder.GetRgId(rg.GetDefaultResourceGroupId(), nil, &hcl.EvalContext{
					Variables: map[string]cty.Value{},
					Functions: map[string]function.Function{},
				}, c.GetResourceType())
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
	for _, cloudR := range allResources {
		translatedResource, err := c.ConvertToMultyResource(r.ResourceId, cloudR, translated)
		if err != nil {
			return err
		}
		translated[translatedResource.GetResourceId()] = translatedResource
	}
	return nil
}
