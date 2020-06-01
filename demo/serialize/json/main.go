package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

type Birthday time.Time

func(b *Birthday)MarshalJSON()(data []byte,err error){
	now := time.Time(*b)
	str := now.Format("2006-01-02")
	data, err = json.Marshal(str)
	return
}

func(b *Birthday)UnmarshalJSON(data []byte)(err error){
	var str string
	err = json.Unmarshal(data,&str)
	now,err := time.Parse("2006-01-02", str)
	*b = Birthday(now)
	return
}

type Person struct {
	Id int64         `json:"id,string"`
	Name string      `json:"name,omitempty"`
	Age int          `json:"-"`
	Day *Birthday    `json:"day"`
}

func main() {
	person := &Person{
		Id:1,
		Age:20,
	}

	var birthday Birthday
	birthday = Birthday(time.Now())
	person.Day = &birthday

	personByte,err := json.Marshal(person)
	if err != nil {
		fmt.Println("json err:",err)
		return
	}
	err = ioutil.WriteFile("./json.txt",personByte,0777)
	if err != nil {
		fmt.Println("write err:",err)
		return
	}

	newPersonByte,err := ioutil.ReadFile("./json.txt")
	if err != nil {
		fmt.Println("read err:",err)
		return
	}
	fmt.Println("json:",string(newPersonByte))
}
