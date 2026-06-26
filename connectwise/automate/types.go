package automate

import "encoding/json"

// Company is an Automate client/company (LabTech.Models.Client). Automate names
// this entity a "Client"; it is named Company here to match the connectwise/psa
// package and to avoid colliding with the connection Client type.
type Company struct {
	ID                            string     `json:"Id,omitempty"`
	Name                          string     `json:"Name,omitempty"`
	Company                       string     `json:"Company,omitempty"`
	FirstName                     string     `json:"FirstName,omitempty"`
	LastName                      string     `json:"LastName,omitempty"`
	Address1                      string     `json:"Address1,omitempty"`
	Address2                      string     `json:"Address2,omitempty"`
	City                          string     `json:"City,omitempty"`
	State                         string     `json:"State,omitempty"`
	ZipCode                       string     `json:"ZipCode,omitempty"`
	PhoneNumber                   string     `json:"PhoneNumber,omitempty"`
	FaxNumber                     string     `json:"FaxNumber,omitempty"`
	Comment                       string     `json:"Comment,omitempty"`
	Country                       string     `json:"Country,omitempty"`
	ExternalId                    string     `json:"ExternalId,omitempty"`
	UsesInHouseSupportStaff       bool       `json:"UsesInHouseSupportStaff,omitempty"`
	NewTicketNotificationEmail    string     `json:"NewTicketNotificationEmail,omitempty"`
	IsHiddenFromAllInclusiveGroup bool       `json:"IsHiddenFromAllInclusiveGroup,omitempty"`
	Locations                     []Location `json:"Locations,omitempty"`
}

// Location is an Automate client location (LabTech.Models.Location for
// POST/PUT request bodies, Automate.Api.Domain.Contracts.Clients.Location for
// GET responses). Deeply nested objects (contact, router, deployment
// templates, etc.) are kept as raw JSON in this scaffold; type them out as
// needed.
//
// Client must be set when creating a location (POST) so the new location is
// associated with a client; setting only Client.ID is sufficient.
type Location struct {
	ID                  int      `json:"Id,omitempty"`
	LocationId          int      `json:"LocationId,omitempty"`
	Client              *Company `json:"Client,omitempty"`
	Name                string   `json:"Name,omitempty"`
	Address1            string   `json:"Address1,omitempty"`
	Address2            string   `json:"Address2,omitempty"`
	City                string   `json:"City,omitempty"`
	State               string   `json:"State,omitempty"`
	ZipCode             string   `json:"ZipCode,omitempty"`
	Country             string   `json:"Country,omitempty"`
	PhoneNumber         string   `json:"PhoneNumber,omitempty"`
	FaxNumber           string   `json:"FaxNumber,omitempty"`
	Comments            string   `json:"Comments,omitempty"`
	RouterPort          int      `json:"RouterPort,omitempty"`
	ScriptDrive         string   `json:"ScriptDrive,omitempty"`
	ScriptUsername      string   `json:"ScriptUsername,omitempty"`
	ScriptPassword      string   `json:"ScriptPassword,omitempty"`
	ScriptRouterAddress string   `json:"ScriptRouterAddress,omitempty"`
	ScriptExtra1        string   `json:"ScriptExtra1,omitempty"`
	ScriptExtra2        string   `json:"ScriptExtra2,omitempty"`
	ProbeId             int      `json:"ProbeId,omitempty"`
	ExternalId          int      `json:"ExternalId,omitempty"`

	Contact                json.RawMessage `json:"Contact,omitempty"`
	Router                 json.RawMessage `json:"Router,omitempty"`
	DeploymentTemplate     json.RawMessage `json:"DeploymentTemplate,omitempty"`
	MaintenanceWindow      json.RawMessage `json:"MaintenanceWindow,omitempty"`
	DefaultDeploymentGroup json.RawMessage `json:"DefaultDeploymentGroup,omitempty"`
	DefaultDeploymentLogin json.RawMessage `json:"DefaultDeploymentLogin,omitempty"`

	// ExtraFields is returned on GET; manage EDFs via the dedicated
	// /Locations/{id}/ExtraFields endpoints (see locations.go).
	ExtraFields []ExtraField `json:"ExtraFields,omitempty"`
}

