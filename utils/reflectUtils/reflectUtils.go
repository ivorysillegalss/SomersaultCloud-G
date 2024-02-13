package reflectUtils

import "reflect"

// 判断一个变量的值是否零值
func IsZero(v any) bool {
	return reflect.ValueOf(v).IsZero()
}

// UpdateStruct 使用反射更新结构体字段，包括嵌套的结构体和指针
func UpdateStruct(v reflect.Value, updates map[string]interface{}) {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// 遍历所有字段，尝试进行更新
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		typeField := v.Type().Field(i)

		// 获取json标签作为键值
		jsonTag := typeField.Tag.Get("json")
		if jsonTag == "" {
			jsonTag = typeField.Name
		}

		if updateValue, ok := updates[jsonTag]; ok {
			// 如果是结构体或结构体指针，需要递归处理
			if field.Kind() == reflect.Struct {
				UpdateStruct(field.Addr(), updateValue.(map[string]interface{}))
			} else if field.Kind() == reflect.Ptr && field.Type().Elem().Kind() == reflect.Struct {
				if field.IsNil() {
					field.Set(reflect.New(field.Type().Elem()))
				}
				UpdateStruct(field, updateValue.(map[string]interface{}))
			} else {
				// 对于基础类型，判断是否为空，部分更新
				newValue := reflect.ValueOf(updateValue)
				//根据类型判断 默认除布尔类型以外 数值不能为零值 不然会导致逻辑混乱 布尔型就不管他了
				if newValue.IsZero() && newValue.Kind() != reflect.Bool {
					continue
					//直接continue到下一轮更新
				}
				if newValue.Type().AssignableTo(field.Type()) {
					field.Set(newValue)
				}
			}
		}
	}
}
