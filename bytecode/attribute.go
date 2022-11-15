package bytecode

import (
	"bytes"
	"encoding/binary"
)

type AttributeCommon struct {
	NameIndex uint16
	Length    uint32
	Info      []byte
}

func (a *AttributeCommon) Parse(data []byte, index int) int {
	binary.Read(bytes.NewBuffer(data[index:index+2]), binary.BigEndian, &a.NameIndex)
	binary.Read(bytes.NewBuffer(data[index+2:index+6]), binary.BigEndian, &a.Length)
	a.Info = data[index+6 : index+6+int(a.Length)]
	return 6 + int(a.Length)
}

func (a *AttributeCommon) GetName(constantPool []ConstantPoolInfo) string {
	result := ""
	result += constantPool[a.NameIndex].String(constantPool)
	return result
}

type AttributeInfo interface {
	Parse(nameIndex uint16, length uint32, data []byte)
	GetName(constantPool []ConstantPoolInfo) string
	String(constantPool []ConstantPoolInfo) string
}

type SourceFile struct {
	NameIndex       uint16
	Length          uint32
	SourceFileIndex uint16
}

func (s *SourceFile) Parse(nameIndex uint16, length uint32, data []byte) {
	s.NameIndex = nameIndex
	s.Length = length
	binary.Read(bytes.NewBuffer(data[0:length]), binary.BigEndian, &s.SourceFileIndex)
}

func (s *SourceFile) GetName(constantPool []ConstantPoolInfo) string {
	return constantPool[s.NameIndex].String(constantPool)
}

func (s *SourceFile) String(constantPool []ConstantPoolInfo) string {
	return constantPool[s.SourceFileIndex].String(constantPool)
}

type InnerClasses struct {
	NameIndex       uint16
	Length          uint32
	NumberOfClasses uint16
	Classes         []InnerClassInfo
}

type InnerClassInfo struct {
	InnerClassIndex       uint16
	OuterClassIndex       uint16
	InnerNameIndex        uint16
	InnerClassAccessFlags uint16
}

func (i *InnerClassInfo) Parse(data []byte, index int) int {
	binary.Read(bytes.NewBuffer(data[index:index+2]), binary.BigEndian, &i.InnerClassIndex)
	binary.Read(bytes.NewBuffer(data[index+2:index+4]), binary.BigEndian, &i.OuterClassIndex)
	binary.Read(bytes.NewBuffer(data[index+4:index+6]), binary.BigEndian, &i.InnerNameIndex)
	binary.Read(bytes.NewBuffer(data[index+6:index+8]), binary.BigEndian, &i.InnerClassAccessFlags)
	return index + 8
}

func (i *InnerClasses) Parse(nameIndex uint16, length uint32, data []byte) {
	i.NameIndex = nameIndex
	i.Length = length
	binary.Read(bytes.NewBuffer(data[0:2]), binary.BigEndian, &i.NumberOfClasses)
	i.Classes = make([]InnerClassInfo, i.NumberOfClasses)
	index := 2
	for n := 0; n < int(i.NumberOfClasses); n++ {
		info := &InnerClassInfo{}
		index = info.Parse(data, index)
		i.Classes[n] = *info
	}
}

func (i *InnerClasses) GetName(constantPool []ConstantPoolInfo) string {
	return constantPool[i.NameIndex].String(constantPool)
}

func (i *InnerClasses) String(constantPool []ConstantPoolInfo) string {
	result := ""
	for _, innerClass := range i.Classes {
		result += constantPool[innerClass.InnerClassIndex].String(constantPool) + "." + constantPool[innerClass.OuterClassIndex].String(constantPool)
	}
	return result
}
