package gospoc

import (
	"context"
	"fmt"
	"net/http"
)

const clientsBasePath = "/clients"

// BackupClient contains the elements that make up a backup client
type BackupClient struct {
	Platform string `json:"PLATFORM,string,omitempty"`
	Domain   string `json:"DOMAIN,string,omitempty"`
	Locked   int    `json:"LOCKED,int,omitempty"`
	Server   string `json:"SERVER,int,omitempty"`
	Version  int    `json:"VERSION,int,omitempty"`
	VMOwner  string `json:"VM_OWNER,int,omitempty"`
	GUID     string `json:"GUID,int,omitempty"`
	Link     string `json:"LINK,string,omitempty"`
	Type     int    `json:"TYPE,int,omitempty"`
	VMType   int    `json:"VM_TYPE,int,omitempty"`
	Name     string `json:"NAME,int,omitempty"`
}

type backupClientsRoot struct {
	Clients      []BackupClient `json:"clients"`
	ClientsCount int            `json:"clients_count"`
}

// RegisterClientRequest represents a request to register a backup client node.
type RegisterClientRequest struct {
	Name              string `json:"name"`
	Authentication    string `json:"authentication"`
	Password          string `json:"password"`
	Domain            string `json:"domain"`
	Contact           string `json:"contact"`
	Email             string `json:"email"`
	Schedule          string `json:"schedule,omitempty"`
	OptionSet         string `json:"optionset,omitempty"`
	Deduplication     string `json:"deduplication,omitempty"`
	SSLRequired       string `json:"sslrequired,omitempty"`
	SessionInitiation string `json:"sessioninitiation,omitempty"`
}

type registerClientRequestRoot struct {
	RegisterClient *RegisterClientRequest `json:"registerclient"`
}

type backupSchedule struct {
	Domain   string `json:"domain,string,omitempty"`
	Schedule string `json:"schedule,string,omitempty"`
}

// UpdateClientRequest represents a request to update a backup client node.
type UpdateClientRequest struct {
	Password    string         `json:"password,omitempty"`
	Schedule    backupSchedule `json:"schedule,omitempty"`
	Lock        string         `json:"lock,omitempty"`
	Decommision string         `json:"decommision,omitempty"`
}

// BackupClientAtRisk contains the elements that make up a backup client at risk response
type BackupClientAtRisk struct {
	Server string `json:"SERVER,string,omitempty"`
	AtRisk string `json:"AT_RISK,string,omitempty"`
	Link   string `json:"LINK,string,omitempty"`
	Type   int    `json:"TYPE,int,omitempty"`
	Name   string `json:"NAME,string,omitempty"`
}

type backupClientAtRiskRoot struct {
	ClientAtRisk *BackupClientAtRisk `json:"clientatrisk"`
}

// BackupClients is an interface for interacting with
// IBM Spectrum Protect backup clients
type BackupClients interface {
	AssignSchedule(ctx context.Context, serverName string, clientName string, scheduleDomain string, scheduleName string) (*http.Response, error)
	Decommission(ctx context.Context, serverName string, clientName string) (*http.Response, error)
	DecommissionVM(ctx context.Context, serverName string, clientName string, vmName string) (*http.Response, error)
	Details(ctx context.Context, serverName string, clientName string) (*BackupClient, *http.Response, error)
	List(ctx context.Context) ([]BackupClient, *http.Response, error)
	Lock(ctx context.Context, serverName string, clientName string) (*http.Response, error)
	RegisterNode(ctx context.Context, serverName string, createRequest *RegisterClientRequest) (*http.Response, error)
	Unlock(ctx context.Context, serverName string, clientName string) (*http.Response, error)
	Update(ctx context.Context, serverName string, clientName string, update *UpdateClientRequest) (*http.Response, error)
	UpdatePassword(ctx context.Context, serverName string, clientName string, password string) (*http.Response, error)
}

// BackupClientsOp handles communication with the backup server related methods of the
// IBM Spectrum Protect Operations Center REST API
type BackupClientsOp struct {
	client *Client
}

type assignScheduleBody struct {
	DefineSchedule struct {
		Schedule string `json:"schedule,string"`
	} `json:"defineschedule"`
}

type updatePasswordBody struct {
	UpdatePassword struct {
		Password string `json:"password,string"`
	} `json:"updatepassword"`
}

