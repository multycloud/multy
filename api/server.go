package api

import (
	"context"
	"fmt"
	"github.com/multycloud/multy/api/deploy"
	"github.com/multycloud/multy/api/errors"
	"github.com/multycloud/multy/api/proto"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/api/services/database"
	"github.com/multycloud/multy/api/services/kubernetes_cluster"
	"github.com/multycloud/multy/api/services/kubernetes_node_pool"
	"github.com/multycloud/multy/api/services/lambda"
	"github.com/multycloud/multy/api/services/network_interface"
	"github.com/multycloud/multy/api/services/network_security_group"
	"github.com/multycloud/multy/api/services/object_storage"
	"github.com/multycloud/multy/api/services/object_storage_object"
	"github.com/multycloud/multy/api/services/public_ip"
	"github.com/multycloud/multy/api/services/route_table"
	"github.com/multycloud/multy/api/services/route_table_association"
	"github.com/multycloud/multy/api/services/subnet"
	"github.com/multycloud/multy/api/services/vault"
	"github.com/multycloud/multy/api/services/virtual_machine"
	"github.com/multycloud/multy/api/services/virtual_network"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/db"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"os"
)

type Server struct {
	proto.UnimplementedMultyResourceServiceServer
	*db.Database
	virtual_network.VnService
	subnet.SubnetService
	network_interface.NetworkInterfaceService
	route_table.RouteTableService
	route_table_association.RouteTableAssociationService
	network_security_group.NetworkSecurityGroupService
	database.DatabaseService
	object_storage.ObjectStorageService
	object_storage_object.ObjectStorageObjectService
	public_ip.PublicIpService
	kubernetes_cluster.KubernetesClusterService
	kubernetes_node_pool.KubernetesNodePoolService
	lambda.LambdaService
	vault.VaultService
	vault.VaultAccessPolicyService
	vault.VaultSecretService
	virtual_machine.VirtualMachineService
}

func RunServer(ctx context.Context, port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	connectionStr, exists := os.LookupEnv("MULTY_DB_CONN_STRING")
	if !exists {
		log.Fatalf("db_connection_string env var is not set")
	}

	endpoint, exists := os.LookupEnv("MULTY_API_ENDPOINT")
	if !exists {
		log.Fatalf("api_endpoint env var is not set")
	}

	certFile := fmt.Sprintf("/etc/letsencrypt/live/%s/fullchain.pem", endpoint)
	keyFile := fmt.Sprintf("/etc/letsencrypt/live/%s/privkey.pem", endpoint)
	creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	var s *grpc.Server
	if err != nil {
		log.Printf("unable to read certificate (%s), running in insecure mode", err.Error())
		s = grpc.NewServer()
	} else {
		s = grpc.NewServer(grpc.Creds(creds))
	}

	go func() {
		<-ctx.Done()
		s.GracefulStop()
		_ = lis.Close()
	}()
	d, err := db.NewDatabase(connectionStr)
	if err != nil {
		log.Fatalf("failed to load db: %v", err)
	}
	defer d.Close()
	server := Server{
		proto.UnimplementedMultyResourceServiceServer{},
		d,
		virtual_network.NewVnService(d),
		subnet.NewSubnetService(d),
		network_interface.NewNetworkInterfaceService(d),
		route_table.NewRouteTableService(d),
		route_table_association.NewRouteTableAssociationService(d),
		network_security_group.NewNetworkSecurityGroupService(d),
		database.NewDatabaseService(d),
		object_storage.NewObjectStorageService(d),
		object_storage_object.NewObjectStorageObjectService(d),
		public_ip.NewPublicIpService(d),
		kubernetes_cluster.NewKubernetesClusterService(d),
		kubernetes_node_pool.NewKubernetesNodePoolService(d),
		lambda.NewLambdaService(d),
		vault.NewVaultService(d),
		vault.NewVaultAccessPolicyService(d),
		vault.NewVaultSecretService(d),
		virtual_machine.NewVirtualMachineService(d),
	}
	proto.RegisterMultyResourceServiceServer(s, &server)
	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *Server) CreateVirtualNetwork(ctx context.Context, in *resourcespb.CreateVirtualNetworkRequest) (*resourcespb.VirtualNetworkResource, error) {
	return s.VnService.Service.Create(ctx, in)
}
func (s *Server) ReadVirtualNetwork(ctx context.Context, in *resourcespb.ReadVirtualNetworkRequest) (*resourcespb.VirtualNetworkResource, error) {
	return s.VnService.Service.Read(ctx, in)
}
func (s *Server) UpdateVirtualNetwork(ctx context.Context, in *resourcespb.UpdateVirtualNetworkRequest) (*resourcespb.VirtualNetworkResource, error) {
	return s.VnService.Service.Update(ctx, in)
}
func (s *Server) DeleteVirtualNetwork(ctx context.Context, in *resourcespb.DeleteVirtualNetworkRequest) (*commonpb.Empty, error) {
	return s.VnService.Service.Delete(ctx, in)
}

