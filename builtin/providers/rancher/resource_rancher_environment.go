package rancher

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/rancher/go-rancher/client"
)

func resourceRancherEnvironment() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancherEnvironmentCreate,
		Read:   resourceRancherEnvironmentRead,
		Update: resourceRancherEnvironmentUpdate,
		Delete: resourceRancherEnvironmentDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"engine": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "cattle",
				ValidateFunc: validateValueFunc([]string{"cattle", "kubernetes", "mesos", "swarm"}),
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"virtual_machine_support": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceRancherEnvironmentCreate(d *schema.ResourceData, meta interface{}) error {
	rClient := meta.(*client.RancherClient)
	settings := &client.Project{
		Name:           d.Get("name").(string),
		Description:    d.Get("description").(string),
		VirtualMachine: d.Get("virtual_machine_support").(bool),
	}
	engine := d.Get("engine").(string)
	if engine == "kubernetes" {
		settings.Kubernetes = true
	}
	if engine == "mesos" {
		settings.Mesos = true
	}
	if engine == "swarm" {
		settings.Swarm = true
	}
	environment, err := rClient.Project.Create(settings)
	if err != nil {
		return err
	}
	d.SetId(environment.Id)
	return resourceRancherEnvironmentRead(d, meta)
}

func resourceRancherEnvironmentRead(d *schema.ResourceData, meta interface{}) error {
	rClient := meta.(*client.RancherClient)
	environment, err := rClient.Project.ById(d.Id())
	if err != nil {
		d.SetId("")
		return nil
	}
	d.Set("description", environment.Description)
	if environment.Kubernetes {
		d.Set("engine", "kubernetes")
	} else if environment.Mesos {
		d.Set("engine", "mesos")
	} else if environment.Swarm {
		d.Set("engine", "swarm")
	} else {
		d.Set("engine", "cattle")
	}
	return nil
}

func resourceRancherEnvironmentUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceRancherEnvironmentDelete(d *schema.ResourceData, meta interface{}) error {
	rClient := meta.(*client.RancherClient)
	environment, err := rClient.Project.ById(d.Id())
	if err != nil {
		return nil
	}
	err = rClient.Project.Delete(environment)
	if err != nil {
		return err
	}
	_, _ = retry(func() (interface{}, error) {
		environment, err := rClient.Project.ById(d.Id())
		if err != nil {
			return nil, err
		}
		if environment.State != "removed" && environment.State != "purged" {
			return nil, fmt.Errorf("Environment[%s] is not deleted.", environment.Id)
		}
		return environment, nil
	}, time.Duration(30*time.Second), time.Duration(3*time.Second))
	return nil
}
