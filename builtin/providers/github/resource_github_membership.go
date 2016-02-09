package github

import (
	"github.com/google/go-github/github"
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
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateRoleValueFunc([]string{"member", "admin"}),
				Default:      "member",
			},
		},
	}
}

func resourceGithubMembershipCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*GithubClient).client
	n := d.Get("username").(string)
	r := d.Get("role").(string)

	membership, _, err := client.Organizations.EditOrgMembership(n, meta.(*GithubClient).organization,
		&github.Membership{Role: &r})
	if err != nil {
		return err
	}

	d.SetId(buildTwoPartId(membership.Organization.Login, membership.User.Login))

	return resourceGithubMembershipRead(d, meta)
}

func resourceGithubMembershipRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*GithubClient).client

	membership, _, err := client.Organizations.GetOrgMembership(d.Get("username").(string), meta.(*GithubClient).organization)
	if err != nil {
		d.SetId("")
		return nil
	}
	username := membership.User.Login
	roleName := membership.Role

	d.Set("username", *username)
	d.Set("role", *roleName)
	return nil
}

func resourceGithubMembershipUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*GithubClient).client
	n := d.Get("username").(string)
	r := d.Get("role").(string)

	_, _, err := client.Organizations.EditOrgMembership(n, meta.(*GithubClient).organization, &github.Membership{
		Role: &r,
	})
	if err != nil {
		return err
	}
	return nil
}

func resourceGithubMembershipDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*GithubClient).client
	n := d.Get("username").(string)

	_, err := client.Organizations.RemoveOrgMembership(n, meta.(*GithubClient).organization)

	return err
}
