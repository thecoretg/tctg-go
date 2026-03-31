package iru

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// DeviceUser handles the Kandji API quirk where the user field is either an
// object or an empty string when no user is assigned.
type DeviceUser struct {
	Email      string `json:"email"`
	Name       string `json:"name"`
	ID         string `json:"id"`
	IsArchived bool   `json:"is_archived"`
}

func (u *DeviceUser) UnmarshalJSON(b []byte) error {
	if len(b) > 0 && b[0] == '"' {
		return nil // empty string — leave zero value
	}
	type plain DeviceUser
	return json.Unmarshal(b, (*plain)(u))
}

type Device struct {
	DeviceID                   string    `json:"device_id"`
	DeviceName                 string    `json:"device_name"`
	Model                      string    `json:"model"`
	SerialNumber               string    `json:"serial_number"`
	Platform                   string    `json:"platform"`
	OsVersion                  string    `json:"os_version"`
	SupplementalBuildVersion   string    `json:"supplemental_build_version"`
	SupplementalOsVersionExtra string    `json:"supplemental_os_version_extra"`
	LastCheckIn                time.Time `json:"last_check_in"`
	User                       DeviceUser `json:"user"`
	AssetTag        string `json:"asset_tag"`
	BlueprintID     string `json:"blueprint_id"`
	MdmEnabled      bool   `json:"mdm_enabled"`
	AgentInstalled  bool   `json:"agent_installed"`
	IsMissing       bool   `json:"is_missing"`
	IsRemoved       bool   `json:"is_removed"`
	AgentVersion    string `json:"agent_version"`
	FirstEnrollment string `json:"first_enrollment"`
	LastEnrollment  string `json:"last_enrollment"`
	BlueprintName   string `json:"blueprint_name"`
	LostModeStatus  string `json:"lost_mode_status"`
	Tags            []any  `json:"tags"`
}

