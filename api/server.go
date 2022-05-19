package api

import (
	"context"
	"fmt"
	aws_client "github.com/multycloud/multy/api/aws"
	"github.com/multycloud/multy/api/deploy"
	"github.com/multycloud/multy/api/errors"
	"github.com/multycloud/multy/api/proto"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/api/service_context"
	"github.com/multycloud/multy/api/services"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
	"github.com/multycloud/multy/flags"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"os"
	"strings"
)

type Server struct {
	proto.UnimplementedMultyResourceServiceServer
	*service_context.ServiceContext
	VnService                    services.Service[*resourcespb.VirtualNetworkArgs, *resourcespb.VirtualNetworkResource]
	SubnetService                services.Service[*resourcespb.SubnetArgs, *resourcespb.SubnetResource]
	NetworkInterfaceService      services.Service[*resourcespb.NetworkInterfaceArgs, *resourcespb.NetworkInterfaceResource]
	RouteTableService            services.Service[*resourcespb.RouteTableArgs, *resourcespb.RouteTableResource]
	RouteTableAssociationService services.Service[*resourcespb.RouteTableAssociationArgs, *resourcespb.RouteTableAssociationResource]
	NetworkSecurityGroupService  services.Service[*resourcespb.NetworkSecurityGroupArgs, *resourcespb.NetworkSecurityGroupResource]
	DatabaseService              services.Service[*resourcespb.DatabaseArgs, *resourcespb.DatabaseResource]
	ObjectStorageService         services.Service[*resourcespb.ObjectStorageArgs, *resourcespb.ObjectStorageResource]
	ObjectStorageObjectService   services.Service[*resourcespb.ObjectStorageObjectArgs, *resourcespb.ObjectStorageObjectResource]
	PublicIpService              services.Service[*resourcespb.PublicIpArgs, *resourcespb.PublicIpResource]
	KubernetesClusterService     services.Service[*resourcespb.KubernetesClusterArgs, *resourcespb.KubernetesClusterResource]
	KubernetesNodePoolService    services.Service[*resourcespb.KubernetesNodePoolArgs, *resourcespb.KubernetesNodePoolResource]
	LambdaService                services.Service[*resourcespb.LambdaArgs, *resourcespb.LambdaResource]
	VaultService                 services.Service[*resourcespb.VaultArgs, *resourcespb.VaultResource]
	VaultAccessPolicyService     services.Service[*resourcespb.VaultAccessPolicyArgs, *resourcespb.VaultAccessPolicyResource]
	VaultSecretService           services.Service[*resourcespb.VaultSecretArgs, *resourcespb.VaultSecretResource]
	VirtualMachineService        services.Service[*resourcespb.VirtualMachineArgs, *resourcespb.VirtualMachineResource]
}

