package common

import "regexp"

var ansiRegexp = regexp.MustCompile(
	`\x1b\[[0-9;]*[a-zA-Z]`,
)

func StripANSI(s string) string {
	return ansiRegexp.ReplaceAllString(s, "")
}
