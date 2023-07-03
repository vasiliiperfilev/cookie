package data

import (
	"testing"
)

type Permission int

const (
	PermissionCreateOrder          Permission = 1
	PermissionAcceptOrder          Permission = 2
	PermissionDeclineOrder         Permission = 3
	PermissionFulfillOrder         Permission = 4
	PermissionConfirmFulfillOrder  Permission = 5
	PermissionSupplierChangesOrder Permission = 6
	PermissionClientChangesOrder   Permission = 7
)

func AssertUserPermissions(t *testing.T, got []Permission, want []Permission) {
	if !EqualArrays(got, want) {
		t.Fatalf("Expected equal arrays: want %v, got %v", want, got)
	}
}
