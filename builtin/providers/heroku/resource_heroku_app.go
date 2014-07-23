package heroku

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/config"
	"github.com/hashicorp/terraform/helper/diff"
	"github.com/hashicorp/terraform/terraform"
	"github.com/bgentry/heroku-go"
)

func resource_heroku_app_create(
	s *terraform.ResourceState,
	d *terraform.ResourceDiff,
	meta interface{}) (*terraform.ResourceState, error) {
	p := meta.(*ResourceProvider)
	client := p.client

	// Merge the diff into the state so that we have all the attributes
	// properly.
	rs := s.MergeDiff(d)

	// Build up our creation options
	opts := heroku.AppCreateOpts{}

	if attr := rs.Attributes["name"]; attr != "" {
		opts.Name = &attr
	}

	if attr := rs.Attributes["region"]; attr != "" {
		opts.Region = &attr
	}

	if attr := rs.Attributes["stack"]; attr != "" {
		opts.Stack = &attr
	}

	log.Printf("[DEBUG] App create configuration: %#v", opts)

	app, err := client.AppCreate(&opts)
	if err != nil {
		return s, err
	}

	rs.ID = app.Name

	log.Printf("[INFO] App ID: %s", rs.ID)

	return resource_heroku_app_update_state(rs, app)
}

func resource_heroku_app_update(
	s *terraform.ResourceState,
	d *terraform.ResourceDiff,
	meta interface{}) (*terraform.ResourceState, error) {

	panic("does not update")

	return nil, nil
}

func resource_heroku_app_destroy(
	s *terraform.ResourceState,
	meta interface{}) error {
	p := meta.(*ResourceProvider)
	client := p.client

	log.Printf("[INFO] Deleting App: %s", s.ID)

	// Destroy the app
	err := client.AppDelete(s.ID)

	if err != nil {
		return fmt.Errorf("Error deleting App: %s", err)
	}

	return nil
}

func resource_heroku_app_refresh(
	s *terraform.ResourceState,
	meta interface{}) (*terraform.ResourceState, error) {
	p := meta.(*ResourceProvider)
	client := p.client

	app, err := resource_heroku_app_retrieve(s.ID, client)
	if err != nil {
		return nil, err
	}

	return resource_heroku_app_update_state(s, app)
}

func resource_heroku_app_diff(
	s *terraform.ResourceState,
	c *terraform.ResourceConfig,
	meta interface{}) (*terraform.ResourceDiff, error) {

	b := &diff.ResourceBuilder{
		Attrs: map[string]diff.AttrType{
			"name":              diff.AttrTypeCreate,
			"region":            diff.AttrTypeUpdate,
			"stack":              diff.AttrTypeCreate,
		},

		ComputedAttrs: []string{
			"name",
			"region",
			"stack",
			"git_url",
			"web_url",
			"id",
		},
	}

	return b.Diff(s, c)
}

func resource_heroku_app_update_state(
	s *terraform.ResourceState,
	app *heroku.App) (*terraform.ResourceState, error) {

	s.Attributes["name"] = app.Name
	s.Attributes["stack"] = app.Stack.Name
	s.Attributes["region"] = app.Region.Name
	s.Attributes["git_url"] = app.GitURL
	s.Attributes["web_url"] = app.WebURL
	s.Attributes["id"] = app.Id

	return s, nil
}

func resource_heroku_app_retrieve(id string, client *heroku.Client) (*heroku.App, error) {
	app, err := client.AppInfo(id)

	if err != nil {
		return nil, fmt.Errorf("Error retrieving app: %s", err)
	}

	return app, nil
}

func resource_heroku_app_validation() *config.Validator {
	return &config.Validator{
		Required: []string{},
		Optional: []string{
			"name",
			"region",
			"stack",
		},
	}
}