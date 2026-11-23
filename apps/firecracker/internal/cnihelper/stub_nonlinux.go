//go:build !linux

package cnihelper

import (
	"context"
	"net"
)

const SnapshotGuestIP = "192.168.127.8"

type Network struct {
	TapDevice string
	GuestIP   net.IP
	MacAddr   string
	NetNS     string
	cleanup   []func() error
}

func (n *Network) Cleanup() {}

type Config struct {
	NetworkName string
	ConfDir     string
	BinPath     string
	GuestIP     string
}

func RefreshGuestNAT() {}

func RepairGuestEgress(confDir, networkName string) error {
	return nil
}

func PrepareCNIEnvironment(confDir, networkName string) error {
	return nil
}

func FlushGuestConntrack() {}

func Add(ctx context.Context, containerID string, cfg Config) (*Network, error) {
	return &Network{
		TapDevice: "tap0",
		GuestIP:   net.ParseIP(SnapshotGuestIP),
		MacAddr:   "06:00:c0:a8:7f:08",
		NetNS:     "/var/run/netns/" + containerID,
	}, nil
}

func Delete(ctx context.Context, containerID string, cfg Config) error {
	return nil
}

func SetLinkUp(device string) error {
	return nil
}

func SetLinkUpInNetNS(netns, device string) error {
	return nil
}

func GuestIPArgs(guestIP string) [][2]string {
	return [][2]string{{"IP", guestIP}}
}

func CleanupStaleAllocations(networkName string) error {
	return nil
}
