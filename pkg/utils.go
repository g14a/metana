package pkg

import (
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
