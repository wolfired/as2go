package utils

import (
	"bytes"
	"encoding/binary"
	"math"
	"testing"
)

func Test_ByteArray_New(t *testing.T) {
	ba := NewByteArray()

	if 0 != len(ba.raw) {
		t.Error("Expect capacity", 0)
		t.Error("Actual capacity", len(ba.raw))
	}

	if binary.BigEndian != ba.endian {
		t.Error("Expect endian", binary.BigEndian)
		t.Error("Actual endian", ba.endian)
	}

	if 0 != ba.position {
		t.Error("Expect position", 0)
		t.Error("Actual position", ba.position)
	}

	if 0 != ba.length {
		t.Error("Expect length", 0)
		t.Error("Actual length", ba.length)
	}
}

func Test_ByteArray_SetLength(t *testing.T) {
	{
		ba := NewByteArray()

		ba.SetLength(4)

		if 4 != ba.length {
			t.Error("Expect length", 4)
			t.Error("Actual length", ba.length)
		}

		if 4 != len(ba.raw) {
			t.Error("Expect raw length", 4)
			t.Error("Actual raw length", len(ba.raw))
		}

		ba.SetLength(8)

		if 8 != ba.length {
			t.Error("Expect length", 8)
			t.Error("Actual length", ba.length)
		}

		if 8 != len(ba.raw) {
			t.Error("Expect raw length", 8)
			t.Error("Actual raw length", len(ba.raw))
		}

		ba.SetLength(13)

		if 13 != ba.length {
			t.Error("Expect length", 13)
			t.Error("Actual length", ba.length)
		}

		if 16 != len(ba.raw) {
			t.Error("Expect raw length", 16)
			t.Error("Actual raw length", len(ba.raw))
		}
	}

	{
		ba := NewByteArray()

		ba.SetLength(4)

		ba.raw[0] = 0x1
		ba.raw[1] = 0x2
		ba.raw[2] = 0x3
		ba.raw[3] = 0x4

		ba.SetLength(8)

		if 0 != bytes.Compare([]byte{0x1, 0x2, 0x3, 0x4, 0x0, 0x0, 0x0, 0x0}, ba.raw) {
			t.Error("Expect raw", []byte{0x1, 0x2, 0x3, 0x4, 0x0, 0x0, 0x0, 0x0})
			t.Error("Actual raw", ba.raw)
		}
	}
}

func Test_ByteArray_BytesAvailable(t *testing.T) {
	{
		ba := NewByteArray()

		ba.SetLength(4)

		if 4 != ba.BytesAvailable() {
			t.Error("Expect bytes available", 4)
			t.Error("Actual bytes available", ba.BytesAvailable())
		}

		ba.SetPosition(2)

		if 2 != ba.BytesAvailable() {
			t.Error("Expect bytes available", 2)
			t.Error("Actual bytes available", ba.BytesAvailable())
		}

		ba.SetPosition(4)

		if 0 != ba.BytesAvailable() {
			t.Error("Expect bytes available", 0)
			t.Error("Actual bytes available", ba.BytesAvailable())
		}
	}
}

func Test_ByteArray_Clear(t *testing.T) {
	{
		ba := NewByteArray()

		ba.SetLength(4)
		ba.SetPosition(4)

		ba.Clear()

		if 0 != ba.length {
			t.Error("Expect length", 0)
			t.Error("Actual length", ba.length)
		}

		if 0 != ba.position {
			t.Error("Expect position", 0)
			t.Error("Actual position", ba.position)
		}

		if 4 != len(ba.raw) {
			t.Error("Expect raw ength", 4)
			t.Error("Actual raw length", len(ba.raw))
		}
	}
}

func Test_ByteArray_ReadBoolean(t *testing.T) {
	{
		ba := NewByteArray()

		ba.SetLength(4)

		ba.raw[0] = 0
		ba.raw[1] = 1

		v, _ := ba.ReadBoolean()

		if false != v {
			t.Error("Expect", false)
			t.Error("Actual ", v)
		}

		v, _ = ba.ReadBoolean()

		if true != v {
			t.Error("Expect", true)
			t.Error("Actual ", v)
		}
	}
}

func Test_ByteArray_WriteBoolean(t *testing.T) {
	{
		ba := NewByteArray()

		ba.WriteBoolean(true)
		ba.WriteBoolean(false)

		ba.SetPosition(0)

		v, _ := ba.ReadBoolean()

		if true != v {
			t.Error("Expect", true)
			t.Error("Actual ", v)
		}

		v, _ = ba.ReadBoolean()

		if false != v {
			t.Error("Expect", false)
			t.Error("Actual ", v)
		}
	}
}

