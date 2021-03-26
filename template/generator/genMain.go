package generator

import (
	toolsBase "myRPC/template/base"
	"myRPC/util"
	"os"
	"path"
	"text/template"
)

//根据目标生成main文件
type generatorMain struct {}

func NewGeneratorMain()(*generatorMain){
	return &generatorMain{}
}

func(g *generatorMain)Run(opt *toolsBase.Option,meta *toolsBase.ServiceMetaData) error{
	filename := path.Join(opt.OutputPath, "main.go")
	if util.IsFileExist(filename) {
		return nil
	}
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	defer file.Close()

	t := template.New("main")
	t, err = t.Parse(mainTemplateFile)
	if err != nil {
		return err
	}
	err = t.Execute(file, meta)
	if err != nil {
		return err
	}
	return nil
}

func(g *generatorMain)Name() string{
	return "gen_main"
}
