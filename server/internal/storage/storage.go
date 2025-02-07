package storage

import "errors"

var (
    ErrUserNameNotFound = errors.New("username not found")
    ErrUserNameExists = errors.New("username already exists")
)

type StorageInfoUsers struct {
    Id int
    Username string
    Openkey int
}


type StorageMessages struct {
    Id int
    Message string
}
