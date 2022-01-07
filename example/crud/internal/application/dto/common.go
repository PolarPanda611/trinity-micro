// Author: Daniel TAN
// Date: 2021-10-04 00:02:51
// LastEditors: Daniel TAN
// LastEditTime: 2021-10-04 00:02:52
// FilePath: /trinity-micro/example/crud/internal/application/dto/common.go
// Description:
package dto

type TenantRequest struct {
	Tenant string `path_param:"tenant"`
}

type PageRequest struct {
	PageSize *uint `query_param:"page_size" validate:"required,min=1,max=500"`
	PageNum  *uint `query_param:"page_num"  validate:"required,min=1"`
}
