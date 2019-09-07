package gospoc

import (
	"context"
	"net/http"
)

const serversBasePath = "/servers"

// BackupServers is an interface for interacting with
// IBM Spectrum Protect backup servers
type BackupServers interface {
	List(ctx context.Context) ([]BackupServer, *http.Response, error)
	Get(ctx context.Context, serverName string) (*BackupServer, *http.Response, error)
}

// BackupServersOp handles communication with the backup server related methods of the
// IBM Spectrum Protect Operations Center REST API
type BackupServersOp struct {
	client *Client
}

// BackupServer contains the elements that make up a backup server
type BackupServer struct {
	CatalogUsedSpace     int     `json:"CATALOG_USED_SPACE,int,omitempty"`
	SurOcc               float64 `json:"SUR_OCC,float64,omitempty"`
	SecLastCatalogBackup int     `json:"SEC_LAST_CATALOG_BACKUP,int,omitempty"`
	NumAlerts            int     `json:"NUMALERTS,int,omitempty"`
	ArchiveLogUsedSpace  int     `json:"ARCHIVELOG_USED_SPACE,int,omitempty"`
	NumClients           int     `json:"NUMCLIENTS,int,omitempty"`
	HasBackup            int     `json:"HAS_BACKUP,int,omitempty"`
	ArchiveLog           string  `json:"ARCHIVELOG,string,omitempty"`
	SecUptime            int     `json:"SEC_UPTIME,int,omitempty"`
	SurOccTimestamp      string  `json:"SUROCC_TIMESTAMP,string,omitempty"`
	ActiveLogUsedSpace   int     `json:"ACTIVELOG_USED_SPACE,int,omitempty"`
	HasSpaceMG           int     `json:"HAS_SPACEMG,int,omitempty"`
	Name                 string  `json:"NAME,string,omitempty"`
	FETimestamp          string  `json:"FE_TIMESTAMP,string,omitempty"`
	Role                 string  `json:"ROLE,string,omitempty"`
	HasArchive           int     `json:"HAS_ARCHIVE,int,omitempty"`
	Status               int     `json:"STATUS,int,omitempty"`
	Catalog              string  `json:"CATALOG,string,omitempty"`
	ActiveLog            string  `json:"ACTIVELOG,string,omitempty"`
	Configured           int     `json:"CONFIGURED,int,omitempty"`
	Link                 string  `json:"LINK,string,omitempty"`
	VRMF                 string  `json:"VRMF,string,omitempty"`
	FECapacityTB         float64 `json:"FE_CAPACITY_TB,float64,omitempty"`
}

type backupServersRoot struct {
	Servers      []BackupServer `json:"servers"`
	ServersCount int            `json:"servers_count"`
}

type serverDetailRoot struct {
	ServerDetail *BackupServer `json:"serverdetail"`
}

type clientDetailRoot struct {
	ClientDetail *BackupClient `json:"clientdetail"`
}

// List all backup servers
func (s *BackupServersOp) List(ctx context.Context) ([]BackupServer, *http.Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, serversBasePath, nil)
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

// Get the details of a specific backup server
func (s *BackupServersOp) Get(ctx context.Context, serverName string) (*BackupServer, *http.Response, error) {
	if serverName == "" {
		return nil, nil, NewArgError("serverName", "cannot be empty")
	}

	path := serversBasePath + serverName + "/details"

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(serverDetailRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.ServerDetail, resp, err
}
