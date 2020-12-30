package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	. "github.com/JasonYangShadow/fakebatchtool/util"
	logrus "github.com/sirupsen/logrus"
)

var (
	logger = logrus.New()
)

func init() {
	lvl, ok := os.LookupEnv("LOG_LEVEL")
	// LOG_LEVEL not set, let's default to debug
	if !ok {
		lvl = "info"
	}
	// parse string, this is built-in feature of logrus
	ll, err := logrus.ParseLevel(lvl)
	if err != nil {
		ll = logrus.InfoLevel
	}
	// set global log level
	logrus.SetLevel(ll)
}

func main() {
	args := os.Args
	var opath string
	for idx, line := range args {
		if line == "-o" {
			opath = args[idx+1]
			break
		}
	}

	epath, err := os.Executable()
	if err != nil {
		logger.WithFields(logrus.Fields{
			"args": args,
		}).Fatal("could not locate current executable path")
	}
	cwd, _ := filepath.Abs(filepath.Dir(epath))

	containerid, cok := os.LookupEnv("ContainerId")
	if !cok {
		logger.WithFields(logrus.Fields{
			"args": args,
		}).Fatal("could not read ContainerId from the process")
	}

	lpmxexe, cok := os.LookupEnv("LPMXEXE")
	if !cok {
		logger.WithFields(logrus.Fields{
			"args": args,
		}).Fatal("could not read LPMXEXE from the process")
	}

	containerroot, cok := os.LookupEnv("ContainerRoot")
	if !cok {
		logger.WithFields(logrus.Fields{
			"args": args,
		}).Fatal("could not read ContainerRoot from the process")
	}

	containercwd, cok := os.LookupEnv("ContainerCWD")
	if !cok {
		logger.WithFields(logrus.Fields{
			"args": args,
		}).Fatal("could not read ContainerCWD from the process")
	}

	//containercwd -> /root/xxx
	//lpmxexe -> /home/xxx/Linux-x86_64-lpmx (absolute path against the host)
	//cwd -> current executable path( /home/xx/bin/qsub)
	//containerroot -> /home/xx/.lpmxdata/ubuntu/18.04/workspaces/xxxxx/rw
	logger.WithFields(logrus.Fields{
		"container cwd":  containercwd,
		"lpmx exe":       lpmxexe,
		"container root": containerroot,
		"cwd":            cwd,
		"args":           args,
	}).Debug("debug the info when calling canuqsub")

	shpath := args[len(args)-1]
	data := []string{
		"#!/bin/bash",
		fmt.Sprintf("FAKECHROOT_USE_SYS_LIB=true %s resume %s \"cd %s && sh -c '%s'\" && echo \"done!\" && exit", lpmxexe, containerid, containercwd, shpath),
	}

	if strings.HasSuffix(shpath, ".sh") {
		logger.WithFields(logrus.Fields{
			"script": shpath,
		}).Debug("debug before readfile")

		opath = containerroot + containercwd + "/" + opath

		//modify output option
		for idx, line := range args {
			if line == "-o" {
				args[idx+1] = opath
				break
			}
		}

		rand_str := CreateShellScript()
		err := WriteToFile(rand_str, data)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"err":  err,
				"file": rand_str,
				"data": data,
			}).Fatal("could not write script to file")
		}

		//modify the submitted shell script
		rand_path := containerroot + containercwd + "/" + rand_str
		args[len(args)-1] = rand_path

		fmt.Println(rand_str, args, rand_path)
		err = Command("qsub", args[1:]...)
		if err != nil {
			if err != nil {
				logger.WithFields(logrus.Fields{
					"err":  err,
					"args": args,
				}).Fatal("could not execute command")
			}
		}

		//remove the submitted shell script
		err = os.Remove(rand_str)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"err":          err,
				"shell script": rand_str,
			}).Fatal("could not remove shell script")
		}
	}

}
