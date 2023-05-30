package main

import "testing"

func Test_getDocekrIdByname(t *testing.T) {
	type args struct {
		podname string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Test Docker filter",
			args:    args{podname: "autoscale-deployment"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := getDocekrIdByname(tt.args.podname); (err != nil) != tt.wantErr {
				t.Errorf("getDocekrIdByname() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
