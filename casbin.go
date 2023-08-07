package main

import (
	"log"

	sqladapter "github.com/Blank-Xu/sql-adapter"
	"github.com/casbin/casbin/v2"
)

func LoadRBAC() *casbin.Enforcer {
	var e, _ = casbin.NewEnforcer("./model.conf", "./policy.csv")
	return e
}

func LoadABAC() *casbin.Enforcer {
	var e, _ = casbin.NewEnforcer("./model.conf", "./policy.csv")
	return e
}

// const RoleEumn = (
// 	ADMIN : iota,
// 	MANAGER,
// 	USER,
// )

func LoadEnforcerAdapter(a *sqladapter.Adapter) *casbin.Enforcer {
	e, err := casbin.NewEnforcer("casbin/model.conf", a)
	if err != nil {
		panic(err)
	}
	e.AddRoleForUser("admin", "ADMIN")
	e.AddRoleForUser("jaxchow", "MANAGER")
	e.AddRoleForUser("jaxchow", "USER")
	// e.addR
	// Load the policy from DB.
	if err = e.LoadPolicy(); err != nil {
		log.Println("LoadPolicy failed, err: ", err)
	}

	// Save the policy back to DB.
	if err = e.SavePolicy(); err != nil {
		log.Println("SavePolicy failed, err: ", err)
	}
	return e
}
