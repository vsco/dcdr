package kv

import (
	"errors"

	"github.com/vsco/dcdr/models"
)

var (
	TypeChangeError = errors.New("cannot change existing feature types.")
)

type StoreIFace interface {
	List(prefix string) (models.Features, error)
	Set(f *models.Feature) error
	Get(key string) (*models.Feature, error)
	Delete(key string) error
}
