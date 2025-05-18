package attribute

type Value struct {
	data interface{}
}

func (v *Value) IsSimple() bool {
	switch {
	case v.IsStr(),
		v.IsBool(),
		v.IsInt(),
		v.IsInt8(),
		v.IsInt16(),
		v.IsInt32(),
		v.IsInt64(),
		v.IsUint(),
		v.IsUint8(),
		v.IsUint16(),
		v.IsUint32(),
		v.IsUint64():
		return true
	}
	return false
}

func (v *Value) IsComplex() bool {
	_, ok := v.data.(map[string]interface{})
	return ok
}

func (v *Value) IsNil() bool {
	return v == nil || v.data == nil
}
func (v *Value) IsStr() bool {
	_, ok := v.data.(string)
	return ok
}
func (v *Value) IsBool() bool {
	_, ok := v.data.(bool)
	return ok
}
func (v *Value) IsFloat32() bool {
	_, ok := v.data.(float32)
	return ok
}
func (v *Value) IsFloat64() bool {
	_, ok := v.data.(float64)
	return ok
}
func (v *Value) IsInt() bool {
	_, ok := v.data.(int)
	return ok
}
func (v *Value) IsInt8() bool {
	_, ok := v.data.(int8)
	return ok
}
func (v *Value) IsInt16() bool {
	_, ok := v.data.(int16)
	return ok
}
func (v *Value) IsInt32() bool {
	_, ok := v.data.(int32)
	return ok
}
func (v *Value) IsInt64() bool {
	_, ok := v.data.(int64)
	return ok
}
func (v *Value) IsUint() bool {
	_, ok := v.data.(uint)
	return ok
}
func (v *Value) IsUint8() bool {
	_, ok := v.data.(uint8)
	return ok
}
func (v *Value) IsUint16() bool {
	_, ok := v.data.(uint16)
	return ok
}
func (v *Value) IsUint32() bool {
	_, ok := v.data.(uint32)
	return ok
}
func (v *Value) IsUint64() bool {
	_, ok := v.data.(uint64)
	return ok
}
