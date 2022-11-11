package errors

import "fmt"

var ErrUnauthorized = fmt.Errorf("unauthorized")
var ErrNoPerm = fmt.Errorf("no permission")
var ErrUnimplemented = fmt.Errorf("unimplemented")
var ErrCannotRevokeLastAdmin = fmt.Errorf("cannot revoke last admin")
