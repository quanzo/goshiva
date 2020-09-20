package queue

/*
Класс реализации очереди команд в файле
Обеспечивает чтение построчное чтение команд из файла
Обеспечивает смену файла очереди
В строке из очереди %queue% будет заменено на имя файла очереди (полный путь)
*/

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

type Queue struct {
	fn   string // имя файла очереди комманд
	off  int64  // кол-во обработанных строк очереди
	lock *sync.Mutex
}

/**
* Создать очередь
*
 */
func CreateQueue(fn string) (*Queue, error) {
	q := new(Queue)
	q.lock = new(sync.Mutex)
	ifCorrectFile := q.ChangeQueue(fn)
	if ifCorrectFile {
		return q, nil
	}
	err := errors.New("Not correct queue file")
	return nil, err

	/*if _, err := os.Stat(fn); os.IsNotExist(err) {
		// file not exists
		return nil, err
	}
	return q, nil*/
}

func (this *Queue) Test() {
	fmt.Println("TEST")
}

/**
* Вернуть следующую команду из очереди
*
 */
func (this *Queue) Next() (string, error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	fstat, err := os.Stat(this.fn)
	if os.IsNotExist(err) {
		// file not exists
		return "", err
	}
	if fstat.Size() > this.off {
		for {
			f, err := os.Open(this.fn)
			if err == nil {
				defer f.Close()
				// читаем строку
				cmd, errRead := this.readString(f)
				if /*errRead == nil && */ cmd != "" {
					return this.replaceParams(cmd), errRead
				}
				return cmd, errRead
				break
			}
		}
	}
	return "", nil
} // end Next

/**
* Изменить файл очереди
 */
func (this *Queue) ChangeQueue(fn string) bool {
	this.lock.Lock()
	defer this.lock.Unlock()
	_, err := os.Stat(fn)
	if os.IsNotExist(err) {
		// file not exists
		return false
	}
	this.fn = fn
	this.off = 0
	return true
} // end ChangeQueue

/**
* Вернуть имя файла
 */
func (this *Queue) GetFile() string {
	return this.fn
}

/**
* Прочитать строку из очереди
* Комментарий начинается с # - они игнорируются
*
 */
func (this *Queue) readString(f *os.File) (string, error) {
	fstat, err := f.Stat()
	if err != nil {
		return "", err
	}
	filesize := fstat.Size()
	//fmt.Println(filesize, this.off)

	partReader := io.NewSectionReader(f, this.off, filesize /*-this.off*/)
	strReader := bufio.NewReader(partReader)
	str, err := strReader.ReadString('\n')
	this.off = this.off + (int64)(len(str))
	// дополнительная обработка строки
	if str != "" {
		str = strings.Trim(strings.Map(func(sym rune) rune {
			if sym == '\r' || sym == '\n' {
				return -1
			}
			return sym
		}, str), " ")
	}
	if strings.Index(str, "#") == 0 { // это комментарий
		return "", nil
	}
	return str, err
} // end readString

/**
* Возможные замены в команде
*
 */
func (this *Queue) replaceParams(s string) string {
	r := strings.NewReplacer("%queue%", this.fn)
	return r.Replace(s)
}
