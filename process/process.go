package process

import (
	"oreshell/ast"
	"os"
)

const (
	FD_DEFAULT_IN  = 0
	FD_DEFAULT_OUT = 1
	FD_DEFAULT_ERR = 2
)

const (
	FD_MIN = 0
	FD_MAX = 9
)

func CreateProcAttrFiles(redirections *[]ast.Redirection) (files []*os.File, err error) {

	// FDをキー、入出力先ファイルを値とした辞書(初期値付き)
	m := map[int]*os.File{FD_DEFAULT_IN: os.Stdin, FD_DEFAULT_OUT: os.Stdout, FD_DEFAULT_ERR: os.Stderr}

	// redirectionsから辞書へ
	for _, v := range *redirections {
		var f *os.File
		if v.Direction == ast.IN {
			// 入力用ファイルオープン
			f, err = os.Open(v.FilePath)
		} else { // ast.OUT
			// 出力用ファイルオープン
			f, err = os.Create(v.FilePath)
		}

		if err != nil {
			return nil, err
		}

		m[v.FdNum] = f
	}

	// 辞書からFileの配列へ
	files = []*os.File{}
	for i := FD_MIN; i <= FD_MAX; i++ {
		v, ok := m[i]
		if ok {
			files = append(files, v)
		}
	}

	return files, nil
}
