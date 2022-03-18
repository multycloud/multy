package deploy

import (
	"fmt"
	"github.com/multycloud/multy/api/converter"
	"github.com/multycloud/multy/api/proto/config"
	"github.com/multycloud/multy/api/proto/resources"
	"github.com/multycloud/multy/decoder"
	"github.com/multycloud/multy/encoder"
	common_resources "github.com/multycloud/multy/resources"
	cloud_providers "github.com/multycloud/multy/resources/common"
	rg "github.com/multycloud/multy/resources/resource_group"
	"google.golang.org/protobuf/proto"
	"os"
	"os/exec"
	"path/filepath"
)

func Deploy(c *config.Config, resourceId string) error {
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
				return err
			}
		} else if resourceMessage.MessageIs(&resources.CloudSpecificSubnetArgs{}) {
			err := addMultyResource(r, translated, &converter.SubnetConverter{})
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("unknown resource type %s", resourceMessage.MessageName())
		}
	}

	decodedResources := decoder.DecodedResources{
		Resources: translated,
		GlobalConfig: decoder.DecodedGlobalConfig{
			Location:      "ireland",
			Clouds:        cloud_providers.GetAllCloudProviders(),
			DefaultRgName: rg.GetDefaultResourceGroupId(),
		},
	}

	hclOutput := encoder.Encode(&decodedResources)

	fmt.Println(hclOutput)

	// TODO: different dirs for different users
	tmpDir := filepath.Join(os.TempDir(), "multy")
	err := os.MkdirAll(tmpDir, os.ModeDir|(os.ModePerm&0775))
	if err != nil {
		return fmt.Errorf("error creating output file: %s", err.Error())
	}
	err = os.WriteFile(filepath.Join(tmpDir, "main.tf"), []byte(hclOutput), os.ModePerm&0664)
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

	fmt.Println("running tf apply")

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
	stateBytes, err := os.ReadFile(filepath.Join(tmpDir, "terraform.tfstate"))
	if err == nil {
		fmt.Println(string(stateBytes))
	}

	return nil
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
