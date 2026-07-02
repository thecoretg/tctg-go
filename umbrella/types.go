package umbrella

// ErrorResponse is the standard error body returned by the Umbrella API.
type ErrorResponse struct {
	Message string `json:"message,omitempty"`
}

// Umbrella package IDs, keyed by the package name from the spec. The write and
// trial-conversion endpoints accept only a subset of these values.
const (
	PackageUmbrellaProfessional  = 99
	PackageUmbrellaPlatform      = 101
	PackageUmbrellaInsights      = 107
	PackageUmbrellaWirelessLAN   = 171
	PackageUmbrellaEDU           = 202
	PackageDNSSecurityEssentials = 246
	PackageDNSSecurityAdvantage  = 248
	PackageSIGEssentials         = 250
	PackageSIGAdvantage          = 252
	PackageNFRMSPDNSAdvantage    = 312
)
