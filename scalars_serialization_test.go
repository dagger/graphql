package graphql_test

import (
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/dagger/graphql"
)

type intSerializationTest struct {
	Value    any
	Expected any
}

type float64SerializationTest struct {
	Value    any
	Expected any
}

type stringSerializationTest struct {
	Value    any
	Expected string
}

type dateTimeSerializationTest struct {
	Value    any
	Expected any
}

type boolSerializationTest struct {
	Value    any
	Expected bool
}

func TestTypeSystem_Scalar_SerializesOutputInt(t *testing.T) {
	tests := []intSerializationTest{
		{1, int64(1)},
		{0, int64(0)},
		{-1, int64(-1)},
		{float32(0.1), int64(0)},
		{float32(1.1), int64(1)},
		{float32(-1.1), int64(-1)},
		{float32(1e5), int64(100000)},
		{9876504321, int64(9876504321)},
		{-9876504321, int64(-9876504321)},
		{float32(math.MaxFloat32), nil},
		{float64(0.1), int64(0)},
		{float64(1.1), int64(1)},
		{float64(-1.1), int64(-1)},
		{float64(1e5), int64(100000)},
		{float64(math.MaxFloat32), nil},
		{float64(math.MaxFloat64), nil},
		// safe Go/Javascript `int`, bigger than 2^32, but more than graphQL Int spec
		{9876504321, int64(9876504321)},
		{-9876504321, int64(-9876504321)},
		// Too big to represent as an Int in Go, JavaScript or GraphQL
		{float64(1e100), nil},
		{float64(-1e100), nil},
		{"-1.1", int64(-1)},
		{"one", nil},
		{false, int64(0)},
		{true, int64(1)},
		{int8(1), int64(1)},
		{int16(1), int64(1)},
		{int32(1), int64(1)},
		{int64(1), int64(1)},
		{uint(1), int64(1)},
		// Maybe a safe Go `uint`, bigger than 2^32, but more than graphQL Int spec
		{uint(math.MaxInt32 + 1), int64(2147483648)},
		{uint8(1), int64(1)},
		{uint16(1), int64(1)},
		{uint32(1), int64(1)},
		{uint32(math.MaxUint32), int64(4294967295)},
		{uint64(1), int64(1)},
		{uint64(math.MaxInt32), int64(math.MaxInt32)},
		{int64(math.MaxInt32) + int64(1), int64(2147483648)},
		{int64(math.MinInt32) - int64(1), int64(-2147483649)},
		{uint64(math.MaxInt64) + uint64(1), nil},
		{byte(127), int64(127)},
		{'世', int64('世')},
		// testing types that don't match a value in the array.
		{[]int{}, nil},
	}

	for i, test := range tests {
		val := graphql.Int.Serialize(test.Value)
		if val != test.Expected {
			reflectedTestValue := reflect.ValueOf(test.Value)
			reflectedExpectedValue := reflect.ValueOf(test.Expected)
			reflectedValue := reflect.ValueOf(val)
			t.Fatalf("Failed test #%d - Int.Serialize(%v(%v)), expected: %v(%v), got %v(%v)",
				i, reflectedTestValue.Type(), test.Value,
				reflectedExpectedValue.Type(), test.Expected,
				reflectedValue.Type(), val,
			)
		}
	}
}

func TestTypeSystem_Scalar_SerializesOutputFloat(t *testing.T) {
	tests := []float64SerializationTest{
		{int(1), 1.0},
		{int(0), 0.0},
		{int(-1), -1.0},
		{float32(0.1), float32(0.1)},
		{float32(1.1), float32(1.1)},
		{float32(-1.1), float32(-1.1)},
		{float64(0.1), float64(0.1)},
		{float64(1.1), float64(1.1)},
		{float64(-1.1), float64(-1.1)},
		{"-1.1", -1.1},
		{"one", nil},
		{false, 0.0},
		{true, 1.0},
	}

	for i, test := range tests {
		val := graphql.Float.Serialize(test.Value)
		if val != test.Expected {
			reflectedTestValue := reflect.ValueOf(test.Value)
			reflectedExpectedValue := reflect.ValueOf(test.Expected)
			reflectedValue := reflect.ValueOf(val)
			t.Fatalf("Failed test #%d - Float.Serialize(%v(%v)), expected: %v(%v), got %v(%v)",
				i, reflectedTestValue.Type(), test.Value,
				reflectedExpectedValue.Type(), test.Expected,
				reflectedValue.Type(), val,
			)
		}
	}
}

func TestTypeSystem_Scalar_SerializesOutputStrings(t *testing.T) {
	tests := []stringSerializationTest{
		{"string", "string"},
		{int(1), "1"},
		{float32(-1.1), "-1.1"},
		{float64(-1.1), "-1.1"},
		{true, "true"},
		{false, "false"},
	}

	for _, test := range tests {
		val := graphql.String.Serialize(test.Value)
		if val != test.Expected {
			reflectedValue := reflect.ValueOf(test.Value)
			t.Fatalf("Failed String.Serialize(%v(%v)), expected: %v, got %v", reflectedValue.Type(), test.Value, test.Expected, val)
		}
	}
}

func TestTypeSystem_Scalar_SerializesOutputBoolean(t *testing.T) {
	tests := []boolSerializationTest{
		{"true", true},
		{"false", false},
		{"string", true},
		{"", false},
		{int(1), true},
		{int(0), false},
		{true, true},
		{false, false},
	}

	for _, test := range tests {
		val := graphql.Boolean.Serialize(test.Value)
		if val != test.Expected {
			reflectedValue := reflect.ValueOf(test.Value)
			t.Fatalf("Failed String.Boolean(%v(%v)), expected: %v, got %v", reflectedValue.Type(), test.Value, test.Expected, val)
		}
	}
}

func TestTypeSystem_Scalar_SerializeOutputDateTime(t *testing.T) {
	now := time.Now()
	nowString, err := now.MarshalText()
	if err != nil {
		t.Fatal(err)
	}

	tests := []dateTimeSerializationTest{
		{"string", nil},
		{int(1), nil},
		{float32(-1.1), nil},
		{float64(-1.1), nil},
		{true, nil},
		{false, nil},
		{now, string(nowString)},
		{&now, string(nowString)},
	}

	for _, test := range tests {
		val := graphql.DateTime.Serialize(test.Value)
		if val != test.Expected {
			reflectedValue := reflect.ValueOf(test.Value)
			t.Fatalf("Failed DateTime.Serialize(%v(%v)), expected: %v, got %v", reflectedValue.Type(), test.Value, test.Expected, val)
		}
	}
}
