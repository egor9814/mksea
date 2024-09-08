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

var wd string

func init() {
	var err error
	wd, err = os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
}

func git(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = wd

	var out, errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut

	err := cmd.Run()
	s := out.String()
	if err != nil {
		return "", fmt.Errorf("%v\noutput:\n%s%s", err, s, errOut.String())
	}

	return strings.TrimSpace(s), nil
}

func branchName() (string, error) {
	s, err := git("branch", "--no-color", "--show-current")
	if err != nil {
		return "", err
	}
	if len(s) == 0 {
		s = "untracked"
	} else {
		s = "dev-" + s
	}
	return s, nil
}

func parseVersion(s string) (major int, minor int, patch int, suffix string, err error) {
	if len(s) == 0 {
		suffix, err = branchName()
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

func getVersion() (int, int, int, string, error) {
	s, err := git("tag", "--points-at", "HEAD")
	if err != nil {
		log.Fatal(err)
	}
	return parseVersion(s)
}

func main() {
	file := filepath.Join(wd, "version_init.go")
	content := `package main

func init() {
	Version.Major = %d
	Version.Minor = %d
	Version.Patch = %d
	Version.Suffix = "%s"
}
`

	major, minor, patch, suffix, err := getVersion()
	if err != nil {
		log.Fatal(err)
	}
	if len(suffix) > 0 {
		suffix = "-" + suffix
	}
	content = fmt.Sprintf(content, major, minor, patch, suffix)

	err = os.WriteFile(file, []byte(content), 0655)
	if err != nil {
		log.Fatal(err)
	}
}
