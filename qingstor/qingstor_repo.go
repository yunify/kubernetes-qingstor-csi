package qingstor

import (
	qs "github.com/yunify/qingstor-sdk-go/service"
	"github.com/yunify/qingstor-sdk-go/config"

)

//QingstorRepository qingstor bucket repository
type QingstorRepository struct {
	qsService *qs.Service
}

//NewQingstorRepository Create QingStor repository object
func NewQingstorRepository(config *config.Config)(*QingstorRepository,error) {
	if service,err:=qs.Init(config);err == nil {
		return &QingstorRepository{
			qsService:service,
		},nil
	} else {
		return nil,err
	}
}

