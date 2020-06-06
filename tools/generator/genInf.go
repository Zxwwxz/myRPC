package generator

import toolsBase "myRPC/tools/base"

type GeneratorInf interface {
	Run(opt *toolsBase.Option,meta *toolsBase.ServiceMetaData) error
	Name() string
}
