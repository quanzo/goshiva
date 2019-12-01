/**
*
*
*
 */
package writecontrol

import (
	"io"
	"os"

	//	"os/exec"
	//"fmt"
	"path"
	"path/filepath"
	"sync"
)

var arFileCache map[string]*io.WriteCloser

func init() {
	arFileCache = make(map[string]*io.WriteCloser)
} // end init

type fileSource struct {
	lock           *sync.Mutex
	filename       string
	fd             *os.File
	modeWriteClose bool // перед записью, открывать файл и затем закрывать
}

func (this *fileSource) open() bool {
	//var err error
	if this.fd != nil {
		return true
	}
	if fd, err := os.OpenFile(this.filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644); err == nil {
		this.fd = fd
		return true
	} else {
		this.fd = nil
		return false
	}
} // end open

func (this *fileSource) close() error {
	if this.modeWriteClose {
		return this.Close()
	}
	return nil
}

func (this *fileSource) Write(p []byte) (int, error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	defer this.close()
	if this.open() {
		//fmt.Println("Write")
		return this.fd.Write(p)
	}
	//fmt.Println("Not Write")
	return 0, nil
} // end Write

func (this *fileSource) Close() error {
	if this.fd != nil {
		return this.fd.Close()
	}
	return nil
} // end Close

func (this *fileSource) FileExist() bool {
	// проверим существование файла
	_, err := os.Stat(this.filename)

	if os.IsNotExist(err) {
		// file not exists
		return false
	}
	return err == nil
} // end FileExist

/**
*
*
 */
func ToFile(filename string) (*io.WriteCloser, error) {
	return ToFileWCM(filename, false)
}

/**
* Создать Writer из файла, или вернуть уже существующий из кеша
* Если modeWriteClose == true, то после записи файл будет закрыт, а перед записью - открыт
 */
func ToFileWCM(filename string, modeWriteClose bool) (*io.WriteCloser, error) {
	// проверим наличие файла в кеше
	if cachePointer, existKey := arFileCache[filename]; existKey {
		return cachePointer, nil
	}

	// коррекция и определение имени файла
	if !path.IsAbs(filename) {
		filename, _ = filepath.Abs(filename)
		// еще раз проверим в кеше по имени
		if cachePointer, existKey := arFileCache[filename]; existKey {
			return cachePointer, nil
		}
	}
	// создание экземпляра
	var fs fileSource
	fs.filename = filename
	fs.lock = new(sync.Mutex)
	fs.modeWriteClose = modeWriteClose
	var w io.WriteCloser = &fs
	arFileCache[filename] = &w
	return arFileCache[filename], nil
} // end ToFileWCM
