package services

import "gorm.io/gorm"

type ServiceContext struct {
	Node int
}

func (s *ServiceContext) DB() *gorm.DB {
	return Db
}

func (s *ServiceContext) Config() {

}
