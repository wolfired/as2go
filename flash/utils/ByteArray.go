package utils

import (
	"encoding/binary"
	"math"

	"github.com/wolfired/as2go/flash/errors"
)

const (
	pointerPosition uint = 0x1
	pointerLength   uint = 0x2
	pointerCapacity uint = 0x4

	byteWide1 uint = 1
	byteWide2 uint = 2
	byteWide4 uint = 4
	byteWide8 uint = 8
)

/*
NewByteArray 创建一个 ByteArray 指针.
*/
func NewByteArray(raw []byte) (b *ByteArray) {
	b = &ByteArray{}

	if nil == raw {
		b.raw = make([]byte, 0)
	} else {
		b.raw = raw
	}

	b.SetEndian(EndianBig)
	b.SetPosition(0)
	b.SetLength(0)

	return
}

/*
ByteArray 提供用于优化读取/写入以及处理二进制数据的方法和属性.
*/
type ByteArray struct {
	raw      []byte
	endian   binary.ByteOrder
	position uint
	length   uint
}

/*
GetEndian 获取 ByteArray 的字节序.
*/
func (b *ByteArray) GetEndian() uint {
	if binary.LittleEndian == b.endian {
		return EndianLittle
	}
	return EndianBig
}

/*
SetEndian 设置 ByteArray 的字节序.
*/
func (b *ByteArray) SetEndian(endian uint) {
	if EndianLittle == endian {
		b.endian = binary.LittleEndian
		return
	}

	b.endian = binary.BigEndian
}

/*
GetLength 获取 ByteArray 的长度
*/
func (b *ByteArray) GetLength() uint {
	return b.length
}

/*
SetLength 设置 ByteArray 的长度
*/
func (b *ByteArray) SetLength(newLen uint) {
	if b.length == newLen {
		return
	}

	if b.length < newLen {
		b.checkCapacity(newLen)
	}

	b.length = newLen

	if b.position > b.length {
		b.SetPosition(b.length)
	}
}

/*
GetPosition 获取 ByteArray 的读写位置
*/
func (b *ByteArray) GetPosition() uint {
	return b.position
}

/*
SetPosition 设置 ByteArray 的读写位置
*/
func (b *ByteArray) SetPosition(newPos uint) {
	if b.position == newPos {
		return
	}

	b.position = newPos
}

/*
BytesAvailable ByteArray 的剩余可读取字节数
*/
func (b *ByteArray) BytesAvailable() uint {
	if b.length > b.position {
		return b.length - b.position
	}

	return 0
}

/*
Clear 清空 ByteArray 的数据
*/
func (b *ByteArray) Clear() {
	b.SetPosition(0)
	b.SetLength(0)
}

/*
ReadBoolean 从字节流中读取布尔值.
读取单个字节, 如果字节非零, 则返回 true, 否则返回 false.
*/
func (b *ByteArray) ReadBoolean() (bool, error) {
	err := b.checkLength(byteWide1)

	if nil != err {
		return false, err
	}

	value := b.raw[b.position]

	b.movePointer(byteWide1, pointerPosition)

	return 0 != value, nil
}

/*
ReadByte 从字节流中读取带符号的字节.
返回值的范围是从 -128 到 127.
*/
func (b *ByteArray) ReadByte() (int8, error) {
	err := b.checkLength(byteWide1)

	if nil != err {
		return 0, err
	}

	value := b.raw[b.position]

	b.movePointer(byteWide1, pointerPosition)

	return int8(value), nil
}

/*
ReadBytes 从字节流中读取 length 参数指定的数据字节数.
从 offset 指定的位置开始, 将字节读入 bytes 参数指定的 ByteArray 对象中, 并将字节写入目标 ByteArray 中.
*/
func (b *ByteArray) ReadBytes(bytes *ByteArray, offset uint, length uint) error {
	if 0 == length {
		length = b.BytesAvailable()

		if 0 == length {
			return nil
		}
	} else {
		err := b.checkLength(length)

		if nil != err {
			return err
		}
	}

	if math.MaxUint32 < offset+length {
		return errors.ErrorRange
	}

	bytes.checkCapacity(offset + length)

	copy(bytes.raw[offset:offset+length], b.raw[b.position:b.position+length])

	if offset+length > bytes.length {
		bytes.movePointer(offset+length-bytes.length, pointerLength)
	}

	b.movePointer(length, pointerPosition)

	return nil
}

