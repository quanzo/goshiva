package main

/*
Запуск команд из очереди, заданной в файле

При чтении, файл очереди не блокируется и его можно дописывать
Окончание команды - перевод строки - \n
%queue% в команде будет заменен на полный путь к файлу очереди
Если строка начинается с : то это внутренняя команда

*/

import (
	"fmt"
	"io"

	//	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"flag"
	"runtime"
	"strconv"
	"strings"
	"time"

	"./control" // ввод и исполнение команд с консоли
	"./process"
	"./queue"
	"./writecontrol"
)

var delayAfterStartCmd_ms int = 1 // задержка в милисекундах
var queueFile string              // файл очереди
var waitQueueMode = true          // когда очередь опустеет, программа завершится. true - не завершится
var maxProcess = 5                // кол-во одновременно запущенных программ
var cmdCounter = 0                // счетчик, отправленных на исполнение, команд
//var logger, output *log.Logger    //
var showProgramMess bool = true // показывать или нет сообщения программы
var delayQueueNext_ms int = 1   // задержка после получения команды из очереди (если команда пуста)
var waitComplete bool = false   // true - для продолжения выполнения очереди, задания в очереди должны быть выполнены полностью и лишь затем очередь будет продолжена
var stop bool = false           // если установить в true - очередь будет прервана
var absDir string
var cmdQueue *queue.Queue

//
var outputFilename string = ""
var cmdoutFilename string = ""
var programOutput *io.WriteCloser = nil // вывод программы
var commandOutput *io.WriteCloser = nil // вывод работы команд

