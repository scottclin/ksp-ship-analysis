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
	name string
	mass float64
	weight float64
	initThrust int64
	twr float64
	stages []Stage	
}

var config map[string]string
var partsmap map[string]Part
var ships[] Ship

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
	ships = make([]Ship, 0)

	err = filepath.Walk(partslocation, walkdirsparts)
	errcheck(err)

	shiplocation, _ := config["shiplocation"]
	err = filepath.Walk(shiplocation, shipFiles)
	errcheck(err)

	fmt.Println(partsmap)
	for i := 0; i < len(ships); i++ {
		tempstages := ships[i].stages
		for j := 0; j < len(ships.stages); j++{
			fmt.Println(tempstages[j])
		}
	}
}


func walkdirsparts(path string , _ os.FileInfo, err error ) (error) {
	errcheck(err)

	err = filepath.Walk(path, readPartFile)
	errcheck(err)

	return nil
}

func walkdirsships(path string , _ os.FileInfo, err error ) (error) {
	errcheck(err)

	err = filepath.Walk(path, shipFiles)
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

func readPartFile(path string, _ os.FileInfo, err error) ( error){

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
		splitstring := strings.Split(scanner.Text(), " ")
		switch  splitstring[0] { 
		default:
			if strings.Contains(splitstring[0], "maxThrust"){
				part.thrust, err = strconv.ParseInt(splitstring[2], 10 ,64)
			}
		case "name":
			name = splitstring[2]
			part.name = name
		case "mass":
			part.mass, err = strconv.ParseFloat(splitstring[2], 64)
		case "maximum_drag":
			part.drag, err = strconv.ParseFloat(splitstring[2], 64)
		}
	}

	partsmap[name] = *part
	return  scanner.Err()
}


func shipFiles(path string, _ os.FileInfo, err error) error {

	if !strings.Contains(path ,".craft") || strings.Contains(path ,"Auto-Save"){
		return err
	}

	shipfile, err := os.Open(path)
	errcheck(err)

	defer shipfile.Close()

	ship := new(Ship)
	ship.stages = make([]Stage, 0)
	part := new(Part)
	currentstagenum := 0
	currentstage := new(Stage)
	currentstage.parts = make([]Part, 0)

	scanner := bufio.NewScanner(shipfile)
	for scanner.Scan() {
		splitstring := strings.Fields(scanner.Text())
		switch splitstring[0] {
		default:
			continue
		case "ship":
			for i := 2; i< len(splitstring);i++ {
				ship.name = ship.name + " " + splitstring[i]
			}
		case "part":
			if part.name != "" {
				currentstage.parts = append(currentstage.parts, *part)
			}
			part = new(Part)

			namesplit := strings.Split(splitstring[2], "_")
			part.name = namesplit[0]
			if part.name == "stackDecoupler" {
				ship.stages = append(ship.stages, *currentstage)
				currentstagenum++
				currentstage = new(Stage)
				currentstage.parts = make([]Part, 1)
			}
		}
	}

	ships = append(ships, *ship)

	return scanner.Err()
}
