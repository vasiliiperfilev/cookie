package data

import (
	"testing"
)

type Permissions []int

const (
	PermissionCreateOrder          = 1
	PermissionAcceptOrder          = 2
	PermissionDeclineOrder         = 3
	PermissionFulfillOrder         = 4
	PermissionConfirmFulfillOrder  = 5
	PermissionSupplierChangesOrder = 6
	PermissionClientChangesOrder   = 7
)

func (p Permissions) Include(code int) bool {
	for i := range p {
		if code == p[i] {
			return true
		}
	}
	return false
}

func AssertUserPermissions(t *testing.T, got Permissions, want Permissions) {
	if !EqualArraysContent(got, want) {
		t.Fatalf("Expected equal arrays: want %v, got %v", want, got)
	}
}
