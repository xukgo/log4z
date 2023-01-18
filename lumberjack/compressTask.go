package lumberjack

import (
	"log"
	"os"
	"sync"
	"time"
)

var once = new(sync.Once)
var delayArrayLocker = new(sync.Mutex)
var delayArray = make([]delayInfo, 0, 10)

type delayInfo struct {
	FileName string
	Delay    int
}

func addWillCompressFile(fileName string, delay int) {
	once.Do(func() {
		go loopTimerCompressFile()
	})

	delayArrayLocker.Lock()
	exist := false
	for idx := range delayArray {
		if delayArray[idx].FileName == fileName {
			exist = true
			break
		}
	}
	if !exist {
		delayArray = append(delayArray, delayInfo{FileName: fileName, Delay: delay})
	}
	//fmt.Println(delayArray)
	delayArrayLocker.Unlock()
}

func loopTimerCompressFile() {
	for {
		compressFileArr := make([]string, 0, 0)

		delayArrayLocker.Lock()
		idx := 0
		for _, info := range delayArray {
			fi, err := os.Stat(info.FileName)
			if err != nil {
				continue
			}
			if time.Since(fi.ModTime()).Seconds() >= float64(info.Delay) {
				compressFileArr = append(compressFileArr, info.FileName)
			} else {
				idx++
			}
		}
		delayArray = delayArray[:idx]
		delayArrayLocker.Unlock()

		for idx := range compressFileArr {
			fn := compressFileArr[idx]
			errCompress := compressLogFile(fn, fn+compressSuffix)
			if errCompress != nil {
				log.Printf("compress log file [%s] error:%v\r\n", fn, errCompress)
			}
		}
		time.Sleep(time.Second)
	}
}