// ExtraField is a custom field on a client
// (Automate.Api.Domain.Contracts.ExtraFields.ExtraField). The settings/format
// sub-objects are kept as raw JSON in this scaffold.
type ExtraField struct {
	TargetId               int    `json:"TargetId,omitempty"`
	ExtraFieldDefinitionId int    `json:"ExtraFieldDefinitionId,omitempty"`
	Title                  string `json:"Title,omitempty"`
	Section                string `json:"Section,omitempty"`
	Tooltip                string `json:"Tooltip,omitempty"`
	IsReadOnly             bool   `json:"IsReadOnly,omitempty"`
	IsEncrypted            bool   `json:"IsEncrypted,omitempty"`
	IsDefaultValue         bool   `json:"IsDefaultValue,omitempty"`

	DisplayFormat     json.RawMessage `json:"DisplayFormat,omitempty"`
	Location          json.RawMessage `json:"Location,omitempty"`
	TitleFormat       json.RawMessage `json:"TitleFormat,omitempty"`
	TextFieldSettings json.RawMessage `json:"TextFieldSettings,omitempty"`
	DropdownSettings  json.RawMessage `json:"DropdownSettings,omitempty"`
	CheckboxSettings  json.RawMessage `json:"CheckboxSettings,omitempty"`
}

// Document is a client document (LabTech.Models.Document).
type Document struct {
	ID           int    `json:"Id,omitempty"`
	ClientId     int    `json:"ClientId,omitempty"`
	Name         string `json:"Name,omitempty"`
	Size         int    `json:"Size,omitempty"`
	Data         string `json:"Data,omitempty"`
	LastUser     string `json:"LastUser,omitempty"`
	LastEditDate Time   `json:"LastEditDate,omitzero"`
}

// ManagedLicense is a managed license on a client (LabTech.Models.ManagedLicense).
type ManagedLicense struct {
	ID             int         `json:"Id,omitempty"`
	ClientId       int         `json:"ClientId,omitempty"`
	Name           string      `json:"Name,omitempty"`
	SearchString   string      `json:"SearchString,omitempty"`
	LicenseCount   int         `json:"LicenseCount,omitempty"`
	InstalledCount int         `json:"InstalledCount,omitempty"`
	ProductKey     *ProductKey `json:"ProductKey,omitempty"`
}

