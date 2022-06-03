package aws_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/multycloud/multy/flags"
	"io/ioutil"
	"net/http"
)

func logAction(userId string, service string, method string) error {
	if flags.DryRun || flags.NoTelemetry {
		return nil
	}
	postBody, _ := json.Marshal(map[string]string{
		"action":  method,
		"service": service,
		"user_id": userId,
		"env":     string(flags.Environment),
	})
	resp, err := http.Post(logUrl, "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		data, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("log call returned status code %s with message %s", resp.Status, string(data))
	}
	return nil
}
