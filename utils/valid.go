package utils

import (
	"path"
	"strings"
)

const (
	PNG = iota
	JPG
	JPEG
	GIF
	MP4
)

var supportMediaTypes = map[int]string{
	PNG:  ".png",
	JPG:  ".jpg",
	JPEG: ".jpeg",
	GIF:  ".gif",
	MP4:  ".mp4",
}

func CheckMediaType(fileName string, supportTypes []int) bool {
	fileExt := strings.ToLower(path.Ext(fileName))
	for _, v := range supportTypes {
		if supportMediaTypes[v] == fileExt {
			return true
		}
	}

	return false
}

func CheckMultiFileMediaType(filesName []string, supportTypes []int) bool {
	for _, v := range filesName {
		if !CheckMediaType(v, supportTypes) {
			return false
		}
	}

	return true
}

func CheckPasswordValid(pasword string) bool {
	if len(pasword) < 6 {
		return false
	}
	return true
}
