package command

import (
	"testing"

	"github.com/autonubil/vault/http"
	"github.com/autonubil/vault/meta"
	"github.com/autonubil/vault/vault"
	"github.com/mitchellh/cli"
)

func TestHelp(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	ui := new(cli.MockUi)
	c := &PathHelpCommand{
		Meta: meta.Meta{
			ClientToken: token,
			Ui:          ui,
		},
	}

	args := []string{
		"-address", addr,
		"sys/mounts",
	}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}
}