func (s *Server) CreateSubnet(ctx context.Context, in *resourcespb.CreateSubnetRequest) (*resourcespb.SubnetResource, error) {
	return s.SubnetService.Service.Create(ctx, in)
}
func (s *Server) ReadSubnet(ctx context.Context, in *resourcespb.ReadSubnetRequest) (*resourcespb.SubnetResource, error) {
	return s.SubnetService.Service.Read(ctx, in)
}
func (s *Server) UpdateSubnet(ctx context.Context, in *resourcespb.UpdateSubnetRequest) (*resourcespb.SubnetResource, error) {
	return s.SubnetService.Service.Update(ctx, in)
}
func (s *Server) DeleteSubnet(ctx context.Context, in *resourcespb.DeleteSubnetRequest) (*commonpb.Empty, error) {
	return s.SubnetService.Service.Delete(ctx, in)
}

func (s *Server) CreateNetworkInterface(ctx context.Context, in *resourcespb.CreateNetworkInterfaceRequest) (*resourcespb.NetworkInterfaceResource, error) {
	return s.NetworkInterfaceService.Service.Create(ctx, in)
}
func (s *Server) ReadNetworkInterface(ctx context.Context, in *resourcespb.ReadNetworkInterfaceRequest) (*resourcespb.NetworkInterfaceResource, error) {
	return s.NetworkInterfaceService.Service.Read(ctx, in)
}
func (s *Server) UpdateNetworkInterface(ctx context.Context, in *resourcespb.UpdateNetworkInterfaceRequest) (*resourcespb.NetworkInterfaceResource, error) {
	return s.NetworkInterfaceService.Service.Update(ctx, in)
}
func (s *Server) DeleteNetworkInterface(ctx context.Context, in *resourcespb.DeleteNetworkInterfaceRequest) (*commonpb.Empty, error) {
	return s.NetworkInterfaceService.Service.Delete(ctx, in)
}

func (s *Server) CreateRouteTable(ctx context.Context, in *resourcespb.CreateRouteTableRequest) (*resourcespb.RouteTableResource, error) {
	return s.RouteTableService.Service.Create(ctx, in)
}
func (s *Server) ReadRouteTable(ctx context.Context, in *resourcespb.ReadRouteTableRequest) (*resourcespb.RouteTableResource, error) {
	return s.RouteTableService.Service.Read(ctx, in)
}
func (s *Server) UpdateRouteTable(ctx context.Context, in *resourcespb.UpdateRouteTableRequest) (*resourcespb.RouteTableResource, error) {
	return s.RouteTableService.Service.Update(ctx, in)
}
func (s *Server) DeleteRouteTable(ctx context.Context, in *resourcespb.DeleteRouteTableRequest) (*commonpb.Empty, error) {
	return s.RouteTableService.Service.Delete(ctx, in)
}

func (s *Server) CreateRouteTableAssociation(ctx context.Context, in *resourcespb.CreateRouteTableAssociationRequest) (*resourcespb.RouteTableAssociationResource, error) {
	return s.RouteTableAssociationService.Service.Create(ctx, in)
}
func (s *Server) ReadRouteTableAssociation(ctx context.Context, in *resourcespb.ReadRouteTableAssociationRequest) (*resourcespb.RouteTableAssociationResource, error) {
	return s.RouteTableAssociationService.Service.Read(ctx, in)
}
func (s *Server) UpdateRouteTableAssociation(ctx context.Context, in *resourcespb.UpdateRouteTableAssociationRequest) (*resourcespb.RouteTableAssociationResource, error) {
	return s.RouteTableAssociationService.Service.Update(ctx, in)
}
func (s *Server) DeleteRouteTableAssociation(ctx context.Context, in *resourcespb.DeleteRouteTableAssociationRequest) (*commonpb.Empty, error) {
	return s.RouteTableService.Service.Delete(ctx, in)
}

