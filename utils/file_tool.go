package utils

import (
	"invar/status"
	"os"
	"path"
	"strings"

	"github.com/google/uuid"
)

func GenerateNewFilePathAndRemoveOld(fileName, savePath, oldFilePath string, fileSupportTypes []int) (string, int) {
	if !CheckMediaType(fileName, fileSupportTypes) {
		return "", status.UnsupportedMediaType
	}

	fileSavePath := FilePathGenerator(fileName, savePath)

	err := RemoveFile(oldFilePath)
	if err != nil {
		return "", status.Unkonwn
	}

	return fileSavePath, status.Success
}

func FilePathGenerator(file string, savePath string) string {
	fileExt := strings.ToLower(path.Ext(file))
	fileName := fileNameGenerator() + fileExt
	for {
		if !checkFileExist(savePath + fileName) {
			break
		}
		fileName = fileNameGenerator() + fileExt
	}
	return savePath + fileName
}

func RemoveFile(file string) error {
	if checkFileExist(file) {
		err := os.Remove(file)
		if err != nil {
			return err
		}
	}
	return nil
}

func checkFileExist(file string) bool {
	_, err := os.Stat(file)
	if err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func fileNameGenerator() string {
	uuid := uuid.New()
	return uuid.String()
}
