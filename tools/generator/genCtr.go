package generator

import (
	toolsBase "myRPC/tools/base"
	"myRPC/util"
	"os"
	"path"
	"text/template"
)

type generatorCtr struct {}

func NewGeneratorCtr()(*generatorCtr){
	return &generatorCtr{}
}

func(g *generatorCtr)Run(opt *toolsBase.Option,meta *toolsBase.ServiceMetaData) error{
	filename := path.Join(opt.OutputPath, "controller/controller.go")
	if util.IsFileExist(filename) {
		return nil
	}
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	t := template.New("Ctr")
	t, err = t.Parse(ctrTemplateFile)
	if err != nil {
		return err
	}
	err = t.Execute(file, meta)
	if err != nil {
		return err
	}
	return nil
}

func(g *generatorCtr)Name() string{
	return "gen_Ctr"
}
