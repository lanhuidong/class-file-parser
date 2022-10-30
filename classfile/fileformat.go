package classfile

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Version struct {
	MajorVersion uint16
	MinorVersion uint16
}

func (v *Version) Parse(data []byte) {
	binary.Read(bytes.NewBuffer(data[0:2]), binary.BigEndian, &v.MinorVersion)
	binary.Read(bytes.NewBuffer(data[2:4]), binary.BigEndian, &v.MajorVersion)
}

func (v *Version) String() string {
	if v.MajorVersion == 45 {
		return fmt.Sprintf("JDK Version 1.1, %d.%d", v.MajorVersion, v.MinorVersion)
	} else if v.MajorVersion > 52 && v.MinorVersion == 0 {
		jdkVersion := v.MajorVersion - 44
		if jdkVersion == 8 || jdkVersion == 11 || jdkVersion == 17 {
			return fmt.Sprintf("JDK Version %d (LTS), %d.%d", jdkVersion, v.MajorVersion, v.MinorVersion)
		} else {
			return fmt.Sprintf("JDK Version %d, %d.%d", jdkVersion, v.MajorVersion, v.MinorVersion)
		}
	} else if v.MajorVersion > 45 && v.MinorVersion == 0 {
		jdkVersion := v.MajorVersion - 44
		return fmt.Sprintf("JDK Version 1.%d (LTS), %d.%d", jdkVersion, v.MajorVersion, v.MinorVersion)
	}
	return "Unknow JDK Version"
}

type Utf8 struct {
	Tag    uint8
	Length uint16
	Value  []byte
}

func (u *Utf8) Parse(data []byte) {
	binary.Read(bytes.NewBuffer(data[0:1]), binary.BigEndian, &u.Tag)
	binary.Read(bytes.NewBuffer(data[1:3]), binary.BigEndian, &u.Length)
}

func (u *Utf8) String() string {
	return fmt.Sprintf("utf8: %s", string(u.Value))
}

type Methodref struct {
	Tag              uint8
	ClassIndex       uint16
	NameAndTypeIndex uint16
}

func (m *Methodref) Parse(data []byte) {
	binary.Read(bytes.NewBuffer(data[0:1]), binary.BigEndian, &m.Tag)
	binary.Read(bytes.NewBuffer(data[1:3]), binary.BigEndian, &m.ClassIndex)
	binary.Read(bytes.NewBuffer(data[3:5]), binary.BigEndian, &m.NameAndTypeIndex)
}

func (m *Methodref) String() string {
	return fmt.Sprintf("class index: %d, name and type index: %d", m.ClassIndex, m.NameAndTypeIndex)
}

type Class struct {
	Tag   uint8
	Index uint16
}

func (c *Class) Parse(data []byte) {
	binary.Read(bytes.NewBuffer(data[0:1]), binary.BigEndian, &c.Tag)
	binary.Read(bytes.NewBuffer(data[1:3]), binary.BigEndian, &c.Index)
}

func (c *Class) String() string {
	return fmt.Sprintf("class name is at index: %d", c.Index)
}

type Integer struct {
	Tag   uint8
	Value int32
}

func (i *Integer) Parse(data []byte) {
	binary.Read(bytes.NewBuffer(data[0:1]), binary.BigEndian, &i.Tag)
	binary.Read(bytes.NewBuffer(data[1:5]), binary.BigEndian, &i.Value)
}

func (i *Integer) String() string {
	return fmt.Sprintf("integer value is: %d", i.Value)
}

type Long struct {
	Tag   uint8
	Value int64
}

func (l *Long) Parse(data []byte) {
	binary.Read(bytes.NewBuffer(data[0:1]), binary.BigEndian, &l.Tag)
	binary.Read(bytes.NewBuffer(data[1:9]), binary.BigEndian, &l.Value)
}

func (l *Long) String() string {
	return fmt.Sprintf("long value is: %d", l.Value)
}

type NameAndType struct {
	Tag       uint8
	NameIndex uint16
	DescIndex uint16
}

func (n *NameAndType) Parse(data []byte) {
	binary.Read(bytes.NewBuffer(data[0:1]), binary.BigEndian, &n.Tag)
	binary.Read(bytes.NewBuffer(data[1:3]), binary.BigEndian, &n.NameIndex)
	binary.Read(bytes.NewBuffer(data[3:5]), binary.BigEndian, &n.DescIndex)
}

func (n *NameAndType) String() string {
	return fmt.Sprintf("name is at index: %d, desc is at index %d", n.NameIndex, n.DescIndex)
}
