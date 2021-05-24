package ipx

import (
	"net"
	"reflect"
	"testing"
)

var parseCIDRTests = []struct {
	in  string
	ip  IP
	net *IPNet
	err error
}{
	{"135.104.0.0/32", IPv4(135, 104, 0, 0), &IPNet{IP: IPv4(135, 104, 0, 0), Mask: IPv4Mask(255, 255, 255, 255)}, nil},
	{"0.0.0.0/24", IPv4(0, 0, 0, 0), &IPNet{IP: IPv4(0, 0, 0, 0), Mask: IPv4Mask(255, 255, 255, 0)}, nil},
	{"135.104.0.0/24", IPv4(135, 104, 0, 0), &IPNet{IP: IPv4(135, 104, 0, 0), Mask: IPv4Mask(255, 255, 255, 0)}, nil},
	{"135.104.0.1/32", IPv4(135, 104, 0, 1), &IPNet{IP: IPv4(135, 104, 0, 1), Mask: IPv4Mask(255, 255, 255, 255)}, nil},
	{"135.104.0.1/24", IPv4(135, 104, 0, 1), &IPNet{IP: IPv4(135, 104, 0, 0), Mask: IPv4Mask(255, 255, 255, 0)}, nil},
}

func TestParseCIDR(t *testing.T) {
	for _, tt := range parseCIDRTests {
		ip, netw, err := ParseCIDR(tt.in)
		if !reflect.DeepEqual(err, tt.err) {
			t.Errorf("ParseCIDR(%q) = %v, %v; want %v, %v", tt.in, ip, netw, tt.ip, tt.net)
		}
		if err == nil && (!tt.ip.Equal(ip) || !tt.net.IP.Equal(netw.IP) || !reflect.DeepEqual(netw.Mask, tt.net.Mask)) {
			t.Errorf("ParseCIDR(%q) = %v, {%v, %v}; want %v, {%v, %v}", tt.in, ip, netw.IP, netw.Mask, tt.ip, tt.net.IP, tt.net.Mask)
		}
	}
}

var ipNetContainsTests = []struct {
	ip  IP
	net *IPNet
	ok  bool
}{
	{IPv4(172, 16, 1, 1), &IPNet{IP: IPv4(172, 16, 0, 0), Mask: CIDRMask(12, 32)}, true},
	{IPv4(172, 24, 0, 1), &IPNet{IP: IPv4(172, 16, 0, 0), Mask: CIDRMask(13, 32)}, false},
	{IPv4(192, 168, 0, 3), &IPNet{IP: IPv4(192, 168, 0, 0), Mask: IPv4Mask(0, 0, 255, 252)}, true},
	{IPv4(192, 168, 0, 4), &IPNet{IP: IPv4(192, 168, 0, 0), Mask: IPv4Mask(0, 255, 0, 252)}, false},
	{MustParseIP("2001:db8:1:2::1"), &IPNet{IP: MustParseIP("2001:db8:1::"), Mask: CIDRMask(47, 128)}, true},
	{MustParseIP("2001:db8:1:2::1"), &IPNet{IP: MustParseIP("2001:db8:2::"), Mask: CIDRMask(47, 128)}, false},
	{MustParseIP("2001:db8:1:2::1"), &IPNet{IP: MustParseIP("2001:db8:1::"), Mask: MustParseCIDR("ffff:0:ffff::/32").Mask}, true},
	{MustParseIP("2001:db8:1:2::1"), &IPNet{IP: MustParseIP("2001:db8:1::"), Mask: MustParseCIDR("0:0:0:ffff::/128").Mask}, false},
}

func TestIPNetContains(t *testing.T) {
	for _, tt := range ipNetContainsTests {
		if ok := tt.net.Contains(tt.ip); ok != tt.ok {
			t.Errorf("IPNet(%v).Contains(%v) = %v, want %v", tt.net, tt.ip, ok, tt.ok)
		}
	}
}

var ipNetStringTests = []struct {
	in  *IPNet
	out string
}{
	{&IPNet{IP: IPv4(192, 168, 1, 0), Mask: CIDRMask(26, 32)}, "192.168.1.0/26"},
	{&IPNet{IP: IPv4(192, 168, 1, 0), Mask: IPv4Mask(255, 0, 255, 0)}, "192.168.1.0/ff00ff00"},
	{&IPNet{IP: MustParseIP("2001:db8::"), Mask: CIDRMask(55, 128)}, "2001:db8::/55"},
}

