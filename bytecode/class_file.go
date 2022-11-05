package bytecode

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const MagicNumber = "CAFEBABE"

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
	Interfaces        uint16
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

}

func (f *ClassFile) String() string {
	result := f.Version() + "\n"
	result += fmt.Sprintf("constant number: %d\n", f.ConstantPoolCount)
	for i, item := range f.ConstantPool {
		if item != nil {
			result += fmt.Sprintf("const #%d = %s\n", i, item)
		}
	}
	return result
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

type FieldInfo struct {
}

type MethodInfo struct {
}

type AttributeInfo struct {
}

type Fieldref struct {
	Tag              uint8
	ClassIndex       uint16
	NameAndTypeIndex uint16
}

func (f *Fieldref) Parse(data []byte) {
	binary.Read(bytes.NewBuffer(data[0:1]), binary.BigEndian, &f.Tag)
	binary.Read(bytes.NewBuffer(data[1:3]), binary.BigEndian, &f.ClassIndex)
	binary.Read(bytes.NewBuffer(data[3:5]), binary.BigEndian, &f.NameAndTypeIndex)
}

func (f *Fieldref) String() string {
	return fmt.Sprintf("class index: %d, name and type index: %d", f.ClassIndex, f.NameAndTypeIndex)
}

type Methodref struct {
	Tag              uint8
	ClassIndex       uint16
	NameAndTypeIndex uint16
}

func (m *Methodref) Parse(data []byte) {
	binary.Read(bytes.NewBuffer(data[0:1]), binary.BigEndian, &m.Tag)
	binary.Read(bytes.NewBuffer(data[1:3]), binary.BigEndian, &m.ClassIndex)
	binary.Read(bytes.NewBuffer(data[3:5]), binary.BigEndian, &m.NameAndTypeIndex)
}

func (m *Methodref) String() string {
	return fmt.Sprintf("class index: %d, name and type index: %d", m.ClassIndex, m.NameAndTypeIndex)
}

type Class struct {
	Tag   uint8
	Index uint16
}

func (c *Class) Parse(data []byte) {
	binary.Read(bytes.NewBuffer(data[0:1]), binary.BigEndian, &c.Tag)
	binary.Read(bytes.NewBuffer(data[1:3]), binary.BigEndian, &c.Index)
}

func (c *Class) String() string {
	return fmt.Sprintf("class name is at index: %d", c.Index)
}

type JString struct {
	Tag   uint8
	Index uint16
}

func (s *JString) Parse(data []byte) {
	binary.Read(bytes.NewBuffer(data[0:1]), binary.BigEndian, &s.Tag)
	binary.Read(bytes.NewBuffer(data[1:3]), binary.BigEndian, &s.Index)
}

func (s *JString) String() string {
	return fmt.Sprintf("string is at index: %d", s.Index)
}

type Integer struct {
	Tag   uint8
	Value int32
}

func (i *Integer) Parse(data []byte) {
	binary.Read(bytes.NewBuffer(data[0:1]), binary.BigEndian, &i.Tag)
	binary.Read(bytes.NewBuffer(data[1:5]), binary.BigEndian, &i.Value)
}

func (i *Integer) String() string {
	return fmt.Sprintf("integer value is: %d", i.Value)
}

type Long struct {
	Tag   uint8
	Value int64
}

func (l *Long) Parse(data []byte) {
	binary.Read(bytes.NewBuffer(data[0:1]), binary.BigEndian, &l.Tag)
	binary.Read(bytes.NewBuffer(data[1:9]), binary.BigEndian, &l.Value)
}

func (l *Long) String() string {
	return fmt.Sprintf("long value is: %d", l.Value)
}

type NameAndType struct {
	Tag       uint8
	NameIndex uint16
	DescIndex uint16
}

func (n *NameAndType) Parse(data []byte) {
	binary.Read(bytes.NewBuffer(data[0:1]), binary.BigEndian, &n.Tag)
	binary.Read(bytes.NewBuffer(data[1:3]), binary.BigEndian, &n.NameIndex)
	binary.Read(bytes.NewBuffer(data[3:5]), binary.BigEndian, &n.DescIndex)
}

func (n *NameAndType) String() string {
	return fmt.Sprintf("name is at index: %d, desc is at index %d", n.NameIndex, n.DescIndex)
}

type MethodHandle struct {
	Tag   uint8
	Kind  uint8
	Index uint16
}

func (i *MethodHandle) Parse(data []byte) {
	binary.Read(bytes.NewBuffer(data[0:1]), binary.BigEndian, &i.Tag)
	binary.Read(bytes.NewBuffer(data[1:2]), binary.BigEndian, &i.Kind)
	binary.Read(bytes.NewBuffer(data[2:4]), binary.BigEndian, &i.Index)
}

func (i *MethodHandle) String() string {
	return fmt.Sprintf("MethodHandle kind is at index: %d, index is %d", i.Kind, i.Index)
}

type InvokeDynamic struct {
	Tag                      uint8
	BootstrapMethodAttrIndex uint16
	NameAndTypeIndex         uint16
}

func (i *InvokeDynamic) Parse(data []byte) {
	binary.Read(bytes.NewBuffer(data[0:1]), binary.BigEndian, &i.Tag)
	binary.Read(bytes.NewBuffer(data[1:3]), binary.BigEndian, &i.BootstrapMethodAttrIndex)
	binary.Read(bytes.NewBuffer(data[3:5]), binary.BigEndian, &i.NameAndTypeIndex)
}

func (i *InvokeDynamic) String() string {
	return fmt.Sprintf("bootstrap method is at index: %d, name and type is at index %d", i.BootstrapMethodAttrIndex, i.NameAndTypeIndex)
}
