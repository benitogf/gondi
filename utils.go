package gondi

import "regexp"

var rgx = regexp.MustCompile(`\((.*?)\)`)

func ExtractSourceName(fullName string) string {
	rs := rgx.FindStringSubmatch(fullName)
	if len(rs) < 2 {
		return ""
	}
	return rs[1]
}
