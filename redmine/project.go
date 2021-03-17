package redmine

import (
	"context"
	rmapi "github.com/cloudogu/go-redmine"
	"github.com/pkg/errors"
	"strconv"
)

type Project struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Identifier     string `json:"identifier"`
	Description    string `json:"description"`
	Homepage       string `json:"homepage"`
	IsPublic       bool   `json:"is_public"`
	ParentID       string `json:"parent_id"`
	InheritMembers bool   `json:"inherit_members"`
	CreatedOn      string `json:"created_on"`
	UpdatedOn      string `json:"updated_on"`
}

func (c *Client) CreateProject(ctx context.Context, project *Project) (*Project, error) {
	apiProj := wrapProject(project)

	actualAPIProject, err := c.redmineAPI.CreateProject(*apiProj)
	if err != nil {
		return nil, errors.Wrapf(err, "error while creating project (identifier: %s)", project.Identifier)
	}

	actualProject := unwrapProject(actualAPIProject)

	return actualProject, nil
}

func (c *Client) ReadProject(ctx context.Context, id string) (project *Project, err error) {
	idInt, err := verifyIDtoInt(id)
	if err != nil {
		return nil, errors.Wrap(err, "could not read project because of malformed input data")
	}

	apiProj, err := c.redmineAPI.Project(idInt)
	if err != nil {
		return project, errors.Wrapf(err, "error while reading project (id: %d)", idInt)
	}

	project = unwrapProject(apiProj)

	return project, nil
}

func (c *Client) UpdateProject(ctx context.Context, project *Project) (updatedProject *Project, err error) {
	apiProj := *wrapProject(project)

	err = c.redmineAPI.UpdateProject(apiProj)
	if err != nil {
		return project, errors.Wrapf(err, "error while updating project (id: %s, identifier: %s)", project.ID, project.Identifier)
	}

	project = unwrapProject(&apiProj)

	return project, nil
}

func (c *Client) DeleteProject(ctx context.Context, id string) error {
	idInt, err := verifyIDtoInt(id)
	if err != nil {
		return errors.Wrap(err, "could not delete project because of malformed input data")
	}

	err = c.redmineAPI.DeleteProject(idInt)
	if err != nil {
		return errors.Wrapf(err, "error while deleteting project (id: %d)", idInt)
	}

	return nil

}

func wrapProject(project *Project) *rmapi.Project {
	apiProj := &rmapi.Project{
		Name:           project.Name,
		Identifier:     project.Identifier,
		Description:    project.Description,
		Homepage:       project.Homepage,
		IsPublic:       project.IsPublic,
		InheritMembers: project.InheritMembers,
		CreatedOn:      project.CreatedOn,
		UpdatedOn:      project.UpdatedOn,
	}

	if project.ID != "" && project.ID != "0" {
		apiProj.Id, _ = strconv.Atoi(project.ID)
	}
	if project.ParentID != "" {
		apiProj.ParentID.Id, _ = strconv.Atoi(project.ParentID)
	}

	return apiProj
}

func unwrapProject(apiProj *rmapi.Project) *Project {
	project := &Project{
		Name:           apiProj.Name,
		Identifier:     apiProj.Identifier,
		Description:    apiProj.Description,
		Homepage:       apiProj.Homepage,
		IsPublic:       apiProj.IsPublic,
		InheritMembers: apiProj.InheritMembers,
		CreatedOn:      apiProj.CreatedOn,
		UpdatedOn:      apiProj.UpdatedOn,
	}

	if apiProj.Id != 0 {
		project.ID = strconv.Itoa(apiProj.Id)
	}
	if apiProj.ParentID.Id != 0 {
		project.ParentID = strconv.Itoa(apiProj.ParentID.Id)
	}

	return project
}
