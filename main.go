package main

import (
	"class-file-parser/classfile"
	"flag"
	"fmt"
	"io"
	"os"
)

const MagicNumber = "CAFEBABE"

func main() {
	var classFileName string
	flag.StringVar(&classFileName, "file", "", "字节码文件")
	flag.Parse()

	classFile, err := os.Open(classFileName)
	if err != nil {
		fmt.Printf("open class file error %s\n", err.Error())
		os.Exit(0)
	}

	data, err := io.ReadAll(classFile)
	if err != nil {
		fmt.Printf("read class file error %s\n", err.Error())
		os.Exit(0)
	}

	fmt.Printf("%s: %dbytes\n", classFileName, len(data))

	var magicNumber = fmt.Sprintf("%X%X%X%X", data[0], data[1], data[2], data[3])
	if magicNumber != MagicNumber {
		fmt.Printf("This is not a class file. Expect magic number %s, but actual is %s\n", MagicNumber, magicNumber)
		os.Exit(0)
	}

	version := &classfile.Version{}
	version.Parse(data[4:8])
	fmt.Println(version)
}
