package hystrix

import "github.com/afex/hystrix-go/hystrix"

func InitHystrix(serviceName string,timeout,maxConcurrentRequests,requestVolumeThreshold,sleepWindow,errorPercentThreshold int)(error){
	hystrix.ConfigureCommand(serviceName, hystrix.CommandConfig{
		Timeout:timeout,
		MaxConcurrentRequests:maxConcurrentRequests,
		RequestVolumeThreshold:requestVolumeThreshold,
		SleepWindow:sleepWindow,
		ErrorPercentThreshold:errorPercentThreshold,
	})
	return nil
}
