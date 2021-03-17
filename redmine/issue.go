package redmine

import (
	"context"
	"fmt"
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
	ParentIssueID int    `json:"parent_issue_id"`
	CreatedOn     string `json:"created_on"`
	UpdatedOn     string `json:"updated_on"`
}

func (i *Issue) String() string {
	return fmt.Sprintf("issue{ID=%s,ProjectID=%d,TrackerID=%d,Subject=%s,Description=%s,ParentIssueID=%d,CreatedOn=%s,UpdatedOn=%s}",
		i.ID, i.ProjectID, i.TrackerID, i.Subject, i.Description, i.ParentIssueID, i.CreatedOn, i.UpdatedOn)
}

func (c *Client) CreateIssue(ctx context.Context, issue *Issue) (*Issue, error) {
	apiIssue := wrapIssue(issue)

	createdAPIIssue, err := c.redmineAPI.CreateIssue(*apiIssue)
	if err != nil {
		return nil, errors.Wrapf(err, "error while creating issue (project id: %d, subject: %s)", issue.ProjectID, issue.Subject)
	}

	actualIssue := unwrapIssue(createdAPIIssue)

	return actualIssue, nil
}

func (c *Client) ReadIssue(ctx context.Context, id string) (Issue *Issue, err error) {
	idInt, err := verifyIDtoInt(id)
	if err != nil {
		return nil, errors.Wrap(err, "could not read issue because of malformed input data")
	}

	apiIssue, err := c.redmineAPI.Issue(idInt)
	if err != nil {
		return Issue, errors.Wrapf(err, "error while reading issue (id: %d)", idInt)
	}

	Issue = unwrapIssue(apiIssue)

	return Issue, nil
}

func (c *Client) UpdateIssue(ctx context.Context, issue *Issue) (updatedIssue *Issue, err error) {
	_, err = verifyIDtoInt(issue.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "could not update issue (id: %s, subject: %s) because of malformed input data",
			issue.ID, issue.Subject)
	}

	apiIssue := *wrapIssue(issue)

	err = c.redmineAPI.UpdateIssue(apiIssue)
	if err != nil {
		return issue, errors.Wrapf(err, "error while updating issue (id: %d, subject: %s)", apiIssue.Id, issue.Subject)
	}

	issue = unwrapIssue(&apiIssue)

	return issue, nil
}

func (c *Client) DeleteIssue(ctx context.Context, id string) error {
	idInt, err := verifyIDtoInt(id)
	if err != nil {
		return errors.Wrap(err, "could not delete issue because of malformed input data")
	}

	err = c.redmineAPI.DeleteIssue(idInt)
	if err != nil {
		return errors.Wrapf(err, "error while deleteting issue (id: %d)", idInt)
	}

	return nil
}

func wrapIssue(issue *Issue) *rmapi.Issue {
	apiIssue := &rmapi.Issue{
		ProjectId:   issue.ProjectID,
		Project:     &rmapi.IdName{Id: issue.ProjectID},
		TrackerId:   issue.TrackerID,
		Tracker:     &rmapi.IdName{Id: issue.TrackerID},
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
	issue := &Issue{
		Subject:     apiIssue.Subject,
		Description: apiIssue.Description,
		CreatedOn:   apiIssue.CreatedOn,
		UpdatedOn:   apiIssue.UpdatedOn,
	}

	if apiIssue.Id != 0 {
		issue.ID = strconv.Itoa(apiIssue.Id)
	}
	if apiIssue.Parent != nil {
		issue.ParentIssueID = apiIssue.Parent.Id
	}
	if apiIssue.Project != nil {
		issue.ProjectID = apiIssue.Project.Id
	}

	if apiIssue.Tracker != nil {
		issue.TrackerID = apiIssue.Tracker.Id
	}

	return issue
}
