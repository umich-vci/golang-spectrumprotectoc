package gospoc

import (
	"context"
	"net/http"
)

const domainsBasePath = "/domains"

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
	NumClients     int    `json:"NUM_CLIENTS,int,omitempty"`
	DefDestArch    string `json:"DEFDESTARCH,string,omitempty"`
	ProvidesBkup   int    `json:"PROVIDES_BKUP,int,omitempty"`
	ServerStatus   int    `json:"SERVERSTATUS,int,omitempty"`
	ScheduleCount  int    `json:"SCHEDULE_COUNT,int,omitempty"`
	DefMC          string `json:"DEF_MC,string,omitempty"`
	DefDestBkup    string `json:"DEFDESTBKUP,string,omitempty"`
	DefDestSPMAN   string `json:"DEFDESTSPMAN,string,omitempty"`
	Name           string `json:"NAME,string,omitempty"`
	MgmtClassCount int    `json:"MGMTCLASS_COUNT,int,omitempty"`
	Server         string `json:"SERVER,string,omitempty"`
	ProvidesArch   int    `json:"PROVIDES_ARCH,int,omitempty"`
	ProvidesSPMG   int    `json:"PROVIDES_SPMG,int,omitempty"`
	Link           string `json:"LINK,string,omitempty"`
	ID             string `json:"ID,string,omitempty"`
	OptSetCount    int    `json:"OPTSET_COUNT,int,omitempty"`
}

type backupDomainsRoot struct {
	Domains      []BackupDomains `json:"domains"`
	DomainsCount int             `json:"domains_count"`
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

	path := serversBasePath + serverName + "/domains/" + domainName + "/details"

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
