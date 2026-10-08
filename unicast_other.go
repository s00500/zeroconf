//go:build !linux

package zeroconf

import (
	"errors"
	"net"
)

// listenUnicastUDP4 is Linux only: elsewhere the system responder (e.g.
// mDNSResponder) owns direct queries to the host anyway.
func listenUnicastUDP4(addr *net.UDPAddr) (*net.UDPConn, error) {
	return nil, errors.New("unicast listener not supported on this platform")
}
