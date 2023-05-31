package resourse

import "testing"

func TestStartcAdvisor(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Valid test should not return error",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := StartcAdvisor(); (err != nil) != tt.wantErr {
				t.Errorf("StartcAdvisor() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