// List all backup clients
func (s *BackupClientsOp) List(ctx context.Context) ([]BackupClient, *http.Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, clientsBasePath, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(backupClientsRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Clients, resp, err
}

// Details of a specific backup client
func (s *BackupClientsOp) Details(ctx context.Context, serverName string, clientName string) (*BackupClient, *http.Response, error) {
	if serverName == "" {
		return nil, nil, NewArgError("serverName", "cannot be empty")
	}

	if clientName == "" {
		return nil, nil, NewArgError("clientName", "cannot be empty")
	}

	path := serversBasePath + serverName + "/clients/" + clientName

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(clientDetailRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.ClientDetail, resp, err
}

// RegisterNode registers a new backup client
func (s *BackupClientsOp) RegisterNode(ctx context.Context, serverName string, createRequest *RegisterClientRequest) (*http.Response, error) {
	if createRequest == nil {
		return nil, NewArgError("createRequest", "cannot be nil")
	}

	if serverName == "" {
		return nil, NewArgError("serverName", "cannot be empty")
	}

	requestRoot := new(registerClientRequestRoot)
	requestRoot.RegisterClient = createRequest

	path := serversBasePath + serverName + "/clients"

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, requestRoot.RegisterClient)
	if err != nil {
		return nil, err
	}

	root := new(registerClientRequestRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return resp, err
	}

	return resp, err
}

// Lock a backup client (7.1.4 URL scheme only)
func (s *BackupClientsOp) Lock(ctx context.Context, serverName string, clientName string) (*http.Response, error) {
	if serverName == "" {
		return nil, NewArgError("serverName", "cannot be empty")
	}

	if clientName == "" {
		return nil, NewArgError("clientName", "cannot be empty")
	}

	if s.client.Config.URLScheme == "7.1.4" {

		path := serversBasePath + serverName + "/clients/" + clientName + "/lock"

		req, err := s.client.NewRequest(ctx, http.MethodPut, path, nil)
		if err != nil {
			return nil, err
		}

		root := new(registerClientRequestRoot)
		resp, err := s.client.Do(ctx, req, root)
		if err != nil {
			return resp, err
		}

		return resp, err
	}

	update := new(UpdateClientRequest)
	update.Lock = "yes"

	return s.Update(ctx, serverName, clientName, update)

}

// Unlock a backup client with 7.1.4 URL scheme
func (s *BackupClientsOp) Unlock(ctx context.Context, serverName string, clientName string) (*http.Response, error) {
	if serverName == "" {
		return nil, NewArgError("serverName", "cannot be empty")
	}

	if clientName == "" {
		return nil, NewArgError("clientName", "cannot be empty")
	}

	if s.client.Config.URLScheme == "7.1.4" {
		path := serversBasePath + serverName + "/clients/" + clientName + "/unlock"

		req, err := s.client.NewRequest(ctx, http.MethodPut, path, nil)
		if err != nil {
			return nil, err
		}

		root := new(registerClientRequestRoot)
		resp, err := s.client.Do(ctx, req, root)
		if err != nil {
			return resp, err
		}

		return resp, err
	}

	update := new(UpdateClientRequest)
	update.Lock = "no"

	return s.Update(ctx, serverName, clientName, update)

}

// AssignSchedule assigns a new schedule to a backup client
// Note tht scheduleDomain cannot be changed with 7.1.4 URL Scheme and will be ignored if present
func (s *BackupClientsOp) AssignSchedule(ctx context.Context, serverName string, clientName string, scheduleDomain string, scheduleName string) (*http.Response, error) {
	if serverName == "" {
		return nil, NewArgError("serverName", "cannot be empty")
	}

	if clientName == "" {
		return nil, NewArgError("clientName", "cannot be empty")
	}

	if scheduleName == "" {
		return nil, NewArgError("scheduleName", "cannot be empty")
	}

	if s.client.Config.URLScheme == "7.1.4" {
		body := new(assignScheduleBody)
		body.DefineSchedule.Schedule = scheduleName

		path := serversBasePath + serverName + "/clients/" + clientName + "/assignschedule"

		req, err := s.client.NewRequest(ctx, http.MethodPut, path, body)
		if err != nil {
			return nil, err
		}

		root := new(registerClientRequestRoot)
		resp, err := s.client.Do(ctx, req, root)
		if err != nil {
			return resp, err
		}

		return resp, err
	}

	if scheduleDomain == "" {
		return nil, NewArgError("scheduleDomain", "cannot be empty")
	}

	update := new(UpdateClientRequest)
	update.Schedule.Domain = scheduleDomain
	update.Schedule.Schedule = scheduleName

	return s.Update(ctx, serverName, clientName, update)

}

// UpdatePassword changes the password of a backup client
func (s *BackupClientsOp) UpdatePassword(ctx context.Context, serverName string, clientName string, password string) (*http.Response, error) {
	if serverName == "" {
		return nil, NewArgError("serverName", "cannot be empty")
	}

	if clientName == "" {
		return nil, NewArgError("clientName", "cannot be empty")
	}

	if password == "" {
		return nil, NewArgError("password", "cannot be empty")
	}

	if s.client.Config.URLScheme == "7.1.4" {
		body := new(updatePasswordBody)
		body.UpdatePassword.Password = password

		path := serversBasePath + serverName + "/clients/" + clientName + "/passwords"

		req, err := s.client.NewRequest(ctx, http.MethodPut, path, body)
		if err != nil {
			return nil, err
		}

		root := new(registerClientRequestRoot)
		resp, err := s.client.Do(ctx, req, root)
		if err != nil {
			return resp, err
		}

		return resp, err
	}

	update := new(UpdateClientRequest)
	update.Password = password

	return s.Update(ctx, serverName, clientName, update)

}

// Decommission a backup client with the 7.1.4 URL Scheme
func (s *BackupClientsOp) Decommission(ctx context.Context, serverName string, clientName string) (*http.Response, error) {
	if serverName == "" {
		return nil, NewArgError("serverName", "cannot be empty")
	}

	if clientName == "" {
		return nil, NewArgError("clientName", "cannot be empty")
	}

	if s.client.Config.URLScheme == "7.1.4" {
		path := serversBasePath + serverName + "/clients/" + clientName + "/decommissionclient"

		req, err := s.client.NewRequest(ctx, http.MethodPut, path, nil)
		if err != nil {
			return nil, err
		}

		root := new(registerClientRequestRoot)
		resp, err := s.client.Do(ctx, req, root)
		if err != nil {
			return resp, err
		}

		return resp, err
	}

	update := new(UpdateClientRequest)
	update.Decommision = "yes"

	return s.Update(ctx, serverName, clientName, update)
}

// DecommissionVM decommisions a VM backup
func (s *BackupClientsOp) DecommissionVM(ctx context.Context, serverName string, clientName string, vmName string) (*http.Response, error) {
	if serverName == "" {
		return nil, NewArgError("serverName", "cannot be empty")
	}

	if clientName == "" {
		return nil, NewArgError("clientName", "cannot be empty")
	}

	path := serversBasePath + serverName + "/clients/" + clientName + "/vms/" + vmName + "/decommissionclient"

	if s.client.Config.URLScheme == "7.1.4" {
		path = serversBasePath + serverName + "/clients/" + clientName + "/vm/" + vmName + "/decommissionclient"
	}

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, nil)
	if err != nil {
		return nil, err
	}

	root := new(registerClientRequestRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return resp, err
	}

	return resp, err
}