// Computer is an Automate managed computer (LabTech.Models.Computer). The
// nested entity references (location, client, contact, etc.) and complex
// collections (logged-in users, tickets, groups) are kept as raw JSON in this
// scaffold; type them out as needed.
type Computer struct {
	ID                       string  `json:"Id,omitempty"`
	ComputerName             string  `json:"ComputerName,omitempty"`
	FriendlyName             string  `json:"FriendlyName,omitempty"`
	DomainName               string  `json:"DomainName,omitempty"`
	OperatingSystemName      string  `json:"OperatingSystemName,omitempty"`
	OperatingSystemVersion   string  `json:"OperatingSystemVersion,omitempty"`
	Type                     string  `json:"Type,omitempty"`
	Status                   string  `json:"Status,omitempty"`
	Comment                  string  `json:"Comment,omitempty"`
	SerialNumber             string  `json:"SerialNumber,omitempty"`
	AssetTag                 string  `json:"AssetTag,omitempty"`
	LastUserName             string  `json:"LastUserName,omitempty"`
	PrimaryContactName       string  `json:"PrimaryContactName,omitempty"`
	LocalIPAddress           string  `json:"LocalIPAddress,omitempty"`
	GatewayIPAddress         string  `json:"GatewayIPAddress,omitempty"`
	MACAddress               string  `json:"MACAddress,omitempty"`
	RemoteAgentVersion       string  `json:"RemoteAgentVersion,omitempty"`
	BiosManufacturer         string  `json:"BiosManufacturer,omitempty"`
	BiosFlash                string  `json:"BiosFlash,omitempty"`
	UTCOffset                int     `json:"UTCOffset,omitempty"`
	TotalMemory              int64   `json:"TotalMemory,omitempty"`
	FreeMemory               int64   `json:"FreeMemory,omitempty"`
	CpuUsage                 int     `json:"CpuUsage,omitempty"`
	SystemUptime             int64   `json:"SystemUptime,omitempty"`
	UserIdleTime             int     `json:"UserIdleTime,omitempty"`
	Bandwidth                int     `json:"Bandwidth,omitempty"`
	BandwidthDisplay         string  `json:"BandwidthDisplay,omitempty"`
	MasterMode               string  `json:"MasterMode,omitempty"`
	TempFiles                string  `json:"TempFiles,omitempty"`
	CurrentPowerProfile      string  `json:"CurrentPowerProfile,omitempty"`
	CpuScore                 float64 `json:"CpuScore,omitempty"`
	D3DScore                 float64 `json:"D3DScore,omitempty"`
	DiskScore                float64 `json:"DiskScore,omitempty"`
	GraphicsScore            float64 `json:"GraphicsScore,omitempty"`
	MemoryScore              float64 `json:"MemoryScore,omitempty"`
	IsFasTalk                bool    `json:"IsFasTalk,omitempty"`
	IsMaster                 bool    `json:"IsMaster,omitempty"`
	IsNetworkProbe           bool    `json:"IsNetworkProbe,omitempty"`
	IsHeartbeatEnabled       bool    `json:"IsHeartbeatEnabled,omitempty"`
	IsHeartbeatRunning       bool    `json:"IsHeartbeatRunning,omitempty"`
	IsTunnelSupported        bool    `json:"IsTunnelSupported,omitempty"`
	IsMaintenanceModeEnabled bool    `json:"IsMaintenanceModeEnabled,omitempty"`
	IsVirtualMachine         bool    `json:"IsVirtualMachine,omitempty"`
	IsVirtualHost            bool    `json:"IsVirtualHost,omitempty"`
	IsLockedDown             bool    `json:"IsLockedDown,omitempty"`
	IsSystemAccount          bool    `json:"IsSystemAccount,omitempty"`
	IsRebootNeeded           bool    `json:"IsRebootNeeded,omitempty"`
	HasIntelVPRO             bool    `json:"HasIntelVPRO,omitempty"`
	HasIntelAMT              bool    `json:"HasIntelAMT,omitempty"`
	HasHPiLO                 bool    `json:"HasHPiLO,omitempty"`
	RemoteAgentLastContact   Time    `json:"RemoteAgentLastContact,omitzero"`
	RemoteAgentLastInventory Time    `json:"RemoteAgentLastInventory,omitzero"`
	LastInventoryReceived    Time    `json:"LastInventoryReceived,omitzero"`
	WindowsUpdateDate        Time    `json:"WindowsUpdateDate,omitzero"`
	AntivirusDefinitionDate  Time    `json:"AntivirusDefinitionDate,omitzero"`
	LastHeartbeat            Time    `json:"LastHeartbeat,omitzero"`
	LastStartup              Time    `json:"LastStartup,omitzero"`
	DateAdded                Time    `json:"DateAdded,omitzero"`
	AssetDate                Time    `json:"AssetDate,omitzero"`
	WarrantyEndDate          Time    `json:"WarrantyEndDate,omitzero"`

	OpenPortsTCP      []int    `json:"OpenPortsTCP,omitempty"`
	OpenPortsUDP      []int    `json:"OpenPortsUDP,omitempty"`
	DomainNameServers []string `json:"DomainNameServers,omitempty"`
	PowerProfiles     []string `json:"PowerProfiles,omitempty"`
	HardwarePorts     []string `json:"HardwarePorts,omitempty"`
	UserAccounts      []string `json:"UserAccounts,omitempty"`
	IRQ               []int    `json:"IRQ,omitempty"`
	Address           []int    `json:"Address,omitempty"`
	DMA               []int    `json:"DMA,omitempty"`

	Location        json.RawMessage `json:"Location,omitempty"`
	Client          json.RawMessage `json:"Client,omitempty"`
	Contact         json.RawMessage `json:"Contact,omitempty"`
	VirusScanner    json.RawMessage `json:"VirusScanner,omitempty"`
	CommentPriority json.RawMessage `json:"CommentPriority,omitempty"`
	LoggedInUsers   json.RawMessage `json:"LoggedInUsers,omitempty"`
	Tickets         json.RawMessage `json:"Tickets,omitempty"`
	Groups          json.RawMessage `json:"Groups,omitempty"`
}

