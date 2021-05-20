package ipx

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
