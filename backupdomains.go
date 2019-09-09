package gospoc

import (
	"context"
	"net/http"
)

const domainsBasePath = "/oc/api/domains"

// BackupDomains is an interface for interacting with
// IBM Spectrum Protect backup domains
type BackupDomains interface {
	List(ctx context.Context) ([]BackupServer, *http.Response, error)
	Get(ctx context.Context, serverName string) (*BackupServer, *http.Response, error)
}

// BackupDomainsOp handles communication with the backup server related methods of the
// IBM Spectrum Protect Operations Center REST API
type BackupDomainsOp struct {
	client *Client
}

// BackupDomain contains the elements that make up a backup domain
type BackupDomain struct {
	NumClients     int    `json:"NUM_CLIENTS"`
	DefDestArch    string `json:"DEFDESTARCH"`
	ProvidesBkup   int    `json:"PROVIDES_BKUP"`
	ServerStatus   int    `json:"SERVERSTATUS"`
	ScheduleCount  int    `json:"SCHEDULE_COUNT"`
	DefMC          string `json:"DEF_MC"`
	DefDestBkup    string `json:"DEFDESTBKUP"`
	DefDestSPMAN   string `json:"DEFDESTSPMAN"`
	Name           string `json:"NAME"`
	MgmtClassCount int    `json:"MGMTCLASS_COUNT"`
	Server         string `json:"SERVER"`
	ProvidesArch   int    `json:"PROVIDES_ARCH"`
	ProvidesSPMG   int    `json:"PROVIDES_SPMG"`
	Link           string `json:"LINK"`
	ID             string `json:"ID"`
	OptSetCount    int    `json:"OPTSET_COUNT"`
}

type backupDomainsRoot struct {
	Domains      []BackupDomains `json:"domains"`
	DomainsCount int             `json:"domains_count,int"`
}

type domainDetailRoot struct {
	DomainDetail *BackupDomain `json:"domaindetail"`
}

// List all backup domains
func (s *BackupDomainsOp) List(ctx context.Context) ([]BackupServer, *http.Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, domainsBasePath, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(backupServersRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Servers, resp, err
}

// Get the details of a specific backup domain
func (s *BackupDomainsOp) Get(ctx context.Context, serverName string, domainName string) (*BackupDomain, *http.Response, error) {
	if serverName == "" {
		return nil, nil, NewArgError("serverName", "cannot be empty")
	}

	path := serversBasePath + "/" + serverName + "/domains/" + domainName + "/details"

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(domainDetailRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.DomainDetail, resp, err
}
