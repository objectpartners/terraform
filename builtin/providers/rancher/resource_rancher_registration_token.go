package rancher

import (
	"fmt"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/rancher/go-rancher/client"
)

func resourceRancherRegistrationToken() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancherRegistrationTokenCreate,
		Read:   resourceRancherRegistrationTokenRead,
		Update: resourceRancherRegistrationTokenUpdate,
		Delete: resourceRancherRegistrationTokenDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"environment_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"command": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"image": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"token": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceRancherRegistrationTokenCreate(d *schema.ResourceData, meta interface{}) error {
	rClient, err := retry(func() (interface{}, error) {
		return meta.(*ClientProvider).clientFor(d.Get("environment_id").(string))
	}, time.Duration(30*time.Second), time.Duration(3*time.Second))
	if err != nil {
		return err
	}
	opts := &client.RegistrationToken{
		AccountId:   d.Get("environment_id").(string),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}
	token, err := rClient.(*client.RancherClient).RegistrationToken.Create(opts)
	if err != nil {
		return err
	}
	d.SetId(token.Id)
	return resourceRancherRegistrationTokenRead(d, meta)
}

func resourceRancherRegistrationTokenRead(d *schema.ResourceData, meta interface{}) error {
	rClient, err := meta.(*ClientProvider).client()
	if err != nil {
		return err
	}
	token, err := rClient.RegistrationToken.ById(d.Id())
	if err != nil {
		d.SetId("")
		return err
	}
	d.Set("environment_id", token.AccountId)
	d.Set("name", token.Name)
	d.Set("description", token.Description)
	d.Set("command", token.Command)
	d.Set("image", token.Image)
	d.Set("token", token.Token)
	return nil
}

func resourceRancherRegistrationTokenUpdate(d *schema.ResourceData, meta interface{}) error {
	rClient, err := meta.(*ClientProvider).clientFor(d.Get("environment_id").(string))
	if err != nil {
		return err
	}
	token, err := rClient.RegistrationToken.ById(d.Id())
	if err != nil {
		return err
	}
	_, err = rClient.RegistrationToken.Update(token, &client.RegistrationToken{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	})
	if err != nil {
		return err
	}
	return nil
}

func resourceRancherRegistrationTokenDelete(d *schema.ResourceData, meta interface{}) error {
	rClient, err := meta.(*ClientProvider).clientFor(d.Get("environment_id").(string))
	if err != nil {
		return err
	}
	token, err := rClient.RegistrationToken.ById(d.Id())
	if err != nil {
		return nil
	}
	_, err = rClient.RegistrationToken.ActionDeactivate(token)
	if err != nil {
		return err
	}
	rErr := resource.Retry(30*time.Second, func() *resource.RetryError {
		t, e := rClient.RegistrationToken.ById(d.Id())
		if e != nil {
			return resource.NonRetryableError(err)
		}
		if t.State != "inactive" {
			return resource.RetryableError(fmt.Errorf("RegistrationToken[%s] is not [inactive].", t.Id))
		}
		return nil
	})
	if rErr != nil {
		return errwrap.Wrapf("{{err}}", err)
	}
	token, err = rClient.RegistrationToken.ById(d.Id())
	if err != nil {
		return nil
	}
	_, err = rClient.RegistrationToken.ActionRemove(token)
	if err != nil {
		return err
	}
	rErr = resource.Retry(30*time.Second, func() *resource.RetryError {
		t, e := rClient.RegistrationToken.ById(d.Id())
		if e != nil {
			return resource.NonRetryableError(err)
		}
		if t.State != "removed" {
			return resource.RetryableError(fmt.Errorf("RegistrationToken[%s] is not [removed].", t.Id))
		}
		return nil
	})
	if rErr != nil {
		return errwrap.Wrapf("{{err}}", err)
	}
	return nil
}
