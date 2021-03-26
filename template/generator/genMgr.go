package generator

import (
	"fmt"
	toolsBase "myRPC/template/base"
)

var GeneratorMgrObj = &GeneratorMgr{}

type GeneratorMgr struct {
	GenList []GeneratorInf
}

//注册所有的生成器
func(g *GeneratorMgr)Register(){
	g.GenList = append(g.GenList,
		NewGeneratorMeta(),
		NewGeneratorDir(),
		NewGeneratorGrpc(),
		NewGeneratorMain(),
		NewGeneratorRouter(),
		NewGeneratorCtr(),
		NewGeneratorConfig(),
		NewGeneratorClient(),
		NewGeneratorTest(),
	)
}

//执行所有生成器
func(g *GeneratorMgr)Run(opt *toolsBase.Option,meta *toolsBase.ServiceMetaData){
	for _,gen := range g.GenList{
		err := gen.Run(opt,meta)
		if err != nil {
			fmt.Println("gen run err:",gen.Name(),err)
		}
	}
}
