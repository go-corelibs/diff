// Copyright (c) 2024  The Go-Curses Authors
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
	_ "embed"
)

//go:embed _testing/odd-group_a.txt
var oddGroupA string

//go:embed _testing/odd-group_b.txt
var oddGroupB string

//go:embed _testing/odd-group_unified.patch
var oddGroupUnified string

//go:embed _testing/odd-group_unified.tango-patch
var oddGroupUnifiedTango string

//go:embed _testing/one-group_a.txt
var oneGroupA string

//go:embed _testing/one-group_b.txt
var oneGroupB string

//go:embed _testing/one-group_unified.patch
var oneGroupUnified string

//go:embed _testing/one-group_unified.tango-patch
var oneGroupUnifiedTango string

//go:embed _testing/two-groups_a.txt
var twoGroupsA string

//go:embed _testing/two-groups_b.txt
var twoGroupsB string

//go:embed _testing/two-groups_unified.patch
var twoGroupsUnified string

//go:embed _testing/two-groups_unified_edit_0.patch
var twoGroupsUnifiedEdit0 string

//go:embed _testing/two-groups_unified_edits_0.patch
var twoGroupsUnifiedEdits0 string

//go:embed _testing/two-groups_unified_edit_group_0.patch
var twoGroupsUnifiedEditGroup0 string
