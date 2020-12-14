package main

import (
  "fmt"
  "bufio"
  "log"
  "os"
  "strings"
)

func main() {

  // 標準入力から文字列を読み取る準備
  reader := bufio.NewReader(os.Stdin)

  // 0.ずっとループ
  for {
    // 1.プロンプトを表示してユーザに入力を促す
    fmt.Printf("(ore) > ")
    
    // 3.標準入力から文字列(コマンド)を読み込む
    line, _, err := reader.ReadLine()
    if err != nil {
      log.Fatalf("ReadLine %v", err)
    }
    words := strings.Split(string(line), " ")

    // これから起動するプログラムの出力と自分の出力をつなげる(6,7)
    var procAttr os.ProcAttr
    procAttr.Files = []*os.File{nil, os.Stdout, os.Stderr}

    // 4.入力文字列に該当するプログラムを探して起動する
    process, err := os.StartProcess(words[0], words, &procAttr)
    if err != nil {
      log.Fatalf("StartProcess %v", err)
    }

    // 起動したプログラムが終了するまで待つ(8を待つ)
    _, err = process.Wait()
    if err != nil {
      log.Fatalf("Wait %v", err)
    }
  }
}