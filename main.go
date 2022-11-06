package main

import (
	"class-file-parser/bytecode"
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

	cf := &bytecode.ClassFile{}
	cf.Parser(data)
	fmt.Println(cf)

}
