package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-logr/stdr"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/llmariner/api-usage/pkg/sender"
	v1 "github.com/llmariner/cluster-manager/api/v1"
	"github.com/llmariner/cluster-manager/server/internal/config"
	"github.com/llmariner/cluster-manager/server/internal/k8s"
	"github.com/llmariner/cluster-manager/server/internal/server"
	"github.com/llmariner/cluster-manager/server/internal/store"
	"github.com/llmariner/common/pkg/db"
	"github.com/llmariner/rbac-manager/pkg/auth"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/protobuf/encoding/protojson"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func runCmd() *cobra.Command {
	var path string
	var logLevel int
	cmd := &cobra.Command{
		Use:   "run",
		Short: "run",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.Parse(path)
			if err != nil {
				return err
			}
			if err := c.Validate(); err != nil {
				return err
			}
			stdr.SetVerbosity(logLevel)
			if err := run(cmd.Context(), &c); err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&path, "config", "", "Path to the config file")
	cmd.Flags().IntVar(&logLevel, "v", 0, "Log level")
	_ = cmd.MarkFlagRequired("config")
	return cmd
}

func run(ctx context.Context, c *config.Config) error {
	logger := stdr.New(log.Default())
	log := logger.WithName("boot")

	var dbInst *gorm.DB
	var err error
	if c.Debug.Standalone {
		dbInst, err = gorm.Open(sqlite.Open(c.Debug.SqlitePath), &gorm.Config{})
	} else {
		dbInst, err = db.OpenDB(c.Database)
	}
	if err != nil {
		return err
	}

	st := store.New(dbInst)
	if err := st.AutoMigrate(); err != nil {
		return err
	}

	k8sClientFactory := k8s.NewClientFactory(c.SessionManagerServerEndpoint)

	addr := fmt.Sprintf("localhost:%d", c.GRPCPort)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	conn, err := grpc.NewClient(addr, opts...)
	if err != nil {
		return err
	}
	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			// Do not use the camel case for JSON fields to follow OpenAI API.
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames:     true,
				EmitDefaultValues: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
		runtime.WithIncomingHeaderMatcher(auth.HeaderMatcher),
		runtime.WithHealthzEndpoint(grpc_health_v1.NewHealthClient(conn)),
	)
	if err := v1.RegisterClustersServiceHandlerFromEndpoint(ctx, mux, addr, opts); err != nil {
		return err
	}

	s := server.New(st, k8sClientFactory, c.NVIDIA, c.ComponentStatusTimeout, logger)

	var usageSetter sender.UsageSetter
	if c.UsageSender.Enable {
		usage, err := sender.New(ctx, c.UsageSender, grpc.WithTransportCredentials(insecure.NewCredentials()), logger)
		if err != nil {
			return err
		}
		go func() { usage.Run(ctx) }()
		usageSetter = usage
	} else {
		usageSetter = sender.NoopUsageSetter{}
	}

	errCh := make(chan error)
	go func() {
		log.Info("Starting HTTP server...", "port", c.HTTPPort)
		errCh <- http.ListenAndServe(fmt.Sprintf(":%d", c.HTTPPort), mux)
	}()

	go func() {
		errCh <- s.Run(ctx, c.GRPCPort, c.AuthConfig, usageSetter)
	}()

	go func() {
		s := server.NewWorkerServiceServer(st, logger)
		errCh <- s.Run(ctx, c.WorkerServiceGRPCPort, c.AuthConfig)
	}()

	go func() {
		s := server.NewInternal(st, logger)
		errCh <- s.Run(c.InternalGRPCPort)
	}()

	if dc := c.DefaultCluster; dc.Name != "" {
		if err := s.CreateDefaultCluster(&dc); err != nil {
			return err
		}
	}

	return <-errCh
}
