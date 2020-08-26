# go-imports-sort
This tool aims to automatically fix the order of golang imports. It will split your imports into three categories.

## Installation
```
$ go get -u github.com/AanZee/goimportssort
```

## Usage
```
usage: goimportssort [flags] [path ...]
  -l    write results to stdout (default false)
  -local string
        put imports beginning with this string after 3rd-party packages; comma-separated list 
(default tries to get module name of current directory)
  -v    verbose logging (default false)
  -w    write result to (source) file (default false)
```
Imports will be sorted according to their categories.
```
$ goimportssort -v -w ./..
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
	"net/http/httptest"
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
    "log"
    "net/http/httptest"

    APA "bitbucket.org/example/package/name"
    APZ "bitbucket.org/example/package/name"
    "bitbucket.org/example/package/name2"
    "bitbucket.org/example/package/name3"
    "bitbucket.org/example/package/name4"

    "github.com/AanZee/goimportssort/package1"
    "github.com/AanZee/goimportssort/package2"
)
```