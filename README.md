# sort
This tool aims to automatically fix the order of golang imports. It will split your imports into three categories.

## Installation
```
$ go get github.com/AanZee/goimportssort
```

## Usage
```
usage: goimportssort [flags] [path ...]
  -l    write results to stdout
  -local string
        put imports beginning with this string after 3rd-party packages; comma-separated list
  -v    verbose logging
  -w    write result to (source) file instead of stdout (default true)
```
Imports will be sorted according to their categories.
```
$ goimportssort -v -l --local "github.com/AanZee/goimportssort" example.go
```

For example:
```go
package main

import (
	"fmt"
	"log"
	APZ "bitbucket.org/example/package/name"
	APA "bitbucket.org/example/package/name"
	"github.com/AanZee/goimportssort/package2"
	"github.com/AanZee/goimportssort/package1"
)
import (
	"fmt2"
)

import "bitbucket.org/example/package/name2"
import "bitbucket.org/example/package/name3"
import "bitbucket.org/example/package/name4"
```

will be transformed into:

```go
package main

import (
    "fmt"
    "fmt2"
    "log"

    APA "bitbucket.org/example/package/name"
    APZ "bitbucket.org/example/package/name"
    "bitbucket.org/example/package/name2"
    "bitbucket.org/example/package/name3"
    "bitbucket.org/example/package/name4"

    "github.com/AanZee/goimportssort/package1"
    "github.com/AanZee/goimportssort/package2"
)
```