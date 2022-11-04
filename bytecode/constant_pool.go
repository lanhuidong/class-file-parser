package bytecode

const ACC_PUBLIC = 0x0001
const ACC_FINAL = 0x0010
const ACC_SUPER = 0x0020
const ACC_INTERFACE = 0x0200
const ACC_ABSTRACT = 0x0400
const ACC_SYNTHETIC = 0x1000
const ACC_ANNOTATION = 0x2000
const ACC_ENUM = 0x4000
const ACC_MODULE = 0x8000

type TagValue uint8

type ConstantPoolInfo interface {
	TagValue() TagValue

	TagName() string
}

type ConstantUtf8 struct {
}

func (c *ConstantUtf8) TagValue() uint8 {
	return 1
}

func (c *ConstantUtf8) TagName() string {
	return "Utf8"
}

type ConstantInteger struct {
}

func (c *ConstantInteger) TagValue() uint8 {
	return 3
}

func (c *ConstantInteger) TagName() string {
	return "Integer"
}

type ConstantFloat struct {
}

func (c *ConstantFloat) TagValue() uint8 {
	return 4
}

func (c *ConstantFloat) TagName() string {
	return "Float"
}

type ConstantLong struct {
}

func (c *ConstantLong) TagValue() uint8 {
	return 5
}

func (c *ConstantLong) TagName() string {
	return "Long"
}

type ConstantDouble struct {
}

func (c *ConstantDouble) TagValue() uint8 {
	return 6
}

func (c *ConstantDouble) TagName() string {
	return "Double"
}

type ConstantClass struct {
}

func (c *ConstantClass) TagValue() uint8 {
	return 7
}

func (c *ConstantClass) TagName() string {
	return "Class"
}

type ConstantString struct {
}

func (c *ConstantString) TagValue() uint8 {
	return 8
}

func (c *ConstantString) TagName() string {
	return "String"
}

type ConstantFieldref struct {
}

func (c *ConstantFieldref) TagValue() uint8 {
	return 9
}

func (c *ConstantFieldref) TagName() string {
	return "Fieldref"
}

type ConstantMethodref struct {
}

func (c *ConstantMethodref) TagValue() uint8 {
	return 10
}

func (c *ConstantMethodref) TagName() string {
	return "Methodref"
}

type ConstantInterfaceMethodref struct {
}

func (c *ConstantInterfaceMethodref) TagValue() uint8 {
	return 11
}

func (c *ConstantInterfaceMethodref) TagName() string {
	return "InterfaceMethodref"
}

type ConstantNameAndType struct {
}

func (c *ConstantNameAndType) TagValue() uint8 {
	return 12
}

func (c *ConstantNameAndType) TagName() string {
	return "NameAndType"
}

type ConstantMethodHandle struct {
}

func (c *ConstantMethodHandle) TagValue() uint8 {
	return 15
}

func (c *ConstantMethodHandle) TagName() string {
	return "MethodHandle"
}

type ConstantMethodType struct {
}

func (c *ConstantMethodType) TagValue() uint8 {
	return 16
}

func (c *ConstantMethodType) TagName() string {
	return "MethodType"
}

type ConstantDynamic struct {
}

func (c *ConstantDynamic) TagValue() uint8 {
	return 17
}

func (c *ConstantDynamic) TagName() string {
	return "Dynamic"
}

type ConstantInvokeDynamic struct {
}

func (c *ConstantInvokeDynamic) TagValue() uint8 {
	return 18
}

func (c *ConstantInvokeDynamic) TagName() string {
	return "InvokeDynamic"
}

type ConstantModule struct {
}

func (c *ConstantModule) TagValue() uint8 {
	return 19
}

func (c *ConstantModule) TagName() string {
	return "Module"
}

type ConstantPackage struct {
}

func (c *ConstantPackage) TagValue() uint8 {
	return 20
}

func (c *ConstantPackage) TagName() string {
	return "Package"
}
