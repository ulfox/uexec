package uexec

// Action holds the first return item
type Action struct {
	Values         []interface{}
	CallBackFunc   interface{}
	CallBackArgs   []interface{}
	CallBackValues interface{}
	Err            interface{}
}

// AddCallBack adds a callback function that can run on Exec errors
func (t Action) AddCallBack(callBackFunc interface{}, callBackArgs ...interface{}) Action {
	t.CallBackFunc = callBackFunc
	t.CallBackArgs = append(t.CallBackArgs, callBackArgs...)

	return t
}

// CallBack function for running after an error has been caught
func (t Action) CallBack() Action {
	switch fn := t.CallBackFunc.(type) {
	case func(...interface{}) interface{}:
		t.CallBackValues = fn(t.CallBackArgs...)
	}
	return t
}

// Get Returns a value
func (t Action) Get(n int) interface{} {
	return t.Values[n]
}

// Byte returns a value of type byte
func (t Action) Byte(n int) byte {
	return t.Values[n].(byte)
}

// ByteS returns an array of type byte
func (t Action) ByteS(n int) []byte {
	return t.Values[n].([]byte)
}

func (t Action) String(n int) string {
	return t.Values[n].(string)
}

// StringS returns an array of strings
func (t Action) StringS(n int) []string {
	return t.Values[n].([]string)
}

// Int returns a value of type int
func (t Action) Int(n int) int {
	return t.Values[n].(int)
}

// Int8 returns a value of type int8
func (t Action) Int8(n int) int8 {
	return t.Values[n].(int8)
}

// Int16 returns a value of type int16
func (t Action) Int16(n int) int16 {
	return t.Values[n].(int16)
}

// Int32 returns a value of type int32
func (t Action) Int32(n int) int32 {
	return t.Values[n].(int32)
}

// Int64 returns a value of type int64
func (t Action) Int64(n int) int64 {
	return t.Values[n].(int64)
}

// IntS returns an array of type int
func (t Action) IntS(n int) []int {
	return t.Values[n].([]int)
}

// Int8S returns an array of type int8
func (t Action) Int8S(n int) []int8 {
	return t.Values[n].([]int8)
}

// Int16S returns an array of type int16
func (t Action) Int16S(n int) []int16 {
	return t.Values[n].([]int16)
}

// Int32S returns an array of type int32
func (t Action) Int32S(n int) []int32 {
	return t.Values[n].([]int32)
}

// Int64S returns an array of type int64
func (t Action) Int64S(n int) []int64 {
	return t.Values[n].([]int64)
}

// Float32 returns a value of type float32
func (t Action) Float32(n int) float32 {
	return t.Values[n].(float32)
}

// Float64 returns a value of type float64
func (t Action) Float64(n int) float64 {
	return t.Values[n].(float64)
}

// Float32S returns an array of type float32
func (t Action) Float32S(n int) []float32 {
	return t.Values[n].([]float32)
}

// Float64S returns an array of type float64
func (t Action) Float64S(n int) []float64 {
	return t.Values[n].([]float64)
}
