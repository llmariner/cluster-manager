package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	v1 "github.com/llm-operator/cluster-manager/api/v1"
	"github.com/llm-operator/cluster-manager/server/internal/config"
	"github.com/llm-operator/cluster-manager/server/internal/server"
	"github.com/llm-operator/cluster-manager/server/internal/store"
	"github.com/llm-operator/common/pkg/db"
	"github.com/llm-operator/rbac-manager/pkg/auth"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/protobuf/encoding/protojson"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const flagConfig = "config"

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := cmd.Flags().GetString(flagConfig)
		if err != nil {
			return err
		}

		c, err := config.Parse(path)
		if err != nil {
			return err
		}

		if err := c.Validate(); err != nil {
			return err
		}

		if err := run(cmd.Context(), &c); err != nil {
			return err
		}
		return nil
	},
}

func run(ctx context.Context, c *config.Config) error {
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

	s := server.New(st)

	errCh := make(chan error)
	go func() {
		log.Printf("Starting HTTP server on port %d", c.HTTPPort)
		errCh <- http.ListenAndServe(fmt.Sprintf(":%d", c.HTTPPort), mux)
	}()

	go func() {
		errCh <- s.Run(ctx, c.GRPCPort, c.AuthConfig)
	}()

	go func() {
		s := server.NewWorkerServiceServer(st)
		errCh <- s.Run(ctx, c.WorkerServiceGRPCPort, c.AuthConfig)
	}()

	go func() {
		s := server.NewInternal(st)
		errCh <- s.Run(c.InternalGRPCPort)
	}()

	if dc := c.DefaultCluster; dc.Name != "" {
		if err := s.CreateDefaultCluster(&dc); err != nil {
			return err
		}
	}

	return <-errCh
}

func init() {
	runCmd.Flags().StringP(flagConfig, "c", "", "Configuration file path")
	_ = runCmd.MarkFlagRequired(flagConfig)
}
