package util

import (
	"os"
	"strings"
)

const (
	//环境变量key
	myRPC_ENV   = "myRPC_env"
	//环境变量正式服value
	PRODUCT_ENV = "product"
	//环境变量测试服value
	TEST_ENV    = "test"
)

var (
	cur_env string = TEST_ENV
)

func init() {
	cur_env = strings.ToLower(os.Getenv(myRPC_ENV))
	cur_env = strings.TrimSpace(cur_env)

	if len(cur_env) == 0 {
		cur_env = TEST_ENV
	}
}

func IsProduct() bool {
	return cur_env == PRODUCT_ENV
}

func IsTest() bool {
	return cur_env == TEST_ENV
}

func GetEnv() string {
	return cur_env
}
