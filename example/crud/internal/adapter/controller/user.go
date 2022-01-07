// Author: Daniel TAN
// Date: 2021-08-18 23:39:51
// LastEditors: Daniel TAN
// LastEditTime: 2021-10-04 01:27:48
// FilePath: /trinity-micro/example/crud/internal/adapter/controller/user.go
// Description:
package controller

import (
	"context"
	"net/http"
	"sync"

	"github.com/PolarPanda611/trinity-micro"
	"github.com/PolarPanda611/trinity-micro/core/dbx"
	"github.com/PolarPanda611/trinity-micro/core/logx"
	"github.com/PolarPanda611/trinity-micro/core/tracerx"
	"github.com/PolarPanda611/trinity-micro/example/crud/internal/application/dto"
	"github.com/PolarPanda611/trinity-micro/example/crud/internal/application/service"
	"github.com/go-chi/chi/v5"
)

func init() {
	UserControllerPool := &sync.Pool{
		New: func() interface{} {
			return new(userControllerImpl)
		},
	}
	trinity.RegisterInstance("UserController", UserControllerPool)
	trinity.RegisterController("/example-api/v1/{tenant}/users", "UserController",
		trinity.NewRequestMapping("GET", "/", "ListUser",
			logx.ChiLoggerRequest,
			tracerx.ChiOpenTracer(
				tracerx.OperationNameFunc(
					func(r *http.Request) string {
						chiCtx := chi.RouteContext(r.Context())
						return r.Method + "=>" + chiCtx.RoutePattern()
					},
				),
			),
			dbx.SessionHandler,
		),
		trinity.NewRequestMapping("GET", "/{id}", "GetUserByID",
			logx.ChiLoggerRequest,
			tracerx.ChiOpenTracer(
				tracerx.OperationNameFunc(
					func(r *http.Request) string {
						chiCtx := chi.RouteContext(r.Context())
						return r.Method + "=>" + chiCtx.RoutePattern()
					},
				),
			),
			dbx.SessionHandler,
		),
	)
}

type userControllerImpl struct {
	UserSrv service.UserService `container:"autowire:true;resource:UserService"`
}

// GetUserByID godoc
// @Summary      Get Single user information
// @Description  get string by ID
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        tenant   path		string	true									"tenant id"
// @Param        id       path      integer	true  									"user id"
// @Success      200      {object}  httpx.Response{result=dto.GetUserByIDResponse}	"success response"
// @Failure      400,500  {object}  httpx.ErrorResponse 							"error response"
// @Router       /example-api/v1/{tenant}/users/{id} [get]
func (c *userControllerImpl) GetUserByID(ctx context.Context, req *dto.GetUserByIDRequest) (*dto.GetUserByIDResponse, error) {
	return c.UserSrv.GetUserID(ctx, req)
}

// ListUser godoc
// @Summary      list user
// @Description  list user information
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        current_user_id 	header 	int		false	"current user id"
// @Param        tenant   			path   	string  true	"tenant id"
// @Param        pageSize 			query	int  	true  	"page size"			minimum(1)    	maximum(500)
// @Param        current  			query	int  	true  	"page number"		minimum(1)
// @Param        username__ilike 	query 	string 	false 	"username ilike"	minlength(1)  	maxlength(100)
// @Param        age 				query 	string 	false 	"username ilike"  	minlength(1)  	maxlength(100)
// @Success      200      	{object}	httpx.Response{result=dto.ListUserResponse}	"success response"
// @Failure      400,500	{object} 	httpx.ErrorResponse 						"error response"
// @Router       /example-api/v1/{tenant}/users [get]
func (c *userControllerImpl) ListUser(ctx context.Context, req *dto.ListUserRequest) (*dto.ListUserResponse, error) {
	return c.UserSrv.ListUser(ctx, req)
}
