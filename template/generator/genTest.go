package generator

import (
    "myRPC/template/base"
    "os"
    "path"
    "text/template"
)

//测试服务生成器
type generatorTest struct {}

func NewGeneratorTest()(*generatorTest){
    return &generatorTest{}
}

func(g *generatorTest)Run(opt *toolsBase.Option,meta *toolsBase.ServiceMetaData) error{
    filename := path.Join(opt.OutputPath, "test/ServiceTest.go")
    file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
    if err != nil {
        return err
    }
    defer file.Close()
    t := template.New("Test")
    t, err = t.Parse(testTemplateFile)
    if err != nil {
        return err
    }
    err = t.Execute(file, meta)
    if err != nil {
        return err
    }
    return nil
}

func(g *generatorTest)Name() string{
    return "gen_test"
}
