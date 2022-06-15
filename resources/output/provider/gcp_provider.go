package provider

const GcpResourceName = "google"

type GcpProvider struct {
	ResourceName string `hcl:",key"`
	Region       string `hcl:"region"`
	Alias        string `hcl:"alias"`
	Credentials  string `hcl:"credentials" hcle:"omitempty"`
	Project      string `hcl:"project" hcle:"omitempty"`
}
