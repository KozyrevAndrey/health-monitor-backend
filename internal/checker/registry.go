package checker

import (
	"fmt"
	"sync"

	"health-monitor/internal/domain"
)

type Registry struct {
	mu       sync.RWMutex
	checkers map[domain.TargetType]domain.Checker
}

func NewRegistry() *Registry {
	return &Registry{
		checkers: make(map[domain.TargetType]domain.Checker),
	}
}

func (r *Registry) Register(checker domain.Checker) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.checkers[checker.Type()] = checker
}

func (r *Registry) Get(targetType domain.TargetType) (domain.Checker, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	checker, ok := r.checkers[targetType]
	if !ok {
		return nil, fmt.Errorf("checker not found for type: %s", targetType)
	}

	return checker, nil
}

func (r *Registry) List() []domain.TargetType {
	r.mu.RLock()
	defer r.mu.RUnlock()

	types := make([]domain.TargetType, 0, len(r.checkers))
	for targetType := range r.checkers {
		types = append(types, targetType)
	}

	return types
}

func NewDefaultRegistry() *Registry {
	registry := NewRegistry()
	registry.Register(NewHTTPChecker())
	registry.Register(NewTCPChecker())
	registry.Register(NewDNSChecker())
	return registry
}
