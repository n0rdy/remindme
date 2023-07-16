package idresolver

import "sync"

type IdResolver struct {
	currentId int
	mu        sync.Mutex
}

func (resolver *IdResolver) Next() int {
	resolver.mu.Lock()
	defer resolver.mu.Unlock()

	id := resolver.currentId
	resolver.currentId += 1
	return id
}

func NewIdResolver() IdResolver {
	return IdResolver{
		currentId: 1,
		mu:        sync.Mutex{},
	}
}
