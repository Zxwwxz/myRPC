package hystrix

import (
	"github.com/afex/hystrix-go/hystrix"
)

const (
	defalut_timeout = 1000
	defalut_max_concurrent_requests = 2
	defalut_sleep_window = 1000
	defalut_error_percent_threshold = 50
	defalut_request_volume_threshold = 1
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