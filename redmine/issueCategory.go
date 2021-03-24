package redmine

import (
	"context"
	"fmt"
	rmapi "github.com/cloudogu/go-redmine"
	"github.com/pkg/errors"
	"strconv"
)

type IssueCategory struct {
	ID        string `json:"id"`
	ProjectID int    `json:"project_id"`
	Name      string `json:"name"`
}

func (i *IssueCategory) String() string {
	return fmt.Sprintf("IssueCategory{ID=%s,ProjectID=%d,Name=%s}", i.ID, i.ProjectID, i.Name)
}

func (c *Client) CreateIssueCategory(_ context.Context, IssueCategory *IssueCategory) (*IssueCategory, error) {
	apiIssueCategory := wrapIssueCategory(IssueCategory)

	createdAPIIssueCategory, err := c.redmineAPI.CreateIssueCategory(*apiIssueCategory)
	if err != nil {
		return nil, errors.Wrapf(err, "error while creating issue category (project id: %d, name: %s)", IssueCategory.ProjectID, IssueCategory.Name)
	}

	actualIssueCategory := unwrapIssueCategory(createdAPIIssueCategory)

	return actualIssueCategory, nil
}

func (c *Client) ReadIssueCategory(_ context.Context, id string) (IssueCategory *IssueCategory, err error) {
	idInt, err := verifyIDtoInt(id)
	if err != nil {
		return nil, errors.Wrap(err, "could not read issue category because of malformed input data")
	}

	apiIssueCategory, err := c.redmineAPI.IssueCategory(idInt)
	if err != nil {
		return IssueCategory, errors.Wrapf(err, "error while reading issue category (id: %d)", idInt)
	}

	return unwrapIssueCategory(apiIssueCategory), nil
}

func (c *Client) UpdateIssueCategory(_ context.Context, IssueCategory *IssueCategory) (updatedIssueCategory *IssueCategory, err error) {
	_, err = verifyIDtoInt(IssueCategory.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "could not update issue category (id: %s, name: %s) because of malformed input data",
			IssueCategory.ID, IssueCategory.Name)
	}

	apiIssueCategory := *wrapIssueCategory(IssueCategory)

	err = c.redmineAPI.UpdateIssueCategory(apiIssueCategory)
	if err != nil {
		return IssueCategory, errors.Wrapf(err, "error while updating issue category (id: %d, name: %s)", apiIssueCategory.Id, IssueCategory.Name)
	}

	return unwrapIssueCategory(&apiIssueCategory), nil
}

func (c *Client) DeleteIssueCategory(_ context.Context, id string) error {
	idInt, err := verifyIDtoInt(id)
	if err != nil {
		return errors.Wrap(err, "could not delete issue category because of malformed input data")
	}

	err = c.redmineAPI.DeleteIssueCategory(idInt)
	if err != nil {
		return errors.Wrapf(err, "error while deleteting issue category (id: %d)", idInt)
	}

	return nil
}

func wrapIssueCategory(IssueCategory *IssueCategory) *rmapi.IssueCategory {
	apiIssueCategory := &rmapi.IssueCategory{
		Project: rmapi.IdName{Id: IssueCategory.ProjectID},
		Name:    IssueCategory.Name,
	}

	if IssueCategory.ID != "" {
		apiIssueCategory.Id, _ = strconv.Atoi(IssueCategory.ID)
	}

	return apiIssueCategory
}

func unwrapIssueCategory(apiIssueCategory *rmapi.IssueCategory) *IssueCategory {
	IssueCategory := &IssueCategory{
		ID:        strconv.Itoa(apiIssueCategory.Id),
		Name:      apiIssueCategory.Name,
		ProjectID: apiIssueCategory.Project.Id,
	}

	return IssueCategory
}
