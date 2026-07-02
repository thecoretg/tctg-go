package huntress

import (
	"context"
	"fmt"
)

// ReportExternalService is an external service listed in a summary report.
type ReportExternalService struct {
	Name  string `json:"name,omitempty"`
	Risky bool   `json:"risky,omitempty"`
}

// ReportIncidentRemediation is a remediation summary within a report incident.
type ReportIncidentRemediation struct {
	Type    string `json:"type,omitempty"`
	Subtype string `json:"subtype,omitempty"`
}

// ReportIncident is an incident summarized within a summary report.
type ReportIncident struct {
	ID           int64                       `json:"id,omitempty"`
	Severity     string                      `json:"severity,omitempty"`
	SentAt       string                      `json:"sent_at,omitempty"`
	EventSummary string                      `json:"event_summary,omitempty"`
	Body         string                      `json:"body,omitempty"`
	Remediations []ReportIncidentRemediation `json:"remediations,omitempty"`
}

// SummaryReport represents a Huntress summary report.
type SummaryReport struct {
	ID                                int64                   `json:"id,omitempty"`
	AgentsCount                       int64                   `json:"agents_count,omitempty"`
	AllowedExclusionsCount            int64                   `json:"allowed_exclusions_count,omitempty"`
	AnalystName                       string                  `json:"analyst_name,omitempty"`
	AnalystNote                       string                  `json:"analyst_note,omitempty"`
	AnalystThreats                    []string                `json:"analyst_threats,omitempty"`
	AnalystTitle                      string                  `json:"analyst_title,omitempty"`
	AntivirusExclusionsCount          int64                   `json:"antivirus_exclusions_count,omitempty"`
	AutorunEvents                     int64                   `json:"autorun_events,omitempty"`
	AutorunSignalsDetected            int64                   `json:"autorun_signals_detected,omitempty"`
	AutorunSignalsReviewed            int64                   `json:"autorun_signals_reviewed,omitempty"`
	AutorunsReviewed                  int64                   `json:"autoruns_reviewed,omitempty"`
	BlockedMalwareCount               int64                   `json:"blocked_malware_count,omitempty"`
	CreatedAt                         string                  `json:"created_at,omitempty"`
	DeployedCanariesCount             int64                   `json:"deployed_canaries_count,omitempty"`
	EventsAnalyzed                    int64                   `json:"events_analyzed,omitempty"`
	ExternalIPsCount                  int64                   `json:"external_ips_count,omitempty"`
	ExternalPortsCount                int64                   `json:"external_ports_count,omitempty"`
	ExternalServices                  []ReportExternalService `json:"external_services,omitempty"`
	FirewallDisabledCount             int64                   `json:"firewall_disabled_count,omitempty"`
	FirewallDisabledWithConflictCount int64                   `json:"firewall_disabled_with_conflict_count,omitempty"`
	FirewallEnabledCount              int64                   `json:"firewall_enabled_count,omitempty"`
	FirewallEnabledWithConflictCount  int64                   `json:"firewall_enabled_with_conflict_count,omitempty"`
	GlobalThreatsNote                 string                  `json:"global_threats_note,omitempty"`
	HostProcessesAnalyzed             int64                   `json:"host_processes_analyzed,omitempty"`
	IncidentIndicatorCounts           map[string]any          `json:"incident_indicator_counts,omitempty"`
	IncidentLog                       []string                `json:"incident_log,omitempty"`
	IncidentProductCounts             map[string]any          `json:"incident_product_counts,omitempty"`
	IncidentSeverityCounts            map[string]any          `json:"incident_severity_counts,omitempty"`
	IncidentsReported                 int64                   `json:"incidents_reported,omitempty"`
	IncidentsResolved                 int64                   `json:"incidents_resolved,omitempty"`
	InvestigatedMAVDetectionCount     int64                   `json:"investigated_mav_detection_count,omitempty"`
	InvestigationsCompleted           int64                   `json:"investigations_completed,omitempty"`
	ITDREntities                      int64                   `json:"itdr_entities,omitempty"`
	ITDREvents                        int64                   `json:"itdr_events,omitempty"`
	ITDRIncidentsReported             int64                   `json:"itdr_incidents_reported,omitempty"`
	ITDRInvestigationsCompleted       int64                   `json:"itdr_investigations_completed,omitempty"`
	ITDRSignals                       int64                   `json:"itdr_signals,omitempty"`
	ITDRBillableIdentityCount         int64                   `json:"itdr_billable_identity_count,omitempty"`
	ITDRNonBillableIdentityCount      int64                   `json:"itdr_non_billable_identity_count,omitempty"`
	ITDRLicenseDistribution           map[string]any          `json:"itdr_license_distribution,omitempty"`
	ITDRUsageLocations                []string                `json:"itdr_usage_locations,omitempty"`
	LinuxAgentCount                   int64                   `json:"linux_agent_count,omitempty"`
	MacOSAgentCount                   int64                   `json:"macos_agent_count,omitempty"`
	MacOSAgents                       bool                    `json:"macos_agents,omitempty"`
	MAVIncidentReportCount            int64                   `json:"mav_incident_report_count,omitempty"`
	NewExclusionsCount                int64                   `json:"new_exclusions_count,omitempty"`
	OnlyMacOSAgents                   bool                    `json:"only_macos_agents,omitempty"`
	OrganizationID                    int64                   `json:"organization_id,omitempty"`
	Period                            string                  `json:"period,omitempty"`
	PowerfulApplicationCount          int64                   `json:"powerful_application_count,omitempty"`
	PotentialThreatIndicators         int64                   `json:"potential_threat_indicators,omitempty"`
	ProcessDetections                 int64                   `json:"process_detections,omitempty"`
	ProcessDetectionsReported         int64                   `json:"process_detections_reported,omitempty"`
	ProcessDetectionsReviewed         int64                   `json:"process_detections_reviewed,omitempty"`
	ProtectedProfilesCount            int64                   `json:"protected_profiles_count,omitempty"`
	RansomwareNote                    string                  `json:"ransomware_note,omitempty"`
	RiskyExclusionsRemovedCount       int64                   `json:"risky_exclusions_removed_count,omitempty"`
	RiskyServicesCount                int64                   `json:"risky_services_count,omitempty"`
	RogueAppIncidents                 []ReportIncident        `json:"rogue_app_incidents,omitempty"`
	ServersAgentCount                 int64                   `json:"servers_agent_count,omitempty"`
	ShadowWorkflowIncidents           []ReportIncident        `json:"shadow_workflow_incidents,omitempty"`
	SIEMIncidentsReported             int64                   `json:"siem_incidents_reported,omitempty"`
	SIEMIngestedLogs                  int64                   `json:"siem_ingested_logs,omitempty"`
	SIEMInvestigationsCompleted       int64                   `json:"siem_investigations_completed,omitempty"`
	SIEMSignals                       int64                   `json:"siem_signals,omitempty"`
	SIEMTotalLogs                     int64                   `json:"siem_total_logs,omitempty"`
	SignalsDetected                   int64                   `json:"signals_detected,omitempty"`
	SignalsInvestigated               int64                   `json:"signals_investigated,omitempty"`
	TopIncidentAVThreats              []string                `json:"top_incident_av_threats,omitempty"`
	TopIncidentHosts                  []string                `json:"top_incident_hosts,omitempty"`
	TotalEntities                     int64                   `json:"total_entities,omitempty"`
	TotalMAVDetectionCount            int64                   `json:"total_mav_detection_count,omitempty"`
	Type                              string                  `json:"type,omitempty"`
	UnwantedAccessIncidents           []ReportIncident        `json:"unwanted_access_incidents,omitempty"`
	UpdatedAt                         string                  `json:"updated_at,omitempty"`
	URL                               string                  `json:"url,omitempty"`
	WindowsAgentCount                 int64                   `json:"windows_agent_count,omitempty"`
	WeakApplicationCount              int64                   `json:"weak_application_count,omitempty"`
	WindowsAgents                     bool                    `json:"windows_agents,omitempty"`
}

