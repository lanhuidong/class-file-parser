package bytecode

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type ConstantPoolInfo interface {
	TagValue() uint8

	TagName() string

	Parse(data []byte, index int) int

	String(constantPool []ConstantPoolInfo) string
}

type ConstantUtf8 struct {
	Tag    uint8
	Length uint16
	Value  []byte
}

func (c *ConstantUtf8) TagValue() uint8 {
	return 1
}

func (c *ConstantUtf8) TagName() string {
	return "Utf8"
}

func (c *ConstantUtf8) String(constantPool []ConstantPoolInfo) string {
	return string(c.Value)
}

func (c *ConstantUtf8) Parse(data []byte, index int) int {
	binary.Read(bytes.NewBuffer(data[index:index+1]), binary.BigEndian, &c.Tag)
	binary.Read(bytes.NewBuffer(data[index+1:index+3]), binary.BigEndian, &c.Length)
	c.Value = data[index+3 : index+3+int(c.Length)]
	return 3 + int(c.Length)
}

type ConstantInteger struct {
	Tag   uint8
	Value int32
}

func (c *ConstantInteger) TagValue() uint8 {
	return 3
}

func (c *ConstantInteger) TagName() string {
	return "Integer"
}

func (c *ConstantInteger) String(constantPool []ConstantPoolInfo) string {
	return fmt.Sprintf("%d", c.Value)
}

func (c *ConstantInteger) Parse(data []byte, index int) int {
	binary.Read(bytes.NewBuffer(data[index:index+1]), binary.BigEndian, &c.Tag)
	binary.Read(bytes.NewBuffer(data[index+1:index+5]), binary.BigEndian, &c.Value)
	return 5
}

type ConstantFloat struct {
	Tag   uint8
	Value float32
}

func (c *ConstantFloat) TagValue() uint8 {
	return 4
}

func (c *ConstantFloat) TagName() string {
	return "Float"
}

func (c *ConstantFloat) String(constantPool []ConstantPoolInfo) string {
	return fmt.Sprintf("%f", c.Value)
}

func (c *ConstantFloat) Parse(data []byte, index int) int {
	binary.Read(bytes.NewBuffer(data[index:index+1]), binary.BigEndian, &c.Tag)
	binary.Read(bytes.NewBuffer(data[index+1:index+5]), binary.BigEndian, &c.Value)
	return 5
}

type ConstantLong struct {
	Tag   uint8
	Value int64
}

func (c *ConstantLong) TagValue() uint8 {
	return 5
}

func (c *ConstantLong) TagName() string {
	return "Long"
}

func (c *ConstantLong) String(constantPool []ConstantPoolInfo) string {
	return fmt.Sprintf("%d", c.Value)
}

func (c *ConstantLong) Parse(data []byte, index int) int {
	binary.Read(bytes.NewBuffer(data[index:index+1]), binary.BigEndian, &c.Tag)
	binary.Read(bytes.NewBuffer(data[index+1:index+9]), binary.BigEndian, &c.Value)
	return 9
}

type ConstantDouble struct {
	Tag   uint8
	Value float64
}

func (c *ConstantDouble) TagValue() uint8 {
	return 6
}

func (c *ConstantDouble) TagName() string {
	return "Double"
}

func (c *ConstantDouble) String(constantPool []ConstantPoolInfo) string {
	return fmt.Sprintf("%f", c.Value)
}

func (c *ConstantDouble) Parse(data []byte, index int) int {
	binary.Read(bytes.NewBuffer(data[index:index+1]), binary.BigEndian, &c.Tag)
	binary.Read(bytes.NewBuffer(data[index+1:index+9]), binary.BigEndian, &c.Value)
	return 9
}

type ConstantClass struct {
	Tag       uint8
	NameIndex uint16
}

func (c *ConstantClass) TagValue() uint8 {
	return 7
}

func (c *ConstantClass) TagName() string {
	return "Class"
}

func (c *ConstantClass) String(constantPool []ConstantPoolInfo) string {
	return fmt.Sprintf("#%d", c.NameIndex) + "	//" + constantPool[c.NameIndex].String(constantPool)
}

func (c *ConstantClass) Parse(data []byte, index int) int {
	binary.Read(bytes.NewBuffer(data[index:index+1]), binary.BigEndian, &c.Tag)
	binary.Read(bytes.NewBuffer(data[index+1:index+3]), binary.BigEndian, &c.NameIndex)
	return 3
}

type ConstantString struct {
	Tag         uint8
	StringIndex uint16
}

func (c *ConstantString) TagValue() uint8 {
	return 8
}

func (c *ConstantString) TagName() string {
	return "String"
}

func (c *ConstantString) String(constantPool []ConstantPoolInfo) string {
	return fmt.Sprintf("#%d", c.StringIndex) + "	//" + constantPool[c.StringIndex].String(constantPool)
}

