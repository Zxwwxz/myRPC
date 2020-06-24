package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"myRPC/demo/serialize/protobuf/person"
)

func main()  {
	phone1 := &person.Phone{
		Type : person.PhoneType_HOME,
		Number:"111",
	}
	phone2 := &person.Phone{
		Type : person.PhoneType_WORK,
		Number:"222",
	}
	personObj := &person.Person{
		Id : 1111,
		Name : "name",
		Phones: []*person.Phone{phone1,phone2},
	}
	personBytes,err := proto.Marshal(personObj)
	if err != nil {
		fmt.Println("Marshal err:",err)
	}
	err = ioutil.WriteFile("./proto.txt",personBytes,0777)
	if err != nil {
		fmt.Println("write err:",err)
		return
	}

	newPersonByte,err := ioutil.ReadFile("./proto.txt")
	if err != nil {
		fmt.Println("read err:",err)
		return
	}

	var newPerson person.Person
	err = proto.Unmarshal(newPersonByte,&newPerson)
	if err != nil {
		fmt.Println("Unmarshal err:",err)
		return
	}

	fmt.Println("newPerson:",newPerson)
}
