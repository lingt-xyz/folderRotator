package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

func main() {
	root := flag.String("input", "", "root folder")
	flag.Parse()
	rotate(*root)
}

// {LIBRARY}-{VERSION}/{ARCHITECTURE}-{COMPILER}-{OPTIMIZATION}{OBFUSCATION}/
// {ARCHITECTURE}-{COMPILER}-{LIBRARY}/{VERSION}/-{OPTIMIZATION}/
func rotate(root string) {
	log.Printf("Scanning the directory `%v`; it should have the strucutre `%v`", root, path.Join(root, "{ARCHITECTURE}-{COMPILER}-{LIBRARY}/{VERSION}/-{OPTIMIZATION}/"))
	libDirs, err := ioutil.ReadDir(root)
	if err != nil {
		log.Fatalf("Cannot open the dataset directory `%v`", root)
	}
	for _, libDir := range libDirs {
		if !libDir.IsDir() {
			continue
		}
		libDirName := libDir.Name()
		libDirPath := path.Join(root, libDirName)
		libInfo := strings.Split(libDirName, "-")
		if len(libInfo) != 3 {
			log.Printf("Input library folder `%v` did not have a format of {ARCHITECTURE}-{COMPILER}-{LIBRARY}", libDirName)
			continue
		}
		architecture, compiler, libName := libInfo[0], libInfo[1], libInfo[2]
		if versionDirs, err := ioutil.ReadDir(libDirPath); err != nil {
			log.Printf("Cannot open the library directory %v", libDirName)
		} else {
			for _, versionDir := range versionDirs {
				if !versionDir.IsDir() {
					continue
				}
				versionDirName := versionDir.Name()
				libVersion := versionDirName
				versionDirPath := path.Join(libDirPath, versionDirName)
				if optimizationDirs, err := ioutil.ReadDir(versionDirPath); err != nil {
					log.Printf("Cannot open the library version directory %v", libDirName)
				} else {
					for _, optimizationDir := range optimizationDirs {
						if !optimizationDir.IsDir() {
							continue
						}
						optimizationDirName := optimizationDir.Name()
						optimizationDirPath := path.Join(versionDirPath, optimizationDirName)
						optimization := optimizationDirName[1:]
						// rename folder to "{LIBRARY}-{VERSION}/{ARCHITECTURE}-{COMPILER}-{OPTIMIZATION}"
						desFolder := path.Join(root, libName+"-"+libVersion)
						os.MkdirAll(desFolder, os.ModePerm)
						err := os.Rename(optimizationDirPath, path.Join(desFolder, architecture+"-"+compiler+"-"+optimization))
						if err != nil{
							log.Printf("%v", err)
						}
					}
				}

			}
		}
		os.RemoveAll(libDirPath)
	}
}
