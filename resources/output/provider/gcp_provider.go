package provider

const GcpResourceName = "google"

type GcpProvider struct {
	ResourceName string `hcl:",key"`
	Region       string `hcl:"region"`
	Credentials  string `hcl:"credentials" hcle:"omitempty"`
	ProjectId    string `hcl:"project_id" hcle:"omitempty"`
}
