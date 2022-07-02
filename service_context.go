package services

import (
	"strings"

	"github.com/eqto/config"
	"github.com/eqto/dbm"
)

type ServiceContext struct {
	Node int
}

func (s *ServiceContext) Database() *dbm.Connection {
	return db
}

func (s *ServiceContext) CfgGet(name string) string {
	return config.Get(name)
}

func (s *ServiceContext) CfgGetInt(name string) int {
	return config.GetInt(name)
}

func (s *ServiceContext) CfgGetBool(name string) bool {
	return strings.ToLower(config.Get(name)) == "true"
}
