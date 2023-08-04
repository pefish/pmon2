package main

import (
	"github.com/pefish/pmon2/app"
	"github.com/pefish/pmon2/app/conf"
	"github.com/pefish/pmon2/app/god"
	"log"
)

func main() {
	config := conf.GetDefaultConf()
	err := app.Instance(config)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("pmon2 daemon is running! \n")

	// start monitor process file
	god.NewMonitor()
}
