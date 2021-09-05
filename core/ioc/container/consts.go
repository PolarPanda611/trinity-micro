package container

type Keyword string

const (
	CONTAINER      Keyword = "container"
	AUTOWIRE       Keyword = "autowire"
	RESOURCE       Keyword = "resource"
	TAG_SPLITER            = ";"
	TAG_KV_SPLITER         = ":"
	CONTEXT                = "CONTEXT"
)
