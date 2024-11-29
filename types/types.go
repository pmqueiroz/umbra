package types

type UmbraType string

const (
	STR     UmbraType = "<str>"
	CHAR    UmbraType = "<char>"
	NUM     UmbraType = "<num>"
	BOOL    UmbraType = "<bool>"
	HASHMAP UmbraType = "<hashmap>"
	ARR     UmbraType = "<arr>"
	FUN     UmbraType = "<fun>"
	ANY     UmbraType = "<any>"
	NULL    UmbraType = "<null>"
	ENUM    UmbraType = "<enum>"
	UNKNOWN UmbraType = "<unknown>"
)
