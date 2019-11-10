package process

import (
	"fmt"
	//	"io"
	//	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
)

var maxProcess int = 5
var processes []*process
var lock *sync.Mutex
var output func(s string)
var logger func(s string)

func Init(maxPrc int) {
	maxProcess = maxPrc
	processes = make([]*process, maxPrc, maxPrc*2)
	lock = new(sync.Mutex)
	logger = nil
	output = nil
}

func SetLogger(_logger func(line string)) {
	logger = _logger
}

type process struct {
	id             string // id процесса
	cmd            string // команда
	running        bool   // признак работы
	processPointer *os.Process
	output         func(line string)
}

/**
* Старт команды. Команда будет запущена если будут свободные слоты.
* Если свободных нет - будет возвращено false
 */
func Start(id string, cmd string, output func(line string)) bool {
	lock.Lock()
	defer lock.Unlock()
	if cmd != "" {
		prc := new(process)
		prc.cmd = cmd
		prc.id = id
		prc.output = output
		if slot := getEmptySlotIndex(); slot != -1 {
			processes[slot] = prc
			go prc.run()
			return true
		}
	}
	return false
}

/**
* Запустить команду
*
 */
func (this *process) run() {
	this.running = true
	defer func() {
		this.running = false
	}()
	// части команды
	parts := strings.Fields(this.cmd)
	// определим путь до команды
	pathCmd, err := exec.LookPath(parts[0])
	if err == nil {
		var args []string
		if len(parts) > 1 { // аргументы команды
			args = parts[1:]
		}
		// создадим команду
		cmd := exec.Command(pathCmd)
		cmd.Args = args
		// запустим команду и получим ее вывод
		cmdOutput, err := cmd.CombinedOutput()

		//runProcess, err := os.StartProcess(pathCmd, args, &os.ProcAttr{Files: []*os.File{os.Stdin, output, os.Stderr}})
		//runPrcPointer, err := os.StartProcess(pathCmd, args, &os.ProcAttr{Files: []*os.File{os.Stdin, os.Stdout, os.Stderr}})

		if err != nil {
			if logger != nil {
				logger(fmt.Sprintf("#%s Command error: %s", this.id, err))
			}
		} /* else {
			this.processPointer = runPrcPointer
		}*/
		if this.output != nil && cmdOutput != nil && len(cmdOutput) > 0 {
			this.output(fmt.Sprintf("#%s OUTPUT %s\n", this.id, this.cmd))
			this.output(fmt.Sprintln(string(cmdOutput)))
			//fmt.Fprintf(*this.output, "#%s OUTPUT %s\n", this.id, this.cmd)
			//fmt.Fprintln(*this.output, string(cmdOutput))
		}
	} else {
		if logger != nil {
			logger(fmt.Sprintf("#%s Not found command %s", this.id, parts[0]))
			//logger.Printf("#%s Not found command %s", this.id, parts[0])
		}
	}
	if logger != nil {
		logger(fmt.Sprintf("#%s STOP", this.id))
		//logger.Printf("#%s STOP", this.id)
	}
}

func (this *process) iRunning() bool {
	return this.running
}

/**
 * Кол-во запущенных процессов
 */
func CountRunning() int {
	var total int = 0
	for i := 0; i < maxProcess; i++ {
		if processes[i] != nil {
			if processes[i].iRunning() {
				total += 1
			}
		}
	}
	return total
} // end func

/**
 * Подобрать пустой слот для процесса
 */
func getEmptySlotIndex() int {
	var emptyIndex int = -1
	var i int
	for i = 0; i < maxProcess; i++ {
		if processes[i] != nil {
			if !processes[i].iRunning() {
				emptyIndex = i
				break
			}
		} else {
			emptyIndex = i
			break
		}
	}
	return emptyIndex
}
