package doc

import "context"

type Repository interface {
	CreateDoc(context.Context, *Documentation) (int, error)
}