// ComputerDevice is a hardware device on a computer (LabTech.Models.ComputerDevice).
type ComputerDevice struct {
	ComputerId    int    `json:"ComputerId,omitempty"`
	PnpDeviceId   string `json:"PnpDeviceId,omitempty"`
	DeviceName    string `json:"DeviceName,omitempty"`
	DeviceType    string `json:"DeviceType,omitempty"`
	DriverVersion string `json:"DriverVersion,omitempty"`
	DriverDate    Time   `json:"DriverDate,omitzero"`
	DriverName    string `json:"DriverName,omitempty"`
	DriverFile    string `json:"DriverFile,omitempty"`
	Manufacturer  string `json:"Manufacturer,omitempty"`
	UpdateDate    Time   `json:"UpdateDate,omitzero"`
}

// ComputerService is a Windows service on a computer (LabTech.Models.ComputerService).
type ComputerService struct {
	ComputerServiceId   int64           `json:"ComputerServiceId,omitempty"`
	ComputerId          int             `json:"ComputerId,omitempty"`
	Name                string          `json:"Name,omitempty"`
	Description         string          `json:"Description,omitempty"`
	State               string          `json:"State,omitempty"`
	Startup             string          `json:"Startup,omitempty"`
	PathName            string          `json:"PathName,omitempty"`
	ServiceType         string          `json:"ServiceType,omitempty"`
	Username            string          `json:"Username,omitempty"`
	DateLastInventoried Time            `json:"DateLastInventoried,omitzero"`
	RunLevels           string          `json:"RunLevels,omitempty"`
	Classification      json.RawMessage `json:"Classification,omitempty"`
}

// ComputerSoftware is an installed application on a computer
// (LabTech.Models.ComputerSoftware).
type ComputerSoftware struct {
	ApplicationId       int64           `json:"ApplicationId,omitempty"`
	ComputerId          int             `json:"ComputerId,omitempty"`
	Name                string          `json:"Name,omitempty"`
	InstallationPath    string          `json:"InstallationPath,omitempty"`
	DateInstalled       Time            `json:"DateInstalled,omitzero"`
	Size                int             `json:"Size,omitempty"`
	UninstallerPath     string          `json:"UninstallerPath,omitempty"`
	Version             string          `json:"Version,omitempty"`
	Classification      json.RawMessage `json:"Classification,omitempty"`
	DateLastInventoried Time            `json:"DateLastInventoried,omitzero"`
	ClientId            int             `json:"ClientId,omitempty"`
	ComputerName        string          `json:"ComputerName,omitempty"`
}