func RunServer(ctx context.Context, port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var s *grpc.Server
	if flags.Environment == flags.Local {
		log.Println("[INFO] running in local mode")
		s = grpc.NewServer()
	} else {
		endpoint, exists := os.LookupEnv("MULTY_API_ENDPOINT")
		if !exists {
			log.Fatalf("api_endpoint env var is not set")
		}

		certFile := fmt.Sprintf("/etc/letsencrypt/live/%s/fullchain.pem", endpoint)
		keyFile := fmt.Sprintf("/etc/letsencrypt/live/%s/privkey.pem", endpoint)
		creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
		if err != nil {
			log.Fatalf("unable to read certificate (%s)", err.Error())
		}
		s = grpc.NewServer(grpc.Creds(creds))
	}

	go func() {
		<-ctx.Done()
		s.GracefulStop()
		_ = lis.Close()
	}()
	awsClient, err := aws_client.NewClient()
	if err != nil {
		log.Fatalf("failed to initialize aws client: %v", err)
	}
	database, err := db.NewDatabase(awsClient)
	if err != nil {
		log.Fatalf("failed to load db: %v", err)
	}
	defer database.Close()
	serviceContext := &service_context.ServiceContext{
		Database:  database,
		AwsClient: awsClient,
	}

	server := CreateServer(serviceContext)
	proto.RegisterMultyResourceServiceServer(s, &server)
	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func CreateServer(serviceContext *service_context.ServiceContext) Server {
	return Server{
		proto.UnimplementedMultyResourceServiceServer{},
		serviceContext,
		services.NewService[*resourcespb.VirtualNetworkArgs, *resourcespb.VirtualNetworkResource]("virtual_network", serviceContext),
		services.NewService[*resourcespb.SubnetArgs, *resourcespb.SubnetResource]("subnet", serviceContext),
		services.NewService[*resourcespb.NetworkInterfaceArgs, *resourcespb.NetworkInterfaceResource]("network_interface", serviceContext),
		services.NewService[*resourcespb.RouteTableArgs, *resourcespb.RouteTableResource]("route_table", serviceContext),
		services.NewService[*resourcespb.RouteTableAssociationArgs, *resourcespb.RouteTableAssociationResource]("route_table_association", serviceContext),
		services.NewService[*resourcespb.NetworkSecurityGroupArgs, *resourcespb.NetworkSecurityGroupResource]("network_security_group", serviceContext),
		services.NewService[*resourcespb.DatabaseArgs, *resourcespb.DatabaseResource]("database", serviceContext),
		services.NewService[*resourcespb.ObjectStorageArgs, *resourcespb.ObjectStorageResource]("object_storage", serviceContext),
		services.NewService[*resourcespb.ObjectStorageObjectArgs, *resourcespb.ObjectStorageObjectResource]("object_storage_object", serviceContext),
		services.NewService[*resourcespb.PublicIpArgs, *resourcespb.PublicIpResource]("public_ip", serviceContext),
		services.NewService[*resourcespb.KubernetesClusterArgs, *resourcespb.KubernetesClusterResource]("kubernetes_cluster", serviceContext),
		services.NewService[*resourcespb.KubernetesNodePoolArgs, *resourcespb.KubernetesNodePoolResource]("kubernetes_node_pool", serviceContext),
		services.NewService[*resourcespb.LambdaArgs, *resourcespb.LambdaResource]("lambda", serviceContext),
		services.NewService[*resourcespb.VaultArgs, *resourcespb.VaultResource]("vault", serviceContext),
		services.NewService[*resourcespb.VaultAccessPolicyArgs, *resourcespb.VaultAccessPolicyResource]("vault_access_policy", serviceContext),
		services.NewService[*resourcespb.VaultSecretArgs, *resourcespb.VaultSecretResource]("vault_secret", serviceContext),
		services.NewService[*resourcespb.VirtualMachineArgs, *resourcespb.VirtualMachineResource]("virtual_machine", serviceContext),
	}
}

func (s *Server) CreateVirtualNetwork(ctx context.Context, in *resourcespb.CreateVirtualNetworkRequest) (*resourcespb.VirtualNetworkResource, error) {
	return s.VnService.Create(ctx, in)
}
func (s *Server) ReadVirtualNetwork(ctx context.Context, in *resourcespb.ReadVirtualNetworkRequest) (*resourcespb.VirtualNetworkResource, error) {
	return s.VnService.Read(ctx, in)
}
func (s *Server) UpdateVirtualNetwork(ctx context.Context, in *resourcespb.UpdateVirtualNetworkRequest) (*resourcespb.VirtualNetworkResource, error) {
	return s.VnService.Update(ctx, in)
}
func (s *Server) DeleteVirtualNetwork(ctx context.Context, in *resourcespb.DeleteVirtualNetworkRequest) (*commonpb.Empty, error) {
	return s.VnService.Delete(ctx, in)
}

func (s *Server) CreateSubnet(ctx context.Context, in *resourcespb.CreateSubnetRequest) (*resourcespb.SubnetResource, error) {
	return s.SubnetService.Create(ctx, in)
}
func (s *Server) ReadSubnet(ctx context.Context, in *resourcespb.ReadSubnetRequest) (*resourcespb.SubnetResource, error) {
	return s.SubnetService.Read(ctx, in)
}
func (s *Server) UpdateSubnet(ctx context.Context, in *resourcespb.UpdateSubnetRequest) (*resourcespb.SubnetResource, error) {
	return s.SubnetService.Update(ctx, in)
}
func (s *Server) DeleteSubnet(ctx context.Context, in *resourcespb.DeleteSubnetRequest) (*commonpb.Empty, error) {
	return s.SubnetService.Delete(ctx, in)
}

func (s *Server) CreateNetworkInterface(ctx context.Context, in *resourcespb.CreateNetworkInterfaceRequest) (*resourcespb.NetworkInterfaceResource, error) {
	return s.NetworkInterfaceService.Create(ctx, in)
}
func (s *Server) ReadNetworkInterface(ctx context.Context, in *resourcespb.ReadNetworkInterfaceRequest) (*resourcespb.NetworkInterfaceResource, error) {
	return s.NetworkInterfaceService.Read(ctx, in)
}
func (s *Server) UpdateNetworkInterface(ctx context.Context, in *resourcespb.UpdateNetworkInterfaceRequest) (*resourcespb.NetworkInterfaceResource, error) {
	return s.NetworkInterfaceService.Update(ctx, in)
}
func (s *Server) DeleteNetworkInterface(ctx context.Context, in *resourcespb.DeleteNetworkInterfaceRequest) (*commonpb.Empty, error) {
	return s.NetworkInterfaceService.Delete(ctx, in)
}

func (s *Server) CreateRouteTable(ctx context.Context, in *resourcespb.CreateRouteTableRequest) (*resourcespb.RouteTableResource, error) {
	return s.RouteTableService.Create(ctx, in)
}
func (s *Server) ReadRouteTable(ctx context.Context, in *resourcespb.ReadRouteTableRequest) (*resourcespb.RouteTableResource, error) {
	return s.RouteTableService.Read(ctx, in)
}
func (s *Server) UpdateRouteTable(ctx context.Context, in *resourcespb.UpdateRouteTableRequest) (*resourcespb.RouteTableResource, error) {
	return s.RouteTableService.Update(ctx, in)
}
func (s *Server) DeleteRouteTable(ctx context.Context, in *resourcespb.DeleteRouteTableRequest) (*commonpb.Empty, error) {
	return s.RouteTableService.Delete(ctx, in)
}

func (s *Server) CreateRouteTableAssociation(ctx context.Context, in *resourcespb.CreateRouteTableAssociationRequest) (*resourcespb.RouteTableAssociationResource, error) {
	return s.RouteTableAssociationService.Create(ctx, in)
}
func (s *Server) ReadRouteTableAssociation(ctx context.Context, in *resourcespb.ReadRouteTableAssociationRequest) (*resourcespb.RouteTableAssociationResource, error) {
	return s.RouteTableAssociationService.Read(ctx, in)
}
func (s *Server) UpdateRouteTableAssociation(ctx context.Context, in *resourcespb.UpdateRouteTableAssociationRequest) (*resourcespb.RouteTableAssociationResource, error) {
	return s.RouteTableAssociationService.Update(ctx, in)
}
func (s *Server) DeleteRouteTableAssociation(ctx context.Context, in *resourcespb.DeleteRouteTableAssociationRequest) (*commonpb.Empty, error) {
	return s.RouteTableService.Delete(ctx, in)
}

func (s *Server) CreateNetworkSecurityGroup(ctx context.Context, in *resourcespb.CreateNetworkSecurityGroupRequest) (*resourcespb.NetworkSecurityGroupResource, error) {
	return s.NetworkSecurityGroupService.Create(ctx, in)
}
func (s *Server) ReadNetworkSecurityGroup(ctx context.Context, in *resourcespb.ReadNetworkSecurityGroupRequest) (*resourcespb.NetworkSecurityGroupResource, error) {
	return s.NetworkSecurityGroupService.Read(ctx, in)
}
func (s *Server) UpdateNetworkSecurityGroup(ctx context.Context, in *resourcespb.UpdateNetworkSecurityGroupRequest) (*resourcespb.NetworkSecurityGroupResource, error) {
	return s.NetworkSecurityGroupService.Update(ctx, in)
}
func (s *Server) DeleteNetworkSecurityGroup(ctx context.Context, in *resourcespb.DeleteNetworkSecurityGroupRequest) (*commonpb.Empty, error) {
	return s.NetworkSecurityGroupService.Delete(ctx, in)
}

func (s *Server) CreateDatabase(ctx context.Context, in *resourcespb.CreateDatabaseRequest) (*resourcespb.DatabaseResource, error) {
	return s.DatabaseService.Create(ctx, in)
}
func (s *Server) ReadDatabase(ctx context.Context, in *resourcespb.ReadDatabaseRequest) (*resourcespb.DatabaseResource, error) {
	return s.DatabaseService.Read(ctx, in)
}
func (s *Server) UpdateDatabase(ctx context.Context, in *resourcespb.UpdateDatabaseRequest) (*resourcespb.DatabaseResource, error) {
	return s.DatabaseService.Update(ctx, in)
}
func (s *Server) DeleteDatabase(ctx context.Context, in *resourcespb.DeleteDatabaseRequest) (*commonpb.Empty, error) {
	return s.DatabaseService.Delete(ctx, in)
}

func (s *Server) CreateObjectStorage(ctx context.Context, in *resourcespb.CreateObjectStorageRequest) (*resourcespb.ObjectStorageResource, error) {
	return s.ObjectStorageService.Create(ctx, in)
}
func (s *Server) ReadObjectStorage(ctx context.Context, in *resourcespb.ReadObjectStorageRequest) (*resourcespb.ObjectStorageResource, error) {
	return s.ObjectStorageService.Read(ctx, in)
}
func (s *Server) UpdateObjectStorage(ctx context.Context, in *resourcespb.UpdateObjectStorageRequest) (*resourcespb.ObjectStorageResource, error) {
	return s.ObjectStorageService.Update(ctx, in)
}
func (s *Server) DeleteObjectStorage(ctx context.Context, in *resourcespb.DeleteObjectStorageRequest) (*commonpb.Empty, error) {
	return s.ObjectStorageService.Delete(ctx, in)
}

func (s *Server) CreateObjectStorageObject(ctx context.Context, in *resourcespb.CreateObjectStorageObjectRequest) (*resourcespb.ObjectStorageObjectResource, error) {
	return s.ObjectStorageObjectService.Create(ctx, in)
}
func (s *Server) ReadObjectStorageObject(ctx context.Context, in *resourcespb.ReadObjectStorageObjectRequest) (*resourcespb.ObjectStorageObjectResource, error) {
	return s.ObjectStorageObjectService.Read(ctx, in)
}
func (s *Server) UpdateObjectStorageObject(ctx context.Context, in *resourcespb.UpdateObjectStorageObjectRequest) (*resourcespb.ObjectStorageObjectResource, error) {
	return s.ObjectStorageObjectService.Update(ctx, in)
}
func (s *Server) DeleteObjectStorageObject(ctx context.Context, in *resourcespb.DeleteObjectStorageObjectRequest) (*commonpb.Empty, error) {
	return s.ObjectStorageObjectService.Delete(ctx, in)
}

func (s *Server) CreatePublicIp(ctx context.Context, in *resourcespb.CreatePublicIpRequest) (*resourcespb.PublicIpResource, error) {
	return s.PublicIpService.Create(ctx, in)
}
func (s *Server) ReadPublicIp(ctx context.Context, in *resourcespb.ReadPublicIpRequest) (*resourcespb.PublicIpResource, error) {
	return s.PublicIpService.Read(ctx, in)
}
func (s *Server) UpdatePublicIp(ctx context.Context, in *resourcespb.UpdatePublicIpRequest) (*resourcespb.PublicIpResource, error) {
	return s.PublicIpService.Update(ctx, in)
}
func (s *Server) DeletePublicIp(ctx context.Context, in *resourcespb.DeletePublicIpRequest) (*commonpb.Empty, error) {
	return s.PublicIpService.Delete(ctx, in)
}

func (s *Server) CreateKubernetesCluster(ctx context.Context, in *resourcespb.CreateKubernetesClusterRequest) (*resourcespb.KubernetesClusterResource, error) {
	return s.KubernetesClusterService.Create(ctx, in)
}
func (s *Server) ReadKubernetesCluster(ctx context.Context, in *resourcespb.ReadKubernetesClusterRequest) (*resourcespb.KubernetesClusterResource, error) {
	return s.KubernetesClusterService.Read(ctx, in)
}
func (s *Server) UpdateKubernetesCluster(ctx context.Context, in *resourcespb.UpdateKubernetesClusterRequest) (*resourcespb.KubernetesClusterResource, error) {
	return s.KubernetesClusterService.Update(ctx, in)
}
func (s *Server) DeleteKubernetesCluster(ctx context.Context, in *resourcespb.DeleteKubernetesClusterRequest) (*commonpb.Empty, error) {
	return s.KubernetesClusterService.Delete(ctx, in)
}

func (s *Server) CreateKubernetesNodePool(ctx context.Context, in *resourcespb.CreateKubernetesNodePoolRequest) (*resourcespb.KubernetesNodePoolResource, error) {
	return s.KubernetesNodePoolService.Create(ctx, in)
}
func (s *Server) ReadKubernetesNodePool(ctx context.Context, in *resourcespb.ReadKubernetesNodePoolRequest) (*resourcespb.KubernetesNodePoolResource, error) {
	return s.KubernetesNodePoolService.Read(ctx, in)
}
func (s *Server) UpdateKubernetesNodePool(ctx context.Context, in *resourcespb.UpdateKubernetesNodePoolRequest) (*resourcespb.KubernetesNodePoolResource, error) {
	return s.KubernetesNodePoolService.Update(ctx, in)
}
func (s *Server) DeleteKubernetesNodePool(ctx context.Context, in *resourcespb.DeleteKubernetesNodePoolRequest) (*commonpb.Empty, error) {
	return s.KubernetesNodePoolService.Delete(ctx, in)
}

func (s *Server) CreateLambda(ctx context.Context, in *resourcespb.CreateLambdaRequest) (*resourcespb.LambdaResource, error) {
	return s.LambdaService.Create(ctx, in)
}
func (s *Server) ReadLambda(ctx context.Context, in *resourcespb.ReadLambdaRequest) (*resourcespb.LambdaResource, error) {
	return s.LambdaService.Read(ctx, in)
}
func (s *Server) UpdateLambda(ctx context.Context, in *resourcespb.UpdateLambdaRequest) (*resourcespb.LambdaResource, error) {
	return s.LambdaService.Update(ctx, in)
}
func (s *Server) DeleteLambda(ctx context.Context, in *resourcespb.DeleteLambdaRequest) (*commonpb.Empty, error) {
	return s.LambdaService.Delete(ctx, in)
}

func (s *Server) CreateVault(ctx context.Context, in *resourcespb.CreateVaultRequest) (*resourcespb.VaultResource, error) {
	return s.VaultService.Create(ctx, in)
}
func (s *Server) ReadVault(ctx context.Context, in *resourcespb.ReadVaultRequest) (*resourcespb.VaultResource, error) {
	return s.VaultService.Read(ctx, in)
}
func (s *Server) UpdateVault(ctx context.Context, in *resourcespb.UpdateVaultRequest) (*resourcespb.VaultResource, error) {
	return s.VaultService.Update(ctx, in)
}
func (s *Server) DeleteVault(ctx context.Context, in *resourcespb.DeleteVaultRequest) (*commonpb.Empty, error) {
	return s.VaultService.Delete(ctx, in)
}

func (s *Server) CreateVaultSecret(ctx context.Context, in *resourcespb.CreateVaultSecretRequest) (*resourcespb.VaultSecretResource, error) {
	return s.VaultSecretService.Create(ctx, in)
}
func (s *Server) ReadVaultSecret(ctx context.Context, in *resourcespb.ReadVaultSecretRequest) (*resourcespb.VaultSecretResource, error) {
	return s.VaultSecretService.Read(ctx, in)
}
func (s *Server) UpdateVaultSecret(ctx context.Context, in *resourcespb.UpdateVaultSecretRequest) (*resourcespb.VaultSecretResource, error) {
	return s.VaultSecretService.Update(ctx, in)
}
func (s *Server) DeleteVaultSecret(ctx context.Context, in *resourcespb.DeleteVaultSecretRequest) (*commonpb.Empty, error) {
	return s.VaultSecretService.Delete(ctx, in)
}

func (s *Server) CreateVaultAccessPolicy(ctx context.Context, in *resourcespb.CreateVaultAccessPolicyRequest) (*resourcespb.VaultAccessPolicyResource, error) {
	return s.VaultAccessPolicyService.Create(ctx, in)
}
func (s *Server) ReadVaultAccessPolicy(ctx context.Context, in *resourcespb.ReadVaultAccessPolicyRequest) (*resourcespb.VaultAccessPolicyResource, error) {
	return s.VaultAccessPolicyService.Read(ctx, in)
}
func (s *Server) UpdateVaultAccessPolicy(ctx context.Context, in *resourcespb.UpdateVaultAccessPolicyRequest) (*resourcespb.VaultAccessPolicyResource, error) {
	return s.VaultAccessPolicyService.Update(ctx, in)
}
func (s *Server) DeleteVaultAccessPolicy(ctx context.Context, in *resourcespb.DeleteVaultAccessPolicyRequest) (*commonpb.Empty, error) {
	return s.VaultAccessPolicyService.Delete(ctx, in)
}

func (s *Server) CreateVirtualMachine(ctx context.Context, in *resourcespb.CreateVirtualMachineRequest) (*resourcespb.VirtualMachineResource, error) {
	return s.VirtualMachineService.Create(ctx, in)
}
func (s *Server) ReadVirtualMachine(ctx context.Context, in *resourcespb.ReadVirtualMachineRequest) (*resourcespb.VirtualMachineResource, error) {
	return s.VirtualMachineService.Read(ctx, in)
}
func (s *Server) UpdateVirtualMachine(ctx context.Context, in *resourcespb.UpdateVirtualMachineRequest) (*resourcespb.VirtualMachineResource, error) {
	return s.VirtualMachineService.Update(ctx, in)
}
func (s *Server) DeleteVirtualMachine(ctx context.Context, in *resourcespb.DeleteVirtualMachineRequest) (*commonpb.Empty, error) {
	return s.VirtualMachineService.Delete(ctx, in)
}

func (s *Server) RefreshState(ctx context.Context, _ *commonpb.Empty) (_ *commonpb.Empty, err error) {
	defer func() {
		if err != nil {
			go s.AwsClient.UpdateErrorMetric("refresh", "refresh", errors.ErrorCode(err))
		}
	}()
	return errors.WrappingErrors(s.refresh)(ctx, &commonpb.Empty{})
}

func (s *Server) refresh(ctx context.Context, _ *commonpb.Empty) (*commonpb.Empty, error) {
	log.Println("[INFO] Refreshing state")
	key, err := util.ExtractApiKey(ctx)
	if err != nil {
		return nil, err
	}

	userId, err := s.Database.GetUserId(ctx, key)
	if err != nil {
		return nil, err
	}

	lock, err := s.Database.LockConfig(ctx, userId)
	if err != nil {
		return nil, err
	}
	defer s.Database.UnlockConfig(ctx, lock)

	c, err := s.Database.LoadUserConfig(userId, lock)
	if err != nil {
		return nil, err
	}
	mconfig, err := resources.LoadConfig(c, types.Metadatas)
	if err != nil {
		return nil, err
	}

	_, err = deploy.EncodeAndStoreTfFile(ctx, mconfig, nil, nil, true)
	if err != nil {
		return nil, err
	}

	err = deploy.MaybeInit(ctx, userId, true)
	if err != nil {
		return nil, err
	}

	err = deploy.RefreshState(ctx, userId, true)
	if err != nil {
		return nil, err
	}

	err = s.Database.StoreUserConfig(c, lock)
	if err != nil {
		return nil, err
	}

	return &commonpb.Empty{}, nil
}

func (s *Server) ListResources(ctx context.Context, _ *commonpb.Empty) (resp *commonpb.ListResourcesResponse, err error) {
	defer func() {
		if err != nil {
			go s.AwsClient.UpdateErrorMetric("list", "list", errors.ErrorCode(err))
		}
	}()
	return errors.WrappingErrors(s.list)(ctx, &commonpb.Empty{})
}

func (s *Server) list(ctx context.Context, _ *commonpb.Empty) (*commonpb.ListResourcesResponse, error) {
	log.Println("[INFO] Listing resources")
	key, err := util.ExtractApiKey(ctx)
	if err != nil {
		return nil, err
	}

	userId, err := s.Database.GetUserId(ctx, key)
	if err != nil {
		return nil, err
	}

	c, err := s.Database.LoadUserConfig(userId, nil)
	if err != nil {
		return nil, err
	}

	resp := &commonpb.ListResourcesResponse{}
	for _, r := range c.Resources {
		name := string(r.ResourceArgs.ResourceArgs.MessageName().Name())
		name = strings.TrimSuffix(name, "Args")

		resp.Resources = append(resp.Resources, &commonpb.ListResourcesResponse_ResourceMetadata{
			ResourceId:   r.ResourceId,
			ResourceType: name,
		})
	}

	return resp, nil
}

func (s *Server) DeleteResource(ctx context.Context, req *proto.DeleteResourceRequest) (_ *commonpb.Empty, err error) {
	defer func() {
		if err != nil {
			go s.AwsClient.UpdateErrorMetric("delete", "delete", errors.ErrorCode(err))
		}
	}()
	return errors.WrappingErrors(s.deleteResource)(ctx, req)
}

func (s *Server) deleteResource(ctx context.Context, req *proto.DeleteResourceRequest) (*commonpb.Empty, error) {
	log.Println("[INFO] Deleting resource")
	key, err := util.ExtractApiKey(ctx)
	if err != nil {
		return nil, err
	}

	userId, err := s.Database.GetUserId(ctx, key)
	if err != nil {
		return nil, err
	}

	lock, err := s.Database.LockConfig(ctx, userId)
	if err != nil {
		return nil, err
	}
	defer s.Database.UnlockConfig(ctx, lock)
	c, err := s.Database.LoadUserConfig(userId, lock)
	if err != nil {
		return nil, err
	}

	err = util.DeleteResourceFromConfig(c, req.ResourceId)
	if err != nil {
		return nil, err
	}

	err = s.Database.StoreUserConfig(c, lock)
	if err != nil {
		return nil, err
	}

	return &commonpb.Empty{}, nil
}
