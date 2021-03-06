package ssh

import (
	"strings"
	"sync"

	"github.com/autonubil/vault/helper/salt"
	"github.com/autonubil/vault/logical"
	"github.com/autonubil/vault/logical/framework"
)

type backend struct {
	*framework.Backend
	view      logical.Storage
	salt      *salt.Salt
	saltMutex sync.RWMutex
}

func Factory(conf *logical.BackendConfig) (logical.Backend, error) {
	b, err := Backend(conf)
	if err != nil {
		return nil, err
	}
	return b.Setup(conf)
}

func Backend(conf *logical.BackendConfig) (*backend, error) {
	var b backend
	b.view = conf.StorageView
	b.Backend = &framework.Backend{
		Help: strings.TrimSpace(backendHelp),

		PathsSpecial: &logical.Paths{
			Unauthenticated: []string{
				"verify",
				"public_key",
			},

			LocalStorage: []string{
				"otp/",
			},
		},

		Paths: []*framework.Path{
			pathConfigZeroAddress(&b),
			pathKeys(&b),
			pathListRoles(&b),
			pathRoles(&b),
			pathCredsCreate(&b),
			pathLookup(&b),
			pathVerify(&b),
			pathConfigCA(&b),
			pathSign(&b),
			pathFetchPublicKey(&b),
		},

		Secrets: []*framework.Secret{
			secretDynamicKey(&b),
			secretOTP(&b),
		},

		Init: b.initialize,

		Invalidate: b.invalidate,
	}
	return &b, nil
}

func (b *backend) initialize() error {
	b.saltMutex.Lock()
	defer b.saltMutex.Unlock()
	salt, err := salt.NewSalt(b.view, &salt.Config{
		HashFunc: salt.SHA256Hash,
		Location: salt.DefaultLocation,
	})
	if err != nil {
		return err
	}
	b.salt = salt
	return nil
}

func (b *backend) invalidate(key string) {
	switch key {
	case salt.DefaultLocation:
		// reread the salt
		b.initialize()
	}
}

const backendHelp = `
The SSH backend generates credentials allowing clients to establish SSH
connections to remote hosts.

There are three variants of the backend, which generate different types of
credentials: dynamic keys, One-Time Passwords (OTPs) and certificate authority. The desired behavior
is role-specific and chosen at role creation time with the 'key_type'
parameter.

Please see the backend documentation for a thorough description of both
types. The Vault team strongly recommends the OTP type.

After mounting this backend, before generating credentials, configure the
backend's lease behavior using the 'config/lease' endpoint and create roles
using the 'roles/' endpoint.
`
