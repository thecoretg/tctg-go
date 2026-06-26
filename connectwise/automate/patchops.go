package automate

// PatchOp is a single RFC 6902 JSON Patch operation, used by the PATCH methods
// (PatchCompany, PatchLocation, PatchLocationExtraField, ...). The OpenAPI spec
// names these fields Op/Path/Value, but the live API (and ConnectWise's own EDF
// documentation) requires the lowercase RFC 6902 keys op/path/value.
type PatchOp struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value any    `json:"value,omitempty"`
}
