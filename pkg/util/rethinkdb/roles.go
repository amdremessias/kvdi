package rethinkdb

import (
	"github.com/tinyzimmer/kvdi/pkg/auth/types"
	"github.com/tinyzimmer/kvdi/pkg/util/errors"

	rdb "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

// GetAllRoles returns a slice of all the roles in the database.
func (d *rethinkDBSession) GetAllRoles() ([]*types.Role, error) {
	cursor, err := rdb.DB(kvdiDB).Table(rolesTable).Run(d.session)
	if err != nil {
		return nil, err
	}
	roles := make([]*types.Role, 0)
	if cursor.IsNil() {
		return roles, nil
	}
	if err := cursorIntoObjSlice(cursor, &roles); err != nil {
		return nil, err
	}
	for _, role := range roles {
		role.GrantNames = role.RoleGrants().Names()
	}
	return roles, nil
}

// GetRole looks up a role by name.
func (d *rethinkDBSession) GetRole(name string) (*types.Role, error) {
	cursor, err := rdb.DB(kvdiDB).Table(rolesTable).Get(name).Run(d.session)
	if err != nil {
		return nil, err
	}
	if cursor.IsNil() {
		return nil, errors.NewRoleNotFoundError(name)
	}
	role := &types.Role{}
	if err := cursorIntoObj(cursor, role); err != nil {
		return nil, err
	}
	role.GrantNames = role.RoleGrants().Names()
	return role, nil
}

// CreateRole creates a new role in the database.
func (d *rethinkDBSession) CreateRole(role *types.Role) error {
	cursor, err := rdb.DB(kvdiDB).Table(rolesTable).Insert(role).Run(d.session)
	if err != nil {
		return err
	}
	return cursor.Err()
}

// UpdateRole destructively updates a role in the database. Make sure any values
// that you want to stay the same are replicated in the given argument.
func (d *rethinkDBSession) UpdateRole(role *types.Role) error {
	cursor, err := rdb.DB(kvdiDB).Table(rolesTable).Get(role.Name).Update(role).Run(d.session)
	if err != nil {
		return err
	}
	return cursor.Err()
}

// DeleteRole removes the given role from the database.
func (d *rethinkDBSession) DeleteRole(role *types.Role) error {
	cursor, err := rdb.DB(kvdiDB).Table(rolesTable).Get(role.Name).Delete().Run(d.session)
	if err != nil {
		return err
	}
	return cursor.Err()
}