func Test_ByteArray_ReadByte(t *testing.T) {
	{
		expect := [...]int8{math.MinInt8, -1, 0, 1, math.MaxInt8}
		actual := [len(expect)]int8{}

		ba := NewByteArray()

		ba.SetLength(uint(len(expect)))

		for i := 0; i < len(expect); i++ {
			ba.raw[i] = byte(expect[i])
		}
		for i := 0; i < len(expect); i++ {
			a, _ := ba.ReadByte()
			actual[i] = a
		}

		if expect != actual {
			t.Error("Expect", expect)
			t.Error("Actual", actual)
		}
	}
}

func Test_ByteArray_WriteByte(t *testing.T) {
	{
		expect := [...]int8{math.MinInt8, -1, 0, 1, math.MaxInt8}
		actual := [len(expect)]int8{}

		ba := NewByteArray()

		for i := 0; i < len(expect); i++ {
			ba.WriteByte(expect[i])
		}

		ba.SetPosition(0)

		for i := 0; i < len(expect); i++ {
			a, _ := ba.ReadByte()
			actual[i] = a
		}

		if expect != actual {
			t.Error("Expect", expect)
			t.Error("Actual", actual)
		}
	}
}

func Test_ByteArray_ReadWriteBytes(t *testing.T) {
	{
		expect := [...]int8{0x00, 0x01, 0x02, 0x03, 0x04}

		src := NewByteArray()
		dst := NewByteArray()

		for _, v := range expect {
			src.WriteByte(v)
		}

		// src.SetPosition(0)

		src.ReadBytes(dst, 0, 0)

		if 0 != dst.GetLength() {
			t.Error("Expect", 0)
			t.Error("Actual", dst.GetLength())
		}
	}

	{
		expect := [...]int8{0x00, 0x01, 0x02, 0x03, 0x04}
		actual := [len(expect)]int8{}

		src := NewByteArray()
		dst := NewByteArray()

		for _, v := range expect {
			src.WriteByte(v)
		}

		src.SetPosition(0)

		src.ReadBytes(dst, 0, 0)

		for i := 0; i < len(expect); i++ {
			a, _ := dst.ReadByte()
			actual[i] = a
		}

		if actual != expect {
			t.Error("Expect", expect)
			t.Error("Actual", actual)
		}
	}

	{
		expect := [...]int8{0x00, 0x01}
		actual := [len(expect)]int8{}

		src := NewByteArray()
		dst := NewByteArray()

		for _, v := range expect {
			src.WriteByte(v)
		}

		src.SetPosition(1)

		src.ReadBytes(dst, 1, 1)

		for i := 0; i < len(expect); i++ {
			a, _ := dst.ReadByte()
			actual[i] = a
		}

		if actual != expect {
			t.Error("Expect", expect)
			t.Error("Actual", actual)
		}
	}
}

func Test_ByteArray_ReadWriteDouble(t *testing.T) {
	ba := NewByteArray()

	expect := [...]float64{math.SmallestNonzeroFloat64, -3.14, 0, 3.14, math.MaxFloat64}

	for _, v := range expect {
		ba.WriteDouble(v)
	}

	ba.SetPosition(0)

	actual := [len(expect)]float64{}

	for i := 0; i < len(expect); i++ {
		a, _ := ba.ReadDouble()
		actual[i] = a
	}

	if actual != expect {
		t.Error("Expect", expect)
		t.Error("Actual", actual)
	}
}

func Test_ByteArray_ReadWriteFloat(t *testing.T) {
	ba := NewByteArray()

	expect := [...]float32{math.SmallestNonzeroFloat32, -3.14, 0, 3.14, math.MaxFloat32}

	for _, v := range expect {
		ba.WriteFloat(v)
	}

	ba.SetPosition(0)

	actual := [len(expect)]float32{}

	for i := 0; i < len(expect); i++ {
		a, _ := ba.ReadFloat()
		actual[i] = a
	}

	if actual != expect {
		t.Error("Expect", expect)
		t.Error("Actual", actual)
	}
}

func Test_ByteArray_ReadWriteUTF(t *testing.T) {
	ba := NewByteArray()

	expect := [...]string{"你好，我叫DayDayUp。", "你好，我是新来的犀利哥。"}

	for _, v := range expect {
		ba.WriteUTF(v)
	}

	ba.SetPosition(0)

	actual := [len(expect)]string{}

	for i := 0; i < len(expect); i++ {
		a, _ := ba.ReadUTF()
		actual[i] = a
	}

	if actual != expect {
		t.Error("Expect", expect)
		t.Error("Actual", actual)
	}
}
