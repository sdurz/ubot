package ubot

import (
	"errors"
	"strings"
)

// V is a wrapper to any JSON value
type V struct {
	Value interface{}
}

// AsInteger unwraps the wrapped number value and return it as a integer
func (v *V) AsInteger() (result int64, err error) {
	var (
		fVal float64
		ok   bool
	)
	if fVal, ok = v.Value.(float64); !ok {
		err = errors.New("value is not a float64")
	} else {
		result = int64(fVal)
	}
	return
}

// AsFloat will unwrap the underlying number value and return it as a float
func (v *V) AsFloat() (result float64, err error) {
	var ok bool
	if result, ok = v.Value.(float64); !ok {
		err = errors.New("value is not a float64")
	}
	return
}

// AsString will unwrap the underlying string value
func (v *V) AsString() (result string, err error) {
	var ok bool
	if result, ok = v.Value.(string); !ok {
		err = errors.New("value is not a string")
	}
	return
}

// AsBool will unwrap the underlying boolean value
func (v *V) AsBool() (result bool, err error) {
	var ok bool
	if result, ok = v.Value.(bool); !ok {
		err = errors.New("value is not a bool")
	}
	return
}

// AsObject will unwrap the underlying objecy value and return it as O
func (v *V) AsObject() (result O, err error) {
	var (
		ok bool
	)
	if result, ok = v.Value.(map[string]interface{}); !ok {
		err = errors.New("value is not a O (map[string]interface{})")
	}
	return
}

// AsArray will unwrap the underlying objecy value and return it as A
func (v *V) AsArray() (result A, err error) {
	var ok bool
	if result, ok = v.Value.([]interface{}); !ok {
		err = errors.New("value is not a A ([]interface{})")
	}
	return
}

// O  represents a JSON Object
type O map[string]interface{}

// Get returns the path property as generic V
func (o *O) Get(path string) (result V, err error) {
	if path == "" {
		return result, errors.New("empty path")
	}
	splits := strings.Split(path, ".")
	nextProp := splits[0]
	remainingProps := splits[1:]

	var ok bool
	if len(remainingProps) > 0 {
		var nextVal interface{}
		if nextVal, ok = (*o)[nextProp]; !ok {
			err = errors.New("trying to traverse path " + path + ", unknown property: " + nextProp)
		} else {
			var nextO O
			if nextO, ok = nextVal.(map[string]interface{}); !ok {
				err = errors.New("trying to traverse path " + path + ", not a O (map[string]interface{}): " + nextProp)
			} else {
				result, err = nextO.Get(strings.Join(remainingProps, "."))
			}
		}
	} else {
		var value interface{}
		if value, ok = (*o)[nextProp]; !ok {
			err = errors.New("unknown key: " + nextProp)
		} else {
			result = V{value}
		}
	}
	return
}

// GetObject returns the path property as O
func (o *O) GetObject(path string) (result O, err error) {
	var value V
	if value, err = o.Get(path); err != nil {
		return
	}
	return value.AsObject()
}

// GetString returns the path property as string
func (o *O) GetString(path string) (result string, err error) {
	var value V
	if value, err = o.Get(path); err != nil {
		return
	}
	return value.AsString()
}

// GetBoolean returns the path property as bool
func (o *O) GetBoolean(path string) (result bool, err error) {
	var value V
	if value, err = o.Get(path); err != nil {
		return
	}
	return value.AsBool()
}

// GetInteger returns the path property as integer
func (o *O) GetInteger(path string) (result int64, err error) {
	var value V
	if value, err = o.Get(path); err != nil {
		return
	}
	return value.AsInteger()
}

// GetFloat returns the path property as float
func (o *O) GetFloat(path string) (result float64, err error) {
	var value V
	if value, err = o.Get(path); err != nil {
		return
	}
	return value.AsFloat()
}

// GetFloat returns the path property as A
func (o *O) GetArray(path string) (result A, err error) {
	var value V
	if value, err = o.Get(path); err != nil {
		return
	}
	return value.AsArray()
}

// A represents a JSON Array
type A []interface{}
