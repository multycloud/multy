//go:build e2e
// +build e2e

package e2e

import (
	"encoding/base64"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources/common"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func testObjectStorageObject(t *testing.T, cloud commonpb.CloudProvider) {
	ctx := getCtx(t, cloud, "obj")

	objectStorageName := "e2eteststorage" + common.RandomString(4)

	createObjStorageRequest := &resourcespb.CreateObjectStorageRequest{Resource: &resourcespb.ObjectStorageArgs{
		CommonParameters: &commonpb.ResourceCommonArgs{
			Location:      commonpb.Location_US_WEST_1,
			CloudProvider: cloud,
		},
		Name:       objectStorageName,
		Versioning: false,
	}}
	storage, err := server.ObjectStorageService.Create(ctx, createObjStorageRequest)
	if err != nil {
		logGrpcErrorDetails(t, err)
		t.Fatalf("unable to create storage: %+v", err)
	}
	cleanup(t, ctx, server.ObjectStorageService, storage)

	read, err := server.ObjectStorageService.Read(ctx, &resourcespb.ReadObjectStorageRequest{ResourceId: storage.CommonParameters.ResourceId})
	if err != nil {
		t.Fatalf("unable to read obj, %s", err)
	}

	assert.Nil(t, read.GetCommonParameters().GetResourceStatus())
	assert.Equal(t, read.GetName(), createObjStorageRequest.GetResource().GetName())
	assert.Equal(t, read.GetVersioning(), createObjStorageRequest.GetResource().GetVersioning())

	createObjStorageObjRequest := &resourcespb.CreateObjectStorageObjectRequest{Resource: &resourcespb.ObjectStorageObjectArgs{
		Name:            "public-text.html",
		Acl:             resourcespb.ObjectStorageObjectAcl_PUBLIC_READ,
		ObjectStorageId: storage.CommonParameters.ResourceId,
		ContentBase64:   base64.StdEncoding.EncodeToString([]byte(`<h1>hello world</h1>`)),
		ContentType:     "text/html",
	}}
	obj, err := server.ObjectStorageObjectService.Create(ctx, createObjStorageObjRequest)
	if err != nil {
		logGrpcErrorDetails(t, err)
		t.Fatalf("unable to create vn: %+v", err)
	}
	cleanup(t, ctx, server.ObjectStorageObjectService, obj)

	readObj, err := server.ObjectStorageObjectService.Read(ctx, &resourcespb.ReadObjectStorageObjectRequest{ResourceId: obj.CommonParameters.ResourceId})
	if err != nil {
		t.Fatalf("unable to read obj, %s", err)
	}

	assert.Nil(t, readObj.GetCommonParameters().GetResourceStatus())
	assert.Equal(t, readObj.GetName(), createObjStorageObjRequest.GetResource().GetName())
	assert.Equal(t, readObj.GetAcl(), createObjStorageObjRequest.GetResource().GetAcl())
	assert.Equal(t, readObj.GetContentBase64(), createObjStorageObjRequest.GetResource().GetContentBase64())

	resp, err := http.Get(obj.GetUrl())
	if err != nil {
		t.Fatalf("unable to do GET request on %s: %s", obj.GetUrl(), err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("unable to read body: %s", err)
	}

	assert.Equal(t, "<h1>hello world</h1>", string(body))
}

func TestAwsObjectStorageObject(t *testing.T) {
	t.Parallel()
	testObjectStorageObject(t, commonpb.CloudProvider_AWS)
}
func TestAzureObjectStorageObject(t *testing.T) {
	t.Parallel()
	testObjectStorageObject(t, commonpb.CloudProvider_AZURE)
}
func TestGcpObjectStorageObject(t *testing.T) {
	t.Parallel()
	testObjectStorageObject(t, commonpb.CloudProvider_GCP)
}
