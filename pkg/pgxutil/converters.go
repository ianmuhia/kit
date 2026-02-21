// Package pgxutil provides utility functions for working with pgx database types.
// This package contains helpers for converting between PostgreSQL types (pgtype)
// and standard Go types, particularly for handling nullable database values.
//
// The package provides two approaches:
//
// 1. Generic functions (ToPointer, ToValue, ToValueOrDefault, FromPointer, FromValue)
//   - Work with any type
//   - Reduce code duplication
//   - Type-safe at compile time
//
// 2. Type-specific wrapper functions (StringFromText, IntFromInt4, etc.)
//   - Backwards compatible
//   - More explicit for pgtype conversions
//   - Delegate to generic implementations internally
//
// Use the generic functions directly for maximum flexibility, or use the wrapper
// functions for clearer intent when working with specific pgtype types.
package pgxutil

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

//
// These generic functions reduce code duplication and work with any type.
// They provide the core conversion logic that type-specific functions use.

// NullableValue represents a wrapper for nullable database values with a common interface.
type NullableValue[T any] struct {
	Value T
	Valid bool
}

// ToPointer converts a nullable database value to a pointer, returning nil for NULL values.
// This is a generic helper that works with any type that has a Valid field and a value.
//
// Example:
//
//	name := pgxutil.ToPointer(row.UserName.Valid, row.UserName.String)
//	if name != nil {
//	    fmt.Println(*name)
//	}
//
//	age := pgxutil.ToPointer(row.Age.Valid, int(row.Age.Int32))
//	price := pgxutil.ToPointer(row.Price.Valid, row.Price.Float64)
func ToPointer[T any](valid bool, value T) *T {
	if !valid {
		return nil
	}
	return &value
}

// ToValue converts a nullable database value to its value type, returning a zero value for NULL.
// This is a generic helper that works with any type.
//
// Example:
//
//	name := pgxutil.ToValue(row.UserName.Valid, row.UserName.String)
//	fmt.Println(name) // prints "" if NULL
//
//	count := pgxutil.ToValue(row.Count.Valid, row.Count.Int32)
//	active := pgxutil.ToValue(row.IsActive.Valid, row.IsActive.Bool)
func ToValue[T any](valid bool, value T) T {
	if !valid {
		var zero T
		return zero
	}
	return value
}

// ToValueOrDefault converts a nullable database value to its value type with a custom default for NULL.
// This is a generic helper that works with any type.
//
// Example:
//
//	status := pgxutil.ToValueOrDefault(row.Status.Valid, row.Status.String, "pending")
//	priority := pgxutil.ToValueOrDefault(row.Priority.Valid, row.Priority.Int32, 1)
func ToValueOrDefault[T any](valid bool, value T, defaultValue T) T {
	if !valid {
		return defaultValue
	}
	return value
}

// FromPointer creates a nullable database value from a pointer.
// Returns a struct with Valid=false if the pointer is nil.
//
// Example:
//
//	var age *int = nil
//	result := pgxutil.FromPointer(age)
//	// result.Valid == false
//
//	name := "John"
//	result := pgxutil.FromPointer(&name)
//	// result.Valid == true, result.Value == "John"
func FromPointer[T any](ptr *T) NullableValue[T] {
	if ptr == nil {
		return NullableValue[T]{Valid: false}
	}
	return NullableValue[T]{Value: *ptr, Valid: true}
}

// FromValue creates a nullable database value from a value and a condition.
// Returns a struct with Valid based on the condition.
//
// Example:
//
//	result := pgxutil.FromValue("", func(s string) bool { return s != "" })
//	// result.Valid == false for empty string
//
//	result := pgxutil.FromValue(0, func(n int) bool { return n > 0 })
//	// result.Valid == false for zero or negative
func FromValue[T any](value T, isValid func(T) bool) NullableValue[T] {
	if !isValid(value) {
		return NullableValue[T]{Valid: false}
	}
	return NullableValue[T]{Value: value, Valid: true}
}

//
// These functions provide backwards compatibility and explicit type conversions.
// They all delegate to the generic functions internally.

