package expansion

import (
	"oreshell/constdef"
	"oreshell/infra"
	"oreshell/log"
	"regexp"
)

func expandShellParameters(src []string) (dst []string) {
	e1 := newEnvVariableNamePartExpander()
	e2 := newEnvVariableNamePartWithBraceExpander()

	for _, arg := range src {
		r1 := e1.expand(arg)
		r2 := e2.expand(r1)
		dst = append(dst, r2)
	}

	return dst
}

type envVariableNamePartExpander struct {
	osService            infra.OSService
	namePartRegexp       *regexp.Regexp
	nameInNamePartRegexp *regexp.Regexp
}

func newEnvVariableNamePartExpander() envVariableNamePartExpander {
	return envVariableNamePartExpander{
		osService:            infra.MyOSService{},
		namePartRegexp:       regexp.MustCompile(constdef.REGEX_VARIABLE_NAME_PART),
		nameInNamePartRegexp: regexp.MustCompile(constdef.REGEX_VARIABLE_NAME_IN_NAME_PART),
	}
}

func newEnvVariableNamePartWithBraceExpander() envVariableNamePartExpander {
	return envVariableNamePartExpander{
		osService:            infra.MyOSService{},
		namePartRegexp:       regexp.MustCompile(constdef.REGEX_VARIABLE_NAME_PART_WITH_BRACE),
		nameInNamePartRegexp: regexp.MustCompile(constdef.REGEX_VARIABLE_NAME_IN_NAME_PART),
	}
}

func (me envVariableNamePartExpander) expand(src string) string {
	log.Logger.Printf("expandShellParameter before: %v\n", src)
	dst := src
	for {
		submatches := me.namePartRegexp.FindStringSubmatch(dst)
		if len(submatches) == 0 {
			break
		}
		log.Logger.Printf("ExpandVariablePartFirst 0: %v\n", submatches[0])
		log.Logger.Printf("ExpandVariablePartFirst 1: %v\n", submatches[1])
		dst = me.namePartRegexp.ReplaceAllStringFunc(src, me.osGetEnvByVariableNamePart)
	}
	log.Logger.Printf("expandShellParameter after: %v\n", dst)
	return dst
}

func (me envVariableNamePartExpander) osGetEnvByVariableNamePart(namePart string) string {
	name := me.nameInNamePartRegexp.FindString(namePart)
	log.Logger.Printf("osGetEnvByVariableNamePartWithBrace before: %v\n", name)
	dst := me.osService.Getenv(name)
	log.Logger.Printf("osGetEnvByVariableNamePartWithBrace after: %v\n", dst)
	return dst
}
