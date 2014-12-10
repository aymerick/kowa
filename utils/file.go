package utils

import "path"

func FileBase(filePath string) string {
	fileName := path.Base(filePath)
	fileExt := path.Ext(filePath)

	return fileName[:len(fileName)-len(fileExt)]
}
