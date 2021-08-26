package util

import "os"

func FileExists(path string) bool {
	stat, err := os.Stat(path)
	if err!=nil {
		return  false;
	}
	return stat.IsDir() == false

}

func PathExists(path string) bool {
	stat, err := os.Stat(path)
	if err!=nil {
		return  false
	}
	return stat.IsDir()
}

func CreatePath(path string) (err error) {
	err = os.Mkdir(path,os.ModeDir)
	return err
}