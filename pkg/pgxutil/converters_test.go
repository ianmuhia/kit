package pgxutil

import (
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func TestStringFromText(t *testing.T) {
	t.Run("valid text", func(t *testing.T) {
		text := pgtype.Text{String: "hello", Valid: true}
		result := StringFromText(text)
		assert.NotNil(t, result)
		assert.Equal(t, "hello", *result)
	})

	t.Run("null text", func(t *testing.T) {
		text := pgtype.Text{Valid: false}
		result := StringFromText(text)
		assert.Nil(t, result)
	})
}

func TestIntFromInt4(t *testing.T) {
	t.Run("valid int", func(t *testing.T) {
		val := pgtype.Int4{Int32: 42, Valid: true}
		result := IntFromInt4(val)
		assert.NotNil(t, result)
		assert.Equal(t, 42, *result)
	})

	t.Run("null int", func(t *testing.T) {
		val := pgtype.Int4{Valid: false}
		result := IntFromInt4(val)
		assert.Nil(t, result)
	})
}

func TestInt64FromInt8(t *testing.T) {
	t.Run("valid int64", func(t *testing.T) {
		val := pgtype.Int8{Int64: 1234567890, Valid: true}
		result := Int64FromInt8(val)
		assert.NotNil(t, result)
		assert.Equal(t, int64(1234567890), *result)
	})

	t.Run("null int64", func(t *testing.T) {
		val := pgtype.Int8{Valid: false}
		result := Int64FromInt8(val)
		assert.Nil(t, result)
	})
}

// func TestBoolFromBool(t *testing.T) {
// 	t.Run("valid true", func(t *testing.T) {
// 		val := pgtype.Bool{Bool: true, Valid: true}
// 		result := BoolFromBool(val)
