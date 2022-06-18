package types

import (
	"github.com/multycloud/multy/api/errors"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/validate"
)

type VaultSecret struct {
	resources.ChildResourceWithId[*Vault, *resourcespb.VaultSecretArgs]

	Vault *Vault
}

func (r *VaultSecret) Create(resourceId string, args *resourcespb.VaultSecretArgs, others *resources.Resources) error {
	return NewVaultSecret(r, resourceId, args, others)
}

func (r *VaultSecret) Update(args *resourcespb.VaultSecretArgs, others *resources.Resources) error {
	return NewVaultSecret(r, r.ResourceId, args, others)
}

func (r *VaultSecret) Import(resourceId string, args *resourcespb.VaultSecretArgs, others *resources.Resources) error {
	return NewVaultSecret(r, resourceId, args, others)
}

func (r *VaultSecret) Export(_ *resources.Resources) (*resourcespb.VaultSecretArgs, bool, error) {
	return r.Args, true, nil
}

func NewVaultSecret(vs *VaultSecret, resourceId string, args *resourcespb.VaultSecretArgs, others *resources.Resources) error {
	v, err := resources.Get[*Vault](resourceId, others, args.VaultId)
	if err != nil {
		return errors.ValidationError(resources.NewError(err, vs.ResourceId, "vault_id"))
	}
	vs.ChildResourceWithId = resources.NewChildResource(resourceId, v, args)
	vs.Parent = v
	vs.Vault = v
	return nil
}

func (r *VaultSecret) Validate(ctx resources.MultyContext) (errs []validate.ValidationError) {
	return errs
}
