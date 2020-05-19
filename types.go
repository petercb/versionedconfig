package versionedconfig

// VersionedConfig is the minimal interface of a versioned config schema
type VersionedConfig interface {
	GetKind() string
	GetVersion() string
}

// Version declares the Kind and SchemaVersion of a VersionedConfig
// as well as the Factory function to create a new one
type Version struct {
	SchemaVersion string
	Kind          string
	Factory       func() VersionedConfig
}

// Versions is a collection of Version structs
type Versions []Version

// Find searches the constructor for a given config version
// It returns the Factory function on success, and a boolean value indicating whether the schemaVersion is present
func (obj *Versions) Find(kind, schemaVersion string) (func() VersionedConfig, bool) {
	for _, ver := range *obj {
		if ver.Kind == kind && ver.SchemaVersion == schemaVersion {
			return ver.Factory, true
		}
	}

	return nil, false
}
