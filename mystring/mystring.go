package mystring

import (
	"oreshell/log"
)

// 入力した文字列がクォートされていたらはずすべきだが、展開した値がクォートされていたらはずすべきではないのでは？ todo
func UnescapeAndUnquote(src string) (dst string) {
	log.Logger.Printf("unescapeAndUnquote start : %s\n", src)

	buf := []rune("")
	foundEscape := false
	for _, c := range src {
		if foundEscape {
			if c == '\\' {
				// bufに追加する
			} else if c == '\'' || c == '"' {
				// bufに追加する
			}
			foundEscape = false
		} else {
			if c == '\\' {
				foundEscape = true
				continue // bufに追加しない
			} else if c == '\'' || c == '"' {
				continue // bufに追加しない
			}
		}
		buf = append(buf, c)
	}
	dst = string(buf)
	log.Logger.Printf("unescapeAndUnquote end : %s\n", dst)

	return dst
}
