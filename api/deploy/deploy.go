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

func Encode(credentials *credspb.CloudCredentials, c *configpb.Config, prev *configpb.Resource, curr *configpb.Resource) (string, error) {
	decodedResources, err := GetResources(c, prev)
	if err != nil {
		return "", err
	}

	hclOutput, errs, err := encoder.Encode(decodedResources, credentials)
	if len(errs) > 0 {
		return hclOutput, errors.ValidationErrors(errs)
	}
	if err != nil {
		return hclOutput, err
	}

	for _, r := range decodedResources.Resources {
		if r.GetCloud() == commonpb.CloudProvider_AWS && (credentials.GetAwsCreds().GetAccessKey() == "" || credentials.GetAwsCreds().GetSecretKey() == "") {
			return hclOutput, AwsCredsNotSetErr
		}
		if r.GetCloud() == commonpb.CloudProvider_AZURE && (credentials.GetAzureCreds().GetSubscriptionId() == "" ||
			credentials.GetAzureCreds().GetClientId() == "" ||
			credentials.GetAzureCreds().GetTenantId() == "" ||
			credentials.GetAzureCreds().GetClientSecret() == "") {
			return hclOutput, AzureCredsNotSetErr
		}
	}

	return hclOutput, nil
}

func GetResources(c *configpb.Config, prev *configpb.Resource) (*encoder.DecodedResources, error) {
	translated := map[string]resources.Resource{}

	for _, r := range c.Resources {
		resourceMessage := r.ResourceArgs.ResourceArgs
		added := false
		for messageType, conv := range converter.Converters {
			if resourceMessage.MessageIs(messageType) {
				err := addMultyResourceNew(r, translated, conv)
				if err != nil {
					return nil, err
				}
				added = true
				break
			}
		}
		if !added {
			return nil, fmt.Errorf("unknown resource type %s", resourceMessage.MessageName())
		}
	}

	provider, err := getExistingProvider(prev)
	if err != nil {
		return nil, err
	}

	// TODO add s3 state file backend
	decodedResources := encoder.DecodedResources{
		Resources: translated,
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
	credentials, err := util.ExtractCloudCredentials(ctx)
	if err != nil {
		return nil, err
	}
	hclOutput, err := Encode(credentials, c, prev, curr)
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
	outputJson := new(bytes.Buffer)
	cmd := exec.Command("terraform", "-chdir="+tmpDir, "apply", "-auto-approve", "--json")
	cmd.Stdout = outputJson
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		outputs, parseErr := parseTfOutputs(outputJson)
		if parseErr != nil {
			return nil, errors.InternalServerErrorWithMessage("error deploying resources", parseErr)
		}
		if parseErr := getFirstError(outputs); parseErr != nil {
			return nil, errors.InternalServerErrorWithMessage("error deploying resources", parseErr)
		}

		fmt.Println(outputs)

		return nil, errors.InternalServerErrorWithMessage("error deploying resources", err)
	}
	log.Printf("tf apply ended in %s", time.Since(startApply))

	state, err := GetState(c.UserId)
	if err != nil {
		return state, errors.InternalServerErrorWithMessage("error parsing state", err)
	}

	return state, nil
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
	GetCommonParameters() *commonpb.ResourceCommonArgs
}

func addMultyResourceNew(r *configpb.Resource, translated map[string]resources.Resource, metadata converter.ResourceMetadata) error {
	m, err := r.ResourceArgs.ResourceArgs.UnmarshalNew()
	if err != nil {
		return err
	}
	// TODO: refactor this
	if commonArgs, ok := m.(hasCommonArgs); ok {
		if commonArgs.GetCommonParameters().ResourceGroupId == "" {
			rgId := rg.GetDefaultResourceGroupIdString(metadata.AbbreviatedName)
			if rgId != "" {
				commonArgs.GetCommonParameters().ResourceGroupId = rgId
				if commonArgs.GetCommonParameters().CloudProvider == commonpb.CloudProvider_AZURE {
					resourceGroup :=
						&rg.Type{
							ResourceId: rgId,
							Name:       rgId,
							Location:   strings.ToLower(commonArgs.GetCommonParameters().Location.String()),
							Cloud:      commonArgs.GetCommonParameters().CloudProvider,
						}

					translated[resourceGroup.GetResourceId()] = resourceGroup
				}
			}
		}
	}

	translatedResource, err := metadata.InitFunc(r.ResourceId, m, translated)
	if err != nil {
		return err
	}
	translated[translatedResource.GetResourceId()] = translatedResource

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
