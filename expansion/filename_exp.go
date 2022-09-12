package expansion

import (
	"oreshell/log"
	"path/filepath"
)

func expandFilename(src string) []string {
	log.Logger.Printf("extractArgs before: %s\n", src)
	files, _ := filepath.Glob(src)
	if files == nil { // argはワイルドカード文字列ではなかった、もしくはワイルドカード文字列だがヒットしなかった
		files = []string{unescapeAndUnquote(src)}
	}
	log.Logger.Printf("extractArgs after: %v\n", files)
	return files
}

func expandFilenames(src []string) (dst []string) {
	for _, arg := range src {
		dst = append(dst, expandFilename(arg)...)
	}
	return dst
}

func unescapeAndUnquote(src string) (dst string) {
	log.Logger.Printf("unescapeAndUnquote before : %s\n", src)

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
	log.Logger.Printf("unescapeAndUnquote after : %s\n", dst)

	return dst
}
