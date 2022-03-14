package io

type Service interface {
	ReadFile(filepath string) (content []byte, err error)
}
