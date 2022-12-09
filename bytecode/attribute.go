package bytecode

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func ParseAttribute(count int, data []byte, index int, constantPool []ConstantPoolInfo) (int, []AttributeInfo) {
	attrs := make([]AttributeInfo, count)
	for i := 0; i < count; i++ {
		len, attr := parse(data, index, constantPool)
		index += len
		attrs[i] = attr
	}
	return index, attrs
}

func parse(data []byte, index int, constantPool []ConstantPoolInfo) (int, AttributeInfo) {
	base := &AttributeBase{}
	binary.Read(bytes.NewBuffer(data[index:index+2]), binary.BigEndian, &base.NameIndex)
	binary.Read(bytes.NewBuffer(data[index+2:index+6]), binary.BigEndian, &base.Length)
	info := data[index+6 : index+6+int(base.Length)]
	base.Name = constantPool[base.NameIndex].String(constantPool)
	var item AttributeInfo
	switch base.Name {
	case "ConstantValue":
		item = &ConstantValue{}
		item.parse(base, info)
	case "InnerClasses":
		item = &InnerClasses{}
		item.parse(base, info)
	case "EnclosingMethod":
		item = &EnclosingMethod{}
		item.parse(base, info)
	case "SourceFile":
		item = &SourceFile{}
		item.parse(base, info)
	case "SourceDebugExtension":
		item = &SourceDebugExtension{}
		item.parse(base, info)
	case "BootstrapMethods":
		item = &BootstrapMethods{}
		item.parse(base, info)
	case "NestMembers":
		item = &NestMembers{}
		item.parse(base, info)
	default:
		fmt.Printf("attribue name is %s\n", base.Name)
	}
	return 6 + int(base.Length), item
}

type AttributeInfo interface {
	parse(base *AttributeBase, data []byte)
	GetName() string
	String(constantPool []ConstantPoolInfo) string
}

type AttributeBase struct {
	NameIndex uint16
	Name      string
	Length    uint32
}

func (a *AttributeBase) GetName() string {
	return a.Name
}

type ConstantValue struct {
	AttributeBase
	ConstantValueIndex uint16
}

func (c *ConstantValue) parse(base *AttributeBase, data []byte) {
	c.AttributeBase = *base
	binary.Read(bytes.NewBuffer(data), binary.BigEndian, &c.ConstantValueIndex)
}

func (c *ConstantValue) String(constantPool []ConstantPoolInfo) string {
	return constantPool[c.ConstantValueIndex].String(constantPool)
}

type InnerClasses struct {
	AttributeBase
	NumberOfClasses uint16
	Classes         []InnerClassInfo
}

type InnerClassInfo struct {
	InnerClassIndex       uint16
	OuterClassIndex       uint16
	InnerNameIndex        uint16
	InnerClassAccessFlags uint16
}

func (i *InnerClassInfo) parse(data []byte, index int) int {
	binary.Read(bytes.NewBuffer(data[index:index+2]), binary.BigEndian, &i.InnerClassIndex)
	binary.Read(bytes.NewBuffer(data[index+2:index+4]), binary.BigEndian, &i.OuterClassIndex)
	binary.Read(bytes.NewBuffer(data[index+4:index+6]), binary.BigEndian, &i.InnerNameIndex)
	binary.Read(bytes.NewBuffer(data[index+6:index+8]), binary.BigEndian, &i.InnerClassAccessFlags)
	return index + 8
}

func (i *InnerClasses) parse(base *AttributeBase, data []byte) {
	i.AttributeBase = *base
	binary.Read(bytes.NewBuffer(data[0:2]), binary.BigEndian, &i.NumberOfClasses)
	i.Classes = make([]InnerClassInfo, i.NumberOfClasses)
	index := 2
	for n := 0; n < int(i.NumberOfClasses); n++ {
		info := &InnerClassInfo{}
		index = info.parse(data, index)
		i.Classes[n] = *info
	}
}

func (i *InnerClasses) String(constantPool []ConstantPoolInfo) string {
	result := ""
	for _, innerClass := range i.Classes {
		result += constantPool[innerClass.InnerClassIndex].String(constantPool) + "." + constantPool[innerClass.OuterClassIndex].String(constantPool)
	}
	return result
}

