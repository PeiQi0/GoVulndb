// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package osv implements the Go OSV vulnerability format
// (https://go.dev/security/vuln/database#schema), which is a subset of
// the OSV shared vulnerability format
// (https://ossf.github.io/osv-schema), with database and
// ecosystem-specific meanings and fields.
//
// As this package is intended for use with the Go vulnerability
// database, only the subset of features which are used by that
// database are implemented (for instance, only the SEMVER affected
// range type is implemented).
package osv

// RangeType specifies the type of version range being recorded and
// defines the interpretation of the RangeEvent object's Introduced
// and Fixed fields.
//
// In this implementation, only the "SEMVER" type is supported.
//
// See https://ossf.github.io/osv-schema/#affectedrangestype-field.
type RangeType string

// RangeTypeSemver indicates a semantic version as defined by
// SemVer 2.0.0, with no leading "v" prefix.
const RangeTypeSemver RangeType = "SEMVER"

// Ecosystem identifies the overall library ecosystem.
// In this implementation, only the "Go" ecosystem is supported.
type Ecosystem string

// GoEcosystem indicates the Go ecosystem.
const GoEcosystem Ecosystem = "Go"

// Pseudo-module paths used to describe vulnerabilities
// in the Go standard library and toolchain.
const (
	// GoStdModulePath is the pseudo-module path string used
	// to describe vulnerabilities in the Go standard library.
	GoStdModulePath = "stdlib"
	// GoCmdModulePath is the pseudo-module path string used
	// to describe vulnerabilities in the go command.
	GoCmdModulePath = "toolchain"
)

// Module identifies the Go module containing the vulnerability.
// Note that this field is called "package" in the OSV specification.
//
// See https://ossf.github.io/osv-schema/#affectedpackage-field.
type Module struct {
	// The Go module path. Required.
	// For the Go standard library, this is "stdlib".
	// For the Go toolchain, this is "toolchain."
	Path string `json:"name"`
	// The ecosystem containing the module. Required.
	// This should always be "Go".
	Ecosystem Ecosystem `json:"ecosystem"`
}

// RangeEvent describes a single module version that either
// introduces or fixes a vulnerability.
//
// Exactly one of Introduced and Fixed must be present. Other range
// event types (e.g, "last_affected" and "limit") are not supported in
// this implementation.
//
// See https://ossf.github.io/osv-schema/#affectedrangesevents-fields.
type RangeEvent struct {
	// Introduced is a version that introduces the vulnerability.
	// A special value, "0", represents a version that sorts before
	// any other version, and should be used to indicate that the
	// vulnerability exists from the "beginning of time".
	Introduced string `json:"introduced,omitempty"`
	// Fixed is a version that fixes the vulnerability.
	Fixed string `json:"fixed,omitempty"`
}

// Range describes the affected versions of the vulnerable module.
//
// See https://ossf.github.io/osv-schema/#affectedranges-field.
type Range struct {
	// Type is the version type that should be used to interpret the
	// versions in Events. Required.
	// In this implementation, only the "SEMVER" type is supported.
	Type RangeType `json:"type"`
	// Events is a list of versions representing the ranges in which
	// the module is vulnerable. Required.
	// The events should be sorted, and MUST represent non-overlapping
	// ranges.
	// There must be at least one RangeEvent containing a value for
	// Introduced.
	// See https://ossf.github.io/osv-schema/#examples for examples.
	Events []RangeEvent `json:"events"`
}

// Reference type is a reference (link) type.
type ReferenceType string

const (
	// ReferenceTypeAdvisory is a published security advisory for
	// the vulnerability.
	ReferenceTypeAdvisory = ReferenceType("ADVISORY")
	// ReferenceTypeArticle is an article or blog post describing the vulnerability.
	ReferenceTypeArticle = ReferenceType("ARTICLE")
	// ReferenceTypeReport is a report, typically on a bug or issue tracker, of
	// the vulnerability.
	ReferenceTypeReport = ReferenceType("REPORT")
	// ReferenceTypeFix is a source code browser link to the fix (e.g., a GitHub commit).
	ReferenceTypeFix = ReferenceType("FIX")
	// ReferenceTypePackage is a home web page for the package.
	ReferenceTypePackage = ReferenceType("PACKAGE")
	// ReferenceTypeEvidence is a demonstration of the validity of a vulnerability claim.
	ReferenceTypeEvidence = ReferenceType("EVIDENCE")
	// ReferenceTypeWeb is a web page of some unspecified kind.
	ReferenceTypeWeb = ReferenceType("WEB")
)

