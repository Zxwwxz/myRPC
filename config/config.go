package config

//服务的配置
import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

//可以通过参数修改的服务配置
type ServiceParams struct {
	ServiceType   int       `yaml:"service_type"`
	ServiceId     int       `yaml:"service_id"`
	ServiceVer    int       `yaml:"service_ver"`
	ServiceName   string    `yaml:"service_name"`
	ServicePort   int       `yaml:"service_port"`
	HttpPort      int       `yaml:"http_port"`
}

//可以通过http修改的服务配置
type ServiceHttpParams struct {
	PrometheusSwitchOn              bool       `yaml:"prometheus_switch_on"`
	LogSwitchOn                     bool       `yaml:"log_switch_on"`
	ClientLimitSwitchOn             bool       `yaml:"client_limit_switch_on"`
	ServerLimitSwitchOn             bool       `yaml:"server_limit_switch_on"`
	TraceSwitchOn                   bool       `yaml:"trace_switch_on"`
	HystrixSwitchOn                 bool       `yaml:"hystrix_switch_on"`
	LogLevel                        string     `yaml:"log_level"`
}

//服务的配置
type ServiceConf struct {
	Base        BaseConf       `yaml:"base"`
	Prometheus  PrometheusConf `yaml:"prometheus"`
	Registry    RegistryConf   `yaml:"registry"`
	Log         LogConf        `yaml:"log"`
	ClientLimit LimitConf      `yaml:"client_limit"`
	ServerLimit LimitConf      `yaml:"server_limit"`
	Trace       TraceConf      `yaml:"trace"`
	Balance     BalanceConf    `yaml:"balance"`
	Hystrix     HystrixConf    `yaml:"hystrix"`
	Http        HttpConf       `yaml:"http"`
	Other       interface{}    `yaml:"other"`

	RootDir    string `yaml:"-"`
	ConfigDir  string `yaml:"-"`
}

// 基础配置
type BaseConf struct {
	ServiceIDC    string    `yaml:"service_idc"`
	ServiceType   int       `yaml:"service_type"`
	ServiceId     int       `yaml:"service_id"`
	ServiceVer    int       `yaml:"service_ver"`
	ServiceName   string    `yaml:"service_name"`
	ServicePort   int       `yaml:"service_port"`
	ServiceWidget int       `yaml:"service_widget"`
	ServiceFuncs  []string  `yaml:"service_funcs"`
}

// 监控配置
type PrometheusConf struct {
	SwitchOn       bool     `yaml:"switch_on"`
	ClientHistogram string  `yaml:"client_histogram"`
	ServerHistogram string  `yaml:"server_histogram"`
}

// 注册配置
type RegistryConf struct {
	Type             string         `yaml:"type"`
	Params           interface{}    `yaml:"params"`
}

// 日志配置
type LogConf struct {
	SwitchOn      bool          `yaml:"switch_on"`
	Level         string        `yaml:"level"`
	ChanSize      int           `yaml:"chan_size"`
	Params        interface{}   `yaml:"params"`
}

// 限流配置
type LimitConf struct {
	SwitchOn   bool             `yaml:"switch_on"`
	Type       string           `yaml:"type"`
	Params     interface{}      `yaml:"params"`
}

// 追踪配置
type TraceConf struct {
	SwitchOn      bool       `yaml:"switch_on"`
	ReportAddr    string     `yaml:"addr"`
	SampleType    string     `yaml:"sample_type"`
	SampleRate    float64    `yaml:"sample_rate"`
}

// 负载配置
type BalanceConf struct {
	Type       string        `yaml:"type"`
}

// 熔断配置
type HystrixConf struct {
	SwitchOn                 bool       `yaml:"switch_on"`
	TimeOut                  int     	`yaml:"timeout"`
	MaxConcurrentRequests    int        `yaml:"max_concurrent_requests"`
	SleepWindow              int        `yaml:"sleep_window"`
	ErrorPercentThreshold    int        `yaml:"error_percent_threshold"`
	RequestVolumeThreshold   int        `yaml:"request_volume_threshold"`
}

type HttpConf struct {
	SwitchOn        bool             `yaml:"switch_on"`
	PprofSwitchOn   bool             `yaml:"pprof_switch_on"`
	HttpPort        int              `yaml:"http_port"`
}

// 初始化配置
func NewConfig(configDir string,serviceParams ServiceParams) (serviceConf *ServiceConf,err error) {
	serviceConf = &ServiceConf{}
	err = serviceConf.initDir(configDir)
	if err != nil {
		return
	}
	//读配置
	data, err := ioutil.ReadFile(serviceConf.ConfigDir)
	if err != nil {
		return
	}
	//解析配置
	err = yaml.Unmarshal(data, &serviceConf)
	if err != nil {
		return
	}
	serviceConf.paramsMergeConfig(serviceParams)
	return
}

//初始化配置路径
func (serviceConf *ServiceConf)initDir(configDir string) (err error) {
	//当前起服务路径
	exeFilePath, err := filepath.Abs(os.Args[0])
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
	//当前服务根路径
	serviceConf.RootDir = path.Join(strings.ToLower(exeFilePath[0:lastIdx]), ".")
	//当前服务配置路径
	if configDir != "" {
		serviceConf.ConfigDir = configDir
	}else {
		serviceConf.ConfigDir = path.Join(serviceConf.RootDir, "./config/", "/config.yaml")
	}
	return
}

func (serviceConf *ServiceConf)paramsMergeConfig(serviceParams ServiceParams) {
	if serviceParams.ServiceType != 0 {
		serviceConf.Base.ServiceType = serviceParams.ServiceType
	}
	if serviceParams.ServiceId != 0 {
		serviceConf.Base.ServiceId = serviceParams.ServiceId
	}
	if serviceParams.ServiceName != ""{
		serviceConf.Base.ServiceName = serviceParams.ServiceName
	}
	if serviceParams.ServiceVer != 0 {
		serviceConf.Base.ServiceVer = serviceParams.ServiceVer
	}
	if serviceParams.ServicePort != 0 {
		serviceConf.Base.ServicePort = serviceParams.ServicePort
	}
	if serviceParams.HttpPort != 0 {
		serviceConf.Http.HttpPort = serviceParams.HttpPort
	}
}

func (serviceConf *ServiceConf)HttpMergeConfig(serviceHttpParams ServiceHttpParams) {
	serviceConf.Prometheus.SwitchOn = serviceHttpParams.PrometheusSwitchOn
	serviceConf.Log.SwitchOn = serviceHttpParams.LogSwitchOn
	serviceConf.Hystrix.SwitchOn = serviceHttpParams.HystrixSwitchOn
	serviceConf.ClientLimit.SwitchOn = serviceHttpParams.ClientLimitSwitchOn
	serviceConf.ServerLimit.SwitchOn = serviceHttpParams.ServerLimitSwitchOn
	serviceConf.Trace.SwitchOn = serviceHttpParams.TraceSwitchOn
	serviceConf.Log.Level = serviceHttpParams.LogLevel
}

func (serviceConf *ServiceConf)GetConfigDir() string {
	return serviceConf.ConfigDir
}

func (serviceConf *ServiceConf)GetRootDir() string {
	return serviceConf.RootDir
}

func (serviceConf *ServiceConf)GetOtherConf() map[interface{}]interface{} {
	if serviceConf != nil {
		return nil
	}
	return serviceConf.Other.(map[interface{}]interface{})
}