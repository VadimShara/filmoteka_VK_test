package api

import (
	"vk-test-task/api/rest/handlers"
)

type Container struct {
	Resolver *handlers.Resolver
}

func NewContainer(
	resolver *handlers.Resolver,
) Container {
	return Container{
		Resolver: resolver,
	}
}
