/*
 * @Author: Daniel TAN
 * @Description:
 * @Date: 2021-08-30 16:26:22
 * @LastEditTime: 2021-09-02 17:51:08
 * @LastEditors: Daniel TAN
 * @FilePath: /fr-price-common-pkg/e/err_code.go
 */
package e

// 400 error
const (
	// ErrInvalidRequest
	/**
	 * @Author: Daniel TAN
	 * @Description: when post body, query param , path param invalid, post body validation error
	 */
	ErrInvalidRequest errorCode = 400001
	// ErrInvalidRequest
	/**
	 * @Author: Daniel TAN
	 * @Description: when calling third party API respond with unsuccessful error and without err msg
	 */
	ErrHttpResponseCodeNotSuccess errorCode = 400002
	// ErrVertexAccessDeniedException
	/**
	 * @Author: Daniel TAN
	 * @Description: when calling vertex , the auth info is wrong caused the access denied
	 */
	ErrVertexAccessDeniedException errorCode = 400003
	// ErrVertexNumberFormatException
	/**
	 * @Author: Daniel TAN
	 * @Description: when calling vertex, the number field invalid
	 */
	ErrVertexNumberFormatException errorCode = 400004
	// ErrVertexInvalidAddressException
	/**
	 * @Author: Daniel TAN
	 * @Description: when calling vertex, the address info is invalid
	 */
	ErrVertexInvalidAddressException errorCode = 400005
	// ErrVertexInvalidTaxAreaIdException
	/**
	 * @Author: Daniel TAN
	 * @Description: when calling vertex, the tax area id info is invalid
	 */
	ErrVertexInvalidTaxAreaIdException errorCode = 400006
	// ErrVertexApplicationException
	/**
	 * @Author: Daniel TAN
	 * @Description: when calling vertex, the content is invalid
	 */
	ErrVertexApplicationException errorCode = 400007
	// ErrVertexInvalidCountryException
	/**
	 * @Author: Daniel TAN
	 * @Description: when calling vertex, there country provided is invalid
	 */
	ErrVertexInvalidCountryException errorCode = 400008
	// ErrVertexUnknownException
	/**
	 * @Author: Daniel TAN
	 * @Description: when calling vertex, there has un-caught error
	 */
	ErrVertexUnknownException errorCode = 400009
	//ErrResourceNotFound
	/**
	 * @Author: Daniel TAN
	 * @Description: when http query, http body has some resource not found , return this error
	 * @Description: for example, post {"store":"112201"}, the store not exist
	 */
	ErrResourceNotFound errorCode = 400010
	//ErrAdvisoryLock
	/**
	 * @Author: Daniel TAN
	 * @Description: when trying to lock with advisory error
	 */
	ErrAdvisoryLock errorCode = 400011
	//ErrAdvisoryLock
	/**
	 * @Author: Daniel TAN
	 * @Description:when trying to unlock with advisory error
	 */
	ErrAdvisoryUnlock errorCode = 400012
	//ErrDIParam
	/**
	 * @Author: Daniel TAN
	 * @Description:error when trying to di the param
	 */
	ErrDIParam errorCode = 400013
	//ErrDecodeRequestBody
	/**
	 * @Author: Daniel TAN
	 * @Description:error when decode the request body
	 */
	ErrDecodeRequestBody errorCode = 400014
)

// 404 error
const (
	// ErrRecordNotFound
	/**
	 * @Author: Daniel TAN
	 * @Description: only when get /{id}, the id not found, we response this error
	 */
	ErrRecordNotFound errorCode = 404001
)

// 500 error
const (
	// ErrInternalServer
	/**
	 * @Author: Daniel TAN
	 * @Description: the default error, which is unexpected from the developers
	 */
	ErrInternalServer errorCode = 500001
	// ErrReadResponseBody
	/**
	 * @Author: Daniel TAN
	 * @Description: when read third party response, we met unexpected error while reading the response body
	 */
	ErrReadResponseBody errorCode = 500002
	// ErrDecodeResponseBody
	/**
	 * @Author: Daniel TAN
	 * @Description: when read third party response, we met unexpected error while decode the response body
	 */
	ErrDecodeResponseBody errorCode = 500003
	// ErrExecuteSQL
	/**
	 * @Author: Daniel TAN
	 * @Description: when execute the sql, met unexpected error
	 */
	ErrExecuteSQL errorCode = 500004
	// ErrReadRequestBody
	/**
	 * @Author: Daniel TAN
	 * @Description: when read request body, met unexpected error
	 */
	ErrReadRequestBody errorCode = 500005
)
