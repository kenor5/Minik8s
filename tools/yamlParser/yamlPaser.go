package yamlParser

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

/*
输入： v：需要parse的struct， filename： yaml文件路径
*/
func parseYaml(v interface{}, filename string) (bool, error) {
	b, _ := fileExists(filename)
	if !b {
		return false, errors.New("file not exist")
	}

	file, err := os.Open(filename)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("file open err")
		}
	}(file)
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	filesize := fileInfo.Size()

	buffer := make([]byte, filesize)
	_, err = file.Read(buffer)

	if err != nil {
		return false, errors.New("read file error")
	}

	err = yaml.Unmarshal(buffer, v)
	if err != nil {
		return false, errors.New("unmarshal yaml error")
	}
	return true, nil
}
