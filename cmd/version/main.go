package main

import (
	"bytes"
	"fmt"
	"log"
	"mksea/common"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func parseVersion(s string) (major int, minor int, patch int, suffix string, err error) {
	if len(s) == 0 {
		return
	}
	if s[0] == 'v' {
		s = s[1:]
	}
	i := strings.Index(s, "-")
	j := strings.Index(s, "+")
	if i == -1 {
		i = j
	} else if j != -1 {
		i = min(i, j)
	}
	if i != -1 {
		suffix = s[i:]
		s = s[:i]
	}
	numbers := strings.Split(s, ".")
	output := []*int{
		&major, &minor, &patch,
	}
	for i, it := range numbers {
		if n, err2 := strconv.ParseInt(it, 10, 64); err2 != nil {
			err = common.NewContextError("version parsing at \""+it+"\"", err2)
			break
		} else {
			*output[i] = int(n)
		}
	}
	return
}

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	file := filepath.Join(wd, "version_init.go")
	content := `package main

func init() {
	Version.Major = %d
	Version.Minor = %d
	Version.Patch = %d
	Version.Suffix = "%s"
}
`

	cmd := exec.Command("git", "tag", "--points-at", "HEAD")
	cmd.Dir = wd

	var out bytes.Buffer
	cmd.Stdout = &out

	err = cmd.Run()
	if err != nil {
		log.Fatalf("%v\noutput:\n%s", err, out.String())
	}

	major, minor, patch, suffix, err := parseVersion(out.String())
	if err != nil {
		log.Fatal(err)
	}
	content = fmt.Sprintf(content, major, minor, patch, suffix)

	err = os.WriteFile(file, []byte(content), 0444)
	if err != nil {
		log.Fatal(err)
	}
}
