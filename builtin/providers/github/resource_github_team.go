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
	s := d.Get("description").(string)
	githubTeam, _, err := client.Organizations.CreateTeam(meta.(*GithubClient).organization, &github.Team{
		Name:        &n,
		Description: &s,
	})
	if err != nil {
		return err
	}
	d.SetId(fromGithubId(githubTeam.ID))
	return resourceGithubTeamRead(d, meta)
}

func resourceGithubTeamRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*GithubClient).client

	t, err := getGithubTeam(d, client)
	if err != nil {
		d.SetId("")
		return nil
	}
	d.Set("name", t.Name)
	d.Set("description", t.Description)
	return nil
}

func resourceGithubTeamUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*GithubClient).client
	n := d.Get("name").(string)
	s := d.Get("description").(string)
	id := d.Id()

	_, _, err := client.Organizations.EditTeam(toGithubId(id), &github.Team{
		Name:        &n,
		Description: &s,
	})
	if err != nil {
		return err
	}
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
