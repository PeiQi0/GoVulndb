// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package semver

import (
	"sort"

	"golang.org/x/mod/semver"
	"golang.org/x/vulndb/internal/osv"
)

func AffectsSemver(ranges []osv.Range, v string) bool {
	if len(ranges) == 0 {
		// No ranges implies all versions are affected
		return true
	}
	var semverRangePresent bool
	for _, r := range ranges {
		if r.Type != osv.RangeTypeSemver {
			continue
		}
		semverRangePresent = true
		if containsSemver(r, v) {
			return true
		}
	}
	// If there were no semver ranges present we
	// assume that all semvers are affected, similarly
	// to how to we assume all semvers are affected
	// if there are no ranges at all.
	return !semverRangePresent
}

// containsSemver checks if semver version v is in the
// range encoded by ar. If ar is not a semver range,
// returns false.
//
// Assumes that
//   - exactly one of Introduced or Fixed fields is set
//   - ranges in ar are not overlapping
//   - beginning of time is encoded with .Introduced="0"
//   - no-fix is not an event, as opposed to being an
//     event where Introduced="" and Fixed=""
func containsSemver(ar osv.Range, v string) bool {
	if ar.Type != osv.RangeTypeSemver {
		return false
	}
	if len(ar.Events) == 0 {
		return true
	}
	// Strip and then add the semver prefix so we can support bare versions,
	// versions prefixed with 'v', and versions prefixed with 'go'.
	v = CanonicalizeSemverPrefix(v)
	// Sort events by semver versions. Event for beginning
	// of time, if present, always comes first.
	sort.SliceStable(ar.Events, func(i, j int) bool {
		e1 := ar.Events[i]
		v1 := e1.Introduced
		if v1 == "0" {
			// -inf case.
			return true
		}
		if e1.Fixed != "" {
			v1 = e1.Fixed
		}
		e2 := ar.Events[j]
		v2 := e2.Introduced
		if v2 == "0" {
			// -inf case.
			return false
		}
		if e2.Fixed != "" {
			v2 = e2.Fixed
		}
		return semver.Compare(CanonicalizeSemverPrefix(v1), CanonicalizeSemverPrefix(v2)) < 0
	})
	var affected bool
	for _, e := range ar.Events {
		if !affected && e.Introduced != "" {
			affected = e.Introduced == "0" || semver.Compare(v, CanonicalizeSemverPrefix(e.Introduced)) >= 0
		} else if affected && e.Fixed != "" {
			affected = semver.Compare(v, CanonicalizeSemverPrefix(e.Fixed)) < 0
		}
	}
	return affected
}
