package main

import (
	"bytes"
	"class-file-parser/classfile"
	"encoding/binary"
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

	index := 0
	var magicNumber = fmt.Sprintf("%X%X%X%X", data[0], data[1], data[2], data[3])
	if magicNumber != MagicNumber {
		fmt.Printf("This is not a class file. Expect magic number %s, but actual is %s\n", MagicNumber, magicNumber)
		os.Exit(0)
	}

	index += 4
	version := &classfile.Version{}
	version.Parse(data[index : index+4])
	index += 4
	fmt.Println(version)

	var constantCount uint16
	binary.Read(bytes.NewBuffer(data[index:index+2]), binary.BigEndian, &constantCount)
	index += 2
	fmt.Printf("constant count is %d\n", constantCount)

	var tag uint8
	var constants []interface{}
	for i := 0; i < int(constantCount)-1; i++ {
		binary.Read(bytes.NewBuffer(data[index:index+1]), binary.BigEndian, &tag)
		switch tag {
		case 1:
			utf8 := &classfile.Utf8{}
			utf8.Parse(data[index : index+3])
			index += 3
			utf8.Value = data[index : index+int(utf8.Length)]
			index += int(utf8.Length)
			fmt.Printf("const #%d, tag %d, %v\n", i+1, tag, utf8)
			constants = append(constants, utf8)
		case 3:
			integer := &classfile.Integer{}
			integer.Parse(data[index : index+3])
			index += 5
			fmt.Printf("const #%d, tag %d, %v\n", i+1, tag, integer)
			constants = append(constants, integer)
		case 5:
			long := &classfile.Long{}
			long.Parse(data[index : index+5])
			index += 9
			fmt.Printf("const #%d, tag %d, %v\n", i+1, tag, long)
			constants = append(constants, long)
			i++
		case 7:
			class := &classfile.Class{}
			class.Parse(data[index : index+3])
			index += 3
			fmt.Printf("const #%d, tag %d, %v\n", i+1, tag, class)
			constants = append(constants, class)
		case 9:
			fieldref := &classfile.Fieldref{}
			fieldref.Parse(data[index : index+5])
			index += 5
			fmt.Printf("const #%d, tag %d, %v\n", i+1, tag, fieldref)
			constants = append(constants, fieldref)
		case 10:
			methodRef := &classfile.Methodref{}
			methodRef.Parse(data[index : index+5])
			index += 5
			fmt.Printf("const #%d, tag %d, %v\n", i+1, tag, methodRef)
			constants = append(constants, methodRef)
		case 12:
			NameAndType := &classfile.NameAndType{}
			NameAndType.Parse(data[index : index+5])
			index += 5
			fmt.Printf("const #%d, tag %d, %v\n", i+1, tag, NameAndType)
			constants = append(constants, NameAndType)
		default:
			fmt.Printf("tag %d\n", tag)
		}
	}

	var accessFlag uint16
	binary.Read(bytes.NewBuffer(data[index:index+2]), binary.BigEndian, &accessFlag)
	index += 2
	flagDesc := ""
	if accessFlag&classfile.ACC_PUBLIC != 0 {
		flagDesc += "public "
	}
	if accessFlag&classfile.ACC_FINAL != 0 {
		flagDesc += "final "
	}
	if accessFlag&classfile.ACC_SUPER != 0 {
		flagDesc += "super "
	}
	if accessFlag&classfile.ACC_INTERFACE != 0 {
		flagDesc += "interface "
	}
	if accessFlag&classfile.ACC_ABSTRACT != 0 {
		flagDesc += "abstract "
	}
	if accessFlag&classfile.ACC_SYNTHETIC != 0 {
		flagDesc += "synthetic "
	}
	if accessFlag&classfile.ACC_ANNOTATION != 0 {
		flagDesc += "@ "
	}
	if accessFlag&classfile.ACC_ENUM != 0 {
		flagDesc += "enum "
	}
	if accessFlag&classfile.ACC_MODULE != 0 {
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
	}

}
