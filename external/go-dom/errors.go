package dom

import (
	"fmt"
)

const (
	INDEX_SIZE_ERR              = "INDEX_SIZE"
	DOMSTRING_SIZE_ERR          = "DOMSTRING_SIZE"
	HIERARCHY_REQUEST_ERR       = "HIERARCHY_REQUEST"
	WRONG_DOCUMENT_ERR          = "WRONG_DOCUMENT"
	INVALID_CHARACTER_ERR       = "INVALID_CHARACTER"
	NO_DATA_ALLOWED_ERR         = "NO_DATA_ALLOWED"
	NO_MODIFICATION_ALLOWED_ERR = "NO_MODIFICATION_ALLOWED"
	NOT_FOUND_ERR               = "NOT_FOUND"
	NOT_SUPPORTED_ERR           = "NOT_SUPPORTED"
	INUSE_ATTRIBUTE_ERR         = "INUSE_ATTRIBUTE"
	INVALID_STATE_ERR           = "INVALID_STATE"
	SYNTAX_ERR                  = "SYNTAX"
	INVALID_MODIFICATION_ERR    = "INVALID_MODIFICATION"
	NAMESPACE_ERR               = "NAMESPACE"
	INVALID_ACCESS_ERR          = "INVALID_ACCESS"
	VALIDATION_ERR              = "VALIDATION"
	TYPE_MISMATCH_ERR           = "TYPE_MISMATCH"
	SECURITY_ERR                = "SECURITY"
	NETWORK_ERR                 = "NETWORK"
	ABORT_ERR                   = "ABORT"
	URL_MISMATCH_ERR            = "URL_MISMATCH"
	QUOTA_EXCEEDED_ERR          = "QUOTA_EXCEEDED"
	TIMEOUT_ERR                 = "TIMEOUT"
	INVALID_NODE_TYPE_ERR       = "INVALID_NODE_TYPE"
	DATA_CLONE_ERR              = "DATA_CLONE"
)

type ErrDOM struct {
	Typ string
	Msg string
	Op  string
}

func (e ErrDOM) Error() string {
	return fmt.Sprintf("%s.%s: %s", e.Op, e.Typ, e.Msg)
}

func ErrHierarchyRequest(op, msg string) ErrDOM {
	return ErrDOM{
		Typ: HIERARCHY_REQUEST_ERR,
		Msg: msg,
		Op:  op,
	}
}
