package flags

type Env string

const (
	Dev   Env = "dev"
	Prod  Env = "prod"
	Local Env = "local"
)

var DryRun = false
var Environment = Prod
