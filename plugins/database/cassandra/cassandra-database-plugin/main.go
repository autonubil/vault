package main

import (
	"log"
	"os"

	"github.com/autonubil/vault/helper/pluginutil"
	"github.com/autonubil/vault/plugins/database/cassandra"
)

func main() {
	apiClientMeta := &pluginutil.APIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(os.Args)

	err := cassandra.Run(apiClientMeta.GetTLSConfig())
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
