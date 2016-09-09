package rancher

import (
	"fmt"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/rancher/go-rancher/client"
)

func resourceRancherStack() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancherStackCreate,
		Read:   resourceRancherStackRead,
		Update: resourceRancherStackUpdate,
		Delete: resourceRancherStackDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"environment_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"docker_compose": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"rancher_compose": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"start_on_create": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
				Default:  true,
			},
		},
	}
}

func resourceRancherStackCreate(d *schema.ResourceData, meta interface{}) error {
	rClient, err := meta.(*ClientProvider).clientFor(d.Get("environment_id").(string))
	if err != nil {
		return err
	}
	settings := &client.Environment{
		AccountId:      d.Get("environment_id").(string),
		Description:    d.Get("description").(string),
		DockerCompose:  d.Get("docker_compose").(string),
		Name:           d.Get("name").(string),
		RancherCompose: d.Get("rancher_compose").(string),
		StartOnCreate:  d.Get("start_on_create").(bool),
	}
	stack, err := rClient.Environment.Create(settings)
	if err != nil {
		return err
	}
	d.SetId(stack.Id)
	rErr := resource.Retry(30*time.Second, func() *resource.RetryError {
		stack, e := rClient.Environment.ById(d.Id())
		if e != nil {
			return resource.NonRetryableError(err)
		}
		if stack.State != "active" {
			return resource.RetryableError(fmt.Errorf("Stack[%s] is not active[%s].", stack.Id, stack.State))
		}
		return nil
	})
	if rErr != nil {
		d.SetId("")
		return errwrap.Wrapf("{{err}}", rErr)
	}
	return resourceRancherStackRead(d, meta)
}

func resourceRancherStackRead(d *schema.ResourceData, meta interface{}) error {
	rClient, err := meta.(*ClientProvider).clientFor(d.Get("environment_id").(string))
	if err != nil {
		return err
	}
	stack, err := rClient.Environment.ById(d.Id())
	if err != nil {
		d.SetId("")
		return err
	}
	d.Set("description", stack.Description)
	d.Set("docker_compose", stack.DockerCompose)
	d.Set("name", stack.Name)
	d.Set("rancher_compose", stack.RancherCompose)
	d.Set("start_on_create", stack.StartOnCreate)
	return nil
}

func resourceRancherStackUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceRancherStackDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
