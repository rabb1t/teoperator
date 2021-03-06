package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/schollz/logger"
	"github.com/schollz/teoperator/src/download"
	"github.com/schollz/teoperator/src/ffmpeg"
	"github.com/schollz/teoperator/src/op1"
	"github.com/schollz/teoperator/src/server"
)

func main() {
	var flagSynth, flagOut, flagDuct, flagServerName string
	var flagDebug, flagServer, flagWorker bool
	var flagPort int
	flag.BoolVar(&flagDebug, "debug", false, "debug mode")
	flag.BoolVar(&flagServer, "serve", false, "make a server")
	flag.BoolVar(&flagWorker, "work", false, "start a download worker")
	flag.IntVar(&flagPort, "port", 8053, "port to use")
	flag.StringVar(&flagSynth, "synth", "", "build synth patch from file")
	flag.StringVar(&flagOut, "out", "", "name of new patch")
	flag.StringVar(&flagDuct, "duct", "", "name of duct")
	flag.StringVar(&flagServerName, "server", "http://localhost:8053", "name of external ip")
	flag.Parse()

	if flagDebug {
		log.SetLevel("debug")
	} else {
		log.SetLevel("info")
	}

	download.Duct = flagDuct
	download.ServerName = flagServerName

	if !ffmpeg.IsInstalled() {
		fmt.Println("ffmpeg not installed")
		fmt.Println("you can install it here: https://www.ffmpeg.org/download.html")
		os.Exit(1)
	}

	var err error
	if flagServer {
		err = server.Run(flagPort, flagServerName)
	} else if flagSynth != "" {
		_, fname := filepath.Split(flagSynth)
		if flagOut == "" {
			flagOut = strings.Split(fname, ".")[0] + ".op1.aif"
		}
		st := time.Now()
		sp := op1.NewSynthPatch()
		err = sp.SaveSample(flagSynth, flagOut, true)
		if err == nil {
			fmt.Printf("converted '%s' to op-1 synth patch '%s' in %s\n", fname, flagOut, time.Since(st))
		}
	} else if flagWorker {
		err = download.Work()
	} else {
		flag.PrintDefaults()
	}
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
