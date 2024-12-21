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

package application

import (
	"fmt"

	"github.com/Leo-MathGuy/calcapi/pkg/calcapi"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

//// MARK: Structs

type CalculateRequest struct {
	Expr string `json:"expression"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type ResultResponse struct {
	Result float64 `json:"result"`
}

//// MARK: Server

// todo: delete
func BindJSON(c *gin.Context, obj interface{}) error {
	if err := binding.JSON.Bind(c.Request, obj); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return err
	}
	return nil
}

func ErrorHandler(c *gin.Context, err error, status int) {
	errr := "Internal Server Error"
	if err != nil {
		errr = err.Error()
	}
	fmt.Println("Error: " + errr)
	c.JSON(status, ErrorResponse{errr})
}

func CalculateHandler(c *gin.Context) {
	var requestBody CalculateRequest
	if err := binding.JSON.Bind(c.Request, &requestBody); err != nil {
		ErrorHandler(c, err, 500)
		return
	}

	if result, err := calcapi.Calc(requestBody.Expr); err != nil {
		ErrorHandler(c, err, 422)
	} else {
		c.JSON(200, ResultResponse{result})
	}
}

func RunServer() {
	r := gin.Default()
	r.POST("/api/v1/calculate", CalculateHandler)

	r.Run(":80")
}
