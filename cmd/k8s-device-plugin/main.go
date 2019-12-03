
package main

import (
"fmt"
"os"

"github.com/asdfsx/k8s-device-plugin/pkg/cmd"
)

func main() {
	command := cmd.NewRootCommand()
	if err := command.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}