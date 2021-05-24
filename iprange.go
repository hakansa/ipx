package ipx

import "math/rand"

// IPRange represents an IP range
// Lower boundary is included and Upper boundary is not included
type IPRange struct {
	Lower IP // LowerIP
	Upper IP // UpperIP
}

// ParseIPRange parses x and y as an IPRange
// Order is not imported on input parameters
func ParseIPRange(x, y string) (*IPRange, error) {
	xIP, err := ParseIP(x)
	if err != nil {
		return nil, err
	}

	yIP, err := ParseIP(y)
	if err != nil {
		return nil, err
	}

	ipRange := NewIPRange(xIP, yIP)

	return ipRange, nil
}

// MustParseIPRange parses x and y as an IPRange
// It throws panic if x or y is not a valid IP address
// Order is not imported on input parameters
func MustParseIPRange(x, y string) *IPRange {
	ipRange, err := ParseIPRange(x, y)
	if err != nil {
		panic(err)
	}
	return ipRange
}

// NewIPRange creates a new IPRange with x and y
// Order is not important
func NewIPRange(x, y IP) *IPRange {
	if x.ToInt() < y.ToInt() {
		return &IPRange{x, y}
	}

	return &IPRange{y, x}
}

// Contains reports whether the IPRange includes ip
func (i *IPRange) Contains(ip IP) bool {
	ipInt := ip.ToInt()
	return ipInt >= i.Lower.ToInt() && ipInt < i.Upper.ToInt()
}

// IPNumber returns the number of ip addresses in IPRange
func (i *IPRange) IPNumber() int {
	return int(i.Upper.ToInt() - i.Lower.ToInt())
}

// FirstIP returns the first ip in IPRange
func (i *IPRange) FirstIP() IP {
	return i.Lower
}

// LastIP returns the last ip in IPRange
func (i *IPRange) LastIP() IP {
	if i.Lower.Equal(i.Upper) {
		return i.Lower
	}
	return i.Upper.GetPrevious()
}

// GetAllIP returns all IP's in IPRange
func (i *IPRange) GetAllIP() []IP {

	ipList := []IP{i.Lower}

	if i.Lower.Equal(i.Upper) {
		return ipList
	}

	ipList = append(ipList, i.Lower.GetAllNextN(uint32(i.IPNumber())-1)...)

	return ipList
}

// RandomIP returns a random ip address in IPRange
func (i *IPRange) RandomIP() IP {
	return FromInt(uint32(rand.Intn(i.IPNumber())) + i.Lower.ToInt())
}
