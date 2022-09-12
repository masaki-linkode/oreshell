package expansion

import (
	"oreshell/log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	log.Logger = log.New()
}

type myOSServiceMock struct {
	m map[string]string
}

func (me myOSServiceMock) Getenv(name string) string {
	return me.m[name]
}

func (me myOSServiceMock) Setenv(name string, val string) error {
	me.m[name] = val
	return nil
}

func newExpanderWithMyOSServiceMock() envVariableNamePartExpander {
	me := newEnvVariableNamePartWithBraceExpander()
	me.osService = myOSServiceMock{m: map[string]string{}}
	return me
}

func Test_osGetEnvByVariableNamePartWithBrace(t *testing.T) {
	me := newExpanderWithMyOSServiceMock()
	me.osService.Setenv("HOGE", "hige")
	assert.Equal(t, "hige", me.osGetEnvByVariableNamePart("${HOGE}"))
}

func Test_ExpandShellParameter(t *testing.T) {
	me := newExpanderWithMyOSServiceMock()
	me.osService.Setenv("HOGE", "hige")
	me.osService.Setenv("_HUGE", "hege")

	assert.Equal(t, "  hige_hege  ", me.expand("  ${HOGE}_${_HUGE}  "))
}
