package util

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
)

//This read function reads from a Json file as a byte array and returns it.
//This function will be called for all the reading from json functions

func ReadByte(filename string) ([]byte, error) {
	jsonFile, err := os.Open(filename)
	// if we os.Open returns an error then handle it
	if err != nil {
		return nil, err
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}
	return byteValue, nil
}

//Writes arbitrary data as a JSON File.
// If the file does not exist, it will be created.
func WriteData(filename string, data interface{}) error {
	jsonFile, err := os.Open(filename)
	// if we os.Open returns an error then handle it
	if err != nil && strings.Contains(err.Error(), "no such file or directory") {
		jsonFile, err = os.Create(filename)
	}
	if err != nil {
		return err
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	//write to the corresponding file
	file, err := json.MarshalIndent(data, " ", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, file, 0644)
	if err != nil {
		return err
	}
	return nil
}

func CreateFile(path string) {
    // check if file exists
    var _, err = os.Stat(path)

    // create file if not exists
    if os.IsNotExist(err) {
        var file, err = os.Create(path)
        if err!=nil {
            return
        }
        defer file.Close()
    }
}
