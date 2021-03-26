package logger

import (
	"fmt"
	"myRPC/util"
	"os"
	"path/filepath"
	"time"
)

const (
	default_path = "../logs/"
	default_max_size = 50000000
)

//文件日志输出器
type FileOutputer struct {
	//文件句柄
	file       		*os.File
	//文件最大
	maxSize         int64
	//文件路径
	path            string
	//原始文件名
	originFileName  string
	//当前文件名
	curFileName     string
	//创建时间
	createTime   	string
	//上次分割的日期
	lastSplitDay  	int
}

//新建
func NewFileOutputer(params map[interface{}]interface{},originFileName string) (*FileOutputer, error) {
	path := params["path"].(string)
	if path == "" {
		path = default_path
	}
	maxSize := params["max_size"].(int)
	if maxSize == 0 {
		maxSize = default_max_size
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	log := &FileOutputer{
		maxSize:int64(maxSize),
		path:absPath,
		originFileName: originFileName,
		createTime:"",
		lastSplitDay:-1,
	}
	err = log.init()
	return log, err
}

//初始化日志
func (f *FileOutputer) init() (err error) {
	//先创建文件夹
	if !util.IsFileExist(f.path) {
		err := os.Mkdir(f.path, os.ModePerm)
		if err != nil {
			return err
		}
	}
	//获取文件名
	curFilename,createTime := f.getCurFilename("create")
	f.createTime = createTime
	f.curFileName = curFilename
	//创建日志文件
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
	createTime := fmt.Sprintf("%d-%d-%d %d-%d-%d",now.Year(),now.Month(),now.Day(),now.Hour(),now.Minute(),now.Second())
	if nameType == "create" {
		curFilename = fmt.Sprintf("[%s][%s].log", f.originFileName,
			createTime)
	}else if nameType == "close"{
		curFilename = fmt.Sprintf("[%s][%s][%s].log", f.originFileName,
			f.createTime,createTime)
	}
	return curFilename,createTime
}

//初始化日志文件
func (f *FileOutputer) initFile(filename string) (file *os.File, err error) {
	file, err = os.OpenFile(fmt.Sprintf("%s\\%s",f.path,filename), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		err = fmt.Errorf("open file %s failed, err:%v", filename, err)
		return
	}
	return
}

//检测是否分割
func (f *FileOutputer) checkSplitFile(curTime time.Time,file *os.File)(error) {
	day := curTime.Day()
	//已经不是同一天了，需要分割
	if day != f.lastSplitDay {
		//关闭原有的文件，并改名字
		err := f.Close()
		if err != nil {
			return err
		}
		//新建一个文件
		err = f.init()
		if err != nil {
			return err
		}
	}
	fileInfo,err := os.Stat(fmt.Sprintf("%s\\%s",f.path,f.curFileName))
	if err != nil {
		return err
	}
	//文件大小过大，也要分割
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

//写日志
func (f *FileOutputer) Write(data *LogData)(error) {
	//分割
	err := f.checkSplitFile(data.curTime,f.file)
	if err != nil {
		return err
	}
	//往文件里面写
	_,err = f.file.Write(data.Bytes())
	return err
}

//关闭原有的文件，进行重命名
func (f *FileOutputer) Close()(error) {
	curFilename,_ := f.getCurFilename("close")
	err := os.Rename(fmt.Sprintf("%s\\%s",f.path,f.curFileName),curFilename)
	if err != nil {
		return err
	}
	return f.file.Close()
}
