package casbinModels

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

// casbin结构体
type CasbinService struct {
	Enforcer *casbin.Enforcer
	Adapter  *gormadapter.Adapter
}

type CasbinJoinAuditService struct {
	Enforcer *casbin.Enforcer
	Adapter  *gormadapter.Adapter
}
