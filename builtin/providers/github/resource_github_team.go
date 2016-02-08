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
	desc := d.Get("description").(string)
	githubTeam, _, err := client.Organizations.CreateTeam(meta.(*GithubClient).organization, &github.Team{
		Name:        &n,
		Description: &desc,
	})
	if err != nil {
		return err
	}
	d.SetId(fromGithubId(githubTeam.ID))
	return resourceGithubTeamRead(d, meta)
}

func resourceGithubTeamRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*GithubClient).client

	team, err := getGithubTeam(d, client)
	if err != nil {
		d.SetId("")
		return nil
	}
	d.Set("description", team.Description)
	d.Set("name", team.Name)
	return nil
}

func resourceGithubTeamUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*GithubClient).client
	team, err := getGithubTeam(d, client)

	if err != nil {
		d.SetId("")
		return nil
	}

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	team.Description = &description
	team.Name = &name

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
