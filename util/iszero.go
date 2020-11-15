package util

import "reflect"

// IsZero checks if the passed value is a zero value
func IsZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Func, reflect.Map, reflect.Slice:
		return v.IsNil()
	case reflect.Array:
		z := true
		for i := 0; i < v.Len(); i++ {
			z = z && IsZero(v.Index(i))
		}
		return z
	case reflect.Struct:
		z := true
		for i := 0; i < v.NumField(); i++ {
			if v.Field(i).CanSet() {
				z = z && IsZero(v.Field(i))
			}
		}
		return z
	case reflect.Ptr:
		return IsZero(reflect.Indirect(v))
	}

	// Compare other types directly:
	z := reflect.Zero(v.Type())
	result := v.Interface() == z.Interface()

	return result
}
