package logger

import (
	"fmt"
	"myRPC/util"
	"os"
	"path/filepath"
	"time"
)

const (
	default_path = "./logs/"
	default_max_size = 50000000
)

type FileOutputer struct {
	file       		*os.File
	maxSize         int64
	path            string
	originFileName  string
	curFileName     string
	createTime   	string
	lastSplitDay  	int
}

func NewFileOutputer(params map[interface{}]interface{},originFileName string) (*FileOutputer, error) {
	path := params["path"].(string)
	if path == "" {
		path = default_path
	}
	maxSize := params["max_size"].(int64)
	if maxSize == 0 {
		maxSize = default_max_size
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	log := &FileOutputer{
		maxSize:maxSize,
		path:absPath,
		originFileName: originFileName,
		createTime:"",
		lastSplitDay:-1,
	}
	err = log.init()
	return log, err
}

func (f *FileOutputer) init() (err error) {
	if !util.IsFileExist(f.path) {
		err := os.Mkdir(f.path, os.ModePerm)
		if err != nil {
			return err
		}
	}
	curFilename,createTime := f.getCurFilename("create")
	f.createTime = createTime
	f.curFileName = curFilename
	f.file, err = f.initFile(curFilename)
	if err != nil {
		return err
	}
	f.lastSplitDay = time.Now().Day()
	return nil
}

//获取文件名
func (f *FileOutputer) getCurFilename(nameType string) (curFilename string,nowTime string) {
	now := time.Now()
	createTime := now.Format("2006-01-02 15:04:05")
	if nameType == "create" {
		curFilename = fmt.Sprintf("[%s][%s].log", f.originFileName,
			createTime)
	}else if nameType == "close"{
		curFilename = fmt.Sprintf("[%s][%s][%s].log", f.originFileName,
			f.createTime,createTime)
	}
	return curFilename,createTime
}

func (f *FileOutputer) initFile(filename string) (file *os.File, err error) {
	file, err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		err = fmt.Errorf("open file %s failed, err:%v", filename, err)
		return
	}
	return
}

func (f *FileOutputer) checkSplitFile(curTime time.Time,file *os.File)(error) {
	day := curTime.Day()
	if day != f.lastSplitDay {
		err := f.Close()
		if err != nil {
			return err
		}
		err = f.init()
		if err != nil {
			return err
		}
	}
	fileInfo,err := os.Stat(fmt.Sprintf("%s%s",f.path,f.curFileName))
	if err != nil {
		return err
	}
	if fileInfo.Size() >= f.maxSize {
		err := f.Close()
		if err != nil {
			return err
		}
		err = f.init()
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *FileOutputer) Write(data *LogData)(error) {
	err := f.checkSplitFile(data.curTime,f.file)
	if err != nil {
		return err
	}
	_,err = f.file.Write(data.Bytes())
	return err
}

func (f *FileOutputer) Close()(error) {
	curFilename,_ := f.getCurFilename("close")
	err := os.Rename(fmt.Sprintf("%s%s",f.path,f.curFileName),curFilename)
	if err != nil {
		return err
	}
	return f.file.Close()
}
