package random_test

import (
	"reflect"
	"testing"

	"github.com/mikhail-alaska/cli-messenger/server/internal/lib/random"
)

func TestNewRandomString(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		length int
		want   string
	}{
        {
            name: "first",
            length: 10,
            want: "firstfirst",
        },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := random.NewRandomString(tt.length)

            gt := reflect.TypeOf(got).Kind()
            if gt != reflect.String{
				t.Errorf("NewRandomString() = %v, want %v", got, tt.want)
            }
		})
	}
}

