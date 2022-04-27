package terraform

import (
	"fmt"
	"github.com/multycloud/multy/resources/output"
)

const LocalFileResourceName = "local_file"

type LocalFile struct {
	*output.TerraformResource `hcl:",squash" default:"name=local_file"`
	ContentBase64             string `hcl:"content_base64"`
	Filename                  string `hcl:"filename"`
}

func NewLocalFile(resourceId string, contentBase64 string) LocalFile {
	return LocalFile{
		TerraformResource: &output.TerraformResource{
			ResourceName: LocalFileResourceName,
			ResourceId:   resourceId,
		},
		ContentBase64: contentBase64,
		Filename:      fmt.Sprintf("./multy/local/%s", resourceId),
	}
}

func (l *LocalFile) GetFilename() string {
	return fmt.Sprintf("%s.%s.filename", LocalFileResourceName, l.ResourceId)
}
