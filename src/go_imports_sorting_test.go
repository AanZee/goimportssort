package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestProcessFile(t *testing.T) {
	asserts := assert.New(t)
	*localPrefix = "github.com/AanZee/go-imports-sort"
	*write = false
	want := `package main

// builtin
// external
// local
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

func main() {
	fmt.Println("Hello!")
}`

	output, err := processFile("../_fixtures/sample.txt", nil, os.Stdout)
	asserts.NotEqual(nil, output)
	asserts.Equal(nil, err)
	asserts.Equal(want, string(output))
}

func TestProcessFile_Equal(t *testing.T) {
	asserts := assert.New(t)
	*localPrefix = "github.com/AanZee/go-imports-sort"
	*write = false

	output, err := processFile("../_fixtures/sample2.txt", nil, os.Stdout)
	asserts.NotEqual(nil, output)
	asserts.Equal(nil, err)
	asserts.Equal([]byte(nil), output)
}