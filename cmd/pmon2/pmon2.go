package main

import (
	"github.com/pefish/pmon2/app"
	"github.com/pefish/pmon2/app/conf"
	"github.com/pefish/pmon2/client/cmd"
	"log"
)

func main() {
	conf := conf.GetDefaultConf()
	err := app.Instance(conf)
	if err != nil {
		log.Fatal(err)
	}

	err = cmd.Exec()
	if err != nil {
		log.Fatal(err)
	}
}
