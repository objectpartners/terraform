package rancher

import "github.com/hashicorp/terraform/helper/schema"

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
			"docker_compose": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"rancher_compose": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceRancherStackCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceRancherStackRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceRancherStackUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceRancherStackDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
