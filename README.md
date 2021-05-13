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

    // MustParseIP throws a panic if given ip is not valid
    ip := ipx.MustParseIP("256.256.256.256") 


}

```
