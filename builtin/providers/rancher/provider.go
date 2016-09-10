package rancher

import (
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/rancher/go-rancher/client"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("RANCHER_URL", nil),
				Description: "Rancher Server URL",
			},
			"access_key": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				DefaultFunc: schema.EnvDefaultFunc("RANCHER_ACCESS_KEY", nil),
				Description: "Rancher API Access Key",
			},
			"secret_key": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				DefaultFunc: schema.EnvDefaultFunc("RANCHER_SECRET_KEY", nil),
				Description: "Rancher API Secret Key",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"rancher_environment":        resourceRancherEnvironment(),
			"rancher_registration_token": resourceRancherRegistrationToken(),
			"rancher_stack":              resourceRancherStack(),
		},

		DataSourcesMap: map[string]*schema.Resource{},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	opts := &client.ClientOpts{
		Url: d.Get("url").(string),
	}
	if key := d.Get("access_key").(string); key != "" {
		opts.AccessKey = key
	}
	if secret := d.Get("secret_key").(string); secret != "" {
		opts.SecretKey = secret
	}
	return &ClientProvider{
		Opts: opts,
	}, nil
}

type ClientProvider struct {
	Opts *client.ClientOpts
}

func (rancherClientProvider *ClientProvider) client() (*client.RancherClient, error) {
	return client.NewRancherClient(rancherClientProvider.Opts)
}

func (rancherClientProvider *ClientProvider) clientFor(environmentID string) (*client.RancherClient, error) {
	opts := &client.ClientOpts{
		Url:       strings.Join([]string{rancherClientProvider.Opts.Url, "v1", "projects", environmentID, "schemas"}, "/"),
		AccessKey: rancherClientProvider.Opts.AccessKey,
		SecretKey: rancherClientProvider.Opts.SecretKey,
	}
	return client.NewRancherClient(opts)
}
