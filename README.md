# ipx

[![Coverage Status](https://coveralls.io/repos/github/hakansa/ipx/badge.svg?branch=main)](https://coveralls.io/github/hakansa/ipx?branch=main) [![Go Report Card](https://goreportcard.com/badge/github.com/hakansa/ipx)](https://goreportcard.com/report/github.com/hakansa/ipx) [![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/hakansa/ipx) [![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/hakansa/ipx/master/LICENSE)

ipx is a library which provides a set of extensions on go's standart IP functions in `net` package.

## compability with net package
ipx is fully compatible with net package.
Also, it implements the necessary functions in net package such as ParseIP, ParseCIDR, CIDRMask etc.
Therefore, don't need to import net package additionally.

## install

    go get github.com/hakansa/ipx

## usage

The following examples shows some common use cases for ipx.

You can access the example in [The Go Playground](https://play.golang.org/p/Hlic8q3BQMw)

## IP 
```go
package main

import (
	"fmt"

	"github.com/hakansa/ipx"
)

func main() {
	// ParseIP throws ErrInvalidIP if given ip is not valid
	ip, err := ipx.ParseIP("256.256.256.256")
	if err != nil {
		// invalid ip
	}

	// MustParseIP throws a panic if the given string is not a valid IP address
	ip = ipx.MustParseIP("172.16.16.1")

	// IPv4 returns the IP address (in 16-byte form) of the IPv4 address a.b.c.d.
	ip = ipx.IPv4(172, 16, 16, 1)

	// IsV4 returns true for ipv4 addresses
	ip.IsV4() // true

	// IsV6 returns true for ipv6 addresses
	ip.IsV6() // false 

	// ToInt returns the decimal representation of ip
	// ToInt returns 0 for ipv6 addresses
	ip.ToInt() // 2886733825

	// ToBigInt returns the decimal representation of ip in math.BigInt format
	// ToBigInt can be used for ipv4 and ipv6 addresses
	ip.ToBigInt() // 2886733825

	// ToBinary returns the binary representation of ip
	ip.ToBinary() // 10101100000100000001000000000001

	// ToHex returns the hex representation of ip
	ip.ToHex() // AC101001

	// IsPrivate returns true if ip is in a private network
	ip.IsPrivate() // true

	// GetNext returns next IP
	ip.GetNext() // 172.16.16.2

	// GetNextN returns n'th next IP
	ip.GetNextN(uint32(10)) // 172.16.16.11

	// GetAllNextN returns an IP array
	// that contains all IP's until n'th next IP
	ip.GetAllNextN(uint32(3)) // []IP{ 172.16.16.2 , 172.16.16.3 , 172.16.16.4 }

	// GetPrevious returns previous IP
	ip.GetPrevious() // 172.16.16.0

	// GetPreviousN returns n'th previous IP
	ip.GetPreviousN(uint32(2)) // 172.16.15.255

	// GetAllPreviousN returns an IP array
	// that contains all IP's from n'th previous IP
	ip.GetAllPreviousN(uint32(3)) // []IP{ 172.16.15.254 , 172.16.15.255 , 172.16.16.0 }

	// FromInt returns IP address for given integer
	ip = ipx.FromInt(uint32(2886733825)) // 172.16.16.1

	// RandomIPv4 returns a random IPv4 address
	ip = ipx.RandomIPv4() // x.y.z.t

	
	// All the other methods that the net package provides can be used with ipx
	ip.DefaultMask()
	ip.Equal(ipx.IPv4(172, 16, 16, 1))
	ip.IsGlobalUnicast()
	ip.IsInterfaceLocalMulticast()
	ip.IsLinkLocalMulticast()
	ip.IsLinkLocalUnicast()
	ip.IsLoopback()
	ip.IsMulticast()
	ip.IsUnspecified()
	ip.MarshalText()
	ip.Mask(IPMask{})
	ip.String()
	ip.To4()
	ip.To16()
	ip.UnmarshalText()
}

```
## IPNet
```go
package main

import (
	"fmt"

	"github.com/hakansa/ipx"
)

func main() {

	// ParseCIDR parses a string in CIDR notation
	_, ipNet, _ := ipx.ParseCIDR("172.16.16.0/24")

	// MustParseIP throws a panic if the given string is not a valid IP Network
	ipNet = ipx.MustParseCIDR("172.16.16.0/24")

	// IPNumber returns the number of ip addresses in the network
	ipNet.IPNumber() // 256

	// UsableIPNumber returns the number of usable ip addresses in the network
	// Basically it excludes the network address and broadcast address
	ipNet.UsableIPNumber() // 254

	// NetworkSize returns the network size
	ipNet.NetworkSize() // 24

	// FirstIP returns the first IP in network 
	ipNet.FirstIP() // 172.16.16.0

	// FirstUsableIP returns the first usable (addressable) IP in network
	ipNet.FirstUsableIP() // 172.16.16.1

	// LastIP returns the list IP in network
	ipNet.LastIP() // 172.16.16.255

	// LastUsableIP returns the last usable (addressable) IP in network
	ipNet.LastUsableIP() // 172.16.16.254

	// GetAllIP returns all IP's in the network as an array
	ipNet.GetAllIP() // []IP{ 172.16.16.0, 172.16.16.1, ... , 172.16.16.255 }

	// GetAllUsableIP returns all usable IP's in the network as an array
	ipNet.GetAllUsableIP() // []IP{ 172.16.16.1, 172.16.16.2, ... , 172.16.16.254 }

	// RandomIP returns a random ip in the network
	ipNet.RandomIP() // 172.16.16.X

	// Intersects whether the networks intersects the other network
	ipNet.Intersects(ipx.MustParseCIDR("172.16.15.0/23")) // true

	// All the other methods that the net package provides can be used with ipx
	ipNet.Contains(ipx.IPv4(172, 16, 16, 23))
	ipNet.Network()
	ipNet.String()

}
```

## IPRange
```go
package main

import (
	"fmt"

	"github.com/hakansa/ipx"
)

func main() {

	// ParseIPRange parses x and y as an IPRange
	// Upper ip's boundary is excluded
	ipRange, _ := ipx.ParseIPRange("172.16.16.0", "172.16.16.100")

	// ParseIPRange throws a panic if the given strings is not valid IP addresses
	ipRange = ipx.MustParseIPRange("172.16.16.0", "172.16.16.100")

	// NewIPRange creates a new IPRange
	ipRange = ipx.NewIPRange(ipx.IPv4(172, 16, 16, 0), ipx.IPv4(172, 16, 16, 100))

	// Order is not important when creating IPRange
	ipRange = ipx.MustParseIPRange("172.16.16.100", "172.16.16.0")

	// Contains checks if ip is in range
	ipRange.Contains(ipx.IPv4(172, 16, 16, 75)) // true

	// IPNumber returns the number of ip addresses in IPRange
	ipRange.IPNumber() // 100

	// FirstIP returns the first IP in IPRange 
	ipRange.FirstIP() // 172.16.16.0

	// LastIP returns the list IP in IPRange
	ipRange.LastIP() // 172.16.16.99

	// GetAllIP returns all IP's in IPRange as an array
	ipRange.GetAllIP() // []IP{ 172.16.16.0, 172.16.16.1, ... , 172.16.16.99 }

	// RandomIP returns a random ip in IPRange
	ipRange.RandomIP() // 172.16.16.X (0 <= X < 100)

	// Intersects whether the IPRange intersects other IPRange
	ipRange.Intersects(ipx.MustParseIPRange("172.16.16.50", "172.16.16.150")) // true
}

```