/*
ReadDouble 从字节流中读取一个 IEEE 754 双精度(64 位)浮点数.
*/
func (b *ByteArray) ReadDouble() (float64, error) {
	err := b.checkLength(byteWide8)

	if nil != err {
		return 0.0, err
	}

	value := math.Float64frombits(b.endian.Uint64(b.raw[b.position:]))

	b.movePointer(byteWide8, pointerPosition)

	return value, nil
}

/*
ReadFloat 从字节流中读取一个 IEEE 754 单精度(32 位)浮点数.
*/
func (b *ByteArray) ReadFloat() (float32, error) {
	err := b.checkLength(byteWide4)

	if nil != err {
		return 0.0, err
	}

	value := math.Float32frombits(b.endian.Uint32(b.raw[b.position:]))

	b.movePointer(byteWide4, pointerPosition)

	return value, nil
}

/*
ReadInt 从字节流中读取一个带符号的 32 位整数.
返回值的范围是从 -2147483648 到 2147483647.
*/
func (b *ByteArray) ReadInt() (int, error) {
	err := b.checkLength(byteWide4)

	if nil != err {
		return 0, err
	}

	value := int(b.endian.Uint32(b.raw[b.position:]))

	b.movePointer(byteWide4, pointerPosition)

	return value, nil
}

/*
ReadShort 从字节流中读取一个带符号的 16 位整数.
返回值的范围是从 -32768 到 32767.
*/
func (b *ByteArray) ReadShort() (int16, error) {
	err := b.checkLength(byteWide2)

	if nil != err {
		return 0, err
	}

	value := int16(b.endian.Uint16(b.raw[b.position:]))

	b.movePointer(byteWide2, pointerPosition)

	return value, nil
}

/*
ReadUnsignedByte 从字节流中读取无符号的字节.
返回值的范围是从 0 到 255.
*/
func (b *ByteArray) ReadUnsignedByte() (uint8, error) {
	err := b.checkLength(byteWide1)

	if nil != err {
		return 0, err
	}

	value := b.raw[b.position]

	b.movePointer(byteWide1, pointerPosition)

	return uint8(value), nil
}

/*
ReadUnsignedInt 从字节流中读取一个无符号的 32 位整数.
返回值的范围是从 0 到 4294967295.
*/
func (b *ByteArray) ReadUnsignedInt() (uint, error) {
	err := b.checkLength(byteWide4)

	if nil != err {
		return 0, err
	}

	value := b.endian.Uint32(b.raw[b.position:])

	b.movePointer(byteWide4, pointerPosition)

	return uint(value), nil
}

/*
ReadUnsignedShort 从字节流中读取一个无符号的 16 位整数.
返回值的范围是从 0 到 65535.
*/
func (b *ByteArray) ReadUnsignedShort() (uint16, error) {
	err := b.checkLength(byteWide2)

	if nil != err {
		return 0, err
	}

	value := b.endian.Uint16(b.raw[b.position:])

	b.movePointer(byteWide2, pointerPosition)

	return value, nil
}

/*
ReadUTF 从字节流中读取一个 UTF-8 字符串. 假定字符串的前缀是无符号的短整型(以字节表示长度).
*/
func (b *ByteArray) ReadUTF() (string, error) {
	length, _ := b.ReadUnsignedShort()
	return b.ReadUTFBytes(length)
}

/*
ReadUTFBytes 从字节流中读取一个由 length 参数指定的 UTF-8 字节序列, 并返回一个字符串.
*/
func (b *ByteArray) ReadUTFBytes(length uint16) (string, error) {
	err := b.checkLength(uint(length))

	if nil != err {
		return "", err
	}

	str := string(b.raw[b.position : b.position+uint(length)])

	b.movePointer(uint(length), pointerPosition)

	return str, nil
}

/*
WriteBoolean 写入布尔值. 根据 value 参数写入单个字节. 如果为 true, 则写入 1, 如果为 false, 则写入 0.
*/
func (b *ByteArray) WriteBoolean(value bool) {
	b.checkCapacity(b.position + byteWide1)

	b.raw[b.position] = 0

	if value {
		b.raw[b.position] = 1
	}

	b.movePointer(byteWide1, pointerPosition|pointerLength)
}

/*
WriteByte 在字节流中写入一个字节.
使用参数的低 8 位. 忽略高 24 位.
*/
func (b *ByteArray) WriteByte(value int8) {
	b.checkCapacity(b.position + byteWide1)

	b.raw[b.position] = byte(value)

	b.movePointer(byteWide1, pointerPosition|pointerLength)
}

