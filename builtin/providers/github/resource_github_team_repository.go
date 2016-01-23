package github

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceGithubTeamRepository() *schema.Resource {

	return &schema.Resource{
		Create: resourceGithubTeamRepositoryCreate,
		Read:   resourceGithubTeamRepositoryRead,
		Update: resourceGithubTeamRepositoryUpdate,
		Delete: resourceGithubTeamRepositoryDelete,

		Schema: map[string]*schema.Schema{
			"team_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"repository": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"permission": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				//TODO validate function (pull, push, admin)
			},
		},
	}
}

func resourceGithubTeamRepositoryCreate(d *schema.ResourceData, meta interface{}) error {
	//TODO the ID to store in the state file needs to be team_id:repository
	return nil
}

func resourceGithubTeamRepositoryRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceGithubTeamRepositoryUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceGithubTeamRepositoryDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
