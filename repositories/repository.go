package repositories

import (
	"errors"
	"gorm.io/gorm"
)

type (
	NotFoundError struct {
	}
)

func (n NotFoundError) Error() string {
	return "record not found"
}

type IRepository[Model interface{}] interface {
	Create(*Model) error
	Update(*Model) error
	Delete(*Model) error
	Get(id uint64) (*Model, error)
	All(data []Model)
}

type gormRepository[m interface{}] struct {
	db *gorm.DB
}

func (g gormRepository[m]) Create(model *m) error {
	res := g.db.Create(model)
	return res.Error
}

func (g gormRepository[m]) Update(model *m) error {
	res := g.db.Save(model)
	return res.Error
}

func (g gormRepository[m]) Delete(model *m) error {
	res := g.db.Delete(model)
	return res.Error
}

func (g gormRepository[m]) Get(id uint64) (*m, error) {
	var x m
	res := g.db.First(&x, id)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, NotFoundError{}
	} else {
		return &x, res.Error
	}
}

func (g gormRepository[m]) All(data []m) {
	g.db.Find(&data)
}
