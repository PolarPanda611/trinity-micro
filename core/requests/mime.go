/*
 * @Author: Daniel TAN
 * @Description:
 * @Date: 2021-08-24 11:26:20
 * @LastEditTime: 2021-08-24 13:13:37
 * @LastEditors: Daniel TAN
 * @FilePath: /fr-price-common-pkg/core/requests/mime.go
 */
package requests

const (
	HeaderMime          = "Content-Type"
	HeaderXFRDate       = "X-FR-Date"
	HeaderDate          = "Date"
	HeaderXFRClientID   = "X-FR-ClientID"
	HeaderAuthorization = "Authorization"
	HeaderXFRPhumacAlgo = "X-FR-Phumac-Algo"
	HeaderKeyContentMD5 = "Content-MD5"
)

const (
	MimeJson    = "application/json"
	MimeTextXML = "text/xml"
	MimeXML     = "application/xml"
)
