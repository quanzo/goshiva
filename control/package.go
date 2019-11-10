package control

import (
	"bufio"
	"fmt"
	"os"
)

var input *os.File
var output *os.File

var exec func(cmd string) string
var stop bool = false

func init() {
	input = os.Stdin
	output = os.Stdout
	exec = func(s string) string {
		return ""
	}
	go listenAndExec()
}

/**
* Установить функцию, которая будет обрабатывать команды с клавиатуры
 */
func SetExecFunc(f func(cmd string) string) {
	exec = f
}

/**
* Остановить ввод команд
*
 */
func Stop() {
	stop = true
}

func listenAndExec() {
	var cmd string = ""
	for !stop {
		in := bufio.NewReader(os.Stdin)
		cmd, _ = in.ReadString('\n')
		if exec != nil {
			res := exec(cmd)
			if output != nil {
				_, _ = fmt.Fprintln(output, res)
			}
		}
	}
}
