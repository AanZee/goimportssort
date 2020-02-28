# go-imports-sort
This tool aims to automatically fix the order of golang imports. It will split your imports into three categories.

## Installation
```
$ go get github.com/AanZee/go-imports-sort
```

## Usage
Imports will be sorted according to their categories.
```
$ goimportssort -v -l --local "github.com/AanZee/go-imports-sort" example.go
```

For example:
```go
package main

import (
	"fmt"
	"log"
	APZ "bitbucket.org/example/package/name"
	APA "bitbucket.org/example/package/name"
	"github.com/AanZee/go-imports-sort/package2"
	"github.com/AanZee/go-imports-sort/package1"
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

    "github.com/AanZee/go-imports-sort/package1"
    "github.com/AanZee/go-imports-sort/package2"
)
```