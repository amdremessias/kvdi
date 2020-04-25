package rethinkdb

import (
	"github.com/tinyzimmer/kvdi/pkg/auth/types"
	"github.com/tinyzimmer/kvdi/pkg/util/common"
	"github.com/tinyzimmer/kvdi/pkg/util/errors"

	rdb "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

var hashFunc = common.HashPassword

func (d *rethinkDBSession) GetAllUsers() ([]types.User, error) {
	cursor, err := rdb.DB(kvdiDB).Table(usersTable).Map(func(row rdb.Term) interface{} {
		return rdb.Branch(row, row.Merge(func(plan rdb.Term) interface{} {
			return map[string]interface{}{
				"role_ids": rdb.DB(kvdiDB).Table(rolesTable).GetAll(rdb.Args(plan.Field("role_ids"))).CoerceTo("array"),
			}
		}), nil)
	}).Run(d.session)
	if err != nil {
		return nil, err
	}
	users := make([]types.User, 0)
	if cursor.IsNil() {
		return users, nil
	}
	if err := cursorIntoObjSlice(cursor, &users); err != nil {
		return nil, err
	}
	for _, user := range users {
		for _, role := range user.Roles {
			role.GrantNames = role.RoleGrants().Names()
		}
	}
	return users, nil
}

func (d *rethinkDBSession) GetUser(name string) (*types.User, error) {
	cursor, err := rdb.DB(kvdiDB).Table(usersTable).Get(name).Do(func(row rdb.Term) interface{} {
		return rdb.Branch(row, row.Merge(func(plan rdb.Term) interface{} {
			return map[string]interface{}{
				"role_ids": rdb.DB(kvdiDB).Table(rolesTable).GetAll(rdb.Args(plan.Field("role_ids"))).CoerceTo("array"),
			}
		}), nil)
	}).Run(d.session)
	if err != nil {
		return nil, err
	}
	if cursor.IsNil() {
		return nil, errors.NewUserNotFoundError(name)
	}
	user := &types.User{}
	if err := cursorIntoObj(cursor, user); err != nil {
		return nil, err
	}
	for _, role := range user.Roles {
		role.GrantNames = role.RoleGrants().Names()
	}
	return user, nil
}

func (d *rethinkDBSession) CreateUser(user *types.User) error {
	hash, err := hashFunc(user.Password)
	if err != nil {
		return err
	}
	user.PasswordSalt = string(hash)
	cursor, err := rdb.DB(kvdiDB).Table(usersTable).Insert(user).Run(d.session)
	if err != nil {
		return err
	}
	return cursor.Err()
}

func (d *rethinkDBSession) UpdateUser(user *types.User) error {
	role_ids := make([]string, 0)
	for _, role := range user.Roles {
		role_ids = append(role_ids, role.Name)
	}
	cursor, err := rdb.DB(kvdiDB).Table(usersTable).Get(user.Name).Update(map[string]interface{}{
		"role_ids": role_ids,
	}).Run(d.session)
	if err != nil {
		return err
	}
	return cursor.Err()
}

func (d *rethinkDBSession) SetUserPassword(user *types.User, password string) error {
	hash, err := hashFunc(password)
	if err != nil {
		return err
	}
	cursor, err := rdb.DB(kvdiDB).Table(usersTable).Get(user.Name).Update(map[string]interface{}{
		"password": string(hash),
	}).Run(d.session)
	if err != nil {
		return err
	}
	return cursor.Err()
}

func (d *rethinkDBSession) DeleteUser(user *types.User) error {
	cursor, err := rdb.DB(kvdiDB).Table(usersTable).Get(user.Name).Delete().Run(d.session)
	if err != nil {
		return err
	}
	return cursor.Err()
}
