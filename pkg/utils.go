package pkg

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
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

func ExecuteCommand(root *cobra.Command, args ...string) (output string, err error) {
	_, output, err = ExecuteCommandC(root, args...)
	return output, err
}

func ExecuteCommandC(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	c, err = root.ExecuteC()

	return c, buf.String(), err
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

func GetExpectedLinesInit() []string {
	return []string{`// This file is auto generated. DO NOT EDIT!`,
		`MigrateUp`,
		`upUntil string, lastRunTS int`,
		`return nil`,
		`MigrateDown`,
		`downUntil string, lastRunTS int`,
		`return nil`,
		`func main()`,
		`upCmd := flag.NewFlagSet`,
		`downCmd := flag.NewFlagSet`,
		`var upUntil, downUntil string`,
		`var lastRunTS int`,
		`upCmd.StringVar`,
		`upCmd.IntVar`,
		`downCmd.StringVar`,
		`downCmd.IntVar`,
		`switch`,
		`case "up"`,
		`err := upCmd.Parse`,
		`if err != nil {`,
		`return`,
		`}`,
		`case "down"`,
		`err := downCmd.Parse`,
		`if err != nil {`,
		`return`,
		`}`,
		`if upCmd.Parsed()`,
		`MigrateUp()`,
		`if err != nil {`,
		`fmt.Fprintf()`,
		`}`,
		`if downCmd.Parsed()`,
		`MigrateDown()`,
		`if err != nil {`,
		`fmt.Fprintf()`,
		`}`,
	}
}
