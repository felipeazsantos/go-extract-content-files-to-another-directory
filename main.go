package main

import (
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const TAG_EXTRACT_OPEN = "<svg"
const TAG_EXTRACT_CLOSE = "</svg>"
const NEW_FILE_SAVE_EXTENSION = ".svg"
const SOURCE_DIR = "/home/felipe/Projetos/CWS/Git/lib-react-components-cws/src/icons"
const DEST_DIR = "/home/felipe/Projetos/CWS/Git/micro-front-trading-desk-root-vue/src/assets/images/cms/icons"

func main() {
	files := readFiles()
	if extractFilesContentToAnotherDirectory(files) {
		fmt.Println("Content extraction from files to another directory was succefully processed")
	}
}

func readFiles() []fs.FileInfo {
	files, err := ioutil.ReadDir(SOURCE_DIR)
	if err != nil {
		log.Printf("Unable to read the given directory (%s): %s", SOURCE_DIR, err.Error())
		return nil
	}
	return files
}

func extractFilesContentToAnotherDirectory(files []fs.FileInfo) (result bool) {
	result = true

	for _, file := range files {
		if !file.IsDir() {
			filePath := filepath.Join(SOURCE_DIR, file.Name())
			content, err := ioutil.ReadFile(filePath)
			if err != nil {
				log.Printf("Error reading file %s", file.Name())
				continue
			}
			contentStr := string(content)
			tagOpen := strings.Index(contentStr, TAG_EXTRACT_OPEN)
			tagClose := strings.Index(contentStr, TAG_EXTRACT_CLOSE) + len(TAG_EXTRACT_CLOSE)

			if tagOpen != -1 && tagClose != -1 {
				contentStr = contentStr[tagOpen:tagClose]
			} else {
				log.Printf("Tag not found in file %s \n", file.Name())
				continue
			}

			// svg files from react lib
			if TAG_EXTRACT_OPEN == "<svg" {
				const FILL, WIDTH, HEIGHT, PROPS = "fill={color}", "width={size}", "height={size}", "{...props}"
				const FILL_REPLACE, WIDTH_REPLACE, HEIGHT_REPLACE = "fill='currentColor'", "width='73'", "height='44'"
				contentStr = strings.Replace(contentStr, FILL, FILL_REPLACE, -1)
				contentStr = strings.Replace(contentStr, WIDTH, WIDTH_REPLACE, -1)
				contentStr = strings.Replace(contentStr, HEIGHT, HEIGHT_REPLACE, -1)
				contentStr = strings.Replace(contentStr, PROPS, "", 1)
			}

			chanDestDir := make(chan bool)
			go writeToDestDir(contentStr, file.Name(), chanDestDir)
			createDirStatus := <-chanDestDir

			// if has an error to create destination directory, break the program
			if !createDirStatus {
				result = false
				break
			}
		}
	}

	return
}

func writeToDestDir(content, filename string, chanDestDir chan bool) {
	// make the destination directory
	err := os.MkdirAll(DEST_DIR, os.ModePerm)
	if err != nil {
		log.Printf("Error to create the destination directory (%s), error: %s \n", DEST_DIR, err.Error())
		chanDestDir <- false
		return
	}

	// create the destination file
	newFileName := strings.Split(filename, ".")[0] + NEW_FILE_SAVE_EXTENSION
	destinationFile, err := os.Create(DEST_DIR + "/" + newFileName)
	if err != nil {
		log.Printf("Error to create the destination file %s, error: %s \n", newFileName, err.Error())
		chanDestDir <- false
		return
	}
	defer destinationFile.Close()

	// write in the destination file
	_, err = io.WriteString(destinationFile, content)
	if err != nil {
		log.Printf("Error to write in the file %s, error: %s", newFileName, err.Error())
		chanDestDir <- false
		return
	}

	chanDestDir <- true
}
