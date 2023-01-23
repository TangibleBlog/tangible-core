package utils

import (
	"fmt"
	"gopkg.in/errgo.v2/fmt/errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

//goland:noinspection GoUnusedExportedFunction
func CopyDir(srcPath string, destPath string) error {

	//Check the correctness of the directory
	if srcInfo, err := os.Stat(srcPath); err != nil {
		fmt.Println(err.Error())
		return err
	} else {
		if !srcInfo.IsDir() {
			e := errors.New("srcPath Not a proper directory!")
			fmt.Println(e.Error())
			return e
		}
	}
	if destInfo, err := os.Stat(destPath); err != nil {
		fmt.Println(err.Error())
		return err
	} else {
		if !destInfo.IsDir() {
			e := errors.New("destInfo Not a proper directory!")
			fmt.Println(e.Error())
			return e
		}
	}

	err := filepath.Walk(srcPath, func(path string, f os.FileInfo, err error) error {
		if f == nil {

			return err
		}
		if !f.IsDir() {
			path := strings.Replace(path, "\\", "/", -1)
			destNewPath := strings.Replace(path, srcPath, destPath, -1)
			log.Println("Copy files: " + path + " to: " + destNewPath)
			_, err := CopyFile(path, destNewPath)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {

		log.Printf(err.Error())
	}
	return err
}

// CopyFile Create directory and copy files
func CopyFile(src, dest string) (w int64, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer func(srcFile *os.File) {
		err := srcFile.Close()
		if err != nil {
		}
	}(srcFile)
	//Split path directory
	destSplitPathDirs := strings.Split(dest, "/")

	//The directory exists at the time of detection
	destSplitPath := ""
	for index, dir := range destSplitPathDirs {

		if index < len(destSplitPathDirs)-1 {

			destSplitPath = destSplitPath + dir + "/"
			b, _ := pathExists(destSplitPath)
			if b == false {

				log.Println("Create a directory:" + destSplitPath)
				//Create a directory
				err := os.Mkdir(destSplitPath, os.ModePerm)
				if err != nil {

					log.Println(err)
				}
			}
		}
	}
	dstFile, err := os.Create(dest)
	if err != nil {

		log.Println(err.Error())
		return
	}
	defer func(dstFile *os.File) {
		err := dstFile.Close()
		if err != nil {

		}
	}(dstFile)

	return io.Copy(dstFile, srcFile)
}

// Exist when detecting folder path
func pathExists(path string) (bool, error) {

	_, err := os.Stat(path)
	if err == nil {

		return true, nil
	}
	if os.IsNotExist(err) {

		return false, nil
	}
	return false, err
}
