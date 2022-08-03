package deploy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/multycloud/multy/api/errors"
	"github.com/multycloud/multy/db"
	"github.com/multycloud/multy/flags"
	"github.com/multycloud/multy/resources/output"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime/trace"
)

type TerraformCommand interface {
	Init(ctx context.Context, dir string) error
	Apply(ctx context.Context, dir string, resources []string) error
	Refresh(ctx context.Context, dir string) error
	GetState(ctx context.Context, userId string, dir db.TfStateReader) (*output.TfState, error)
}

type terraformCmd struct {
}

type tfOutput struct {
	Level      string `json:"@level"`
	Message    string `json:"@message"`
	Diagnostic struct {
		Summary string `json:"summary"`
		Detail  string `json:"detail"`
		Address string `json:"address"`
	} `json:"diagnostic"`
}

func (c terraformCmd) Apply(ctx context.Context, dir string, resources []string) error {
	if len(resources) == 0 {
		return nil
	}

	region := trace.StartRegion(ctx, "tf apply")
	defer region.End()

	var targetArgs []string

	idsToPrint := ""
	for _, id := range resources {
		idsToPrint += id + ", "
		targetArgs = append(targetArgs, "-target="+id)
	}
	log.Printf("[INFO] Running apply for targets: %s", idsToPrint)
	cmd := exec.CommandContext(ctx, "terraform", append([]string{"-chdir=" + dir, "apply", "-refresh=false", "-auto-approve", "--json"}, targetArgs...)...)
	if flags.DryRun {
		cmd = exec.CommandContext(ctx, "terraform", append([]string{"-chdir=" + dir, "plan", "-refresh=false", "--json"}, targetArgs...)...)
	}
	outputJson := new(bytes.Buffer)
	cmd.Stdout = outputJson
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		outputs, parseErr := parseTfOutputs(outputJson)
		if parseErr != nil {
			return errors.InternalServerErrorWithMessage("error deploying resources", parseErr)
		}
		if parseErr := getFirstError(outputs); parseErr != nil {
			return errors.DeployError(parseErr)
		}
		return errors.InternalServerErrorWithMessage("error deploying resources", err)
	}

	return nil
}

func (c terraformCmd) Init(ctx context.Context, dir string) error {
	region := trace.StartRegion(ctx, "tf init")
	defer region.End()

	cmd := exec.CommandContext(ctx, "terraform", "-chdir="+dir, "init", "-reconfigure", "-lock-timeout", "1m")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("unable to initialize terraform")
	}
	return nil
}

func (c terraformCmd) Refresh(ctx context.Context, dir string) error {
	region := trace.StartRegion(ctx, "tf refresh")
	defer region.End()

	outputJson := new(bytes.Buffer)
	cmd := exec.CommandContext(ctx, "terraform", "-chdir="+dir, "refresh", "-json")
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

func (c terraformCmd) GetState(ctx context.Context, userId string, client db.TfStateReader) (*output.TfState, error) {
	region := trace.StartRegion(ctx, "tf show")
	defer region.End()

	terraformState, err := client.LoadTerraformState(ctx, userId)
	if err != nil {
		return nil, err
	}

	state := output.TfState{}

	err = json.Unmarshal([]byte(terraformState), &state)
	if err != nil {
		return nil, err
	}
	return &state, err
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
