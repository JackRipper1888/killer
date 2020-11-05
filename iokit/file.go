package iokit

import(
	"os"
	"io/ioutil"
)

func ReadAll(confDir string) ([]byte, error) {
	file, err := os.OpenFile(confDir, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
		return nil, err
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
		return nil, err
	}
	return data, nil
}

func Write(confDir string,targer []byte) error {
	file, err := os.Create(confDir)
	if err != nil {
		panic(err)
		return err
	}

	defer file.Close()
	_, err = file.Write(targer)
	if err != nil {
		panic(err)
		return err
	}
	return err
}