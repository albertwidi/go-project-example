package project

import (
	"github.com/albertwidi/go_project_example/internal/kothak"
	"github.com/albertwidi/go_project_example/internal/repository/image"
)

// Repositories list
type Repositories struct {
	Image *image.Repository
}

func newRepositories(resources *kothak.Kothak) (*Repositories, error) {
	r := Repositories{}

	// iamge repository
	imageRepo := image.New(resources.MustGetRedis("image"))
	r.Image = imageRepo
	return &r, nil
}
