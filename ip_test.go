package ipx

import (
	"bytes"
	"net"
	"reflect"
	"testing"
)

var parseIPTests = []struct {
	in  string
	out IP
}{
	{"127.0.1.2", IPv4(127, 0, 1, 2)},
	{"127.0.0.1", IPv4(127, 0, 0, 1)},
	{"127.001.002.003", IPv4(127, 1, 2, 3)},
	{"::ffff:127.1.2.3", IPv4(127, 1, 2, 3)},
	{"::ffff:127.001.002.003", IPv4(127, 1, 2, 3)},
	{"::ffff:7f01:0203", IPv4(127, 1, 2, 3)},
	{"0:0:0:0:0000:ffff:127.1.2.3", IPv4(127, 1, 2, 3)},
	{"0:0:0:0:000000:ffff:127.1.2.3", IPv4(127, 1, 2, 3)},
	{"0:0:0:0::ffff:127.1.2.3", IPv4(127, 1, 2, 3)},

	{"2001:4860:0:2001::68", IP{net.IP{0x20, 0x01, 0x48, 0x60, 0, 0, 0x20, 0x01, 0, 0, 0, 0, 0, 0, 0x00, 0x68}}},
	{"2001:4860:0000:2001:0000:0000:0000:0068", IP{net.IP{0x20, 0x01, 0x48, 0x60, 0, 0, 0x20, 0x01, 0, 0, 0, 0, 0, 0, 0x00, 0x68}}},

	{"-0.0.0.0", IP{}},
	{"0.-1.0.0", IP{}},
	{"0.0.-2.0", IP{}},
	{"0.0.0.-3", IP{}},
	{"127.0.0.256", IP{}},
	{"abc", IP{}},
	{"123:", IP{}},
	{"fe80::1%lo0", IP{}},
	{"fe80::1%911", IP{}},
	{"a1:a2:a3:a4::b1:b2:b3:b4", IP{}},
}

func TestParseIP(t *testing.T) {
	for _, tt := range parseIPTests {
		if out, _ := ParseIP(tt.in); !reflect.DeepEqual(out, tt.out) {
			t.Errorf("ParseIP(%q) = %v, want %v", tt.in, out, tt.out)
		}

		var out IP
		if err := out.UnmarshalText([]byte(tt.in)); !reflect.DeepEqual(out, tt.out) || (tt.out.IP == nil) != (err != nil) {
			t.Errorf("IP.UnmarshalText(%q) = %v, %v, want %v", tt.in, out, err, tt.out)
		}
	}
}

func BenchmarkParseIP(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tt := range parseIPTests {
			ParseIP(tt.in)
		}
	}
}

var ipStringTests = []*struct {
	in  IP     // see RFC 791 and RFC 4291
	str string // see RFC 791, RFC 4291 and RFC 5952
	byt []byte
	error
}{
	// IPv4 address
	{
		IP{net.IP{192, 0, 2, 1}},
		"192.0.2.1",
		[]byte("192.0.2.1"),
		nil,
	},
	{
		IP{net.IP{0, 0, 0, 0}},
		"0.0.0.0",
		[]byte("0.0.0.0"),
		nil,
	},

	// IPv4-mapped IPv6 address
	{
		IP{net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xff, 0xff, 192, 0, 2, 1}},
		"192.0.2.1",
		[]byte("192.0.2.1"),
		nil,
	},
	{
		IP{net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xff, 0xff, 0, 0, 0, 0}},
		"0.0.0.0",
		[]byte("0.0.0.0"),
		nil,
	},

	// IPv6 address
	{
		IP{net.IP{0x20, 0x1, 0xd, 0xb8, 0, 0, 0, 0, 0, 0, 0x1, 0x23, 0, 0x12, 0, 0x1}},
		"2001:db8::123:12:1",
		[]byte("2001:db8::123:12:1"),
		nil,
	},
	{
		IP{net.IP{0x20, 0x1, 0xd, 0xb8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x1}},
		"2001:db8::1",
		[]byte("2001:db8::1"),
		nil,
	},
	{
		IP{net.IP{0x20, 0x1, 0xd, 0xb8, 0, 0, 0, 0x1, 0, 0, 0, 0x1, 0, 0, 0, 0x1}},
		"2001:db8:0:1:0:1:0:1",
		[]byte("2001:db8:0:1:0:1:0:1"),
		nil,
	},
	{
		IP{net.IP{0x20, 0x1, 0xd, 0xb8, 0, 0x1, 0, 0, 0, 0x1, 0, 0, 0, 0x1, 0, 0}},
		"2001:db8:1:0:1:0:1:0",
		[]byte("2001:db8:1:0:1:0:1:0"),
		nil,
	},
	{
		IP{net.IP{0x20, 0x1, 0, 0, 0, 0, 0, 0, 0, 0x1, 0, 0, 0, 0, 0, 0x1}},
		"2001::1:0:0:1",
		[]byte("2001::1:0:0:1"),
		nil,
	},
	{
		IP{net.IP{0x20, 0x1, 0xd, 0xb8, 0, 0, 0, 0, 0, 0x1, 0, 0, 0, 0, 0, 0}},
		"2001:db8:0:0:1::",
		[]byte("2001:db8:0:0:1::"),
		nil,
	},
	{
		IP{net.IP{0x20, 0x1, 0xd, 0xb8, 0, 0, 0, 0, 0, 0x1, 0, 0, 0, 0, 0, 0x1}},
		"2001:db8::1:0:0:1",
		[]byte("2001:db8::1:0:0:1"),
		nil,
	},
	{
		IP{net.IP{0x20, 0x1, 0xd, 0xb8, 0, 0, 0, 0, 0, 0xa, 0, 0xb, 0, 0xc, 0, 0xd}},
		"2001:db8::a:b:c:d",
		[]byte("2001:db8::a:b:c:d"),
		nil,
	},
	{
		IPv6unspecified,
		"::",
		[]byte("::"),
		nil,
	},

	// IP wildcard equivalent address in Dial/Listen API
	{
		IP{},
		"<nil>",
		nil,
		nil,
	},
}

