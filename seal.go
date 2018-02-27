package main

import (
	"flag"
	"log"
	"os"
	"seal/conf"
	"sync"
	"time"

	"github.com/calabashdad/utiltools"
)

const sealVersion = "seal: 1.0.0"

var (
	configFile  = flag.String("c", "./seal.yaml", "configure filename")
	showVersion = flag.Bool("v", false, "show version of seal")
)

var (
	gGuards sync.WaitGroup
)

func init() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)
	flag.Parse()
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
			time.Sleep(1 * time.Second)
		}
	}()

	if len(os.Args) < 2 {
		log.Println("Show usage : ./seal --help.")
		return
	}

	if *showVersion {
		log.Println(sealVersion)
		return
	}

	err := conf.GlobalConfInfo.Loads(*configFile)
	if err != nil {
		log.Println("conf loads failed.err=", err)
		return
	}

	log.Printf("load conf file success, conf=%+v\n", conf.GlobalConfInfo)

	gGuards.Add(1)
	if true {
		rtmpSrv := rtmpServer{}
		rtmpSrv.Start()
	}

	gGuards.Wait()
	log.Println("seal quit gracefully.")
}
