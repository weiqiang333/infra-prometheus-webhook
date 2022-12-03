package str

import (
	"fmt"
	"strings"
)

// EscapeStrings 转义字符. 将 str 中含有 character 的所有特定字符进行转义. 增加 "\"
func EscapeStrings(str string, character []string) string {
	for _, old := range character {
		str = strings.ReplaceAll(str, old, fmt.Sprintf("\\\\%s", old))
	}
	return str
}
