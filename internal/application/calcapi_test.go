package application_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Leo-MathGuy/calcapi/internal/application"
	"github.com/gin-gonic/gin"
)

func TestCalcapiHandler(t *testing.T) {
	returns422 := []struct {
		name string
		expr *strings.Reader
	}{
		{
			name: "end",
			expr: strings.NewReader(`{"expression":"2+5/"}`),
		},
		{
			name: "middle",
			expr: strings.NewReader(`{"expression":"1++2+5"}`),
		},
		{
			name: "parentheses",
			expr: strings.NewReader(`{"expression":"(1)+(2))"}`),
		},
		{
			name: "empty",
			expr: strings.NewReader(`{"expression":""}`),
		},
		{
			name: "spaces",
			expr: strings.NewReader(`{"expression":"2 +  2 2"}`),
		},
	}

	for _, testCase := range returns422 {
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		req := httptest.NewRequest(http.MethodPost, "/", testCase.expr)
		c.Request = req

		application.CalculateHandler(c)
		res := recorder.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Errorf("ERR: 422 not given for %s", testCase.name)
		}
	}

	testCases := []struct {
		name   string
		expr   *strings.Reader
		expect float64
	}{
		{
			name:   "sanity",
			expr:   strings.NewReader(`{"expression":"2+2"}`),
			expect: 4,
		},
		{
			name:   "pemdas1",
			expr:   strings.NewReader(`{"expression":"1+4*(5/4+4*2/1)"}`),
			expect: 38,
		},
		{
			name:   "pemdas2",
			expr:   strings.NewReader(`{"expression":"1+2*4"}`),
			expect: 9,
		}, {
			name:   "negatives 1",
			expr:   strings.NewReader(`{"expression":"1*-1"}`),
			expect: -1,
		},
		{
			name:   "negative 2",
			expr:   strings.NewReader(`{"expression":"1+-5"}`),
			expect: -4,
		},
		{
			name:   "negative 3",
			expr:   strings.NewReader(`{"expression":"-5/2"}`),
			expect: -2.5,
		},
		{
			name:   "spaces",
			expr:   strings.NewReader(`{"expression":" 2  + 2    + 0"}`),
			expect: 4,
		},
	}

	for _, testCase := range testCases {
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		req := httptest.NewRequest(http.MethodPost, "/", testCase.expr)
		c.Request = req

		application.CalculateHandler(c)
		res := recorder.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("ERR: 200 not given for %s", testCase.name)
		}

		var resJson application.ResultResponse

		buf := make([]byte, 1024)
		if readLen, err := res.Body.Read(buf); err != nil {
			t.Errorf("ERR: %s", err.Error())
		} else {
			buf = buf[:readLen]
		}

		if err := json.Unmarshal(buf, &resJson); err != nil {
			t.Errorf("ERR: %s", err.Error())
		}

		if resJson.Result != testCase.expect {
			t.Errorf("ERR: %f != returned %f", testCase.expect, resJson.Result)
		}
	}
}
