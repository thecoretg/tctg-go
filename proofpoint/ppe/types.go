package ppe

// Types follow the schemas in reference/proofpoint/essentials-v1.json. Every
// field is omitempty and bools and numbers are pointers, so a struct used as a
// request body sends only the fields that were set. Leaving a field nil never
// sends false or 0 over a stored value.

// EntityType is the kind of entity an Organization is.
type EntityType string

const (
	EntityOEMPartner       EntityType = "oem_partner"
	EntityStrategicPartner EntityType = "strategic_partner"
	EntityChannel          EntityType = "channel"
	EntityOrganization     EntityType = "organization"
)

// UserType is a user's role.
type UserType string

const (
	UserOEMPartnerAdmin       UserType = "oem_partner_admin"
	UserStrategicPartnerAdmin UserType = "strategic_partner_admin"
	UserChannelAdmin          UserType = "channel_admin"
	UserOrganizationAdmin     UserType = "organization_admin"
	UserEndUser               UserType = "end_user"
	UserSilentUser            UserType = "silent_user"
	UserFunctionalAccount     UserType = "functional_account"
)

// LicensingPackage is a Proofpoint Essentials package tier.
type LicensingPackage string

const (
	PackageBeginner         LicensingPackage = "beginner"
	PackageBusiness         LicensingPackage = "business"
	PackageBusinessPlus     LicensingPackage = "business_plus"
	PackageAdvanced         LicensingPackage = "advanced"
	PackageAdvancedPlus     LicensingPackage = "advanced_plus"
	PackageProfessional     LicensingPackage = "professional"
	PackageProfessionalPlus LicensingPackage = "professional_plus"
)

// Me describes the account the client is authenticated as.
type Me struct {
	UserID              *int       `json:"user_id,omitempty"`
	Firstname           string     `json:"firstname,omitempty"`
	Surname             string     `json:"surname,omitempty"`
	PrimaryEmail        string     `json:"primary_email,omitempty"`
	Type                UserType   `json:"type,omitempty"`
	IsAdmin             *bool      `json:"is_admin,omitempty"`
	IsPartnerAdmin      *bool      `json:"is_partner_admin,omitempty"`
	EntityID            *int       `json:"entity_id,omitempty"`
	EntityPrimaryDomain string     `json:"entity_primary_domain,omitempty"`
	EntityType          EntityType `json:"entity_type,omitempty"`
	ReadOnlyUser        *bool      `json:"read_only_user,omitempty"`

	// The fields below are returned by the API but not in the spec.
	EntityName     string   `json:"entity_name,omitempty"`
	EntityTypeID   *int     `json:"entity_type_id,omitempty"`
	UserTypeID     *int     `json:"user_type_id,omitempty"`
	IsDefaultAdmin *bool    `json:"is_default_admin,omitempty"`
	HasBetaAccess  *bool    `json:"has_beta_access,omitempty"`
	Aliases        []string `json:"aliases,omitempty"`
	// Permissions maps a permission such as "users.update" to whether the
	// account holds it.
	Permissions map[string]bool `json:"permissions,omitempty"`
	Config      map[string]any  `json:"config,omitempty"`
}

