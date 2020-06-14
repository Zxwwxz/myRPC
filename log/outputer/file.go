package logOutputer

import (
	"fmt"
	logBase "myRPC/log/base"
	"os"
	"path/filepath"
	"time"
)

type FileOutputer struct {
	file       *os.File
	accessFile *os.File
	filename      string
	lastSplitDay  int
}

func NewFileOutputer(filename string) (OutputerInterface, error) {
	filename, err := filepath.Abs(filename)
	if err != nil {
		return nil, err
	}
	log := &FileOutputer{
		filename: filename,
	}
	err = log.init("all")
	return log, err
}

func (f *FileOutputer) getCurFilename() (curFilename, originFilename string) {
	now := time.Now()
	curFilename = fmt.Sprintf("%s.%04d%02d%02d", f.filename,
		now.Year(), now.Month(), now.Day())
	originFilename = f.filename
	return
}

func (f *FileOutputer) getCurAccessFilename() (curAccessFilename, originAccessFilename string) {
	now := time.Now()
	curAccessFilename = fmt.Sprintf("%s.access.%04d%02d%02d", f.filename,
		now.Year(), now.Month(), now.Day())
	originAccessFilename = fmt.Sprintf("%s.acccess", f.filename)
	return
}

func (f *FileOutputer) initFile(filename, originFilename string) (file *os.File, err error) {
	file, err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		err = fmt.Errorf("open faile %s failed, err:%v", filename, err)
		return
	}
	_ = os.Symlink(filename, originFilename)
	return
}

func (f *FileOutputer) init(tag string) (err error) {
	if tag == "other" || tag == "all" {
		curFilename, originFilename := f.getCurFilename()
		f.file, err = f.initFile(curFilename, originFilename)
		if err != nil {
			return
		}
	}
	if tag == "access" || tag == "all" {
		accessFilename, originAccessFilename := f.getCurAccessFilename()
		f.accessFile, err = f.initFile(accessFilename, originAccessFilename)
		if err != nil {
			return
		}
	}
	f.lastSplitDay = time.Now().Day()
	return
}

func (f *FileOutputer) checkSplitFile(curTime time.Time,file *os.File) {
	day := curTime.Day()
	if day != f.lastSplitDay {
		_ = f.init("all")
		return
	}
}

func (f *FileOutputer) Write(data *logBase.LogData) {
	file := f.file
	if data.Level == "access" {
		file = f.accessFile
	}
	f.checkSplitFile(data.CurTime,file)
	_,_ = file.Write(data.Bytes())
}

func (f *FileOutputer) Close() {
	_ = f.file.Close()
	_ = f.accessFile.Close()
}
