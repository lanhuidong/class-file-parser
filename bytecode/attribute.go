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
