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

Below is an example which shows some common use cases for ipx.

```go
package main

import (
    "fmt"
        
    "github.com/hakansa/ipx"
)

func main(){

    // ParseIP throws ErrInvalidIP if given ip is not valid
    ip, err := ipx.ParseIP("256.256.256.256")
    if err != nil {
        // invalid ip
    }

    // MustParseIP throws a panic if the given string is not a valid IP address
    ip = ipx.MustParseIP("172.16.16.1") 

    // IsV4 returns true if the ip is v4
    fmt.Printf("Is %v V4: %v\n", ip.String(), ip.IsV4()) // true

    // IsPrivate returns true if the ip is in a private network
    fmt.Printf("Is %v Private: %v\n", ip.String(), ip.IsPrivate()) // true

    // IPv4 returns the IP address (in 16-byte form) of the IPv4 address 255.255.224.0
    ip = ipx.IPv4(255, 255, 224, 0)


}

```
