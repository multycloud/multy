package common

import (
	"github.com/multycloud/multy/api/proto/commonpb"
)

var VMSIZE = map[commonpb.VmSize_Enum]map[commonpb.CloudProvider]string{
	commonpb.VmSize_MICRO: {
		commonpb.CloudProvider_AWS:   "t2.nano",
		commonpb.CloudProvider_AZURE: "Standard_B1ls",
	},
	commonpb.VmSize_MEDIUM: {
		commonpb.CloudProvider_AWS:   "t2.medium",
		commonpb.CloudProvider_AZURE: "Standard_A2_v2",
	},
}

const (
	// MICRO - 1 core and 0.5 gb ram
	MICRO = "nano"
	// MEDIUM - 2 cores and 4 gb ram
	MEDIUM = "medium"
)

var DBSIZE = map[commonpb.DatabaseSize_Enum]map[commonpb.CloudProvider]string{
	commonpb.DatabaseSize_MICRO: {
		commonpb.CloudProvider_AWS:   "db.t2.micro",
		commonpb.CloudProvider_AZURE: "GP_Gen5_2",
	},
}
