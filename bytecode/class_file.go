package bytecode

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const MagicNumber = "CAFEBABE"

const ACC_PUBLIC = 0x0001
const ACC_FINAL = 0x0010
const ACC_SUPER = 0x0020
const ACC_INTERFACE = 0x0200
const ACC_ABSTRACT = 0x0400
const ACC_SYNTHETIC = 0x1000
const ACC_ANNOTATION = 0x2000
const ACC_ENUM = 0x4000
const ACC_MODULE = 0x8000

type ClassFile struct {
	Magic             uint32
	MinorVersion      uint16
	MajorVersion      uint16
	ConstantPoolCount uint16
	ConstantPool      []ConstantPoolInfo
	AccessFlags       uint16
	ThisClass         uint16
	SuperClass        uint16
	InterfacesCount   uint16
	Interfaces        []uint16
	FieldsCount       uint16
	Fields            []FieldInfo
	MethodsCount      uint16
	Methods           []MethodInfo
	AttributesCount   uint16
	Attributes        []AttributeInfo
}

func (f *ClassFile) Parser(data []byte) {
	index := 0
	binary.Read(bytes.NewBuffer(data[index:index+4]), binary.BigEndian, &f.Magic)
	index += 4
	magicNumber := fmt.Sprintf("%X%X%X%X", data[0], data[1], data[2], data[3])
	if magicNumber != MagicNumber {
		fmt.Printf("Invalid class file. Expect magic number %s, but actual is %s\n", MagicNumber, magicNumber)
		return
	}

	binary.Read(bytes.NewBuffer(data[index:index+2]), binary.BigEndian, &f.MinorVersion)
	index += 2

	binary.Read(bytes.NewBuffer(data[index:index+2]), binary.BigEndian, &f.MajorVersion)
	index += 2

	binary.Read(bytes.NewBuffer(data[index:index+2]), binary.BigEndian, &f.ConstantPoolCount)
	index += 2

	f.ConstantPool = make([]ConstantPoolInfo, f.ConstantPoolCount)
	//常量池从1开始计数，long和double占2个位置
	for i := 1; i < int(f.ConstantPoolCount); i++ {
		var tag uint8
		var item ConstantPoolInfo
		addOne := false
		binary.Read(bytes.NewBuffer(data[index:index+1]), binary.BigEndian, &tag)
		switch tag {
		case 1:
			item = &ConstantUtf8{}
		case 3:
			item = &ConstantInteger{}
		case 4:
			item = &ConstantFloat{}
		case 5:
			item = &ConstantLong{}
			addOne = true
		case 6:
			item = &ConstantDouble{}
			addOne = true
		case 7:
			item = &ConstantClass{}
		case 8:
			item = &ConstantString{}
		case 9:
			item = &ConstantFieldref{}
		case 10:
			item = &ConstantMethodref{}
		case 11:
			item = &ConstantInterfaceMethodref{}
		case 12:
			item = &ConstantNameAndType{}
		case 15:
			item = &ConstantMethodHandle{}
		case 16:
			item = &ConstantMethodType{}
		case 17:
			item = &ConstantDynamic{}
		case 18:
			item = &ConstantInvokeDynamic{}
		case 19:
			item = &ConstantModule{}
		case 20:
			item = &ConstantPackage{}
		default:
			fmt.Printf("unkown constant type, tag: %d", tag)
		}
		index += item.Parse(data, index)
		f.ConstantPool[i] = item
		if addOne {
			i++
		}
	}

	binary.Read(bytes.NewBuffer(data[index:index+2]), binary.BigEndian, &f.AccessFlags)
	index += 2

	binary.Read(bytes.NewBuffer(data[index:index+2]), binary.BigEndian, &f.ThisClass)
	index += 2

	binary.Read(bytes.NewBuffer(data[index:index+2]), binary.BigEndian, &f.SuperClass)
	index += 2

	binary.Read(bytes.NewBuffer(data[index:index+2]), binary.BigEndian, &f.InterfacesCount)
	index += 2

	f.Interfaces = make([]uint16, f.InterfacesCount)
	for i := 0; i < int(f.InterfacesCount); i++ {
		binary.Read(bytes.NewBuffer(data[index:index+2]), binary.BigEndian, &f.Interfaces[i])
		index += 2
	}

	binary.Read(bytes.NewBuffer(data[index:index+2]), binary.BigEndian, &f.FieldsCount)
	index += 2

	f.Fields = make([]FieldInfo, 0)
	for i := 0; i < int(f.FieldsCount); i++ {
		field := &FieldInfo{}
		index = field.Parse(data, index)
		f.Fields = append(f.Fields, *field)
	}

	binary.Read(bytes.NewBuffer(data[index:index+2]), binary.BigEndian, &f.MethodsCount)
	index += 2

	f.Methods = make([]MethodInfo, 0)
	for i := 0; i < int(f.MethodsCount); i++ {
		method := &MethodInfo{}
		index = method.Parse(data, index)
		f.Methods = append(f.Methods, *method)
	}

	binary.Read(bytes.NewBuffer(data[index:index+2]), binary.BigEndian, &f.AttributesCount)
	index += 2

	for i := 0; i < int(f.AttributesCount); i++ {
		attr := &AttributeInfo{}
		index += attr.Parse(data, index)
		f.Attributes = append(f.Attributes, *attr)
	}

}

