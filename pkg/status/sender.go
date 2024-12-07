package status

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	v1 "github.com/llmariner/cluster-manager/api/v1"
	"github.com/llmariner/rbac-manager/pkg/auth"
	"google.golang.org/grpc"
)

// componentStatusComposer is a function to compose a component status message.
type componentStatusComposer func(ctx context.Context) (*v1.ComponentStatus, error)

// New creates a new ComponentStatusSender.
func New(c Config, fn componentStatusComposer, opt grpc.DialOption, log logr.Logger) (*ComponentStatusSender, error) {
	cc, err := grpc.NewClient(c.ClusterManagerServerWorkerServiceAddr, opt)
	if err != nil {
		return nil, fmt.Errorf("create cluster manager server worker service client: %s", err)
	}

	return &ComponentStatusSender{
		client:       v1.NewClustersWorkerServiceClient(cc),
		name:         c.Name,
		fn:           fn,
		logger:       log.WithName("component-status-sender"),
		interval:     c.Interval,
		initialDelay: c.InitalDelay,
	}, nil
}

// ComponentStatusSender sends component status data to the cluster manager.
type ComponentStatusSender struct {
	client       v1.ClustersWorkerServiceClient
	name         string
	fn           componentStatusComposer
	logger       logr.Logger
	interval     time.Duration
	initialDelay time.Duration
}

// Run starts the component status sender.
func (s *ComponentStatusSender) Run(ctx context.Context) {
	s.logger.Info("Starting component status sender...", "name", s.name, "interval", s.interval)

	time.Sleep(s.initialDelay)
	if err := s.send(ctx); err != nil {
		s.logger.Error(err, "Failed to send component status message")
	}

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := s.send(ctx); err != nil {
				s.logger.Error(err, "Failed to send component status message")
			}
		case <-ctx.Done():
			s.logger.Info("Stopping component status sender...")
			return
		}
	}
}

func (s *ComponentStatusSender) send(ctx context.Context) error {
	status, err := s.fn(ctx)
	if err != nil {
		return fmt.Errorf("compose status message: %s", err)
	}
	req := &v1.UpdateComponentStatusRequest{
		Name:   s.name,
		Status: status,
	}
	ctx = auth.AppendWorkerAuthorization(ctx)
	if _, err := s.client.UpdateComponentStatus(ctx, req); err != nil {
		return fmt.Errorf("update component status: %s", err)
	}
	s.logger.Info("Sent UpdateComponentStatus message", "req", req)
	return nil
}
