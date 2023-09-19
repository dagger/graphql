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
	Fails    bool
}

type float64SerializationTest struct {
	Value    any
	Expected any
	Fails    bool
}

type stringSerializationTest struct {
	Value    any
	Expected string
}

type dateTimeSerializationTest struct {
	Value    any
	Expected any
	Fails    bool
}

type boolSerializationTest struct {
	Value    any
	Expected bool
}

func TestTypeSystem_Scalar_SerializesOutputInt(t *testing.T) {
	tests := []intSerializationTest{
		{1, int64(1), false},
		{0, int64(0), false},
		{-1, int64(-1), false},
		{float32(0.1), int64(0), false},
		{float32(1.1), int64(1), false},
		{float32(-1.1), int64(-1), false},
		{float32(1e5), int64(100000), false},
		{9876504321, int64(9876504321), false},
		{-9876504321, int64(-9876504321), false},
		{float32(math.MaxFloat32), nil, true},
		{float64(0.1), int64(0), false},
		{float64(1.1), int64(1), false},
		{float64(-1.1), int64(-1), false},
		{float64(1e5), int64(100000), false},
		{float64(math.MaxFloat32), nil, true},
		{float64(math.MaxFloat64), nil, true},
		// safe Go/Javascript `int`, bigger than 2^32, but more than graphQL Int spec
		{9876504321, int64(9876504321), false},
		{-9876504321, int64(-9876504321), false},
		// Too big to represent as an Int in Go, JavaScript or GraphQL
		{float64(1e100), nil, true},
		{float64(-1e100), nil, true},
		{"-1.1", int64(-1), false},
		{"one", nil, true},
		{false, int64(0), false},
		{true, int64(1), false},
		{int8(1), int64(1), false},
		{int16(1), int64(1), false},
		{int32(1), int64(1), false},
		{int64(1), int64(1), false},
		{uint(1), int64(1), false},
		// Maybe a safe Go `uint`, bigger than 2^32, but more than graphQL Int spec
		{uint(math.MaxInt32 + 1), int64(2147483648), false},
		{uint8(1), int64(1), false},
		{uint16(1), int64(1), false},
		{uint32(1), int64(1), false},
		{uint32(math.MaxUint32), int64(4294967295), false},
		{uint64(1), int64(1), false},
		{uint64(math.MaxInt32), int64(math.MaxInt32), false},
		{int64(math.MaxInt32) + int64(1), int64(2147483648), false},
		{int64(math.MinInt32) - int64(1), int64(-2147483649), false},
		{uint64(math.MaxInt64) + uint64(1), nil, true},
		{byte(127), int64(127), false},
		{'世', int64('世'), false},
		// testing types that don't match a value in the array.
		{[]int{}, nil, true},
	}

	for i, test := range tests {
		val, err := graphql.Int.Serialize(test.Value)
		if err != nil && !test.Fails {
			t.Fatalf("Failed test #%d - Int.Serialize(%v(%v)), expected: %v, got error: %v", i, reflect.TypeOf(test.Value), test.Value, test.Expected, err)
		}
		if err == nil && test.Fails {
			t.Fatalf("Failed test #%d - Int.Serialize(%v(%v)), should have failed", i, reflect.TypeOf(test.Value), test.Value)
		}
		if !test.Fails && val != test.Expected {
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
		{int(1), 1.0, false},
		{int(0), 0.0, false},
		{int(-1), -1.0, false},
		{float32(0.1), float32(0.1), false},
		{float32(1.1), float32(1.1), false},
		{float32(-1.1), float32(-1.1), false},
		{float64(0.1), float64(0.1), false},
		{float64(1.1), float64(1.1), false},
		{float64(-1.1), float64(-1.1), false},
		{"-1.1", -1.1, false},
		{"one", nil, true},
		{false, 0.0, false},
		{true, 1.0, false},
	}

	for i, test := range tests {
		val, err := graphql.Float.Serialize(test.Value)
		if err != nil && !test.Fails {
			t.Fatalf("Failed test #%d - Float.Serialize(%v(%v)), expected: %v, got error: %v", i, reflect.TypeOf(test.Value), test.Value, test.Expected, err)
		}
		if err == nil && test.Fails {
			t.Fatalf("Failed test #%d - Float.Serialize(%v(%v)), should have failed", i, reflect.TypeOf(test.Value), test.Value)
		}
		if !test.Fails && val != test.Expected {
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

	for i, test := range tests {
		val, err := graphql.String.Serialize(test.Value)
		if err != nil {
			t.Fatalf("Failed test #%d - String.Serialize(%v(%v)), expected: %v, got error: %v", i, reflect.TypeOf(test.Value), test.Value, test.Expected, err)
		}
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

	for i, test := range tests {
		val, err := graphql.Boolean.Serialize(test.Value)
		if err != nil {
			t.Fatalf("Failed test #%d - Boolean.Serialize(%v(%v)), expected: %v, got error: %v", i, reflect.TypeOf(test.Value), test.Value, test.Expected, err)
		}
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
		{"string", nil, true},
		{int(1), nil, true},
		{float32(-1.1), nil, true},
		{float64(-1.1), nil, true},
		{true, nil, true},
		{false, nil, true},
		{now, string(nowString), false},
		{&now, string(nowString), false},
	}

	for i, test := range tests {
		val, err := graphql.DateTime.Serialize(test.Value)
		if err != nil && !test.Fails {
			t.Fatalf("Failed test #%d - DateTime.Serialize(%v(%v)), expected: %v, got error: %v", i, reflect.TypeOf(test.Value), test.Value, test.Expected, err)
		}
		if err == nil && test.Fails {
			t.Fatalf("Failed test #%d - DateTime.Serialize(%v(%v)), should have failed", i, reflect.TypeOf(test.Value), test.Value)
		}
		if !test.Fails && val != test.Expected {
			reflectedValue := reflect.ValueOf(test.Value)
			t.Fatalf("Failed DateTime.Serialize(%v(%v)), expected: %v, got %v", reflectedValue.Type(), test.Value, test.Expected, val)
		}
	}
}
