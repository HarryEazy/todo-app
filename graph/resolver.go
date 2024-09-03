package graph

import "todo_app/graph/generated"

// Resolver struct for dependency injection
type Resolver struct{}

// Query returns the implementation of the generated.QueryResolver interface
func (r *Resolver) Query() generated.QueryResolver {
	return &queryResolver{r}
}

// Mutation returns the implementation of the generated.MutationResolver interface
func (r *Resolver) Mutation() generated.MutationResolver {
	return &mutationResolver{r}
}

type queryResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
