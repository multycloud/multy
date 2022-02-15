package output

type isTfBlock struct {
}

type TerraformBlock interface {
	isTerraformBlock()
	GetR() any
}

func (*isTfBlock) isTerraformBlock() {
}

// ResourceWrapper just to add a resource {} around when encoding into hcl
type ResourceWrapper struct {
	*isTfBlock `hcle:"omit"`
	R          any `hcl:"resource"`
}

type DataSourceWrapper struct {
	*isTfBlock `hcle:"omit"`
	R          any `hcl:"data"`
}

func (r DataSourceWrapper) GetR() any {
	return r.R
}

func (r ResourceWrapper) GetR() any {
	return r.R
}

func IsTerraformBlock(r any) bool {
	if _, ok := r.(TerraformBlock); ok {
		return true
	}
	return false
}
