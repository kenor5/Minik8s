package scale

import "testing"

func TestStartPrometheusServer(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "Create PrometheusServer test",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := StartPrometheusServer(); (err != nil) != tt.wantErr {
				t.Errorf("StartPrometheusServer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
