package pgxutil

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestBoolFromBool(t *testing.T) {
	t.Run("valid true", func(t *testing.T) {
		val := pgtype.Bool{Bool: true, Valid: true}
		result := BoolFromBool(val)
		assert.NotNil(t, result)
		assert.True(t, *result)
	})

	t.Run("valid false", func(t *testing.T) {
		val := pgtype.Bool{Bool: false, Valid: true}
		result := BoolFromBool(val)
		assert.NotNil(t, result)
		assert.False(t, *result)
	})

	t.Run("null bool", func(t *testing.T) {
		val := pgtype.Bool{Valid: false}
		result := BoolFromBool(val)
		assert.Nil(t, result)
	})
}

func TestFloat64FromFloat8(t *testing.T) {
	t.Run("valid float", func(t *testing.T) {
		val := pgtype.Float8{Float64: 3.14, Valid: true}
		result := Float64FromFloat8(val)
		assert.NotNil(t, result)
		assert.InDelta(t, 3.14, *result, 0.001)
	})

	t.Run("null float", func(t *testing.T) {
		val := pgtype.Float8{Valid: false}
		result := Float64FromFloat8(val)
		assert.Nil(t, result)
	})
}

func TestPgNumericToFloat64(t *testing.T) {
	t.Run("null numeric returns 0", func(t *testing.T) {
		result := PgNumericToFloat64(pgtype.Numeric{Valid: false})
		assert.Equal(t, float64(0), result)
	})

	t.Run("valid numeric", func(t *testing.T) {
		n := pgtype.Numeric{}
		require.NoError(t, n.Scan("99.99"))
		result := PgNumericToFloat64(n)
		assert.InDelta(t, 99.99, result, 0.001)
	})
}

func TestPgNumericToFloat64E(t *testing.T) {
	t.Run("null numeric returns 0 with no error", func(t *testing.T) {
		result, err := PgNumericToFloat64E(pgtype.Numeric{Valid: false})
		assert.NoError(t, err)
		assert.Equal(t, float64(0), result)
	})

	t.Run("valid numeric returns value", func(t *testing.T) {
		n := pgtype.Numeric{}
		require.NoError(t, n.Scan("12.345"))
		result, err := PgNumericToFloat64E(n)
		assert.NoError(t, err)
		assert.InDelta(t, 12.345, result, 0.001)
	})
}

func TestTimeFromTimestamp(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Microsecond)

	t.Run("valid timestamp", func(t *testing.T) {
		val := pgtype.Timestamptz{Time: now, Valid: true}
		result := TimeFromTimestamp(val)
		assert.NotNil(t, result)
		assert.Equal(t, now, *result)
	})

	t.Run("null timestamp", func(t *testing.T) {
		val := pgtype.Timestamptz{Valid: false}
		result := TimeFromTimestamp(val)
		assert.Nil(t, result)
	})
}

func TestUUIDConversions(t *testing.T) {
	id := uuid.New()

	t.Run("UUIDToPgUUID round-trip", func(t *testing.T) {
		pg := UUIDToPgUUID(id)
		assert.True(t, pg.Valid)
		back := PgUUIDToUUID(pg)
		assert.Equal(t, id, back)
	})

	t.Run("UUIDToPgUUIDPtr with value", func(t *testing.T) {
		pg := UUIDToPgUUIDPtr(&id)
		assert.True(t, pg.Valid)
	})

	t.Run("UUIDToPgUUIDPtr with nil", func(t *testing.T) {
		pg := UUIDToPgUUIDPtr(nil)
		assert.False(t, pg.Valid)
	})

	t.Run("PgUUIDToUUIDPtr with value", func(t *testing.T) {
		pg := pgtype.UUID{Bytes: id, Valid: true}
		result := PgUUIDToUUIDPtr(pg)
		assert.NotNil(t, result)
		assert.Equal(t, id, *result)
	})

	t.Run("PgUUIDToUUIDPtr null returns nil", func(t *testing.T) {
		pg := pgtype.UUID{Valid: false}
		result := PgUUIDToUUIDPtr(pg)
		assert.Nil(t, result)
	})
}