func (c *ConstantString) Parse(data []byte, index int) int {
	binary.Read(bytes.NewBuffer(data[index:index+1]), binary.BigEndian, &c.Tag)
	binary.Read(bytes.NewBuffer(data[index+1:index+3]), binary.BigEndian, &c.StringIndex)
	return 3
}

type ConstantFieldref struct {
	Tag              uint8
	ClassIndex       uint16
	NameAndTypeIndex uint16
}

func (c *ConstantFieldref) TagValue() uint8 {
	return 9
}

func (c *ConstantFieldref) TagName() string {
	return "Fieldref"
}

func (c *ConstantFieldref) String(constantPool []ConstantPoolInfo) string {
	return fmt.Sprintf("#%d.#%d", c.ClassIndex, c.NameAndTypeIndex) + "	//" + constantPool[c.ClassIndex].String(constantPool) + "." + constantPool[c.NameAndTypeIndex].String(constantPool)
}

func (c *ConstantFieldref) Parse(data []byte, index int) int {
	binary.Read(bytes.NewBuffer(data[index:index+1]), binary.BigEndian, &c.Tag)
	binary.Read(bytes.NewBuffer(data[index+1:index+3]), binary.BigEndian, &c.ClassIndex)
	binary.Read(bytes.NewBuffer(data[index+3:index+5]), binary.BigEndian, &c.NameAndTypeIndex)
	return 5
}

type ConstantMethodref struct {
	Tag              uint8
	ClassIndex       uint16
	NameAndTypeIndex uint16
}

func (c *ConstantMethodref) TagValue() uint8 {
	return 10
}

func (c *ConstantMethodref) TagName() string {
	return "Methodref"
}

func (c *ConstantMethodref) String(constantPool []ConstantPoolInfo) string {
	return fmt.Sprintf("#%d.#%d", c.ClassIndex, c.NameAndTypeIndex) + "	//" + constantPool[c.ClassIndex].String(constantPool) + "." + constantPool[c.NameAndTypeIndex].String(constantPool)
}

func (c *ConstantMethodref) Parse(data []byte, index int) int {
	binary.Read(bytes.NewBuffer(data[index:index+1]), binary.BigEndian, &c.Tag)
	binary.Read(bytes.NewBuffer(data[index+1:index+3]), binary.BigEndian, &c.ClassIndex)
	binary.Read(bytes.NewBuffer(data[index+3:index+5]), binary.BigEndian, &c.NameAndTypeIndex)
	return 5
}

type ConstantInterfaceMethodref struct {
	Tag              uint8
	ClassIndex       uint16
	NameAndTypeIndex uint16
}

func (c *ConstantInterfaceMethodref) TagValue() uint8 {
	return 11
}

func (c *ConstantInterfaceMethodref) TagName() string {
	return "InterfaceMethodref"
}

func (c *ConstantInterfaceMethodref) String(constantPool []ConstantPoolInfo) string {
	return fmt.Sprintf("#%d.#%d", c.ClassIndex, c.NameAndTypeIndex) + "	//" + constantPool[c.ClassIndex].String(constantPool) + "." + constantPool[c.NameAndTypeIndex].String(constantPool)
}

func (c *ConstantInterfaceMethodref) Parse(data []byte, index int) int {
	binary.Read(bytes.NewBuffer(data[index:index+1]), binary.BigEndian, &c.Tag)
	binary.Read(bytes.NewBuffer(data[index+1:index+3]), binary.BigEndian, &c.ClassIndex)
	binary.Read(bytes.NewBuffer(data[index+3:index+5]), binary.BigEndian, &c.NameAndTypeIndex)
	return 5
}

type ConstantNameAndType struct {
	Tag             uint8
	NameIndex       uint16
	DescriptorIndex uint16
}

func (c *ConstantNameAndType) TagValue() uint8 {
	return 12
}

func (c *ConstantNameAndType) TagName() string {
	return "NameAndType"
}

func (c *ConstantNameAndType) String(constantPool []ConstantPoolInfo) string {
	return fmt.Sprintf("#%d.#%d", c.NameIndex, c.DescriptorIndex) + "	//" + constantPool[c.NameIndex].String(constantPool) + "." + constantPool[c.DescriptorIndex].String(constantPool)
}

func (c *ConstantNameAndType) Parse(data []byte, index int) int {
	binary.Read(bytes.NewBuffer(data[index:index+1]), binary.BigEndian, &c.Tag)
	binary.Read(bytes.NewBuffer(data[index+1:index+3]), binary.BigEndian, &c.NameIndex)
	binary.Read(bytes.NewBuffer(data[index+3:index+5]), binary.BigEndian, &c.DescriptorIndex)
	return 5
}

type ConstantMethodHandle struct {
	Tag            uint8
	ReferenceKind  uint8
	ReferenceIndex uint16
}

func (c *ConstantMethodHandle) TagValue() uint8 {
	return 15
}

func (c *ConstantMethodHandle) TagName() string {
	return "MethodHandle"
}

