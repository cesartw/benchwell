package main

import (
	"flag"
	"io/ioutil"
	"os"

	"bitbucket.org/goreorto/sqlhero"
	"bitbucket.org/goreorto/sqlhero/config"
	"bitbucket.org/goreorto/sqlhero/logger"
	"bitbucket.org/goreorto/sqlhero/sqlengine"
)

func main() {
	var (
		phrase  string
		logFile string
	)
	flag.StringVar(&phrase, "phrase", "", "En/Decrypting phrase for connection passwords")
	flag.StringVar(&logFile, "logfile", "", "Enable log")
	flag.Parse()

	conf, err := config.New(phrase)
	if err != nil {
		panic(err)
	}

	output := ioutil.Discard
	if logFile != "" {
		output, err = os.OpenFile(logFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
	}

	eng := sqlengine.New(conf, logger.NewLogger(output))
	defer eng.Dispose()

	app := sqlhero.New(conf, eng)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
