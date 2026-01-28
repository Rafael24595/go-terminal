package assert

import (
	"cmp"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func NotNil(t *testing.T, item any, message ...any) {
	t.Helper()

	custom := ""
	if len(message) > 0 {
		if format, ok := message[0].(string); ok {
			custom = fmt.Sprintf(format+" - ", message[1:]...)
		}
	}

	if item == nil {
		t.Errorf("%sUnexpected nil value", custom)
	}

	v := reflect.ValueOf(item)
	switch v.Kind() {
	case reflect.Func, reflect.Pointer, reflect.Map, reflect.Slice, reflect.Interface, reflect.Chan:
		if v.IsNil() {
			t.Errorf("%sUnexpected nil value", custom)
		}
	}
}

func Nil(t *testing.T, item any, message ...any) {
	t.Helper()

	custom := ""
	if len(message) > 0 {
		if format, ok := message[0].(string); ok {
			custom = fmt.Sprintf(format+" - ", message[1:]...)
		}
	}

	if item != nil {
		t.Errorf("%sExpected nil value", custom)
	}

	v := reflect.ValueOf(item)
	switch v.Kind() {
	case reflect.Func, reflect.Pointer, reflect.Map, reflect.Slice, reflect.Interface, reflect.Chan:
		if !v.IsNil() {
			t.Errorf("%sUnexpected nil value", custom)
		}
	}
}

func True(t *testing.T, result bool, message ...any) {
	t.Helper()

	if result {
		return
	}

	custom := ""
	if len(message) > 0 {
		if format, ok := message[0].(string); ok {
			custom = fmt.Sprintf(format+" - ", message[1:]...)
		}
	}

	t.Errorf("%sExpected true, but got false", custom)
}

func False(t *testing.T, result bool, message ...any) {
	t.Helper()

	if !result {
		return
	}

	custom := ""
	if len(message) > 0 {
		if format, ok := message[0].(string); ok {
			custom = fmt.Sprintf(format+" - ", message[1:]...)
		}
	}

	t.Errorf("%sExpected false, but got true", custom)
}

func Equal[T comparable](t *testing.T, want, have T, message ...any) {
	t.Helper()

	if want == have {
		return
	}

	custom := ""
	if len(message) > 0 {
		if format, ok := message[0].(string); ok {
			custom = fmt.Sprintf(format+" - ", message[1:]...)
		}
	}

	t.Errorf("%sExpected '%v', but got '%v'", custom, want, have)
}

func GreaterOrEqual[T cmp.Ordered](t *testing.T, want, have T, message ...any) {
	t.Helper()

	if have >= want {
		return
	}

	custom := ""
	if len(message) > 0 {
		if format, ok := message[0].(string); ok {
			custom = fmt.Sprintf(format+" - ", message[1:]...)
		}
	}

	t.Errorf("%sExpected greater or equal than %v, but got %v", custom, want, have)
}

func LessOrEqual[T cmp.Ordered](t *testing.T, want, have T, message ...any) {
	t.Helper()

	if have <= want {
		return
	}

	custom := ""
	if len(message) > 0 {
		if format, ok := message[0].(string); ok {
			custom = fmt.Sprintf(format+" - ", message[1:]...)
		}
	}

	t.Errorf("%sExpected less or equal than %v, but got %v", custom, want, have)
}

func Error(t *testing.T, err error, message ...any) {
	t.Helper()

	if err != nil {
		return
	}

	custom := ""
	if len(message) > 0 {
		if format, ok := message[0].(string); ok {
			custom = fmt.Sprintf(format+" - ", message[1:]...)
		}
	}

	t.Errorf("%sExpected error found but nothing found", custom)
}

func NotError(t *testing.T, err error, message ...any) {
	t.Helper()

	if err == nil {
		return
	}

	custom := ""
	if len(message) > 0 {
		if format, ok := message[0].(string); ok {
			custom = fmt.Sprintf(format+" - ", message[1:]...)
		}
	}

	t.Errorf("%sUnexpected error found: '%s'", custom, err.Error())
}

func Len[T any](t *testing.T, want int, have []T, message ...any) {
	t.Helper()

	if want == len(have) {
		return
	}

	custom := ""
	if len(message) > 0 {
		if format, ok := message[0].(string); ok {
			custom = fmt.Sprintf(format+" - ", message[1:]...)
		}
	}

	t.Fatalf("%sExpected %v, but got %v", custom, want, len(have))
}

func Contains[T comparable](t *testing.T, container any, item T, message ...any) {
	t.Helper()

	custom := ""
	if len(message) > 0 {
		if format, ok := message[0].(string); ok {
			custom = fmt.Sprintf(format+" - ", message[1:]...)
		}
	}

	val := reflect.ValueOf(container)

	switch val.Kind() {
	case reflect.String:
		substr, ok := any(item).(string)
		if !ok {
			t.Fatalf("%sCannot search non-string in string container", custom)
		}

		if !strings.Contains(val.String(), substr) {
			t.Errorf("%sExpected '%s' to contain '%s'", custom, val.String(), substr)
		}

		return
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			elem := val.Index(i).Interface()
			if elem == item {
				return
			}
		}

		t.Errorf("%sExpected slice/array to contain '%v'", custom, item)

		return
	}

	t.Fatalf("%sContains not supported for type %s", custom, val.Kind())
}

func NotContains[T comparable](t *testing.T, container any, item T, message ...any) {
	t.Helper()

	custom := ""
	if len(message) > 0 {
		if format, ok := message[0].(string); ok {
			custom = fmt.Sprintf(format+" - ", message[1:]...)
		}
	}

	val := reflect.ValueOf(container)
	switch val.Kind() {
	case reflect.String:
		substr, ok := any(item).(string)
		if !ok {
			t.Fatalf("%sCannot search non-string in string container", custom)
		}

		if strings.Contains(val.String(), substr) {
			t.Errorf("%sExpected '%s' NOT to contain '%s'", custom, val.String(), substr)
		}

		return
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			elem := val.Index(i).Interface()
			if elem == item {
				t.Errorf("%sExpected slice/array NOT to contain '%v'", custom, item)
				return
			}
		}

		return
	}

	t.Fatalf("%sNotContains not supported for type %s", custom, val.Kind())
}
