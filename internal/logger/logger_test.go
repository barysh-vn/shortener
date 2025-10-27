package logger

import (
	"testing"

	"go.uber.org/zap"
)

func TestInitialize(t *testing.T) {
	tests := []struct {
		name      string
		level     string
		wantErr   bool
		checkFunc func(t *testing.T)
	}{
		{
			name:    "Test valid level - info",
			level:   "info",
			wantErr: false,
			checkFunc: func(t *testing.T) {
				if BaseLogger == zap.NewNop() {
					t.Error("expected BaseLogger to be initialized, got NOP")
				}
				BaseLogger.Info("info log test")
			},
		},
		{
			name:    "Test valid level - debug (should enable debug logs)",
			level:   "debug",
			wantErr: false,
			checkFunc: func(t *testing.T) {
				if !BaseLogger.Core().Enabled(zap.DebugLevel) {
					t.Error("expected debug level to be enabled")
				}
			},
		},
		{
			name:    "Test invalid level",
			level:   "notalevel",
			wantErr: true,
		},
		{
			name:    "Test reinitialize",
			level:   "warn",
			wantErr: false,
			checkFunc: func(t *testing.T) {
				err := Initialize("error")
				if err != nil {
					t.Errorf("unexpected error during reinitialization: %v", err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Initialize(tt.level)

			if (err != nil) != tt.wantErr {
				t.Errorf("Initialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.checkFunc != nil {
				tt.checkFunc(t)
			}
		})
	}
}

func TestBaseSugarLogger(t *testing.T) {
	if err := Initialize("info"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	sugar := BaseSugarLogger()
	if sugar == nil {
		t.Fatal("expected non-nil sugared logger")
	}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("BaseSugarLogger panicked: %v", r)
		}
	}()

	sugar.Infow("sugar test", "key", "value")
}