// Update the settings of a backup client
func (s *BackupClientsOp) Update(ctx context.Context, serverName string, clientName string, update *UpdateClientRequest) (*http.Response, error) {
	if s.client.Config.URLScheme == "7.1.4" {
		return nil, fmt.Errorf("The Update method is not supported with the 7.1.4 URL Scheme")
	}

	if serverName == "" {
		return nil, NewArgError("serverName", "cannot be empty")
	}

	if clientName == "" {
		return nil, NewArgError("clientName", "cannot be empty")
	}

	if update == nil {
		return nil, NewArgError("update", "cannot be nil")
	}

	path := serversBasePath + serverName + "/clients/" + clientName

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, update)
	if err != nil {
		return nil, err
	}

	root := new(registerClientRequestRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return resp, err
	}

	return resp, err
}

// AtRisk status of a specific backup client
func (s *BackupClientsOp) AtRisk(ctx context.Context, serverName string, clientName string) (*BackupClientAtRisk, *http.Response, error) {
	if serverName == "" {
		return nil, nil, NewArgError("serverName", "cannot be empty")
	}

	if clientName == "" {
		return nil, nil, NewArgError("clientName", "cannot be empty")
	}

	path := serversBasePath + serverName + "/clients/" + clientName + "/atrisk"

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(backupClientAtRiskRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.ClientAtRisk, resp, err
}
