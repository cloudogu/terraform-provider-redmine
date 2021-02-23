package redmine

import (
	"context"
	rmapi "github.com/cloudogu/go-redmine"
	"github.com/pkg/errors"
	"strconv"
)

type Issue struct {
	ID            string `json:"id"`
	ProjectID     int    `json:"project_id"`
	TrackerID     int    `json:"tracker_id"`
	Subject       string `json:"subject"`
	Description   string `json:"description"`
	CreatedOn     string `json:"created_on"`
	UpdatedOn     string `json:"updated_on"`
	ParentIssueID int    `json:"parent_issue_id"`
}

func (c *Client) CreateIssue(ctx context.Context, Issue *Issue) (*Issue, error) {
	apiIssue := wrapIssue(Issue)

	actualAPIIssue, err := c.redmineAPI.CreateIssue(*apiIssue)
	if err != nil {
		return nil, errors.Wrapf(err, "error while creating Issue (id %s)", Issue.ID)
	}

	actualIssue := unwrapIssue(actualAPIIssue)

	return actualIssue, nil
}

func (c *Client) ReadIssue(ctx context.Context, id string) (Issue *Issue, err error) {
	idInt, _ := strconv.Atoi(id)
	apiIssue, err := c.redmineAPI.Issue(idInt)
	if err != nil {
		return Issue, errors.Wrapf(err, "error while reading Issue (id %s)", id)
	}

	Issue = unwrapIssue(apiIssue)

	return Issue, nil
}

func (c *Client) UpdateIssue(ctx context.Context, Issue *Issue) (updatedIssue *Issue, err error) {
	apiIssue := *wrapIssue(Issue)

	err = c.redmineAPI.UpdateIssue(apiIssue)
	if err != nil {
		return Issue, errors.Wrapf(err, "error while updating Issue (id %d)", apiIssue.Id)
	}

	Issue = unwrapIssue(&apiIssue)

	return Issue, nil
}

func (c *Client) DeleteIssue(ctx context.Context, name string) error {
	return nil
}

func wrapIssue(issue *Issue) *rmapi.Issue {
	apiIssue := &rmapi.Issue{
		ProjectId:   issue.ProjectID,
		TrackerId:   issue.TrackerID,
		Subject:     issue.Subject,
		Description: issue.Description,
		CreatedOn:   issue.CreatedOn,
		UpdatedOn:   issue.UpdatedOn,
	}

	if issue.ID != "" {
		apiIssue.Id, _ = strconv.Atoi(issue.ID)
	}
	if issue.ParentIssueID != 0 {
		apiIssue.Parent.Id = issue.ParentIssueID
		apiIssue.ParentId = issue.ParentIssueID
	}

	return apiIssue
}

func unwrapIssue(apiIssue *rmapi.Issue) *Issue {
	Issue := &Issue{
		ProjectID:   apiIssue.Project.Id,
		TrackerID:   apiIssue.TrackerId,
		Subject:     apiIssue.Subject,
		Description: apiIssue.Description,
		CreatedOn:   apiIssue.CreatedOn,
		UpdatedOn:   apiIssue.UpdatedOn,
	}

	if apiIssue.Id != 0 {
		Issue.ID = strconv.Itoa(apiIssue.Id)
	}
	if apiIssue.ParentId != 0 {
		Issue.ParentIssueID = apiIssue.Parent.Id
	}

	return Issue
}