// Organization is any entity in the hierarchy: a partner or a customer
// organization. It is identified in URLs by its primary domain. Creating one
// requires PrimaryDomain, Name, Type, AdminUser, Domains and UserLicenses.
type Organization struct {
	PrimaryDomain          string           `json:"primary_domain,omitempty"`
	Name                   string           `json:"name,omitempty"`
	IsActive               *bool            `json:"is_active,omitempty"`
	ConfigurationType      string           `json:"configuration_type,omitempty"` // automatic, manual, skipped
	DeploymentMethod       string           `json:"deployment_method,omitempty"`  // integrated, mx_record
	Type                   EntityType       `json:"type,omitempty"`
	EID                    *int             `json:"eid,omitempty"`
	OrganizationHierarchy  []string         `json:"organization_hierarchy,omitempty"`
	LicensingPackage       LicensingPackage `json:"licensing_package,omitempty"`
	IsOnTrial              *bool            `json:"is_on_trial,omitempty"`
	UserLicenses           *int             `json:"user_licenses,omitempty"`
	ActiveUsers            *int             `json:"active_users,omitempty"`
	WhenRenewal            string           `json:"when_renewal,omitempty"`
	IsBeginnerPlus         *bool            `json:"is_beginner_plus,omitempty"`
	IsBeginnerPlusEnabled  *bool            `json:"is_beginner_plus_enabled,omitempty"`
	WWW                    string           `json:"www,omitempty"`
	Address                string           `json:"address,omitempty"`
	Postcode               string           `json:"postcode,omitempty"`
	StateProv              string           `json:"stateprov,omitempty"`
	Country                string           `json:"country,omitempty"`
	Phone                  string           `json:"phone,omitempty"`
	ExternalID             string           `json:"external_id,omitempty"`
	AdminUser              *User            `json:"admin_user,omitempty"`
	Domains                []Domain         `json:"domains,omitempty"`
	OutgoingServers        []string         `json:"outgoing_servers,omitempty"`
	SafeListSenders        []string         `json:"safe_list_senders,omitempty"`
	BlockListSenders       []string         `json:"block_list_senders,omitempty"`
	WhiteListSenders       []string         `json:"white_list_senders,omitempty"` // deprecated term for SafeListSenders
	BlackListSenders       []string         `json:"black_list_senders,omitempty"` // deprecated term for BlockListSenders
	LDAPURL                string           `json:"ldap_url,omitempty"`
	LDAPUsername           string           `json:"ldap_username,omitempty"`
	LDAPBaseDN             string           `json:"ldap_basedn,omitempty"`
	IsSMTPDiscoveryEnabled *bool            `json:"is_smtp_discovery_enabled,omitempty"`
	AccountTemplateID      string           `json:"account_template_id,omitempty"`
	SendWelcomeEmail       *bool            `json:"send_welcome_email,omitempty"`
	OdinCapabilities       string           `json:"odin_capabilities,omitempty"` // odin_master, odin_child, none
	BillingMethod          string           `json:"billing_method,omitempty"`
	UUID                   string           `json:"uuid,omitempty"` // returned by the API but not in the spec
}

// OrgUpdate is the body of UpdateOrg. is_active is the only field the API
// documents as updatable.
type OrgUpdate struct {
	IsActive *bool `json:"is_active,omitempty"`
}

// User belongs to one Organization and is identified in URLs by its primary
// email. Creating one requires PrimaryEmail.
type User struct {
	UID              *int              `json:"uid,omitempty"`
	PrimaryEmail     string            `json:"primary_email,omitempty"`
	IsActive         *bool             `json:"is_active,omitempty"`
	Type             UserType          `json:"type,omitempty"`
	SendWelcomeEmail *bool             `json:"send_welcome_email,omitempty"`
	Firstname        string            `json:"firstname,omitempty"`
	Surname          string            `json:"surname,omitempty"`
	AliasEmails      []string          `json:"alias_emails,omitempty"`
	SafeListSenders  []string          `json:"safe_list_senders,omitempty"`
	BlockListSenders []string          `json:"block_list_senders,omitempty"`
	WhiteListSenders []string          `json:"white_list_senders,omitempty"` // deprecated term for SafeListSenders
	BlackListSenders []string          `json:"black_list_senders,omitempty"` // deprecated term for BlockListSenders
	IsBillable       *bool             `json:"is_billable,omitempty"`
	OdinSettings     *UserOdinSettings `json:"odin_settings,omitempty"`
	ReadOnlyUser     *bool             `json:"read_only_user,omitempty"`
	CreationDate     string            `json:"creation_date,omitempty"`
	LastUpdate       string            `json:"last_update,omitempty"`
	LastLogin        string            `json:"last_login,omitempty"`
	IsDefaultAdmin   *bool             `json:"is_default_admin,omitempty"` // returned by the API but not in the spec
	// Password is write-only.
	Password string `json:"password,omitempty"`
}

