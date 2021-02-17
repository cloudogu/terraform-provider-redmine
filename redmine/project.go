package redmine

import (
	"context"
)

type Project struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	Identifier         string `json:"identifier"`
	Description        string `json:"description"`
	IsPublic           string `json:"is_public"`
	ParentID           string `json:"parent_id"`
	InheritMembers     string `json:"inherit_members"`
	TrackerIDs         []int  `json:"tracker_ids"`
	EnabledModuleNames string `json:"enabled_module_names"`
}

func (prj *Project) GetID() string {
	return prj.ID
}

func (c *Client) CreateProject(ctx context.Context, project Project) error {
	return nil
}

func (c *Client) GetProject(ctx context.Context, name string) (Project, error) {
	return Project{}, nil
}

func (c *Client) UpdateProject(ctx context.Context, name string, project Project) error {
	return nil
}

func (c *Client) DeleteProject(ctx context.Context, name string) error {
	return nil
}

func (c *Client) setProject(ctx context.Context, project Project, method string, url string) error {
	return nil
}
