package functions

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"
)

func CombineStringArray(arrays ...[]string) []string {
	result := []string{}
	for _, a := range arrays {
		result = append(result, a...)
	}
	return result
}

func GetPathFromDirtyPath(dirtyPath string) (string, error) {
	return GetPathFromArray(strings.Split(dirtyPath, "/"))
}

func GetPathFromArray(pathArray []string) (string, error) {

	sanitizedArray := []string{}
	for _, pa := range pathArray {
		if pa == "GOPATH" {
			gopath := os.Getenv("GOPATH")
			if gopath == "" {
				return "", fmt.Errorf("GOPATH environment is empty")
			} else {
				sanitizedArray = append(sanitizedArray, gopath)
			}
		} else {
			sanitizedArray = append(sanitizedArray, pa)
		}
	}
	return path.Join(sanitizedArray...), nil
}

func IsUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
