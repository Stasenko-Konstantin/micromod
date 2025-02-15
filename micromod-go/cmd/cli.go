package main

import (
	"fmt"
	"log"
	"micromod/internal/micromod"
	"os"
)

const (
	version = "" +
		"20180204 (c)2018 mumart@gmail.com\n" +
		"golang fork (2024) - stasenko.kost@yandex.ru"
)

func main() {
	var (
		args          = os.Args[1:]
		modFile       string
		interpolation = false
	)
	if len(args) == 2 && args[0] == "-int" {
		interpolation = true
		modFile = args[1]
	} else if len(args) == 1 {
		modFile = args[0]
	} else {
		log.Fatalln(
			"Micromod Go ProTracker Replay " + version + "\n" +
				"Usage: ./scripts/run.sh [-int] modfile",
		)
	}

	mm := micromod.New(
		modFile,
		interpolation, false,
	)
	fmt.Println(mm.ModuleInfo())

	clCh := make(chan struct{})
	go mm.Run(clCh)

	<-clCh
}
