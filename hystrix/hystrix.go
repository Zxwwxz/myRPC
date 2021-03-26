package hystrix

import (
	"github.com/afex/hystrix-go/hystrix"
)

//熔断
const (
	defalut_timeout = 1000                 // 执行超时时间
	defalut_max_concurrent_requests = 2    // 最大并发数
	defalut_sleep_window = 1000            // 开启熔断后隔多久再次检测，单位毫秒
	defalut_error_percent_threshold = 50   // 超过阈值，开始计算错误率，达到错误率开启熔断
	defalut_request_volume_threshold = 1   // 开启熔断前调用次数先超过这个值
)

func InitHystrix(serviceName string,timeout,maxConcurrentRequests,requestVolumeThreshold,sleepWindow,errorPercentThreshold int)(error){
	if timeout == 0 {
		timeout = defalut_timeout
	}
	if maxConcurrentRequests == 0 {
		maxConcurrentRequests = defalut_max_concurrent_requests
	}
	if requestVolumeThreshold == 0 {
		requestVolumeThreshold = defalut_sleep_window
	}
	if sleepWindow == 0 {
		sleepWindow = defalut_error_percent_threshold
	}
	if errorPercentThreshold == 0 {
		errorPercentThreshold = defalut_request_volume_threshold
	}
	//全局配置
	hystrix.ConfigureCommand(serviceName, hystrix.CommandConfig{
		Timeout:timeout,
		MaxConcurrentRequests:maxConcurrentRequests,
		RequestVolumeThreshold:requestVolumeThreshold,
		SleepWindow:sleepWindow,
		ErrorPercentThreshold:errorPercentThreshold,
	})
	return nil
}

func Stop()  {

}