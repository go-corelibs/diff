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

var (
	// TangoRender is a Renderer preset for the go-curses Tango markup
	// format used in the ctk.Label and other widgets
	TangoRender Renderer = &CRender{
		Line: AddRemTags{
			Add: MarkupTag{
				Open:  `<span foreground="#ffffff" background="#007700">`,
				Close: `</span>`,
			},
			Rem: MarkupTag{
				Open:  `<span foreground="#eeeeee" background="#770000">`,
				Close: `</span>`,
			},
		},
		Text: AddRemTags{
			Add: MarkupTag{
				Open:  `<span background="#004400" weight="bold">`,
				Close: `</span>`,
			},
			Rem: MarkupTag{
				Open:  `<span background="#440000" weight="dim" strikethrough="true">`,
				Close: `</span>`,
			},
		},
		Normal: MarkupTag{
			Open:  `<span weight="dim">`,
			Close: `</span>`,
		},
		Comment: MarkupTag{
			Open:  `<span style="italic" weight="dim">`,
			Close: `</span>`,
		},
	}

	// HTMLRender is a Renderer preset for browser presentation
	HTMLRender Renderer = &CRender{
		File: MarkupTag{
			Open:  `<ul style="list-style-type:none;margin:0;padding:0;">` + "\n",
			Close: `</ul>`,
		},
		Normal: MarkupTag{
			Open:  `<li style="opacity:0.77;">`,
			Close: `</li>`,
		},
		Comment: MarkupTag{
			Open:  `<li style="font-style:italic;opacity:0.77;">`,
			Close: `</li>`,
		},
		Line: AddRemTags{
			Add: MarkupTag{
				Open:  `<li style="color:#ffffff;background-color:#007700;">`,
				Close: `</li>`,
			},
			Rem: MarkupTag{
				Open:  `<li style="color:#eeeeee;background-color:#770000;">`,
				Close: `</li>`,
			},
		},
		Text: AddRemTags{
			Add: MarkupTag{
				Open:  `<span style="background-color:#004400;font-weight:bold;">`,
				Close: `</span>`,
			},
			Rem: MarkupTag{
				Open:  `<span style="background-color:#440000;opacity:0.77;text-decoration:line-through;">`,
				Close: `</span>`,
			},
		},
	}
)
