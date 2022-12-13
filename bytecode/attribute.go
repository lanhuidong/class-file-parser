package bytecode

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
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
		item.parse(base, info, constantPool)
	case "Code":
		item = &Code{}
		item.parse(base, info, constantPool)
	/*case "StackMapTable":
	item = &StackMapTable{}
	item.parse(base, info, constantPool)*/
	case "Exceptions":
		item = &Exceptions{}
		item.parse(base, info, constantPool)
	case "InnerClasses":
		item = &InnerClasses{}
		item.parse(base, info, constantPool)
	case "EnclosingMethod":
		item = &EnclosingMethod{}
		item.parse(base, info, constantPool)
	case "Synthetic":
		item = &Synthetic{}
		item.parse(base, info, constantPool)
	case "Signature":
		item = &Signature{}
		item.parse(base, info, constantPool)
	case "SourceFile":
		item = &SourceFile{}
		item.parse(base, info, constantPool)
	case "SourceDebugExtension":
		item = &SourceDebugExtension{}
		item.parse(base, info, constantPool)
	case "LineNumberTable":
		item = &LineNumberTable{}
		item.parse(base, info, constantPool)
	case "LocalVariableTable":
		item = &LocalVariableTable{}
		item.parse(base, info, constantPool)
	case "LocalVariableTypeTable":
		item = &LocalVariableTypeTable{}
		item.parse(base, info, constantPool)
	case "Deprecated":
		item = &Deprecated{}
		item.parse(base, info, constantPool)
	case "RuntimeVisibleAnnotations", "RuntimeInvisibleAnnotations":
		item = &RuntimeVisibleAnnotations{}
		item.parse(base, info, constantPool)
	case "RuntimeVisibleParameterAnnotations", "RuntimeInvisibleParameterAnnotations":
		item = &RuntimeVisibleParameterAnnotations{}
		item.parse(base, info, constantPool)
	case "AnnotationDefault":
		item = &AnnotationDefault{}
		item.parse(base, info, constantPool)
	case "BootstrapMethods":
		item = &BootstrapMethods{}
		item.parse(base, info, constantPool)
	case "MethodParameters":
		item = &MethodParameters{}
		item.parse(base, info, constantPool)
	case "Module":
		item = &Module{}
		item.parse(base, info, constantPool)
	case "ModulePackages":
		item = &ModulePackages{}
		item.parse(base, info, constantPool)
	case "ModuleMainClass":
		item = &ModuleMainClass{}
		item.parse(base, info, constantPool)
	case "NestHost":
		item = &NestHost{}
		item.parse(base, info, constantPool)
	case "NestMembers":
		item = &NestMembers{}
		item.parse(base, info, constantPool)
	case "Record":
		item = &Record{}
		item.parse(base, info, constantPool)
	case "PermittedSubclasses":
		item = &PermittedSubclasses{}
		item.parse(base, info, constantPool)
	default:
		fmt.Printf("attribue name is %s\n", base.Name)
	}
	return 6 + int(base.Length), item
}

