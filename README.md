# Parse strings à la shell

[![GoDoc](https://godoc.org/github.com/arkady-emelyanov/go-shellparse?status.svg)](https://godoc.org/github.com/arkady-emelyanov/go-shellparse)
[![Go Report Card](https://goreportcard.com/badge/github.com/arkady-emelyanov/go-shellparse)](https://goreportcard.com/report/github.com/arkady-emelyanov/go-shellparse)
[![Codecov](https://codecov.io/gh/arkady-emelyanov/go-shellparse/branch/master/graph/badge.svg)](https://codecov.io/gh/arkady-emelyanov/go-shellparse)
![CodeQL](https://github.com/arkady-emelyanov/go-shellparse/workflows/CodeQL/badge.svg)
![CI Build](https://github.com/arkady-emelyanov/go-shellparse/actions/workflows/go.yml/badge.svg)

Whenever you need parse command and arguments from a config file,
you facing quotes/escaping problem. Library hides complexity 
of parsing such strings.

## Features

* Ability to parse complex and multiline strings
* Comments for multiline strings via `#` char
* Useful helpers to parse string into:
    * command and arguments
    * map
    * slice
* Ability to expand variables like`${VAR}` with provided k/v map
* DotEnv-like file parser
* Remove unnecessary quotes from commands and arguments

## Installation

`go get -u github.com/arkady-emelyanov/go-shellparse`

## Quick Start

Simple usage:
```
bin, args, err := shellparse.Command(`bash -c 'echo "It\'s awesome"'`)
// bin: bash
// args: []string{"-c", "echo \"It's awesome\""}
```

With current environment
```
// export SLEEP=1

bin, args, err := shellparse.CommandWithEnv(`bash -c 'sleep ${SLEEP}'`)
// bin: bash
// args: []string{"-c", "sleep 1"}

```

With custom variables:
```
vars := map[string]string{}{
    "USER": "johndoe",
}
bin, args, err := shellparse.CommandWithMap(`echo ${USER}`, vars)
// bin: echo
// args: []string{"johndoe"}
```

If string contains ${VAR} which is not present in provided map,
error will be raised. To avoid, escape variable in the following way: `\${VAR}`.

> Please note, `*WithVars` functions will never lookup current environment directly. 
All `key`=`value` var replacements expected to be present in a provided map.

## Other helpers

```
parts, err := shellparse.StringToSlice(`one 'two' "three"`)
// []string{"one", "two", "three"}

parts, err := shellparse.StringToMap(`foo=bar foobar=baz`)
// []map[string]string{"foo": "bar", "foobar": "baz"}
```

> Library supports comments defined as `#` char. The rest of the line
will be ignored. 

## License

Licensed under the [MIT License](http://www.opensource.org/licenses/MIT).
