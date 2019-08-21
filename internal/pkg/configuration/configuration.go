package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"reflect"
)

type Object map[string]interface{}

func main() {
	data, err := read()

	if err != nil {
		log.Panic(1, err.Error())
	}

	var js Object

	json.Unmarshal(data, &js)

	name, err := js.Get("name")

	if err != nil {
		log.Panic(2, err.Error())
	}

	nameAsserted, ok := name.(string)

	if !ok {
		log.Panic("Assertion failed")
	}

	log.Println(nameAsserted)
}

func (o Object) Get(key string) (interface{}, error) {
	val := reflect.ValueOf(o)

	if val.Kind() != reflect.Map {
		return nil, errors.New("Not a map")
	}

	m := map[string]interface{}(o)

	return m[key], nil
}

func read() ([]byte, error) {
	file, err := os.Open("test.json")

	if err != nil {
		return nil, errors.New(err.Error())
	}

	data, err := ioutil.ReadAll(file)

	if err != nil {
		return nil, errors.New(err.Error())
	}

	return data, nil
}
