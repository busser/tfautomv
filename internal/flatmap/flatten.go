package flatmap

import (
	"fmt"
	"reflect"
)

// Flatten takes any object and turns it into a flat map[string]interface{}.
//
// With "obj", map keys must be strings. Values must be slices, maps,
// primitives, or any combination of those together.
func Flatten(obj map[string]interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	for k, raw := range obj {
		err := flatten(result, k, reflect.ValueOf(raw))
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func flatten(result map[string]interface{}, prefix string, v reflect.Value) error {
	if v.Kind() == reflect.Interface {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Map:
		err := flattenMap(result, prefix, v)
		if err != nil {
			return err
		}
	case reflect.Slice:
		err := flattenSlice(result, prefix, v)
		if err != nil {
			return err
		}
	default:
		if !v.IsValid() { // nil values
			result[prefix] = nil
			return nil
		}
		result[prefix] = v.Interface()
	}

	return nil
}

func flattenMap(result map[string]interface{}, prefix string, v reflect.Value) error {
	for _, k := range v.MapKeys() {
		if k.Kind() != reflect.String {
			return fmt.Errorf("%s: map key is not string: %s", prefix, k)
		}

		flatten(result, fmt.Sprintf("%s.%s", prefix, k.String()), v.MapIndex(k))
	}

	return nil
}

func flattenSlice(result map[string]interface{}, prefix string, v reflect.Value) error {
	prefix = prefix + "."

	result[prefix+"#"] = v.Len()
	for i := 0; i < v.Len(); i++ {
		flatten(result, fmt.Sprintf("%s%d", prefix, i), v.Index(i))
	}

	return nil
}
