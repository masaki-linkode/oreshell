package myvariables

import (
	"oreshell/constdef"
	"regexp"
	"strings"
)

func NewAssignVariableParser() AssignVariableParser {
	return AssignVariableParser{
		myRegexp: regexp.MustCompile(constdef.REGEX_VARIABLE_NAME),
	}
}

type AssignVariableParser struct {
	myRegexp *regexp.Regexp
}

func (me AssignVariableParser) TryParse(s string) (ok bool, variable_name string, value string) {
	pair := strings.SplitN(s, "=", 2)
	if len(pair) == 1 {
		return false, "", ""
	}
	// 変数名が空
	if len(pair[0]) == 0 {
		return false, "", ""
	}

	// 変数の文字種のチェック
	if !me.myRegexp.MatchString(pair[0]) {
		return false, "", ""
	}

	//return true, pair[0], mystring.UnescapeAndUnquote(pair[1]) // todo ここでアンクォートしてよいのか？
	return true, pair[0], pair[1]
}
