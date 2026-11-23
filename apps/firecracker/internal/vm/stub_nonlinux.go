//go:build !linux

package vm

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/muqtadir09/devin/apps/firecracker/internal/config"
	"github.com/muqtadir09/devin/apps/firecracker/internal/snapshot"
)

type Instance struct {
	ID         string
	Name       string
	Runtime    string
	IP         net.IP
	RuntimeURL string
	Phase      string
	Message    string
}

func (i *Instance) HealthCheck(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, i.RuntimeURL+"/health", nil)
	if err != nil {
		return err
	}
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (i *Instance) Shutdown(ctx context.Context) error {
	return nil
}

type Launcher struct {
	cfg       config.Config
	snapshot  *snapshot.Store
	activeIDs func() map[string]struct{}
}

func NewLauncher(cfg config.Config, snapshotStore *snapshot.Store) *Launcher {
	return &Launcher{cfg: cfg, snapshot: snapshotStore}
}

func (l *Launcher) SetActiveIDs(fn func() map[string]struct{}) {
	l.activeIDs = fn
}

func (l *Launcher) Restore(ctx context.Context, vmID, name, runtime string, cpu int32, memory string) (*Instance, error) {
	return &Instance{
		ID:         vmID,
		Name:       name,
		Runtime:    runtime,
		RuntimeURL: l.cfg.RuntimeFallback,
		Phase:      "Running",
		Message:    "mock restore on non-linux",
	}, nil
}
