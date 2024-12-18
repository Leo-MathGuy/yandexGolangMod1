/*
Calculator Web API
Copyright (C) 2024 @Leo-MathGuy

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"fmt"

	"errors"
	"strconv"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
)

//// MARK: Structs

type CalculateRequest struct {
	Expr string `json:"expression"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type ResultResponse struct {
	Error string `json:"result"`
}

//// MARK: Server

func ErrorHandler(c *gin.Context, err error, status int) {
	var errr string = "Internal Server Error"
	if err != nil {
		errr = err.Error()
	}
	c.JSON(status, ErrorResponse{errr})
}

func CalculateHandler(c *gin.Context) {
	var requestBody CalculateRequest
	if err := c.BindJSON(&requestBody); err != nil {
		ErrorHandler(c, err, 500)
		return
	}

	if result, err := Calc(requestBody.Expr); err != nil {
		ErrorHandler(c, err, 422)
		return
	} else {
		c.JSON(200, ResultResponse{fmt.Sprint(result)})
	}
}

func RunServer() {
	r := gin.Default()
	r.POST("/api/v1/calculate", CalculateHandler)

	r.Run(":80")
}

func main() {
	RunServer()
}

//// MARK: Calculator

// Constants for lookup
const NUMBER, OP, SPACE, PAR int = 0, 1, 2, 3
const OPS string = "+-*/"
const PARS string = "()"

// // MARK: Helper functions

// Lookup
func get_op(x float64) string {
	return string(OPS[int(x)])
}
func get_par(x float64) string {
	return string(PARS[int(x)])
}

// Check if a float number equals a set int (for readability)
func is_op(x float64) bool {
	return int(x) == OP
}
func is_par(x float64) bool {
	return int(x) == PAR
}

// Symbol type (really needed doc I know)
func SymbolType(x rune) int {
	if unicode.IsDigit(x) || x == '.' {
		return NUMBER
	}

	switch x {
	case ' ':
		return SPACE
	case '+', '-', '/', '*':
		return OP
	case '(', ')':
		return PAR
	}

	return -1
}

// // MARK: The stuff™
func Calc(expr string) (float64, error) {
	if len(expr) == 0 {
		return 0, errors.New("empty expression")
	}

	/// Separate symbols
	rexpr := make([]string, 0)
	expr = strings.TrimSpace(expr)

	for _, v := range expr {
		t := SymbolType(v)

		if t == -1 {
			return 0, errors.New("invalid character " + string(v))
		}

		rexpr = append(rexpr, string(v))
	}

	// DEBUG: for _, v := range rexpr {
	// DEBUG: 	fmt.Printf("Type %d - %s\n", SymbolType(rune(v[0])), v)
	// DEBUG: }

	/// Validate operators at ends of equations
	l := len(rexpr)

	if rexpr[0] != "-" && SymbolType(rune(rexpr[0][0])) == OP {
		return 0, errors.New("Operator at start: " + rexpr[0])
	}

	if SymbolType(rune(rexpr[l-1][0])) == OP {
		return 0, errors.New("Operator at end: " + rexpr[l-1])
	}

	/// Check spaces (because spaces are guaranteed not to be at the ends now)
	cleared := make([]string, 0)

	for i := 0; i < len(rexpr)-1; i++ {
		if rexpr[i] == " " && rexpr[i+1] == " " {
			// DEBUG: fmt.Printf("Removed space at %d\n", i)
			rexpr = append(rexpr[:i], rexpr[i+1:]...)
			i--
		}
	}

	for i, v := range rexpr {
		if v == " " {

			before := SymbolType(rune(rexpr[i-1][0]))
			after := SymbolType(rune(rexpr[i+1][0]))

			if before == after && (before != PAR) && rexpr[i+1] != "-" {
				return 0, errors.New("Invalid space: \"" + rexpr[i-1] + " " + rexpr[i+1] + "\"")
			}
		} else {
			cleared = append(cleared, v)
		}
	}
	rexpr = cleared

	/// Check parentheses
	count := 0

	for _, v := range rexpr {
		if v == "(" {
			count++
		}

		if v == ")" {
			count--

			if count < 0 {
				return 0, errors.New("invalid parentheses")
			}
		}
	}

	if count != 0 {
		return 0, errors.New("invalid parentheses")
	}

	// DEBUG: fmt.Println()

	/// Join numbers
	l = len(rexpr)

	cleared = make([]string, 0)

	skip := 0

	for i, v := range rexpr {
		if skip > 0 {
			skip--
			continue
		}

		minus_is_num := i == 0 || SymbolType(rune(rexpr[i-1][0])) == OP || rune(rexpr[i-1][0]) == '('

		if (SymbolType(rune(v[0])) != 0 && !(minus_is_num && v == "-")) || i == l-1 {
			cleared = append(cleared, v)
			// DEBUG: fmt.Printf("Added: %s\n", v)
			continue
		}

		j := 1
		cleared = append(cleared, v)
		// DEBUG: fmt.Printf("Adding: %s", v)

		for SymbolType(rune(rexpr[i+j][0])) == 0 {
			cleared[len(cleared)-1] += rexpr[i+j]
			// DEBUG: fmt.Printf(", %s", rexpr[i+j])
			skip++
			j++

			if i+j >= l {
				break
			}
		}
		// DEBUG: fmt.Println()
	}

	rexpr = cleared

	/// Number dot check
	for index, v := range rexpr {
		if v == "." {
			return 0, errors.New("invalid number: \".\"")
		}

		if strings.Count(v, ".") > 1 {
			return 0, errors.New("invalid number: " + v)
		}

		i := strings.Index(v, ".")
		if i == -1 {
			continue
		}

		if i == len(v)-1 {
			return 0, errors.New("invalid number: " + v)
		}

		if i == 0 {
			rexpr[index] = "0" + v
		}
	}

	// Turn into numbers
	fexpr := make([][]float64, 0)

	for _, v := range rexpr {
		if len(v) > 1 || SymbolType(rune(v[0])) == NUMBER {
			num, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return 0, errors.New("Error during number parse: " + v + " - " + err.Error())
			}
			fexpr = append(fexpr, []float64{float64(NUMBER), num})
		} else if SymbolType(rune(v[0])) == OP {
			fexpr = append(fexpr, []float64{float64(OP), float64(strings.Index(OPS, v))})
		} else if SymbolType(rune(v[0])) == PAR {
			fexpr = append(fexpr, []float64{float64(PAR), float64(strings.Index(PARS, v))})
		}
	}

	// DEBUG: fmt.Println("\nFinished final transform:")

	// DEBUG: for i, v := range fexpr {
	// DEBUG: fmt.Println(strings.TrimRight(fmt.Sprintf(" - %d - Type %d - value %f", i, int(v[0]), v[1]), "0"))
	// DEBUG: }

	return Eval(fexpr)
}

