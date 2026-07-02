package huntress

import (
	"context"
	"fmt"
)

// Agent represents a Huntress agent installed on a host.
type Agent struct {
	ID                         int64    `json:"id,omitempty"`
	AccountID                  int64    `json:"account_id,omitempty"`
	Arch                       string   `json:"arch,omitempty"`
	CreatedAt                  string   `json:"created_at,omitempty"`
	DomainName                 string   `json:"domain_name,omitempty"`
	EDRVersion                 string   `json:"edr_version,omitempty"`
	ExternalIP                 string   `json:"external_ip,omitempty"`
	Hostname                   string   `json:"hostname,omitempty"`
	DefenderPolicyStatus       string   `json:"defender_policy_status,omitempty"`
	DefenderStatus             string   `json:"defender_status,omitempty"`
	DefenderSubstatus          string   `json:"defender_substatus,omitempty"`
	FirewallStatus             string   `json:"firewall_status,omitempty"`
	TamperProtectionConfigured bool     `json:"tamper_protection_configured,omitempty"`
	TamperProtectionActual     bool     `json:"tamper_protection_actual,omitempty"`
	IPv4Address                string   `json:"ipv4_address,omitempty"`
	LastCallbackAt             string   `json:"last_callback_at,omitempty"`
	LastSurveyAt               string   `json:"last_survey_at,omitempty"`
	MACAddresses               []string `json:"mac_addresses,omitempty"`
	ServicePackMajor           int      `json:"service_pack_major,omitempty"`
	ServicePackMinor           int      `json:"service_pack_minor,omitempty"`
	OrganizationID             int64    `json:"organization_id,omitempty"`
	OS                         string   `json:"os,omitempty"`
	OSBuildVersion             string   `json:"os_build_version,omitempty"`
	OSMajor                    int      `json:"os_major,omitempty"`
	OSMinor                    int      `json:"os_minor,omitempty"`
	OSPatch                    int      `json:"os_patch,omitempty"`
	Platform                   string   `json:"platform,omitempty"`
	SerialNumber               string   `json:"serial_number,omitempty"`
	Tags                       []string `json:"tags,omitempty"`
	UpdatedAt                  string   `json:"updated_at,omitempty"`
	Version                    string   `json:"version,omitempty"`
	VersionNumber              int      `json:"version_number,omitempty"`
	WinBuildNumber             int      `json:"win_build_number,omitempty"`
}

// UpdateAgent is the request body for updating an agent. Tags is the full set
// of user classifications on the host; the supplied list replaces any existing
// tags.
type UpdateAgent struct {
	Tags                       []string `json:"tags,omitempty"`
	TamperProtectionConfigured *bool    `json:"tamper_protection_configured,omitempty"`
}

// IsolateAgent is the request body for isolating an agent's host.
type IsolateAgent struct {
	Reason          string `json:"reason,omitempty"`
	StrictIsolation *bool  `json:"strict_isolation,omitempty"`
}

// AgentsResponse is a page of agents.
type AgentsResponse struct {
	Agents     []Agent    `json:"agents"`
	Pagination Pagination `json:"pagination"`
}

type agentEnvelope struct {
	Agent Agent `json:"agent"`
}

// ListAgents returns a single page of agents.
func (c *Client) ListAgents(ctx context.Context, params map[string]string) (*AgentsResponse, error) {
	result, err := get[AgentsResponse](ctx, c, endpointURL("agents"), params)
	if err != nil {
		return nil, fmt.Errorf("list agents: %w", err)
	}
	return result, nil
}

// GetAgent returns a single agent by ID.
func (c *Client) GetAgent(ctx context.Context, id int64) (*Agent, error) {
	result, err := get[agentEnvelope](ctx, c, endpointURL(fmt.Sprintf("agents/%d", id)), nil)
	if err != nil {
		return nil, fmt.Errorf("get agent: %w", err)
	}
	return &result.Agent, nil
}

