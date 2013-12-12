package main

import (
	"fmt"
	"bufio"
	"os"
	"strings"
	"strconv"
)

type Part struct {
	mass float64
	thrust int
	name string
	drag float64
}

type Stage struct {
	mass float64
	thrust int
	drag float64
	parts []Part
}

type Ship struct {
	mass float64
	weight float64
	initThrust int
	twr float64
	stages []Stage	
}

func main() {
	fmt.Println("Hello")
	
	partfile, err := readconfig("config.stfu")
	if err != nil {
		panic(err)
	}


	testPart, err := readPartFile(partfile)
	if err != nil {
		panic(err)
	}
	
	fmt.Println(testPart)
}

func readconfig(path string) (string, error) {
	configfile, err := os.Open(path)
	
	if err != nil {
		return "", err
	}
	
	defer configfile.Close()

	testfile := ""

	scanner := bufio.NewScanner(configfile)
	for scanner.Scan() {
		if( strings.Contains(scanner.Text(), "testfile") ){
			splitstring := strings.Split(scanner.Text(), "=")
			testfile = splitstring[1]	
		}
	}

	return testfile, scanner.Err()
}

func readPartFile(path string) (*Part, error){

	part := new(Part)

	partfile, err := os.Open(path)
	if err != nil { 
		return &Part{}, err
	}



	defer partfile.Close()

	scanner := bufio.NewScanner(partfile)
	for scanner.Scan() {
		if( strings.Contains(scanner.Text(), "mass") ){
			splitstring := strings.Split(scanner.Text(), " ")
			part.mass, err = strconv.ParseFloat(splitstring[2], 64)
		}
	}

	return part, scanner.Err()
}


