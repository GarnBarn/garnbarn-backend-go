package service

import (
	"github.com/GarnBarn/garnbarn-backend-go/model"
	"github.com/GarnBarn/garnbarn-backend-go/repository"
)

type Example interface {
	Example() error
}

type example struct {
	exampleRepository repository.Example
}

func NewExampleService(exampleRepository repository.Example) Example {
	return &example{
		exampleRepository: exampleRepository,
	}
}

func (e *example) Example() error {
	return e.exampleRepository.Example(&model.ExampleDB{
		Name:  "Hello",
		Value: 1,
	})
}
