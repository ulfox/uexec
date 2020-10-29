package uexec

import (
	"os"

	"github.com/sirupsen/logrus"
)

// ErrorHandler struct for handling errors
type ErrorHandler struct {
	Logger logrus.FieldLogger
	Try    Try
}

// NewErrorHandler factory method to create a new errorhandler
func NewErrorHandler() *ErrorHandler {
	errorHandler := &ErrorHandler{
		Logger: logrus.New(),
		Try:    Try{},
	}
	return errorHandler
}

// Exec factory for creating a new try
func (e *ErrorHandler) Exec(cmd ...interface{}) Try {
	if len(cmd) == 0 {
		return Try{}
	}

	try := Try{}

	var erP int = -1
	for i, j := range cmd {
		if isErr, err := try.getErr(j); isErr {
			try.Err = append(try.Err, err)
			erP = i
		}

	}

	if erP >= 0 {
		e.CheckE(try.Err[0])
	}

	for i, j := range cmd {
		if i == erP {
			continue
		}
		try.Values = append(try.Values, j)
	}

	return try
}

// CheckE for checking the an error result and exiting if it is not nil
func (e *ErrorHandler) CheckE(args interface{}) {
	if args == nil {
		return
	}
	var err error = args.(error)

	e.Logger.Error(err)
	os.Exit(1)
}

// Try holds the first return item
type Try struct {
	Values       []interface{}
	Err          []interface{}
	optFunctions []func(...interface{}) []interface{}
}

func (t *Try) getErr(x interface{}) (bool, error) {
	switch x.(type) {
	case error:
		return true, x.(error)
	default:
		return false, nil
	}
}

// Get Returns a value
func (t Try) Get(n int) interface{} {
	return t.Values[n]
}

// Byte returns a value of type byte
func (t Try) Byte(n int) byte {
	return t.Values[n].(byte)
}

// ByteS returns an array of type byte
func (t Try) ByteS(n int) []byte {
	return t.Values[n].([]byte)
}

func (t Try) String(n int) string {
	return t.Values[n].(string)
}

// StringS returns an array of strings
func (t Try) StringS(n int) []string {
	return t.Values[n].([]string)
}

// Int returns a value of type int
func (t Try) Int(n int) int {
	return t.Values[n].(int)
}

// Int8 returns a value of type int8
func (t Try) Int8(n int) int8 {
	return t.Values[n].(int8)
}

// Int16 returns a value of type int16
func (t Try) Int16(n int) int16 {
	return t.Values[n].(int16)
}

// Int32 returns a value of type int32
func (t Try) Int32(n int) int32 {
	return t.Values[n].(int32)
}

// Int64 returns a value of type int64
func (t Try) Int64(n int) int64 {
	return t.Values[n].(int64)
}

// IntS returns an array of type int
func (t Try) IntS(n int) []int {
	return t.Values[n].([]int)
}

// Int8S returns an array of type int8
func (t Try) Int8S(n int) []int8 {
	return t.Values[n].([]int8)
}

// Int16S returns an array of type int16
func (t Try) Int16S(n int) []int16 {
	return t.Values[n].([]int16)
}

// Int32S returns an array of type int32
func (t Try) Int32S(n int) []int32 {
	return t.Values[n].([]int32)
}

// Int64S returns an array of type int64
func (t Try) Int64S(n int) []int64 {
	return t.Values[n].([]int64)
}

// Float32 returns a value of type float32
func (t Try) Float32(n int) float32 {
	return t.Values[n].(float32)
}

// Float64 returns a value of type float64
func (t Try) Float64(n int) float64 {
	return t.Values[n].(float64)
}

// Float32S returns an array of type float32
func (t Try) Float32S(n int) []float32 {
	return t.Values[n].([]float32)
}

// Float64S returns an array of type float64
func (t Try) Float64S(n int) []float64 {
	return t.Values[n].([]float64)
}