// ComputerOperatingSystem describes a computer's OS
// (LabTech.Models.ComputerOperatingSystem).
type ComputerOperatingSystem struct {
	ComputerId           int    `json:"ComputerId,omitempty"`
	Name                 string `json:"Name,omitempty"`
	MajorVersion         int    `json:"MajorVersion,omitempty"`
	MinorVersion         int    `json:"MinorVersion,omitempty"`
	Version              string `json:"Version,omitempty"`
	DotNetVersion        string `json:"DotNetVersion,omitempty"`
	ServicePack          string `json:"ServicePack,omitempty"`
	ServicePackName      string `json:"ServicePackName,omitempty"`
	ProductType          int    `json:"ProductType,omitempty"`
	Suite                int    `json:"Suite,omitempty"`
	ProductInfo          int    `json:"ProductInfo,omitempty"`
	IsLicensed           bool   `json:"IsLicensed,omitempty"`
	IsTablet             bool   `json:"IsTablet,omitempty"`
	IsStarter            bool   `json:"IsStarter,omitempty"`
	IsMediaCenter        bool   `json:"IsMediaCenter,omitempty"`
	BaseFolder           string `json:"BaseFolder,omitempty"`
	SystemDrive          string `json:"SystemDrive,omitempty"`
	HasGui               bool   `json:"HasGui,omitempty"`
	Is64Bit              bool   `json:"Is64Bit,omitempty"`
	Domain               string `json:"Domain,omitempty"`
	IsDomainController   bool   `json:"IsDomainController,omitempty"`
	IsServer             bool   `json:"IsServer,omitempty"`
	InstallDate          Time   `json:"InstallDate,omitzero"`
	DateUpdated          Time   `json:"DateUpdated,omitzero"`
	ReleaseId            int    `json:"ReleaseId,omitempty"`
	Edition              string `json:"Edition,omitempty"`
	BranchReadinessLevel int    `json:"BranchReadinessLevel,omitempty"`
}

// ComputerBios describes a computer's BIOS (LabTech.Models.ComputerBios).
type ComputerBios struct {
	ComputerId        int     `json:"ComputerId,omitempty"`
	Vendor            string  `json:"Vendor,omitempty"`
	Version           string  `json:"Version,omitempty"`
	Date              Time    `json:"Date,omitzero"`
	Size              int     `json:"Size,omitempty"`
	SystemBiosVersion float64 `json:"SystemBiosVersion,omitempty"`
	EMBCTLVersion     float64 `json:"EMBCTLVersion,omitempty"`
	SmBiosVersion     float64 `json:"SmBiosVersion,omitempty"`
	DmiVersion        int     `json:"DmiVersion,omitempty"`
	BiosChar          int     `json:"BiosChar,omitempty"`
	SupportsAcpi      bool    `json:"SupportsAcpi,omitempty"`
	SupportsApm       bool    `json:"SupportsApm,omitempty"`
	SupportsAgp       bool    `json:"SupportsAgp,omitempty"`
	SupportsPcmcia    bool    `json:"SupportsPcmcia,omitempty"`
	HasSmartBattery   bool    `json:"HasSmartBattery,omitempty"`
	SupportsUefi      bool    `json:"SupportsUefi,omitempty"`
	SupportsLegacyUsb bool    `json:"SupportsLegacyUsb,omitempty"`
	SupportsPci       bool    `json:"SupportsPci,omitempty"`
	SupportsVlvesa    bool    `json:"SupportsVlvesa,omitempty"`
	SupportsEscd      bool    `json:"SupportsEscd,omitempty"`
	SupportsNetBoot   bool    `json:"SupportsNetBoot,omitempty"`
	SupportsI2OBoot   bool    `json:"SupportsI2OBoot,omitempty"`
	IsVirtualMachine  bool    `json:"IsVirtualMachine,omitempty"`
	PowerOnReason     string  `json:"PowerOnReason,omitempty"`
	IsPortable        bool    `json:"IsPortable,omitempty"`
	VmHost            string  `json:"VmHost,omitempty"`
	VmType            string  `json:"VmType,omitempty"`
	VmName            string  `json:"VmName,omitempty"`
	DateUpdated       Time    `json:"DateUpdated,omitzero"`
}

