package casbinModels

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

func NewCasbinService(db *gorm.DB) (*CasbinService, error) {
	// 创建适配器
	a, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}

	// 创建模型
	//m, err := model.NewModelFromFile("models/model.conf")
	m, err := model.NewModelFromString(`
[request_definition]
r = sub, obj


[policy_definition]
p = sub, obj

[role_definition]
g = _,_

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && r.obj == p.obj`)
	if err != nil {
		return nil, err
	}

	// 创建执行器
	e, err := casbin.NewEnforcer(m, a)
	if err != nil {
		return nil, err
	}
	return &CasbinService{Enforcer: e, Adapter: a}, nil
}

// RolePolicy 对应于 `CasbinRule` 表中的 (v0, v1)
type RolePolicy struct {
	RoleName string `gorm:"column:v0"`
	MenuId   string `gorm:"column:v1"`
}

// GetRoles 获取所有角色组
func (c *CasbinService) GetRoles() ([]string, error) {
	return c.Enforcer.GetAllRoles()
}

// GetRolePolicy 获取所有角色组权限
func (c *CasbinService) GetRolePolicy() (roles []RolePolicy, err error) {
	err = c.Adapter.GetDb().Model(&gormadapter.CasbinRule{}).Where("ptype = 'p'").Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return
}

// CreateRolePolicy 创建角色组权限, 已有的会忽略
func (c *CasbinService) CreateRolePolicy(r RolePolicy) error {
	// 不直接操作数据库，利用enforcer简化操作
	err := c.Enforcer.LoadPolicy()
	if err != nil {
		return err
	}
	_, err = c.Enforcer.AddPolicy(r.RoleName, r.MenuId)
	if err != nil {
		return err
	}
	return c.Enforcer.SavePolicy()
}

// UpdateRolePolicy 修改角色组权限
func (c *CasbinService) UpdateRolePolicy(old, new RolePolicy) error {
	_, err := c.Enforcer.UpdatePolicy([]string{old.RoleName, old.MenuId},
		[]string{new.RoleName, new.MenuId})
	if err != nil {
		return err
	}
	return c.Enforcer.SavePolicy()
}

// DeleteRolePolicy 删除角色组权限
func (c *CasbinService) DeleteRolePolicy(r RolePolicy) error {
	_, err := c.Enforcer.RemovePolicy(r.RoleName, r.MenuId)
	if err != nil {
		return err
	}
	return c.Enforcer.SavePolicy()
}

// CanAccess 验证用户权限
func (c *CasbinService) CanAccess(username, menuId string) (ok bool, err error) {
	return c.Enforcer.Enforce(username, menuId)
}

// AddPolicy 添加角色
func (c *CasbinService) AddPolicy(code, role string) (ok bool, err error) {
	ok, err = c.Enforcer.AddRoleForUser(code, role)
	if err != nil {
		return false, err
	}
	return true, nil
}
