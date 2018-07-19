package main

import (
	"fmt"

	"github.com/arkady-emelyanov/go-shellparse"
)

func main() {
	fmt.Println(">>> Command ")
	parseCommandSimple()

	fmt.Println("")
	fmt.Println("")
	fmt.Println(">>> CommandWithEnv")
	parseCommandWithEnv()
}

func parseCommandSimple() {
	cmd := `bash -c 'echo "it\'s complex command" && sleep 3 && exit 1'`
	bin, args, err := shellparse.Command(cmd)
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

	bin, args, err := shellparse.CommandWithEnv(cmd, env)
	if err != nil {
		panic(err)
	}

	fmt.Println("src:", cmd)
	fmt.Println("bin:", bin)
	fmt.Printf("args: %#v\n", args)
}
