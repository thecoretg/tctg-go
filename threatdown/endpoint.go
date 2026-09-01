package threatdown

import (
	"context"
	"fmt"
	"maps"
)

type endpointSearchResp struct {
	Endpoints  []Endpoint `json:"endpoints"`
	NextCursor string     `json:"next_cursor"`
}

type Endpoint struct {
	Link             string `json:"link"`
	ProtectionStatus string `json:"protection_status"`
	DisplayName      string `json:"display_name"`
	Agent            struct {
		SerialNumber              string `json:"serial_number"`
		HasAlerts                 bool   `json:"has_alerts"`
		IsSoftwareUpdateAvailable bool   `json:"is_software_update_available"`
		LastUser                  string `json:"last_user"`
		At                        string `json:"at"`
		MachineID                 string `json:"machine_id"`
		AccountID                 string `json:"account_id"`
		GroupID                   string `json:"group_id"`
		OsInfo                    struct {
			OsType         string `json:"os_type"`
			OsVersion      string `json:"os_version"`
			OsPlatform     string `json:"os_platform"`
			OsArchitecture string `json:"os_architecture"`
			OsReleaseName  string `json:"os_release_name"`
		} `json:"os_info"`
		DomainName             string `json:"domain_name"`
		HostName               string `json:"host_name"`
		FullyQualifiedHostName string `json:"fully_qualified_host_name"`
		EngineVersion          string `json:"engine_version"`
		PolicyEtag             string `json:"policy_etag"`
		Version                int    `json:"version"`
	} `json:"agent"`
	Machine struct {
		ID                      string   `json:"id"`
		Job                     struct{} `json:"job"`
		Account                 struct{} `json:"account"`
		Online                  bool     `json:"online"`
		AccountID               string   `json:"account_id"`
		GroupID                 string   `json:"group_id"`
		RootGroupID             string   `json:"root_group_id"`
		GroupName               string   `json:"group_name"`
		PolicyID                string   `json:"policy_id"`
		PolicyName              string   `json:"policy_name"`
		LastDaySeen             string   `json:"last_day_seen"`
		LastActive              string   `json:"last_active"`
		Isolated                bool     `json:"isolated"`
		ScanAgeDays             int      `json:"scan_age_days"`
		SuspiciousActivityCount int      `json:"suspicious_activity_count"`
		InfectionCount          int      `json:"infection_count"`
		RebootRequired          int      `json:"reboot_required"`
		LastScannedAt           string   `json:"last_scanned_at"`
		IsDeleted               bool     `json:"is_deleted"`
		Version                 int      `json:"version"`
		IsDefaultGroup          bool     `json:"is_default_group"`
		DocumentID              string   `json:"document_id"`
	} `json:"machine"`
	MachineVersion int `json:"machineVersion"`
}

func (c *Client) SearchEndpoints(ctx context.Context, nebulaAccountID string, query map[string]string) ([]Endpoint, error) {
	url := endpointURLV1(fmt.Sprintf("accounts/%s/endpoints", nebulaAccountID))
	var all []Endpoint
	cursor := ""
	for {
		body := make(map[string]string, len(query)+1)
		maps.Copy(body, query)
		if cursor != "" {
			body["next_cursor"] = cursor
		}

		result, err := c.Post[endpointSearchResp](ctx, url, body)
		if err != nil {
			return nil, fmt.Errorf("search endpoints: %w", err)
		}

		all = append(all, result.Endpoints...)
		if result.NextCursor == "" {
			return all, nil
		}
		cursor = result.NextCursor
	}
}
