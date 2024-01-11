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
	"html"
)

type renderBatch struct {
	d []string
	a []string
}

func (b *renderBatch) add(line string) {
	b.a = append(b.a, line)
}

func (b *renderBatch) rem(line string) {
	b.d = append(b.d, line)
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
		if len(line) == 0 {
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
