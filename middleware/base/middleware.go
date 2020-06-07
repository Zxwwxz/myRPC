package mwBase

import (
	"context"
)

type MiddleWareFunc func(context.Context,interface{})(interface{},error)

type MiddleWare func(MiddleWareFunc)(MiddleWareFunc)

func Chain(list ...MiddleWare)(MiddleWare){
	return func(next MiddleWareFunc)(MiddleWareFunc){
		for i:=len(list)-1;i>=0;i--{
			next = list[i](next)
		}
		return next
	}
}