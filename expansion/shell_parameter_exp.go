package expansion

import (
	"oreshell/constdef"
	"oreshell/infra"
	"oreshell/log"
	"oreshell/myvariables"
	"regexp"
)

/*
func expandShellParametersForArray(src []string) (dst []string) {
	e1 := newVariableNamePartExpander()
	e2 := newVariableNamePartWithBraceExpander()

	for _, arg := range src {
		dst = append(dst, expandShellParameter(
			[]variableNamePartExpander{e1, e2},
			arg))
	}

	return dst
}
*/
/*
	func expandShellParametersForMap(src map[string]string) (dst map[string]string) {
		e1 := newVariableNamePartExpander()
		e2 := newVariableNamePartWithBraceExpander()

		dst = map[string]string{}
		for variable_name, value := range src {
			dst[variable_name] = expandShellParameter(
				[]variableNamePartExpander{e1, e2},
				value)
		}

		return dst
	}
*/
func expandShellParameter(expanders []variableNamePartExpander, src string) (dst string) {

	for _, e := range expanders {
		src = e.expand(src)
	}
	dst = src

	return dst
}

type variableNamePartExpander struct {
	osService            infra.OSService
	namePartRegexp       *regexp.Regexp
	nameInNamePartRegexp *regexp.Regexp
}

func newVariableNamePartExpander() variableNamePartExpander {
	return variableNamePartExpander{
		osService:            infra.MyOSService{},
		namePartRegexp:       regexp.MustCompile(constdef.REGEX_VARIABLE_NAME_PART),
		nameInNamePartRegexp: regexp.MustCompile(constdef.REGEX_VARIABLE_NAME_IN_NAME_PART),
	}
}

func newVariableNamePartWithBraceExpander() variableNamePartExpander {
	return variableNamePartExpander{
		osService:            infra.MyOSService{},
		namePartRegexp:       regexp.MustCompile(constdef.REGEX_VARIABLE_NAME_PART_WITH_BRACE),
		nameInNamePartRegexp: regexp.MustCompile(constdef.REGEX_VARIABLE_NAME_IN_NAME_PART),
	}
}

func (me variableNamePartExpander) expand(src string) string {
	log.Logger.Printf("expandShellParameter before: %v\n", src)
	dst := src
	for {
		submatches := me.namePartRegexp.FindStringSubmatch(dst)
		if len(submatches) == 0 {
			break
		}
		log.Logger.Printf("expandShellParameter FindStringSubmatch 0: %v\n", submatches[0])
		log.Logger.Printf("expandShellParameter FindStringSubmatch 1: %v\n", submatches[1])
		dst = me.namePartRegexp.ReplaceAllStringFunc(src, me.lookupVariables)
	}
	log.Logger.Printf("expandShellParameter after: %v\n", dst)
	return dst
}

func (me variableNamePartExpander) lookupVariables(namePart string) string {
	name := me.nameInNamePartRegexp.FindString(namePart)
	log.Logger.Printf("lookupVariables before: %v\n", name)
	dst := me.osService.Getenv(name)
	if len(dst) == 0 {
		dst = myvariables.Variables().GetValue(name)
	}
	log.Logger.Printf("lookupVariables after: %v\n", dst)
	return dst
}
