package deploy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/multycloud/multy/api/errors"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/credspb"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/encoder"
	"github.com/multycloud/multy/flags"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/types"
	"golang.org/x/exp/maps"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	tfFile = "main.tf"
	tfDir  = ".terraform"
)

var (
	AwsCredsNotSetErr   = status.Error(codes.InvalidArgument, "aws credentials are required but not set")
	AzureCredsNotSetErr = status.Error(codes.InvalidArgument, "azure credentials are required but not set")
)

type tfOutput struct {
	Level      string `json:"@level"`
	Message    string `json:"@message"`
	Diagnostic struct {
		Summary string `json:"summary"`
		Detail  string `json:"detail"`
		Address string `json:"address"`
	} `json:"diagnostic"`
}

type EncodedResources struct {
	HclString         string
	affectedResources []string
}

func Encode(credentials *credspb.CloudCredentials, c *resources.MultyConfig, prev resources.Resource, curr resources.Resource) (EncodedResources, error) {
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

	result.affectedResources = maps.Keys(affectedResources)
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

func Deploy(ctx context.Context, c *resources.MultyConfig, prev resources.Resource, curr resources.Resource) (*output.TfState, error) {
	tmpDir := getTempDirForUser(false, c.GetUserId())
	encoded, err := EncodeAndStoreTfFile(ctx, c, prev, curr, false)
	if err != nil {
		return nil, err
	}

	err = MaybeInit(ctx, c.GetUserId(), false)
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

	if len(encoded.affectedResources) != 0 {
		cmd := exec.CommandContext(ctx, "terraform", append([]string{"-chdir=" + tmpDir, "apply", "-refresh=false", "-auto-approve", "--json"}, targetArgs...)...)
		if flags.DryRun {
			cmd = exec.CommandContext(ctx, "terraform", append([]string{"-chdir=" + tmpDir, "plan", "-refresh=false", "--json"}, targetArgs...)...)
		}
		outputJson := new(bytes.Buffer)
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
	}

	state, err := GetState(ctx, c.GetUserId(), false)
	if err != nil {
		return state, errors.InternalServerErrorWithMessage("error parsing state", err)
	}

	return state, nil
}

func EncodeAndStoreTfFile(ctx context.Context, c *resources.MultyConfig, prev resources.Resource, curr resources.Resource, readonly bool) (EncodedResources, error) {
	credentials, err := util.ExtractCloudCredentials(ctx)
	if err != nil {
		return EncodedResources{}, err
	}
	encoded, err := Encode(credentials, c, prev, curr)
	if err != nil {
		return encoded, err
	}

	tfBlock, err := GetTerraformBlock(c.GetUserId())
	if err != nil {
		return encoded, err
	}

	// TODO: move this to a proper place
	hclOutput := tfBlock + encoded.HclString

	tmpDir := getTempDirForUser(readonly, c.GetUserId())
	err = os.MkdirAll(tmpDir, os.ModeDir|(os.ModePerm&0775))
	if err != nil {
		return EncodedResources{}, err
	}
	err = os.WriteFile(filepath.Join(tmpDir, tfFile), []byte(hclOutput), os.ModePerm&0664)
	return encoded, err
}

func getFirstError(outputs []tfOutput) error {
	var err error
	for _, o := range outputs {
		if o.Level == "error" {
			log.Printf("[ERROR] %s\n%s\n", o.Diagnostic.Summary, o.Diagnostic.Detail)
			if err == nil {
				err = fmt.Errorf(o.Diagnostic.Summary)
			}
		}
	}
	return err
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

func MaybeInit(ctx context.Context, userId string, readonly bool) error {
	tmpDir := getTempDirForUser(readonly, userId)
	_, err := os.Stat(filepath.Join(tmpDir, tfDir))
	if os.IsNotExist(err) {
		start := time.Now()
		defer func() {
			log.Printf("[DEBUG] init finished in %s", time.Since(start))
		}()

		err := os.MkdirAll(tmpDir, os.ModeDir|(os.ModePerm&0775))
		if err != nil {
			return errors.InternalServerErrorWithMessage("error creating output file", err)
		}

		cmd := exec.CommandContext(ctx, "terraform", "-chdir="+tmpDir, "init", "-reconfigure")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("unable to initialize terraform")
		}
	} else if err != nil {
		return err
	}
	return nil
}

func GetState(ctx context.Context, userId string, readonly bool) (*output.TfState, error) {
	tmpDir := getTempDirForUser(readonly, userId)
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

func RefreshState(ctx context.Context, userId string, readonly bool) error {
	start := time.Now()
	defer func() {
		log.Printf("[DEBUG] refresh finished in %s", time.Since(start))
	}()

	tmpDir := getTempDirForUser(readonly, userId)
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

func getTempDirForUser(readonly bool, userId string) string {
	tmpDir := filepath.Join(os.TempDir(), "multy", userId)

	if flags.Environment == flags.Local {
		tmpDir = filepath.Join(tmpDir, "local")
	}

	if readonly {
		tmpDir = filepath.Join(tmpDir, "readonly")
	}

	return tmpDir
}
