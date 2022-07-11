package util

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func ValidatePassword(password string, hashedPassword string) bool {
	// Comparing the password with the hash
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func GenerateHash(value string) (string, error) {
	// Generate "hash" to store from user password
	hash, err := bcrypt.GenerateFromPassword([]byte(value), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func ValidateEmail(Email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(Email)
}

// Bool stores v in a new bool value and returns a pointer to it.
func Bool(v bool) *bool { return &v }

// Int32 stores v in a new int value and returns a pointer to it.
func Int(v int) *int { return &v }

// Int32 stores v in a new int32 value and returns a pointer to it.
func Int32(v int32) *int32 { return &v }

// Int64 stores v in a new int64 value and returns a pointer to it.
func Int64(v int64) *int64 { return &v }

// Float32 stores v in a new float32 value and returns a pointer to it.
func Float32(v float32) *float32 { return &v }

// Float64 stores v in a new float64 value and returns a pointer to it.
func Float64(v float64) *float64 { return &v }

// Uint32 stores v in a new uint32 value and returns a pointer to it.
func Uint32(v uint32) *uint32 { return &v }

// Uint64 stores v in a new uint64 value and returns a pointer to it.
func Uint64(v uint64) *uint64 { return &v }

// String stores v in a new string value and returns a pointer to it.
func String(v string) *string { return &v }

func IsNull(v interface{}, t string) bool {
	if t == "string" && v.(string) == "" {
		return true
	} else if t == "int" && v == 0 {
		return true
	} else if t == "datetime" && v.(time.Time).IsZero() {
		return true
	}
	return false
}

func IsValidDataType(dataType string, v interface{}) bool {
	return GetDataType(v) == dataType
}

func GetDataType(v interface{}) string {
	switch v.(type) {
	case bool:
		return "bool"
	case int32:
		return "int32"
	case int64:
		return "int64"
	case float32:
		return "float32"
	case float64:
		return "float64"
	case uint32:
		return "uint32"
	case uint64:
		return "uint64"
	case string:
		return "string"
	}
	return ""
}

func InterfaceToInt(i interface{}) (*int, error) {
	switch v := i.(type) {
	case int:
		return &v, nil
	case string:
		v1ID, err := strconv.Atoi(i.(string))
		if err != nil {
			return nil, fmt.Errorf("wrong type")
		}
		return &v1ID, nil
	case float64:
		val := int(v)
		return &val, nil
	}
	return nil, fmt.Errorf("can't parse, please check request")

}

func GetEnvWithKey(key string) string {
	return os.Getenv(key)
}
