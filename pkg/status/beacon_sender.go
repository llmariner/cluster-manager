package status

import (
	"context"

	"github.com/go-logr/logr"
	v1 "github.com/llmariner/cluster-manager/api/v1"
	"google.golang.org/grpc"
)

// NewBeaconSender creates a new BeaconSender for sending beacon messages.
func NewBeaconSender(
	c Config,
	opt grpc.DialOption,
	log logr.Logger,
) (*BeaconSender, error) {
	var err error
	bs := &BeaconSender{}
	bs.sender, err = New(c, bs.composeStatusMessage, opt, log)
	return bs, err
}

// BeaconSender sends beacon messages to the cluster manager.
type BeaconSender struct {
	sender *ComponentStatusSender
}

// Run starts the beacon sender.
func (s *BeaconSender) Run(ctx context.Context) {
	s.sender.Run(ctx)
}

func (s *BeaconSender) composeStatusMessage(ctx context.Context) (*v1.ComponentStatus, error) {
	return &v1.ComponentStatus{
		IsHealthy: true,
	}, nil
}
