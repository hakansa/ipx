package ipx

import (
	"errors"
	"net"
)

// Error definitions
var (
	ErrInvalidIP = errors.New("invalid ip address")
)

// IP address lengths (bytes).
const (
	IPv4len = 4
	IPv6len = 16
)

// Well-known IPv4 addresses
var (
	IPv4bcast     = IPv4(255, 255, 255, 255) // limited broadcast
	IPv4allsys    = IPv4(224, 0, 0, 1)       // all systems
	IPv4allrouter = IPv4(224, 0, 0, 2)       // all routers
	IPv4zero      = IPv4(0, 0, 0, 0)         // all zeros
)

// Well-known IPv6 addresses
var (
	IPv6zero                   = IP{net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}
	IPv6unspecified            = IP{net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}
	IPv6loopback               = IP{net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}}
	IPv6interfacelocalallnodes = IP{net.IP{0xff, 0x01, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x01}}
	IPv6linklocalallnodes      = IP{net.IP{0xff, 0x02, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x01}}
	IPv6linklocalallrouters    = IP{net.IP{0xff, 0x02, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x02}}
)

// IP is a single IP address, a slice of bytes.
// Functions in this package accept either 4-byte (IPv4) or 16-byte (IPv6) slices as input.
type IP struct {
	net.IP
}

// IsV4 returns true if ip address is v4
func (i *IP) IsV4() bool {
	return len(i.IP) == IPv4len
}

// IsV6 returns true if ip address is v6
func (i *IP) IsV6() bool {
	return len(i.IP) == IPv6len
}

// IPv4 returns the IP address (in 16-byte form) of the IPv4 address a.b.c.d.
func IPv4(a, b, c, d byte) IP {
	ip := net.IPv4(a, b, c, d)
	return IP{ip}
}

// LookupIP looks up host using the local resolver.
// It returns a slice of that host's IPv4 and IPv6 addresses.
func LookupIP(host string) ([]IP, error) {
	// array for returning ip's type of ipx.IP
	var ipList []IP

	netIPList, err := net.LookupIP(host)

	for _, netIP := range netIPList {
		ipList = append(ipList, IP{netIP})
	}

	return ipList, err
}

// ParseIP parses s as an IP address, returning the result.
// The string s can be in IPv4 dotted decimal ("192.0.2.1"),
// IPv6 ("2001:db8::68"), or IPv4-mapped IPv6 ("::ffff:192.0.2.1") form.
// If s is not a valid textual representation of an IP address, ParseIP returns nil.
func ParseIP(s string) IP {
	ip := net.ParseIP(s)

	return IP{ip}
}

// MustParseIP parses s as an IP address as exactly as ParseIP.
// Differently, MustParseIP returns ErrInvalidIP if s is not valid
// textual reprenstation of an IP address
func MustParseIP(s string) (IP, error) {
	ip := net.ParseIP(s)

	if ip == nil {
		return IP{}, ErrInvalidIP
	}

	return IP{ip}, nil
}

// ipEmptyString returns an empty string when ip is unset.
func ipEmptyString(ip net.IP) string {
	if len(ip) == 0 {
		return ""
	}
	return ip.String()
}
