// Author: Daniel TAN
// Date: 2021-10-04 00:02:51
// LastEditors: Daniel TAN
// LastEditTime: 2021-10-04 00:02:52
// FilePath: /trinity-micro/example/crud/internal/application/dto/common.go
// Description:
package dto

type CommonRequest struct {
	Tenant string `path_param:"tenant"`
}
