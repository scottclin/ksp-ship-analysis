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
	drag float64
	weight float64
	initThrust int64
	twr float64
	stages []Stage	
}

var config map[string]string
var partsMap map[string]Part
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
	
	partsMap = make(map[string]Part)
	ships = make([]Ship, 0)

	err = filepath.Walk(partslocation, walkdirsparts)
	errcheck(err)

	shiplocation, _ := config["shiplocation"]
	err = filepath.Walk(shiplocation, shipFiles)
	errcheck(err)

	for i := 0; i < len(ships); i++ {
		matchmaking(ships[i])
		ships[i] = calcshipstats(ships[i])
		fmt.Println(ships[i])
		fmt.Println()
	}
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


func walkdirsparts(path string , _ os.FileInfo, err error ) (error) {
	errcheck(err)

	err = filepath.Walk(path, readPartFile)
	errcheck(err)

	return nil
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

	partsMap[name] = *part
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
			part.name = strings.Replace(namesplit[0], ".", "_",-1)

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

func matchmaking(shipToMatch Ship){
	for i := 0; i < len(shipToMatch.stages);i++{

		currentStage := shipToMatch.stages[i]

		for j := 0; j < len(currentStage.parts);j++{

			currentPart := currentStage.parts[j]
			partFound, haveWe := partsMap[currentPart.name]

			if !haveWe {
				fmt.Println("Part not found : ", currentPart.name)
				continue
			}

			currentPart.mass = partFound.mass
			currentPart.thrust = partFound.thrust
			currentPart.drag = partFound.drag
			
			currentStage.parts[j] = currentPart
		}
	}
}

func calcshipstats(shipToCalc Ship) (Ship) {

	shipMass := 0.0
	shipDrag := 0.0
	for i := 0; i< len(shipToCalc.stages);i++{
		currentStage := shipToCalc.stages[i]
		currentMass := 0.0
		var currentThrust int64
		currentThrust = 0
		currentDrag := 0.0

		for j := 0; j < len(currentStage.parts);j++{
			currentPart := currentStage.parts[j]
			currentDrag += currentPart.drag
			currentThrust += currentPart.thrust
			currentMass += currentPart.mass
		}

		currentStage.mass = currentMass
		shipMass += currentMass
		currentStage.thrust = currentThrust
		currentStage.drag = currentDrag
		shipDrag += currentDrag
		shipToCalc.stages[i] = currentStage
	} 

	shipToCalc.mass = shipMass
	shipToCalc.drag = shipDrag
	shipToCalc.weight = shipMass * 9.801
	shipToCalc.initThrust = shipToCalc.stages[len(shipToCalc.stages) - 1].thrust
	shipToCalc.twr = (float64(shipToCalc.initThrust) / shipToCalc.weight)

	return shipToCalc
}

