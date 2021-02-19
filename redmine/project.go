package redmine

import (
	"context"
	rmapi "github.com/cloudogu/go-redmine"
	"github.com/pkg/errors"
	"strconv"
)

type Project struct {
	ID                 string   `json:"id"`
	Name               string   `json:"name"`
	Identifier         string   `json:"identifier"`
	Description        string   `json:"description"`
	IsPublic           bool     `json:"is_public"`
	ParentID           string   `json:"parent_id"`
	InheritMembers     bool     `json:"inherit_members"`
	TrackerIDs         []string `json:"tracker_ids"`
	EnabledModuleNames []string `json:"enabled_module_names"`
	UpdatedOn          string   `json:"updated_on"`
}

func (c *Client) CreateProject(ctx context.Context, project *Project) (*Project, error) {
	apiProj := wrapProject(project)

	actualAPIProject, err := c.redmineAPI.CreateProject(*apiProj)
	if err != nil {
		return nil, errors.Wrapf(err, "error while creating project (id %s)", project.ID)
	}

	actualProject := unwrapProject(actualAPIProject)

	return actualProject, nil
}

func (c *Client) ReadProject(ctx context.Context, id string) (project *Project, err error) {
	idInt, _ := strconv.Atoi(id)
	apiProj, err := c.redmineAPI.Project(idInt)
	if err != nil {
		return project, errors.Wrapf(err, "error while reading project (id %s)", id)
	}

	project = unwrapProject(apiProj)

	return project, nil
}

func (c *Client) UpdateProject(ctx context.Context, project *Project) (updatedProject *Project, err error) {
	apiProj := *wrapProject(project)

	err = c.redmineAPI.UpdateProject(apiProj)
	if err != nil {
		return project, errors.Wrapf(err, "error while updating project (id %d)", apiProj.Id)
	}

	project = unwrapProject(&apiProj)

	return project, nil
}

func (c *Client) DeleteProject(ctx context.Context, name string) error {
	return nil
}

func wrapProject(project *Project) *rmapi.Project {
	apiProj := &rmapi.Project{
		Name:        project.Name,
		Identifier:  project.Identifier,
		Description: project.Description,
		UpdatedOn:   project.UpdatedOn,
	}

	if project.ID != "" {
		apiProj.Id, _ = strconv.Atoi(project.ID)
	}
	if project.ParentID != "" {
		apiProj.ParentID.Id, _ = strconv.Atoi(project.ParentID)
	}

	return apiProj
}

func unwrapProject(apiProj *rmapi.Project) *Project {
	project := &Project{
		ID:                 strconv.Itoa(apiProj.Id),
		Name:               apiProj.Name,
		Identifier:         apiProj.Identifier,
		Description:        apiProj.Description,
		IsPublic:           true,
		InheritMembers:     false,
		TrackerIDs:         []string{},
		EnabledModuleNames: []string{},
		UpdatedOn:          apiProj.UpdatedOn,
	}

	if apiProj.ParentID.Id != 0 {
		project.ParentID = strconv.Itoa(apiProj.ParentID.Id)
	}

	return project
}
