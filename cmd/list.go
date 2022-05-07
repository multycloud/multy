package cmd

import (
	"context"
	"crypto/x509"
	"fmt"
	"github.com/multycloud/multy/api/proto"
	"github.com/multycloud/multy/api/proto/commonpb"
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

// ListCommand - Lists all resource ids.
type ListCommand struct {
	ApiKey         string
	ServerEndpoint string
}

func (c *ListCommand) ParseFlags(f *flag.FlagSet, args []string) {
	f.StringVar(&c.ApiKey, "api_key", os.Getenv("MULTY_API_KEY"), "Multy api key. Can also be set with env var MULTY_API_KEY")
	f.StringVar(&c.ServerEndpoint, "endpoint", "api.multy.dev:443", "Endpoint where server is running at")
	_ = f.Parse(args)
}

func (c *ListCommand) Description() CommandDesc {
	return CommandDesc{
		Name:        "list",
		Description: "lists all resources",
		Usage:       "multy serve [options]",
	}
}

func (c *ListCommand) Execute(ctx context.Context) error {
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
	resources, err := client.ListResources(ctx, &commonpb.Empty{})
	if err != nil {
		if s, ok := status.FromError(err); ok && s.Code() != codes.Internal {
			fmt.Println(s.Message())
			os.Exit(1)
		}
		return err
	}

	format := "%-24s\t%s\n"
	fmt.Printf(format, "Resource type", "Resource id")
	for _, r := range resources.Resources {
		fmt.Printf(format, r.ResourceType, r.ResourceId)
	}

	return nil
}
