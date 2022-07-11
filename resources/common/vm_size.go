package common

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
)

var VMSIZE = map[commonpb.VmSize_Enum]map[commonpb.CloudProvider]string{
	commonpb.VmSize_GENERAL_NANO: {
		commonpb.CloudProvider_AWS:   "t2.nano",
		commonpb.CloudProvider_AZURE: "Standard_B1ls",
	},
	commonpb.VmSize_GENERAL_MICRO: {
		commonpb.CloudProvider_AWS:   "t2.micro",
		commonpb.CloudProvider_AZURE: "Standard_B1s",
		commonpb.CloudProvider_GCP:   "e2-micro",
	},
	commonpb.VmSize_GENERAL_SMALL: {
		commonpb.CloudProvider_AWS:   "t2.small",
		commonpb.CloudProvider_AZURE: "Standard_B1ms",
		commonpb.CloudProvider_GCP:   "e2-small",
	},
	commonpb.VmSize_GENERAL_MEDIUM: {
		commonpb.CloudProvider_AWS:   "t2.medium",
		commonpb.CloudProvider_AZURE: "Standard_B2s",
		commonpb.CloudProvider_GCP:   "e2-medium",
	},
	commonpb.VmSize_GENERAL_LARGE: {
		commonpb.CloudProvider_AWS:   "t2.large",
		commonpb.CloudProvider_AZURE: "Standard_B2ms",
		commonpb.CloudProvider_GCP:   "e2-standard-2",
	},
	commonpb.VmSize_GENERAL_XLARGE: {
		commonpb.CloudProvider_AWS:   "t2.xlarge",
		commonpb.CloudProvider_AZURE: "Standard_B4ms",
		commonpb.CloudProvider_GCP:   "e2-standard-4",
	},
	commonpb.VmSize_GENERAL_2XLARGE: {
		commonpb.CloudProvider_AWS:   "t2.2xlarge",
		commonpb.CloudProvider_AZURE: "Standard_B8ms",
		commonpb.CloudProvider_GCP:   "e2-standard-8",
	},
	commonpb.VmSize_COMPUTE_LARGE: {
		commonpb.CloudProvider_AWS:   "c4.large",
		commonpb.CloudProvider_AZURE: "Standard_F2s_v2",
	},
	commonpb.VmSize_COMPUTE_XLARGE: {
		commonpb.CloudProvider_AWS:   "c4.xlarge",
		commonpb.CloudProvider_AZURE: "Standard_F4s_v2",
	},
	commonpb.VmSize_COMPUTE_2XLARGE: {
		commonpb.CloudProvider_AWS:   "c4.2xlarge",
		commonpb.CloudProvider_AZURE: "Standard_F8s_v2",
	},
	commonpb.VmSize_COMPUTE_4XLARGE: {
		commonpb.CloudProvider_AWS:   "c4.4xlarge",
		commonpb.CloudProvider_AZURE: "Standard_F16s_v2",
	},
	commonpb.VmSize_COMPUTE_8XLARGE: {
		commonpb.CloudProvider_AWS:   "c4.8xlarge",
		commonpb.CloudProvider_AZURE: "Standard_F48s_v2",
	},
	commonpb.VmSize_MEMORY_LARGE: {
		commonpb.CloudProvider_AWS:   "r6g.large",
		commonpb.CloudProvider_AZURE: "Standard_E2s_v3",
	},
	commonpb.VmSize_MEMORY_XLARGE: {
		commonpb.CloudProvider_AWS:   "r6g.xlarge",
		commonpb.CloudProvider_AZURE: "Standard_E4s_v3",
	},
	commonpb.VmSize_MEMORY_2XLARGE: {
		commonpb.CloudProvider_AWS:   "r6g.2xlarge",
		commonpb.CloudProvider_AZURE: "Standard_E8s_v3",
	},
	commonpb.VmSize_MEMORY_4XLARGE: {
		commonpb.CloudProvider_AWS:   "r6g.4xlarge",
		commonpb.CloudProvider_AZURE: "Standard_E16s_v3",
	},
	commonpb.VmSize_MEMORY_8XLARGE: {
		commonpb.CloudProvider_AWS:   "r6g.8xlarge",
		commonpb.CloudProvider_AZURE: "Standard_E32s_v3",
	},
	commonpb.VmSize_MEMORY_12XLARGE: {
		commonpb.CloudProvider_AWS:   "r6g.12xlarge",
		commonpb.CloudProvider_AZURE: "Standard_E48s_v3",
	},
	commonpb.VmSize_MEMORY_16XLARGE: {
		commonpb.CloudProvider_AWS:   "r6g.16xlarge",
		commonpb.CloudProvider_AZURE: "Standard_E64a_v4",
	},
}

var DBSIZE = map[commonpb.DatabaseSize_Enum]map[commonpb.CloudProvider]string{
	commonpb.DatabaseSize_MICRO: {
		commonpb.CloudProvider_AWS:   "db.t2.micro",
		commonpb.CloudProvider_AZURE: "GP_Gen5_2",
		commonpb.CloudProvider_GCP:   "db-f1-micro",
	},
	commonpb.DatabaseSize_SMALL: {
		commonpb.CloudProvider_AWS:   "db.t2.small",
		commonpb.CloudProvider_AZURE: "GP_Gen5_4",
		commonpb.CloudProvider_GCP:   "db-g1-small",
	},
	commonpb.DatabaseSize_MEDIUM: {
		commonpb.CloudProvider_AWS:   "db.t2.medium",
		commonpb.CloudProvider_AZURE: "GP_Gen5_6",
		commonpb.CloudProvider_GCP:   "db-n1-standard-1",
	},
}

func GetVmSize(s commonpb.VmSize_Enum, c commonpb.CloudProvider) (string, error) {
	if size, ok := VMSIZE[s]; !ok {
		return "", fmt.Errorf("size %s is not available yet", s.String())
	} else if result, ok := size[c]; !ok {
		return "", fmt.Errorf("size %s is not available in cloud %s", s.String(), c.String())
	} else {
		return result, nil
	}
}
