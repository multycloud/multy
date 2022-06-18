package types

import (
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/validate"
)

type Vault struct {
	resources.ResourceWithId[*resourcespb.VaultArgs]
}

func (r *Vault) Create(resourceId string, args *resourcespb.VaultArgs, others *resources.Resources) error {
	return CreateVault(r, resourceId, args, others)
}

func (r *Vault) Update(args *resourcespb.VaultArgs, others *resources.Resources) error {
	return NewVault(r, r.ResourceId, args, others)
}

func (r *Vault) Import(resourceId string, args *resourcespb.VaultArgs, others *resources.Resources) error {
	return NewVault(r, resourceId, args, others)
}

func (r *Vault) Export(_ *resources.Resources) (*resourcespb.VaultArgs, bool, error) {
	return r.Args, true, nil
}

func CreateVault(r *Vault, resourceId string, args *resourcespb.VaultArgs, others *resources.Resources) error {
	if args.CommonParameters.ResourceGroupId == "" {
		rgId, err := NewRg("kv", others, args.GetCommonParameters().GetLocation(), args.GetCommonParameters().GetCloudProvider())
		if err != nil {
			return err
		}
		args.CommonParameters.ResourceGroupId = rgId
	}

	return NewVault(r, resourceId, args, others)
}

func NewVault(r *Vault, resourceId string, args *resourcespb.VaultArgs, _ *resources.Resources) error {
	r.ResourceWithId = resources.NewResource(resourceId, args)
	return nil
}

func (r *Vault) Validate(ctx resources.MultyContext) (errs []validate.ValidationError) {
	errs = append(errs, r.ResourceWithId.Validate()...)
	return errs
}
