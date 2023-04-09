package expansion

import (
	"oreshell/log"
	"path/filepath"
)

func expandFilename(src string) (files []string) {
	log.Logger.Printf("expandFilename start: %s\n", src)
	files, err := filepath.Glob(src) // todo bashの挙動と異なるから変えたほうが良い
	log.Logger.Printf("filepath.Glob failed.: %+v\n", err)

	log.Logger.Printf("expandFilename end: %v\n", files)
	return files
}

/*
func expandFilenames(src []string) (dst []string) {
	for _, arg := range src {
		dst = append(dst, expandFilename(arg)...)
	}
	return dst
}
*/
