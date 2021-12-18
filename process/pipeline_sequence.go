package process

import (
	"oreshell/ast"
	"oreshell/log"
)

type PipelineSequence struct {
	processes []*Process
}

func NewPipelineSequence(pipelineSequence *ast.PipelineSequence) (me *PipelineSequence, err error) {

	me = &PipelineSequence{}

	for _, sc := range pipelineSequence.SimpleCommands {
		p, err := NewProcess(sc)
		if err != nil {
			return nil, err
		}
		log.Logger.Printf("Process New %+v\n", p)

		me.processes = append(me.processes, p)
	}
	log.Logger.Printf("processes size %d\n", me.size())

	for i, p := range me.processes {
		//プロセス数が1つの場合はパイプ設定対象外
		//プロセスが複数でも末尾のプロセスはパイプ設定対象外。
		if i+1 < me.size() {
			log.Logger.Printf("pipe %+v\n", p)
			err = p.PipeWithNext(me.processes[i+1])
			if err != nil {
				return nil, err
			}
		}
	}

	return me, nil
}

func (me *PipelineSequence) size() int {
	if me.processes == nil {
		return 0
	} else {
		return len(me.processes)
	}
}

func (me *PipelineSequence) Exec() (err error) {

	for _, p := range me.processes {
		log.Logger.Printf("Process Start %+v\n", p)
		err = p.Start()
		if err != nil {
			return err
		}
	}

	// 起動したプログラムが終了するまで待つ
	for _, p := range me.processes {
		log.Logger.Printf("Process Wait %+v\n", p)
		_ = p.Wait()
	}

	log.Logger.Printf("processes done \n")

	return nil
}
