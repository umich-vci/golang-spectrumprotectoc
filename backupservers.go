package gospoc

import (
	"context"
	"net/http"
)

const serversBasePath = "/oc/api/servers"

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
	CatalogUsedSpace     int     `json:"CATALOG_USED_SPACE"`
	SurOcc               float64 `json:"SUR_OCC"`
	SecLastCatalogBackup int     `json:"SEC_LAST_CATALOG_BACKUP"`
	NumAlerts            int     `json:"NUMALERTS"`
	ArchiveLogUsedSpace  int     `json:"ARCHIVELOG_USED_SPACE"`
	NumClients           int     `json:"NUMCLIENTS"`
	HasBackup            int     `json:"HAS_BACKUP"`
	ArchiveLog           string  `json:"ARCHIVELOG"`
	SecUptime            int     `json:"SEC_UPTIME"`
	SurOccTimestamp      string  `json:"SUROCC_TIMESTAMP"`
	ActiveLogUsedSpace   int     `json:"ACTIVELOG_USED_SPACE"`
	HasSpaceMG           int     `json:"HAS_SPACEMG,int"`
	Name                 string  `json:"NAME"`
	FETimestamp          string  `json:"FE_TIMESTAMP"`
	Role                 string  `json:"ROLE"`
	HasArchive           int     `json:"HAS_ARCHIVE"`
	Status               int     `json:"STATUS"`
	Catalog              string  `json:"CATALOG"`
	ActiveLog            string  `json:"ACTIVELOG"`
	Configured           int     `json:"CONFIGURED"`
	Link                 string  `json:"LINK"`
	VRMF                 string  `json:"VRMF"`
	FECapacityTB         float64 `json:"FE_CAPACITY_TB"`
}

type backupServersRoot struct {
	Servers      []BackupServer `json:"servers"`
	ServersCount int            `json:"servers_count"`
}

type serverDetailRoot struct {
	ServerDetail *BackupServer `json:"serverdetail"`
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

	path := serversBasePath + "/" + serverName + "/details"

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
