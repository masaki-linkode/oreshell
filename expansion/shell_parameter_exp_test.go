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

func (me myOSServiceMock) Hasenv(name string) bool {
	return len(me.m[name]) > 0
}

func (me myOSServiceMock) Setenv(name string, val string) error {
	me.m[name] = val
	return nil
}

func newExpanderWithMyOSServiceMock() variableNamePartExpander {
	me := newVariableNamePartWithBraceExpander()
	me.osService = myOSServiceMock{m: map[string]string{}}
	return me
}

func Test_osGetEnvByVariableNamePartWithBrace(t *testing.T) {
	me := newExpanderWithMyOSServiceMock()
	me.osService.Setenv("HOGE", "hige")
	assert.Equal(t, "hige", me.lookupVariables("${HOGE}"))
}

func Test_ExpandShellParameter(t *testing.T) {
	me := newExpanderWithMyOSServiceMock()
	me.osService.Setenv("HOGE", "hige")
	me.osService.Setenv("_HUGE", "hege")

	assert.Equal(t, "  hige_hege  ", me.expand("  ${HOGE}_${_HUGE}  "))
}
