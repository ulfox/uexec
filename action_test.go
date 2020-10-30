package uexec

import (
	"errors"
	"fmt"
	"testing"

	"github.com/likexian/gokit/assert"
)

// func someCallBack(v ...interface{})
func someFunc(v interface{}, e interface{}) (interface{}, interface{}) {
	return v, e
}

func someCallback(args ...interface{}) interface{} {
	return args
}

//TestOutputs run unit tests on action outputs
func testOutputs(t *testing.T) {
	t.Parallel()

	erH := NewErrorHandler().SetLogLevel("trace").
		SetElasticity(true).
		AddGenericCallBack(someCallback).
		OnErr("callback")

	testCases := []struct {
		value           interface{}
		err             interface{}
		RValues         []interface{}
		RCallBackFunc   interface{}
		RCallBackArgs   []interface{}
		RCallBackValues interface{}
		RErr            interface{}
	}{
		{"someString", nil, []interface{}{"someString", nil}, nil, []interface{}{"someString", nil}, nil, nil},
		{"someString", 1, []interface{}{"someString", 1}, nil, []interface{}{"someString", 1}, nil, nil},
		{"someString", errors.New("someError"), []interface{}{"someString"}, nil, []interface{}{"someString", errors.New("someError")}, nil, errors.New("someError")},
		{"someString", 1, []interface{}{"someString", 1}, nil, []interface{}{"someString", 1}, nil, nil},
		{[]string{"someString"}, nil, []interface{}{[]string{"someString"}, nil}, nil, []interface{}{[]string{"someString"}, nil}, nil, nil},
	}

	for _, testCase := range testCases {
		try := erH.Exec(someFunc(testCase.value, testCase.err))
		assert.Equal(t, testCase.RValues, try.Values, fmt.Sprintf("Expected: %v", testCase.RValues))
		assert.Equal(t, testCase.RErr, try.Err, fmt.Sprintf("Expected: %v", testCase.RErr))
		assert.Equal(t, testCase.RCallBackFunc, try.CallBackFunc, fmt.Sprintf("Expected: %v", testCase.RCallBackFunc))
		assert.Equal(t, testCase.RCallBackArgs, try.CallBackArgs, fmt.Sprintf("Expected: %v", testCase.RCallBackArgs))
		assert.Equal(t, testCase.RCallBackValues, try.CallBackValues, fmt.Sprintf("Expected: %v", testCase.RCallBackValues))
	}

	try := erH.Exec(someFunc([]string{"someString"}, nil))
	assert.Equal(t, []string{"someString"}, try.StringS(0))

	try = erH.Exec(someFunc("someString", nil))
	assert.Equal(t, "someString", try.String(0))

	try = erH.Exec(someFunc([]byte{}, nil))
	assert.Equal(t, []byte{}, try.ByteS(0))

	try = erH.Exec(someFunc(byte(0), nil))
	assert.Equal(t, byte(0), try.Byte(0))

	try = erH.Exec(someFunc(errors.New("someError"), nil))
	assert.Equal(t, errors.New("someError"), try.Err)

	try = erH.Exec(someFunc(errors.New("someError"), errors.New("someError")))
	assert.Equal(t, errors.New("someError"), try.GetError(0))

	try = erH.ErP(1).Exec(someFunc(errors.New("someError"), nil))
	assert.Equal(t, errors.New("someError"), try.GetError(0))

	try = erH.Exec(someFunc([]int{1}, nil))
	assert.Equal(t, []int{1}, try.IntS(0))
}
