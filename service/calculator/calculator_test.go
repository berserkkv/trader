package calculator

import "testing"

func TestCalculateStopLossWithPercent(t *testing.T) {
	tests := []struct {
		name     string
		price    float64
		percent  float64
		isShort  bool
		expected float64
	}{
		{"Long position - 5%", 100.0, 5.0, false, 95.0},
		{"Short position - 5%", 100.0, 5.0, true, 105.0},
		{"Long position - 10%", 200.0, 10.0, false, 180.0},
		{"Short position - 10%", 200.0, 10.0, true, 220.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateStopLossWithPercent(tt.price, tt.percent, tt.isShort)
			if result != tt.expected {
				t.Errorf("got %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCalculateTakeProfitWithPercent(t *testing.T) {
	tests := []struct {
		name     string
		price    float64
		percent  float64
		isShort  bool
		expected float64
	}{
		{"Long position - 5%", 100.0, 5.0, false, 105.0},
		{"Short position - 5%", 100.0, 5.0, true, 95.0},
		{"Long position - 10%", 200.0, 10.0, false, 220.0},
		{"Short position - 10%", 200.0, 10.0, true, 180.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateTakeProfitWithPercent(tt.price, tt.percent, tt.isShort)
			if result != tt.expected {
				t.Errorf("got %v, want %v", result, tt.expected)
			}
		})
	}
}
