package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"myRPC/util"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	serviceConf = &ServiceConf{}
)

type ServiceConf struct {
	ServiceType int            `yaml:"service_type"`
	ServiceId   int            `yaml:"service_id"`
	ServiceVer   int           `yaml:"service_ver"`
	ServiceName string         `yaml:"service_name"`
	Port        int            `yaml:"port"`
	Prometheus  PrometheusConf `yaml:"prometheus"`
	Regiser     RegisterConf   `yaml:"register"`
	Log         LogConf        `yaml:"log"`
	Limit       LimitConf      `yaml:"limit"`
	Trace       TraceConf      `yaml:"trace"`

	RootDir    string `yaml:"-"`
	ConfigDir  string `yaml:"-"`
}

type PrometheusConf struct {
	SwitchOn bool `yaml:"switch_on"`
	Port     int  `yaml:"port"`
}

type RegisterConf struct {
	SwitchOn     bool          `yaml:"switch_on"`
	RegisterPath string        `yaml:"register_path"`
	Timeout      int64         `yaml:"timeout"`
	HeartBeat    int64         `yaml:"heart_beat"`
	RegisterName string        `yaml:"register_name"`
	RegisterAddr string        `yaml:"register_addr"`
}

type LogConf struct {
	Level      string `yaml:"level"`
	Dir        string `yaml:"path"`
	ChanSize   int    `yaml:"chan_size"`
	ConsoleLog bool   `yaml:"console_log"`
}

type LimitConf struct {
	SwitchOn bool     `yaml:"switch_on"`
	QPSLimit float64  `yaml:"qps"`
	AllWater int      `yaml:"all_water"`

}

type TraceConf struct {
	SwitchOn   bool    `yaml:"switch_on"`
	ReportAddr string  `yaml:"report_addr"`
	SampleType string  `yaml:"sample_type"`
	SampleRate float64 `yaml:"sample_rate"`
}

func initDir() (err error) {
	exeFilePath, err := filepath.Abs(os.Args[0])
	fmt.Println("lujing---------->:",exeFilePath)
	if err != nil {
		return
	}
	if runtime.GOOS == "windows" {
		exeFilePath = strings.Replace(exeFilePath, "\\", "/", -1)
	}
	lastIdx := strings.LastIndex(exeFilePath, "/")
	if lastIdx < 0 {
		err = fmt.Errorf("invalid exe path:%v", exeFilePath)
		return
	}
	serviceConf.RootDir = path.Join(strings.ToLower(exeFilePath[0:lastIdx]), "..")
	serviceConf.ConfigDir = path.Join(serviceConf.RootDir, "./config/", util.GetEnv(), "/config.yaml")
	return
}

func InitConfig() (err error) {
	err = initDir()
	if err != nil {
		return
	}

	data, err := ioutil.ReadFile(serviceConf.ConfigDir)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(data, &serviceConf)
	if err != nil {
		return
	}

	fmt.Printf("init conf succ, conf:%#v\n", serviceConf)
	return
}

func GetConfigDir() string {
	return serviceConf.ConfigDir
}

func GetRootDir() string {
	return serviceConf.RootDir
}

func GetServerPort() int {
	return serviceConf.Port
}

func GetConf() *ServiceConf {
	return serviceConf
}