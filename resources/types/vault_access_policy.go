package types

import (
	"fmt"
	"github.com/multycloud/multy/api/errors"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/validate"
)

type VaultAccessPolicy struct {
	resources.ChildResourceWithId[*Vault, *resourcespb.VaultAccessPolicyArgs]
	Vault *Vault
}

func (r *VaultAccessPolicy) Create(resourceId string, args *resourcespb.VaultAccessPolicyArgs, others *resources.Resources) error {
	return NewVaultAccessPolicy(r, resourceId, args, others)
}

func (r *VaultAccessPolicy) Update(args *resourcespb.VaultAccessPolicyArgs, others *resources.Resources) error {
	return NewVaultAccessPolicy(r, r.ResourceId, args, others)
}

func (r *VaultAccessPolicy) Import(resourceId string, args *resourcespb.VaultAccessPolicyArgs, others *resources.Resources) error {
	return NewVaultAccessPolicy(r, resourceId, args, others)
}

func (r *VaultAccessPolicy) Export(_ *resources.Resources) (*resourcespb.VaultAccessPolicyArgs, bool, error) {
	return r.Args, true, nil
}

func NewVaultAccessPolicy(vap *VaultAccessPolicy, resourceId string, args *resourcespb.VaultAccessPolicyArgs, others *resources.Resources) error {
	v, err := resources.Get[*Vault](resourceId, others, args.VaultId)
	if err != nil {
		return errors.ValidationError(resources.NewError(err, vap.ResourceId, "vault_id"))
	}
	vap.ChildResourceWithId = resources.NewChildResource(resourceId, v, args)
	vap.Parent = v
	vap.Vault = v

	// This dependency is not modelled directly by the id, so we have to do it this way
	for _, res := range others.ResourceMap {
		if vm, ok := res.(*VirtualMachine); ok && vm.GetAwsIdentity() == args.Identity {
			others.AddDependency(resourceId, vm.GetResourceId())
		}
	}

	return nil
}

func (r *VaultAccessPolicy) Validate(ctx resources.MultyContext) (errs []validate.ValidationError) {
	if r.Args.Access == resourcespb.VaultAccess_UNKNOWN {
		errs = append(errs, r.NewValidationError(fmt.Errorf("unknown vault access"), "access"))
	}
	return errs
}
