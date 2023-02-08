package matcher

import (
	"strings"
)

func MatchPath(pathParts []string, pattern string) bool {

	if pattern == "/" && pathParts[0] == "" && pathParts[1] == "" {
		return true
	}

	patternPaths := strings.Split(pattern, "/")

	if len(pathParts) != len(patternPaths) {
		return false
	}

	mountedMatch := ""

	for i, pathPart := range pathParts {
		patternPath := patternPaths[i]

		if i == len(pathParts)-1 {
			pathPart = strings.Split(pathPart, "?")[0]
		}

		if patternPath == pathPart && len(patternPath) == 0 {
			continue
		}

		if pathPart == patternPath && !strings.HasPrefix(patternPath, ":") {
			mountedMatch += "/" + pathPart
		} else if strings.HasPrefix(patternPath, ":") {
			mountedMatch += "/" + patternPath
		} else {
			return false
		}
	}

	if mountedMatch == pattern {
		return true
	}

	return false
}