// ReportsResponse is a page of summary reports.
type ReportsResponse struct {
	Reports    []SummaryReport `json:"reports"`
	Pagination Pagination      `json:"pagination"`
}

// ListReports returns a single page of summary reports.
func (c *Client) ListReports(ctx context.Context, params map[string]string) (*ReportsResponse, error) {
	result, err := get[ReportsResponse](ctx, c, endpointURL("reports"), params)
	if err != nil {
		return nil, fmt.Errorf("list reports: %w", err)
	}
	return result, nil
}

// GetReport returns a single summary report by ID.
func (c *Client) GetReport(ctx context.Context, id int64) (*SummaryReport, error) {
	result, err := get[SummaryReport](ctx, c, endpointURL(fmt.Sprintf("reports/%d", id)), nil)
	if err != nil {
		return nil, fmt.Errorf("get report: %w", err)
	}
	return result, nil
}

// ListAccountReports returns a single page of summary reports for an account
// (Reseller credentials only).
func (c *Client) ListAccountReports(ctx context.Context, accountID int64, params map[string]string) (*ReportsResponse, error) {
	result, err := get[ReportsResponse](ctx, c, endpointURL(fmt.Sprintf("accounts/%d/reports", accountID)), params)
	if err != nil {
		return nil, fmt.Errorf("list account reports: %w", err)
	}
	return result, nil
}

// GetAccountReport returns a single summary report for an account (Reseller
// credentials only).
func (c *Client) GetAccountReport(ctx context.Context, accountID, id int64) (*SummaryReport, error) {
	result, err := get[SummaryReport](ctx, c, endpointURL(fmt.Sprintf("accounts/%d/reports/%d", accountID, id)), nil)
	if err != nil {
		return nil, fmt.Errorf("get account report: %w", err)
	}
	return result, nil
}