// UpdateAgent updates an agent's tags or tamper-protection setting.
func (c *Client) UpdateAgent(ctx context.Context, id int64, body UpdateAgent) (*Agent, error) {
	result, err := patch[Agent](ctx, c, endpointURL(fmt.Sprintf("agents/%d", id)), body)
	if err != nil {
		return nil, fmt.Errorf("update agent: %w", err)
	}
	return result, nil
}

// UninstallAgent schedules an agent for uninstall.
func (c *Client) UninstallAgent(ctx context.Context, id int64) error {
	if err := del(ctx, c, endpointURL(fmt.Sprintf("agents/%d", id))); err != nil {
		return fmt.Errorf("uninstall agent: %w", err)
	}
	return nil
}

// IsolateAgent isolates an agent's host from the network.
func (c *Client) IsolateAgent(ctx context.Context, id int64, body IsolateAgent) (*Agent, error) {
	result, err := post[Agent](ctx, c, endpointURL(fmt.Sprintf("agents/%d/isolation", id)), body)
	if err != nil {
		return nil, fmt.Errorf("isolate agent: %w", err)
	}
	return result, nil
}

// ReleaseAgentIsolation releases an agent's host from isolation.
func (c *Client) ReleaseAgentIsolation(ctx context.Context, id int64) error {
	if err := del(ctx, c, endpointURL(fmt.Sprintf("agents/%d/isolation", id))); err != nil {
		return fmt.Errorf("release agent isolation: %w", err)
	}
	return nil
}

// ListAccountAgents returns a single page of agents for an account (Reseller
// credentials only).
func (c *Client) ListAccountAgents(ctx context.Context, accountID int64, params map[string]string) (*AgentsResponse, error) {
	result, err := get[AgentsResponse](ctx, c, endpointURL(fmt.Sprintf("accounts/%d/agents", accountID)), params)
	if err != nil {
		return nil, fmt.Errorf("list account agents: %w", err)
	}
	return result, nil
}

// GetAccountAgent returns a single agent for an account (Reseller credentials only).
func (c *Client) GetAccountAgent(ctx context.Context, accountID, id int64) (*Agent, error) {
	result, err := get[agentEnvelope](ctx, c, endpointURL(fmt.Sprintf("accounts/%d/agents/%d", accountID, id)), nil)
	if err != nil {
		return nil, fmt.Errorf("get account agent: %w", err)
	}
	return &result.Agent, nil
}

// UpdateAccountAgent updates an agent within an account (Reseller credentials only).
func (c *Client) UpdateAccountAgent(ctx context.Context, accountID, id int64, body UpdateAgent) (*Agent, error) {
	result, err := patch[Agent](ctx, c, endpointURL(fmt.Sprintf("accounts/%d/agents/%d", accountID, id)), body)
	if err != nil {
		return nil, fmt.Errorf("update account agent: %w", err)
	}
	return result, nil
}

// UninstallAccountAgent schedules an account's agent for uninstall (Reseller
// credentials only).
func (c *Client) UninstallAccountAgent(ctx context.Context, accountID, id int64) error {
	if err := del(ctx, c, endpointURL(fmt.Sprintf("accounts/%d/agents/%d", accountID, id))); err != nil {
		return fmt.Errorf("uninstall account agent: %w", err)
	}
	return nil
}

// IsolateAccountAgent isolates an account's agent host (Reseller credentials only).
func (c *Client) IsolateAccountAgent(ctx context.Context, accountID, id int64, body IsolateAgent) (*Agent, error) {
	result, err := post[Agent](ctx, c, endpointURL(fmt.Sprintf("accounts/%d/agents/%d/isolation", accountID, id)), body)
	if err != nil {
		return nil, fmt.Errorf("isolate account agent: %w", err)
	}
	return result, nil
}

// ReleaseAccountAgentIsolation releases an account's agent host from isolation
// (Reseller credentials only).
func (c *Client) ReleaseAccountAgentIsolation(ctx context.Context, accountID, id int64) error {
	if err := del(ctx, c, endpointURL(fmt.Sprintf("accounts/%d/agents/%d/isolation", accountID, id))); err != nil {
		return fmt.Errorf("release account agent isolation: %w", err)
	}
	return nil
}
