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

	}
}
