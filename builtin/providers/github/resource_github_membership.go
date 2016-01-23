package github

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceGithubMembership() *schema.Resource {

	return &schema.Resource{
		Create: resourceGithubMembershipCreate,
		Read:   resourceGithubMembershipRead,
		Update: resourceGithubMembershipUpdate,
		Delete: resourceGithubMembershipDelete,

		Schema: map[string]*schema.Schema{
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"role": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				//TODO validate function (member, admin)
			},
		},
	}
}

func resourceGithubMembershipCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceGithubMembershipRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceGithubMembershipUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceGithubMembershipDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
