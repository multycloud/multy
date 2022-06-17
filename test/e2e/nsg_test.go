//go:build e2e
// +build e2e

package e2e

import (
	"context"
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ssh"
	"testing"
	"time"
)

func testNsgRules(t *testing.T, ctx context.Context, nsg *resourcespb.NetworkSecurityGroupResource, ip string, config *ssh.ClientConfig) {
	updateRequest := &resourcespb.UpdateNetworkSecurityGroupRequest{
		ResourceId: nsg.CommonParameters.ResourceId,
		Resource: &resourcespb.NetworkSecurityGroupArgs{
			CommonParameters: &commonpb.ResourceCommonArgs{
				Location:        nsg.CommonParameters.Location,
				CloudProvider:   nsg.CommonParameters.CloudProvider,
				ResourceGroupId: nsg.CommonParameters.ResourceGroupId,
			},
			Name:             nsg.Name,
			VirtualNetworkId: nsg.VirtualNetworkId,
			Rules: []*resourcespb.NetworkSecurityRule{{
				Protocol: "tcp",
				Priority: 110,
				PortRange: &resourcespb.PortRange{
					From: 443,
					To:   443,
				},
				CidrBlock: "0.0.0.0/0",
				Direction: resourcespb.Direction_EGRESS,
			}, {
				Protocol: "tcp",
				Priority: 120,
				PortRange: &resourcespb.PortRange{
					From: 80,
					To:   80,
				},
				CidrBlock: "0.0.0.0/0",
				Direction: resourcespb.Direction_EGRESS,
			}},
		},
	}
	_, err := server.NetworkSecurityGroupService.Update(ctx, updateRequest)
	if err != nil {
		logGrpcErrorDetails(t, err)
		t.Fatal(fmt.Errorf("error updating nsg: %+v", err))
	}

	t.Cleanup(func() {
		updateRequest.Resource.Rules = nsg.Rules
		_, err := server.NetworkSecurityGroupService.Update(ctx, updateRequest)
		if err != nil {
			logGrpcErrorDetails(t, err)
			t.Logf("unable to update resource %+v", err)
		}
	})

	time.Sleep(1 * time.Minute)

	conn, err := ssh.Dial("tcp", ip+":22", config)
	assert.Error(t, err)
	if err == nil {
		t.Cleanup(func() {
			conn.Close()
		})
	}
}
