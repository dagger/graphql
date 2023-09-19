package graphql

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/dagger/graphql/language/ast"
)

const (
	// Aligned with JavaScript integers that are limited to the range -(2^53 - 1) to 2^53 - 1,
	// due to being encoded as double-precision floating-point numbers
	MaxInt int64 = 9007199254740991
	MinInt int64 = -9007199254740991
)

// As per the GraphQL Spec, Integers are only treated as valid when a valid
// 32-bit signed integer, providing the broadest support across platforms.
//
// n.b. JavaScript's integers are safe between -(2^53 - 1) and 2^53 - 1 because
// they are internally represented as IEEE 754 doubles.
func coerceInt(value any) (any, error) {
	switch value := value.(type) {
	case bool:
		if value {
			return int64(1), nil
		}
		return int64(0), nil
	case *bool:
		if value == nil {
			return nil, nil
		}
		return coerceInt(*value)
	case int:
		return int64(value), nil
	case *int:
		if value == nil {
			return nil, nil
		}
		return coerceInt(*value)
	case int8:
		return int64(value), nil
	case *int8:
		if value == nil {
			return nil, nil
		}
		return int64(*value), nil
	case int16:
		return int64(value), nil
	case *int16:
		if value == nil {
			return nil, nil
		}
		return int64(*value), nil
	case int32:
		return int64(value), nil
	case *int32:
		if value == nil {
			return nil, nil
		}
		return int64(*value), nil
	case int64:
		return intOrNil(value)
	case *int64:
		if value == nil {
			return nil, nil
		}
		return coerceInt(*value)
	case uint:
		return int64(value), nil
	case *uint:
		if value == nil {
			return nil, nil
		}
		return coerceInt(*value)
	case uint8:
		return int64(value), nil
	case *uint8:
		if value == nil {
			return nil, nil
		}
		return int64(*value), nil
	case uint16:
		return int64(value), nil
	case *uint16:
		if value == nil {
			return nil, nil
		}
		return int64(*value), nil
	case uint32:
		return int64(value), nil
	case *uint32:
		if value == nil {
			return nil, nil
		}
		return coerceInt(*value)
	case uint64:
		if value > uint64(math.MaxInt64) {
			// bypass intOrNil
			return nil, fmt.Errorf("value %d is greater than MaxInt64", value)
		}
		return intOrNil(int64(value))
	case *uint64:
		if value == nil {
			return nil, nil
		}
		return coerceInt(*value)
	case float32:
		return intOrNil(int64(value))
	case *float32:
		if value == nil {
			return nil, nil
		}
		return coerceInt(*value)
	case float64:
		return intOrNil(int64(value))
	case *float64:
		if value == nil {
			return nil, nil
		}
		return coerceInt(*value)
	case string:
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, err
		}
		return coerceInt(val)
	case *string:
		if value == nil {
			return nil, nil
		}
		return coerceInt(*value)
	}

	// If the value cannot be transformed into an int, return nil instead of '0'
	// to denote 'no integer found'
	return nil, fmt.Errorf("cannot coerce %T to int", value)
}

// Integers are only safe when between -(2^53 - 1) and 2^53 - 1 due to being
// encoded in JavaScript and represented in JSON as double-precision floating
// point numbers, as specified by IEEE 754.
func intOrNil(value int64) (int64, error) {
	if MinInt <= value && value <= MaxInt {
		return value, nil
	}
	return 0, fmt.Errorf("value %d is out of range (%d, %d) for Int64", value, MinInt, MaxInt)
}

// Int is the GraphQL Integer type definition.
var Int = NewScalar(ScalarConfig{
	Name: "Int",
	Description: "The `Int` scalar type represents non-fractional signed whole numeric " +
		"values. Int can represent values between -(2^31) and 2^31 - 1. ",
	Serialize:  coerceInt,
	ParseValue: coerceInt,
	ParseLiteral: func(valueAST ast.Value) (any, error) {
		switch valueAST := valueAST.(type) {
		case *ast.IntValue:
			return strconv.ParseInt(valueAST.Value, 10, 64)
		}
		return nil, fmt.Errorf("cannot parse %T to int", valueAST)
	},
})

