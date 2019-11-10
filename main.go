// shiva project main.go
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"flag"
	"runtime"
	"strconv"
	"strings"
	"time"

	"./control"
	"./process"
	"./queue"
)

var delayAfterStartCmd_ms int = 1 // задержка в милисекундах
var queueFile string              // файл очереди
var waitQueueMode = true          // когда очередь опустеет, программа завершится. true - не завершится
var maxProcess = 5                // кол-во одновременно запущенных программ
var cmdCounter = 0                // счетчик, отправленных на исполнение, команд
var logger, output *log.Logger    //
var showProgramMess bool = true   // показывать или нет сообщения программы
var delayQueueNext_ms int = 1     // задержка после получения команды из очереди (если команда пуста)
var stop bool = false
var cmdOutput *io.Writer = nil
var absDir string
var cmdQueue *queue.Queue

func main() {
	var err error
	startTime := time.Now()

	config()
	verifyConfig()

	// установим в пакет контроля, функцию исполнения команд
	control.SetExecFunc(controller)

	if showProgramMess {
		logger.Println("Start", time.Now())
	}

	// создаем очередь из файла
	cmdQueue, err = queue.CreateQueue(queueFile)
	if err != nil {
		if os.IsNotExist(err) { // файл с очередью не найден
			logger.Println("ERROR: Queue file not exists")
		} else {
			logger.Println(err)
		}
	} else {
		var notEmptyQueue bool = true
		var ifStartCmd = false

		// инициализация процессов
		process.Init(maxProcess)
		process.SetLogger(func(s string) {
			logger.Println(s)
		})

		for !stop && notEmptyQueue { // перебираем очередь
			cmd, err := cmdQueue.Next() // выбираем команду из очереди
			if len(cmd) > 0 {
				if strings.Index(cmd, ":") == 0 { // в очереди команд обнаружена внутренняя команда
					result := controller(string([]rune(cmd)[1:]))
					if cmdOutput != nil && result != "" {
						fmt.Fprintln(*cmdOutput, result)
					}
				} else {
					ifStartCmd = false
					for !stop && !ifStartCmd {
						// если есть свободные слоты, то будет выполнена команда
						ifStartCmd = process.Start(strconv.Itoa(cmdCounter), cmd, func(s string) {
							fmt.Fprintln(*cmdOutput, s)
						})
						if ifStartCmd {
							if showProgramMess {
								logger.Printf("#%d Start %s", cmdCounter, cmd)
							}
							if delayAfterStartCmd_ms > 0 {
								time.Sleep(time.Duration(delayAfterStartCmd_ms) * time.Millisecond)
							}
							cmdCounter++
						}
						runtime.Gosched() // выделим таймслот для горутин
					}
				}

			} else if delayQueueNext_ms > 0 { // задержка опроса очереди (если там ничего нет)
				time.Sleep(time.Duration(delayQueueNext_ms) * time.Millisecond)
			}
			if err != nil {
				if !waitQueueMode { // если установлена опция ожидания данных в очереди
					notEmptyQueue = false // продолжаем цмкл и ожидаем данные в очереди
				}
			}
		}

		// подождем завершения всех задач
		for process.CountRunning() > 0 {
			time.Sleep(time.Duration(100) * time.Microsecond) // задержка позволяет снизить загрузку процессора при ожидании выполнения все задач
			runtime.Gosched()
		}
	}
	if showProgramMess {
		logger.Printf("Execution time %s", time.Now().Sub(startTime))
	}
} // end main

/**
* Установка параметров конфигурации из консоли
*
*
 */
func config() {
	absDir, _ = os.Getwd()
	absDir = absDir + string(os.PathSeparator)

	logger = log.New(os.Stdout, "shiva: ", log.LstdFlags)
	queueFile = "./queue.txt"

	flag.StringVar(&queueFile, "queue-file", "./queue.txt", "Queue filename")
	flag.StringVar(&queueFile, "q", "./queue.txt", "Queue filename")

	flag.IntVar(&maxProcess, "max-process", 5, "Max process")
	flag.IntVar(&maxProcess, "p", 5, "Max process")

	flag.BoolVar(&waitQueueMode, "wait-queue", true, "Wait command in queue")
	flag.BoolVar(&waitQueueMode, "w", true, "Wait command in queue")

	flag.IntVar(&delayAfterStartCmd_ms, "sleep-time", 0, "Sleep time (in msec) after start command")
	flag.IntVar(&delayAfterStartCmd_ms, "s", 0, "Sleep time (in msec) after start command")

	flag.BoolVar(&showProgramMess, "messages", true, "Show program messages")
	flag.BoolVar(&showProgramMess, "m", true, "Show program messages")

	var outputFilename string
	var w io.Writer = os.Stdout
	flag.StringVar(&outputFilename, "out-file", "", "Queue command output file name")
	flag.StringVar(&outputFilename, "of", "", "Queue command output file name")
	if outputFilename != "" {
		fo, err := os.OpenFile(outputFilename, os.O_APPEND|os.O_WRONLY, 0644)
		if err == nil {
			w = fo
		} else {
			logger.Printf("ERROR: Not found output file %s\n", outputFilename)
		}
	}
	cmdOutput = &w
}

/**
* Проверка конфигурации программы
*
 */
func verifyConfig() {
	var errCounter = 0
	if queueFile != "" {
		pathQueue, err := exec.LookPath(queueFile)
		if err == nil {
			if !path.IsAbs(pathQueue) {
				pathQueue, _ = filepath.Abs(pathQueue)
			}
			queueFile = pathQueue
		} else {
			logger.Printf("ERROR: Not found queue file %s", queueFile)
			errCounter++
		}
	}
	if errCounter > 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}
} // end verifyConfig

/**
* Исполнение команд
*
 */
func controller(cmd string) string {
	var result string = ""

	// части команды
	parts := strings.Fields(cmd)
	//parts := strings.Split(cmd, " ")
	//fmt.Println(parts, cmd)
	countParts := len(parts)

	if countParts > 0 {
		switch parts[0] {
		case "rc", "running-count":
			result = "MESSAGE: Count running = " + strconv.Itoa(process.CountRunning())
		case "q", "queue":
			result = "MESSAGE: Queue file = " + cmdQueue.GetFile()
		case "s", "status":
			if process.CountRunning() > 0 {
				result = "MESSAGE: Processing queue..."
			} else {
				result = "MESSAGE: Waiting queue..."
			}
		case "stop", "exit", "quit":
			stop = true
			result = "MESSAGE: Stopped... Wait..."
		case "delay", "sleep":
			var timeMsec = 1000
			if countParts > 1 {
				paramTimeMsec, _ := strconv.Atoi(parts[1])
				if paramTimeMsec > 0 {
					timeMsec = paramTimeMsec
				}
			}
			time.Sleep(time.Duration(timeMsec) * time.Millisecond)
		case "change-queue", "cq":
			fmt.Println(cmd, parts, countParts)
			if countParts > 1 {
				newQueueFile := strings.Trim(strings.Join(parts[1:], " "), " ")
				if cmdQueue.ChangeQueue(newQueueFile) {
					result = "OK: Success change"
				} else {
					result = "ERROR: Not change queue file"
				}
			} else {
				result = "ERROR: Not set queue file param from command"
			}
		}
	}
	return result
}