func TestDecimalConversions(t *testing.T) {
	d := decimal.NewFromFloat(123.45)

	t.Run("NumericFromDecimal round-trip", func(t *testing.T) {
		n := NumericFromDecimal(d)
		assert.True(t, n.Valid)
		back := DecimalFromNumeric(n)
		assert.True(t, d.Equal(back))
	})

	t.Run("NumericFromDecimalPtr with nil", func(t *testing.T) {
		n := NumericFromDecimalPtr(nil)
		assert.False(t, n.Valid)
	})

	t.Run("DecimalFromNumericPtr null returns nil", func(t *testing.T) {
		result := DecimalFromNumericPtr(pgtype.Numeric{Valid: false})
		assert.Nil(t, result)
	})

	t.Run("DecimalFromNumericPtr with value", func(t *testing.T) {
		n := NumericFromDecimal(d)
		result := DecimalFromNumericPtr(n)
		assert.NotNil(t, result)
		assert.True(t, d.Equal(*result))
	})
}

func TestIntervalConversions(t *testing.T) {
	dur := 3*time.Hour + 30*time.Minute

	t.Run("IntervalFromDuration round-trip", func(t *testing.T) {
		iv := IntervalFromDuration(dur)
		assert.True(t, iv.Valid)
		back := DurationFromInterval(iv)
		assert.Equal(t, dur, back)
	})

	t.Run("DurationFromInterval null returns 0", func(t *testing.T) {
		result := DurationFromInterval(pgtype.Interval{Valid: false})
		assert.Equal(t, time.Duration(0), result)
	})

	t.Run("IntervalFromDurationPtr with nil", func(t *testing.T) {
		iv := IntervalFromDurationPtr(nil)
		assert.False(t, iv.Valid)
	})

	t.Run("IntervalFromDurationPtr with value", func(t *testing.T) {
		iv := IntervalFromDurationPtr(&dur)
		assert.True(t, iv.Valid)
		assert.Equal(t, dur, DurationFromInterval(iv))
	})
}

func TestSliceConversions(t *testing.T) {
	t.Run("Int4SliceFromInts round-trip", func(t *testing.T) {
		ints := []int32{1, 2, 3}
		arr := Int4SliceFromInts(ints)
		back := IntsFromInt4Slice(arr)
		assert.Equal(t, ints, back)
	})

	t.Run("Int4SliceFromInts nil", func(t *testing.T) {
		arr := Int4SliceFromInts(nil)
		assert.Empty(t, arr)
	})

	t.Run("TextSliceFromStrings round-trip", func(t *testing.T) {
		words := []string{"a", "b", "c"}
		arr := TextSliceFromStrings(words)
		back := StringsFromTextSlice(arr)
		assert.Equal(t, words, back)
	})

	t.Run("TextSliceFromStrings nil", func(t *testing.T) {
		arr := TextSliceFromStrings(nil)
		assert.Empty(t, arr)
	})
}

func TestOrDefaultFunctions(t *testing.T) {
	t.Run("StringOrDefault present", func(t *testing.T) {
		assert.Equal(t, "hello", StringOrDefault(pgtype.Text{String: "hello", Valid: true}, "default"))
	})

	t.Run("StringOrDefault null", func(t *testing.T) {
		assert.Equal(t, "default", StringOrDefault(pgtype.Text{Valid: false}, "default"))
	})

	t.Run("IntOrDefault present", func(t *testing.T) {
		assert.Equal(t, 7, IntOrDefault(pgtype.Int4{Int32: 7, Valid: true}, 0))
	})

	t.Run("IntOrDefault null", func(t *testing.T) {
		assert.Equal(t, 99, IntOrDefault(pgtype.Int4{Valid: false}, 99))
	})

	t.Run("BoolOrDefault present", func(t *testing.T) {
		assert.True(t, BoolOrDefault(pgtype.Bool{Bool: true, Valid: true}, false))
	})

	t.Run("BoolOrDefault null", func(t *testing.T) {
		assert.True(t, BoolOrDefault(pgtype.Bool{Valid: false}, true))
	})
}
