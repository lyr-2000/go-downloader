package main

import (
	"fmt"
	"strconv"
	"testing"
)

func TestFname(t*testing.T) {
	fmt.Println(fName("/app/fexe"))
	fmt.Println(fName("/app/app.exe"))
	fmt.Println(fName("appexe"))
}
func TestFSize(t*testing.T) {
	fmt.Printf("%v",strconv.FormatInt(1024,2))
	//fmt.Println(fName("/app/fexe"))
	//fmt.Println(fName("/app/app.exe"))
	//fmt.Println(fName("appexe"))
}