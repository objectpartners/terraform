package github

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceGithubTeamMembership() *schema.Resource {

	return &schema.Resource{
		Create: resourceGithubTeamMembershipCreate,
		Read:   resourceGithubTeamMembershipRead,
		Update: resourceGithubTeamMembershipUpdate,
		Delete: resourceGithubTeamMembershipDelete,

		Schema: map[string]*schema.Schema{
			"team_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"role": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "member",
				//TODO validate function (member, maintainer)
			},
		},
	}
}

func resourceGithubTeamMembershipCreate(d *schema.ResourceData, meta interface{}) error {
	//TODO the ID to store in the state file needs to be team_id:username
	return nil
}

func resourceGithubTeamMembershipRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceGithubTeamMembershipUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceGithubTeamMembershipDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