// StringFromText converts pgtype.Text to *string, returning nil for NULL values.
// This is useful when working with nullable VARCHAR/TEXT columns in PostgreSQL.
//
// Example:
//
//	name := pgxutil.StringFromText(row.UserName)
//	if name != nil {
//	    fmt.Println(*name)
//	}
//
// Note: You can also use the generic ToPointer function:
//
//	name := pgxutil.ToPointer(row.UserName.Valid, row.UserName.String)
func StringFromText(t pgtype.Text) *string {
	return ToPointer(t.Valid, t.String)
}

// IntFromInt4 converts pgtype.Int4 to *int, returning nil for NULL values.
// This is useful when working with nullable INTEGER columns in PostgreSQL.
//
// Example:
//
//	age := pgxutil.IntFromInt4(row.UserAge)
//	if age != nil {
//	    fmt.Println(*age)
//	}
//
// Note: You can also use the generic ToPointer function:
//
//	age := pgxutil.ToPointer(row.UserAge.Valid, int(row.UserAge.Int32))
func IntFromInt4(val pgtype.Int4) *int {
	result := int(val.Int32)
	return ToPointer(val.Valid, result)
}

// Int64FromInt8 converts pgtype.Int8 to *int64, returning nil for NULL values.
// This is useful when working with nullable BIGINT columns in PostgreSQL.
//
// Example:
//
//	amount := pgxutil.Int64FromInt8(row.TotalAmount)
//	if amount != nil {
//	    fmt.Println(*amount)
//	}
//
// Note: You can also use the generic ToPointer function:
//
//	amount := pgxutil.ToPointer(row.TotalAmount.Valid, row.TotalAmount.Int64)
func Int64FromInt8(val pgtype.Int8) *int64 {
	return ToPointer(val.Valid, val.Int64)
}

// BoolFromBool converts pgtype.Bool to *bool, returning nil for NULL values.
// This is useful when working with nullable BOOLEAN columns in PostgreSQL.
//
// Example:
//
//	active := pgxutil.BoolFromBool(row.IsActive)
//	if active != nil {
//	    fmt.Println(*active)
//	}
//
// Note: You can also use the generic ToPointer function:
//
//	active := pgxutil.ToPointer(row.IsActive.Valid, row.IsActive.Bool)
func BoolFromBool(val pgtype.Bool) *bool {
	return ToPointer(val.Valid, val.Bool)
}

// Float64FromFloat8 converts pgtype.Float8 to *float64, returning nil for NULL values.
// This is useful when working with nullable DOUBLE PRECISION columns in PostgreSQL.
//
// Example:
//
//	price := pgxutil.Float64FromFloat8(row.UnitPrice)
//	if price != nil {
//	    fmt.Println(*price)
//	}
//
// Note: You can also use the generic ToPointer function:
//
//	price := pgxutil.ToPointer(row.UnitPrice.Valid, row.UnitPrice.Float64)
func Float64FromFloat8(val pgtype.Float8) *float64 {
	return ToPointer(val.Valid, val.Float64)
}

// Float64FromNumeric converts pgtype.Numeric (decimal.Decimal) to float64.
// This is useful when working with NUMERIC/DECIMAL columns in PostgreSQL.
// Note: This conversion may lose precision for very large or precise decimal values.
//
// Example:
//
//	amount := pgxutil.Float64FromNumeric(row.Amount)
//	fmt.Printf("Amount: %.2f\n", amount)
func Float64FromNumeric(d decimal.Decimal) float64 {
	f, _ := d.Float64()
	return f
}

// StringOrDefault returns the string value from pgtype.Text or a default value if NULL.
// This is useful when you need a non-pointer string with a fallback.
//
// Example:
//
//	status := pgxutil.StringOrDefault(row.Status, "pending")
//
// Note: You can also use the generic ToValueOrDefault function:
//
//	status := pgxutil.ToValueOrDefault(row.Status.Valid, row.Status.String, "pending")
func StringOrDefault(t pgtype.Text, defaultValue string) string {
	return ToValueOrDefault(t.Valid, t.String, defaultValue)
}

