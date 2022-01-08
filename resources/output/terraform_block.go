package output

type isTfBlock struct {
}

type TerraformBlock interface {
	isTerraformBlock()
}

func (*isTfBlock) isTerraformBlock() {
}

// ResourceWrapper just to add a resource {} around when encoding into hcl
type ResourceWrapper struct {
	*isTfBlock `hcle:"omit"`
	R          interface{} `hcl:"resource"`
}

type DataSourceWrapper struct {
	*isTfBlock `hcle:"omit"`
	R          interface{} `hcl:"data"`
}

func IsTerraformBlock(r interface{}) bool {
	if _, ok := r.(TerraformBlock); ok {
		return true
	}
	return false
}
