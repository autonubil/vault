package rabbitmq

import (
	"fmt"

	"github.com/hashicorp/go-uuid"
	"github.com/autonubil/vault/logical"
	"github.com/autonubil/vault/logical/framework"
	"github.com/michaelklishin/rabbit-hole"
)

func pathCreds(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "creds/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the role.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathCredsRead,
		},

		HelpSynopsis:    pathRoleCreateReadHelpSyn,
		HelpDescription: pathRoleCreateReadHelpDesc,
	}
}

// Issues the credential based on the role name
func (b *backend) pathCredsRead(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	if name == "" {
		return logical.ErrorResponse("missing name"), nil
	}

	// Get the role
	role, err := b.Role(req.Storage, name)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return logical.ErrorResponse(fmt.Sprintf("unknown role: %s", name)), nil
	}

	// Ensure username is unique
	uuidVal, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}
	username := fmt.Sprintf("%s-%s", req.DisplayName, uuidVal)

	password, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}

	// Get the client configuration
	client, err := b.Client(req.Storage)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return logical.ErrorResponse("failed to get the client"), nil
	}

	// Register the generated credentials in the backend, with the RabbitMQ server
	if _, err = client.PutUser(username, rabbithole.UserSettings{
		Password: password,
		Tags:     role.Tags,
	}); err != nil {
		return nil, fmt.Errorf("failed to create a new user with the generated credentials")
	}

	// If the role had vhost permissions specified, assign those permissions
	// to the created username for respective vhosts.
	for vhost, permission := range role.VHosts {
		if _, err := client.UpdatePermissionsIn(vhost, username, rabbithole.Permissions{
			Configure: permission.Configure,
			Write:     permission.Write,
			Read:      permission.Read,
		}); err != nil {
			// Delete the user because it's in an unknown state
			if _, rmErr := client.DeleteUser(username); rmErr != nil {
				return nil, fmt.Errorf("failed to delete user:%s, err: %s. %s", username, err, rmErr)
			}
			return nil, fmt.Errorf("failed to update permissions to the %s user. err:%s", username, err)
		}
	}

	// Return the secret
	resp := b.Secret(SecretCredsType).Response(map[string]interface{}{
		"username": username,
		"password": password,
	}, map[string]interface{}{
		"username": username,
	})

	// Determine if we have a lease
	lease, err := b.Lease(req.Storage)
	if err != nil {
		return nil, err
	}
	if lease != nil {
		resp.Secret.TTL = lease.TTL
	}

	return resp, nil
}

const pathRoleCreateReadHelpSyn = `
Request RabbitMQ credentials for a certain role.
`

const pathRoleCreateReadHelpDesc = `
This path reads RabbitMQ credentials for a certain role. The
RabbitMQ credentials will be generated on demand and will be automatically
revoked when the lease is up.
`