// ReferenceTypes is the set of reference types defined in OSV.
var ReferenceTypes = []ReferenceType{
	ReferenceTypeAdvisory,
	ReferenceTypeArticle,
	ReferenceTypeReport,
	ReferenceTypeFix,
	ReferenceTypePackage,
	ReferenceTypeEvidence,
	ReferenceTypeWeb,
}

// Reference is a reference URL containing additional information,
// advisories, issue tracker entries, etc., about the vulnerability.
//
// See https://ossf.github.io/osv-schema/#references-field.
type Reference struct {
	// The type of reference. Required.
	Type ReferenceType `json:"type"`
	// The fully-qualified URL of the reference. Required.
	URL string `json:"url"`
}

// Affected gives details about a module affected by the vulnerability.
//
// See https://ossf.github.io/osv-schema/#affected-fields.
type Affected struct {
	// The affected Go module. Required.
	// Note that this field is called "package" in the OSV specification.
	Module Module `json:"package"`
	// The module version ranges affected by the vulnerability.
	Ranges []Range `json:"ranges,omitempty"`
	// Details on the affected packages and symbols within the module.
	EcosystemSpecific *EcosystemSpecific `json:"ecosystem_specific,omitempty"`
}

// Package contains additional information about an affected package.
// This is an ecosystem-specific field for the Go ecosystem.
type Package struct {
	// Path is the package import path. Required.
	Path string `json:"path"`
	// GOOS is the execution operating system where the symbols appear, if
	// known.
	GOOS []string `json:"goos,omitempty"`
	// GOARCH specifies the execution architecture where the symbols appear, if
	// known.
	GOARCH []string `json:"goarch,omitempty"`
	// Symbols is a list of function and method names affected by
	// this vulnerability. Methods are listed as <recv>.<method>.
	//
	// If included, only programs which use these symbols will be marked as
	// vulnerable by `govulncheck`. If omitted, any program which imports this
	// package will be marked vulnerable.
	Symbols []string `json:"symbols,omitempty"`
}

// EcosystemSpecific contains additional information about the vulnerable
// module for the Go ecosystem.
//
// See https://go.dev/security/vuln/database#schema.
type EcosystemSpecific struct {
	// Packages is the list of affected packages within the module.
	Packages []Package `json:"imports,omitempty"`
}

// Entry represents a vulnerability in the Go OSV format, documented
// in https://go.dev/security/vuln/database#schema.
// It is a subset of the OSV schema (https://ossf.github.io/osv-schema).
// Only fields that are published in the Go Vulnerability Database
// are supported.
type Entry struct {
	// SchemaVersion is the OSV schema version used to encode this
	// vulnerability.
	SchemaVersion string `json:"schema_version,omitempty"`
	// ID is a unique identifier for the vulnerability. Required.
	// The Go vulnerability database issues IDs of the form
	// GO-<YEAR>-<ENTRYID>.
	ID string `json:"id"`
	// Modified is the time the entry was last modified. Required.
	Modified Time `json:"modified,omitempty"`
	// Published is the time the entry should be considered to have
	// been published.
	Published Time `json:"published,omitempty"`
	// Withdrawn is the time the entry should be considered to have
	// been withdrawn. If the field is missing, then the entry has
	// not been withdrawn.
	Withdrawn *Time `json:"withdrawn,omitempty"`
	// Aliases is a list of IDs for the same vulnerability in other
	// databases.
	Aliases []string `json:"aliases,omitempty"`
	// Details contains English textual details about the vulnerability.
	Details string `json:"details"`
	// Affected contains information on the modules and versions
	// affected by the vulnerability.
	Affected []Affected `json:"affected"`
	// References contains links to more information about the
	// vulnerability.
	References []Reference `json:"references,omitempty"`
	// Credits contains credits to entities that helped find or fix the
	// vulnerability.
	Credits []Credit `json:"credits,omitempty"`
	// DatabaseSpecific contains additional information about the
	// vulnerability, specific to the Go vulnerability database.
	DatabaseSpecific *DatabaseSpecific `json:"database_specific,omitempty"`
}

// Credit represents a credit for the discovery, confirmation, patch, or
// other event in the life cycle of a vulnerability.
//
// See https://ossf.github.io/osv-schema/#credits-fields.
type Credit struct {
	// Name is the name, label, or other identifier of the individual or
	// entity being credited. Required.
	Name string `json:"name"`
}

// DatabaseSpecific contains additional information about the
// vulnerability, specific to the Go vulnerability database.
//
// See https://go.dev/security/vuln/database#schema.
type DatabaseSpecific struct {
	// The URL of the Go advisory for this vulnerability, of the form
	// "https://pkg.go.dev/GO-YYYY-XXXX".
	URL string `json:"url,omitempty"`
}