type EnclosingMethod struct {
	AttributeBase
	ClassIndex  uint16
	MethodIndex uint16
}

func (e *EnclosingMethod) parse(base *AttributeBase, data []byte) {
	e.AttributeBase = *base
	binary.Read(bytes.NewBuffer(data[0:2]), binary.BigEndian, &e.ClassIndex)
	binary.Read(bytes.NewBuffer(data[2:4]), binary.BigEndian, &e.MethodIndex)
}

func (e *EnclosingMethod) String(constantPool []ConstantPoolInfo) string {
	return constantPool[e.ClassIndex].String(constantPool) + "." + constantPool[e.MethodIndex].String(constantPool)
}

type SourceFile struct {
	AttributeBase
	SourceFileIndex uint16
}

func (s *SourceFile) parse(base *AttributeBase, data []byte) {
	s.AttributeBase = *base
	binary.Read(bytes.NewBuffer(data), binary.BigEndian, &s.SourceFileIndex)
}

func (s *SourceFile) String(constantPool []ConstantPoolInfo) string {
	return constantPool[s.SourceFileIndex].String(constantPool)
}

type NestMembers struct {
	AttributeBase
	NumberOfClasses uint16
	Classes         []uint16
}

func (n *NestMembers) parse(base *AttributeBase, data []byte) {
	n.AttributeBase = *base
	binary.Read(bytes.NewBuffer(data[0:2]), binary.BigEndian, &n.NumberOfClasses)
	n.Classes = make([]uint16, n.NumberOfClasses)
	index := 2
	for i := 0; i < int(n.NumberOfClasses); i++ {
		var classIndex uint16
		binary.Read(bytes.NewBuffer(data[index:index+2]), binary.BigEndian, &classIndex)
		n.Classes[i] = classIndex
		index += 2
	}
}

func (n *NestMembers) String(constantPool []ConstantPoolInfo) string {
	result := ""
	for _, nestMember := range n.Classes {
		result += constantPool[nestMember].String(constantPool)
	}
	return result
}

type SourceDebugExtension struct {
	AttributeBase
	DebugExtension []uint8
}

func (s *SourceDebugExtension) parse(base *AttributeBase, data []byte) {
	s.AttributeBase = *base
	binary.Read(bytes.NewBuffer(data), binary.BigEndian, &s.DebugExtension)
}

func (n *SourceDebugExtension) String(constantPool []ConstantPoolInfo) string {
	return string(n.DebugExtension)
}

type BootStrapMethod struct {
	BootstrapMethodRef uint16
	ArgumentsNum       uint16
	Arguments          []uint16
}

type BootstrapMethods struct {
	AttributeBase
	Num     uint16
	Methods []BootStrapMethod
}

func (b *BootStrapMethod) parse(data []byte, index int) int {
	binary.Read(bytes.NewBuffer(data[index:index+2]), binary.BigEndian, &b.BootstrapMethodRef)
	binary.Read(bytes.NewBuffer(data[index+2:index+4]), binary.BigEndian, &b.ArgumentsNum)
	b.Arguments = make([]uint16, b.ArgumentsNum)
	index += 4
	for i := 0; i < int(b.ArgumentsNum); i++ {
		var argument uint16
		binary.Read(bytes.NewBuffer(data[index:index+2]), binary.BigEndian, &argument)
		b.Arguments[i] = argument
		index += 2
	}
	return index
}

func (b *BootstrapMethods) parse(base *AttributeBase, data []byte) {
	b.AttributeBase = *base
	binary.Read(bytes.NewBuffer(data[0:2]), binary.BigEndian, &b.Num)
	b.Methods = make([]BootStrapMethod, b.Num)
	index := 2
	for n := 0; n < int(b.Num); n++ {
		method := &BootStrapMethod{}
		index = method.parse(data, index)
		b.Methods[n] = *method
	}
}

func (b *BootstrapMethods) String(constantPool []ConstantPoolInfo) string {
	result := ""
	for _, method := range b.Methods {
		result += constantPool[method.BootstrapMethodRef].String(constantPool)
	}
	return result
}
