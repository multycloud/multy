package common

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"strings"
)

func GetResourceId(prefix string, cloud commonpb.CloudProvider) string {
	if cloud == commonpb.CloudProvider_UNKNOWN_PROVIDER {
		return prefix
	}
	return fmt.Sprintf("%s_%s", prefix, strings.ToLower(cloud.String()))
}

func ParseCloudFromResourceId(resourceId string) commonpb.CloudProvider {
	split := strings.Split(resourceId, "_")
	cloud := strings.ToUpper(split[len(split)-1])
	if cloudValue, ok := commonpb.CloudProvider_value[cloud]; ok {
		return commonpb.CloudProvider(cloudValue)
	}
	return commonpb.CloudProvider_UNKNOWN_PROVIDER
}
