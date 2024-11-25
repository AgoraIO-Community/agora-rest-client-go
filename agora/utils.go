package agora

import "reflect"

func Contains(S []string, E string) bool {
	for _, s := range S {
		if s == E {
			return true
		}
	}
	return false
}

func IsNil(v interface{}) bool {
	if v == nil {
		return true
	}

	rv := reflect.ValueOf(v)
	return rv.Kind() == reflect.Ptr && rv.IsNil()
}