func (c *ConstantMethodHandle) String(constantPool []ConstantPoolInfo) string {
	return c.TagName() + "	" + fmt.Sprintf("kind: %d.#%d", c.ReferenceKind, c.ReferenceKind)
}

func (c *ConstantMethodHandle) Parse(data []byte, index int) int {
	binary.Read(bytes.NewBuffer(data[index:index+1]), binary.BigEndian, &c.Tag)
	binary.Read(bytes.NewBuffer(data[index+1:index+2]), binary.BigEndian, &c.ReferenceKind)
	binary.Read(bytes.NewBuffer(data[index+2:index+4]), binary.BigEndian, &c.ReferenceIndex)
	return 4
}

type ConstantMethodType struct {
	Tag             uint8
	DescriptorIndex uint16
}

func (c *ConstantMethodType) TagValue() uint8 {
	return 16
}

func (c *ConstantMethodType) TagName() string {
	return "MethodType"
}

func (c *ConstantMethodType) String(constantPool []ConstantPoolInfo) string {
	return c.TagName() + "	" + fmt.Sprintf("#%d", c.DescriptorIndex)
}

func (c *ConstantMethodType) Parse(data []byte, index int) int {
	binary.Read(bytes.NewBuffer(data[index:index+1]), binary.BigEndian, &c.Tag)
	binary.Read(bytes.NewBuffer(data[index+1:index+3]), binary.BigEndian, &c.DescriptorIndex)
	return 3
}

type ConstantDynamic struct {
	Tag                      uint8
	BootstrapMethodAttrIndex uint16
	NameAndTypeIndex         uint16
}

func (c *ConstantDynamic) TagValue() uint8 {
	return 17
}

func (c *ConstantDynamic) TagName() string {
	return "Dynamic"
}

func (c *ConstantDynamic) String(constantPool []ConstantPoolInfo) string {
	return c.TagName() + "	" + fmt.Sprintf("#%d.#%d", c.BootstrapMethodAttrIndex, c.NameAndTypeIndex)
}

func (c *ConstantDynamic) Parse(data []byte, index int) int {
	binary.Read(bytes.NewBuffer(data[index:index+1]), binary.BigEndian, &c.Tag)
	binary.Read(bytes.NewBuffer(data[index+1:index+3]), binary.BigEndian, &c.BootstrapMethodAttrIndex)
	binary.Read(bytes.NewBuffer(data[index+3:index+5]), binary.BigEndian, &c.NameAndTypeIndex)
	return 5
}

type ConstantInvokeDynamic struct {
	Tag                      uint8
	BootstrapMethodAttrIndex uint16
	NameAndTypeIndex         uint16
}

func (c *ConstantInvokeDynamic) TagValue() uint8 {
	return 18
}

func (c *ConstantInvokeDynamic) TagName() string {
	return "InvokeDynamic"
}

func (c *ConstantInvokeDynamic) String(constantPool []ConstantPoolInfo) string {
	return c.TagName() + "	" + fmt.Sprintf("#%d.#%d", c.BootstrapMethodAttrIndex, c.NameAndTypeIndex)
}

func (c *ConstantInvokeDynamic) Parse(data []byte, index int) int {
	binary.Read(bytes.NewBuffer(data[index:index+1]), binary.BigEndian, &c.Tag)
	binary.Read(bytes.NewBuffer(data[index+1:index+3]), binary.BigEndian, &c.BootstrapMethodAttrIndex)
	binary.Read(bytes.NewBuffer(data[index+3:index+5]), binary.BigEndian, &c.NameAndTypeIndex)
	return 5
}

type ConstantModule struct {
	Tag       uint8
	NameIndex uint16
}

func (c *ConstantModule) TagValue() uint8 {
	return 19
}

func (c *ConstantModule) TagName() string {
	return "Module"
}

func (c *ConstantModule) String(constantPool []ConstantPoolInfo) string {
	return c.TagName() + "	" + fmt.Sprintf("#%d", c.NameIndex)
}

func (c *ConstantModule) Parse(data []byte, index int) int {
	binary.Read(bytes.NewBuffer(data[index:index+1]), binary.BigEndian, &c.Tag)
	binary.Read(bytes.NewBuffer(data[index+1:index+3]), binary.BigEndian, &c.NameIndex)
	return 3
}

type ConstantPackage struct {
	Tag       uint8
	NameIndex uint16
}

func (c *ConstantPackage) TagValue() uint8 {
	return 20
}

func (c *ConstantPackage) TagName() string {
	return "Package"
}

func (c *ConstantPackage) String(constantPool []ConstantPoolInfo) string {
	return c.TagName() + "	" + fmt.Sprintf("#%d", c.NameIndex)
}

func (c *ConstantPackage) Parse(data []byte, index int) int {
	binary.Read(bytes.NewBuffer(data[index:index+1]), binary.BigEndian, &c.Tag)
	binary.Read(bytes.NewBuffer(data[index+1:index+3]), binary.BigEndian, &c.NameIndex)
	return 3
}
