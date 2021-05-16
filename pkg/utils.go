package pkg

import (
	"regexp"
	"strconv"
	"strings"
)

func GetComponents(filename string) (timestamp int, migrationName string, err error) {
	trimmedName := strings.TrimSuffix(filename, ".go")
	components := strings.Split(trimmedName, "-")

	timestamp, err = strconv.Atoi(components[0])
	if err != nil {
		return 0, "", err
	}

	return timestamp, components[1], nil
}

type T interface {
	Helper()
	Errorf(string, ...interface{})
}

func ExpectLines(t T, output string, lines ...string) {
	t.Helper()
	var r *regexp.Regexp
	for _, l := range lines {
		r = regexp.MustCompile(l)
		if !r.MatchString(output) {
			t.Errorf("output did not match regexp /%s/\n> output\n%s\n", r, output)
			return
		}
	}
}
