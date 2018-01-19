package main

import (
	"encoding/binary"
	"fmt"
	"math"
)

type Amf0Object struct {
	propertyName string
	value        interface{}
}

type Amf0EcmaArray struct {
	count     uint32
	anyObject []Amf0Object
}

type Amf0StrictArray struct {
	count     uint32
	anyObject []interface{}
}

//this function do not affect the offset parsed in data.
func Amf0ObjectEof(data []uint8, offset *uint32) (res bool) {
	if len(data) < 3 {
		res = false
		return
	}

	objectEofFlag := uint32(data[*offset])<<16 + uint32(data[*offset+1])<<8 + uint32(data[*offset+2])
	if RTMP_AMF0_ObjectEnd == objectEofFlag {
		res = true
		*offset += 3
	}

	return
}

func Amf0ReadUtf8(data []uint8, offset *uint32) (err error, value string) {
	if (uint32(len(data)) - *offset) < 2 {
		err = fmt.Errorf("Amf0ReadUtf8: 1, data len is not enough")
		return
	}

	dataLen := binary.BigEndian.Uint16(data[*offset : *offset+2])
	*offset += 2

	if (uint32(len(data)) - *offset) < uint32(dataLen) {
		err = fmt.Errorf("Amf0ReadUtf8: 2, data len is not enough")
		return
	}

	if 0 == dataLen {
		return
	}

	value = string(data[*offset : *offset+uint32(dataLen)])
	*offset += uint32(dataLen)

	return
}

func Amf0ReadAny(data []uint8, offset *uint32) (err error, value interface{}) {

	if Amf0ObjectEof(data, offset) {
		return
	}

	if (uint32(len(data)) - *offset) < 1 {
		err = fmt.Errorf("Amf0ReadAny: 0, data len is not enough")
		return
	}

	marker := data[*offset]

	switch marker {
	case RTMP_AMF0_String:
		err, value = Amf0ReadString(data, offset)
	case RTMP_AMF0_Boolean:
		err, value = Amf0ReadBool(data, offset)
	case RTMP_AMF0_Number:
		err, value = Amf0ReadNumber(data, offset)
	case RTMP_AMF0_Null:
		err = Amf0ReadNull(data, offset)
	case RTMP_AMF0_Undefined:
		err = Amf0ReadUndefined(data, offset)
	case RTMP_AMF0_Object:
		err, value = Amf0ReadObjects(data, offset)
	case RTMP_AMF0_LongString:
		err, value = Amf0ReadLongString(data, offset)
	case RTMP_AMF0_EcmaArray:
		err, value = Amf0ReadEcmaArray(data, offset)
	case RTMP_AMF0_StrictArray:
		err, value = Amf0ReadStrictArray(data, offset)
	default:
		err = fmt.Errorf("Amf0ReadAny: unknown marker value, marker=", marker)
	}

	if err != nil {
		return
	}

	return
}

func Amf0ReadString(data []uint8, offset *uint32) (err error, value string) {

	if (uint32(len(data)) - *offset) < 1 {
		err = fmt.Errorf("Amf0ReadString: 0, data len is not enough")
		return
	}

	marker := data[*offset]
	*offset += 1

	if RTMP_AMF0_String != marker {
		err = fmt.Errorf("Amf0ReadString: RTMP_AMF0_String != marker")
		return
	}

	err, value = Amf0ReadUtf8(data, offset)

	return
}

func Amf0ReadNumber(data []uint8, offset *uint32) (err error, value float64) {

	if (uint32(len(data)) - *offset) < 1 {
		err = fmt.Errorf("Amf0ReadNumber: 0, data len is not enough")
		return
	}

	marker := data[*offset]
	*offset += 1

	if RTMP_AMF0_Number != marker {
		err = fmt.Errorf("Amf0ReadNumber: RTMP_AMF0_Number != marker")
		return
	}

	if (uint32(len(data)) - *offset) < 8 {
		err = fmt.Errorf("Amf0ReadNumber: 1, data len is not enough")
		return
	}

	value_tmp := binary.BigEndian.Uint64(data[*offset : *offset+8])
	*offset += 8

	value = math.Float64frombits(value_tmp)

	return
}

