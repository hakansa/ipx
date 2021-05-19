package ipx

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
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

// Well-known IPv4 addresses' decimal represntation
var (
	IPv4bcastInt = IPv4bcast.ToInt()
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
// Functions in this package accept
// either 4-byte (IPv4) or 16-byte (IPv6) slices as input.
type IP struct {
	net.IP
}

// IsV4 returns true if ip address is v4
func (i IP) IsV4() bool {
	i = i.To4()
	return i.IP != nil
}

// IsV6 returns true if ip address is v6
func (i IP) IsV6() bool {
	// check if ip is v4
	if i.IsV4() {
		return false
	}
	i = i.To16()
	return i.IP != nil
}

// ToInt returns the int reprenstation of IP
// If i is an IPv6 address, it returns zero
func (i IP) ToInt() uint32 {
	i = i.To4()

	if i.IP == nil {
		return uint32(0)
	}
	return binary.BigEndian.Uint32(i.IP)
}

// ToBigInt returns the bigint reprenstation of IP
func (i IP) ToBigInt() *big.Int {

	num := big.NewInt(0)
	// bigint for ipv4
	ip := i.To4()
	if ip.IP == nil {
		ip = i.To16()
	}
	num.SetBytes(ip.IP)

	return num
}

// ToBinary returns the binary reprenstation of IP
func (i IP) ToBinary() string {
	// ipv4
	ip := i.To4()
	if ip.IP != nil {
		return fmt.Sprintf("%08b%08b%08b%08b", ip.IP[0], ip.IP[1], ip.IP[2], ip.IP[3])
	}
	// ipv6
	// TODO: implement toBinary for ipv6
	return ""
}

// ToHex returns the hex reprenstation of IP
func (i IP) ToHex() string {
	// ipv4
	ip := i.To4()
	if ip.IP != nil {
		return fmt.Sprintf("%02X%02X%02X%02X", ip.IP[0], ip.IP[1], ip.IP[2], ip.IP[3])
	}
	// ipv6
	// TODO: implement toHex for ipv6
	return ""
}

// IsPrivate returns whether i is in a private network
func (i IP) IsPrivate() bool {
	for _, net := range privateNetworks {
		if net.Contains(i) {
			return true
		}
	}

	return false
}

// GetNext returns the next IP
func (i IP) GetNext() IP {
	// TODO: implement GetNext for IPv6
	return i.GetNextN(uint32(1))
}

// GetNextN returns the n'th next IP
func (i IP) GetNextN(n uint32) IP {
	// TODO: implement GetNextN for IPv6
	val := i.ToInt()
	val += n
	i = FromInt(val)
	return i
}

// GetAllNextN returns all IP's until n'th next IP
func (i IP) GetAllNextN(n uint32) []IP {
	var ipList []IP

	for j := 0; uint32(j) < n; j++ {
		ipList = append(ipList, i.GetNextN(uint32(j+1)))
	}
	return ipList
}

// GetPrevious returns the previous IP
func (i IP) GetPrevious() IP {
	// TODO: implement GetPrevious for IPv6
	return i.GetPreviousN(uint32(1))
}

// GetPreviousN returns the n'th next IP
func (i IP) GetPreviousN(n uint32) IP {
	// TODO: implement GetPreviousN for IPv6
	val := i.ToInt()

	if n > val {
		n = n - val - 1
		val = IPv4bcastInt
	}
	val -= n

	i = FromInt(val)
	return i
}

// GetAllPreviousN returns all IP's until n'th previous IP
func (i IP) GetAllPreviousN(n uint32) []IP {
	var ipList []IP

	for n > 0 {
		ipList = append(ipList, i.GetPreviousN(uint32(n)))
		n--
	}
	return ipList
}

// Equal reports whether ip and x are the same IP address.
// An IPv4 address and that same address in IPv6 form are
// considered to be equal.
func (i IP) Equal(x IP) bool {
	return i.IP.Equal(x.IP)
}

// Mask returns the result of masking the IP address ip with mask.
func (i IP) Mask(mask IPMask) IP {
	return IP{i.IP.Mask(mask.IPMask)}
}

// To4 converts the IPv4 address ip to a 4-byte representation.
// If ip is not an IPv4 address, To4 returns nil.
func (i IP) To4() IP {
	return IP{i.IP.To4()}
}

// To16 converts the IP address ip to a 16-byte representation.
// If ip is not an IP address (it is the wrong length), To16 returns nil.
func (i IP) To16() IP {
	return IP{i.IP.To16()}
}

// DefaultMask returns the default IP mask for the IP address ip.
// Only IPv4 addresses have default masks; DefaultMask returns
// nil if ip is not a valid IPv4 address.
func (i IP) DefaultMask() IPMask {
	return IPMask{i.IP.DefaultMask()}
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

// MustParseIP parses s as an IP address, returning the result.
// The string s can be in IPv4 dotted decimal ("192.0.2.1"),
// IPv6 ("2001:db8::68"), or IPv4-mapped IPv6 ("::ffff:192.0.2.1") form.
// If s is not a valid textual representation of an IP address,
// MustParseIP throws a panic
func MustParseIP(s string) IP {
	ip := net.ParseIP(s)

	if ip == nil {
		panic(ErrInvalidIP)
	}
	return IP{ip}
}

// ParseIP parses s as an IP address as exactly as MustParseIP.
// Differently, ParseIP returns ErrInvalidIP if s is not valid
// textual reprenstation of an IP address
func ParseIP(s string) (IP, error) {
	ip := net.ParseIP(s)

	if ip == nil {
		return IP{}, ErrInvalidIP
	}

	return IP{ip}, nil
}

// FromInt parses i as an IP address
// Receieved from https://gist.github.com/ammario/649d4c0da650162efd404af23e25b86b
func FromInt(i uint32) IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, i)
	return IP{ip}
}

// RandomIPv4 returns a random IPv4 address
func RandomIPv4() IP {
	return FromInt(uint32(rand.Intn(int(IPv4bcastInt))))
}

// ipEmptyString returns an empty string when ip is unset.
func ipEmptyString(ip net.IP) string {
	if len(ip) == 0 {
		return ""
	}
	return ip.String()
}
