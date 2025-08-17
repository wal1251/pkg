package conversion

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPointerStrToEmptyStrOrValue(t *testing.T) {
	ptrStr := PointerStrToEmptyStrOrValue(nil)
	assert.Equal(t, "", ptrStr)

	s := "test"
	ptrStr = PointerStrToEmptyStrOrValue(&s)
	assert.Equal(t, "test", ptrStr)
}

func TestPointerBoolToBool(t *testing.T) {
	ptrBool := PointerBoolToBool(nil)
	assert.False(t, ptrBool)

	b := true
	ptrBool = PointerBoolToBool(&b)
	assert.True(t, ptrBool)
}

func TestPointerIntToEmptyStrOrValue(t *testing.T) {
	ptrInt := PointerIntToEmptyStrOrValue(nil)
	assert.Equal(t, "", ptrInt)

	i := 123
	ptrInt = PointerIntToEmptyStrOrValue(&i)
	assert.Equal(t, "123", ptrInt)
}

func TestStringFrom(t *testing.T) {
	str := StringFrom(nil)
	assert.Equal(t, "", str)

	str = StringFrom(123)
	assert.Equal(t, "123", str)

	str = StringFrom("test")
	assert.Equal(t, "test", str)
}

func TestStringFromNonZ(t *testing.T) {
	str := StringFromNonZ(0)
	assert.Equal(t, "", str)

	str = StringFromNonZ(123)
	assert.Equal(t, "123", str)

	str = StringFromNonZ("test")
	assert.Equal(t, "test", str)
}

func TestBoolToString(t *testing.T) {
	str := BoolToString(true)
	assert.Equal(t, "true", str)

	str = BoolToString(false)
	assert.Equal(t, "false", str)
}

func TestStringToBool(t *testing.T) {
	b := StringToBool("true")
	assert.True(t, b)

	b = StringToBool("false")
	assert.False(t, b)

	b = StringToBool("test")
	assert.False(t, b)
}

func TestPtr(t *testing.T) {
	strValue := "hello"
	strPtr := Ptr(strValue)
	assert.NotNil(t, strPtr)
	assert.Equal(t, strValue, *strPtr)

	boolValue := true
	boolPtr := Ptr(boolValue)
	assert.NotNil(t, boolPtr)
	assert.Equal(t, boolValue, *boolPtr)

	intValue := 42
	intPtr := Ptr(intValue)
	assert.NotNil(t, intPtr)
	assert.Equal(t, intValue, *intPtr)

	type TestStruct struct {
		Field string
	}
	structValue := TestStruct{Field: "test"}
	structPtr := Ptr(structValue)
	assert.NotNil(t, structPtr)
	assert.Equal(t, structValue, *structPtr)

	literalStrPtr := Ptr("literalHello")
	assert.NotNil(t, literalStrPtr)
	assert.Equal(t, "literalHello", *literalStrPtr)

	literalBoolPtr := Ptr(false)
	assert.NotNil(t, literalBoolPtr)
	assert.Equal(t, false, *literalBoolPtr)

	literalIntPtr := Ptr(123)
	assert.NotNil(t, literalIntPtr)
	assert.Equal(t, 123, *literalIntPtr)

	literalFloatPtr := Ptr(3.14)
	assert.NotNil(t, literalFloatPtr)
	assert.Equal(t, 3.14, *literalFloatPtr)

	literalStructPtr := Ptr(struct {
		Field1 int
		Field2 string
	}{
		Field1: 1,
		Field2: "test",
	})
	assert.NotNil(t, literalStructPtr)
	assert.Equal(t, struct {
		Field1 int
		Field2 string
	}{
		Field1: 1,
		Field2: "test",
	}, *literalStructPtr)
}
