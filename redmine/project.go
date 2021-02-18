package redmine

import (
	"context"
	rmapi "github.com/mattn/go-redmine"
	"github.com/pkg/errors"
	"strconv"
)

type Project struct {
	ID                 int      `json:"id"`
	Name               string   `json:"name"`
	Identifier         string   `json:"identifier"`
	Description        string   `json:"description"`
	IsPublic           bool     `json:"is_public"`
	ParentID           string   `json:"parent_id"`
	InheritMembers     bool     `json:"inherit_members"`
	TrackerIDs         []string `json:"tracker_ids"`
	EnabledModuleNames []string `json:"enabled_module_names"`
}

func (prj *Project) GetID() string {
	return strconv.Itoa(prj.ID)
}

func (c *Client) CreateProject(ctx context.Context, project *Project) error {
	apiProj := wrapProject(project)

	_, err := c.redmineAPI.CreateProject(*apiProj)
	if err != nil {
		return errors.Wrapf(err, "error while creating project (id %d)", project.ID)
	}

	return nil
}

func (c *Client) ReadProject(ctx context.Context, id string) (project *Project, err error) {
	projectID, _ := strconv.Atoi(id)
	apiProj, err := c.redmineAPI.Project(projectID)
	if err != nil {
		return project, errors.Wrapf(err, "error while reading project (id %s)", id)
	}

	project = unwrapProject(apiProj)

	return project, nil
}

func (c *Client) UpdateProject(ctx context.Context, name string, project *Project) error {
	return nil
}

func (c *Client) DeleteProject(ctx context.Context, name string) error {
	return nil
}

func (c *Client) setProject(ctx context.Context, project *Project, method string, url string) error {
	return nil
}

func wrapProject(project *Project) *rmapi.Project {
	apiProj := &rmapi.Project{
		Id:          project.ID,
		Name:        project.Name,
		Identifier:  project.Identifier,
		Description: project.Description,
	}

	if project.ParentID != "" {
		apiProj.Parent.Id, _ = strconv.Atoi(project.ParentID)
	}

	return apiProj
}

func unwrapProject(apiProj *rmapi.Project) *Project {
	project := &Project{
		ID:                 apiProj.Id,
		Name:               apiProj.Name,
		Identifier:         apiProj.Identifier,
		Description:        apiProj.Description,
		IsPublic:           true,
		InheritMembers:     false,
		TrackerIDs:         []string{},
		EnabledModuleNames: []string{},
	}

	if apiProj.Parent.Name != "" {
		project.ParentID = strconv.Itoa(apiProj.Parent.Id)
	}

	return project
}
