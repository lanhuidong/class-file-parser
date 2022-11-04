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
	magicNumber := fmt.Sprintf("%X%X%X%X", data[0], data[1], data[2], data[3])
	if magicNumber != MagicNumber {
		fmt.Printf("Invalid class file. Expect magic number %s, but actual is %s\n", MagicNumber, magicNumber)
		return
	}
}

func (f *ClassFile) String() {
}

type FieldInfo struct {
}

type MethodInfo struct {
}

type AttributeInfo struct {
}

type Version struct {
	MajorVersion uint16
	MinorVersion uint16
}

func (v *Version) Parse(data []byte) {
	binary.Read(bytes.NewBuffer(data[0:2]), binary.BigEndian, &v.MinorVersion)
	binary.Read(bytes.NewBuffer(data[2:4]), binary.BigEndian, &v.MajorVersion)
}

func (v *Version) String() string {
	if v.MajorVersion == 45 {
		return fmt.Sprintf("JDK Version 1.0.2 or 1.1, %d.%d", v.MajorVersion, v.MinorVersion)
	} else if v.MajorVersion > 52 {
		if v.MajorVersion >= 56 && v.MinorVersion != 0 && v.MinorVersion != 65535 {
			return "Unknow JDK Version"
		}
		jdkVersion := v.MajorVersion - 44
		if jdkVersion == 8 || jdkVersion == 11 || jdkVersion == 17 {
			return fmt.Sprintf("JDK Version %d (LTS), %d.%d", jdkVersion, v.MajorVersion, v.MinorVersion)
		} else {
			return fmt.Sprintf("JDK Version %d, %d.%d", jdkVersion, v.MajorVersion, v.MinorVersion)
		}
	} else if v.MajorVersion > 45 && v.MinorVersion == 0 {
		jdkVersion := v.MajorVersion - 44
		return fmt.Sprintf("JDK Version 1.%d (LTS), %d.%d", jdkVersion, v.MajorVersion, v.MinorVersion)
	}
	return "Unknow JDK Version"
}

type Utf8 struct {
	Tag    uint8
	Length uint16
	Value  []byte
}

func (u *Utf8) Parse(data []byte) {
	binary.Read(bytes.NewBuffer(data[0:1]), binary.BigEndian, &u.Tag)
	binary.Read(bytes.NewBuffer(data[1:3]), binary.BigEndian, &u.Length)
}

func (u *Utf8) String() string {
	return fmt.Sprintf("utf8: %s", string(u.Value))
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
