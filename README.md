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

	// IsV4 returns true for ipv4 addresses
	ip.IsV4() // true

	// IsV6 returns true for ipv6 addresses
	ip.IsV6() // false 

	// ToInt returns the decimal representation of ip
	ip.ToInt() // 2886733825

	// ToBinary returns the binary representation of ip
	ip.ToBinary() // 10101100000100000001000000000001

	// ToHex returns the hex representation of ip
	ip.ToHex() // AC101001

	// IsPrivate returns true if ip is in a private network
	ip.IsPrivate() // true

	// GetNext returns next IP
	ip = ip.GetNext() // returns 172.16.16.2

	// GetNextN returns n'th next IP
	ip = ip.GetNextN(uint32(10)) // returns 172.16.16.12

	
	// All the other methods that the net package provides can be used with ipx
	// ip.DefaultMask()
	// ip.Equal(x)
	// ip.IsGlobalUnicast()
	// ip.IsInterfaceLocalMulticast()
	// ip.IsLinkLocalMulticast()
	// ip.IsLinkLocalUnicast()
	// ip.IsLoopback()
	// ip.IsMulticast()
	// ip.IsUnspecified()
	// ip.MarshalText()
	// ip.Mask(IPMask{})
	// ip.String()
	// ip.To4()
	// ip.To16()
	// ip.UnmarshalText()
}

```

## Other Examples

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

	// IsPrivate returns true if the IP is in a private network
	fmt.Printf("Is %v Private: %v\n", ip.String(), ip.IsPrivate()) // true

	// IPv4 returns the IP address (in 16-byte form) of the IPv4 address 10.99.99.1
	ip = ipx.IPv4(10, 99, 99, 1)

	// ParseCIDR parses a string in CIDR notation
	_, ipNet, _ := ipx.ParseCIDR("10.99.99.0/24")

	// Containts returns true if the given IP is in the IP network
	if ipNet.Contains(ip) {
		fmt.Printf("%v is in %v network\n", ip.String(), ipNet.String())
	}

	_, ipNet2, _ := ipx.ParseCIDR("10.99.98.0/23")

	// Intersects returns true if ip networks intersects with each other
	if ipNet.Intersects(*ipNet2) {
		fmt.Printf("%v network is intersects with %v network\n", ipNet.String(), ipNet2.String())
	}

}
```
