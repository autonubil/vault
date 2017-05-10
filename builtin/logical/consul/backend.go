package consul

import (
	"github.com/autonubil/vault/logical"
	"github.com/autonubil/vault/logical/framework"
)

func Factory(conf *logical.BackendConfig) (logical.Backend, error) {
	return Backend().Setup(conf)
}

func Backend() *backend {
	var b backend
	b.Backend = &framework.Backend{
		Paths: []*framework.Path{
			pathConfigAccess(),
			pathListRoles(&b),
			pathRoles(),
			pathToken(&b),
		},

		Secrets: []*framework.Secret{
			secretToken(&b),
		},
	}

	return &b
}

type backend struct {
	*framework.Backend
}
