package pkg

import "testing"

func TestContextKey_String(t *testing.T) {
	tests := []struct {
		name string
		c    ContextKey
		want string
	}{
		{
			name: "adds prefix to context key when converting to string",
			c:    "sample-key",
			want: "gophermart-sample-key",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
