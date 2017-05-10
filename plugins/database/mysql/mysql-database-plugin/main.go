package main

import (
	"log"
	"os"

	"github.com/autonubil/vault/helper/pluginutil"
	"github.com/autonubil/vault/plugins/database/mysql"
)

func main() {
	apiClientMeta := &pluginutil.APIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(os.Args)

	err := mysql.Run(apiClientMeta.GetTLSConfig())
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
