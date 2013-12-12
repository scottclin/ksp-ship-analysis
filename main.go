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

	testPart, err := readPartFile("/home/tox/.local/share/Steam/SteamApps/common/Kerbal Space Program/GameData/Squad/Parts/Aero/advancedCanard/part.cfg")
	if err != nil {
		panic(err)
	}
	
	fmt.Println(testPart)
}

func readPartFile(path string) (*Part, error){

	part := new(Part)

	partfile, err := os.Open(path)
	if err != nil { 
		return &Part{}, err
	}

	scanner := bufio.NewScanner(partfile)
	for scanner.Scan() {
		if( strings.Contains(scanner.Text(), "mass") ){
			splitstring := strings.Split(scanner.Text(), " ")
			part.mass, err = strconv.ParseFloat(splitstring[2], 64)
		}
	}

	return part, scanner.Err()
}
