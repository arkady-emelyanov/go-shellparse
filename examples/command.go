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
