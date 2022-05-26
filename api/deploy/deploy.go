package deploy

import (
	"context"
	"github.com/multycloud/multy/api/errors"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/flags"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"os"
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

type EncodedResources struct {
	HclString         string
	affectedResources []string
}

type DeploymentExecutor struct {
	TfCmd TerraformCommand
}

func NewDeploymentExecutor() DeploymentExecutor {
	return DeploymentExecutor{TfCmd: terraformCmd{}}
}

func (d DeploymentExecutor) Deploy(ctx context.Context, c *resources.MultyConfig, prev resources.Resource, curr resources.Resource) (state *output.TfState, err error) {
	tmpDir := GetTempDirForUser(false, c.GetUserId())
	encoded, err := d.EncodeAndStoreTfFile(ctx, c, prev, curr, false)
	if err != nil {
		return nil, err
	}

	err = d.MaybeInit(ctx, c.GetUserId(), false)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	defer func() {
		log.Printf("[DEBUG] apply finished in %s", time.Since(start))
	}()

	defer func() {
		if flags.DryRun {
			return
		}
		// rollback if something goes wrong
		if err != nil {
			log.Println("[ERROR] Something went wrong, rolling back")
			originalC, err2 := c.GetOriginalConfig(types.Metadatas)
			if err2 != nil {
				log.Printf("[ERROR] Rollback unsuccessful: %s\n", err2)
				return
			}
			_, err2 = d.EncodeAndStoreTfFile(ctx, originalC, curr, prev, false)
			if err2 != nil {
				log.Printf("[ERROR] Rollback unsuccessful: %s\n", err2)
				return
			}

			err2 = d.TfCmd.Apply(ctx, tmpDir, encoded.affectedResources)
			if err2 != nil {
				log.Printf("[ERROR] Rollback unsuccessful: %s\n", err2)
				return
			}
		}
	}()

	err = d.TfCmd.Apply(ctx, tmpDir, encoded.affectedResources)
	if err != nil {
		return nil, err
	}

	state, err = d.GetState(ctx, c.GetUserId(), false)
	if err != nil {
		return state, errors.InternalServerErrorWithMessage("error parsing state", err)
	}

	return state, nil
}

func (d DeploymentExecutor) EncodeAndStoreTfFile(ctx context.Context, c *resources.MultyConfig, prev resources.Resource, curr resources.Resource, readonly bool) (EncodedResources, error) {
	credentials, err := util.ExtractCloudCredentials(ctx)
	if err != nil {
		return EncodedResources{}, err
	}
	encoded, err := EncodeTfFile(credentials, c, prev, curr)
	if err != nil {
		return encoded, err
	}

	tfBlock, err := GetTerraformBlock(c.GetUserId())
	if err != nil {
		return encoded, err
	}

	// TODO: move this to a proper place
	hclOutput := tfBlock + encoded.HclString

	tmpDir := GetTempDirForUser(readonly, c.GetUserId())
	err = os.MkdirAll(tmpDir, os.ModeDir|(os.ModePerm&0775))
	if err != nil {
		return EncodedResources{}, err
	}
	err = os.WriteFile(filepath.Join(tmpDir, tfFile), []byte(hclOutput), os.ModePerm&0664)
	return encoded, err
}

func (d DeploymentExecutor) MaybeInit(ctx context.Context, userId string, readonly bool) error {
	tmpDir := GetTempDirForUser(readonly, userId)
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

		err = d.TfCmd.Init(ctx, tmpDir)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}

func (d DeploymentExecutor) GetState(ctx context.Context, userId string, readonly bool) (*output.TfState, error) {
	tmpDir := GetTempDirForUser(readonly, userId)
	return d.TfCmd.GetState(ctx, tmpDir)
}

func (d DeploymentExecutor) RefreshState(ctx context.Context, userId string, c *resources.MultyConfig) error {
	_, err := d.EncodeAndStoreTfFile(ctx, c, nil, nil, true)
	if err != nil {
		return err
	}

	err = d.MaybeInit(ctx, userId, true)
	if err != nil {
		return err
	}

	return d.refresh(ctx, userId)
}

func (d DeploymentExecutor) refresh(ctx context.Context, userId string) error {
	start := time.Now()
	defer func() {
		log.Printf("[DEBUG] refresh finished in %s", time.Since(start))
	}()

	tmpDir := GetTempDirForUser(true, userId)
	return d.TfCmd.Refresh(ctx, tmpDir)
}

func GetTempDirForUser(readonly bool, userId string) string {
	tmpDir := filepath.Join(os.TempDir(), "multy", userId)

	if flags.Environment == flags.Local {
		tmpDir = filepath.Join(tmpDir, "local")
	}

	if readonly {
		tmpDir = filepath.Join(tmpDir, "readonly")
	}

	return tmpDir
}