type UserOdinSettings struct {
	CanImpersonate *bool `json:"can_impersonate,omitempty"`
	IsImitable     *bool `json:"is_imitable,omitempty"`
}

// Domain is a mail domain owned by an Organization. Creating one requires
// Name.
type Domain struct {
	Name        string   `json:"name,omitempty"`
	IsActive    *bool    `json:"is_active,omitempty"`
	IsRelay     *bool    `json:"is_relay,omitempty"`
	Destination string   `json:"destination,omitempty"`
	Failovers   []string `json:"failovers,omitempty"`
}

// Licensing is an Organization's license count and package.
type Licensing struct {
	LicenseCount   *int   `json:"license_count,omitempty"`
	Package        string `json:"package,omitempty"`
	IsOnTrial      *bool  `json:"is_on_trial,omitempty"`
	IsBeginnerPlus *bool  `json:"is_beginner_plus,omitempty"`
}

// Package is the body of UpdatePackage. Package and IsEULAConfirmed are
// required.
type Package struct {
	Package         LicensingPackage `json:"package,omitempty"`
	IsActivated     *bool            `json:"is_activated,omitempty"`
	IsEULAConfirmed *bool            `json:"is_eula_confirmed,omitempty"`
	IsTrialExtended *bool            `json:"is_trial_extended,omitempty"`
}

// Features maps a feature name, such as "url_defense" or "instant_replay", to
// its value. Values are bools or strings depending on the feature, and the set
// of features depends on the Organization's package.
type Features map[string]any

// Product is an add-on product, available or purchased.
type Product struct {
	// Status is one of available, provisioning, requires_user_action, error,
	// confirmed, cancelled, deleting.
	Status        string               `json:"status,omitempty"`
	Label         string               `json:"label,omitempty"`
	Name          string               `json:"name,omitempty"`
	Description   string               `json:"description,omitempty"`
	BillingMethod string               `json:"billing_method,omitempty"`
	MinTerm       *int                 `json:"min_term,omitempty"`
	IsRenewable   *bool                `json:"is_renewable,omitempty"`
	MinLicenses   *int                 `json:"min_licenses,omitempty"`
	HelpURL       string               `json:"help_url,omitempty"`
	Variants      []ProductVariant     `json:"variants,omitempty"`
	PurchaseInfo  *ProductPurchaseInfo `json:"purchase_info,omitempty"`
}

type ProductVariant struct {
	Name          string `json:"name,omitempty"`
	Ref           string `json:"ref,omitempty"`
	OrderIndex    *int   `json:"order_index,omitempty"`
	Description   string `json:"description,omitempty"`
	IsPurchasable *bool  `json:"is_purchasable,omitempty"`
}

type ProductPurchaseInfo struct {
	WhenAdded         string `json:"when_added,omitempty"`
	Variant           string `json:"variant,omitempty"`
	WhenRenewal       string `json:"when_renewal,omitempty"`
	WhenStatusUpdated string `json:"when_status_updated,omitempty"`
	AutoRenew         *bool  `json:"auto_renew,omitempty"`
	TrialVariant      string `json:"trial_variant,omitempty"`
	IsTrial           *bool  `json:"is_trial,omitempty"`
	TrialEndDate      string `json:"trial_end_date,omitempty"`
}

// ProductPurchase is the body of PurchaseProduct.
type ProductPurchase struct {
	Label     string `json:"label,omitempty"`
	Variant   string `json:"variant,omitempty"`
	AutoRenew *bool  `json:"auto_renew,omitempty"`
	IsTrial   *bool  `json:"is_trial,omitempty"`
}

