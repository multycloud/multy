package aws_client

import (
	"log"
	"os"
	"path"
)

type LocalClient struct {
}

func newLocalClient() (*LocalClient, error) {
	return &LocalClient{}, nil
}

func (c LocalClient) SaveFile(configPrefix string, fileName string, content string) error {
	filePath, err := c.getFilePath(configPrefix, fileName)
	if err != nil {
		return err
	}
	err = os.WriteFile(filePath, []byte(content), os.ModePerm&0664)
	if err != nil {
		log.Printf("[ERROR] error saving file locally: %s\n", err.Error())
		return err
	}
	return nil
}

func (c LocalClient) ReadFile(configPrefix string, fileName string) (string, error) {
	filePath, err := c.getFilePath(configPrefix, fileName)
	if err != nil {
		return "", err
	}
	file, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		log.Printf("[ERROR] error reading file locally:  %s\n", err.Error())
		return "", err
	}
	return string(file), nil
}

func (c LocalClient) getFilePath(configPrefix string, fileName string) (string, error) {
	tmpDir := path.Join(os.TempDir(), "multy", configPrefix, "local")
	err := os.MkdirAll(tmpDir, os.ModeDir|(os.ModePerm&0775))
	if err != nil {
		return "", err
	}
	return path.Join(tmpDir, fileName), nil
}

func (c LocalClient) UpdateQPSMetric(_ string, service string, method string) error {
	err := logAction("local#", service, method)
	if err != nil {
		log.Printf("[WARNING] Logging error ocurred: %s", err)
	}
	return err
}

func (c LocalClient) UpdateErrorMetric(_ string, _ string, _ string) error {
	return nil
}