func TestIPNetString(t *testing.T) {
	for _, tt := range ipNetStringTests {
		if out := tt.in.String(); out != tt.out {
			t.Errorf("IPNet.String(%v) = %q, want %q", tt.in, out, tt.out)
		}
	}
}

var cidrMaskTests = []struct {
	ones int
	bits int
	out  IPMask
}{
	{0, 32, IPv4Mask(0, 0, 0, 0)},
	{12, 32, IPv4Mask(255, 240, 0, 0)},
	{24, 32, IPv4Mask(255, 255, 255, 0)},
	{32, 32, IPv4Mask(255, 255, 255, 255)},
	{0, 128, IPMask{net.IPMask{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}},
	{4, 128, IPMask{net.IPMask{0xf0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}},
	{48, 128, IPMask{net.IPMask{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}},
	{128, 128, IPMask{net.IPMask{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}}},
	{33, 32, IPMask{}},
	{32, 33, IPMask{}},
	{-1, 128, IPMask{}},
	{128, -1, IPMask{}},
}

func TestCIDRMask(t *testing.T) {
	for _, tt := range cidrMaskTests {
		if out := CIDRMask(tt.ones, tt.bits); !reflect.DeepEqual(out, tt.out) {
			t.Errorf("CIDRMask(%v, %v) = %v, want %v", tt.ones, tt.bits, out, tt.out)
		}
	}
}

var (
	v4addr         = IP{net.IP{192, 168, 0, 1}}
	v4mappedv6addr = IP{net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xff, 0xff, 192, 168, 0, 1}}
	v6addr         = IP{net.IP{0x20, 0x1, 0xd, 0xb8, 0, 0, 0, 0, 0, 0, 0x1, 0x23, 0, 0x12, 0, 0x1}}
	v4mask         = IPMask{net.IPMask{255, 255, 255, 0}}
	v4mappedv6mask = IPMask{net.IPMask{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 255, 255, 255, 0}}
	v6mask         = IPMask{net.IPMask{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0, 0, 0, 0, 0, 0, 0, 0}}
	badaddr        = IP{net.IP{192, 168, 0}}
	badmask        = IPMask{net.IPMask{255, 255, 0}}
	v4maskzero     = IPMask{net.IPMask{0, 0, 0, 0}}
)

var networkNumberAndMaskTests = []struct {
	in  IPNet
	out IPNet
}{
	{IPNet{IP: v4addr, Mask: v4mask}, IPNet{IP: v4addr, Mask: v4mask}},
	{IPNet{IP: v4addr, Mask: v4mappedv6mask}, IPNet{IP: v4addr, Mask: v4mask}},
	{IPNet{IP: v4mappedv6addr, Mask: v4mappedv6mask}, IPNet{IP: v4addr, Mask: v4mask}},
	{IPNet{IP: v4mappedv6addr, Mask: v6mask}, IPNet{IP: v4addr, Mask: v4maskzero}},
	{IPNet{IP: v4addr, Mask: v6mask}, IPNet{IP: v4addr, Mask: v4maskzero}},
	{IPNet{IP: v6addr, Mask: v6mask}, IPNet{IP: v6addr, Mask: v6mask}},
	{IPNet{IP: v6addr, Mask: v4mappedv6mask}, IPNet{IP: v6addr, Mask: v4mappedv6mask}},
	{in: IPNet{IP: v6addr, Mask: v4mask}},
	{in: IPNet{IP: v4addr, Mask: badmask}},
	{in: IPNet{IP: v4mappedv6addr, Mask: badmask}},
	{in: IPNet{IP: v6addr, Mask: badmask}},
	{in: IPNet{IP: badaddr, Mask: v4mask}},
	{in: IPNet{IP: badaddr, Mask: v4mappedv6mask}},
	{in: IPNet{IP: badaddr, Mask: v6mask}},
	{in: IPNet{IP: badaddr, Mask: badmask}},
}

func TestNetworkNumberAndMask(t *testing.T) {
	for _, tt := range networkNumberAndMaskTests {
		ip, m := networkNumberAndMask(&tt.in)
		out := &IPNet{IP: IP{ip}, Mask: IPMask{m}}
		if !reflect.DeepEqual(&tt.out, out) {
			t.Errorf("networkNumberAndMask(%v) = %v, want %v", tt.in, out, &tt.out)
		}
	}
}

var ipNumberTests = []struct {
	in  IPNet
	out int
}{
	{IPNet{IP: IPv4(135, 104, 0, 0), Mask: IPv4Mask(255, 255, 255, 255)}, 1},
	{IPNet{IP: IPv4(0, 0, 0, 0), Mask: IPv4Mask(255, 255, 254, 0)}, 512},
	{IPNet{IP: IPv4(135, 104, 0, 0), Mask: IPv4Mask(255, 255, 252, 0)}, 1024},
	{IPNet{IP: IPv4(135, 104, 0, 1), Mask: IPv4Mask(0, 0, 0, 0)}, 4294967296},
	{IPNet{IP: IPv4(135, 104, 0, 0), Mask: IPv4Mask(255, 255, 255, 0)}, 256},
}

func TestIPNumber(t *testing.T) {
	for _, tt := range ipNumberTests {
		out := tt.in.IPNumber()
		if out != tt.out {
			t.Errorf("IPNet.IPNumber(%v) = %v, want %v", tt.in, out, tt.out)
		}
	}
}

var usableIPNumberTests = []struct {
	in  IPNet
	out int
}{
	{IPNet{IP: IPv4(135, 104, 0, 0), Mask: IPv4Mask(255, 255, 255, 255)}, 1},
	{IPNet{IP: IPv4(135, 104, 0, 0), Mask: IPv4Mask(255, 255, 255, 254)}, 2},
	{IPNet{IP: IPv4(0, 0, 0, 0), Mask: IPv4Mask(255, 255, 254, 0)}, 510},
	{IPNet{IP: IPv4(135, 104, 0, 0), Mask: IPv4Mask(255, 255, 252, 0)}, 1022},
	{IPNet{IP: IPv4(135, 104, 0, 1), Mask: IPv4Mask(0, 0, 0, 0)}, 4294967294},
	{IPNet{IP: IPv4(135, 104, 0, 0), Mask: IPv4Mask(255, 255, 255, 0)}, 254},
}

func TestUsableIPNumber(t *testing.T) {
	for _, tt := range usableIPNumberTests {
		out := tt.in.UsableIPNumber()
		if out != tt.out {
			t.Errorf("IPNet.UsableIPNumber(%v) = %v, want %v", tt.in, out, tt.out)
		}
	}
}

var firstIPTests = []struct {
	in        IPNet
	out       IP
	outUsable IP
}{
	{
		*MustParseCIDR("172.16.16.0/24"),
		IPv4(172, 16, 16, 0),
		IPv4(172, 16, 16, 1),
	},
	{
		*MustParseCIDR("172.16.16.11/20"),
		IPv4(172, 16, 16, 0),
		IPv4(172, 16, 16, 1),
	},
	{
		*MustParseCIDR("192.168.97.24/26"),
		IPv4(192, 168, 97, 0),
		IPv4(192, 168, 97, 1),
	},
	{
		*MustParseCIDR("192.168.97.24/20"),
		IPv4(192, 168, 96, 0),
		IPv4(192, 168, 96, 1),
	},
}

func TestFirstIP(t *testing.T) {
	for _, tt := range firstIPTests {
		out := tt.in.FirstIP()
		if out.String() != tt.out.String() {
			t.Errorf("IPNet.FirstIP(%v) = %v, want %v", tt.in, out, tt.out)
		}
	}
}

func TestFirstUsableIP(t *testing.T) {
	for _, tt := range firstIPTests {
		out := tt.in.FirstUsableIP()
		if out.String() != tt.outUsable.String() {
			t.Errorf("IPNet.FirstUsableIP(%v) = %v, want %v", tt.in, out, tt.outUsable)
		}
	}
}

var lastIPTests = []struct {
	in        IPNet
	out       IP
	outUsable IP
}{
	{
		*MustParseCIDR("172.16.16.0/24"),
		IPv4(172, 16, 16, 255),
		IPv4(172, 16, 16, 254),
	},
	{
		*MustParseCIDR("172.16.16.11/20"),
		IPv4(172, 16, 31, 255),
		IPv4(172, 16, 31, 254),
	},
	{
		*MustParseCIDR("192.168.97.24/26"),
		IPv4(192, 168, 97, 63),
		IPv4(192, 168, 97, 62),
	},
	{
		*MustParseCIDR("192.168.97.24/20"),
		IPv4(192, 168, 111, 255),
		IPv4(192, 168, 111, 254),
	},
}

func TestLastIP(t *testing.T) {
	for _, tt := range lastIPTests {
		out := tt.in.LastIP()
		if out.String() != tt.out.String() {
			t.Errorf("IPNet.LastIP(%v) = %v, want %v", tt.in, out, tt.out)
		}
	}
}

func TestLastUsableIP(t *testing.T) {
	for _, tt := range lastIPTests {
		out := tt.in.LastUsableIP()
		if out.String() != tt.outUsable.String() {
			t.Errorf("IPNet.LastUsableIP(%v) = %v, want %v", tt.in, out, tt.outUsable)
		}
	}
}

func TestRandomIP(t *testing.T) {
	for _, tt := range lastIPTests {
		out := tt.in.RandomIP()
		if !tt.in.Contains(out) {
			t.Errorf("IPNet.RandomIP(%v) = %v is not in %v network", tt.in, out, tt.in)
		}
	}
}

var getAllIPTests = []struct {
	in        IPNet
	out       []IP
	outUsable []IP
}{
	{
		*MustParseCIDR("172.16.16.0/30"),
		[]IP{
			IPv4(172, 16, 16, 0),
			IPv4(172, 16, 16, 1),
			IPv4(172, 16, 16, 2),
			IPv4(172, 16, 16, 3),
		},
		[]IP{
			IPv4(172, 16, 16, 1),
			IPv4(172, 16, 16, 2),
		},
	},
	{
		*MustParseCIDR("172.16.16.0/31"),
		[]IP{
			IPv4(172, 16, 16, 0),
			IPv4(172, 16, 16, 1),
		},
		[]IP{
			IPv4(172, 16, 16, 0),
			IPv4(172, 16, 16, 1),
		},
	},
	{
		*MustParseCIDR("172.16.16.0/32"),
		[]IP{
			IPv4(172, 16, 16, 0),
		},
		[]IP{
			IPv4(172, 16, 16, 0),
		},
	},
}

func TestGetAllIP(t *testing.T) {
	for _, tt := range getAllIPTests {
		out := tt.in.GetAllIP()
		for i, outIP := range out {
			if !reflect.DeepEqual(outIP.To4(), tt.out[i].To4()) {
				t.Errorf("IPNet.GetAllIP(%v) = %v, want %v", tt.in, outIP, tt.out[i])
			}
		}

	}
}

func TestGetAllUsableIP(t *testing.T) {
	for _, tt := range getAllIPTests {
		out := tt.in.GetAllUsableIP()

		for i, outIP := range out {
			if !reflect.DeepEqual(outIP.To4(), tt.outUsable[i].To4()) {
				t.Errorf("IPNet.GetAllIP(%v) = %v, want %v", tt.in, outIP, tt.out[i])
			}
		}
	}
}

var ipNetIntersectsTests = []struct {
	net1 *IPNet
	net2 *IPNet
	out  bool
}{
	{
		MustParseCIDR("172.16.16.0/24"),
		MustParseCIDR("172.16.16.0/23"),
		true,
	},
	{
		MustParseCIDR("172.16.16.0/24"),
		MustParseCIDR("172.16.14.0/24"),
		false,
	},

	{
		MustParseCIDR("0.0.0.0/0"),
		MustParseCIDR("172.16.14.0/24"),
		true,
	},
}

func TestIPNetIntersects(t *testing.T) {
	for _, tt := range ipNetIntersectsTests {
		out := tt.net1.Intersects(*tt.net2)
		if out != tt.out {
			t.Errorf("IPNet.Intersects(%v)(%v) = %v, want %v", tt.net1, tt.net2, out, tt.out)
		}

	}
}
