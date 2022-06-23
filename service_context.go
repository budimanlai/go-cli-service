package services

import (
	"github.com/eqto/dbm"
)

type ServiceContext struct {
	Node int
}

func (s *ServiceContext) Database() *dbm.Connection {
	return db
}

func (s *ServiceContext) Config() {

}
