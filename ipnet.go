package ipx

import "net"

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