func (f *ClassFile) String() string {
	result := f.Version() + "\n"
	result += fmt.Sprintf("constant number: %d\n", f.ConstantPoolCount)
	for i, item := range f.ConstantPool {
		if item != nil {
			result += fmt.Sprintf("const #%d = %s	%s\n", i, item.TagName(), item.String(f.ConstantPool))
		}
	}

	if f.AccessFlags&ACC_PUBLIC != 0 {
		result += "public "
	}
	if f.AccessFlags&ACC_FINAL != 0 {
		result += "final "
	}
	if f.AccessFlags&ACC_SUPER != 0 {
		result += "super "
	}
	if f.AccessFlags&ACC_INTERFACE != 0 {
		result += "interface "
	}
	if f.AccessFlags&ACC_ABSTRACT != 0 {
		result += "abstract "
	}
	if f.AccessFlags&ACC_SYNTHETIC != 0 {
		result += "synthetic "
	}
	if f.AccessFlags&ACC_ANNOTATION != 0 {
		result += "@ "
	}
	if f.AccessFlags&ACC_ENUM != 0 {
		result += "enum "
	}
	if f.AccessFlags&ACC_MODULE != 0 {
		result += "module "
	}

	result += "\n"

	thisClassName := f.getClassName(f.ThisClass)
	result += thisClassName

	superClassName := f.getClassName(f.SuperClass)
	result += " extends " + superClassName

	for i := 0; i < int(f.InterfacesCount); i++ {
		interfaceName := f.getClassName(f.Interfaces[i])
		if i == 0 {
			result += " implements " + interfaceName
		} else {
			result += ", " + interfaceName
		}
	}
	result += "\n"
	result += fmt.Sprintf("字段个数: %d\n", f.FieldsCount)
	for _, field := range f.Fields {
		result += field.String(f.ConstantPool) + "\n"
	}

	result += "\n"
	result += fmt.Sprintf("方法个数: %d\n", f.MethodsCount)
	for _, method := range f.Methods {
		result += method.String(f.ConstantPool) + "\n"
	}

	result += "\n"
	result += fmt.Sprintf("属性个数: %d\n", f.AttributesCount)
	for _, attr := range f.Attributes {
		result += attr.String(f.ConstantPool) + "\n"
	}
	return result
}

func (f *ClassFile) getClassName(index uint16) (className string) {
	item := f.ConstantPool[index]
	constClazz, ok := item.(*ConstantClass)
	if ok {
		item = f.ConstantPool[constClazz.NameIndex]
		constUtf8, ok := item.(*ConstantUtf8)
		if ok {
			className = string(constUtf8.Value)
		}
	}
	return className
}

func (f *ClassFile) Version() string {
	if f.MajorVersion == 45 {
		return fmt.Sprintf("JDK Version 1.0.2 or 1.1, %d.%d", f.MajorVersion, f.MinorVersion)
	} else if f.MajorVersion > 52 {
		if f.MajorVersion >= 56 && f.MinorVersion != 0 && f.MinorVersion != 65535 {
			return "Unknow JDK Version"
		}
		jdkVersion := f.MajorVersion - 44
		if jdkVersion == 8 || jdkVersion == 11 || jdkVersion == 17 {
			return fmt.Sprintf("JDK Version %d (LTS), %d.%d", jdkVersion, f.MajorVersion, f.MinorVersion)
		} else {
			return fmt.Sprintf("JDK Version %d, %d.%d", jdkVersion, f.MajorVersion, f.MinorVersion)
		}
	} else if f.MajorVersion > 45 && f.MinorVersion == 0 {
		jdkVersion := f.MajorVersion - 44
		return fmt.Sprintf("JDK Version 1.%d (LTS), %d.%d", jdkVersion, f.MajorVersion, f.MinorVersion)
	}
	return "Unknow JDK Version"
}
