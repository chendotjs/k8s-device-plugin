package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/spf13/pflag"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/component-base/logs"

	"github.com/chendotjs/k8s-device-plugin/pkg/cmd"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	pflag.CommandLine.SetNormalizeFunc(cliflag.WordSepNormalizeFunc)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	logs.InitLogs()
	defer logs.FlushLogs()

	command := cmd.NewRootCommand()
	if err := command.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
