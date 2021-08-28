package file

import (
	"errors"
	"fmt"
	"io"
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
	fmt.Printf("\nLength: %d bytes, file:%s\n", len(data), filename)

	//fmt.Printf("\nData: %s", data)
	//fmt.Printf("\nError: %v", err)

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

func PrintFilesInDirectory()  {
	files, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fmt.Println(f.Name())
	}
}

func GetCurrentDirectory()  {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Current DIR:" +path)
}

func CreateACopyOfFile(srcFile string, outFile string)  {
	sourceFile, err := os.Open(srcFile)
	if err != nil {
		log.Fatal(err)
	}
	defer sourceFile.Close()

	// Create new file
	newFile, err := os.Create(outFile)
	if err != nil {
		log.Fatal(err)
	}
	defer newFile.Close()

	bytesCopied, err := io.Copy(newFile, sourceFile)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Copied %d bytes.", bytesCopied)
}

func WriteFile(fileName string, contents string)  {
	f, err := os.Create(fileName)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err2 := f.WriteString(contents)

	if err2 != nil {
		log.Fatal(err2)
	}

	fmt.Printf("Write file to %s done...\n", fileName)
}

func CreateDirectoryIfNotExists(path string)  {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}
}