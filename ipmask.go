package ipx

import "net"

// An IPMask is a bitmask that can be used to manipulate
// IP addresses for IP addressing and routing.
//
// See type IPNet and func ParseCIDR for details.
type IPMask struct {
	net.IPMask
}

// CIDRMask returns an IPMask consisting of 'ones' 1 bits
// followed by 0s up to a total length of 'bits' bits.
// For a mask of this form, CIDRMask is the inverse of IPMask.Size.
func CIDRMask(ones, bits int) IPMask {
	ipmask := net.CIDRMask(ones, bits)
	return IPMask{ipmask}
}

// IPv4Mask returns the IP mask (in 4-byte form) of the
// IPv4 mask a.b.c.d.
func IPv4Mask(a, b, c, d byte) IPMask {
	ipmask := net.IPv4Mask(a, b, c, d)

	return IPMask{ipmask}
}
