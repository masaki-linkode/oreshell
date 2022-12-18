package expansion

import (
	"oreshell/log"
	"oreshell/myfiles"
	"oreshell/mystring"
)

func ExpandWords(src []string) (dst []string) {
	log.Logger.Printf("ExpandWords start: %+v\n", src)

	e1 := newVariableNamePartExpander()
	e2 := newVariableNamePartWithBraceExpander()

	for _, arg := range src {
		v1 := arg
		v2 := expandShellParameter([]variableNamePartExpander{e1, e2}, v1)
		if myfiles.Exists(v2) {
			v3 := mystring.UnescapeAndUnquote(v2)
			dst = append(dst, v3)
		} else {
			filenames := expandFilename(v2)
			if filenames != nil {
				dst = append(dst, filenames...)
			} else {
				v3 := mystring.UnescapeAndUnquote(v2)
				dst = append(dst, v3)
			}
		}
	}

	log.Logger.Printf("ExpandWords end: %+v\n", dst)
	return dst
}

func ExpandVarableValues(src map[string]string) (dst map[string]string) {
	log.Logger.Printf("ExpandVarableValues start: %+v\n", src)

	e1 := newVariableNamePartExpander()
	e2 := newVariableNamePartWithBraceExpander()

	dst = map[string]string{}
	for variable_name, value := range src {

		v1 := value
		v2 := expandShellParameter([]variableNamePartExpander{e1, e2}, v1)
		v3 := mystring.UnescapeAndUnquote(v2)

		dst[variable_name] = v3
	}

	log.Logger.Printf("ExpandVarableValues end: %+v\n", dst)
	return dst
}
