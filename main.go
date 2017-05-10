package main // import "github.com/autonubil/vault"

import (
	"os"

	"github.com/autonubil/vault/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:]))
}
