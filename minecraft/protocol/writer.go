package protocol

import (
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/google/uuid"
	"github.com/eren5960/gophertunnel408/minecraft/nbt"
	"image/color"
	"io"
	"reflect"
)

// Writer implements writing methods for data types from Minecraft packets. Each Packet implementation has one
// passed to it when writing.
// Writer implements methods where values are passed using a pointer, so that Reader and Writer have a
// synonymous interface and both implement the IO interface.
type Writer struct {
	w interface {
		io.Writer
		io.ByteWriter
	}
}

// NewWriter creates a new initialised Writer with an underlying io.ByteWriter to write to.
func NewWriter(w interface {
	io.Writer
	io.ByteWriter
}) *Writer {
	return &Writer{w: w}
}

// Uint8 writes a uint8 to the underlying buffer.
func (w *Writer) Uint8(x *uint8) {
	_ = w.w.WriteByte(*x)
}

// Bool writes a bool as either 0 or 1 to the underlying buffer.
func (w *Writer) Bool(x *bool) {
	if *x {
		_ = w.w.WriteByte(1)
		return
	}
	_ = w.w.WriteByte(0)
}

// String writes a string, prefixed with a varuint32, to the underlying buffer.
func (w *Writer) String(x *string) {
	l := uint32(len(*x))
	w.Varuint32(&l)
	_, _ = w.w.Write([]byte(*x))
}

// ByteSlice writes a []byte, prefixed with a varuint32, to the underlying buffer.
func (w *Writer) ByteSlice(x *[]byte) {
	l := uint32(len(*x))
	w.Varuint32(&l)
	_, _ = w.w.Write(*x)
}

// Bytes appends a []byte to the underlying buffer.
func (w *Writer) Bytes(x *[]byte) {
	_, _ = w.w.Write(*x)
}

// ByteFloat writes a rotational float32 as a single byte to the underlying buffer.
func (w *Writer) ByteFloat(x *float32) {
	_ = w.w.WriteByte(byte(*x / (360.0 / 256.0)))
}

// Vec3 writes an mgl32.Vec3 as 3 float32s to the underlying buffer.
func (w *Writer) Vec3(x *mgl32.Vec3) {
	w.Float32(&x[0])
	w.Float32(&x[1])
	w.Float32(&x[2])
}

// Vec2 writes an mgl32.Vec2 as 2 float32s to the underlying buffer.
func (w *Writer) Vec2(x *mgl32.Vec2) {
	w.Float32(&x[0])
	w.Float32(&x[1])
}

// BlockPos writes a BlockPos as 3 varint32s to the underlying buffer.
func (w *Writer) BlockPos(x *BlockPos) {
	w.Varint32(&x[0])
	w.Varint32(&x[1])
	w.Varint32(&x[2])
}

// UBlockPos writes a BlockPos as 2 varint32s and a varuint32 to the underlying buffer.
func (w *Writer) UBlockPos(x *BlockPos) {
	w.Varint32(&x[0])
	y := uint32(x[1])
	w.Varuint32(&y)
	w.Varint32(&x[2])
}

// VarRGBA writes a color.RGBA x as a varuint32 to the underlying buffer.
func (w *Writer) VarRGBA(x *color.RGBA) {
	val := uint32(x.R) | uint32(x.G)<<8 | uint32(x.B)<<16 | uint32(x.A)<<24
	w.Varuint32(&val)
}

// UUID writes a UUID to the underlying buffer.
func (w *Writer) UUID(x *uuid.UUID) {
	b := append((*x)[8:], (*x)[:8]...)
	for i, j := 0, 15; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	_, _ = w.w.Write(b)
}

// EntityMetadata writes an entity metadata map x to the underlying buffer.
func (w *Writer) EntityMetadata(x *map[uint32]interface{}) {
	l := uint32(len(*x))
	w.Varuint32(&l)
	for key, value := range *x {
		w.Varuint32(&key)
		switch v := value.(type) {
		case byte:
			w.Varuint32(&entityDataByte)
			w.Uint8(&v)
		case int16:
			w.Varuint32(&entityDataInt16)
			w.Int16(&v)
		case int32:
			w.Varuint32(&entityDataInt32)
			w.Varint32(&v)
		case float32:
			w.Varuint32(&entityDataFloat32)
			w.Float32(&v)
		case string:
			w.Varuint32(&entityDataString)
			w.String(&v)
		case map[string]interface{}:
			w.Varuint32(&entityDataCompoundTag)
			w.NBT(&v, nbt.NetworkLittleEndian)
		case BlockPos:
			w.Varuint32(&entityDataBlockPos)
			w.BlockPos(&v)
		case int64:
			w.Varuint32(&entityDataInt64)
			w.Varint64(&v)
		case mgl32.Vec3:
			w.Varuint32(&entityDataVec3)
			w.Vec3(&v)
		default:
			w.UnknownEnumOption(reflect.TypeOf(value), "entity metadata")
		}
	}
}

