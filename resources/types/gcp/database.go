package gcp_resources

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/flags"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/database"
	"github.com/multycloud/multy/resources/types"
	"strings"
)

type GcpDatabase struct {
	*types.Database
}

func InitDatabase(r *types.Database) resources.ResourceTranslator[*resourcespb.DatabaseResource] {
	return GcpDatabase{r}
}

func (r GcpDatabase) FromState(state *output.TfState, plan *output.TfPlan) (*resourcespb.DatabaseResource, error) {
	out := &resourcespb.DatabaseResource{
		CommonParameters: &commonpb.CommonResourceParameters{
			ResourceId:      r.ResourceId,
			ResourceGroupId: r.Args.CommonParameters.ResourceGroupId,
			Location:        r.Args.CommonParameters.Location,
			CloudProvider:   r.Args.CommonParameters.CloudProvider,
			NeedsUpdate:     false,
		},
		Name:               r.Args.Name,
		Engine:             r.Args.Engine,
		EngineVersion:      r.Args.EngineVersion,
		StorageGb:          r.Args.StorageGb,
		Size:               r.Args.Size,
		Username:           r.Args.Username,
		Password:           r.Args.Password,
		SubnetIds:          r.Args.SubnetIds,
		Port:               r.Args.Port,
		SubnetId:           r.Args.SubnetId,
		Host:               "",
		ConnectionUsername: r.Args.Username,
		GcpOverride:        r.Args.GcpOverride,
	}

	if flags.DryRun {
		return out, nil
	}

	statuses := map[string]commonpb.ResourceStatus_Status{}
	out.GcpOutputs = &resourcespb.DatabaseGcpOutputs{}

	if db, exists, err := output.MaybeGetParsedById[database.GoogleSqlDatabaseInstance](state, r.ResourceId); exists {
		if err != nil {
			return nil, err
		}
		out.Host = db.PublicIpAddress
		out.GcpOutputs.SqlDatabaseInstanceId = db.SelfLink
		out.Name = db.Name
		out.StorageGb = int64(db.Settings[0].DiskSize)
		engine, version, versionErr := parseDbVersion(db.DatabaseVersion)
		if versionErr != nil {
			statuses["gcp_sql_database_instance"] = commonpb.ResourceStatus_NEEDS_RECREATE
		}
		out.Engine = engine
		out.EngineVersion = version

	} else {
		statuses["gcp_sql_database_instance"] = commonpb.ResourceStatus_NEEDS_CREATE
	}

	if dbUser, exists, err := output.MaybeGetParsedById[database.GoogleSqlUser](state, r.ResourceId); exists {
		if err != nil {
			return nil, err
		}
		out.Username = dbUser.Name
		out.Password = dbUser.Password
	} else {
		statuses["gcp_sql_user"] = commonpb.ResourceStatus_NEEDS_CREATE
	}

	if len(statuses) > 0 {
		out.CommonParameters.ResourceStatus = &commonpb.ResourceStatus{Statuses: statuses}
	}
	return out, nil
}

func (r GcpDatabase) getDbVersion() (string, error) {
	engineMap := map[resourcespb.DatabaseEngine]string{
		resourcespb.DatabaseEngine_MYSQL:    "MYSQL",
		resourcespb.DatabaseEngine_POSTGRES: "POSTGRES",
	}

	if _, ok := engineMap[r.Args.Engine]; !ok {
		return "", fmt.Errorf("database engine %s not available in GCP", r.Args.Engine.String())
	}

	version := strings.Replace(r.Args.EngineVersion, ".", "_", 1)
	return fmt.Sprintf("%s_%s", engineMap[r.Args.Engine], version), nil
}

func parseDbVersion(dbVersion string) (resourcespb.DatabaseEngine, string, error) {
	engine := resourcespb.DatabaseEngine_UNKNOWN_ENGINE
	version := ""
	engineAndVersion := strings.SplitN(dbVersion, "_", 2)
	if len(engineAndVersion) != 2 {
		return resourcespb.DatabaseEngine_UNKNOWN_ENGINE, "", fmt.Errorf("unknown version format")
	}
	if e, ok := resourcespb.DatabaseEngine_value[engineAndVersion[0]]; ok {
		engine = resourcespb.DatabaseEngine(e)
	}

	version = strings.Replace(engineAndVersion[1], "_", ".", 1)
	return engine, version, nil
}

func (r GcpDatabase) Translate(_ resources.MultyContext) ([]output.TfBlock, error) {
	dbVersion, err := r.getDbVersion()
	if err != nil {
		return nil, err
	}

	db := &database.GoogleSqlDatabaseInstance{
		GcpResource:        common.NewGcpResource(r.ResourceId, r.Args.Name, r.Args.GetGcpOverride().GetProject()),
		DatabaseVersion:    dbVersion,
		DeletionProtection: false,
		Settings: []database.GoogleSqlDatabaseInstanceSettings{{
			Tier:             common.DBSIZE[r.Args.Size][commonpb.CloudProvider_GCP],
			AvailabilityType: "ZONAL",
			DiskAutoResize:   false,
			DiskSize:         int(r.Args.StorageGb),
			IpConfiguration: database.GoogleIpConfiguration{
				AuthorizedNetworks: []database.GoogleAuthorizedNetwork{
					{
						Value: "0.0.0.0/0",
					},
				},
			},
		}},
	}

	instance := fmt.Sprintf("%s.%s.name", output.GetResourceName(*db), r.ResourceId)
	user := &database.GoogleSqlUser{
		GcpResource: common.NewGcpResource(r.ResourceId, r.Args.Username, r.Args.GetGcpOverride().GetProject()),
		Instance:    instance,
		Password:    r.Args.Password,
	}

	return []output.TfBlock{db, user}, nil

}

func (r GcpDatabase) GetMainResourceName() (string, error) {
	return output.GetResourceName(database.GoogleSqlDatabaseInstance{}), nil
}
