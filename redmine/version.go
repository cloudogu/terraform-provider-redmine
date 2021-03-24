package redmine

import (
	"context"
	"fmt"
	rmapi "github.com/cloudogu/go-redmine"
	"github.com/pkg/errors"
	"strconv"
)

type Version struct {
	ID          string `json:"id"`
	ProjectID   int    `json:"project_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	DueDate     string `json:"due_date"`
	CreatedOn   string `json:"created_on"`
	UpdatedOn   string `json:"updated_on"`
}

func (i *Version) String() string {
	return fmt.Sprintf("Version{ID=%s,ProjectID=%d,Name=%s}", i.ID, i.ProjectID, i.Name)
}

func (c *Client) CreateVersion(_ context.Context, Version *Version) (*Version, error) {
	apiVersion := wrapVersion(Version)

	createdAPIVersion, err := c.redmineAPI.CreateVersion(*apiVersion)
	if err != nil {
		return nil, errors.Wrapf(err, "error while creating version (project id: %d, name: %s)", Version.ProjectID, Version.Name)
	}

	actualVersion := unwrapVersion(createdAPIVersion)

	return actualVersion, nil
}

func (c *Client) ReadVersion(_ context.Context, id string) (Version *Version, err error) {
	idInt, err := verifyIDtoInt(id)
	if err != nil {
		return nil, errors.Wrap(err, "could not read version because of malformed input data")
	}

	apiVersion, err := c.redmineAPI.Version(idInt)
	if err != nil {
		return Version, errors.Wrapf(err, "error while reading version (id: %d)", idInt)
	}

	return unwrapVersion(apiVersion), nil
}

func (c *Client) UpdateVersion(_ context.Context, Version *Version) (updatedVersion *Version, err error) {
	_, err = verifyIDtoInt(Version.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "could not update version (id: %s, name: %s) because of malformed input data",
			Version.ID, Version.Name)
	}

	apiVersion := *wrapVersion(Version)

	err = c.redmineAPI.UpdateVersion(apiVersion)
	if err != nil {
		return Version, errors.Wrapf(err, "error while updating version (id: %d, name: %s)", apiVersion.Id, Version.Name)
	}

	return unwrapVersion(&apiVersion), nil
}

func (c *Client) DeleteVersion(_ context.Context, id string) error {
	idInt, err := verifyIDtoInt(id)
	if err != nil {
		return errors.Wrap(err, "could not delete version because of malformed input data")
	}

	err = c.redmineAPI.DeleteVersion(idInt)
	if err != nil {
		return errors.Wrapf(err, "error while deleteting version (id: %d)", idInt)
	}

	return nil
}

func wrapVersion(Version *Version) *rmapi.Version {
	apiVersion := &rmapi.Version{
		Project:     rmapi.IdName{Id: Version.ProjectID},
		Name:        Version.Name,
		Description: Version.Description,
		Status:      Version.Status,
		DueDate:     Version.DueDate,
		CreatedOn:   Version.CreatedOn,
		UpdatedOn:   Version.UpdatedOn,
	}

	if Version.ID != "" {
		apiVersion.Id, _ = strconv.Atoi(Version.ID)
	}

	return apiVersion
}

func unwrapVersion(apiVersion *rmapi.Version) *Version {
	Version := &Version{
		ID:          strconv.Itoa(apiVersion.Id),
		ProjectID:   apiVersion.Project.Id,
		Name:        apiVersion.Name,
		Description: apiVersion.Description,
		Status:      apiVersion.Status,
		DueDate:     apiVersion.DueDate,
		CreatedOn:   apiVersion.CreatedOn,
		UpdatedOn:   apiVersion.UpdatedOn,
	}

	return Version
}