func main() {
	defer stopped()
	var err error
	startTime := time.Now()

	config()
	verifyConfig()
	showConfig()

	// установим в пакет контроля, функцию исполнения команд
	control.SetExecFunc(controller)

	echo(fmt.Sprint("Start ", time.Now()))

	// создаем очередь из файла
	cmdQueue, err = queue.CreateQueue(queueFile)
	if err != nil {
		if os.IsNotExist(err) { // файл с очередью не найден
			echo(fmt.Sprint("ERROR: Queue file not exists"))
		} else {
			echo(fmt.Sprint(err))
		}
	} else {
		var notEmptyQueue bool = true
		var ifStartCmd = false

		// инициализация процессов
		process.Init(maxProcess)
		process.SetLogger(func(s string) {
			echo(s)
		})

		for !stop && notEmptyQueue { // перебираем очередь
			/*
				Если ждем исполнение очереди, то кол-во запущенных процессов для продолжения == 0. И пока ждем - очередь не читаем
				Или не ждем исполнение очереди, то продолжаем
			*/
			if (waitComplete && process.CountRunning() == 0) || !waitComplete {
				waitComplete = false
				cmd, err := cmdQueue.Next() // выбираем команду из очереди
				if len(cmd) > 0 {           // команда выбрана
					if strings.Index(cmd, ":") == 0 { // в очереди команд обнаружена внутренняя команда
						echo(controller(string([]rune(cmd)[1:])))
					} else {
						ifStartCmd = false
						for !stop && !ifStartCmd {
							// если есть свободные слоты, то будет выполнена команда
							ifStartCmd = process.Start(strconv.Itoa(cmdCounter), cmd, func(s string) {
								fmt.Fprintln(*commandOutput, s)
							})
							if ifStartCmd {
								echo(fmt.Sprintf("#%d Start %s", cmdCounter, cmd))
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

		}

		// подождем завершения всех задач
		for process.CountRunning() > 0 {
			time.Sleep(time.Duration(100) * time.Microsecond) // задержка позволяет снизить загрузку процессора при ожидании выполнения все задач
			runtime.Gosched()
		}
	}
	echo(fmt.Sprintf("Execution time %s", time.Now().Sub(startTime)))
} // end main

func stopped() {
	/*if programOutput != nil {
		_ = (*programOutput).(*os.File).Close()
	}
	if commandOutput != nil {
		_ = (*commandOutput).(*os.File).Close()
	}*/
	if programOutput != nil {
		_ = (*programOutput).Close()
	}
	if commandOutput != nil {
		_ = (*commandOutput).Close()
	}
}

/**
* Установка параметров конфигурации из консоли
*
*
 */
func config() {
	/*var outputFilename string = ""
	var cmdoutFilename string = ""*/

	absDir, _ = os.Getwd()
	absDir = absDir + string(os.PathSeparator)
	queueFile = "./queue.txt"

	//logger = log.New(os.Stdout, "shiva: ", log.LstdFlags)

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

	flag.StringVar(&outputFilename, "output", "", "Programm message output file name")
	flag.StringVar(&outputFilename, "o", "", "Programm message output file name")

	flag.StringVar(&cmdoutFilename, "cmd-output", "", "Command output file name")
	flag.StringVar(&cmdoutFilename, "co", "", "Command output file name")

	flag.Parse()

	//***
	// настройка вывода для программы
	if outputFilename != "" && programOutput == nil {
		if fo, err := writecontrol.ToFile(outputFilename); err == nil {
			programOutput = fo
		} else {
			echo(fmt.Sprintf("ERROR: Not found output file %s\n", outputFilename))
		}
	}
	if programOutput == nil {
		var w io.WriteCloser = &(*os.Stdout)
		programOutput = &w
	}
	//***
	// настройка вывода для результатов работы комманд
	if cmdoutFilename != "" && commandOutput == nil {
		if fo, err := writecontrol.ToFile(cmdoutFilename); err == nil {
			commandOutput = fo
		} else {
			echo(fmt.Sprintf("ERROR: Not found command output file %s\n", cmdoutFilename))
		}
	}
	if commandOutput == nil {
		var w io.WriteCloser = &(*os.Stdout)
		commandOutput = &w
	}

	//***
	/*var w io.Writer = &(*os.Stdout)
	if outputFilename != "" {
		fo, err := os.OpenFile(outputFilename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err == nil {
			w = fo
		} else {
			echo(fmt.Sprintf("ERROR: Not found output file %s\n", outputFilename))
		}
	}
	programOutput = &w
	//***
	var wc io.Writer = &(*os.Stdout)
	if cmdoutFilename != "" {
		fo, err := os.OpenFile(cmdoutFilename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err == nil {
			wc = fo
		} else {
			echo(fmt.Sprintf("ERROR: Not found command output file %s\n", cmdoutFilename))
		}
	}
	commandOutput = &wc
	*/
	fmt.Println("Out file |" + outputFilename + "|" + cmdoutFilename)
} // end config

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
			echo(fmt.Sprintf("ERROR: Not found queue file %s", queueFile))
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
* Используется для обработки встреченных команд из очереди и набранных в консоли
* Это просто функция-исполнитель. Принимают (определяют) другие части программы.
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
		case "rc", "running-count": // вывод кол-ва запущенных команд
			result = "MESSAGE: Count running = " + strconv.Itoa(process.CountRunning())
		case "wait_complete", "wc": // команда предписывает ждать исполнения всех уже запущенных команд для продолжения очереди
			if process.CountRunning() > 0 {
				result = "MESSAGE: Wait complete..."
			}
		case "q", "queue": // показать имя файла-очереди
			result = "MESSAGE: Queue file = " + cmdQueue.GetFile()
		case "s", "status": // текущий статус обработки очереди
			if process.CountRunning() > 0 {
				result = "MESSAGE: Processing queue..."
			} else {
				result = "MESSAGE: Waiting queue..."
			}
		case "stop", "exit", "quit": // остановить обработку очереди ивыйти
			stop = true
			result = "MESSAGE: Stopped... Wait..."
		case "delay", "sleep": // задержка в обработке очереди
			var timeMsec = 1000
			if countParts > 1 {
				paramTimeMsec, _ := strconv.Atoi(parts[1])
				if paramTimeMsec > 0 {
					timeMsec = paramTimeMsec
				}
			}
			time.Sleep(time.Duration(timeMsec) * time.Millisecond)
		case "change-queue", "cq": // начать обработку другого файла-очереди
			//fmt.Println(cmd, parts, countParts)
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
		case "echo": // напечатать сообщение в выводе программы
			message := strings.Trim(strings.Join(parts[1:], " "), " ")
			echo(message)
		}

	}
	return result
}

/**
* Напечать сообщение программы
*
 */
func echo(result string) {
	if showProgramMess && result != "" {
		var writer io.Writer
		if programOutput != nil {
			writer = *programOutput
		} else {
			writer = os.Stdout
		}
		fmt.Fprintln(writer, result)
	}
}

/**
* Напечатать конфигурацию
*
 */
func showConfig() {
	fmt.Println(fmt.Sprintf("Queue file: %s\nMax command running: %d\nWait new command in queue: %t\nDelay after start command: %d(msec)\nShow program message: %t\nSave program message file: %s\nSave queue command output to file: %s\n", queueFile, maxProcess, waitQueueMode, delayAfterStartCmd_ms, showProgramMess, outputFilename, cmdoutFilename))
} // end showConfig