// ProductUpdate is the body of UpdateProduct.
type ProductUpdate struct {
	Variant         string `json:"variant,omitempty"`
	AutoRenew       *bool  `json:"auto_renew,omitempty"`
	IsTrial         *bool  `json:"is_trial,omitempty"`
	IsTrialExtended *bool  `json:"is_trial_extended,omitempty"`
}

// BillingReport is one Organization's billing data.
type BillingReport struct {
	EID                *int             `json:"eid,omitempty"`
	Name               string           `json:"name,omitempty"`
	IsActive           *bool            `json:"is_active,omitempty"`
	WhenCreated        string           `json:"when_created,omitempty"`
	CustomerType       EntityType       `json:"customer_type,omitempty"`
	URLDefenseFeature  *bool            `json:"url_defense_feature,omitempty"`
	Address            string           `json:"address,omitempty"`
	StateProv          string           `json:"state_prov,omitempty"`
	Postcode           string           `json:"postcode,omitempty"`
	LicensingPackage   LicensingPackage `json:"licensing_package,omitempty"`
	IsBeginnerPlus     *bool            `json:"is_beginner_plus,omitempty"`
	IsBasicPackage     *bool            `json:"is_basic_package,omitempty"`
	CountryCode        string           `json:"country_code,omitempty"`
	ConfirmedUsers     *int             `json:"confirmed_users,omitempty"`
	WhenRenewal        string           `json:"when_renewal,omitempty"`
	IsTestAccount      *bool            `json:"is_test_account,omitempty"`
	WhenConfirmed      string           `json:"when_confirmed,omitempty"`
	WhenUpdated        string           `json:"when_updated,omitempty"`
	WhenPackageChanged string           `json:"when_package_changed,omitempty"`
	WhenDeactivated    string           `json:"when_deactivated,omitempty"`
	ActiveUsers        *int             `json:"active_users,omitempty"`
	NonBillableUsers   *int             `json:"non_billable_users,omitempty"`
	IsOnTrial          *bool            `json:"is_on_trial,omitempty"`
	TrialPackage       string           `json:"trial_package,omitempty"`
	TrialEndDate       string           `json:"trial_end_date,omitempty"`
	PrimaryDomain      string           `json:"primary_domain,omitempty"`
	MonthlyBilling     string           `json:"monthly_billing,omitempty"`
	DataCenter         string           `json:"data_center,omitempty"`
	BillablePackage    LicensingPackage `json:"billable_package,omitempty"`
	InTotal            *int             `json:"in_total,omitempty"`
	BillingMethod      string           `json:"billing_method,omitempty"`
	MigrationStatus    *int             `json:"migration_status,omitempty"`
	WhenMigrated       string           `json:"when_migrated,omitempty"`
	ExternalID         string           `json:"external_id,omitempty"`

	// The fields below are returned by the API but not in the spec.

	// Bundle is the package's display name, such as "Business +".
	Bundle string `json:"bundle,omitempty"`
	// Products maps a product to its tier's display name, such as
	// "email_security": "Advanced".
	Products map[string]string `json:"products,omitempty"`
	// The ancestors of the Organization, nearest first. Levels above the top
	// of the hierarchy are nil.
	ParentEID                 *int   `json:"parent_eid,omitempty"`
	ParentName                string `json:"parent_name,omitempty"`
	GrandparentEID            *int   `json:"parents_parent_eid,omitempty"`
	GrandparentName           string `json:"parents_parent_name,omitempty"`
	GreatGrandparentEID       *int   `json:"parents_parents_parent_eid,omitempty"`
	GreatGrandparentName      string `json:"parents_parents_parent_name,omitempty"`
	GreatGreatGrandparentEID  *int   `json:"parents_parents_parents_parent_eid,omitempty"`
	GreatGreatGrandparentName string `json:"parents_parents_parents_parent_name,omitempty"`
}
