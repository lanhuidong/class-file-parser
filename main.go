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

	/*index := 10
	var accessFlag uint16
	binary.Read(bytes.NewBuffer(data[index:index+2]), binary.BigEndian, &accessFlag)
	index += 2
	flagDesc := ""
	if accessFlag&bytecode.ACC_PUBLIC != 0 {
		flagDesc += "public "
	}
	if accessFlag&bytecode.ACC_FINAL != 0 {
		flagDesc += "final "
	}
	if accessFlag&bytecode.ACC_SUPER != 0 {
		flagDesc += "super "
	}
	if accessFlag&bytecode.ACC_INTERFACE != 0 {
		flagDesc += "interface "
	}
	if accessFlag&bytecode.ACC_ABSTRACT != 0 {
		flagDesc += "abstract "
	}
	if accessFlag&bytecode.ACC_SYNTHETIC != 0 {
		flagDesc += "synthetic "
	}
	if accessFlag&bytecode.ACC_ANNOTATION != 0 {
		flagDesc += "@ "
	}
	if accessFlag&bytecode.ACC_ENUM != 0 {
		flagDesc += "enum "
	}
	if accessFlag&bytecode.ACC_MODULE != 0 {
		flagDesc += "module "
	}
	fmt.Printf("access flag is %s\n", flagDesc)

	var thisClass uint16
	binary.Read(bytes.NewBuffer(data[index:index+2]), binary.BigEndian, &thisClass)
	index += 2
	fmt.Printf("this_class is %d\n", thisClass)

	var superClass uint16
	binary.Read(bytes.NewBuffer(data[index:index+2]), binary.BigEndian, &superClass)
	index += 2
	fmt.Printf("super_class is %d\n", superClass)

	var interfaceCount uint16
	binary.Read(bytes.NewBuffer(data[index:index+2]), binary.BigEndian, &interfaceCount)
	index += 2
	fmt.Printf("this class implement is %d interfaces\n", interfaceCount)
	for i := 0; i < int(interfaceCount); i++ {
		var interfaceIndex uint16
		binary.Read(bytes.NewBuffer(data[index:index+2]), binary.BigEndian, &interfaceIndex)
		index += 2
		fmt.Printf("interface #%d is at %d\n", i+1, interfaceIndex)
	}*/

}
