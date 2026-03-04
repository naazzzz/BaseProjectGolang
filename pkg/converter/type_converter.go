package converter

import (
	"encoding/json"
	"reflect"

	"github.com/rotisserie/eris"
)

func TypeConverter[outType any](inType any) (out outType, err error) {
	marshalObject, err := json.Marshal(&inType)
	if err != nil {
		err = eris.Wrap(err, "failed to marshal converter")
		return
	}

	if err = json.Unmarshal(marshalObject, &out); err != nil {
		err = eris.Wrap(err, "failed to unmarshal converter")
		return
	}

	return
}

func GetType(v any) string {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	}

	return t.Name()
}

func Val[T any, P *T](p P) T {
	if p != nil {
		return *p
	}

	var def T

	return def
}

func Ptr[T comparable](t T) *T {
	var def T
	if t == def {
		return nil
	}

	return &t
}
