package gzlog

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

func InitGZLogger(fileName string, rollingFileSize int, rollingNum int) {
	if !strings.Contains(fileName, string(os.PathSeparator)) {
		wd, err := os.Getwd()
		if err != nil {
			log.Println("get current directory failed! user system std instead")
		}
		fileName = wd + string(os.PathSeparator) + fileName
	}
	if _, err := os.Stat(fileName); !os.IsNotExist(err) {
		fmt.Printf("fileName already existed, so we remove the old one!")
		os.Remove(fileName)
	}
	f, err := os.Create(fileName)
	if err != nil {
		log.Printf("create log file failed!%s \n", err)
		log.Println("user system std instead")
		return
	}
	log.SetOutput(&GZLogger{FileName: fileName, RollingFileSize: rollingFileSize, RollingNum: rollingNum,
		f: f, mu: &sync.Mutex{}})

}

type GZLogger struct {
	FileName        string
	RollingFileSize int
	RollingNum      int
	mu              *sync.Mutex
	f               *os.File
}

func (gz *GZLogger) Write(p []byte) (n int, err error) {
	gz.mu.Lock()
	defer gz.mu.Unlock()
	fmt.Println("start to call Write method!")
	s, _ := gz.f.Stat()
	if s.Size() > int64(gz.RollingFileSize) {
		gz.f.Close()
		RecurseRenameFile(gz.FileName, gz.RollingNum, 0)
		f, err := os.Create(gz.FileName)
		if err != nil {
			return 0, err
		}
		gz.f = f
		return f.Write(p)
	}
	fmt.Println("Do write!")
	return gz.f.Write(p)
}

func RecurseRenameFile(fileName string, rollingNum int, num int) error {
	num++
	_, err := os.Stat(fileName + "." + strconv.Itoa(num))
	fmt.Printf("stat fileName:%s \n", fileName+"."+strconv.Itoa(num))
	if !os.IsNotExist(err) {
		fmt.Printf("exists!")
		if num != rollingNum {
			fmt.Printf("not exceed, we recurse fileName:%s, num %d \n", fileName, num)
			err = RecurseRenameFile(fileName, rollingNum, num)
			if err != nil {
				return err
			}
		} else {
			fmt.Printf("exceed!,just remove:%s \n", fileName+"."+strconv.Itoa(num))
			os.Remove(fileName + "." + strconv.Itoa(num))
		}
	}
	suffix := ""
	if num != 1 {
		suffix = "." + strconv.Itoa(num-1)
	}
	fmt.Printf("Now we are renaming file:%s to %s \n", fileName+suffix, fileName+"."+strconv.Itoa(num))
	return os.Rename(fileName+suffix, fileName+"."+strconv.Itoa(num))
}
