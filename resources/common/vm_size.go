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
	commonpb.VmSize_LARGE: {
		commonpb.CloudProvider_AWS:   "t2.large",
		commonpb.CloudProvider_AZURE: "Standard_D2as_v4",
	},
	commonpb.VmSize_GENERAL_NANO: {
		commonpb.CloudProvider_AWS:   "t2.nano",
		commonpb.CloudProvider_AZURE: "Standard_B1ls",
	},
	commonpb.VmSize_GENERAL_LARGE: {
		commonpb.CloudProvider_AWS:   "t2.large",
		commonpb.CloudProvider_AZURE: "Standard_B2ms",
	},
	commonpb.VmSize_GENERAL_XLARGE: {
		commonpb.CloudProvider_AWS:   "t2.xlarge",
		commonpb.CloudProvider_AZURE: "Standard_B4ms",
	},
	commonpb.VmSize_GENERAL_2XLARGE: {
		commonpb.CloudProvider_AWS:   "t2.2xlarge",
		commonpb.CloudProvider_AZURE: "Standard_B8ms",
	},
	commonpb.VmSize_COMPUTE_LARGE: {
		commonpb.CloudProvider_AWS:   "c7g.large",
		commonpb.CloudProvider_AZURE: "Standard_F2s_v2",
	},
	commonpb.VmSize_COMPUTE_XLARGE: {
		commonpb.CloudProvider_AWS:   "c7g.xlarge",
		commonpb.CloudProvider_AZURE: "Standard_F4s_v2",
	},
	commonpb.VmSize_COMPUTE_2XLARGE: {
		commonpb.CloudProvider_AWS:   "c7g.2xlarge",
		commonpb.CloudProvider_AZURE: "Standard_F8s_v2",
	},
	commonpb.VmSize_COMPUTE_4XLARGE: {
		commonpb.CloudProvider_AWS:   "c7g.4xlarge",
		commonpb.CloudProvider_AZURE: "Standard_F16s_v2",
	},
	commonpb.VmSize_COMPUTE_8XLARGE: {
		commonpb.CloudProvider_AWS:   "c7g.8xlarge",
		commonpb.CloudProvider_AZURE: "Standard_F32s_v2",
	},
	commonpb.VmSize_COMPUTE_12XLARGE: {
		commonpb.CloudProvider_AWS:   "c7g.12xlarge",
		commonpb.CloudProvider_AZURE: "Standard_F48s_v2",
	},
	commonpb.VmSize_COMPUTE_16XLARGE: {
		commonpb.CloudProvider_AWS:   "c7g.16xlarge",
		commonpb.CloudProvider_AZURE: "Standard_F64s_v2",
	},
}

const (
	// MICRO - 1 core and 0.5 gb ram
	MICRO = "nano"
	// MEDIUM - 2 cores and 4 gb ram
	MEDIUM = "medium"
	// LARGE - 2 cores and 8 gb ram
	LARGE = "large"
)

var DBSIZE = map[commonpb.DatabaseSize_Enum]map[commonpb.CloudProvider]string{
	commonpb.DatabaseSize_MICRO: {
		commonpb.CloudProvider_AWS:   "db.t2.micro",
		commonpb.CloudProvider_AZURE: "GP_Gen5_2",
	},
	commonpb.DatabaseSize_SMALL: {
		commonpb.CloudProvider_AWS:   "db.t2.small",
		commonpb.CloudProvider_AZURE: "GP_Gen5_4",
	},
	commonpb.DatabaseSize_MEDIUM: {
		commonpb.CloudProvider_AWS:   "db.t2.medium",
		commonpb.CloudProvider_AZURE: "GP_Gen5_6",
	},
}
