package provider

const AwsResourceName = "aws"

type AwsProvider struct {
	ResourceName string `hcl:",key"`
	Region       string `hcl:"region"`
	Alias        string `hcl:"alias" hcle:"omitempty"`
	AccessKey    string `hcl:"access_key" hcle:"omitempty"`
	SecretKey    string `hcl:"secret_key" hcle:"omitempty"`
}