func TestIPString(t *testing.T) {
	for _, tt := range ipStringTests {
		if out := tt.in.String(); out != tt.str {
			t.Errorf("IP.String(%v) = %q, want %q", tt.in, out, tt.str)
		}
		if out, err := tt.in.MarshalText(); !bytes.Equal(out, tt.byt) || !reflect.DeepEqual(err, tt.error) {
			t.Errorf("IP.MarshalText(%v) = %v, %v, want %v, %v", tt.in, out, err, tt.byt, tt.error)
		}
	}
}

var sink string

func BenchmarkIPString(b *testing.B) {

	b.Run("IPv4", func(b *testing.B) {
		benchmarkIPString(b, IPv4len)
	})

	b.Run("IPv6", func(b *testing.B) {
		benchmarkIPString(b, IPv6len)
	})
}

func benchmarkIPString(b *testing.B, size int) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, tt := range ipStringTests {
			if tt.in.IP != nil && len(tt.in.IP) == size {
				sink = tt.in.String()
			}
		}
	}
}

var ipTypeTests = []*struct {
	in   IP
	isV4 bool
	isV6 bool
}{
	// IPv4 address
	{
		IP{net.IP{192, 0, 2, 1}},
		true,
		false,
	},
	{
		IP{net.IP{0, 0, 0, 0}},
		true,
		false,
	},

	// IPv4-mapped IPv6 address
	{
		IP{net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xff, 0xff, 192, 0, 2, 1}},
		false,
		true,
	},
	{
		IP{net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xff, 0xff, 0, 0, 0, 0}},
		false,
		true,
	},

	// IPv6 address
	{
		IP{net.IP{0x20, 0x1, 0xd, 0xb8, 0, 0, 0, 0, 0, 0, 0x1, 0x23, 0, 0x12, 0, 0x1}},
		false,
		true,
	},
	{
		IP{net.IP{0x20, 0x1, 0xd, 0xb8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x1}},
		false,
		true,
	},
	{
		IP{net.IP{0x20, 0x1, 0xd, 0xb8, 0, 0, 0, 0x1, 0, 0, 0, 0x1, 0, 0, 0, 0x1}},
		false,
		true,
	},
	{
		IP{net.IP{0x20, 0x1, 0xd, 0xb8, 0, 0x1, 0, 0, 0, 0x1, 0, 0, 0, 0x1, 0, 0}},
		false,
		true,
	},
	{
		IP{net.IP{0x20, 0x1, 0, 0, 0, 0, 0, 0, 0, 0x1, 0, 0, 0, 0, 0, 0x1}},
		false,
		true,
	},
	{
		IP{net.IP{0x20, 0x1, 0xd, 0xb8, 0, 0, 0, 0, 0, 0x1, 0, 0, 0, 0, 0, 0}},
		false,
		true,
	},
	{
		IP{net.IP{0x20, 0x1, 0xd, 0xb8, 0, 0, 0, 0, 0, 0x1, 0, 0, 0, 0, 0, 0x1}},
		false,
		true,
	},
	{
		IP{net.IP{0x20, 0x1, 0xd, 0xb8, 0, 0, 0, 0, 0, 0xa, 0, 0xb, 0, 0xc, 0, 0xd}},
		false,
		true,
	},
	{
		IPv6unspecified,
		false,
		true,
	},

	// IP wildcard equivalent address in Dial/Listen API
	{
		IP{},
		false,
		false,
	},
}

func TestIPTypes(t *testing.T) {
	for _, tt := range ipTypeTests {
		if out := tt.in.IsV4(); out != tt.isV4 {
			t.Errorf("IP.IsV4(%v) = %v, want %v", tt.in, out, tt.isV4)
		}
		if out := tt.in.IsV6(); out != tt.isV6 {
			t.Errorf("IP.IsV6(%v) = %v, want %v", tt.in, out, tt.isV6)
		}
	}
}

var ipEmptyStringTests = []*struct {
	in  IP
	out string
}{
	// IPv4 address
	{
		IP{net.IP{192, 0, 2, 1}},
		"192.0.2.1",
	},
	{
		IP{},
		"",
	},
}

func TestIPEmptyString(t *testing.T) {
	for _, tt := range ipEmptyStringTests {
		if out := ipEmptyString(tt.in.IP); out != tt.out {
			t.Errorf("ipEmptyString(%v) = %v, want %v", tt.in, out, tt.out)
		}
	}
}