func coerceFloat(value any) (any, error) {
	switch value := value.(type) {
	case bool:
		if value {
			return 1.0, nil
		}
		return 0.0, nil
	case *bool:
		if value == nil {
			return nil, nil
		}
		return coerceFloat(*value)
	case int:
		return float64(value), nil
	case *int:
		if value == nil {
			return nil, nil
		}
		return coerceFloat(*value)
	case int8:
		return float64(value), nil
	case *int8:
		if value == nil {
			return nil, nil
		}
		return coerceFloat(*value)
	case int16:
		return float64(value), nil
	case *int16:
		if value == nil {
			return nil, nil
		}
		return coerceFloat(*value)
	case int32:
		return float64(value), nil
	case *int32:
		if value == nil {
			return nil, nil
		}
		return coerceFloat(*value)
	case int64:
		return float64(value), nil
	case *int64:
		if value == nil {
			return nil, nil
		}
		return coerceFloat(*value)
	case uint:
		return float64(value), nil
	case *uint:
		if value == nil {
			return nil, nil
		}
		return coerceFloat(*value)
	case uint8:
		return float64(value), nil
	case *uint8:
		if value == nil {
			return nil, nil
		}
		return coerceFloat(*value)
	case uint16:
		return float64(value), nil
	case *uint16:
		if value == nil {
			return nil, nil
		}
		return coerceFloat(*value)
	case uint32:
		return float64(value), nil
	case *uint32:
		if value == nil {
			return nil, nil
		}
		return coerceFloat(*value)
	case uint64:
		return float64(value), nil
	case *uint64:
		if value == nil {
			return nil, nil
		}
		return coerceFloat(*value)
	case float32:
		return value, nil
	case *float32:
		if value == nil {
			return nil, nil
		}
		return coerceFloat(*value)
	case float64:
		return value, nil
	case *float64:
		if value == nil {
			return nil, nil
		}
		return coerceFloat(*value)
	case string:
		return strconv.ParseFloat(value, 64)
	case *string:
		if value == nil {
			return nil, nil
		}
		return coerceFloat(*value)
	}

	// If the value cannot be transformed into an float, return nil instead of '0.0'
	// to denote 'no float found'
	return nil, fmt.Errorf("cannot coerce %T to float", value)
}

// Float is the GraphQL float type definition.
var Float = NewScalar(ScalarConfig{
	Name: "Float",
	Description: "The `Float` scalar type represents signed double-precision fractional " +
		"values as specified by " +
		"[IEEE 754](http://en.wikipedia.org/wiki/IEEE_floating_point). ",
	Serialize:  coerceFloat,
	ParseValue: coerceFloat,
	ParseLiteral: func(valueAST ast.Value) (any, error) {
		switch valueAST := valueAST.(type) {
		case *ast.FloatValue:
			return strconv.ParseFloat(valueAST.Value, 64)
		case *ast.IntValue:
			return strconv.ParseFloat(valueAST.Value, 64)
		}
		return nil, fmt.Errorf("cannot parse %T to float", valueAST)
	},
})

func coerceString(value any) (any, error) {
	if v, ok := value.(*string); ok {
		if v == nil {
			return nil, nil
		}
		return *v, nil
	}
	return fmt.Sprintf("%v", value), nil
}

// String is the GraphQL string type definition
var String = NewScalar(ScalarConfig{
	Name: "String",
	Description: "The `String` scalar type represents textual data, represented as UTF-8 " +
		"character sequences. The String type is most often used by GraphQL to " +
		"represent free-form human-readable text.",
	Serialize:  coerceString,
	ParseValue: coerceString,
	ParseLiteral: func(valueAST ast.Value) (any, error) {
		switch valueAST := valueAST.(type) {
		case *ast.StringValue:
			return valueAST.Value, nil
		}
		return nil, fmt.Errorf("cannot parse %T to string", valueAST)
	},
})