type AttributeInfo interface {
	parse(base *AttributeBase, data []byte, constantPool []ConstantPoolInfo)
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

func (c *ConstantValue) parse(base *AttributeBase, data []byte, constantPool []ConstantPoolInfo) {
	c.AttributeBase = *base
	binary.Read(bytes.NewBuffer(data), binary.BigEndian, &c.ConstantValueIndex)
}

func (c *ConstantValue) String(constantPool []ConstantPoolInfo) string {
	return constantPool[c.ConstantValueIndex].String(constantPool)
}

type ExceptionTable struct {
	StartPc   uint16
	EndPc     uint16
	HandlerPc uint16
	CatchType uint16
}

func (e *ExceptionTable) Parse(data []byte, index int) {
	binary.Read(bytes.NewBuffer(data[index:index+2]), binary.BigEndian, &e.StartPc)
	binary.Read(bytes.NewBuffer(data[index+2:index+4]), binary.BigEndian, &e.EndPc)
	binary.Read(bytes.NewBuffer(data[index+4:index+6]), binary.BigEndian, &e.HandlerPc)
	binary.Read(bytes.NewBuffer(data[index+6:index+8]), binary.BigEndian, &e.CatchType)
}

type Code struct {
	AttributeBase
	MaxStack             uint16
	MaxLocals            uint16
	CodeLength           uint32
	Code                 []byte
	ExceptionTableLength uint16
	Table                []ExceptionTable
	AttributesCount      uint16
	Attributes           []AttributeInfo
}

func (c *Code) parse(base *AttributeBase, data []byte, constantPool []ConstantPoolInfo) {
	c.AttributeBase = *base
	binary.Read(bytes.NewBuffer(data[0:2]), binary.BigEndian, &c.MaxStack)
	binary.Read(bytes.NewBuffer(data[2:4]), binary.BigEndian, &c.MaxLocals)
	binary.Read(bytes.NewBuffer(data[4:8]), binary.BigEndian, &c.CodeLength)
	c.Code = data[8 : 8+c.CodeLength]
	index := 8 + int(c.CodeLength)
	binary.Read(bytes.NewBuffer(data[index:index+2]), binary.BigEndian, &c.ExceptionTableLength)
	index += 2
	for i := 0; i < int(c.ExceptionTableLength); i++ {
		table := &ExceptionTable{}
		table.Parse(data, index)
		index += 8
		c.Table = append(c.Table, *table)
	}
	binary.Read(bytes.NewBuffer(data[index:index+2]), binary.BigEndian, &c.AttributesCount)
	index += 2
	_, attrs := ParseAttribute(int(c.AttributesCount), data, index, constantPool)
	c.Attributes = attrs
}

func (c *Code) String(constantPool []ConstantPoolInfo) string {
	result := fmt.Sprintf("max stack: %d, max locals: %d\n", c.MaxStack, c.MaxLocals)
	for _, attr := range c.Attributes {
		if attr != nil { //TODO 等所有解析完成后这个if应该移除
			result += attr.GetName() + "\n"
			result += attr.String(constantPool)
		}
	}
	return result
}

type StackMapTable struct {
	AttributeBase
	NumberOfEntries uint16
}

func (s *StackMapTable) parse(base *AttributeBase, data []byte, constantPool []ConstantPoolInfo) {
	s.AttributeBase = *base
	binary.Read(bytes.NewBuffer(data[0:2]), binary.BigEndian, &s.NumberOfEntries)
}

func (s *StackMapTable) String(constantPool []ConstantPoolInfo) string {
	return ""
}

type Exceptions struct {
	AttributeBase
	NumberOfExceptions  uint16
	ExceptionIndexTable []uint16
}

func (e *Exceptions) parse(base *AttributeBase, data []byte, constantPool []ConstantPoolInfo) {
	e.AttributeBase = *base
	binary.Read(bytes.NewBuffer(data[0:2]), binary.BigEndian, &e.NumberOfExceptions)
	index := 2
	for i := 0; i < int(e.NumberOfExceptions); i++ {
		var exceptionIndex uint16
		binary.Read(bytes.NewBuffer(data[index:index+2]), binary.BigEndian, &exceptionIndex)
		e.ExceptionIndexTable = append(e.ExceptionIndexTable, exceptionIndex)
		index += 2
	}
}

func (e *Exceptions) String(constantPool []ConstantPoolInfo) string {
	result := ""
	for _, index := range e.ExceptionIndexTable {
		result += constantPool[index].String(constantPool)
	}
	return result
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

func (i *InnerClasses) parse(base *AttributeBase, data []byte, constantPool []ConstantPoolInfo) {
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

func (e *EnclosingMethod) parse(base *AttributeBase, data []byte, constantPool []ConstantPoolInfo) {
	e.AttributeBase = *base
	binary.Read(bytes.NewBuffer(data[0:2]), binary.BigEndian, &e.ClassIndex)
	binary.Read(bytes.NewBuffer(data[2:4]), binary.BigEndian, &e.MethodIndex)
}

func (e *EnclosingMethod) String(constantPool []ConstantPoolInfo) string {
	return constantPool[e.ClassIndex].String(constantPool) + "." + constantPool[e.MethodIndex].String(constantPool)
}

type Synthetic struct {
	AttributeBase
}

func (s *Synthetic) parse(base *AttributeBase, data []byte, constantPool []ConstantPoolInfo) {
	s.AttributeBase = *base
}

func (s *Synthetic) String(constantPool []ConstantPoolInfo) string {
	return "Synthetic"
}

type Signature struct {
	AttributeBase
	SignatureIndex uint16
}

func (s *Signature) parse(base *AttributeBase, data []byte, constantPool []ConstantPoolInfo) {
	s.AttributeBase = *base
	binary.Read(bytes.NewBuffer(data), binary.BigEndian, &s.SignatureIndex)
}

func (s *Signature) String(constantPool []ConstantPoolInfo) string {
	return constantPool[s.SignatureIndex].String(constantPool)
}

type SourceFile struct {
	AttributeBase
	SourceFileIndex uint16
}

func (s *SourceFile) parse(base *AttributeBase, data []byte, constantPool []ConstantPoolInfo) {
	s.AttributeBase = *base
	binary.Read(bytes.NewBuffer(data), binary.BigEndian, &s.SourceFileIndex)
}

func (s *SourceFile) String(constantPool []ConstantPoolInfo) string {
	return constantPool[s.SourceFileIndex].String(constantPool)
}

type SourceDebugExtension struct {
	AttributeBase
	DebugExtension []uint8
}

func (s *SourceDebugExtension) parse(base *AttributeBase, data []byte, constantPool []ConstantPoolInfo) {
	s.AttributeBase = *base
	binary.Read(bytes.NewBuffer(data), binary.BigEndian, &s.DebugExtension)
}

func (s *SourceDebugExtension) String(constantPool []ConstantPoolInfo) string {
	return string(s.DebugExtension)
}

type LineNumber struct {
	StartPc    uint16
	LineNumber uint16
}

func (l *LineNumber) Parse(data []byte, index int) {
	binary.Read(bytes.NewBuffer(data[index:index+2]), binary.BigEndian, &l.StartPc)
	binary.Read(bytes.NewBuffer(data[index+2:index+4]), binary.BigEndian, &l.LineNumber)
}

type LineNumberTable struct {
	AttributeBase
	LineNumberTableLength uint16
	LineNumber            []LineNumber
}

func (l *LineNumberTable) parse(base *AttributeBase, data []byte, constantPool []ConstantPoolInfo) {
	l.AttributeBase = *base
	binary.Read(bytes.NewBuffer(data[0:2]), binary.BigEndian, &l.LineNumberTableLength)
	index := 2
	for i := 0; i < int(l.LineNumberTableLength); i++ {
		line := &LineNumber{}
		line.Parse(data, index)
		index += 4
		l.LineNumber = append(l.LineNumber, *line)
	}

}

func (l *LineNumberTable) String(constantPool []ConstantPoolInfo) string {
	result := ""
	for _, line := range l.LineNumber {
		result += fmt.Sprintf("start pc: %d, line number: %d\n", line.StartPc, line.LineNumber)
	}
	return result
}

type LocalVariable struct {
	StartPc         uint16
	Length          uint16
	NameIndex       uint16
	DescriptorIndex uint16
	Index           uint16
}

func (l *LocalVariable) Parse(data []byte, index int) {
	binary.Read(bytes.NewBuffer(data[index:index+2]), binary.BigEndian, &l.StartPc)
	binary.Read(bytes.NewBuffer(data[index+2:index+4]), binary.BigEndian, &l.Length)
	binary.Read(bytes.NewBuffer(data[index+4:index+6]), binary.BigEndian, &l.NameIndex)
	binary.Read(bytes.NewBuffer(data[index+6:index+8]), binary.BigEndian, &l.DescriptorIndex)
	binary.Read(bytes.NewBuffer(data[index+8:index+10]), binary.BigEndian, &l.Index)
}

type LocalVariableTable struct {
	AttributeBase
	LocalVariableTableLength uint16
	LocalVariable            []LocalVariable
}

func (l *LocalVariableTable) parse(base *AttributeBase, data []byte, constantPool []ConstantPoolInfo) {
	l.AttributeBase = *base
	binary.Read(bytes.NewBuffer(data[0:2]), binary.BigEndian, &l.LocalVariableTableLength)
	index := 2
	for i := 0; i < int(l.LocalVariableTableLength); i++ {
		localVar := &LocalVariable{}
		localVar.Parse(data, index)
		index += 10
		l.LocalVariable = append(l.LocalVariable, *localVar)
	}

}

func (l *LocalVariableTable) String(constantPool []ConstantPoolInfo) string {
	result := ""
	for _, localVar := range l.LocalVariable {
		result += fmt.Sprintf("start pc: %d, length: %d, name index: %d, descriptor index: %d, index: %d\n", localVar.StartPc, localVar.Length, localVar.NameIndex, localVar.DescriptorIndex, localVar.Index)
	}
	return result
}

type LocalVariableType struct {
	StartPc        uint16
	Length         uint16
	NameIndex      uint16
	SignatureIndex uint16
	Index          uint16
}

func (l *LocalVariableType) Parse(data []byte, index int) {
	binary.Read(bytes.NewBuffer(data[index:index+2]), binary.BigEndian, &l.StartPc)
	binary.Read(bytes.NewBuffer(data[index+2:index+4]), binary.BigEndian, &l.Length)
	binary.Read(bytes.NewBuffer(data[index+4:index+6]), binary.BigEndian, &l.NameIndex)
	binary.Read(bytes.NewBuffer(data[index+6:index+8]), binary.BigEndian, &l.SignatureIndex)
	binary.Read(bytes.NewBuffer(data[index+8:index+10]), binary.BigEndian, &l.Index)
}

type LocalVariableTypeTable struct {
	AttributeBase
	LocalVariableTypeTableLength uint16
	LocalVariableType            []LocalVariableType
}

func (l *LocalVariableTypeTable) parse(base *AttributeBase, data []byte, constantPool []ConstantPoolInfo) {
	l.AttributeBase = *base
	binary.Read(bytes.NewBuffer(data[0:2]), binary.BigEndian, &l.LocalVariableTypeTableLength)
	index := 2
	for i := 0; i < int(l.LocalVariableTypeTableLength); i++ {
		localVar := &LocalVariableType{}
		localVar.Parse(data, index)
		index += 10
		l.LocalVariableType = append(l.LocalVariableType, *localVar)
	}

}

func (l *LocalVariableTypeTable) String(constantPool []ConstantPoolInfo) string {
	result := ""
	for _, localVar := range l.LocalVariableType {
		result += fmt.Sprintf("start pc: %d, length: %d, name index: %d, signature index: %d, index: %d\n", localVar.StartPc, localVar.Length, localVar.NameIndex, localVar.SignatureIndex, localVar.Index)
	}
	return result
}

type Deprecated struct {
	AttributeBase
}

func (d *Deprecated) parse(base *AttributeBase, data []byte, constantPool []ConstantPoolInfo) {
	d.AttributeBase = *base
	if d.Length != 0 {
		panic("attribute deprecated's length must be 0, but actual is " + strconv.Itoa(int(d.Length)))
	}
}

func (d *Deprecated) String(constantPool []ConstantPoolInfo) string {
	return ""
}

type ArrayValue struct {
	NumValues uint16
	Values    []ElementValue
}

func (a *ArrayValue) parse(data []byte, index int) int {
	binary.Read(bytes.NewBuffer(data[index:index+2]), binary.BigEndian, &a.NumValues)
	index += 2
	for i := 0; i < int(a.NumValues); i++ {
		elem := &ElementValue{}
		index = elem.parse(data, index)
		a.Values = append(a.Values, *elem)
	}
	return index
}

type EnumConstValue struct {
	TypeNameIndex  uint16
	ConstNameIndex uint16
}

func (e *EnumConstValue) parse(data []byte) {
	binary.Read(bytes.NewBuffer(data[0:2]), binary.BigEndian, &e.TypeNameIndex)
	binary.Read(bytes.NewBuffer(data[2:4]), binary.BigEndian, &e.ConstNameIndex)
}

type ElementValue struct {
	Tag             uint8
	ConstValueIndex uint16
	EnumConstValue
	ClassInfoIndex  uint16
	AnnotationValue Annotation
	ArrayValue
}

func (e *ElementValue) parse(data []byte, index int) int {
	binary.Read(bytes.NewBuffer(data[index:index+1]), binary.BigEndian, &e.Tag)
	switch e.Tag {
	case 'B', 'C', 'D', 'F', 'I', 'J', 'S', 'Z', 's':
		binary.Read(bytes.NewBuffer(data[1:3]), binary.BigEndian, &e.ConstValueIndex)
		index += 3
	case 'e':
		value := &EnumConstValue{}
		value.parse(data[index+1 : index+5])
		e.EnumConstValue = *value
		index += 5
	case 'c':
		binary.Read(bytes.NewBuffer(data[1:3]), binary.BigEndian, &e.ClassInfoIndex)
		index += 3
	case '@':
		ann := &Annotation{}
		index = ann.parse(data, index+1)
		e.AnnotationValue = *ann
	case '[':
		arr := &ArrayValue{}
		index = arr.parse(data, index+1)
		e.ArrayValue = *arr
	default:
		panic(fmt.Sprintf("unknown element value tag: %d(%c)\n", e.Tag, e.Tag))
	}
	return index
}

type ElementValuePairs struct {
	ElementNameIndex uint16
	ElementValue
}

func (e *ElementValuePairs) parse(data []byte, index int) int {
	binary.Read(bytes.NewBuffer(data[index:index+2]), binary.BigEndian, &e.ElementNameIndex)
	elem := &ElementValue{}
	index = elem.parse(data, index+2)
	e.ElementValue = *elem
	return index
}

type Annotation struct {
	TypeIndex            uint16
	NumElementValuePairs uint16
	ValuePairs           []ElementValuePairs
}

func (a *Annotation) parse(data []byte, index int) int {
	binary.Read(bytes.NewBuffer(data[index:index+2]), binary.BigEndian, &a.TypeIndex)
	binary.Read(bytes.NewBuffer(data[index+2:index+4]), binary.BigEndian, &a.NumElementValuePairs)
	index += 4
	for i := 0; i < int(a.NumElementValuePairs); i++ {
		pair := &ElementValuePairs{}
		index = pair.parse(data, index)
		a.ValuePairs = append(a.ValuePairs, *pair)
	}
	return index
}

type RuntimeVisibleAnnotations struct {
	AttributeBase
	NumAnnotations uint16
	Annotations    []Annotation
}

func (r *RuntimeVisibleAnnotations) parse(base *AttributeBase, data []byte, constantPool []ConstantPoolInfo) {
	r.AttributeBase = *base
	binary.Read(bytes.NewBuffer(data[0:2]), binary.BigEndian, &r.NumAnnotations)
	index := 2
	for i := 0; i < int(r.NumAnnotations); i++ {
		ann := &Annotation{}
		index = ann.parse(data, index)
		r.Annotations = append(r.Annotations, *ann)
	}
}

func (r *RuntimeVisibleAnnotations) String(constantPool []ConstantPoolInfo) string {
	return fmt.Sprintf("%d runtime visible annotation\n", r.NumAnnotations)
}

type ParameterAnnotation struct {
	NumAnnotations uint16
	Annotations    []Annotation
}

func (p *ParameterAnnotation) parse(data []byte, index int) int {
	binary.Read(bytes.NewBuffer(data[index:index+2]), binary.BigEndian, &p.NumAnnotations)
	index += 2
	for i := 0; i < int(p.NumAnnotations); i++ {
		ann := &Annotation{}
		index = ann.parse(data, index)
		p.Annotations = append(p.Annotations, *ann)
	}
	return index
}

type RuntimeVisibleParameterAnnotations struct {
	AttributeBase
	NumParameters        uint8
	ParameterAnnotations []ParameterAnnotation
}

func (r *RuntimeVisibleParameterAnnotations) parse(base *AttributeBase, data []byte, constantPool []ConstantPoolInfo) {
	r.AttributeBase = *base
	binary.Read(bytes.NewBuffer(data[0:1]), binary.BigEndian, &r.NumParameters)
	index := 1
	for i := 0; i < int(r.NumParameters); i++ {
		param := &ParameterAnnotation{}
		index = param.parse(data, index)
		r.ParameterAnnotations = append(r.ParameterAnnotations, *param)
	}
}

func (r *RuntimeVisibleParameterAnnotations) String(constantPool []ConstantPoolInfo) string {
	return ""
}

type AnnotationDefault struct {
	AttributeBase
	DefaultValue ElementValue
}

func (a *AnnotationDefault) parse(base *AttributeBase, data []byte, constantPool []ConstantPoolInfo) {
	a.AttributeBase = *base
	a.DefaultValue = ElementValue{}
	a.DefaultValue.parse(data, 0)
}

func (a *AnnotationDefault) String(constantPool []ConstantPoolInfo) string {
	return ""
}

type BootStrapMethod struct {
	BootstrapMethodRef uint16
	ArgumentsNum       uint16
	Arguments          []uint16
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

type BootstrapMethods struct {
	AttributeBase
	Num     uint16
	Methods []BootStrapMethod
}

func (b *BootstrapMethods) parse(base *AttributeBase, data []byte, constantPool []ConstantPoolInfo) {
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

type MethodParameter struct {
	NameIndex   uint16
	AccessFlags uint16
}

func (m *MethodParameter) parse(data []byte) {
	binary.Read(bytes.NewBuffer(data[0:2]), binary.BigEndian, &m.NameIndex)
	binary.Read(bytes.NewBuffer(data[2:4]), binary.BigEndian, &m.AccessFlags)
}

type MethodParameters struct {
	AttributeBase
	ParametersCount uint8
	parameter       []MethodParameter
}

func (m *MethodParameters) parse(base *AttributeBase, data []byte, constantPool []ConstantPoolInfo) {
	m.AttributeBase = *base
	binary.Read(bytes.NewBuffer(data[0:1]), binary.BigEndian, &m.ParametersCount)
	index := 1
	for n := 0; n < int(m.ParametersCount); n++ {
		param := &MethodParameter{}
		param.parse(data[index:])
		m.parameter = append(m.parameter, *param)
		index += 4
	}
}

func (m *MethodParameters) String(constantPool []ConstantPoolInfo) string {
	result := ""
	for _, param := range m.parameter {
		result += constantPool[param.NameIndex].String(constantPool) + " "
	}
	return result
}

type Require struct {
	RequiresIndex        uint16
	RequiresFlags        uint16
	RequiresVersionIndex uint16
}

func (r *Require) parse(data []byte, index int) {
	binary.Read(bytes.NewBuffer(data[index:index+2]), binary.BigEndian, &r.RequiresIndex)
	binary.Read(bytes.NewBuffer(data[index+2:index+4]), binary.BigEndian, &r.RequiresFlags)
	binary.Read(bytes.NewBuffer(data[index+4:index+6]), binary.BigEndian, &r.RequiresVersionIndex)
}

type Export struct {
	ExportsIndex   uint16
	ExportsFlags   uint16
	ExportsToCount uint16
	ExportsToIndex []uint16
}

func (e *Export) parse(data []byte, index int) int {
	binary.Read(bytes.NewBuffer(data[index:index+2]), binary.BigEndian, &e.ExportsIndex)
	binary.Read(bytes.NewBuffer(data[index+2:index+4]), binary.BigEndian, &e.ExportsFlags)
	binary.Read(bytes.NewBuffer(data[index+4:index+6]), binary.BigEndian, &e.ExportsToCount)
	index += 6
	for i := 0; i < int(e.ExportsToCount); i++ {
		var idx uint16
		binary.Read(bytes.NewBuffer(data[index:index+2]), binary.BigEndian, &idx)
		e.ExportsToIndex = append(e.ExportsToIndex, idx)
		index += 2
	}
	return index
}

type Open struct {
	OpenIndex   uint16
	OpenFlags   uint16
	OpenToCount uint16
	OpenToIndex []uint16
}

func (o *Open) parse(data []byte, index int) int {
	binary.Read(bytes.NewBuffer(data[index:index+2]), binary.BigEndian, &o.OpenIndex)
	binary.Read(bytes.NewBuffer(data[index+2:index+4]), binary.BigEndian, &o.OpenFlags)
	binary.Read(bytes.NewBuffer(data[index+4:index+6]), binary.BigEndian, &o.OpenToCount)
	index += 6
	for i := 0; i < int(o.OpenToCount); i++ {
		var idx uint16
		binary.Read(bytes.NewBuffer(data[index:index+2]), binary.BigEndian, &idx)
		o.OpenToIndex = append(o.OpenToIndex, idx)
		index += 2
	}
	return index
}

type Provide struct {
	ProvidesIndex     uint16
	ProvidesWithCount uint16
	ProvidesWithIndex []uint16
}

func (p *Provide) parse(data []byte, index int) int {
	binary.Read(bytes.NewBuffer(data[index:index+2]), binary.BigEndian, &p.ProvidesIndex)
	binary.Read(bytes.NewBuffer(data[index+2:index+4]), binary.BigEndian, &p.ProvidesWithCount)
	index += 4
	for i := 0; i < int(p.ProvidesWithCount); i++ {
		var idx uint16
		binary.Read(bytes.NewBuffer(data[index:index+2]), binary.BigEndian, &idx)
		p.ProvidesWithIndex = append(p.ProvidesWithIndex, idx)
		index += 2
	}
	return index
}

type Module struct {
	AttributeBase
	ModuleNameIndex    uint16
	ModuleFlags        uint16
	ModuleVersionIndex uint16
	RequiresCount      uint16
	Requires           []Require
	ExportsCount       uint16
	Exports            []Export
	OpenCount          uint16
	Opens              []Open
	UsesCount          uint16
	UsesIndex          []uint16
	ProvidesCount      uint16
	Provides           []Provide
}

func (m *Module) parse(base *AttributeBase, data []byte, constantPool []ConstantPoolInfo) {
	m.AttributeBase = *base
	index := 0
	binary.Read(bytes.NewBuffer(data[index:index+2]), binary.BigEndian, &m.ModuleNameIndex)
	binary.Read(bytes.NewBuffer(data[index+2:index+4]), binary.BigEndian, &m.ModuleFlags)
	binary.Read(bytes.NewBuffer(data[index+4:index+6]), binary.BigEndian, &m.ModuleVersionIndex)

	binary.Read(bytes.NewBuffer(data[index+6:index+8]), binary.BigEndian, &m.RequiresCount)
	index += 8
	for i := 0; i < int(m.RequiresCount); i++ {
		r := &Require{}
		r.parse(data, index)
		m.Requires = append(m.Requires, *r)
		index += 6
	}

	binary.Read(bytes.NewBuffer(data[index:index+2]), binary.BigEndian, &m.ExportsCount)
	index += 2
	for i := 0; i < int(m.ExportsCount); i++ {
		e := &Export{}
		index = e.parse(data, index)
		m.Exports = append(m.Exports, *e)
	}

	binary.Read(bytes.NewBuffer(data[index:index+2]), binary.BigEndian, &m.OpenCount)
	index += 2
	for i := 0; i < int(m.OpenCount); i++ {
		o := &Open{}
		index = o.parse(data, index)
		m.Opens = append(m.Opens, *o)
	}

	binary.Read(bytes.NewBuffer(data[index:index+2]), binary.BigEndian, &m.UsesCount)
	index += 2
	for i := 0; i < int(m.UsesCount); i++ {
		var idx uint16
		binary.Read(bytes.NewBuffer(data[index:index+2]), binary.BigEndian, &idx)
		m.UsesIndex = append(m.UsesIndex, idx)
		idx += 2
	}

	binary.Read(bytes.NewBuffer(data[index:index+2]), binary.BigEndian, &m.ProvidesCount)
	index += 2
	for i := 0; i < int(m.ProvidesCount); i++ {
		p := &Provide{}
		index = p.parse(data, index)
		m.Provides = append(m.Provides, *p)
	}
}

func (m *Module) String(constantPool []ConstantPoolInfo) string {
	return constantPool[m.ModuleNameIndex].String(constantPool)
}

type ModulePackages struct {
	AttributeBase
	PackageCount uint16
	PackageIndex []uint16
}

func (m *ModulePackages) parse(base *AttributeBase, data []byte, constantPool []ConstantPoolInfo) {
	m.AttributeBase = *base
	binary.Read(bytes.NewBuffer(data[0:2]), binary.BigEndian, &m.PackageCount)
	index := 2
	for i := 0; i < int(m.PackageCount); i++ {
		var packageIndex uint16
		binary.Read(bytes.NewBuffer(data[index:index+2]), binary.BigEndian, &packageIndex)
		index += 2
		m.PackageIndex = append(m.PackageIndex, packageIndex)
	}
}

func (m *ModulePackages) String(constantPool []ConstantPoolInfo) string {
	result := ""
	for _, index := range m.PackageIndex {
		result += "\n" + constantPool[index].String(constantPool)
	}
	return result
}

type ModuleMainClass struct {
	AttributeBase
	MainClassIndex uint16
}

func (m *ModuleMainClass) parse(base *AttributeBase, data []byte, constantPool []ConstantPoolInfo) {
	m.AttributeBase = *base
	binary.Read(bytes.NewBuffer(data[0:2]), binary.BigEndian, &m.MainClassIndex)
}

func (m *ModuleMainClass) String(constantPool []ConstantPoolInfo) string {
	return constantPool[m.MainClassIndex].String(constantPool)
}

type NestHost struct {
	AttributeBase
	HostClassIndex uint16
}

func (n *NestHost) parse(base *AttributeBase, data []byte, constantPool []ConstantPoolInfo) {
	n.AttributeBase = *base
	binary.Read(bytes.NewBuffer(data[0:2]), binary.BigEndian, &n.HostClassIndex)
}

func (n *NestHost) String(constantPool []ConstantPoolInfo) string {
	return constantPool[n.HostClassIndex].String(constantPool)
}

type NestMembers struct {
	AttributeBase
	NumberOfClasses uint16
	Classes         []uint16
}

func (n *NestMembers) parse(base *AttributeBase, data []byte, constantPool []ConstantPoolInfo) {
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

type RecordComponent struct {
	NameIndex       uint16
	DescriptorIndex uint16
	AttributesCount uint16
	Attributes      []AttributeInfo
}

func (r *RecordComponent) parse(data []byte, index int, constantPool []ConstantPoolInfo) int {
	binary.Read(bytes.NewBuffer(data[index:index+2]), binary.BigEndian, &r.NameIndex)
	binary.Read(bytes.NewBuffer(data[index+2:index+4]), binary.BigEndian, &r.DescriptorIndex)
	binary.Read(bytes.NewBuffer(data[index+4:index+6]), binary.BigEndian, &r.AttributesCount)
	index += 6
	_, attrs := ParseAttribute(int(r.AttributesCount), data, index, constantPool)
	r.Attributes = attrs
	return index
}

type Record struct {
	AttributeBase
	ComponentsCount     uint16
	RecordComponentInfo []RecordComponent
}

func (r *Record) parse(base *AttributeBase, data []byte, constantPool []ConstantPoolInfo) {
	r.AttributeBase = *base
	binary.Read(bytes.NewBuffer(data[0:2]), binary.BigEndian, &r.ComponentsCount)
	index := 2
	for i := 0; i < int(r.ComponentsCount); i++ {
		component := &RecordComponent{}
		index = component.parse(data, index, constantPool)
		r.RecordComponentInfo = append(r.RecordComponentInfo, *component)
	}
}

func (b *Record) String(constantPool []ConstantPoolInfo) string {
	result := ""
	for _, component := range b.RecordComponentInfo {
		result += " name: " + constantPool[component.NameIndex].String(constantPool)
	}
	return result
}

type PermittedSubclasses struct {
	AttributeBase
	NumberOfClasses uint16
	Classes         []uint16
}

func (p *PermittedSubclasses) parse(base *AttributeBase, data []byte, constantPool []ConstantPoolInfo) {
	p.AttributeBase = *base
	binary.Read(bytes.NewBuffer(data[0:2]), binary.BigEndian, &p.NumberOfClasses)
	index := 2
	for i := 0; i < int(p.NumberOfClasses); i++ {
		var classIndex uint16
		binary.Read(bytes.NewBuffer(data[index:index+2]), binary.BigEndian, &classIndex)
		p.Classes = append(p.Classes, classIndex)
		index += 2
	}
}

func (p *PermittedSubclasses) String(constantPool []ConstantPoolInfo) string {
	result := ""
	for _, class := range p.Classes {
		result += constantPool[class].String(constantPool) + " "
	}
	return result
}
