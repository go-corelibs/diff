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

// RenderBuilder is the buildable interface for constructing new Renderer instances
type RenderBuilder interface {
	// SetFile specifies the markup wrapping the entire unified diff
	SetFile(open, close string) RenderBuilder
	// SetNormal specifies the markup wrapping each normal line of unified diff
	SetNormal(open, close string) RenderBuilder
	// SetComment specifies the markup wrapping comment lines (starting with a backslash `\` or a hash `#`)
	SetComment(open, close string) RenderBuilder
	// SetLineAdded specifies the markup wrapping lines that were added
	SetLineAdded(open, close string) RenderBuilder
	// SetTextAdded specifies the markup wrapping additions within a line
	SetTextAdded(open, close string) RenderBuilder
	// SetLineRemoved specifies the markup wrapping lines that were removed
	SetLineRemoved(open, close string) RenderBuilder
	// SetTextRemoved specifies the markup wrapping removals within a line
	SetTextRemoved(open, close string) RenderBuilder
	// Make returns the built Renderer instance
	Make() Renderer
}

// Renderer is the interface used to actual render unified diff content
type Renderer interface {
	// RenderLine compares to single-line strings and highlights the differences
	// with the CRender.Text markup strings
	RenderLine(a, b string) (ma, mb string)
	// RenderDiff parses a unified diff string and highlights the interesting
	// details using the CRender.Line, CRender.Comment and CRender.Normal
	// markup strings
	RenderDiff(unified string) (markup string)
	// Clone returns a new RenderBuilder instance configured exactly the same
	// as the one the Clone method is called upon
	Clone() RenderBuilder
}

// CRender implements the RenderBuilder and Renderer interfaces
type CRender struct {
	File    MarkupTag
	Normal  MarkupTag
	Comment MarkupTag
	Line    AddRemTags
	Text    AddRemTags
}

// NewRenderer returns a new RenderBuilder instance
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

func (r *CRender) processRenderDiffBatch(lastIdx int, lines *[]string, batch **renderBatch) {
	if *batch == nil {
		return
	}

	if numDel := len((*batch).d); numDel > 0 {
		if numAdd := len((*batch).a); numAdd > 0 {
			for idx := range (*batch).d {
				if idx < numAdd {
					a, b := r.RenderLine((*batch).d[idx], (*batch).a[idx])
					(*lines)[lastIdx-numDel-numAdd+idx] = "-" + a
					(*lines)[lastIdx-numAdd+idx] = "+" + b
				}
			}
		}
	}

	*batch = nil
}

func (r *CRender) prepareRenderDiff(original []string) (lines []string) {
	var batch *renderBatch
	for idx, line := range original {
		if idx < 2 {
			// skip the patch header lines
			lines = append(lines, line)
			continue
		}
		size := len(line)
		if size == 0 {
			lines = append(lines, "")
			r.processRenderDiffBatch(idx, &lines, &batch)
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
				r.processRenderDiffBatch(idx, &lines, &batch)
				batch = &renderBatch{}
			}
			batch.rem(line[1:])
		} else if line[0] == '+' {
			batch.add(line[1:])
		} else {
			r.processRenderDiffBatch(idx, &lines, &batch)
		}
	}
	r.processRenderDiffBatch(len(original), &lines, &batch)
	return
}

func (r *CRender) RenderDiff(unified string) (markup string) {
	original := strings.Split(unified, "\n")
	lines := r.prepareRenderDiff(original)

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
