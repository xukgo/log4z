package log4z

import (
	"io"
	"sync/atomic"
)

type BreakWriter struct {
	enable int32
	writer io.Writer
}

func InitBreakWriter(enable bool, writer io.Writer) BreakWriter {
	var ev int32 = 0
	if enable {
		ev = 1
	} else {
		ev = 0
	}
	return BreakWriter{
		enable: ev,
		writer: writer,
	}
}

func (thiz *BreakWriter) SetEnable(enable bool) {
	if enable {
		atomic.StoreInt32(&thiz.enable, 1)
	} else {
		atomic.StoreInt32(&thiz.enable, 0)
	}
}

func (thiz *BreakWriter) Write(p []byte) (n int, err error) {
	if atomic.LoadInt32(&thiz.enable) == 0 {
		//fmt.Println("write return")
		return len(p), nil
	}

	return thiz.writer.Write(p)
}