func (s *Server) CreateNetworkSecurityGroup(ctx context.Context, in *resourcespb.CreateNetworkSecurityGroupRequest) (*resourcespb.NetworkSecurityGroupResource, error) {
	return s.NetworkSecurityGroupService.Service.Create(ctx, in)
}
func (s *Server) ReadNetworkSecurityGroup(ctx context.Context, in *resourcespb.ReadNetworkSecurityGroupRequest) (*resourcespb.NetworkSecurityGroupResource, error) {
	return s.NetworkSecurityGroupService.Service.Read(ctx, in)
}
func (s *Server) UpdateNetworkSecurityGroup(ctx context.Context, in *resourcespb.UpdateNetworkSecurityGroupRequest) (*resourcespb.NetworkSecurityGroupResource, error) {
	return s.NetworkSecurityGroupService.Service.Update(ctx, in)
}
func (s *Server) DeleteNetworkSecurityGroup(ctx context.Context, in *resourcespb.DeleteNetworkSecurityGroupRequest) (*commonpb.Empty, error) {
	return s.NetworkSecurityGroupService.Service.Delete(ctx, in)
}

func (s *Server) CreateDatabase(ctx context.Context, in *resourcespb.CreateDatabaseRequest) (*resourcespb.DatabaseResource, error) {
	return s.DatabaseService.Service.Create(ctx, in)
}
func (s *Server) ReadDatabase(ctx context.Context, in *resourcespb.ReadDatabaseRequest) (*resourcespb.DatabaseResource, error) {
	return s.DatabaseService.Service.Read(ctx, in)
}
func (s *Server) UpdateDatabase(ctx context.Context, in *resourcespb.UpdateDatabaseRequest) (*resourcespb.DatabaseResource, error) {
	return s.DatabaseService.Service.Update(ctx, in)
}
func (s *Server) DeleteDatabase(ctx context.Context, in *resourcespb.DeleteDatabaseRequest) (*commonpb.Empty, error) {
	return s.DatabaseService.Service.Delete(ctx, in)
}

func (s *Server) CreateObjectStorage(ctx context.Context, in *resourcespb.CreateObjectStorageRequest) (*resourcespb.ObjectStorageResource, error) {
	return s.ObjectStorageService.Service.Create(ctx, in)
}
func (s *Server) ReadObjectStorage(ctx context.Context, in *resourcespb.ReadObjectStorageRequest) (*resourcespb.ObjectStorageResource, error) {
	return s.ObjectStorageService.Service.Read(ctx, in)
}
func (s *Server) UpdateObjectStorage(ctx context.Context, in *resourcespb.UpdateObjectStorageRequest) (*resourcespb.ObjectStorageResource, error) {
	return s.ObjectStorageService.Service.Update(ctx, in)
}
func (s *Server) DeleteObjectStorage(ctx context.Context, in *resourcespb.DeleteObjectStorageRequest) (*commonpb.Empty, error) {
	return s.ObjectStorageService.Service.Delete(ctx, in)
}

func (s *Server) CreateObjectStorageObject(ctx context.Context, in *resourcespb.CreateObjectStorageObjectRequest) (*resourcespb.ObjectStorageObjectResource, error) {
	return s.ObjectStorageObjectService.Service.Create(ctx, in)
}
func (s *Server) ReadObjectStorageObject(ctx context.Context, in *resourcespb.ReadObjectStorageObjectRequest) (*resourcespb.ObjectStorageObjectResource, error) {
	return s.ObjectStorageObjectService.Service.Read(ctx, in)
}
func (s *Server) UpdateObjectStorageObject(ctx context.Context, in *resourcespb.UpdateObjectStorageObjectRequest) (*resourcespb.ObjectStorageObjectResource, error) {
	return s.ObjectStorageObjectService.Service.Update(ctx, in)
}
func (s *Server) DeleteObjectStorageObject(ctx context.Context, in *resourcespb.DeleteObjectStorageObjectRequest) (*commonpb.Empty, error) {
	return s.ObjectStorageObjectService.Service.Delete(ctx, in)
}

