package ipx

import (
	"reflect"
	"testing"
)

var parseIPRangeTests = []struct {
	inLower string
	inUpper string
	out     IPRange
	reverse bool
	err     error
}{
	{"172.16.16.1", "172.16.16.100", IPRange{IPv4(172, 16, 16, 1), IPv4(172, 16, 16, 100)}, false, nil},
	{"172.16.16.100", "172.16.16.1", IPRange{IPv4(172, 16, 16, 1), IPv4(17, 16, 16, 100)}, true, nil},
}

func TestParseIPRange(t *testing.T) {
	for _, tt := range parseIPRangeTests {
		out, err := ParseIPRange(tt.inLower, tt.inUpper)
		if !reflect.DeepEqual(err, tt.err) {
			t.Errorf("ParseIPRange(%v,%v) = %v, %v; want %v, %v", tt.inLower, tt.inUpper, out, err, tt.out, tt.err)
		}
		if tt.reverse {
			if !reflect.DeepEqual(tt.inLower, out.Upper.String()) || !reflect.DeepEqual(tt.inUpper, out.Lower.String()) {
				t.Errorf("ParseIPRange(%v,%v) = %v, %v; want %v, %v", tt.inLower, tt.inUpper, out, err, tt.out, tt.err)
			}
		} else {
			if !reflect.DeepEqual(tt.inLower, out.Lower.String()) || !reflect.DeepEqual(tt.inUpper, out.Upper.String()) {
				t.Errorf("ParseIPRange(%v,%v) = %v, %v; want %v, %v", tt.inLower, tt.inUpper, out, err, tt.out, tt.err)
			}
		}

		// Test MustParseIPRange
		out = MustParseIPRange(tt.inLower, tt.inUpper)

		if tt.reverse {
			if !reflect.DeepEqual(tt.inLower, out.Upper.String()) || !reflect.DeepEqual(tt.inUpper, out.Lower.String()) {
				t.Errorf("ParseIPRange(%v,%v) = %v, %v; want %v, %v", tt.inLower, tt.inUpper, out, err, tt.out, tt.err)
			}
		} else {
			if !reflect.DeepEqual(tt.inLower, out.Lower.String()) || !reflect.DeepEqual(tt.inUpper, out.Upper.String()) {
				t.Errorf("ParseIPRange(%v,%v) = %v, %v; want %v, %v", tt.inLower, tt.inUpper, out, err, tt.out, tt.err)
			}
		}

		// Test NewIPRange
		out = NewIPRange(MustParseIP(tt.inLower), MustParseIP(tt.inUpper))

		if tt.reverse {
			if !reflect.DeepEqual(tt.inLower, out.Upper.String()) || !reflect.DeepEqual(tt.inUpper, out.Lower.String()) {
				t.Errorf("ParseIPRange(%v,%v) = %v, %v; want %v, %v", tt.inLower, tt.inUpper, out, err, tt.out, tt.err)
			}
		} else {
			if !reflect.DeepEqual(tt.inLower, out.Lower.String()) || !reflect.DeepEqual(tt.inUpper, out.Upper.String()) {
				t.Errorf("ParseIPRange(%v,%v) = %v, %v; want %v, %v", tt.inLower, tt.inUpper, out, err, tt.out, tt.err)
			}
		}
	}
}

var ipRangeContainsTests = []struct {
	ip      IP
	ipRange *IPRange
	ok      bool
}{
	{IPv4(172, 16, 16, 1), &IPRange{Lower: IPv4(172, 16, 16, 0), Upper: IPv4(172, 16, 16, 100)}, true},
	{IPv4(172, 16, 16, 100), &IPRange{Lower: IPv4(172, 16, 16, 0), Upper: IPv4(172, 16, 16, 100)}, false},
	{IPv4(172, 16, 15, 254), &IPRange{Lower: IPv4(172, 16, 16, 0), Upper: IPv4(172, 16, 16, 100)}, false},
	{IPv4(172, 16, 16, 0), &IPRange{Lower: IPv4(172, 16, 16, 0), Upper: IPv4(172, 16, 16, 100)}, true},
}

func TestIPRangeContains(t *testing.T) {
	for _, tt := range ipRangeContainsTests {
		if ok := tt.ipRange.Contains(tt.ip); ok != tt.ok {
			t.Errorf("IPRange(%v).Contains(%v) = %v, want %v", tt.ipRange, tt.ip, ok, tt.ok)
		}
	}
}

var ipRangeIPNumberTests = []struct {
	in  *IPRange
	out int
}{
	{&IPRange{Lower: IPv4(172, 16, 16, 0), Upper: IPv4(172, 16, 16, 100)}, 100},
	{&IPRange{Lower: IPv4(172, 16, 15, 254), Upper: IPv4(172, 16, 16, 4)}, 6},
}

func TestIPRangeIPNumber(t *testing.T) {
	for _, tt := range ipRangeIPNumberTests {
		out := tt.in.IPNumber()
		if out != tt.out {
			t.Errorf("IPRange.IPNumber(%v) = %v, want %v", tt.in, out, tt.out)
		}
	}
}

var ipRangeFirstLastIPTests = []struct {
	in       *IPRange
	outFirst IP
	outLast  IP
}{
	{
		MustParseIPRange("172.16.16.0", "172.16.16.100"),
		IPv4(172, 16, 16, 0),
		IPv4(172, 16, 16, 99),
	},
	{
		MustParseIPRange("0.0.0.0", "10.10.10.10"),
		IPv4(0, 0, 0, 0),
		IPv4(10, 10, 10, 9),
	},
	{
		MustParseIPRange("192.168.100.200", "192.168.100.50"),
		IPv4(192, 168, 100, 50),
		IPv4(192, 168, 100, 199),
	},
}

func TestIPRangeFirstIP(t *testing.T) {
	for _, tt := range ipRangeFirstLastIPTests {
		out := tt.in.FirstIP()
		if out.String() != tt.outFirst.String() {
			t.Errorf("IPRange.FirstIP(%v) = %v, want %v", tt.in, out, tt.outFirst)
		}
	}
}

func TestIPRangeLastIP(t *testing.T) {
	for _, tt := range ipRangeFirstLastIPTests {
		out := tt.in.LastIP()
		if out.String() != tt.outLast.String() {
			t.Errorf("IPRange.LastIP(%v) = %v, want %v", tt.in, out, tt.outLast)
		}
	}
}

var getIPRangeAllIPTests = []struct {
	in  IPRange
	out []IP
}{
	{
		*MustParseIPRange("172.16.16.0", "172.16.16.5"),
		[]IP{
			IPv4(172, 16, 16, 0),
			IPv4(172, 16, 16, 1),
			IPv4(172, 16, 16, 2),
			IPv4(172, 16, 16, 3),
			IPv4(172, 16, 16, 4),
		},
	},
	{
		*MustParseIPRange("172.16.16.0", "172.16.16.0"),
		[]IP{
			IPv4(172, 16, 16, 0),
		},
	},
}

func TestGetIPRangeAllIP(t *testing.T) {
	for _, tt := range getIPRangeAllIPTests {
		out := tt.in.GetAllIP()
		for i, outIP := range out {
			if !reflect.DeepEqual(outIP.To4(), tt.out[i].To4()) {
				t.Errorf("IPRange.GetAllIP(%v) = %v, want %v", tt.in, outIP, tt.out[i])
			}
		}

	}
}
