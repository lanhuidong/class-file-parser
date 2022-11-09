package bytecode

import (
	"bytes"
	"encoding/binary"
)

type AttributeInfo struct {
	NameIndex uint16
	Length    uint32
	Info      []byte
}

func (a *AttributeInfo) Parse(data []byte, index int) int {
	binary.Read(bytes.NewBuffer(data[index:index+2]), binary.BigEndian, &a.NameIndex)
	binary.Read(bytes.NewBuffer(data[index+2:index+6]), binary.BigEndian, &a.Length)
	a.Info = data[index+6 : index+6+int(a.Length)]
	return 6 + int(a.Length)
}

func (a *AttributeInfo) String(constantPool []ConstantPoolInfo) string {
	result := ""
	result += constantPool[a.NameIndex].String(constantPool)
	return result
}
