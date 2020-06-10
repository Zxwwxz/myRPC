package util

import (
	"os"
	"strings"
)

const (
	myRPC_ENV   = "myRPC_env"
	PRODUCT_ENV = "product"
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
