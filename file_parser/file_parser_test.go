package file_parser

import (
	"encoding/json"
	"fmt"
	"github.com/shirley981128/swift_file_parser/file_parser/mt940_parser"
	"github.com/shirley981128/swift_file_parser/file_parser/mt942_parser"
	"io/ioutil"
	"strings"
	"testing"
)

func TestParseFile(t *testing.T) {
	file, err := ioutil.ReadFile("mt942_parser/MT942_test.txt")
	if err != nil {
		fmt.Println("read file error")
	}
	origin := mt942_parser.ParseMT942File(file)
	jsonOrigin, err := json.Marshal(origin)
	if err != nil {
		fmt.Println("json marshal fail")
	}
	fmt.Println(string(jsonOrigin))
}

func TestMT940ParseField(t *testing.T) {
	const DIR = "./mt940_parser"
	dir, err := ioutil.ReadDir(DIR)
	if err != nil {
		fmt.Println("read dir error")
	}
	for _, fileInfo := range dir {
		//if strings.HasSuffix(fileInfo.Name(),"txt"){
		if strings.HasPrefix(fileInfo.Name(), "MT940") && strings.HasSuffix(fileInfo.Name(), ".txt") {
			file, err := ioutil.ReadFile(DIR + "/" + fileInfo.Name())
			if err != nil {
				fmt.Println("read file error")
			}
			mt940, ok := mt940_parser.ParseMT940Field(file)
			if !ok {
				fmt.Println("parse fail", mt940)
			}
			jsonMt940, err := json.Marshal(mt940)
			if err != nil {
				fmt.Println("json marshal fail")
			}
			fmt.Println(string(jsonMt940))
		}
	}
}

func TestMT942ParseField(t *testing.T) {
	const DIR = "./mt942_parser"
	dir, err := ioutil.ReadDir(DIR)
	if err != nil {
		fmt.Println("read dir error")
	}
	for _, fileInfo := range dir {
		if strings.HasPrefix(fileInfo.Name(), "MT942") && strings.HasSuffix(fileInfo.Name(), ".txt") {
			file, err := ioutil.ReadFile(DIR + "/" + fileInfo.Name())
			if err != nil {
				fmt.Println("read file error")
			}
			mt942, ok := mt942_parser.ParseMT942Field(file)
			if !ok {
				fmt.Println("parse fail", mt942)
			}
			jsonMt942, err := json.Marshal(mt942)
			if err != nil {
				fmt.Println("json marshal fail")
			}
			fmt.Println(string(jsonMt942))
		}
	}
}
