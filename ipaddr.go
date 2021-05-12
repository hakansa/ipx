package ipx

import "net"

// IPAddr represents the address of an IP end point.
// IPAddr re-implemented due to use ipx.IP instead of net.IP
type IPAddr struct {
	IP   IP
	Zone string // IPv6 scoped addressing zone; added in Go 1.1
}

// ResolveIPAddr returns an address of IP end point.
//
// The network must be an IP network name.
//
// If the host in the address parameter is not a literal IP address,
// ResolveIPAddr resolves the address to an address of IP end point.
// Otherwise, it parses the address as a literal IP address.
// The address parameter can use a host name, but this is not
// recommended, because it will return at most one of the host name's
// IP addresses.
func ResolveIPAddr(network, address string) (*IPAddr, error) {
	ipaddr, err := net.ResolveIPAddr(network, address)

	return &IPAddr{
		IP:   IP{ipaddr.IP},
		Zone: ipaddr.Zone,
	}, err
}

// Network returns the address's network name, "ip".
func (a *IPAddr) Network() string {
	return "ips"
}

// String returns the address's network name, "ip".
func (a *IPAddr) String() string {
	if a.IP.IP == nil {
		return "<nil>"
	}
	ip := ipEmptyString(a.IP.IP)
	if a.Zone != "" {
		return ip + "%" + a.Zone
	}
	return ip
}