// ComputerDrive is a disk volume on a computer (LabTech.Models.ComputerDrive).
type ComputerDrive struct {
	DriveId             int    `json:"DriveId,omitempty"`
	ComputerId          int    `json:"ComputerId,omitempty"`
	Letter              string `json:"Letter,omitempty"`
	Size                int    `json:"Size,omitempty"`
	FreeSpace           int    `json:"FreeSpace,omitempty"`
	FileSystem          string `json:"FileSystem,omitempty"`
	Model               string `json:"Model,omitempty"`
	SmartStatus         string `json:"SmartStatus,omitempty"`
	IsMissing           bool   `json:"IsMissing,omitempty"`
	DateLastInventoried Time   `json:"DateLastInventoried,omitzero"`
	VolumeName          string `json:"VolumeName,omitempty"`
	BackupFlag          int    `json:"BackupFlag,omitempty"`
	IsSolidState        bool   `json:"IsSolidState,omitempty"`
	IsInternal          bool   `json:"IsInternal,omitempty"`

	MaximumHistoryDaysAvailable int64 `json:"MaximumHistoryDaysAvailable,omitempty"`
}

// ComputerPatchingStats summarizes a computer's patch compliance
// (Automate.Api.Domain.Contracts.Patching.ComputerPatchingStats).
type ComputerPatchingStats struct {
	ComputerId                int     `json:"ComputerId,omitempty"`
	OverallCompliance         float64 `json:"OverallCompliance,omitempty"`
	InstalledPatchCount       int     `json:"InstalledPatchCount,omitempty"`
	MissingPatchCount         int     `json:"MissingPatchCount,omitempty"`
	FailedPatchCount          int     `json:"FailedPatchCount,omitempty"`
	CompliantSoftwareCount    int     `json:"CompliantSoftwareCount,omitempty"`
	NonCompliantSoftwareCount int     `json:"NonCompliantSoftwareCount,omitempty"`
	FailedSoftwareCount       int     `json:"FailedSoftwareCount,omitempty"`
	IncorrectSoftwareCount    int     `json:"IncorrectSoftwareCount,omitempty"`
	Stage                     string  `json:"Stage,omitempty"`
	NoPatchInventory          bool    `json:"NoPatchInventory,omitempty"`
	WSUSEnabled               bool    `json:"WSUSEnabled,omitempty"`
	PatchJobRunning           bool    `json:"PatchJobRunning,omitempty"`
	DaytimePatchingEnabled    bool    `json:"DaytimePatchingEnabled,omitempty"`
	WUAOutOfDate              bool    `json:"WUAOutOfDate,omitempty"`
	MissingBaselinePatches    bool    `json:"MissingBaselinePatches,omitempty"`
	WUAVersion                string  `json:"WUAVersion,omitempty"`
	LastInstallWindow         Time    `json:"LastInstallWindow,omitzero"`
	NextInstallWindow         Time    `json:"NextInstallWindow,omitzero"`
	LastSoftwareWindow        Time    `json:"LastSoftwareWindow,omitzero"`
	NextSoftwareWindow        Time    `json:"NextSoftwareWindow,omitzero"`
	LastPatchedDate           Time    `json:"LastPatchedDate,omitzero"`
	LastMicrosoftPatchedDate  Time    `json:"LastMicrosoftPatchedDate,omitzero"`
	LastThirdPartyPatchedDate Time    `json:"LastThirdPartyPatchedDate,omitzero"`
	LastPatchInventory        Time    `json:"LastPatchInventory,omitzero"`
	IsMicrosoftManaged        bool    `json:"IsMicrosoftManaged,omitempty"`
	IsThirdPartyManaged       bool    `json:"IsThirdPartyManaged,omitempty"`
}

// ProductKey is a product key on a client (LabTech.Models.ProductKey).
type ProductKey struct {
	ID             int    `json:"Id,omitempty"`
	ClientId       int    `json:"ClientId,omitempty"`
	ProductName    string `json:"ProductName,omitempty"`
	SerialNumber   string `json:"SerialNumber,omitempty"`
	LicenseKey     string `json:"LicenseKey,omitempty"`
	DoesExpire     bool   `json:"DoesExpire,omitempty"`
	ExpirationDate Time   `json:"ExpirationDate,omitzero"`
	Notes          string `json:"Notes,omitempty"`
	ComputerId     int    `json:"ComputerId,omitempty"`
	ComputerName   string `json:"ComputerName,omitempty"`
}
