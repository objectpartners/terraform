package github

import (
	"fmt"
	"github.com/google/go-github/github"
	"github.com/hashicorp/terraform/helper/schema"
	"strings"
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
				Required:     true,
				ValidateFunc: validateRoleValue,
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

	d.SetId(buildMembershipId(membership.Organization.Login, membership.User.Login))

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

// return the pieces of the id as org, user
func parseMembershipId(id string) (string, string) {
	parts := strings.SplitN(id, ":", 2)
	return parts[0], parts[1]
}

// Since there is no id for memberships, we are storing in form organization:user
func buildMembershipId(org, user *string) string {
	return fmt.Sprintf("%s:%s", *org, *user)
}

func validateRoleValue(v interface{}, k string) (we []string, errors []error) {
	value := v.(string)
	viewTypes := map[string]bool{
		"member": true,
		"admin":  true,
	}

	if !viewTypes[value] {
		errors = append(errors, fmt.Errorf("%q is an invalid Github role type", k))
	}
	return
}