func (s *Server) CreatePublicIp(ctx context.Context, in *resourcespb.CreatePublicIpRequest) (*resourcespb.PublicIpResource, error) {
	return s.PublicIpService.Service.Create(ctx, in)
}
func (s *Server) ReadPublicIp(ctx context.Context, in *resourcespb.ReadPublicIpRequest) (*resourcespb.PublicIpResource, error) {
	return s.PublicIpService.Service.Read(ctx, in)
}
func (s *Server) UpdatePublicIp(ctx context.Context, in *resourcespb.UpdatePublicIpRequest) (*resourcespb.PublicIpResource, error) {
	return s.PublicIpService.Service.Update(ctx, in)
}
func (s *Server) DeletePublicIp(ctx context.Context, in *resourcespb.DeletePublicIpRequest) (*commonpb.Empty, error) {
	return s.PublicIpService.Service.Delete(ctx, in)
}

func (s *Server) CreateKubernetesCluster(ctx context.Context, in *resourcespb.CreateKubernetesClusterRequest) (*resourcespb.KubernetesClusterResource, error) {
	return s.KubernetesClusterService.Service.Create(ctx, in)
}
func (s *Server) ReadKubernetesCluster(ctx context.Context, in *resourcespb.ReadKubernetesClusterRequest) (*resourcespb.KubernetesClusterResource, error) {
	return s.KubernetesClusterService.Service.Read(ctx, in)
}
func (s *Server) UpdateKubernetesCluster(ctx context.Context, in *resourcespb.UpdateKubernetesClusterRequest) (*resourcespb.KubernetesClusterResource, error) {
	return s.KubernetesClusterService.Service.Update(ctx, in)
}
func (s *Server) DeleteKubernetesCluster(ctx context.Context, in *resourcespb.DeleteKubernetesClusterRequest) (*commonpb.Empty, error) {
	return s.KubernetesClusterService.Service.Delete(ctx, in)
}

func (s *Server) CreateKubernetesNodePool(ctx context.Context, in *resourcespb.CreateKubernetesNodePoolRequest) (*resourcespb.KubernetesNodePoolResource, error) {
	return s.KubernetesNodePoolService.Service.Create(ctx, in)
}
func (s *Server) ReadKubernetesNodePool(ctx context.Context, in *resourcespb.ReadKubernetesNodePoolRequest) (*resourcespb.KubernetesNodePoolResource, error) {
	return s.KubernetesNodePoolService.Service.Read(ctx, in)
}
func (s *Server) UpdateKubernetesNodePool(ctx context.Context, in *resourcespb.UpdateKubernetesNodePoolRequest) (*resourcespb.KubernetesNodePoolResource, error) {
	return s.KubernetesNodePoolService.Service.Update(ctx, in)
}
func (s *Server) DeleteKubernetesNodePool(ctx context.Context, in *resourcespb.DeleteKubernetesNodePoolRequest) (*commonpb.Empty, error) {
	return s.KubernetesNodePoolService.Service.Delete(ctx, in)
}

func (s *Server) CreateLambda(ctx context.Context, in *resourcespb.CreateLambdaRequest) (*resourcespb.LambdaResource, error) {
	return s.LambdaService.Service.Create(ctx, in)
}
func (s *Server) ReadLambda(ctx context.Context, in *resourcespb.ReadLambdaRequest) (*resourcespb.LambdaResource, error) {
	return s.LambdaService.Service.Read(ctx, in)
}
func (s *Server) UpdateLambda(ctx context.Context, in *resourcespb.UpdateLambdaRequest) (*resourcespb.LambdaResource, error) {
	return s.LambdaService.Service.Update(ctx, in)
}
func (s *Server) DeleteLambda(ctx context.Context, in *resourcespb.DeleteLambdaRequest) (*commonpb.Empty, error) {
	return s.LambdaService.Service.Delete(ctx, in)
}

