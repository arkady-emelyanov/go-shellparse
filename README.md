# Parse strings à la shell

[![GoDoc](https://godoc.org/github.com/arkady-emelyanov/go-shellparse?status.svg)](https://godoc.org/github.com/arkady-emelyanov/go-shellparse)
[![Build Status](https://travis-ci.org/arkady-emelyanov/go-shellparse.svg?branch=master)](https://travis-ci.org/arkady-emelyanov/go-shellparse)
[![Go Report Card](https://goreportcard.com/badge/github.com/arkady-emelyanov/go-shellparse)](https://goreportcard.com/report/github.com/arkady-emelyanov/go-shellparse)
[![codecov](https://codecov.io/gh/arkady-emelyanov/go-shellparse/branch/master/graph/badge.svg)](https://codecov.io/gh/arkady-emelyanov/go-shellparse)

## Features

* No dependencies
* Ability to parse complex strings
* Expand variables `${VAR}` with custom provided k/v map
* Useful helpers to parse strings into
    * slices
    * maps
    * command and arguments

## Sample

Program
```
package main

import (
	"fmt"

	"github.com/arkady-emelyanov/go-shellparse"
)

func main() {
	fmt.Println(">>> ParseCommand ")
	parseCommandSimple()

	fmt.Println("")
	fmt.Println("")
	fmt.Println(">>> ParseCommandWithEnv")
	parseCommandWithEnv()
}

func parseCommandSimple() {
	cmd := `bash -c 'echo "it\'s complex command" && sleep 3 && exit 1'`
	bin, args, err := shellparse.ParseCommand(cmd)
	if err != nil {
		panic(err)
	}

	fmt.Println("src:", cmd)
	fmt.Println("bin:", bin)
	fmt.Printf("args: %#v\n", args)
}

func parseCommandWithEnv() {
	cmd := `bash -c 'echo "it\'s complex command for user=${USER}" && sleep 3 && exit 1'`
	env := map[string]string{
		"USER": "joe",
	}

	bin, args, err := shellparse.ParseCommandWithEnv(cmd, env)
	if err != nil {
		panic(err)
	}

	fmt.Println("src:", cmd)
	fmt.Println("bin:", bin)
	fmt.Printf("args: %#v\n", args)
}
```

Will output:
```
>>> ParseCommand
src: bash -c 'echo "it\'s complex command" && sleep 3 && exit 1'
bin: bash
args: []string{"-c", "echo \"it's complex command\" && sleep 3 && exit 1"}

>>> ParseCommandWithEnv
src: bash -c 'echo "it\'s complex command for user=${USER}" && sleep 3 && exit 1'
bin: bash
args: []string{"-c", "echo \"it's complex command for user=joe\" && sleep 3 && exit 1"}
```

> Please note, library will never lookup current environment directly. Instead for 
`*WithEnv` functions it expects map with safe key=value replacements.

