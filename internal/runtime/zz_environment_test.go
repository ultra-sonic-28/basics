package runtime

import (
	"testing"

	"basics/testutils"
)

func TestEnv_SetGet(t *testing.T) {
	tests := []struct {
		name       string
		setName    string
		setValue   Value
		getName    string
		wantValue  Value
		wantExists bool
	}{
		{
			name:       "Get after set number",
			setName:    "X",
			setValue:   Value{Type: NUMBER, Num: 42},
			getName:    "X",
			wantValue:  Value{Type: NUMBER, Num: 42},
			wantExists: true,
		},
		{
			name:       "Get after set string",
			setName:    "MSG",
			setValue:   Value{Type: STRING, Str: "Hello"},
			getName:    "MSG",
			wantValue:  Value{Type: STRING, Str: "Hello"},
			wantExists: true,
		},
		{
			name:       "Get after set boolean to true",
			setName:    "MSG",
			setValue:   Value{Type: BOOLEAN, Flag: true},
			getName:    "MSG",
			wantValue:  Value{Type: BOOLEAN, Flag: true},
			wantExists: true,
		},
		{
			name:       "Get after set boolean to false",
			setName:    "MSG",
			setValue:   Value{Type: BOOLEAN, Flag: false},
			getName:    "MSG",
			wantValue:  Value{Type: BOOLEAN, Flag: false},
			wantExists: true,
		},
		{
			name:       "Get unset variable returns default number 0",
			setName:    "",
			setValue:   Value{},
			getName:    "UNDEF",
			wantValue:  Value{Type: NUMBER, Num: 0},
			wantExists: false,
		},
		{
			name:       "Overwrite variable",
			setName:    "A",
			setValue:   Value{Type: NUMBER, Num: 5},
			getName:    "A",
			wantValue:  Value{Type: NUMBER, Num: 10},
			wantExists: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewEnvironment()

			// Si setName est non vide, on définit la variable
			if tt.setName != "" {
				env.Set(tt.setName, tt.setValue)
			}

			// Cas "Overwrite variable"
			if tt.name == "Overwrite variable" {
				env.Set("A", Value{Type: NUMBER, Num: 10})
			}

			got, exists := env.Get(tt.getName)

			// Vérification des valeurs
			testutils.Equal(t, "exists check", exists, tt.wantExists)
			testutils.Equal(t, "Type check", got.Type, tt.wantValue.Type)
			testutils.Equal(t, "Num check", got.Num, tt.wantValue.Num)
			testutils.Equal(t, "Str check", got.Str, tt.wantValue.Str)
		})
	}
}

func TestValue_String(t *testing.T) {
	tests := []struct {
		name     string
		value    Value
		expected string
	}{
		{
			name: "STRING value",
			value: Value{
				Type: STRING,
				Str:  "HELLO",
			},
			expected: "HELLO",
		},
		{
			name: "INTEGER positive",
			value: Value{
				Type: INTEGER,
				Int:  42,
			},
			expected: "42",
		},
		{
			name: "INTEGER negative",
			value: Value{
				Type: INTEGER,
				Int:  -7,
			},
			expected: "-7",
		},
		{
			name: "NUMBER integer-like",
			value: Value{
				Type: NUMBER,
				Num:  5,
			},
			expected: "5.000000",
		},
		{
			name: "NUMBER float",
			value: Value{
				Type: NUMBER,
				Num:  3.14159,
			},
			expected: "3.141590",
		},
		{
			name: "NUMBER negative float",
			value: Value{
				Type: NUMBER,
				Num:  -1.75,
			},
			expected: "-1.750000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.value.String()
			testutils.Equal(t, "string representation", tt.expected, result)
		})
	}
}
