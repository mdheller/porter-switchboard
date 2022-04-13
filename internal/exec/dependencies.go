package exec

import (
	"fmt"

	"github.com/porter-dev/switchboard/pkg/models"
)

type DependencyResolver struct {
	resources  []*models.Resource
	graph      map[string][]string
	resolved   map[string]bool
	unresolved map[string]bool
}

func NewDependencyResolver(resources []*models.Resource) *DependencyResolver {
	return &DependencyResolver{
		resources:  resources,
		graph:      make(map[string][]string),
		resolved:   make(map[string]bool),
		unresolved: make(map[string]bool),
	}
}

func (r *DependencyResolver) Resolve() error {
	// construct dependency graph
	for _, resource := range r.resources {
		// check for duplicate resource
		if _, ok := r.graph[resource.Name]; ok {
			return fmt.Errorf("duplicate resource detected: '%s'", resource.Name)
		}

		r.graph[resource.Name] = append(r.graph[resource.Name], resource.Dependencies...)
	}

	err := r.depResolve(r.resources[0].Name)

	if err != nil {
		return err
	}

	return nil
}

func (r *DependencyResolver) depResolve(name string) error {
	r.unresolved[name] = true

	for _, dep := range r.graph[name] {
		if _, ok := r.graph[dep]; !ok {
			return fmt.Errorf("no such resource as: '%s'", dep)
		}

		if _, ok := r.resolved[dep]; !ok {
			if _, ok = r.unresolved[dep]; ok {
				return fmt.Errorf("circular depedency detected: '%s' -> '%s'", name, dep)
			}
			err := r.depResolve(dep)
			if err != nil {
				return err
			}
		}
	}

	r.resolved[name] = true
	delete(r.unresolved, name)

	return nil
}
