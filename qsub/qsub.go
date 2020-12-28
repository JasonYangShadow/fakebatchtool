package main

import (
	"os"
	"path/filepath"
	"strings"

	. "github.com/JasonYangShadow/fakebatchtool/util"
	logrus "github.com/sirupsen/logrus"
)

var (
	logger = logrus.New()
)

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

	_, cok := os.LookupEnv("ContainerId")
	if !cok {
		logger.WithFields(logrus.Fields{
			"args": args,
		}).Fatal("could not read ContainerId from the process")
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

	shpath := args[len(args)-1]
	if strings.HasSuffix(shpath, ".sh") {
		logger.WithFields(logrus.Fields{
			"script": shpath,
		}).Debug("debug before readfile")
		shpath = containerroot + containercwd + "/" + shpath
		data, derr := ReadFile(shpath)
		if derr != nil {
			logger.WithFields(logrus.Fields{
				"err":    derr,
				"shpath": shpath,
			}).Fatal("read file from path encounters error")
		}

		//starts pass content to filters
		data, _ = ProcessShell(data, cwd)

		//write data to the same folder
		parentfolder := filepath.Dir(shpath)
		newname := "." + filepath.Base(shpath)
		newpath := parentfolder + "/" + newname
		err := WriteToFile(newpath, data)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"err":  err,
				"path": newpath,
			}).Fatal("could not write file to new path")
		}

		//change output path
		opath = containerroot + containercwd + "/" + opath

		//modify args
		for idx, line := range args {
			if line == "-o" {
				args[idx+1] = opath
				break
			}
		}
		args[len(args)-1] = newpath

		err = Command("qsub", args[1:]...)
		if err != nil {
			if err != nil {
				logger.WithFields(logrus.Fields{
					"err":  err,
					"args": args,
				}).Fatal("could not execute command")
			}
		}
	}

}
