# go-corelibs/diff - Unified diff utilities

diff is a Go package for computing the differences between two strings, can
interact with the result to selectively pick groups of changes and generate
unified diff output.

# Installation

``` shell
> go get github.com/go-corelibs/diff@latest
```

# Description

## Diff

Get a unified diff between two strings.

``` go
original := `This is the first line
This is the second line`
modified := strings.Replace(original, "the first", "one", 1)
delta := diff.New("filename.txt", original, modified)
if unified, err := delta.Unified(); err != nil {
    panic(err)
} else {
    fmt.Println(unified)
}
```

With the output being:

``` diff
--- a/filename.txt
+++ b/filename.txt
@@ -1,2 +1,2 @@
-This is the first line
+This is one line
 This is the second line
\ No newline at end of file
```

## Renderer

Sometimes we like to render unified diffs for users in ways beyond plain text,
such as HTML or the Go-Curses Tango markup format.

``` go
// using the unified variable from the Diff example
output := diff.HTMLRenderer.RenderDiff(unified)
fmt.Println(output)
```

Produces the following:

``` html
<ul style="list-style-type:none;margin:0;padding:0;">
<li style="color:#eeeeee;background-color:#770000;">--- a/filename.txt</li>
<li style="color:#ffffff;background-color:#007700;">+++ b/filename.txt</li>
<li style="font-style:italic;opacity:0.77;">@@ -1,2 +1,2 @@</li>
<li style="color:#eeeeee;background-color:#770000;">-This is <span style="background-color:#440000;opacity:0.77;text-decoration:line-through;">th</span>e<span style="background-color:#440000;opacity:0.77;text-decoration:line-through;"> first</span> line</li>
<li style="color:#ffffff;background-color:#007700;">+This is <span style="background-color:#004400;font-weight:bold;">on</span>e line</li>
<li style="opacity:0.77;"> This is the second line</li>
<li style="font-style:italic;opacity:0.77;">\ No newline at end of file</li>
</ul>
```

# Go-CoreLibs

[Go-CoreLibs] is a repository of shared code between the [Go-Curses] and
[Go-Enjin] projects.

# License

``` 
Copyright 2023 The Go-CoreLibs Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use file except in compliance with the License.
You may obtain a copy of the license at

 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```

[Go-CoreLibs]: https://github.com/go-corelibs
[Go-Curses]: https://github.com/go-curses
[Go-Enjin]: https://github.com/go-enjin