func Amf0ReadObjects(data []uint8, offset *uint32) (err error, amf0objects []Amf0Object) {

	if (uint32(len(data)) - *offset) < 1 {
		err = fmt.Errorf("Amf0ReadObjects: 0, data len is not enough")
		return
	}

	marker := data[*offset]
	*offset += 1

	if RTMP_AMF0_Object != marker {
		err = fmt.Errorf("error: Amf0ReadObjects:RTMP_AMF0_Object != marker")
		return
	}

	for {
		if *offset >= uint32(len(data)) {
			break
		}

		if Amf0ObjectEof(data, offset) {
			break
		}

		var amf0object Amf0Object

		err, amf0object.propertyName = Amf0ReadUtf8(data, offset)
		if err != nil {
			break
		}

		err, amf0object.value = Amf0ReadAny(data, offset)
		if err != nil {
			break
		}

		amf0objects = append(amf0objects, amf0object)
	}

	return
}

func Amf0ReadBool(data []uint8, offset *uint32) (err error, value bool) {
	if (uint32(len(data)) - *offset) < 1 {
		err = fmt.Errorf("Amf0ReadBool:  0, data len is not enough")
		return
	}

	marker := data[*offset]
	*offset += 1

	if RTMP_AMF0_Boolean != marker {
		err = fmt.Errorf("Amf0ReadBool: RTMP_AMF0_Boolean != marker")
		return
	}

	if (uint32(len(data)) - *offset) < 1 {
		err = fmt.Errorf("Amf0ReadBool:  1, data len is not enough")
		return
	}

	v := data[*offset]
	*offset += 1

	if v != 0 {
		value = true
	} else {
		value = false
	}

	return
}

func Amf0ReadNull(data []uint8, offset *uint32) (err error) {
	if (uint32(len(data)) - *offset) < 1 {
		err = fmt.Errorf("Amf0ReadNull:  0, data len is not enough")
		return
	}

	marker := data[*offset]
	*offset += 1

	if RTMP_AMF0_Null != marker {
		err = fmt.Errorf("Amf0ReadNull: RTMP_AMF0_Null != marker")
		return
	}

	return
}

func Amf0ReadUndefined(data []uint8, offset *uint32) (err error) {
	if (uint32(len(data)) - *offset) < 1 {
		err = fmt.Errorf("Amf0ReadUndefined:  0, data len is not enough")
		return
	}

	marker := data[*offset]
	*offset += 1

	if RTMP_AMF0_Undefined != marker {
		err = fmt.Errorf("Amf0ReadUndefined: RTMP_AMF0_Undefined != marker")
		return
	}

	return
}

func Amf0ReadLongString(data []uint8, offset *uint32) (err error, value string) {

	if (uint32(len(data)) - *offset) < 1 {
		err = fmt.Errorf("Amf0ReadLongString: 0, data len is not enough")
		return
	}

	marker := data[*offset]
	*offset += 1

	if RTMP_AMF0_LongString != marker {
		err = fmt.Errorf("Amf0ReadLongString: RTMP_AMF0_LongString != marker")
		return
	}

	if (uint32(len(data)) - *offset) < 4 {
		err = fmt.Errorf("Amf0ReadLongString: 1, data len is not enough")

		return
	}

	dataLen := binary.BigEndian.Uint32(data[*offset : *offset+4])
	*offset += 4
	if dataLen <= 0 {
		err = fmt.Errorf("Amf0ReadLongString: data len is <= 0, dataLen=", dataLen)
		return
	}

	if (uint32(len(data)) - *offset) < dataLen {
		err = fmt.Errorf("Amf0ReadLongString: 2, data len is not enough")
		return
	}

	value = string(data[*offset : *offset+dataLen])
	*offset += dataLen

	return
}

