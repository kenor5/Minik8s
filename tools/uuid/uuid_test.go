package uuid

import (
	"fmt"
	"testing"
)

func TestUUID(t *testing.T) {
	fmt.Println(UUID())
	fmt.Println(UUID())
	fmt.Println(UUID())
}

func TestUUID1(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "UUID",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UUID(); got != tt.want {
				t.Errorf("UUID() = %v, want %v", got, tt.want)
			}
		})
	}
}
