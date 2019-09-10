package gospoc

import (
	"context"
	"net/http"
)

const (
	cliBasePath    = "/oc/api/cli"
	cliContentType = "text/plain"
)

// CLI is an interface for interacting with
// IBM Spectrum Protect CLI
type CLI interface {
	IssueCommand(ctx context.Context, serverName string, command string) (*http.Response, error)
	IssueConfirmCommand(ctx context.Context, serverName string, command string) (*http.Response, error)
}

// CLIOp handles communication with the cli related methods of the
// IBM Spectrum Protect Operations Center REST API
type CLIOp struct {
	client *Client
}

// IssueCommand issues a TSM Command
func (s *CLIOp) IssueCommand(ctx context.Context, serverName string, command string) (*http.Response, error) {
	path := cliBasePath + "/issueCommand"

	if serverName != "" {
		path = path + "/" + serverName
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, command)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", cliContentType)

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, err
}

// IssueConfirmCommand issues a confirmed TSM Command
func (s *CLIOp) IssueConfirmCommand(ctx context.Context, serverName string, command string) (*http.Response, error) {
	path := cliBasePath + "/issueConfirmedCommand"

	if serverName != "" {
		path = path + "/" + serverName
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, command)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", cliContentType)

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, err
}