// IntOrDefault returns the int value from pgtype.Int4 or a default value if NULL.
// This is useful when you need a non-pointer int with a fallback.
//
// Example:
//
//	count := pgxutil.IntOrDefault(row.Count, 0)
//
// Note: You can also use the generic ToValueOrDefault function:
//
//	count := pgxutil.ToValueOrDefault(row.Count.Valid, int(row.Count.Int32), 0)
func IntOrDefault(val pgtype.Int4, defaultValue int) int {
	return ToValueOrDefault(val.Valid, int(val.Int32), defaultValue)
}

// BoolOrDefault returns the bool value from pgtype.Bool or a default value if NULL.
// This is useful when you need a non-pointer bool with a fallback.
//
// Example:
//
//	active := pgxutil.BoolOrDefault(row.IsActive, false)
//
// Note: You can also use the generic ToValueOrDefault function:
//
//	active := pgxutil.ToValueOrDefault(row.IsActive.Valid, row.IsActive.Bool, false)
func BoolOrDefault(val pgtype.Bool, defaultValue bool) bool {
	return ToValueOrDefault(val.Valid, val.Bool, defaultValue)
}

// TextFromString converts a Go string pointer to pgtype.Text.
// This is useful when inserting/updating nullable VARCHAR/TEXT columns.
//
// Example:
//
//	var name *string = nil
//	params := db.CreateUserParams{
//	    Name: pgxutil.TextFromString(name),
//	}
func TextFromString(s string) pgtype.Text {
	if s == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: s, Valid: true}
}