func coerceBool(value any) (any, error) {
	switch value := value.(type) {
	case bool:
		return value, nil
	case *bool:
		if value == nil {
			return nil, nil
		}
		return *value, nil
	case string:
		switch value {
		case "", "false":
			return false, nil
		}
		return true, nil
	case *string:
		if value == nil {
			return nil, nil
		}
		return coerceBool(*value)
	case float64:
		if value != 0 {
			return true, nil
		}
		return false, nil
	case *float64:
		if value == nil {
			return nil, nil
		}
		return coerceBool(*value)
	case float32:
		if value != 0 {
			return true, nil
		}
		return false, nil
	case *float32:
		if value == nil {
			return nil, nil
		}
		return coerceBool(*value)
	case int:
		if value != 0 {
			return true, nil
		}
		return false, nil
	case *int:
		if value == nil {
			return nil, nil
		}
		return coerceBool(*value)
	case int8:
		if value != 0 {
			return true, nil
		}
		return false, nil
	case *int8:
		if value == nil {
			return nil, nil
		}
		return coerceBool(*value)
	case int16:
		if value != 0 {
			return true, nil
		}
		return false, nil
	case *int16:
		if value == nil {
			return nil, nil
		}
		return coerceBool(*value)
	case int32:
		if value != 0 {
			return true, nil
		}
		return false, nil
	case *int32:
		if value == nil {
			return nil, nil
		}
		return coerceBool(*value)
	case int64:
		if value != 0 {
			return true, nil
		}
		return false, nil
	case *int64:
		if value == nil {
			return nil, nil
		}
		return coerceBool(*value)
	case uint:
		if value != 0 {
			return true, nil
		}
		return false, nil
	case *uint:
		if value == nil {
			return nil, nil
		}
		return coerceBool(*value)
	case uint8:
		if value != 0 {
			return true, nil
		}
		return false, nil
	case *uint8:
		if value == nil {
			return nil, nil
		}
		return coerceBool(*value)
	case uint16:
		if value != 0 {
			return true, nil
		}
		return false, nil
	case *uint16:
		if value == nil {
			return nil, nil
		}
		return coerceBool(*value)
	case uint32:
		if value != 0 {
			return true, nil
		}
		return false, nil
	case *uint32:
		if value == nil {
			return nil, nil
		}
		return coerceBool(*value)
	case uint64:
		if value != 0 {
			return true, nil
		}
		return false, nil
	case *uint64:
		if value == nil {
			return nil, nil
		}
		return coerceBool(*value)
	}
	return false, fmt.Errorf("cannot coerce %T to bool", value)
}

// Boolean is the GraphQL boolean type definition
var Boolean = NewScalar(ScalarConfig{
	Name:        "Boolean",
	Description: "The `Boolean` scalar type represents `true` or `false`.",
	Serialize:   coerceBool,
	ParseValue:  coerceBool,
	ParseLiteral: func(valueAST ast.Value) (any, error) {
		switch valueAST := valueAST.(type) {
		case *ast.BooleanValue:
			return valueAST.Value, nil
		}
		return nil, fmt.Errorf("cannot parse %T to bool", valueAST)
	},
})

// ID is the GraphQL id type definition
var ID = NewScalar(ScalarConfig{
	Name: "ID",
	Description: "The `ID` scalar type represents a unique identifier, often used to " +
		"refetch an object or as key for a cache. The ID type appears in a JSON " +
		"response as a String; however, it is not intended to be human-readable. " +
		"When expected as an input type, any string (such as `\"4\"`) or integer " +
		"(such as `4`) input value will be accepted as an ID.",
	Serialize:  coerceString,
	ParseValue: coerceString,
	ParseLiteral: func(valueAST ast.Value) (any, error) {
		switch valueAST := valueAST.(type) {
		case *ast.IntValue:
			return valueAST.Value, nil
		case *ast.StringValue:
			return valueAST.Value, nil
		}
		return nil, fmt.Errorf("cannot parse %T to string", valueAST)
	},
})

func serializeDateTime(value any) (any, error) {
	switch value := value.(type) {
	case time.Time:
		buff, err := value.MarshalText()
		if err != nil {
			return nil, err
		}

		return string(buff), nil
	case *time.Time:
		if value == nil {
			return nil, nil
		}
		return serializeDateTime(*value)
	default:
		if value == nil {
			return nil, nil
		}
		return nil, fmt.Errorf("cannot serialize %T to time.Time", value)
	}
}

func unserializeDateTime(value any) (any, error) {
	switch value := value.(type) {
	case []byte:
		t := time.Time{}
		err := t.UnmarshalText(value)
		if err != nil {
			return nil, err
		}

		return t, nil
	case string:
		return unserializeDateTime([]byte(value))
	case *string:
		if value == nil {
			return nil, nil
		}
		return unserializeDateTime([]byte(*value))
	case time.Time:
		return value, nil
	default:
		if value == nil {
			return nil, nil
		}
		return nil, fmt.Errorf("cannot unserialize %T to time.Time", value)
	}
}

var DateTime = NewScalar(ScalarConfig{
	Name: "DateTime",
	Description: "The `DateTime` scalar type represents a DateTime." +
		" The DateTime is serialized as an RFC 3339 quoted string",
	Serialize:  serializeDateTime,
	ParseValue: unserializeDateTime,
	ParseLiteral: func(valueAST ast.Value) (any, error) {
		switch valueAST := valueAST.(type) {
		case *ast.StringValue:
			return unserializeDateTime(valueAST.Value)
		}
		return nil, fmt.Errorf("cannot parse %T to string", valueAST)
	},
})
