package http

import (
	"context"
	"sync"

	"github.com/PolarPanda611/trinity-micro"
	"github.com/PolarPanda611/trinity-micro/core/httpx"
	"github.com/PolarPanda611/trinity-micro/example/crud/internal/application/dto"
	"github.com/PolarPanda611/trinity-micro/example/crud/internal/application/service"
)

func init() {
	UserControllerPool := &sync.Pool{
		New: func() interface{} {
			return new(userControllerImpl)
		},
	}
	trinity.RegisterInstance("UserController", UserControllerPool)
	trinity.RegisterController("/example-api/v1/{tenant}/users", "UserController",
		trinity.NewRequestMapping("GET", "", "ListUser",
			apiHandler...,
		),
		trinity.NewRequestMapping("GET", "/{id}", "GetUserByID",
			apiHandler...,
		),
		trinity.NewRequestMapping("POST", "", "CreateUser",
			apiHandler...,
		),
		trinity.NewRequestMapping("PATCH", "/{id}", "UpdateUser",
			apiHandler...,
		),
	)
}

type userControllerImpl struct {
	UserSrv service.UserService `container:"autowire:true;resource:UserService"`
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
// @Success      200      	{object}	httpx.SuccessResponse{result=dto.ListUserResponse}	"success response"
// @Failure      400,500	{object} 	httpx.ErrorResponse 						"error response"
// @Router       /example-api/v1/{tenant}/users [get]
func (c *userControllerImpl) ListUser(ctx context.Context, req *dto.ListUserRequest) (*dto.ListUserResponse, error) {
	return c.UserSrv.ListUser(ctx, req)
}

// GetUserByID godoc
// @Summary      Get Single user information
// @Description  get string by ID
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        tenant   path		string	true									"tenant id"
// @Param        id       path      integer	true  									"user id"
// @Success      200      {object}  httpx.SuccessResponse{result=dto.UserInfoResponse}	"success response"
// @Failure      400,500  {object}  httpx.ErrorResponse 							"error response"
// @Router       /example-api/v1/{tenant}/users/{id} [get]
func (c *userControllerImpl) GetUserByID(ctx context.Context, req *dto.GetUserByIDRequest) (*dto.UserInfoResponse, error) {
	return c.UserSrv.GetUserID(ctx, req)
}

// CreateUser godoc
// @Summary      Get Single user information
// @Description  get string by ID
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        tenant   path		string	true										"tenant id"
// @Success      201      {object}  httpx.SuccessResponse{result=dto.UserInfoResponse}	"success response"
// @Failure      400,500  {object}  httpx.ErrorResponse 								"error response"
// @Router       /example-api/v1/{tenant}/users [post]
func (c *userControllerImpl) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserInfoResponse, error) {
	res, err := c.UserSrv.CreateUser(ctx, req)
	if err != nil {
		return nil, err
	}
	httpx.SetHttpStatusCode(ctx, 201)
	return res, nil
}

// CreateUser godoc
// @Summary      Get Single user information
// @Description  get string by ID
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        tenant   path		string	true										"tenant id"
// @Param        id       path      integer	true  									"user id"
// @Success      200      {object}  httpx.SuccessResponse{result=dto.UserInfoResponse}	"success response"
// @Failure      400,500  {object}  httpx.ErrorResponse 								"error response"
// @Router       /example-api/v1/{tenant}/users{id} [patch]
func (c *userControllerImpl) UpdateUser(ctx context.Context, req *dto.UpdateUserRequest) error {
	return c.UserSrv.UpdateUser(ctx, req)

}
