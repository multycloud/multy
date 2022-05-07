package cmd

import (
	"context"
	"crypto/x509"
	"fmt"
	"github.com/multycloud/multy/api/proto"
	flag "github.com/spf13/pflag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"os"
	"strings"
)

// DeleteCommand - Deletes a resource by id.
type DeleteCommand struct {
	ApiKey         string
	ServerEndpoint string
	ResourceId     string
}

func (c *DeleteCommand) ParseFlags(f *flag.FlagSet, args []string) {
	f.StringVar(&c.ApiKey, "api_key", os.Getenv("MULTY_API_KEY"), "Multy api key. Can also be set with env var MULTY_API_KEY")
	f.StringVar(&c.ServerEndpoint, "endpoint", "api.multy.dev:443", "Endpoint where server is running at")
	_ = f.Parse(args)
	if len(f.Args()) < 1 {
		fmt.Println("resource_id argument is mandatory")
		f.Usage()
		os.Exit(1)
	}
	c.ResourceId = f.Args()[0]
}

func (c *DeleteCommand) Description() CommandDesc {
	return CommandDesc{
		Name:        "delete",
		Description: "deletes a resource by id, without destroying it in the underlying cloud provider",
		Usage:       "multy delete <resource_id> [options]",
	}
}

func (c *DeleteCommand) Execute(ctx context.Context) error {
	creds := insecure.NewCredentials()
	if !strings.HasPrefix(c.ServerEndpoint, "localhost") {
		cp, err := x509.SystemCertPool()
		if err != nil {
			return err
		}
		creds = credentials.NewClientTLSFromCert(cp, "")
	}

	conn, err := grpc.Dial(c.ServerEndpoint, grpc.WithTransportCredentials(creds))
	if err != nil {
		return err
	}

	client := proto.NewMultyResourceServiceClient(conn)
	if len(c.ApiKey) == 0 {
		fmt.Println("api_key must be set")
		os.Exit(1)
	}
	ctx = metadata.AppendToOutgoingContext(ctx, "api_key", c.ApiKey)
	_, err = client.DeleteResource(ctx, &proto.DeleteResourceRequest{ResourceId: c.ResourceId})
	if err != nil {
		if s, ok := status.FromError(err); ok && s.Code() != codes.Internal {
			fmt.Println(s.Message())
			os.Exit(1)
		}
		return err
	}

	fmt.Printf("resource %s was deleted\n", c.ResourceId)

	return nil
}
