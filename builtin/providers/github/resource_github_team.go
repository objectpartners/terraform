package github

import "github.com/hashicorp/terraform/helper/schema"

func resourceGithubTeam() *schema.Resource {

	return &schema.Resource{
		Create: resourceGithubTeamCreate,
		Read:   resourceGithubTeamRead,
		Update: resourceGithubTeamUpdate,
		Delete: resourceGithubTeamDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceGithubTeamCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*GithubClient).client
	return nil
}

func resourceGithubTeamRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceGithubTeamUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceGithubTeamDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
