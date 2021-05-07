package service

import (
	"mime/multipart"
)

func (svc *Service) UploadContracts(fileHeaders []*multipart.FileHeader,id string) error {
	return svc.pool.UploadContracts(fileHeaders,id)
}