func Amf0ReadEcmaArray(data []uint8, offset *uint32) (err error, value Amf0EcmaArray) {
	if (uint32(len(data)) - *offset) < 1 {
		err = fmt.Errorf("Amf0ReadEcmaArray: 0, data len is not enough")
		return
	}

	marker := data[*offset]
	*offset += 1

	if RTMP_AMF0_EcmaArray != marker {
		err = fmt.Errorf("error: Amf0ReadEcmaArray: RTMP_AMF0_EcmaArray != marker")
		return
	}

	if (uint32(len(data)) - *offset) < 4 {
		err = fmt.Errorf("Amf0ReadEcmaArray: 1, data len is not enough")
		return
	}

	value.count = binary.BigEndian.Uint32(data[*offset : *offset+4])
	*offset += 4

	for {
		if *offset >= uint32(len(data)) {
			break
		}

		if Amf0ObjectEof(data, offset) {
			break
		}

		var amf Amf0Object
		err, amf.propertyName = Amf0ReadUtf8(data, offset)
		if err != nil {
			break
		}

		err, amf.value = Amf0ReadAny(data, offset)
		if err != nil {
			break
		}

		value.anyObject = append(value.anyObject, amf)
	}

	return
}

func Amf0ReadStrictArray(data []uint8, offset *uint32) (err error, value Amf0StrictArray) {
	if (uint32(len(data)) - *offset) < 1 {
		err = fmt.Errorf("Amf0ReadStrictArray: 0, data len is not enough")
		return
	}

	marker := data[*offset]
	*offset += 1

	if RTMP_AMF0_StrictArray != marker {
		err = fmt.Errorf("Amf0ReadStrictArray: error: RTMP_AMF0_StrictArray != marker")
		return
	}

	if (uint32(len(data)) - *offset) < 4 {
		err = fmt.Errorf("Amf0ReadStrictArray: 1, data len is not enough")
		return
	}

	value.count = binary.BigEndian.Uint32(data[*offset : *offset+4])
	*offset += 4

	for i := 0; uint32(i) < value.count; i++ {
		if *offset >= uint32(len(data)) {
			break
		}

		var obj interface{}

		err, obj = Amf0ReadAny(data, offset)
		if err != nil {
			break
		}

		value.anyObject = append(value.anyObject, obj)
	}

	return
}

func Amf0Discovery(data []uint8, offset *uint32) (err error, value interface{}, marker uint8) {

	if Amf0ObjectEof(data, offset) {
		return
	}

	if (uint32(len(data)) - *offset) < 1 {
		err = fmt.Errorf("Amf0Discovery: 0, data len is not enough")
		return
	}

	marker = data[*offset]

	switch marker {
	case RTMP_AMF0_String:
		err, value = Amf0ReadString(data, offset)
	case RTMP_AMF0_Boolean:
		err, value = Amf0ReadBool(data, offset)
	case RTMP_AMF0_Number:
		err, value = Amf0ReadNumber(data, offset)
	case RTMP_AMF0_Null:
		err = Amf0ReadNull(data, offset)
	case RTMP_AMF0_Undefined:
		err = Amf0ReadUndefined(data, offset)
	case RTMP_AMF0_Object:
		err, value = Amf0ReadObjects(data, offset)
	case RTMP_AMF0_LongString:
		err, value = Amf0ReadLongString(data, offset)
	case RTMP_AMF0_EcmaArray:
		err, value = Amf0ReadEcmaArray(data, offset)
	case RTMP_AMF0_StrictArray:
		err, value = Amf0ReadStrictArray(data, offset)
	default:
		err = fmt.Errorf("Amf0Discovery: unknown marker type, marker=", marker)
	}

	if err != nil {
		return
	}

	return
}