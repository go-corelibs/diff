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
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRender(t *testing.T) {

	Convey("RenderBuilder", t, func() {
		rb := NewRenderer()
		ptr, _ := rb.(*CRender)

		// empty
		So(*ptr, ShouldEqual, CRender{})

		// line added
		So(rb.SetLineAdded("[+", "+]"), ShouldEqual, rb)
		So(*ptr, ShouldEqual, CRender{Line: AddRemTags{
			Add: MarkupTag{Open: "[+", Close: "+]"},
		}})

		// line removed
		So(rb.SetLineRemoved("[-", "-]"), ShouldEqual, rb)
		So(*ptr, ShouldEqual, CRender{Line: AddRemTags{
			Add: MarkupTag{Open: "[+", Close: "+]"},
			Rem: MarkupTag{Open: "[-", Close: "-]"},
		}})

		// text added
		So(rb.SetTextAdded("(+", "+)"), ShouldEqual, rb)
		So(*ptr, ShouldEqual, CRender{Line: AddRemTags{
			Add: MarkupTag{Open: "[+", Close: "+]"},
			Rem: MarkupTag{Open: "[-", Close: "-]"},
		}, Text: AddRemTags{
			Add: MarkupTag{Open: "(+", Close: "+)"},
		}})

		// text removed
		So(rb.SetTextRemoved("(-", "-)"), ShouldEqual, rb)
		So(*ptr, ShouldEqual, CRender{Line: AddRemTags{
			Add: MarkupTag{Open: "[+", Close: "+]"},
			Rem: MarkupTag{Open: "[-", Close: "-]"},
		}, Text: AddRemTags{
			Add: MarkupTag{Open: "(+", Close: "+)"},
			Rem: MarkupTag{Open: "(-", Close: "-)"},
		}})

		// comment line
		So(rb.SetComment("#:", ":#"), ShouldEqual, rb)
		So(*ptr, ShouldEqual, CRender{Line: AddRemTags{
			Add: MarkupTag{Open: "[+", Close: "+]"},
			Rem: MarkupTag{Open: "[-", Close: "-]"},
		}, Text: AddRemTags{
			Add: MarkupTag{Open: "(+", Close: "+)"},
			Rem: MarkupTag{Open: "(-", Close: "-)"},
		}, Comment: MarkupTag{
			Open: "#:", Close: ":#",
		}})

		// everything else line
		So(rb.SetNormal("==", "=="), ShouldEqual, rb)
		So(*ptr, ShouldEqual, CRender{Line: AddRemTags{
			Add: MarkupTag{Open: "[+", Close: "+]"},
			Rem: MarkupTag{Open: "[-", Close: "-]"},
		}, Text: AddRemTags{
			Add: MarkupTag{Open: "(+", Close: "+)"},
			Rem: MarkupTag{Open: "(-", Close: "-)"},
		}, Comment: MarkupTag{
			Open: "#:", Close: ":#",
		}, Normal: MarkupTag{
			Open: "==", Close: "==",
		}})

		// file opener, closer lines
		So(rb.SetFile("[begin]", "[end]"), ShouldEqual, rb)
		So(*ptr, ShouldEqual, CRender{Line: AddRemTags{
			Add: MarkupTag{Open: "[+", Close: "+]"},
			Rem: MarkupTag{Open: "[-", Close: "-]"},
		}, Text: AddRemTags{
			Add: MarkupTag{Open: "(+", Close: "+)"},
			Rem: MarkupTag{Open: "(-", Close: "-)"},
		}, Comment: MarkupTag{
			Open: "#:", Close: ":#",
		}, Normal: MarkupTag{
			Open: "==", Close: "==",
		}, File: MarkupTag{Open: "[begin]", Close: "[end]"}})

		// Make checks
		r := rb.Make() // r is just a retyped rb, both are *CRender
		So(r, ShouldEqual, rb)
		crb := r.Clone()
		So(crb, ShouldEqual, rb)
	})

	Convey("RenderLine", t, func() {
		rb := NewRenderer()
		r := rb.Make()
		a := `This is the first line`
		b := `This is the second line`
		ma, mb := r.RenderLine(a, b)
		So(ma, ShouldEqual, a)
		So(mb, ShouldEqual, b)
		rb.SetTextAdded("+", "+")
		rb.SetTextRemoved("-", "-")
		ea := `This is the -fir-s-t- line`
		eb := `This is the s+econd+ line`
		ma, mb = r.RenderLine(a, b)
		So(ma, ShouldEqual, ea)
		So(mb, ShouldEqual, eb)
	})

	Convey("RenderDiff", t, func() {
		rb := TangoRender.Clone()
		r := rb.Make()
		markup := r.RenderDiff(oneGroupUnified)
		So(markup, ShouldEqual, oneGroupUnifiedTango)
		markup = r.RenderDiff(oddGroupUnified)
		So(markup, ShouldEqual, oddGroupUnifiedTango)
	})

}
