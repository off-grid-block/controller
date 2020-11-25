package controller

import (
	"testing"
)

func TestConnection(t *testing.T) {

	ac := NewAdminController()
	cc := NewClientController()

	t.Run("Create Connection Invitation", func(t *testing.T) {

		invitation, err := CreateInvitation(ac)
		if err != nil {
			t.Errorf("Error while creating invitation: %v\n", err)
			return
		}
		t.Logf("Invitation: %+v\n", invitation)

		conn, err := ReceiveInvitation(cc, invitation)
		if err != nil {
			t.Errorf("Receive invitation failed: %v\n", err)
		}
		t.Logf("Connection: %+v\n", *conn)
	})

}