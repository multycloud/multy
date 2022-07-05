package kubernetes_service

type KubeConfig struct {
	ApiVersion     string                   `yaml:"apiVersion"`
	Clusters       []NamedKubeConfigCluster `yaml:"clusters"`
	Contexts       []NamedKubeConfigContext `yaml:"contexts"`
	CurrentContext string                   `yaml:"current-context"`
	Users          []KubeConfigUser         `yaml:"users"`
	Kind           string                   `yaml:"kind"`
}

type NamedKubeConfigCluster struct {
	Name    string            `yaml:"name"`
	Cluster KubeConfigCluster `yaml:"cluster"`
}

type KubeConfigCluster struct {
	CertificateAuthorityData string `yaml:"certificate-authority-data"`
	Server                   string `yaml:"server"`
}

type NamedKubeConfigContext struct {
	Name    string            `yaml:"name"`
	Context KubeConfigContext `yaml:"context"`
}

type KubeConfigContext struct {
	User    string `yaml:"user"`
	Cluster string `yaml:"cluster"`
}

type KubeConfigUser struct {
	Name string `yaml:"name"`
	User struct {
		Exec KubeConfigExec `yaml:"exec"`
	} `yaml:"user"`
}

type KubeConfigExec struct {
	ApiVersion         string   `yaml:"apiVersion"`
	Command            string   `yaml:"command"`
	Args               []string `yaml:"args,omitempty"`
	InteractiveMode    string   `yaml:"interactiveMode"`
	ProvideClusterInfo bool     `yaml:"provideClusterInfo,omitempty"`
	InstallHint        string   `yaml:"installHint,omitempty"`
}
