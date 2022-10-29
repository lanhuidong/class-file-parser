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
		return fmt.Sprintf("JDK Version 1.1, mimorVersion: %d", v.MinorVersion)
	} else if v.MajorVersion > 45 && v.MinorVersion == 0 {
		jdkVersion := v.MajorVersion - 44
		return fmt.Sprintf("JDK Version 1.%d (LTS), mimorVersion: %d", jdkVersion, v.MinorVersion)
	} else if v.MajorVersion > 52 && v.MinorVersion == 0 {
		jdkVersion := v.MajorVersion - 44
		if jdkVersion == 8 || jdkVersion == 11 || jdkVersion == 17 {
			return fmt.Sprintf("JDK Version %d (LTS), mimorVersion: %d", jdkVersion, v.MinorVersion)
		} else {
			return fmt.Sprintf("JDK Version %d, mimorVersion: %d", jdkVersion, v.MinorVersion)
		}
	}
	return "Unknow JDK Version"
}
