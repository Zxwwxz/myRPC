package client

import (
	mwBase "myRPC/middleware/base"
)

func BuildClientMiddleware(handle mwBase.MiddleWareFunc,frontMiddles,backMiddles []mwBase.MiddleWare) mwBase.MiddleWareFunc {
	var middles []mwBase.MiddleWare
	middles = append(middles,frontMiddles...)
	middles = append(middles,backMiddles...)
	m := mwBase.Chain(middles...)
	return m(handle)
}
