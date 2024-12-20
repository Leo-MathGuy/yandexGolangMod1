package calcapi_test

import (
	"testing"

	"github.com/Leo-MathGuy/calcapi/pkg/calcapi"
)

func TestCalculator(t *testing.T) {
	testCases := []struct {
		name   string
		expr   string
		expect float64
	}{
		{
			name:   "sanity",
			expr:   "2+2",
			expect: 4,
		},
		{
			name:   "pemdas1",
			expr:   "1+4*(5/4+4*2/1)",
			expect: 38,
		},
		{
			name:   "pemdas2",
			expr:   "1+2*4",
			expect: 9,
		}, {
			name:   "negatives 1",
			expr:   "1*-1",
			expect: -1,
		},
		{
			name:   "negative 2",
			expr:   "1+-5",
			expect: -4,
		},
		{
			name:   "negative 3",
			expr:   "-5/2",
			expect: -2.5,
		},
		{
			name:   "spaces",
			expr:   " 2  + 2    + 0",
			expect: 4,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if val, err := calcapi.Calc(testCase.expr); err != nil {
				t.Fatalf("ERR: %s on test %s", err.Error(), testCase.name)
			} else if val != testCase.expect {
				t.Fatalf("%f should be equal %f", val, testCase.expect)
			}
		})
	}

	invalidTestCase := []struct {
		name string
		expr string
	}{
		{
			name: "end",
			expr: "2+5/",
		},
		{
			name: "middle",
			expr: "1++2+5",
		},
		{
			name: "parentheses",
			expr: "(1)+(2))",
		},
		{
			name: "empty",
			expr: "",
		},
		{
			name: "spaces",
			expr: "2 +  2 2",
		},
	}

	for _, testCase := range invalidTestCase {
		t.Run(testCase.name, func(t *testing.T) {
			if val, err := calcapi.Calc(testCase.expr); err == nil {
				t.Fatalf("ERR: %s passed with value %f", testCase.name, val)
			}
		})
	}
}
