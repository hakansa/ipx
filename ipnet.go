package ipx

import (
	"net"
	"strconv"
)

// IPNet represents an IP network.
type IPNet struct {
	IP   IP     // network number
	Mask IPMask // network mask
}

// Contains reports whether the network includes ip.
func (n *IPNet) Contains(ip IP) bool {
	nn, m := networkNumberAndMask(n)
	if x := ip.To4(); x != nil {
		ip.IP = x
	}
	l := len(ip.IP)
	if l != len(nn) {
		return false
	}
	for i := 0; i < l; i++ {
		if nn[i]&m[i] != ip.IP[i]&m[i] {
			return false
		}
	}
	return true
}

// Network returns the address's network name, "ip+net".
func (n *IPNet) Network() string { return "ip+net" }

// String returns the CIDR notation of n like "192.0.2.0/24"
// or "2001:db8::/48" as defined in RFC 4632 and RFC 4291.
// If the mask is not in the canonical form, it returns the
// string which consists of an IP address, followed by a slash
// character and a mask expressed as hexadecimal form with no
// punctuation like "198.51.100.0/c000ff00".
func (n *IPNet) String() string {
	nn, m := networkNumberAndMask(n)
	if nn == nil || m == nil {
		return "<nil>"
	}
	l := simpleMaskLength(m)
	if l == -1 {
		return nn.String() + "/" + m.String()
	}
	return nn.String() + "/" + strconv.FormatUint(uint64(l), 10)
}

func networkNumberAndMask(n *IPNet) (ip net.IP, m net.IPMask) {
	if ip = n.IP.To4(); ip == nil {
		ip = n.IP.IP
		if len(ip) != net.IPv6len {
			return nil, nil
		}
	}
	m = n.Mask.IPMask
	switch len(m) {
	case IPv4len:
		if len(ip) != IPv4len {
			return nil, nil
		}
	case IPv6len:
		if len(ip) == IPv4len {
			m = m[12:]
		}
	default:
		return nil, nil
	}
	return
}

// If mask is a sequence of 1 bits followed by 0 bits,
// return the number of 1 bits.
func simpleMaskLength(mask net.IPMask) int {
	var n int
	for i, v := range mask {
		if v == 0xff {
			n += 8
			continue
		}
		// found non-ff byte
		// count 1 bits
		for v&0x80 != 0 {
			n++
			v <<= 1
		}
		// rest must be 0 bits
		if v != 0 {
			return -1
		}
		for i++; i < len(mask); i++ {
			if mask[i] != 0 {
				return -1
			}
		}
		break
	}
	return n
}
