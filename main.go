package main

import (
	"fmt"
	"bufio"
	"os"
	"strings"
	"strconv"
	"path/filepath"
)

type Part struct {
	mass float64
	thrust int64
	name string
	drag float64
}

type Stage struct {
	mass float64
	thrust int64
	drag float64
	parts []Part
}

type Ship struct {
	mass float64
	weight float64
	initThrust int64
	twr float64
	stages []Stage	
}

var config map[string]string
var partsmap map[string]Part

func errcheck(err error){
	if err != nil {
		panic(err)
	}
}

func main() {

	err := readconfig("config.stfu")
	errcheck(err)

	partslocation, _ := config["partlocation"]
	
	partsmap = make(map[string]Part)

	err = filepath.Walk(partslocation, walkdirs)
	errcheck(err)

	fmt.Println(partsmap)
}


func walkdirs(path string , _ os.FileInfo, err error ) (error) {
	errcheck(err)

	err = filepath.Walk(path, readPartFile)
	errcheck(err)

	return nil
}

func readconfig(path string) (error) {
	
	configfile, err := os.Open(path)

	config = make(map[string]string)
	
	errcheck(err)
	
	defer configfile.Close()

	scanner := bufio.NewScanner(configfile)
	for scanner.Scan() {
		splitstring := strings.Split(scanner.Text(), "=")
		config[splitstring[0]] = splitstring[1]	
	}

	return  scanner.Err()
}

func readPartFile(path string, f os.FileInfo, err error) ( error){

	if !strings.Contains(path , "part.cfg"){
		return err
	}
	
	partfile, err := os.Open(path)
	errcheck(err)

	defer partfile.Close()

	part := new(Part)
	name := ""
	
	scanner := bufio.NewScanner(partfile)
	for scanner.Scan() {
//		fmt.Println(scanner.Text())
		splitstring := strings.Split(scanner.Text(), " ")
		switch  splitstring[0] { 
		default:
			continue
		case "name":
			name = splitstring[2]
			part.name = name
		case "mass":
			part.mass, err = strconv.ParseFloat(splitstring[2], 64)
		case "maximum_drag":
			part.drag, err = strconv.ParseFloat(splitstring[2], 64)
		}
//		if len(splitstring) > 3 {
		fmt.Println(splitstring)
		if strings.Contains(splitstring[0], "maxThrust"){
			part.thrust, err = strconv.ParseInt(splitstring[2], 10 ,64)
		}
//		}
	}

	partsmap[name] = *part
	return  scanner.Err()
}


