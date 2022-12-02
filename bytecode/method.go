package bytecode

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const METHOD_ACC_PUBLIC = 0x0001
const METHOD_ACC_PRIVATE = 0x0002
const METHOD_ACC_PROTECTED = 0x0004
const METHOD_ACC_STATIC = 0x0008
const METHOD_ACC_FINAL = 0x0010
const METHOD_ACC_SYNCHRONIZED = 0x0020
const METHOD_ACC_BRIDGE = 0x0040
const METHOD_ACC_VARARGS = 0x0080
const METHOD_ACC_NATIVE = 0x0100
const METHOD_ACC_ABSTRACT = 0x0400
const METHOD_ACC_STRICT = 0x0800
const METHOD_ACC_SYNTHETIC = 0x1000

type MethodInfo struct {
	AccessFlags     uint16
	NameIndex       uint16
	DescriptorIndex uint16
	AttributesCount uint16
	Attributes      []AttributeInfo
}

func (m *MethodInfo) Parse(data []byte, index int, constantPool []ConstantPoolInfo) int {
	binary.Read(bytes.NewBuffer(data[index:index+2]), binary.BigEndian, &m.AccessFlags)
	binary.Read(bytes.NewBuffer(data[index+2:index+4]), binary.BigEndian, &m.NameIndex)
	binary.Read(bytes.NewBuffer(data[index+4:index+6]), binary.BigEndian, &m.DescriptorIndex)
	binary.Read(bytes.NewBuffer(data[index+6:index+8]), binary.BigEndian, &m.AttributesCount)

	index += 8
	indexInc := 8
	m.Attributes = make([]AttributeInfo, m.AttributesCount)
	for i := 0; i < int(m.AttributesCount); i++ {
		attr := &AttributeCommon{}
		indexInc = attr.Parse(data, index)
		index += indexInc
		var item AttributeInfo
		switch attr.GetName(constantPool) {
		case "SourceFile":
			item = &SourceFile{}
			item.Parse(attr.NameIndex, attr.Length, attr.Info)
		case "EnclosingMethod":
			item = &EnclosingMethod{}
			item.Parse(attr.NameIndex, attr.Length, attr.Info)
		default:
			fmt.Printf("method attr %s \n", attr.GetName(constantPool))
		}
		m.Attributes[i] = item
	}
	return index
}

func (m *MethodInfo) String(constantPool []ConstantPoolInfo) string {
	result := ""
	if METHOD_ACC_PRIVATE&m.AccessFlags != 0 {
		result += "private "
	} else if METHOD_ACC_PROTECTED&m.AccessFlags != 0 {
		result += "protected "
	} else if METHOD_ACC_PUBLIC&m.AccessFlags != 0 {
		result += "public "
	}
	if METHOD_ACC_STATIC&m.AccessFlags != 0 {
		result += "static "
	}
	if METHOD_ACC_FINAL&m.AccessFlags != 0 {
		result += "final "
	}
	if METHOD_ACC_SYNCHRONIZED&m.AccessFlags != 0 {
		result += "synchronized  "
	}
	if METHOD_ACC_NATIVE&m.AccessFlags != 0 {
		result += "native "
	}
	if METHOD_ACC_ABSTRACT&m.AccessFlags != 0 {
		result += "abstract "
	}
	result += constantPool[m.DescriptorIndex].String(constantPool) + " " + constantPool[m.NameIndex].String(constantPool)
	result += fmt.Sprintf("\n属性个数: %d\n", m.AttributesCount)
	for _, attr := range m.Attributes {
		if attr != nil {
			result += attr.GetName(constantPool) + ": " + attr.String(constantPool) + "\n"
		}
	}
	return result
}