/*
WriteBytes 将指定字节数组 bytes(起始偏移量为 offset, 从零开始的索引)中包含 length 个字节的字节序列写入字节流.
如果省略 length 参数, 则使用默认长度 0; 该方法将从 offset 开始写入整个缓冲区. 如果还省略了 offset 参数, 则写入整个缓冲区.
如果 offset 或 length 超出范围, 它们将被锁定到 bytes 数组的开头和结尾.
*/
func (b *ByteArray) WriteBytes(bytes *ByteArray, offset uint, length uint) {
	if bytes.length < offset {
		offset = 0
	}

	if 0 == length || bytes.length < offset+length {
		length = bytes.length - offset
	}

	b.checkCapacity(b.position + length)

	copy(b.raw[b.position:b.position+length], bytes.raw[offset:offset+length])

	b.movePointer(length, pointerPosition)
}

/*
WriteDouble 在字节流中写入一个 IEEE 754 双精度（64 位）浮点数.
*/
func (b *ByteArray) WriteDouble(value float64) {
	b.checkCapacity(b.position + byteWide8)

	b.endian.PutUint64(b.raw[b.position:], math.Float64bits(value))

	b.movePointer(byteWide8, pointerPosition|pointerLength)
}

/*
WriteFloat 在字节流中写入一个 IEEE 754 单精度(32 位)浮点数.
*/
func (b *ByteArray) WriteFloat(value float32) {
	b.checkCapacity(b.position + byteWide4)

	b.endian.PutUint32(b.raw[b.position:], math.Float32bits(value))

	b.movePointer(byteWide4, pointerPosition|pointerLength)
}

/*
WriteInt 在字节流中写入一个带符号的 32 位整数.
*/
func (b *ByteArray) WriteInt(value int32) {
	b.checkCapacity(b.position + byteWide4)

	b.endian.PutUint32(b.raw[b.position:], uint32(value))

	b.movePointer(byteWide4, pointerPosition|pointerLength)
}

/*
WriteShort 在字节流中写入一个 16 位整数. 使用参数的低 16 位. 忽略高 16 位.
*/
func (b *ByteArray) WriteShort(value int16) {
	b.checkCapacity(b.position + byteWide2)

	b.endian.PutUint16(b.raw[b.position:], uint16(value))

	b.movePointer(byteWide2, pointerPosition|pointerLength)
}

/*
WriteUnsignedInt 在字节流中写入一个无符号的 32 位整数.
*/
func (b *ByteArray) WriteUnsignedInt(value uint32) {
	b.checkCapacity(b.position + byteWide4)

	b.endian.PutUint32(b.raw[b.position:], value)

	b.movePointer(byteWide4, pointerPosition|pointerLength)
}

/*
WriteUTF 将 UTF-8 字符串写入字节流. 先写入以字节表示的 UTF-8 字符串长度(作为 16 位整数), 然后写入表示字符串字符的字节.
*/
func (b *ByteArray) WriteUTF(value string) {
	bs := []byte(value)
	length := uint(len(bs))

	b.WriteShort(int16(length))

	b.checkCapacity(b.position + length)

	copy(b.raw[b.position:], bs)

	b.movePointer(length, pointerPosition|pointerLength)
}

/*
WriteUTFBytes 将 UTF-8 字符串写入字节流. 类似于 writeUTF() 方法，但 writeUTFBytes() 不使用 16 位长度的词为字符串添加前缀.
*/
func (b *ByteArray) WriteUTFBytes(value string) {
	bs := []byte(value)
	length := uint(len(bs))

	b.checkCapacity(b.position + length)

	copy(b.raw[b.position:], bs)

	b.movePointer(length, pointerPosition|pointerLength)
}

func (b *ByteArray) checkLength(needBytes uint) error {
	if b.BytesAvailable() < needBytes {
		return errors.ErrorEOF
	}

	return nil
}

func (b *ByteArray) checkCapacity(newCap uint) {
	oldCap := uint(len(b.raw))

	if 0 == oldCap {
		b.raw = make([]byte, newCap)
	} else if oldCap < newCap {
		oldCap += oldCap

		for oldCap < newCap {
			oldCap += oldCap
		}

		oldRaw := b.raw[:b.length]
		b.raw = make([]byte, oldCap)
		copy(b.raw, oldRaw)
	}
}

func (b *ByteArray) movePointer(moveBytes uint, pointerType uint) {
	newPos := b.position + moveBytes

	if 0 < pointerPosition&pointerType {
		b.position = newPos
	}

	if 0 < pointerLength&pointerType && b.length < newPos {
		b.length = newPos
	}
}
