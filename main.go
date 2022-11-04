package main

import (
	"bytes"
	"class-file-parser/bytecode"
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

	cf := bytecode.ClassFile{}
	cf.Parser(data)
	fmt.Println(cf)

	index := 4
	version := &bytecode.Version{}
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
			utf8 := &bytecode.Utf8{}
			utf8.Parse(data[index : index+3])
			index += 3
			utf8.Value = data[index : index+int(utf8.Length)]
			index += int(utf8.Length)
			fmt.Printf("const #%d, tag %d, %v\n", i+1, tag, utf8)
			constants = append(constants, utf8)
		case 3:
			integer := &bytecode.Integer{}
			integer.Parse(data[index : index+3])
			index += 5
			fmt.Printf("const #%d, tag %d, %v\n", i+1, tag, integer)
			constants = append(constants, integer)
		case 5:
			long := &bytecode.Long{}
			long.Parse(data[index : index+5])
			index += 9
			fmt.Printf("const #%d, tag %d, %v\n", i+1, tag, long)
			constants = append(constants, long)
			i++
		case 7:
			class := &bytecode.Class{}
			class.Parse(data[index : index+3])
			index += 3
			fmt.Printf("const #%d, tag %d, %v\n", i+1, tag, class)
			constants = append(constants, class)
		case 8:
			jstring := &bytecode.JString{}
			jstring.Parse(data[index : index+3])
			index += 3
			fmt.Printf("const #%d, tag %d, %v\n", i+1, tag, jstring)
			constants = append(constants, jstring)
		case 9:
			fieldref := &bytecode.Fieldref{}
			fieldref.Parse(data[index : index+5])
			index += 5
			fmt.Printf("const #%d, tag %d, %v\n", i+1, tag, fieldref)
			constants = append(constants, fieldref)
		case 10:
			methodRef := &bytecode.Methodref{}
			methodRef.Parse(data[index : index+5])
			index += 5
			fmt.Printf("const #%d, tag %d, %v\n", i+1, tag, methodRef)
			constants = append(constants, methodRef)
		case 12:
			NameAndType := &bytecode.NameAndType{}
			NameAndType.Parse(data[index : index+5])
			index += 5
			fmt.Printf("const #%d, tag %d, %v\n", i+1, tag, NameAndType)
			constants = append(constants, NameAndType)
		case 15:
			methodHandle := &bytecode.MethodHandle{}
			methodHandle.Parse(data[index : index+4])
			index += 4
			fmt.Printf("const #%d, tag %d, %v\n", i+1, tag, methodHandle)
			constants = append(constants, methodHandle)
		case 18:
			invokeDynamic := &bytecode.InvokeDynamic{}
			invokeDynamic.Parse(data[index : index+5])
			index += 5
			fmt.Printf("const #%d, tag %d, %v\n", i+1, tag, invokeDynamic)
			constants = append(constants, invokeDynamic)
		default:
			fmt.Printf("tag %d\n", tag)
			return
		}
	}

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
	}

}