type DeviceDetails struct {
	General struct {
		DeviceID                   string `json:"device_id"`
		DeviceName                 string `json:"device_name"`
		LastEnrollment             string `json:"last_enrollment"`
		FirstEnrollment            string `json:"first_enrollment"`
		Model                      string `json:"model"`
		Platform                   string `json:"platform"`
		OsVersion                  string `json:"os_version"`
		SupplementalBuildVersion   string `json:"supplemental_build_version"`
		SupplementalOsVersionExtra string `json:"supplemental_os_version_extra"`
		SystemVersion              string `json:"system_version"`
		BootVolume                 string `json:"boot_volume"`
		TimeSinceBoot              string `json:"time_since_boot"`
		LastUser                   string `json:"last_user"`
		AssetTag                   string `json:"asset_tag"`
		AssignedUser               DeviceUser `json:"assigned_user"`
		BlueprintName string `json:"blueprint_name"`
		BlueprintUUID string `json:"blueprint_uuid"`
	} `json:"general"`
	Mdm struct {
		MdmEnabled     string   `json:"mdm_enabled"`
		Supervised     string   `json:"supervised"`
		InstallDate    string   `json:"install_date"`
		LastCheckIn    string   `json:"last_check_in"`
		MdmEnabledUser []string `json:"mdm_enabled_user"`
	} `json:"mdm"`
	ActivationLock struct {
		BypassCodeFailed                     bool `json:"bypass_code_failed"`
		UserActivationLockEnabled            bool `json:"user_activation_lock_enabled"`
		DeviceActivationLockEnabled          bool `json:"device_activation_lock_enabled"`
		ActivationLockAllowedWhileSupervised bool `json:"activation_lock_allowed_while_supervised"`
		ActivationLockSupported              bool `json:"activation_lock_supported"`
	} `json:"activation_lock"`
	Filevault struct {
		FilevaultEnabled         bool   `json:"filevault_enabled"`
		FilevaultRecoverykeyType string `json:"filevault_recoverykey_type"`
		FilevaultPrkEscrowed     bool   `json:"filevault_prk_escrowed"`
		FilevaultNextRotation    string `json:"filevault_next_rotation"`
		FilevaultRegenRequired   bool   `json:"filevault_regen_required"`
	} `json:"filevault"`
	LostMode                  struct{} `json:"lost_mode"`
	AutomatedDeviceEnrollment struct{} `json:"automated_device_enrollment"`
	KandjiAgent               struct {
		AgentInstalled string    `json:"agent_installed"`
		InstallDate    string    `json:"install_date"`
		LastCheckIn    time.Time `json:"last_check_in"`
		AgentVersion   string    `json:"agent_version"`
	} `json:"kandji_agent"`
	HardwareOverview struct {
		ModelName          string `json:"model_name"`
		ModelIdentifier    string `json:"model_identifier"`
		ProcessorName      string `json:"processor_name"`
		ProcessorSpeed     string `json:"processor_speed"`
		NumberOfProcessors string `json:"number_of_processors"`
		TotalNumberOfCores string `json:"total_number_of_cores"`
		Memory             string `json:"memory"`
		Udid               string `json:"udid"`
		SerialNumber       string `json:"serial_number"`
	} `json:"hardware_overview"`
	Volumes []struct {
		Name        string `json:"name"`
		Format      string `json:"format"`
		PercentUsed string `json:"percent_used"`
		Identifier  string `json:"identifier"`
		Capacity    string `json:"capacity"`
		Available   string `json:"available"`
		Encrypted   string `json:"encrypted"`
	} `json:"volumes"`
	Network struct {
		LocalHostname string `json:"local_hostname"`
		MacAddress    string `json:"mac_address"`
		IPAddress     string `json:"ip_address"`
		PublicIP      string `json:"public_ip"`
	} `json:"network"`
	RecoveryInformation struct {
		RecoveryLockEnabled       bool      `json:"recovery_lock_enabled"`
		FirmwarePasswordExist     bool      `json:"firmware_password_exist"`
		FirmwarePasswordPending   bool      `json:"firmware_password_pending"`
		PasswordRotationScheduled time.Time `json:"password_rotation_scheduled"`
		PasswordHasBeenSet        bool      `json:"password_has_been_set"`
	} `json:"recovery_information"`
	Users struct {
		RegularUsers []struct {
			Username string `json:"username"`
			UID      string `json:"uid"`
			Path     string `json:"path"`
			Admin    string `json:"admin"`
			Name     string `json:"name"`
		} `json:"regular_users"`
		SystemUsers []struct {
			Username string `json:"username"`
			UID      string `json:"uid"`
			Path     string `json:"path"`
			Admin    string `json:"admin"`
		} `json:"system_users"`
	} `json:"users"`
	InstalledProfiles []struct {
		Name         string   `json:"name"`
		UUID         string   `json:"uuid"`
		Verified     string   `json:"verified"`
		Identifier   string   `json:"identifier"`
		Organization string   `json:"organization"`
		PayloadTypes []string `json:"payload_types"`
		InstallDate  string   `json:"install_date"`
	} `json:"installed_profiles"`
	AppleBusinessManager struct{} `json:"apple_business_manager"`
	SecurityInformation  struct {
		RemoteDesktopEnabled bool `json:"remote_desktop_enabled"`
	} `json:"security_information"`
	Cellular struct{} `json:"cellular"`
	Tags     []string `json:"tags"`
}

func (c *Client) ListDevices(ctx context.Context) ([]Device, error) {
	res, err := Get[[]Device](ctx, c, "/api/v1/devices", nil)
	if err != nil {
		return nil, err
	}
	return *res, nil
}

func (c *Client) GetDevice(ctx context.Context, deviceID string) (*Device, error) {
	return Get[Device](ctx, c, fmt.Sprintf("/api/v1/devices/%s", deviceID), nil)
}

func (c *Client) GetDeviceDetails(ctx context.Context, deviceID string) (*DeviceDetails, error) {
	return Get[DeviceDetails](ctx, c, fmt.Sprintf("/api/v1/devices/%s/details", deviceID), nil)
}