// TextFromStringPtr converts a Go string pointer to pgtype.Text.
// This is useful when inserting/updating nullable VARCHAR/TEXT columns.
//
// Example:
//
//	var name *string = nil
//	params := db.CreateUserParams{
//	    Name: pgxutil.TextFromStringPtr(name),
//	}
func TextFromStringPtr(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

// Int4FromInt converts a Go int pointer to pgtype.Int4.
// This is useful when inserting/updating nullable INTEGER columns.
//
// Example:
//
//	var age *int = nil
//	params := db.CreateUserParams{
//	    Age: pgxutil.Int4FromInt(age),
//	}
func Int4FromInt(i *int) pgtype.Int4 {
	if i == nil {
		return pgtype.Int4{Valid: false}
	}
	return pgtype.Int4{Int32: int32(*i), Valid: true}
}

// Int8FromInt64 converts a Go int64 pointer to pgtype.Int8.
// This is useful when inserting/updating nullable BIGINT columns.
//
// Example:
//
//	var amount *int64 = nil
//	params := db.CreateTransactionParams{
//	    Amount: pgxutil.Int8FromInt64(amount),
//	}
func Int8FromInt64(i *int64) pgtype.Int8 {
	if i == nil {
		return pgtype.Int8{Valid: false}
	}
	return pgtype.Int8{Int64: *i, Valid: true}
}

// BoolFromBoolPtr converts a Go bool pointer to pgtype.Bool.
// This is useful when inserting/updating nullable BOOLEAN columns.
//
// Example:
//
//	var active *bool = nil
//	params := db.UpdateUserParams{
//	    IsActive: pgxutil.BoolFromBoolPtr(active),
//	}
func BoolFromBoolPtr(b *bool) pgtype.Bool {
	if b == nil {
		return pgtype.Bool{Valid: false}
	}
	return pgtype.Bool{Bool: *b, Valid: true}
}

// Float8FromFloat64 converts a Go float64 pointer to pgtype.Float8.
// This is useful when inserting/updating nullable DOUBLE PRECISION columns.
//
// Example:
//
//	var price *float64 = nil
//	params := db.UpdateProductParams{
//	    Price: pgxutil.Float8FromFloat64(price),
//	}
func Float8FromFloat64(f *float64) pgtype.Float8 {
	if f == nil {
		return pgtype.Float8{Valid: false}
	}
	return pgtype.Float8{Float64: *f, Valid: true}
}

// TimestampFromTime converts a time.Time pointer to pgtype.Timestamptz.
// This is useful when inserting/updating nullable TIMESTAMP columns.
//
// Example:
//
//	var createdAt *time.Time = nil
//	params := db.CreateRecordParams{
//	    CreatedAt: pgxutil.TimestampFromTime(createdAt),
//	}
func TimestampFromTime(t *time.Time) pgtype.Timestamptz {
	if t == nil {
		return pgtype.Timestamptz{Valid: false}
	}
	return pgtype.Timestamptz{Time: *t, Valid: true}
}

// TimeFromTimestamp converts pgtype.Timestamptz to *time.Time.
// This is useful when reading nullable TIMESTAMP columns.
//
// Example:
//
//	createdAt := pgxutil.TimeFromTimestamp(row.CreatedAt)
//	if createdAt != nil {
//	    fmt.Println(createdAt.Format(time.RFC3339))
//	}
//
// Note: You can also use the generic ToPointer function:
//
//	createdAt := pgxutil.ToPointer(row.CreatedAt.Valid, row.CreatedAt.Time)
func TimeFromTimestamp(t pgtype.Timestamptz) *time.Time {
	return ToPointer(t.Valid, t.Time)
}

// PgxTextToString converts pgtype.Text to string, returning empty string for NULL values.
// This is useful when working with VARCHAR/TEXT columns in PostgreSQL.
//
// Example:
//
//	name := pgxutil.PgxTextToString(row.UserName)
//
// Note: You can also use the generic ToValue function:
//
//	name := pgxutil.ToValue(row.UserName.Valid, row.UserName.String)
func PgxTextToString(t pgtype.Text) string {
	return ToValue(t.Valid, t.String)
}

// PgxInt4ToInt32 converts pgtype.Int4 to int32, returning 0 for NULL values.
// This is useful when working with INTEGER columns in PostgreSQL.
//
// Example:
//
//	count := pgxutil.PgxInt4ToInt32(row.Count)
//
// Note: You can also use the generic ToValue function:
//
//	count := pgxutil.ToValue(row.Count.Valid, row.Count.Int32)
func PgxInt4ToInt32(val pgtype.Int4) int32 {
	return ToValue(val.Valid, val.Int32)
}

func Int32ToPgxInt4(val int32) pgtype.Int4 {
	return pgtype.Int4{Int32: val, Valid: true}
}

// PgxBoolToBool converts pgtype.Bool to bool, returning false for NULL values.
// This is useful when working with BOOLEAN columns in PostgreSQL.
//
// Example:
//
//	active := pgxutil.PgxBoolToBool(row.IsActive)
//
// Note: You can also use the generic ToValue function:
//
//	active := pgxutil.ToValue(row.IsActive.Valid, row.IsActive.Bool)
func PgxBoolToBool(val pgtype.Bool) bool {
	return ToValue(val.Valid, val.Bool)
}

// PgxTimestamptzToTime converts pgtype.Timestamptz to time.Time, returning zero time for NULL values.
// This is useful when working with TIMESTAMP columns in PostgreSQL.
//
// Example:
//
//	createdAt := pgxutil.PgxTimestamptzToTime(row.CreatedAt)
//
// Note: You can also use the generic ToValue function:
//
//	createdAt := pgxutil.ToValue(row.CreatedAt.Valid, row.CreatedAt.Time)
func PgxTimestamptzToTime(t pgtype.Timestamptz) time.Time {
	return ToValue(t.Valid, t.Time)
}
// PgxTimestamptzToTimePtr converts pgtype.Timestamptz to *time.Time, returning nil for NULL values.
// This is useful when working with nullable TIMESTAMP columns in PostgreSQL.
//
// Example:
//
//	createdAt := pgxutil.PgxTimestamptzToTimePtr(row.CreatedAt)
//	if createdAt != nil {
//	    fmt.Println(createdAt.Format(time.RFC3339))
//	}
//
// Note: You can also use the generic ToPointer function:
//
//	createdAt := pgxutil.ToPointer(row.CreatedAt.Valid, row.CreatedAt.Time)
func PgxTimestamptzToTimePtr(t pgtype.Timestamptz) *time.Time {
	return ToPointer(t.Valid, t.Time)
}

// UUIDToPgUUID converts a uuid.UUID to pgtype.UUID.
// This is useful when inserting/updating UUID columns in PostgreSQL.
//
// Example:
//
//	id := uuid.New()
//	params := db.CreateUserParams{
//	    ID: pgxutil.UUIDToPgUUID(id),
//	}
func UUIDToPgUUID(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: id, Valid: true}
}

