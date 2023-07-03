package data

type StubPermissionModel struct {
	permissions map[int64]Permissions
}

func NewStubPermissionsModel() *StubPermissionModel {
	permissions := map[int64]Permissions{
		1: {
			PermissionAcceptOrder,
			PermissionDeclineOrder,
			PermissionFulfillOrder,
			PermissionSupplierChangesOrder,
		},
		2: {
			PermissionCreateOrder,
			PermissionClientChangesOrder,
			PermissionConfirmFulfillOrder,
		},
	}
	return &StubPermissionModel{permissions: permissions}
}

func (s *StubPermissionModel) GetAllForType(typeId int64) (Permissions, error) {
	if permissions, ok := s.permissions[typeId]; ok {
		return permissions, nil
	} else {
		return nil, ErrRecordNotFound
	}
}