func (s *Server) CreateVault(ctx context.Context, in *resourcespb.CreateVaultRequest) (*resourcespb.VaultResource, error) {
	return s.VaultService.Service.Create(ctx, in)
}
func (s *Server) ReadVault(ctx context.Context, in *resourcespb.ReadVaultRequest) (*resourcespb.VaultResource, error) {
	return s.VaultService.Service.Read(ctx, in)
}
func (s *Server) UpdateVault(ctx context.Context, in *resourcespb.UpdateVaultRequest) (*resourcespb.VaultResource, error) {
	return s.VaultService.Service.Update(ctx, in)
}
func (s *Server) DeleteVault(ctx context.Context, in *resourcespb.DeleteVaultRequest) (*commonpb.Empty, error) {
	return s.VaultService.Service.Delete(ctx, in)
}

func (s *Server) CreateVaultSecret(ctx context.Context, in *resourcespb.CreateVaultSecretRequest) (*resourcespb.VaultSecretResource, error) {
	return s.VaultSecretService.Service.Create(ctx, in)
}
func (s *Server) ReadVaultSecret(ctx context.Context, in *resourcespb.ReadVaultSecretRequest) (*resourcespb.VaultSecretResource, error) {
	return s.VaultSecretService.Service.Read(ctx, in)
}
func (s *Server) UpdateVaultSecret(ctx context.Context, in *resourcespb.UpdateVaultSecretRequest) (*resourcespb.VaultSecretResource, error) {
	return s.VaultSecretService.Service.Update(ctx, in)
}
func (s *Server) DeleteVaultSecret(ctx context.Context, in *resourcespb.DeleteVaultSecretRequest) (*commonpb.Empty, error) {
	return s.VaultSecretService.Service.Delete(ctx, in)
}

func (s *Server) CreateVaultAccessPolicy(ctx context.Context, in *resourcespb.CreateVaultAccessPolicyRequest) (*resourcespb.VaultAccessPolicyResource, error) {
	return s.VaultAccessPolicyService.Service.Create(ctx, in)
}
func (s *Server) ReadVaultAccessPolicy(ctx context.Context, in *resourcespb.ReadVaultAccessPolicyRequest) (*resourcespb.VaultAccessPolicyResource, error) {
	return s.VaultAccessPolicyService.Service.Read(ctx, in)
}
func (s *Server) UpdateVaultAccessPolicy(ctx context.Context, in *resourcespb.UpdateVaultAccessPolicyRequest) (*resourcespb.VaultAccessPolicyResource, error) {
	return s.VaultAccessPolicyService.Service.Update(ctx, in)
}
func (s *Server) DeleteVaultAccessPolicy(ctx context.Context, in *resourcespb.DeleteVaultAccessPolicyRequest) (*commonpb.Empty, error) {
	return s.VaultAccessPolicyService.Service.Delete(ctx, in)
}

func (s *Server) CreateVirtualMachine(ctx context.Context, in *resourcespb.CreateVirtualMachineRequest) (*resourcespb.VirtualMachineResource, error) {
	return s.VirtualMachineService.Service.Create(ctx, in)
}
func (s *Server) ReadVirtualMachine(ctx context.Context, in *resourcespb.ReadVirtualMachineRequest) (*resourcespb.VirtualMachineResource, error) {
	return s.VirtualMachineService.Service.Read(ctx, in)
}
func (s *Server) UpdateVirtualMachine(ctx context.Context, in *resourcespb.UpdateVirtualMachineRequest) (*resourcespb.VirtualMachineResource, error) {
	return s.VirtualMachineService.Service.Update(ctx, in)
}
func (s *Server) DeleteVirtualMachine(ctx context.Context, in *resourcespb.DeleteVirtualMachineRequest) (*commonpb.Empty, error) {
	return s.VirtualMachineService.Service.Delete(ctx, in)
}

func (s *Server) RefreshState(ctx context.Context, _ *commonpb.Empty) (_ *commonpb.Empty, err error) {
	defer func() {
		if err != nil {
			go s.Database.AwsClient.UpdateErrorMetric("refresh", "refresh", errors.ErrorCode(err))
		}
	}()
	return errors.WrappingErrors(s.refresh)(ctx, &commonpb.Empty{})
}

func (s *Server) refresh(ctx context.Context, _ *commonpb.Empty) (*commonpb.Empty, error) {
	fmt.Println("refreshing state")
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

	_, err = deploy.EncodeAndStoreTfFile(ctx, c, nil, nil, true)
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