// PgUUIDToUUID converts pgtype.UUID to uuid.UUID, returning uuid.Nil for NULL values.
// This is useful when reading UUID columns from PostgreSQL.
//
// Example:
//
//	userID := pgxutil.PgUUIDToUUID(row.UserID)
//	if userID == uuid.Nil {
//	    fmt.Println("User ID is NULL")
//	}
//
// Note: You can also use the generic ToValue function:
//
//	userID := pgxutil.ToValue(row.UserID.Valid, row.UserID.Bytes)
func PgUUIDToUUID(pgUUID pgtype.UUID) uuid.UUID {
	return ToValue(pgUUID.Valid, pgUUID.Bytes)
}

// StringToPgText converts a Go string to pgtype.Text, treating empty strings as NULL.
// This is useful when inserting/updating VARCHAR/TEXT columns in PostgreSQL.
// Note: Empty strings are converted to NULL. Use TextFromString if you need different behavior.
//
// Example:
//
//	params := db.UpdateUserParams{
//	    Bio: pgxutil.StringToPgText(userBio),
//	}
func StringToPgText(s string) pgtype.Text {
	if s == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: s, Valid: true}
}

// PgTextToString converts pgtype.Text to string, returning empty string for NULL values.
// This is useful when reading VARCHAR/TEXT columns from PostgreSQL.
// Note: This is an alias for PgxTextToString for consistency.
//
// Example:
//
//	name := pgxutil.PgTextToString(row.UserName)
//	if name == "" {
//	    fmt.Println("Name is NULL or empty")
//	}
//
// Note: You can also use the generic ToValue function:
//
//	name := pgxutil.ToValue(row.UserName.Valid, row.UserName.String)
func PgTextToString(pgText pgtype.Text) string {
	return ToValue(pgText.Valid, pgText.String)
}

// PgTextToStringPtr converts pgtype.Text to *string, returning nil for NULL values.
// This is useful when reading nullable VARCHAR/TEXT columns from PostgreSQL.
// Note: This is an alias for StringFromText for consistency.
//
// Example:
//
//	name := pgxutil.PgTextToStringPtr(row.UserName)
//	if name != nil {
//	    fmt.Println(*name)
//	}
//
// Note: You can also use the generic ToPointer function:
//
//	name := pgxutil.ToPointer(row.UserName.Valid, row.UserName.String)
func PgTextToStringPtr(pgText pgtype.Text) *string {
	return ToPointer(pgText.Valid, pgText.String)
}

// PgNumericToFloat64 converts pgtype.Numeric to float64, returning 0 for NULL values.
// This is useful when reading NUMERIC/DECIMAL columns from PostgreSQL.
// Note: This conversion may lose precision for very large or precise decimal values.
//
// Example:
//
//	price := pgxutil.PgNumericToFloat64(row.Price)
//	fmt.Printf("Price: %.2f\n", price)
func PgNumericToFloat64(n pgtype.Numeric) float64 {
	if !n.Valid {
		return 0
	}
	f, _ := n.Float64Value()
	return f.Float64
}

// PgUUIDToUUIDPtr converts a pgtype.UUID to a *uuid.UUID, returning nil for NULL values.
//
// Example:
//   var pgUUID pgtype.UUID
//   var uuidPtr *uuid.UUID = PgUUIDToUUIDPtr(pgUUID)
func PgUUIDToUUIDPtr(pgUUID pgtype.UUID) *uuid.UUID {
    if !pgUUID.Valid {
        return nil
    }
    return uuidPtr(pgUUID.UUID)
}

// UUIDToPgUUIDPtr converts a *uuid.UUID pointer to pgtype.UUID. If the pointer is nil, it returns a pgtype.UUID with Valid=false.
func UUIDToPgUUIDPtr(u *uuid.UUID) pgtype.UUID {
    if u == nil {
        return pgtype.UUID{Valid: false}
    }
    return pgtype.UUID{Bytes: u[:], Valid: true}
}
