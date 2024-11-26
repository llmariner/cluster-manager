package status

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	v1 "github.com/llmariner/cluster-manager/api/v1"
	"google.golang.org/grpc"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// NewPodStatusSender creates a new PodStatusSender for sending pod status.
func NewPodStatusSender(
	c Config,
	ns, label string,
	opt grpc.DialOption,
	log logr.Logger,
) (*PodStatusSender, error) {
	restConfig, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	kc, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("create kubernetes client: %s", err)
	}
	pss := &PodStatusSender{
		namespace: ns,
		label:     label,
		kc:        kc,
	}
	pss.sender, err = New(c, pss.composeStatusMessage, opt, log)
	return pss, err
}

// PodStatusSender sends pod status messages to the cluster manager.
type PodStatusSender struct {
	sender    *ComponentStatusSender
	namespace string
	label     string
	kc        *kubernetes.Clientset
}

// Run starts the pod status sender.
func (s *PodStatusSender) Run(ctx context.Context) {
	s.sender.Run(ctx)
}

func (s *PodStatusSender) composeStatusMessage(ctx context.Context) (*v1.ComponentStatus, error) {
	podList, err := s.kc.CoreV1().Pods(s.namespace).List(ctx, metav1.ListOptions{
		LabelSelector: s.label,
	})
	if err != nil {
		return nil, err
	}
	var msgs []string
	isHealthy := true
	for _, pod := range podList.Items {
		msgs = append(msgs, fmt.Sprintf("%s: %s", pod.Name, pod.Status.Phase))
		if pod.Status.Phase != "Running" {
			isHealthy = false
		}
	}
	return &v1.ComponentStatus{
		IsHealthy: isHealthy,
		Message:   strings.Join(msgs, "; "),
	}, nil
}
