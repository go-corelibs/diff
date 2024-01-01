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
	_ "embed"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDigitValues(t *testing.T) {

	Convey("no differences", t, func() {
		delta := New("/nope", "nope", "nope")
		So(delta, ShouldNotEqual, nil)
		So(len(delta.groups), ShouldEqual, 0)
		unified, err := delta.Unified()
		So(err, ShouldEqual, nil)
		So(unified, ShouldEqual, "")
	})

	Convey("one group", t, func() {
		delta := New("one-group.txt", oneGroupA, oneGroupB)
		So(delta, ShouldNotEqual, nil)
		Convey("correct group len", func() {
			So(delta.EditGroupsLen(), ShouldEqual, 1)
		})
		Convey("correct len", func() {
			So(delta.Len(), ShouldEqual, 3)
		})
		Convey("correct patch", func() {
			unified, err := delta.Unified()
			So(err, ShouldEqual, nil)
			So(unified, ShouldEqual, oneGroupUnified)
		})
	})

	Convey("two groups", t, func() {
		delta := New("two-groups.txt", twoGroupsA, twoGroupsB)
		So(delta, ShouldNotEqual, nil)
		Convey("correct group len", func() {
			So(delta.EditGroupsLen(), ShouldEqual, 2)
		})
		Convey("correct len", func() {
			So(delta.Len(), ShouldEqual, 11)
		})
		Convey("correct patch", func() {
			unified, err := delta.Unified()
			So(err, ShouldEqual, nil)
			So(unified, ShouldEqual, twoGroupsUnified)
		})
	})

	Convey("keep/skip groups", t, func() {
		delta := New("two-groups.txt", twoGroupsA, twoGroupsB)
		So(delta, ShouldNotEqual, nil)
		Convey("correct keep len", func() {
			So(delta.KeepLen(), ShouldEqual, 0)
			delta.KeepAll()
			So(delta.KeepLen(), ShouldEqual, 11)
			delta.SkipAll()
			So(delta.KeepLen(), ShouldEqual, 0)
		})
		Convey("correct keep patches", func() {
			unified := delta.UnifiedEdits()
			So(unified, ShouldEqual, "")
			delta.KeepGroup(0)
			unified = delta.UnifiedEdits()
			So(unified, ShouldEqual, twoGroupsUnifiedEdits0)
			delta.KeepGroup(0)
			unified = delta.UnifiedEdits()
			So(unified, ShouldEqual, twoGroupsUnifiedEdits0)
			delta.SkipGroup(0)
			unified = delta.UnifiedEdits()
			So(unified, ShouldEqual, "")
			unified = delta.EditGroup(0)
			So(unified, ShouldEqual, twoGroupsUnifiedEditGroup0)
			delta.SkipAll()
			modified, err := delta.ModifiedEdits()
			So(err, ShouldEqual, nil)
			So(modified, ShouldEqual, twoGroupsA)
			delta.KeepAll()
			modified, err = delta.ModifiedEdits()
			So(err, ShouldEqual, nil)
			So(modified, ShouldEqual, twoGroupsB)
		})
		Convey("correct unified edit", func() {
			unified := delta.UnifiedEdit(0)
			So(unified, ShouldEqual, twoGroupsUnifiedEdit0)
		})
	})

}
