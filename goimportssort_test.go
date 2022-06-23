package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func TestProcessFile(t *testing.T) {
	asserts := assert.New(t)
	*localPrefix = "github.com/AanZee/goimportssort"
	reader := strings.NewReader(`package main

// builtin
// external
// local
import (
	"fmt"
	"log"
	
	APA "bitbucket.org/example/package/name"
	APZ "bitbucket.org/example/package/name"
	"bitbucket.org/example/package/name2"
	"bitbucket.org/example/package/name3" // foopsie
	"bitbucket.org/example/package/name4"
	
	"github.com/AanZee/goimportssort/package1"
	// a
	"github.com/AanZee/goimportssort/package2"
	
	/*
		mijn comment
	*/
	"net/http/httptest"
	"database/sql/driver"
)
// klaslkasdko

func main() {
	fmt.Println("Hello!")
}`)
	want := `package main

import (
	"database/sql/driver"
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

func main() {
	fmt.Println("Hello!")
}
`

	output, err := processFile("", reader, os.Stdout)
	asserts.NotEqual(nil, output)
	asserts.Equal(nil, err)
	asserts.Equal(want, string(output))
}

func TestProcessFile_SingleImport(t *testing.T) {
	asserts := assert.New(t)
	*localPrefix = "github.com/AanZee/goimportssort"

	reader := strings.NewReader(
		`package main


import "github.com/AanZee/goimportssort/package1"


func main() {
	fmt.Println("Hello!")
}`)
	output, err := processFile("", reader, os.Stdout)
	asserts.NotEqual(nil, output)
	asserts.Equal(nil, err)
	asserts.Equal(
		`package main

import (
	"github.com/AanZee/goimportssort/package1"
)

func main() {
	fmt.Println("Hello!")
}
`, string(output))
}

func TestProcessFile_EmptyImport(t *testing.T) {
	asserts := assert.New(t)
	*localPrefix = "github.com/AanZee/goimportssort"

	reader := strings.NewReader(`package main

func main() {
	fmt.Println("Hello!")
}`)
	output, err := processFile("", reader, os.Stdout)
	asserts.NotEqual(nil, output)
	asserts.Equal(nil, err)
	asserts.Equal(`package main

func main() {
	fmt.Println("Hello!")
}`, string(output))
}

func TestProcessFile_ReadMeExample(t *testing.T) {
	asserts := assert.New(t)
	*localPrefix = "github.com/AanZee/goimportssort"

	reader := strings.NewReader(`package main

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
import "bitbucket.org/example/package/name4"`)
	output, err := processFile("", reader, os.Stdout)
	asserts.NotEqual(nil, output)
	asserts.Equal(nil, err)
	asserts.Equal(`package main

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
`, string(output))
}

func TestProcessFile_WronglyFormattedGo(t *testing.T) {
	asserts := assert.New(t)
	*localPrefix = "github.com/AanZee/goimportssort"

	reader := strings.NewReader(
		`package main
import "github.com/AanZee/goimportssort/package1"


func main() {
	fmt.Println("Hello!")
}`)
	output, err := processFile("", reader, os.Stdout)
	asserts.NotEqual(nil, output)
	asserts.Equal(nil, err)
	asserts.Equal(
		`package main

import (
	"github.com/AanZee/goimportssort/package1"
)

func main() {
	fmt.Println("Hello!")
}
`, string(output))
}

func TestGetModuleName(t *testing.T) {
	asserts := assert.New(t)

	name := getModuleName()

	asserts.Equal("github.com/AanZee/goimportssort", name)
}

func TestWithEmptyLines(t *testing.T) {
	asserts := assert.New(t)
	*localPrefix = "github.com/AanZee/goimportssort"
	*winLines = true
	reader := strings.NewReader("package main\r\n\r\nimport (\r\n\t\"fmt\"\r\n\t\"log\"\r\n\r\n\t\"bitbucket.org/example/package/name2\"\r\n\t\r\n\t\"net/http\"\r\n)\r\n\r\nfunc main() {\r\n\tfmt.Println(\"Hello!\")\r\n\r\n\tif true {\r\n\t\tfmt.Println(\"Its true\")\r\n\t}\r\n}\r\n")
	want := "package main\n\nimport (\n\t\"fmt\"\n\t\"log\"\n\t\"net/http\"\n\n\t\"bitbucket.org/example/package/name2\"\n)\n\nfunc main() {\n\tfmt.Println(\"Hello!\")\n\n\tif true {\n\t\tfmt.Println(\"Its true\")\n\t}\n}\n"

	output, err := processFile("", reader, os.Stdout)
	asserts.NotEqual(nil, output)
	asserts.Equal(nil, err)

	asserts.Equal(want, string(output))

	*winLines = false
}
