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
	"html"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

type RenderBuilder interface {
	SetFile(open, close string) RenderBuilder
	SetNormal(open, close string) RenderBuilder
	SetComment(open, close string) RenderBuilder
	SetLineAdded(open, close string) RenderBuilder
	SetTextAdded(open, close string) RenderBuilder
	SetLineRemoved(open, close string) RenderBuilder
	SetTextRemoved(open, close string) RenderBuilder
	Make() Renderer
}

type Renderer interface {
	RenderLine(a, b string) (ma, mb string)
	RenderDiff(unified string) (markup string)
	Clone() RenderBuilder
}

type CRender struct {
	File    MarkupTag
	Normal  MarkupTag
	Comment MarkupTag
	Line    AddRemTags
	Text    AddRemTags
}

func NewRenderer() (tb RenderBuilder) {
	tb = &CRender{}
	return
}

func (r *CRender) SetFile(open, close string) RenderBuilder {
	r.File.Open = open
	r.File.Close = close
	return r
}

func (r *CRender) SetNormal(open, close string) RenderBuilder {
	r.Normal.Open = open
	r.Normal.Close = close
	return r
}

func (r *CRender) SetComment(open, close string) RenderBuilder {
	r.Comment.Open = open
	r.Comment.Close = close
	return r
}

func (r *CRender) SetLineAdded(open, close string) RenderBuilder {
	r.Line.Add.Open = open
	r.Line.Add.Close = close
	return r
}

func (r *CRender) SetTextAdded(open, close string) RenderBuilder {
	r.Text.Add.Open = open
	r.Text.Add.Close = close
	return r
}

func (r *CRender) SetLineRemoved(open, close string) RenderBuilder {
	r.Line.Rem.Open = open
	r.Line.Rem.Close = close
	return r
}

func (r *CRender) SetTextRemoved(open, close string) RenderBuilder {
	r.Text.Rem.Open = open
	r.Text.Rem.Close = close
	return r
}

func (r *CRender) Make() Renderer {
	return r
}

func (r *CRender) Clone() RenderBuilder {
	clone := *r   // copy the struct
	return &clone // return a pointer
}

func (r *CRender) RenderLine(a, b string) (ma, mb string) {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(a, b, false)
	for _, diff := range diffs {
		text := html.EscapeString(diff.Text)
		switch diff.Type {
		case diffmatchpatch.DiffDelete:
			ma += r.Text.Rem.Open + text + r.Text.Rem.Close

		case diffmatchpatch.DiffInsert:
			mb += r.Text.Add.Open + text + r.Text.Add.Close

		case diffmatchpatch.DiffEqual:
			fallthrough
		default:
			ma += text
			mb += text
		}
	}
	return
}

func (r *CRender) RenderDiff(unified string) (markup string) {
	var lines []string
	original := strings.Split(unified, "\n")

	var batch *renderBatch

	processBatch := func(lastIdx int) {
		if batch == nil {
			return
		}

		if numDel := len(batch.d); numDel > 0 {
			if numAdd := len(batch.a); numAdd > 0 {
				for idx := range batch.d {
					if idx < numAdd {
						a, b := r.RenderLine(batch.d[idx], batch.a[idx])
						lines[lastIdx-numDel-numAdd+idx] = "-" + a
						lines[lastIdx-numAdd+idx] = "+" + b
					}
				}
			}
		}

		batch = nil
	}

	for idx, line := range original {
		if idx < 2 {
			// skip the patch header lines
			lines = append(lines, line)
			continue
		}
		size := len(line)
		if size == 0 {
			lines = append(lines, "")
			processBatch(idx)
			continue
		}
		lines = append(lines, string(line[0])+html.EscapeString(line[1:]))
		if batch == nil {
			if line[0] == '-' {
				// new batch starting
				batch = &renderBatch{}
				batch.rem(line[1:])
			}
			continue
		}
		// batch in progress
		if line[0] == '-' {
			if len(batch.a) > 0 {
				processBatch(idx)
				batch = &renderBatch{}
			}
			batch.rem(line[1:])
		} else if line[0] == '+' {
			batch.add(line[1:])
		} else {
			processBatch(idx)
		}
	}
	processBatch(len(original))

	for _, line := range lines {
		if size := len(line); size > 0 {
			switch line[0] {
			case '+':
				// line additions
				markup += r.Line.Add.Open + line + r.Line.Add.Close
			case '-':
				// line removals
				markup += r.Line.Rem.Open + line + r.Line.Rem.Close
			case '@', '\\', '#':
				// diff info, comments
				markup += r.Comment.Open + line + r.Comment.Close
			default:
				// unmodified lines and everything else
				markup += r.Normal.Open + line + r.Normal.Close
			}
			markup += "\n"
		} else {
			// can this even happen with a unified diff?
			// every line is supposed to start with at least one char?
		}
	}

	markup = r.File.Open + markup + r.File.Close
	return
}
