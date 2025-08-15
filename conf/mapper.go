package conf

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// mapToStruct maps rawMap to struct using reflection
func mapToStruct(rawMap map[string]any, target any, option *Option, isUpdate bool) error {
	return mapToStructWithPath(rawMap, target, option, "", isUpdate, false)
}

// mapToStructWithPath maps rawMap to struct with field path tracking
func mapToStructWithPath(rawMap map[string]any, target any, option *Option, basePath string, isUpdate bool, parentOptional bool) error {
	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("target must be a pointer to struct")
	}

	v = v.Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// Skip unexported fields
		if !field.CanSet() {
			continue
		}

		// Parse tag
		tagValue := fieldType.Tag.Get(option.TagName)
		tagInfo := parseTag(tagValue)

		// Skip field if tag says so
		if tagInfo.Skip {
			continue
		}

		// Determine field name
		fieldName := fieldType.Name
		if tagInfo.FieldName != "" {
			fieldName = tagInfo.FieldName
		}

		// Build field path
		fieldPath := fieldName
		if basePath != "" {
			fieldPath = basePath + "." + fieldName
		}

		// Handle anonymous struct (embedded)
		if fieldType.Anonymous {
			if field.Kind() == reflect.Struct {
				// Use current rawMap for anonymous struct
				if err := mapToStructWithPath(rawMap, field.Addr().Interface(), option, basePath, isUpdate, parentOptional || tagInfo.Optional); err != nil {
					return err
				}
			}
			continue
		}

		// Get value from map
		value, exists := findValueInMap(rawMap, fieldName, option.MatchMode)
		if !exists {
			// For update mode, skip missing fields
			if isUpdate {
				continue
			}

			// Handle struct fields
			if field.Kind() == reflect.Struct || (field.Kind() == reflect.Ptr && field.Type().Elem().Kind() == reflect.Struct) {
				if err := handleStructField(field, fieldType, option, fieldPath, isUpdate, parentOptional || tagInfo.Optional); err != nil {
					return err
				}
				continue
			}

			// Use default value if available
			if tagInfo.Default != "" {
				processedDefault, err := processEnvVars(tagInfo.Default, option.UseEnv)
				if err != nil {
					return fmt.Errorf("field %s default value error: %w", fieldPath, err)
				}
				value = processedDefault
				exists = true
			} else {
				// Check if field is required
				if !tagInfo.Optional && !parentOptional {
					return fmt.Errorf("field %s is required but not found", fieldPath)
				}
				continue
			}
		}

		// Process environment variables in value
		processedValue, err := processEnvValue(value, option.UseEnv)
		if err != nil {
			return fmt.Errorf("field %s environment variable error: %w", fieldPath, err)
		}

		// Validate value
		if err := validateValue(processedValue, tagInfo, fieldPath); err != nil {
			return err
		}

		// Handle watch callback for updates
		if isUpdate && tagInfo.Watch && option.WatchCallback != nil {
			oldValue := field.Interface()
			if err := option.WatchCallback(fieldPath, fieldName, oldValue, processedValue); err != nil {
				return fmt.Errorf("watch callback error for field %s: %w", fieldPath, err)
			}
		}

		// Handle struct fields with value
		if field.Kind() == reflect.Struct || (field.Kind() == reflect.Ptr && field.Type().Elem().Kind() == reflect.Struct) {
			if valueMap, ok := processedValue.(map[string]any); ok {
				if field.Kind() == reflect.Ptr {
					if field.IsNil() {
						field.Set(reflect.New(field.Type().Elem()))
					}
					if err := mapToStructWithPath(valueMap, field.Interface(), option, fieldPath, isUpdate, parentOptional || tagInfo.Optional); err != nil {
						return err
					}
				} else {
					if err := mapToStructWithPath(valueMap, field.Addr().Interface(), option, fieldPath, isUpdate, parentOptional || tagInfo.Optional); err != nil {
						return err
					}
				}
			} else {
				return fmt.Errorf("field %s expected map for struct, got %T", fieldPath, processedValue)
			}
		} else {
			// Set field value for non-struct types
			if err := setFieldValue(field, processedValue, fieldPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// handleStructField handles struct and pointer to struct fields
func handleStructField(field reflect.Value, fieldType reflect.StructField, option *Option, fieldPath string, isUpdate bool, parentOptional bool) error {
	if field.Kind() == reflect.Ptr {
		// Handle pointer to struct
		if field.IsNil() {
			// Create new instance
			newStruct := reflect.New(field.Type().Elem())
			field.Set(newStruct)
		}
		// Process the struct that pointer points to
		return mapToStructWithPath(make(map[string]any), field.Interface(), option, fieldPath, isUpdate, parentOptional)
	} else if field.Kind() == reflect.Struct {
		// Handle direct struct
		return mapToStructWithPath(make(map[string]any), field.Addr().Interface(), option, fieldPath, isUpdate, parentOptional)
	}
	return nil
}

// setFieldValue sets field value with type conversion
func setFieldValue(field reflect.Value, value any, fieldPath string) error {
	if value == nil {
		return nil
	}

	fieldType := field.Type()
	valueType := reflect.TypeOf(value)

	// Handle time.Duration specially
	if fieldType == reflect.TypeOf(time.Duration(0)) {
		duration, err := parseDuration(value)
		if err != nil {
			return fmt.Errorf("field %s duration parse error: %w", fieldPath, err)
		}
		field.Set(reflect.ValueOf(duration))
		return nil
	}

	// Handle pointer types
	if fieldType.Kind() == reflect.Ptr {
		if field.IsNil() {
			field.Set(reflect.New(fieldType.Elem()))
		}
		return setFieldValue(field.Elem(), value, fieldPath)
	}

	// Handle struct types
	if fieldType.Kind() == reflect.Struct {
		if valueMap, ok := value.(map[string]any); ok {
			// This should be handled by mapToStructWithPath
			return fmt.Errorf("struct field %s should be handled by mapToStructWithPath, got map: %v", fieldPath, valueMap)
		}
	}

	// Direct assignment if types match
	if valueType.AssignableTo(fieldType) {
		field.Set(reflect.ValueOf(value))
		return nil
	}

	// Type conversion
	return convertAndSetValue(field, value, fieldPath)
}

// convertAndSetValue converts value to field type and sets it
func convertAndSetValue(field reflect.Value, value any, fieldPath string) error {
	fieldType := field.Type()

	switch fieldType.Kind() {
	case reflect.String:
		if str, ok := value.(string); ok {
			field.SetString(str)
		} else {
			field.SetString(fmt.Sprintf("%v", value))
		}

	case reflect.Bool:
		if b, ok := value.(bool); ok {
			field.SetBool(b)
		} else if str, ok := value.(string); ok {
			if boolVal, err := strconv.ParseBool(str); err == nil {
				field.SetBool(boolVal)
			} else {
				return fmt.Errorf("field %s cannot convert '%v' to bool", fieldPath, value)
			}
		} else {
			return fmt.Errorf("field %s cannot convert '%v' to bool", fieldPath, value)
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if err := setIntValue(field, value, fieldPath); err != nil {
			return err
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if err := setUintValue(field, value, fieldPath); err != nil {
			return err
		}

	case reflect.Float32, reflect.Float64:
		if err := setFloatValue(field, value, fieldPath); err != nil {
			return err
		}

	case reflect.Slice:
		return setSliceValue(field, value, fieldPath)

	case reflect.Map:
		return setMapValue(field, value, fieldPath)

	default:
		return fmt.Errorf("field %s unsupported type conversion from %T to %s", fieldPath, value, fieldType.Kind())
	}

	return nil
}

// setIntValue sets integer value
func setIntValue(field reflect.Value, value any, fieldPath string) error {
	var intVal int64
	var err error

	switch v := value.(type) {
	case int:
		intVal = int64(v)
	case int8:
		intVal = int64(v)
	case int16:
		intVal = int64(v)
	case int32:
		intVal = int64(v)
	case int64:
		intVal = v
	case uint:
		intVal = int64(v)
	case uint8:
		intVal = int64(v)
	case uint16:
		intVal = int64(v)
	case uint32:
		intVal = int64(v)
	case uint64:
		intVal = int64(v)
	case float32:
		intVal = int64(v)
	case float64:
		intVal = int64(v)
	case string:
		intVal, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return fmt.Errorf("field %s cannot convert '%v' to int", fieldPath, value)
		}
	default:
		return fmt.Errorf("field %s cannot convert '%v' to int", fieldPath, value)
	}

	if field.OverflowInt(intVal) {
		return fmt.Errorf("field %s value %v overflows %s", fieldPath, value, field.Type().Kind())
	}
	field.SetInt(intVal)
	return nil
}

// setUintValue sets unsigned integer value
func setUintValue(field reflect.Value, value any, fieldPath string) error {
	var uintVal uint64
	var err error

	switch v := value.(type) {
	case uint:
		uintVal = uint64(v)
	case uint8:
		uintVal = uint64(v)
	case uint16:
		uintVal = uint64(v)
	case uint32:
		uintVal = uint64(v)
	case uint64:
		uintVal = v
	case int:
		if v < 0 {
			return fmt.Errorf("field %s cannot convert negative value %v to uint", fieldPath, value)
		}
		uintVal = uint64(v)
	case int8, int16, int32, int64:
		intVal := reflect.ValueOf(v).Int()
		if intVal < 0 {
			return fmt.Errorf("field %s cannot convert negative value %v to uint", fieldPath, value)
		}
		uintVal = uint64(intVal)
	case float32, float64:
		floatVal := reflect.ValueOf(v).Float()
		if floatVal < 0 {
			return fmt.Errorf("field %s cannot convert negative value %v to uint", fieldPath, value)
		}
		uintVal = uint64(floatVal)
	case string:
		uintVal, err = strconv.ParseUint(v, 10, 64)
		if err != nil {
			return fmt.Errorf("field %s cannot convert '%v' to uint", fieldPath, value)
		}
	default:
		return fmt.Errorf("field %s cannot convert '%v' to uint", fieldPath, value)
	}

	if field.OverflowUint(uintVal) {
		return fmt.Errorf("field %s value %v overflows %s", fieldPath, value, field.Type().Kind())
	}
	field.SetUint(uintVal)
	return nil
}

// setFloatValue sets floating point value
func setFloatValue(field reflect.Value, value any, fieldPath string) error {
	var floatVal float64
	var err error

	switch v := value.(type) {
	case float32:
		floatVal = float64(v)
	case float64:
		floatVal = v
	case int, int8, int16, int32, int64:
		floatVal = float64(reflect.ValueOf(v).Int())
	case uint, uint8, uint16, uint32, uint64:
		floatVal = float64(reflect.ValueOf(v).Uint())
	case string:
		floatVal, err = strconv.ParseFloat(v, 64)
		if err != nil {
			return fmt.Errorf("field %s cannot convert '%v' to float", fieldPath, value)
		}
	default:
		return fmt.Errorf("field %s cannot convert '%v' to float", fieldPath, value)
	}

	if field.OverflowFloat(floatVal) {
		return fmt.Errorf("field %s value %v overflows %s", fieldPath, value, field.Type().Kind())
	}
	field.SetFloat(floatVal)
	return nil
}

// setSliceValue sets slice field value
func setSliceValue(field reflect.Value, value any, fieldPath string) error {
	valueSlice, ok := value.([]any)
	if !ok {
		return fmt.Errorf("field %s expected slice, got %T", fieldPath, value)
	}

	sliceType := field.Type()

	// Pre-allocate slice capacity
	newSlice := reflect.MakeSlice(sliceType, len(valueSlice), len(valueSlice))

	for i, item := range valueSlice {
		elem := newSlice.Index(i)
		if err := setFieldValue(elem, item, fmt.Sprintf("%s[%d]", fieldPath, i)); err != nil {
			return err
		}
	}

	field.Set(newSlice)
	return nil
}

// setMapValue sets map field value
func setMapValue(field reflect.Value, value any, fieldPath string) error {
	valueMap, ok := value.(map[string]any)
	if !ok {
		return fmt.Errorf("field %s expected map, got %T", fieldPath, value)
	}

	mapType := field.Type()
	keyType := mapType.Key()
	valueType := mapType.Elem()

	// Only string keys are supported
	if keyType.Kind() != reflect.String {
		return fmt.Errorf("field %s only string keys are supported for maps", fieldPath)
	}

	// Pre-allocate map capacity
	newMap := reflect.MakeMapWithSize(mapType, len(valueMap))

	for k, v := range valueMap {
		mapValue := reflect.New(valueType).Elem()
		if err := setFieldValue(mapValue, v, fmt.Sprintf("%s[%s]", fieldPath, k)); err != nil {
			return err
		}
		newMap.SetMapIndex(reflect.ValueOf(k), mapValue)
	}

	field.Set(newMap)
	return nil
}
