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

type EnclosingMethod struct {
	NameIndex   uint16
	Length      uint32
	ClassIndex  uint16
	MethodIndex uint16
}

func (e *EnclosingMethod) Parse(nameIndex uint16, length uint32, data []byte) {
	e.NameIndex = nameIndex
	e.Length = length
	binary.Read(bytes.NewBuffer(data[0:2]), binary.BigEndian, &e.ClassIndex)
	binary.Read(bytes.NewBuffer(data[2:4]), binary.BigEndian, &e.MethodIndex)
}

func (e *EnclosingMethod) GetName(constantPool []ConstantPoolInfo) string {
	return constantPool[e.NameIndex].String(constantPool)
}

func (e *EnclosingMethod) String(constantPool []ConstantPoolInfo) string {
	return constantPool[e.ClassIndex].String(constantPool) + "." + constantPool[e.MethodIndex].String(constantPool)
}

type NestMembers struct {
	NameIndex       uint16
	Length          uint32
	NumberOfClasses uint16
	Classes         []uint16
}

func (n *NestMembers) Parse(nameIndex uint16, length uint32, data []byte) {
	n.NameIndex = nameIndex
	n.Length = length
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

func (n *NestMembers) GetName(constantPool []ConstantPoolInfo) string {
	return constantPool[n.NameIndex].String(constantPool)
}

func (n *NestMembers) String(constantPool []ConstantPoolInfo) string {
	result := ""
	for _, nestMember := range n.Classes {
		result += constantPool[nestMember].String(constantPool)
	}
	return result
}

type SourceDebugExtension struct {
	NameIndex      uint16
	Length         uint32
	DebugExtension []uint8
}

func (s *SourceDebugExtension) Parse(nameIndex uint16, length uint32, data []byte) {
	s.NameIndex = nameIndex
	s.Length = length
	binary.Read(bytes.NewBuffer(data[0:length]), binary.BigEndian, &s.DebugExtension)
}

func (s *SourceDebugExtension) GetName(constantPool []ConstantPoolInfo) string {
	return constantPool[s.NameIndex].String(constantPool)
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
	NameIndex uint16
	Length    uint32
	Num       uint16
	Methods   []BootStrapMethod
}

func (b *BootStrapMethod) Parse(data []byte, index int) int {
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

func (b *BootstrapMethods) Parse(nameIndex uint16, length uint32, data []byte) {
	b.NameIndex = nameIndex
	b.Length = length
	binary.Read(bytes.NewBuffer(data[0:2]), binary.BigEndian, &b.Num)
	b.Methods = make([]BootStrapMethod, b.Num)
	index := 2
	for n := 0; n < int(b.Num); n++ {
		method := &BootStrapMethod{}
		index = method.Parse(data, index)
		b.Methods[n] = *method
	}
}

func (b *BootstrapMethods) GetName(constantPool []ConstantPoolInfo) string {
	return constantPool[b.NameIndex].String(constantPool)
}

func (b *BootstrapMethods) String(constantPool []ConstantPoolInfo) string {
	result := ""
	for _, method := range b.Methods {
		result += constantPool[method.BootstrapMethodRef].String(constantPool)
	}
	return result
}