func Eval(expr [][]float64) (float64, error) {
	priority := false

	for i, v := range expr {
		kind := v[0]
		value := v[1]

		if is_par(kind) && get_par(value) == "(" {
			count := 1
			x := 0

			for x = i; !(is_par(expr[x][0]) && get_par(expr[x][1]) == ")" && count == 0); x++ {
				next := expr[x+1]
				if is_par(next[0]) {
					if get_par(next[1]) == "(" {
						count++
					} else {
						count--
					}
				}
			}

			// DEBUG: fmt.Printf("Found parentheses - %d and %d\n", i, x)
			inbetween := expr[i+1 : x]

			// DEBUG: fmt.Println(" - Inside:")
			// DEBUG: for i, v := range inbetween {
			// DEBUG: fmt.Println(strings.TrimRight(fmt.Sprintf(" - - %d - Type %d - value %f", i, int(v[0]), v[1]), "0"))
			// DEBUG: }

			result, err := Eval(inbetween)

			if err != nil {
				return 0, err
			}

			expr = append(append(expr[:i], []float64{float64(NUMBER), result}), expr[x+1:]...)
			// DEBUG: fmt.Println(" - Result:")
			// DEBUG: for i, v := range expr {
			// DEBUG: fmt.Println(strings.TrimRight(fmt.Sprintf(" - - %d - Type %d - value %f", i, int(v[0]), v[1]), "0"))
			// DEBUG: }

			return Eval(expr)
		}

		if is_op(kind) && (get_op(value) == "*" || get_op(value) == "/") {
			priority = !priority
		}
	}

	processed := make([][]float64, 0)
	skip := 0
	for i := 0; i < len(expr); i++ {
		if skip > 0 {
			skip--
			continue
		}

		cur := expr[i]
		kind := cur[0]
		value := cur[1]

		if kind == float64(NUMBER) ||
			(is_op(kind) && (get_op(value) == "+" || get_op(value) == "-")) {
			processed = append(processed, cur)
			continue
		}

		if !is_op(kind) {
			return 0, errors.New("we ding dong messed up real time")
		}

		value1 := processed[len(processed)-1][1]
		value2 := expr[i+1][1]

		result := 0.0
		if get_op(value) == "*" {
			result = value1 * value2
			// DEBUG: fmt.Printf("Simplified %f * %f\n", value1, value2)
		} else {
			if value2 == 0 {
				return 0, errors.New("division by zero")
			}
			result = value1 / value2
			// DEBUG: fmt.Printf("Simplified %f / %f\n", value1, value2)
		}

		processed = processed[:len(processed)-1]
		processed = append(processed, []float64{float64(NUMBER), result})
		skip++
	}
	expr = processed

	processed = make([][]float64, 0)
	skip = 0
	for i := 0; i < len(expr); i++ {
		if skip > 0 {
			skip--
			continue
		}

		cur := expr[i]
		kind := cur[0]
		value := cur[1]

		if kind == float64(NUMBER) {
			processed = append(processed, cur)
			continue
		}

		if !is_op(kind) {
			return 0, errors.New("we ding dong messed up real time")
		}

		value1 := processed[len(processed)-1][1]
		value2 := expr[i+1][1]

		result := 0.0
		if get_op(value) == "+" {
			result = value1 + value2
			// DEBUG: fmt.Printf("Simplified %f + %f\n", value1, value2)

		} else {
			result = value1 - value2
			// DEBUG: fmt.Printf("Simplified %f - %f\n", value1, value2)
		}

		processed = processed[:len(processed)-1]
		processed = append(processed, []float64{float64(NUMBER), result})
		skip++
	}

	if len(processed) != 1 || processed[0][0] != float64(NUMBER) {
		// DEBUG: for _, v := range processed {
		// DEBUG: fmt.Printf(" - %d, %f\n", int(v[0]), v[1])
		// DEBUG: }
		return 0, errors.New("not simplified")
	}

	return processed[0][1], nil
}

/*
func main() {
	fmt.Print("Enter expression: ")
	in := bufio.NewReader(os.Stdin)
	line, err := in.ReadString('\n')
	line = line[:len(line)-1]

	res, err := Calc(line)

	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	} else {
		fmt.Printf("Result: %f\n", res)
	}
}
*/
