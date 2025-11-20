package utils

import "testing"

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple camelCase",
			input:    "camelCase",
			expected: "camel_case",
		},
		{
			name:     "PascalCase",
			input:    "PascalCase",
			expected: "pascal_case",
		},
		{
			name:     "single word lowercase",
			input:    "lowercase",
			expected: "lowercase",
		},
		{
			name:     "single word uppercase",
			input:    "UPPERCASE",
			expected: "uppercase",
		},
		{
			name:     "multiple words camelCase",
			input:    "thisIsALongVariableName",
			expected: "this_is_a_long_variable_name",
		},
		{
			name:     "acronym at end",
			input:    "orderID",
			expected: "order_id",
		},
		{
			name:     "acronym at beginning",
			input:    "IDOrder",
			expected: "id_order",
		},
		{
			name:     "acronym in middle",
			input:    "orderIDNumber",
			expected: "order_id_number",
		},
		{
			name:     "consecutive capitals",
			input:    "HTTPSConnection",
			expected: "https_connection",
		},
		{
			name:     "all caps acronym",
			input:    "URL",
			expected: "url",
		},
		{
			name:     "mixed case with numbers",
			input:    "variable123Name",
			expected: "variable123_name",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "single character lowercase",
			input:    "a",
			expected: "a",
		},
		{
			name:     "single character uppercase",
			input:    "A",
			expected: "a",
		},
		{
			name:     "already snake_case",
			input:    "already_snake_case",
			expected: "already_snake_case",
		},
		{
			name:     "with underscores",
			input:    "some_Mixed_Case",
			expected: "some__mixed__case",
		},
		{
			name:     "DroneID",
			input:    "DroneID",
			expected: "drone_id",
		},
		{
			name:     "UserID",
			input:    "UserID",
			expected: "user_id",
		},
		{
			name:     "OrderNumber",
			input:    "OrderNumber",
			expected: "order_number",
		},
		{
			name:     "DeliveredByDroneID",
			input:    "DeliveredByDroneID",
			expected: "delivered_by_drone_id",
		},
		{
			name:     "UpdatedByID",
			input:    "UpdatedByID",
			expected: "updated_by_id",
		},
		{
			name:     "APIKey",
			input:    "APIKey",
			expected: "api_key",
		},
		{
			name:     "XMLParser",
			input:    "XMLParser",
			expected: "xml_parser",
		},
		{
			name:     "IOStream",
			input:    "IOStream",
			expected: "io_stream",
		},
		{
			name:     "getHTTPResponse",
			input:    "getHTTPResponse",
			expected: "get_http_response",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToSnakeCase(tt.input)
			if result != tt.expected {
				t.Errorf("ToSnakeCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func BenchmarkToSnakeCase(b *testing.B) {
	inputs := []string{
		"camelCase",
		"PascalCase",
		"thisIsALongVariableName",
		"HTTPSConnection",
		"DeliveredByDroneID",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, input := range inputs {
			ToSnakeCase(input)
		}
	}
}
