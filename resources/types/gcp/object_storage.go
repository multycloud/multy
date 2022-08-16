package gcp_resources

import (
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/flags"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/object_storage"
	"github.com/multycloud/multy/resources/types"
)

type GcpObjectStorage struct {
	*types.ObjectStorage
}

func InitObjectStorage(o *types.ObjectStorage) resources.ResourceTranslator[*resourcespb.ObjectStorageResource] {
	return GcpObjectStorage{o}
}

func (r GcpObjectStorage) FromState(state *output.TfState) (*resourcespb.ObjectStorageResource, error) {
	out := &resourcespb.ObjectStorageResource{
		CommonParameters: &commonpb.CommonResourceParameters{
			ResourceId:      r.ResourceId,
			ResourceGroupId: r.Args.CommonParameters.ResourceGroupId,
			Location:        r.Args.CommonParameters.Location,
			CloudProvider:   r.Args.CommonParameters.CloudProvider,
			NeedsUpdate:     false,
		},
		Name:        r.Args.Name,
		Versioning:  r.Args.Versioning,
		GcpOverride: r.Args.GcpOverride,
	}

	if flags.DryRun {
		return out, nil
	}

	statuses := map[string]commonpb.ResourceStatus_Status{}

	if stateResource, exists, err := output.MaybeGetParsedById[object_storage.GoogleStorageBucket](state, r.ResourceId); exists {
		if err != nil {
			return nil, err
		}

		out.GcpOutputs = &resourcespb.ObjectStorageGcpOutputs{StorageBucketId: stateResource.SelfLink}
		out.Name = stateResource.Name
		out.Versioning = len(stateResource.Versioning) > 0 && stateResource.Versioning[0].Enabled
	} else {
		statuses["gcp_storage_bucket"] = commonpb.ResourceStatus_NEEDS_CREATE
	}

	if len(statuses) > 0 {
		out.CommonParameters.ResourceStatus = &commonpb.ResourceStatus{Statuses: statuses}
	}
	return out, nil
}

func (r GcpObjectStorage) Translate(resources.MultyContext) ([]output.TfBlock, error) {
	o := &object_storage.GoogleStorageBucket{
		GcpResource:              common.NewGcpResource(r.ResourceId, r.Args.Name, r.Args.GetGcpOverride().GetProject()),
		UniformBucketLevelAccess: false,
		Location:                 r.GetCloudSpecificLocation(),
		// this is needed, otherwise buckets with versioning can't be deleted normally
		ForceDestroy: true,
	}
	if r.Args.Versioning {
		o.Versioning = []object_storage.GoogleStorageBucketVersioning{{Enabled: true}}
	}
	return []output.TfBlock{o}, nil
}

func (r GcpObjectStorage) GetMainResourceName() (string, error) {
	return output.GetResourceName(object_storage.GoogleStorageBucket{}), nil
}
