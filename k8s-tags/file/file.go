package file

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func CheckFileExists(filename string) bool {
	if _, err := os.Stat(filename); err == nil {
		fmt.Printf("File exists:%s\n", filename)
		return true
	}
	fmt.Printf("File does not exist:%s\n", filename)
	return false

}

func ReadFile(tags string, filename string) string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Panicf("failed reading data from file: %s", err)
	}
	fmt.Printf("\nLength: %d bytes", len(data))

	//fmt.Printf("\nData: %s", data)
	fmt.Printf("\nError: %v", err)

	contents := string(data)
	contents += tags


	return contents
}



func AppendFile(tags string, filename string) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}
	defer file.Close()

	len, err := file.WriteString(tags)
	if err != nil {
		log.Fatalf("failed writing to file: %s", err)
	}
	fmt.Printf("\nLength: %d bytes", len)
	fmt.Printf("\nFile Name: %s", file.Name())
}