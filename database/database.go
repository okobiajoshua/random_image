package database

import "random_image/model"

// Database interface
type Database interface {
	GetImageByID(id int) (*model.Image, error)
	GetAll() ([]model.Image, error)
	AddNewImage(id int, img *model.Image) error
}
