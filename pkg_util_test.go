package wo

import (
	"reflect"
	"runtime"
	"testing"
)

func funcName(t *testing.T, f interface{}) string {
	t.Helper()
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}
