package formatter

import (
	"testing"
)

func TestIntToString(t *testing.T) {
	tests := []struct {
		name     string
		input    int64
		expected string
	}{
		{
			name:     "positive integer",
			input:    12345,
			expected: "12345",
		},
		{
			name:     "negative integer",
			input:    -6789,
			expected: "-6789",
		},
		{
			name:     "zero",
			input:    0,
			expected: "0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IntToString(tt.input)
			if got != tt.expected {
				t.Errorf("IntToString() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFloatToString(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected string
	}{
		{
			name:     "positive float",
			input:    123.456,
			expected: "123.456",
		},
		{
			name:     "negative float",
			input:    -78.9,
			expected: "-78.9",
		},
		{
			name:     "zero",
			input:    0.0,
			expected: "0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FloatToString(tt.input)
			if got != tt.expected {
				t.Errorf("FloatToString() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestStringToInt(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
		wantErr  bool
	}{
		{
			name:     "valid positive integer",
			input:    "12345",
			expected: 12345,
			wantErr:  false,
		},
		{
			name:     "valid negative integer",
			input:    "-6789",
			expected: -6789,
			wantErr:  false,
		},
		{
			name:     "invalid integer",
			input:    "abc",
			expected: 0,
			wantErr:  true,
		},
		{
			name:     "empty string",
			input:    "",
			expected: 0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StringToInt(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringToInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("StringToInt() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestStringToFloat(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
		wantErr  bool
	}{
		{
			name:     "valid positive float",
			input:    "123.456",
			expected: 123.456,
			wantErr:  false,
		},
		{
			name:     "valid negative float",
			input:    "-78.9",
			expected: -78.9,
			wantErr:  false,
		},
		{
			name:     "invalid float",
			input:    "abc",
			expected: 0,
			wantErr:  true,
		},
		{
			name:     "empty string",
			input:    "",
			expected: 0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StringToFloat(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringToFloat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("StringToFloat() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCheckSchemeFormat(t *testing.T) {
	tests := []struct {
		name    string
		scheme  string
		wantErr bool
	}{
		{
			name:    "valid http scheme",
			scheme:  "http",
			wantErr: false,
		},
		{
			name:    "valid https scheme",
			scheme:  "https",
			wantErr: false,
		},
		{
			name:    "invalid scheme",
			scheme:  "ftp",
			wantErr: true,
		},
		{
			name:    "empty scheme",
			scheme:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckSchemeFormat(tt.scheme)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckSchemeFormat() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCheckAddrFormat(t *testing.T) {
	tests := []struct {
		name    string
		addr    string
		wantErr bool
	}{
		{
			name:    "valid address",
			addr:    "localhost:8080",
			wantErr: false,
		},
		{
			name:    "invalid address (missing port)",
			addr:    "localhost",
			wantErr: true,
		},
		{
			name:    "invalid address (non-numeric port)",
			addr:    "localhost:abc",
			wantErr: true,
		},
		{
			name:    "empty address",
			addr:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckAddrFormat(tt.addr)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckAddrFormat() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
