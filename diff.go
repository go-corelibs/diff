// Copyright (c) 2023  The Go-Curses Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package diff

import (
	"fmt"
	"os"
	"sort"

	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"
)

type Diff struct {
	path    string
	source  string
	changed string
	edits   []gotextdiff.TextEdit
	keep    []int
	groups  [][]int
}

// New constructs a new Diff instance with the given source and changed
// strings computed into a set of "edits" which can be selectively
// included in the Diff.UnifiedEdits and Diff.ModifiedEdits outputs
func New(path, source, changed string) (delta *Diff) {
	delta = new(Diff)
	delta.path = path
	delta.source = source
	delta.changed = changed
	delta.keep = nil
	delta.init()
	return
}

func (d *Diff) init() {
	d.edits = myers.ComputeEdits(span.URIFromPath(d.path), d.source, d.changed)
	d.groups = make([][]int, 0)
	previousLine := -1
	var group []int
	for idx, edit := range d.edits {
		end := edit.Span.End()
		thisLine := end.Line()
		if previousLine == -1 {
			previousLine = thisLine
			group = append(group, idx)
			continue
		}
		if thisLine == previousLine || thisLine == previousLine+1 {
			previousLine = thisLine
			group = append(group, idx)
			continue
		}
		d.groups = append(d.groups, group)
		group = append([]int{}, idx)
		previousLine = thisLine
	}
	if len(group) > 0 {
		d.groups = append(d.groups, group)
	}
}

func (d *Diff) abPaths() (a, b string) {
	a = fmt.Sprintf("a%c%v", os.PathSeparator, d.path)
	b = fmt.Sprintf("b%c%v", os.PathSeparator, d.path)
	return
}

// Unified returns the source content modified by all edits
func (d *Diff) Unified() (unified string, err error) {
	ap, bp := d.abPaths()
	unified = fmt.Sprint(gotextdiff.ToUnified(ap, bp, d.source, d.edits))
	return
}

// Len returns the total number of edits (regardless of keep/skip state)
func (d *Diff) Len() (length int) {
	length = len(d.edits)
	return
}

// KeepLen returns the total number of edits flagged to be included
// in the UnifiedEdits and ModifiedEdits output
func (d *Diff) KeepLen() (count int) {
	count = len(d.keep)
	return
}

// KeepAll flags all edits to be included in the UnifiedEdits and
// ModifiedEdits output
func (d *Diff) KeepAll() {
	d.keep = nil
	for idx, _ := range d.edits {
		d.keep = append(d.keep, idx)
	}
}

// KeepEdit flags a particular edit to be included in the UnifiedEdits() and
// ModifiedEdits() output
func (d *Diff) KeepEdit(index int) (ok bool) {
	numEdits := len(d.edits)
	if numEdits > 0 && index >= 0 && index < numEdits {
		var found bool
		for _, kid := range d.keep {
			if kid == index {
				found = true
				break
			}
		}
		if !found {
			d.keep = append(d.keep, index)
			sort.Ints(d.keep)
		}
		ok = true
	}
	return
}

// SkipAll flags all edits to be excluded in the UnifiedEdits() and
// ModifiedEdits() output
func (d *Diff) SkipAll() {
	d.keep = nil
}

// SkipEdit flags a particular edit to be excluded in the UnifiedEdits() output
func (d *Diff) SkipEdit(index int) (ok bool) {
	numEdits := len(d.edits)
	if numEdits > 0 && index >= 0 && index < numEdits {
		idx := -1
		for i, v := range d.keep {
			if index == v {
				idx = i
				break
			}
		}
		if idx > -1 {
			d.keep = append(d.keep[:idx], d.keep[idx+1:]...)
		}
		ok = true
	}
	return
}

// UnifiedEdit returns the unified diff output for just the given edit
func (d *Diff) UnifiedEdit(index int) (unified string) {
	ap, bp := d.abPaths()
	unified = fmt.Sprint(gotextdiff.ToUnified(ap, bp, d.source, []gotextdiff.TextEdit{d.edits[index]}))
	return
}

// UnifiedEdits returns the unified diff output for all kept edits
func (d *Diff) UnifiedEdits() (unified string) {
	ap, bp := d.abPaths()
	var edits []gotextdiff.TextEdit
	for _, index := range d.keep {
		edits = append(edits, d.edits[index])
	}
	unified = fmt.Sprint(gotextdiff.ToUnified(ap, bp, d.source, edits))
	return
}

// ModifiedEdits returns the source content modified by only kept edits
func (d *Diff) ModifiedEdits() (modified string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("gotextdiff error: %v", r)
		}
	}()
	var edits []gotextdiff.TextEdit
	for _, index := range d.keep {
		edits = append(edits, d.edits[index])
	}
	modified = gotextdiff.ApplyEdits(d.source, edits)
	return
}

// EditGroupsLen returns the count of edit groups present
func (d *Diff) EditGroupsLen() (count int) {
	count = len(d.groups)
	return
}

// EditGroup returns the unified diff of the edit group at the given index
func (d *Diff) EditGroup(index int) (unified string) {
	ap, bp := d.abPaths()
	if index >= 0 && index < len(d.groups) {
		var edits []gotextdiff.TextEdit
		for _, gid := range d.groups[index] {
			edits = append(edits, d.edits[gid])
		}
		unified = fmt.Sprint(gotextdiff.ToUnified(ap, bp, d.source, edits))
	}
	return
}

// KeepGroup flags the given group index for including in the UnifiedEdits and
// ModifiedEdits outputs
func (d *Diff) KeepGroup(index int) {
	if index >= 0 && index < len(d.groups) {
		for _, gid := range d.groups[index] {
			d.KeepEdit(gid)
		}
	}
}

// SkipGroup flags the given group index for exclusion from the UnifiedEdits
// and ModifiedEdits outputs
func (d *Diff) SkipGroup(index int) {
	if index >= 0 && index < len(d.groups) {
		for _, gid := range d.groups[index] {
			d.SkipEdit(gid)
		}
	}
}
