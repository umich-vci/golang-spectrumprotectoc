package gospoc

import (
	"context"
	"fmt"
	"net/http"
)

const clientsBasePath = "/oc/api/clients"

// BackupClient contains the elements that make up a backup client
type BackupClient struct {
	Platform string `json:"PLATFORM"`
	Domain   string `json:"DOMAIN"`
	Locked   int    `json:"LOCKED"`
	Server   string `json:"SERVER"`
	Version  int    `json:"VERSION"`
	VMOwner  string `json:"VM_OWNER"`
	GUID     string `json:"GUID"`
	Link     string `json:"LINK"`
	Type     int    `json:"TYPE"`
	VMType   int    `json:"VM_TYPE"`
	Name     string `json:"NAME"`
}

type backupClientsRoot struct {
	Clients      []BackupClient `json:"clients"`
	ClientsCount int            `json:"clients_count"`
}

// BackupClientDetail contains the elements that make up backup client details
type BackupClientDetail struct {
	Domain            string `json:"DOMAIN"`
	Contact           string `json:"CONTACT"`
	Locked            string `json:"LOCKED"`
	Deduplication     string `json:"DEDUPLICATION"`
	Email             string `json:"EMAIL"`
	Name              string `json:"NAME"`
	Authentication    string `json:"AUTHENTICATION"`
	SessionInitiation string `json:"SESSIONINITIATION"`
	Decommissioned    string `json:"DECOMMISSIONED"`
	SSLRequired       string `json:"SSLREQUIRED"`
	Link              string `json:"LINK"`
	OptionSet         string `json:"OPTIONSET"`
	SplitLargeObjects string `json:"SPLITLARGEOBJECTS"`
}

type clientDetailRoot struct {
	ClientDetail *BackupClientDetail `json:"clientdetail"`
}

// RegisterClientRequest represents a request to register a backup client node.
type RegisterClientRequest struct {
	Name              string `json:"name,string"`
	Authentication    string `json:"authentication,string"`
	Password          string `json:"password,string"`
	Domain            string `json:"domain,string"`
	Contact           string `json:"contact,string"`
	Email             string `json:"email,string"`
	Schedule          string `json:"schedule,string,omitempty"`
	OptionSet         string `json:"optionset,string,omitempty"`
	Deduplication     string `json:"deduplication,string,omitempty"`
	SSLRequired       string `json:"sslrequired,string,omitempty"`
	SessionInitiation string `json:"sessioninitiation,string,omitempty"`
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
	Password    string         `json:"password,string,omitempty"`
	Schedule    backupSchedule `json:"schedule,omitempty"`
	Lock        string         `json:"lock,string,omitempty"`
	Decommision string         `json:"decommision,string,omitempty"`
}

// BackupClientAtRisk contains the elements that make up a backup client at risk response
type BackupClientAtRisk struct {
	Server string `json:"SERVER"`
	AtRisk string `json:"AT_RISK"`
	Link   string `json:"LINK"`
	Type   int    `json:"TYPE"`
	Name   string `json:"NAME"`
}

type backupClientAtRiskRoot struct {
	ClientAtRisk *BackupClientAtRisk `json:"clientatrisk"`
}

// BackupClientSchedule contains the elements that make up a backup client schedule response
type BackupClientSchedule struct {
	StartTime    string `json:"START_TIME"`
	ServerName   string `json:"SERVER_NAME"`
	ScheduleName string `json:"SCHEDULE_NAME"`
	RunTime      int    `json:"RUNTIME"`
	DomainName   string `json:"DOMAIN_NAME"`
}

type backupClientSchedulesRoot struct {
	ClientSchedules      []BackupClientSchedule `json:"clientschedules"`
	ClientSchedulesCount int                    `json:"clientschedules_count"`
}

// BackupClients is an interface for interacting with
// IBM Spectrum Protect backup clients
type BackupClients interface {
	AssignSchedule(ctx context.Context, serverName string, clientName string, scheduleDomain string, scheduleName string) (*http.Response, error)
	AtRisk(ctx context.Context, serverName string, clientName string) (*BackupClientAtRisk, *http.Response, error)
	Decommission(ctx context.Context, serverName string, clientName string) (*http.Response, error)
	DecommissionVM(ctx context.Context, serverName string, clientName string, vmName string) (*http.Response, error)
	Details(ctx context.Context, serverName string, clientName string) (*BackupClientDetail, *http.Response, error)
	List(ctx context.Context) ([]BackupClient, *http.Response, error)
	Lock(ctx context.Context, serverName string, clientName string) (*http.Response, error)
	RegisterNode(ctx context.Context, serverName string, createRequest *RegisterClientRequest) (*http.Response, error)
	Schedules(ctx context.Context, serverName string, domain string, clientName string) ([]BackupClientSchedule, *http.Response, error)
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
func (s *BackupClientsOp) Details(ctx context.Context, serverName string, clientName string) (*BackupClientDetail, *http.Response, error) {
	if serverName == "" {
		return nil, nil, NewArgError("serverName", "cannot be empty")
	}

	if clientName == "" {
		return nil, nil, NewArgError("clientName", "cannot be empty")
	}

	path := serversBasePath + "/" + serverName + "/clients/" + clientName + "/details"

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

	path := serversBasePath + "/" + serverName + "/clients"

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

		path := serversBasePath + "/" + serverName + "/clients/" + clientName + "/lock"

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
		path := serversBasePath + "/" + serverName + "/clients/" + clientName + "/unlock"

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

		path := serversBasePath + "/" + serverName + "/clients/" + clientName + "/assignschedule"

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

		path := serversBasePath + "/" + serverName + "/clients/" + clientName + "/passwords"

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
		path := serversBasePath + "/" + serverName + "/clients/" + clientName + "/decommissionclient"

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

	path := serversBasePath + "/" + serverName + "/clients/" + clientName + "/vms/" + vmName + "/decommissionclient"

	if s.client.Config.URLScheme == "7.1.4" {
		path = serversBasePath + "/" + serverName + "/clients/" + clientName + "/vm/" + vmName + "/decommissionclient"
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

	path := serversBasePath + "/" + serverName + "/clients/" + clientName

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

	path := serversBasePath + "/" + serverName + "/clients/" + clientName + "/atrisk"

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

// Schedules of a specific backup client
func (s *BackupClientsOp) Schedules(ctx context.Context, serverName string, domain string, clientName string) ([]BackupClientSchedule, *http.Response, error) {
	if serverName == "" {
		return nil, nil, NewArgError("serverName", "cannot be empty")
	}

	if domain == "" {
		return nil, nil, NewArgError("domain", "cannot be empty")
	}

	if clientName == "" {
		return nil, nil, NewArgError("clientName", "cannot be empty")
	}

	path := serversBasePath + "/" + serverName + "/domains/" + domain + "/clients/" + clientName + "/schedules"

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(backupClientSchedulesRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.ClientSchedules, resp, err
}
