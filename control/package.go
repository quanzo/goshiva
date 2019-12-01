package control

/**
* Пакет запускает слушателя ввода команд с консоли.
* Введенные команды передаются во внешнеюю функцию исполнения.
* Функция-исполнитель задается функцией SetExecFunc
 */

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

var inputSteam io.Reader   // поток приема комманд
var outputStream io.Writer // поток печати результата

var exec func(cmd string) string
var stop bool = false

func init() {
	inputSteam = os.Stdin
	outputStream = os.Stdout
	// функция-исполнитель по умолчанию
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
		in := bufio.NewReader(inputSteam)
		cmd, _ = in.ReadString('\n')
		if exec != nil {
			res := exec(cmd)
			if outputStream != nil {
				_, _ = fmt.Fprintln(outputStream, res)
			}
		}
	}
}
