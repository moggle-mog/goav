package table

import (
	"bytes"
	"encoding/binary"
)

// Descriptor 描述表
type Descriptor struct {
	data *bytes.Buffer
}

// NewDescriptor 描述表
func NewDescriptor() *Descriptor {
	return &Descriptor{
		data: bytes.NewBuffer(nil),
	}
}

// GetBuffer is the buffer of descriptor
func (d *Descriptor) GetBuffer() *bytes.Buffer {
	return d.data
}

// Service serviceType: pmt表的program pid
func (d *Descriptor) Service(serviceType byte, serviceProviderName string, serviceName string) error {
	serviceProviderNameLen := byte(len(serviceProviderName))

	serviceNameLen := byte(len(serviceName))
	descriptorLen := serviceProviderNameLen + serviceNameLen + 3

	_, err := d.data.Write([]byte{0x48, descriptorLen, serviceType})
	if err != nil {
		return err
	}

	err = d.data.WriteByte(serviceProviderNameLen)
	if err != nil {
		return err
	}

	_, err = d.data.WriteString(serviceProviderName)
	if err != nil {
		return err
	}

	err = d.data.WriteByte(serviceNameLen)
	if err != nil {
		return err
	}

	_, err = d.data.WriteString(serviceName)
	if err != nil {
		return err
	}

	return nil
}

// NetworkName 网络名称
func (d *Descriptor) NetworkName(name string) error {
	nameLen := len(name)
	descriptorLen := byte(nameLen)

	_, err := d.data.Write([]byte{0x40, descriptorLen})
	if err != nil {
		return err
	}

	_, err = d.data.WriteString(name)
	if err != nil {
		return err
	}

	return nil
}

// CountryAvailability 区域可用性
func (d *Descriptor) CountryAvailability(isAvailable bool, countryCode uint32) error {
	var countryAvailabilityFlag byte
	if isAvailable {
		countryAvailabilityFlag = 0x80
	}

	ret := [7]byte{0x49, 4, countryAvailabilityFlag}
	binary.BigEndian.PutUint32(ret[3:], countryCode&0x0FFF)

	// 写入
	_, err := d.data.Write(ret[:])
	if err != nil {
		return err
	}

	return nil
}
