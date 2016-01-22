package github

import (
	"github.com/google/go-github/github"
	"github.com/hashicorp/terraform/helper/schema"
)

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
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceGithubTeamCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*GithubClient).client
	n := d.Get("name").(string)
	githubTeam, _, err := client.Organizations.CreateTeam(meta.(*GithubClient).organization, &github.Team{
		Name: &n,
	})
	if err != nil {
		return err
	}
	d.SetId(fromGithubId(githubTeam.ID))
	return resourceGithubTeamRead(d, meta)
}

func resourceGithubTeamRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*GithubClient).client

	_, err := getGithubTeam(d, client)
	if err != nil {
		d.SetId("")
		return nil
	}
	//TODO need to update Description but go-github doesn't expose it
	return nil
}

func resourceGithubTeamUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*GithubClient).client
	team, err := getGithubTeam(d, client)

	if err != nil {
		d.SetId("")
		return nil
	}

	//TODO need to set Description on team before sending it, but go-github doesn't expose it

	team, _, err = client.Organizations.EditTeam(*team.ID, team)
	if err != nil {
		return err
	}
	d.SetId(fromGithubId(team.ID))
	return resourceGithubTeamRead(d, meta)
}

func resourceGithubTeamDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*GithubClient).client
	id := toGithubId(d.Id())
	_, err := client.Organizations.DeleteTeam(id)
	return err
}

func getGithubTeam(d *schema.ResourceData, github *github.Client) (*github.Team, error) {
	id := toGithubId(d.Id())
	team, _, err := github.Organizations.GetTeam(id)
	return team, err
}
