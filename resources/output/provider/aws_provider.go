package provider

const AwsResourceName = "aws"

type AwsProvider struct {
	ResourceName string `hcl:",key"`
	Region       string `hcl:"region"`
	Alias        string `hcl:"alias" hcle:"omitempty"`
}
