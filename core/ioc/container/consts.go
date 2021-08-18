package container

type Keyword string

const (
	CONTAINER      Keyword = "container"
	AUTOWIRE       Keyword = "autowire"
	AUTOFREE       Keyword = "autofree"
	RESOURCE       Keyword = "resource"
	TAG_SPLITER            = ";"
	TAG_KV_SPLITER         = ":"
)
