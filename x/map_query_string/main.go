package main

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
)

func structToValues(ptr any) url.Values {
	values := url.Values{}
	v := reflect.Indirect(reflect.ValueOf(ptr))

	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// Convert the field value to a string
		var strValue string
		switch value.Kind() {
		case reflect.String:
			strValue = url.QueryEscape(value.String())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			strValue = strconv.FormatInt(value.Int(), 10)
		case reflect.Float32, reflect.Float64:
			strValue = strconv.FormatFloat(value.Float(), 'f', -1, 64)
		case reflect.Bool:
			strValue = strconv.FormatBool(value.Bool())
		default:
			continue
		}

		// Add the key-value pair to the values
		if strValue != "" {
			key := field.Tag.Get("json")
			if key == "" {
				key = field.Name
			}
			values.Add(key, strValue)
		}
	}
	return values
}

func valuesToStruct(values url.Values, ptr any) error {
	dType := reflect.TypeOf(ptr)
	dhVal := reflect.ValueOf(ptr)

	for i := 0; i < dType.Elem().NumField(); i++ {
		field := dType.Elem().Field(i)
		key := field.Tag.Get("json")
		kind := field.Type.Kind()

		val := values.Get(key)

		result := dhVal.Elem().Field(i)

		switch kind {
		case reflect.String:
			val, err := url.QueryUnescape(val)
			if err != nil {
				return err
			}
			result.SetString(val)
		case reflect.Int:
			v, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return err
			}
			result.SetInt(v)
		case reflect.Float64:
			v, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return err
			}
			result.SetFloat(v)
		case reflect.Bool:
			v, err := strconv.ParseBool(val)
			if err != nil {
				return err
			}
			result.SetBool(v)
		default:
			return fmt.Errorf("unsupported type %s", kind)
		}
	}
	return nil
}

func structToURL(host string, ptr any) (string, error) {
	u, err := url.Parse(host)
	if err != nil {
		return "", err
	}

	u.RawQuery = structToValues(ptr).Encode()

	return u.String(), nil
}

func main() {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	u, err := url.Parse("http://example.com")
	if err != nil {
		panic(err)
	}

	p := Person{
		Name: "John",
		Age:  30,
	}

	values := structToValues(&p)

	u.RawQuery = values.Encode()

	fmt.Println(u.String())

	fmt.Println(values)

	var p2 Person
	valuesToStruct(values, &p2)

	fmt.Println(p2)
}
