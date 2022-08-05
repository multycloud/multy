package services

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/types/metadata"
	"google.golang.org/protobuf/proto"
	"strings"
)

func GetConfigPrefixForCloud(userId string, cloud commonpb.CloudProvider) string {
	return fmt.Sprintf("%s/%s", userId, strings.ToLower(cloud.String()))
}

func GetConfigPrefix(req WithResourceId, userId string) string {
	cloud := common.ParseCloudFromResourceId(req.GetResourceId())
	if cloud == commonpb.CloudProvider_UNKNOWN_PROVIDER {
		return userId
	}

	return GetConfigPrefixForCloud(userId, cloud)
}

func getConfigPrefixForCreateReq(r proto.Message, userId string) string {
	converter, err := resources.ResourceMetadatas(metadata.Metadatas).GetConverter(proto.MessageName(r))
	if err != nil {
		return ""
	}
	cloud := converter.ParseCloud(r)
	if cloud == commonpb.CloudProvider_UNKNOWN_PROVIDER {
		return userId
	}

	return GetConfigPrefixForCloud(userId, cloud)
}
