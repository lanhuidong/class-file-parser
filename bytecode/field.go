package bytecode

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const Field_ACC_PUBLIC = 0x0001
const Field_ACC_PRIVATE = 0x0002
const Field_ACC_PROTECTED = 0x0004
const Field_ACC_STATIC = 0x0008
const Field_ACC_FINAL = 0x0010
const Field_ACC_VOLATILE = 0x0040
const Field_ACC_TRANSIENT = 0x0080
const Field_ACC_SYNTHETIC = 0x1000
const Field_ACC_ENUM = 0x4000

type FieldInfo struct {
	AccessFlags     uint16
	NameIndex       uint16
	DescriptorIndex uint16
	AttributesCount uint16
	Attributes      []AttributeInfo
}

func (f *FieldInfo) Parse(data []byte, index int, constantPool []ConstantPoolInfo) int {
	binary.Read(bytes.NewBuffer(data[index:index+2]), binary.BigEndian, &f.AccessFlags)
	binary.Read(bytes.NewBuffer(data[index+2:index+4]), binary.BigEndian, &f.NameIndex)
	binary.Read(bytes.NewBuffer(data[index+4:index+6]), binary.BigEndian, &f.DescriptorIndex)
	binary.Read(bytes.NewBuffer(data[index+6:index+8]), binary.BigEndian, &f.AttributesCount)

	index += 8
	index, attrs := ParseAttribute(int(f.AttributesCount), data, index, constantPool)
	f.Attributes = attrs
	return index
}

func (f *FieldInfo) String(constantPool []ConstantPoolInfo) string {
	result := ""
	if Field_ACC_PRIVATE&f.AccessFlags != 0 {
		result += "private "
	} else if Field_ACC_PROTECTED&f.AccessFlags != 0 {
		result += "protected "
	} else if Field_ACC_PUBLIC&f.AccessFlags != 0 {
		result += "public "
	}
	if Field_ACC_STATIC&f.AccessFlags != 0 {
		result += "static "
	}
	if Field_ACC_FINAL&f.AccessFlags != 0 {
		result += "final "
	}
	if Field_ACC_VOLATILE&f.AccessFlags != 0 {
		result += "volatile "
	}
	if Field_ACC_TRANSIENT&f.AccessFlags != 0 {
		result += "transient "
	}
	result += constantPool[f.DescriptorIndex].String(constantPool) + " " + constantPool[f.NameIndex].String(constantPool)
	result += fmt.Sprintf("\n属性个数: %d\n", f.AttributesCount)
	for _, attr := range f.Attributes {
		if attr != nil {
			result += attr.GetName() + ": " + attr.String(constantPool) + "\n"
		}
	}
	return result
}
