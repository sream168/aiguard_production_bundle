package strutil

import "strings"

// FirstNonEmpty 返回第一个非空（trim 后）字符串，全部为空时返回空字符串。
func FirstNonEmpty(values ...string) string {
	for _, v := range values {
		v = strings.TrimSpace(v)
		if v != "" {
			return v
		}
	}
	return ""
}
