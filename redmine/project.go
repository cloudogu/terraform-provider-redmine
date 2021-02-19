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
		return nil, errors.Wrapf(err, "error while creating project (identifier %d)", project.Identifier)
	}

	actualProject := unwrapProject(actualAPIProject)

	return actualProject, nil
}

func (c *Client) ReadProject(ctx context.Context, identifier string) (project *Project, err error) {
	project, err = c.getProjectByIdentifier(ctx, identifier)
	if err != nil {
		return project, errors.Wrapf(err, "error while reading project (identifier %s)", identifier)
	}

	return project, nil
}

func (c *Client) UpdateProject(ctx context.Context, project *Project) (updatedProject *Project, err error) {
	apiProj := *wrapProject(project)

	err = c.redmineAPI.UpdateProject(apiProj)
	if err != nil {
		return project, errors.Wrapf(err, "error while updating project (identifier %s)", apiProj.Identifier)
	}

	project = unwrapProject(&apiProj)

	return project, nil
}

func (c *Client) DeleteProject(ctx context.Context, name string) error {
	return nil
}

func (c *Client) getProjectByIdentifier(ctx context.Context, identifier string) (*Project, error) {
	apiProjects, err := c.redmineAPI.Projects()
	if err != nil {
		return nil, errors.Wrapf(err, "error while fetching project by identifier %s", identifier)
	}

	for _, apiProj := range apiProjects {
		if apiProj.Identifier == identifier {
			project := unwrapProject(&apiProj)
			return project, nil
		}
	}

	return nil, errors.Errorf("could not find project by identifier %s", identifier)
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
