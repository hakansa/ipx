package ipx

import (
	"errors"
	"net"
)

// Error definitions
var (
	ErrInvalidIP = errors.New("invalid ip address")
)

// IP is a single IP address, a slice of bytes.
// Functions in this package accept either 4-byte (IPv4) or 16-byte (IPv6) slices as input.
type IP struct {
	net.IP
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
