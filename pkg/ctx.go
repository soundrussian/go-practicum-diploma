package pkg

import "fmt"

// ContextKey defines type for context value keys
type ContextKey string

// String converts ContextKey to string adding prefix from contextKeyPrefix const
func (c ContextKey) String() string {
	return fmt.Sprintf("%s%s", contextKeyPrefix, string(c))
}

const (
	// contextKeyPrefix defines prefix for values stored in context
	// to prevent collision with other packages
	contextKeyPrefix = "gophermart-"
)