// Item writes an ItemStack x to the underlying buffer.
func (w *Writer) Item(x *ItemStack) {
	w.Varint32(&x.NetworkID)
	if x.NetworkID == 0 {
		// The item was air, so there's no more data to follow. Return immediately.
		return
	}
	aux := int32(x.MetadataValue<<8) | int32(x.Count)
	w.Varint32(&aux)
	if len(x.NBTData) != 0 {
		userDataMarker := int16(-1)
		userDataVer := uint8(1)

		w.Int16(&userDataMarker)
		w.Uint8(&userDataVer)
		w.NBT(&x.NBTData, nbt.NetworkLittleEndian)
	} else {
		userDataMarker := int16(0)

		w.Int16(&userDataMarker)
	}
	placeOnLen := int32(len(x.CanBePlacedOn))
	canBreak := int32(len(x.CanBreak))

	w.Varint32(&placeOnLen)
	for _, block := range x.CanBePlacedOn {
		w.String(&block)
	}
	w.Varint32(&canBreak)
	for _, block := range x.CanBreak {
		w.String(&block)
	}

	const shieldID = 513
	if x.NetworkID == shieldID {
		var blockingTick int64
		w.Varint64(&blockingTick)
	}
}

// Varint64 writes an int64 as 1-10 bytes to the underlying buffer.
func (w *Writer) Varint64(x *int64) {
	u := *x
	ux := uint64(u) << 1
	if u < 0 {
		ux = ^ux
	}
	for ux >= 0x80 {
		_ = w.w.WriteByte(byte(ux) | 0x80)
		ux >>= 7
	}
	_ = w.w.WriteByte(byte(ux))
}

// Varuint64 writes a uint64 as 1-10 bytes to the underlying buffer.
func (w *Writer) Varuint64(x *uint64) {
	u := *x
	for u >= 0x80 {
		_ = w.w.WriteByte(byte(u) | 0x80)
		u >>= 7
	}
	_ = w.w.WriteByte(byte(u))
}

// Varint32 writes an int32 as 1-5 bytes to the underlying buffer.
func (w *Writer) Varint32(x *int32) {
	u := *x
	ux := uint32(u) << 1
	if u < 0 {
		ux = ^ux
	}
	for ux >= 0x80 {
		_ = w.w.WriteByte(byte(ux) | 0x80)
		ux >>= 7
	}
	_ = w.w.WriteByte(byte(ux))
}

// Varuint32 writes a uint32 as 1-5 bytes to the underlying buffer.
func (w *Writer) Varuint32(x *uint32) {
	u := *x
	for u >= 0x80 {
		_ = w.w.WriteByte(byte(u) | 0x80)
		u >>= 7
	}
	_ = w.w.WriteByte(byte(u))
}

// NBT writes a map as NBT to the underlying buffer using the encoding passed.
func (w *Writer) NBT(x *map[string]interface{}, encoding nbt.Encoding) {
	if err := nbt.NewEncoderWithEncoding(w.w, encoding).Encode(*x); err != nil {
		panic(err)
	}
}

// NBTList writes a slice as NBT to the underlying buffer using the encoding passed.
func (w *Writer) NBTList(x *[]interface{}, encoding nbt.Encoding) {
	if err := nbt.NewEncoderWithEncoding(w.w, encoding).Encode(*x); err != nil {
		panic(err)
	}
}

// UnknownEnumOption panics with an unknown enum option error.
func (w *Writer) UnknownEnumOption(value interface{}, enum string) {
	w.panicf("unknown value '%v' for enum type '%v'", value, enum)
}

// InvalidValue panics with an invalid value error.
func (w *Writer) InvalidValue(value interface{}, forField, reason string) {
	w.panicf("invalid value '%v' for %v: %v", value, forField, reason)
}

// panicf panics with the format and values passed.
func (w *Writer) panicf(format string, a ...interface{}) {
	panic(fmt.Errorf(format, a...))
}
