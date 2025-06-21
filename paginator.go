package pokesdk

import (
	"context"
	"sync"
)

// PageableResource is an interface that defines a resource that can be paginated.
// NOTE: In the future this could be expanded to support retrieving the previous page as well.
type PageableResource interface {
	GetNextURL() string
}

// Paginator is a generic type that can be used to paginate through resources that implement the PageableResource interface.
type Paginator[T PageableResource] struct {
	sync.Mutex

	next  string
	fetch func(context.Context, string) (T, error)
	done  bool
}

// Page represents a single page of results from a paginated resource.
// It contains the result of type T and an error if one occurred during fetching.
type Page[T PageableResource] struct {
	Result T
	Error  error
}

// NewPaginator creates a new Paginator instance for the given type.
func NewPaginator[T PageableResource](start string, fetch func(context.Context, string) (T, error)) *Paginator[T] {
	return &Paginator[T]{
		next:  start,
		fetch: fetch,
	}
}

// All returns a channel that will yield all pages of results from the paginator.
// The channel will be closed when all pages have been fetched or if the context is done.
func (p *Paginator[T]) All(ctx context.Context) <-chan *Page[T] {
	pages := make(chan *Page[T], 1)

	go func() {
		defer close(pages)

		select {
		case <-ctx.Done():
			return
		default:
		}

		for {
			page := p.Next(ctx)
			if page == nil {
				return
			}
			pages <- page
		}
	}()

	return pages
}

// Next fetches the next page of results from the paginator.
// It returns a Page containing the results or an error if one occurred.
// Note that it will return nil if there are no more pages to fetch, so the caller should check for nil
func (p *Paginator[T]) Next(ctx context.Context) *Page[T] {
	p.Lock()
	defer p.Unlock()

	if p.next == "" {
		return nil
	}
	page, err := p.fetch(ctx, p.next)
	if err != nil {
		return &Page[T]{Error: err}
	}

	next := page.GetNextURL()
	p.next = next

	return &Page[T]{Result: page}
}
