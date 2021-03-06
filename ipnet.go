package ipx

import (
	"math/rand"
	"net"
	"strconv"
)

var privateNetworks = []*IPNet{
	MustParseCIDR("10.0.0.0/8"),         // RFC1918
	MustParseCIDR("172.16.0.0/12"),      // private
	MustParseCIDR("192.168.0.0/16"),     // private
	MustParseCIDR("127.0.0.0/8"),        // RFC5735
	MustParseCIDR("0.0.0.0/8"),          // RFC1122 Section 3.2.1.3
	MustParseCIDR("169.254.0.0/16"),     // RFC3927
	MustParseCIDR("192.0.0.0/24"),       // RFC 5736
	MustParseCIDR("192.0.2.0/24"),       // RFC 5737
	MustParseCIDR("198.51.100.0/24"),    // Assigned as TEST-NET-2
	MustParseCIDR("203.0.113.0/24"),     // Assigned as TEST-NET-3
	MustParseCIDR("192.88.99.0/24"),     // RFC 3068
	MustParseCIDR("192.18.0.0/15"),      // RFC 2544
	MustParseCIDR("224.0.0.0/4"),        // RFC 3171
	MustParseCIDR("240.0.0.0/4"),        // RFC 1112
	MustParseCIDR("255.255.255.255/32"), // RFC 919 Section 7
	MustParseCIDR("100.64.0.0/10"),      // RFC 6598
	MustParseCIDR("::/128"),             // RFC 4291: Unspecified Address
	MustParseCIDR("::1/128"),            // RFC 4291: Loopback Address
	MustParseCIDR("100::/64"),           // RFC 6666: Discard Address Block
	MustParseCIDR("2001::/23"),          // RFC 2928: IETF Protocol Assignments
	MustParseCIDR("2001:2::/48"),        // RFC 5180: Benchmarking
	MustParseCIDR("2001:db8::/32"),      // RFC 3849: Documentation
	MustParseCIDR("2001::/32"),          // RFC 4380: TEREDO
	MustParseCIDR("fc00::/7"),           // RFC 4193: Unique-Local
	MustParseCIDR("fe80::/10"),          // RFC 4291: Section 2.5.6 Link-Scoped Unicast
	MustParseCIDR("ff00::/8"),           // RFC 4291: Section 2.7
	MustParseCIDR("2002::/16"),          // RFC 7526: 6to4 anycast prefix deprecated
}

// IPNet represents an IP network.
type IPNet struct {
	IP   IP     // network number
	Mask IPMask // network mask
}

// ParseCIDR parses s as a CIDR notation IP address and prefix length,
// like "192.0.2.0/24" or "2001:db8::/32", as defined in
// RFC 4632 and RFC 4291.
//
// It returns the IP address and the network implied by the IP and
// prefix length.
// For example, ParseCIDR("192.0.2.1/24") returns the IP address
// 192.0.2.1 and the network 192.0.2.0/24.
func ParseCIDR(s string) (IP, *IPNet, error) {
	ip, ipNet, err := net.ParseCIDR(s)
	if err != nil {
		return IP{}, nil, err
	}

	return IP{ip}, &IPNet{IP{ipNet.IP}, IPMask{ipNet.Mask}}, nil
}

// MustParseCIDR parses s as a CIDR notation
// if an error ocurred, it throws a panic
func MustParseCIDR(s string) *IPNet {
	_, ipNet, err := ParseCIDR(s)
	if err != nil {
		panic(err)
	}

	return ipNet
}

// Contains reports whether the network includes ip.
func (n *IPNet) Contains(ip IP) bool {
	nn, m := networkNumberAndMask(n)
	if x := ip.To4(); x.IP != nil {
		ip = x
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

// IPNumber returns the number of ip addresses in the network
func (n *IPNet) IPNumber() int {
	return 1 << (32 - simpleMaskLength(n.Mask.IPMask))
}

// UsableIPNumber returns the number of usable ip addresses in the network
// Basically it excludes the network address and broadcast address
func (n *IPNet) UsableIPNumber() int {
	num := n.IPNumber()

	// return the exact network size for /31 and /32
	if num <= 2 {
		return num
	}

	// exclude network address and broadcast address
	return num - 2
}

// NetworkSize returns the network size
func (n *IPNet) NetworkSize() int {
	return simpleMaskLength(n.Mask.IPMask)
}

// FirstIP returns the first ip in the network
func (n *IPNet) FirstIP() IP {
	return n.IP
}

// FirstUsableIP returns the first usable ip in the network
func (n *IPNet) FirstUsableIP() IP {
	return n.IP.GetNext()
}

// LastIP returns the last ip in the network
func (n *IPNet) LastIP() IP {
	if n.IPNumber() == 1 {
		return n.IP
	}
	return n.IP.GetNextN(uint32(n.IPNumber() - 1))
}

// LastUsableIP returns the last usable ip in the network
// If n is a /31 or /32 network, returns the firstIP
func (n *IPNet) LastUsableIP() IP {
	if n.IPNumber() <= 2 {
		return n.IP
	}
	return n.IP.GetNextN(uint32(n.IPNumber() - 2))
}

// GetAllIP returns all ip addresses in network
func (n *IPNet) GetAllIP() []IP {
	var ipList []IP
	ip := n.FirstIP()
	for i := 0; i < n.IPNumber(); i++ {
		ipList = append(ipList, ip)
		ip = ip.GetNext()
	}

	return ipList
}

// GetAllUsableIP returns all usable (adressable) ip addresses in network
func (n *IPNet) GetAllUsableIP() []IP {
	var ipList []IP

	num := n.IPNumber()

	if num == 1 {
		ipList = append(ipList, n.FirstIP())
		return ipList
	} else if num == 2 {
		ipList = append(ipList, n.FirstIP(), n.FirstIP().GetNext())
		return ipList
	}

	ip := n.FirstIP().GetNext()
	for i := 0; i < n.IPNumber()-2; i++ {
		ipList = append(ipList, ip)
		ip = ip.GetNext()
	}

	return ipList
}

// RandomIP returns a random ip address in n network
func (n *IPNet) RandomIP() IP {
	return FromInt(uint32(rand.Intn(n.IPNumber())) + n.FirstIP().ToInt())
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

// Intersects whether the networks intersects the other network
func (n *IPNet) Intersects(n2 IPNet) bool {
	return n.Contains(n2.IP) || n2.Contains(n.IP)
}

func networkNumberAndMask(n *IPNet) (ip net.IP, m net.IPMask) {
	if ip = n.IP.To4().IP; ip == nil {
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
