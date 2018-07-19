# Parse strings à la shell

[![GoDoc](https://godoc.org/github.com/arkady-emelyanov/go-shellparse?status.svg)](https://godoc.org/github.com/arkady-emelyanov/go-shellparse)
[![Build Status](https://travis-ci.org/arkady-emelyanov/go-shellparse.svg?branch=master)](https://travis-ci.org/arkady-emelyanov/go-shellparse)
[![Go Report Card](https://goreportcard.com/badge/github.com/arkady-emelyanov/go-shellparse)](https://goreportcard.com/report/github.com/arkady-emelyanov/go-shellparse)
[![codecov](https://codecov.io/gh/arkady-emelyanov/go-shellparse/branch/master/graph/badge.svg)](https://codecov.io/gh/arkady-emelyanov/go-shellparse)


Whenever you need parse command and arguments from config file,
you facing quotes/escaping problem.

Library solves complexity of parsing such strings.

## Features

* No dependencies
* Ability to parse complex and multiline strings
* Useful helpers to parse string into
    * command and arguments
    * map
    * slice
* Ability to expand variables like`${VAR}` with provided k/v map
* DotEnv-like file parser
* Remove unnecessary quotes

## Installation

`go get -u github.com/arkady-emelyanov/go-shellparse`

## Quick Start

Simple:
```
bin, args, err := shellparse.Command(`bash -c 'echo "It\'s awesome"'`)
// bin: bash
// args: []string{"-c", "echo \"It's awesome\""}
```

With custom environment:
```
env := map[string]string{}{
    "USER": "johndoe",
}
bin, args, err := shellparse.CommandWithEnv(`echo ${USER}`, env)
// bin: echo
// args: []string{"johndoe"}
```

If string contains ${VAR} which is not present in provided map,
error will be raised.

> Please note, `*WithEnv` functions will never lookup current environment directly. 
Instead for they expect a map with safe key=value replacements.

## Other helpers

```
parts, err := shellparse.StringToSlice(`one 'two' "three"`)
// []string{"one", "two", "three"}

parts, err := shellparse.StringToMap(`foo=bar`)
// []map[string]string{"foo": "bar"}
```

> Library supports comments defined as `#` char. The rest of the line
will be ignored. 

## License

Licensed under the [MIT License](http://www.opensource.org/licenses/MIT).
