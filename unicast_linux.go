//go:build linux

package zeroconf

import (
	"context"
	"net"
	"syscall"
)

// listenUnicastUDP4 binds addr next to the 0.0.0.0:5353 sockets of the other
// mDNS stacks on the host. Linux only allows that when every socket on the
// port has SO_REUSEADDR set, which avahi and Go's multicast listeners do.
func listenUnicastUDP4(addr *net.UDPAddr) (*net.UDPConn, error) {
	lc := net.ListenConfig{Control: func(network, address string, c syscall.RawConn) error {
		var opErr error
		err := c.Control(func(fd uintptr) {
			opErr = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
		})
		if err != nil {
			return err
		}
		return opErr
	}}
	conn, err := lc.ListenPacket(context.Background(), "udp4", addr.String())
	if err != nil {
		return nil, err
	}
	return conn.(*net.UDPConn), nil
}
