package logger

import (
	"reflect"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestGetLogger(t *testing.T) {
	type args struct {
		level string
	}
	tests := []struct {
		name      string
		args      args
		want      *zap.Logger
		wantLevel zapcore.Level
		wantErr   bool
	}{
		{
			name: "Get info logger",
			args: args{
				level: "INFO",
			},
			wantLevel: zapcore.InfoLevel,
			wantErr:   false,
		},
		{
			name: "Get warn logger",
			args: args{
				level: "WARN",
			},
			wantLevel: zapcore.WarnLevel,
			wantErr:   false,
		},
		{
			name: "Get error logger (invalid level)",
			args: args{
				level: "INVALID_LEVEL",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetLogger(tt.args.level)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLogger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("GetLogger() got = %v, want %v", got, tt.want)
				return
			}
			if !tt.wantErr && got.Level() != tt.wantLevel {
				t.Errorf("GetLogger() got level = %v, want %v", got.Level(), tt.wantLevel)
				return
			}
		})
	}
}
