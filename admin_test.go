package controller

import (
	"testing"
	"time"
)

func TestAdminController_RegisterPublicDid(t *testing.T) {

	var seed = Seed()

	// Create a client controller for signing
	cc := NewClientController()

	// Create an admin controller for verifying the signature
	ac := NewAdminController()

	// Register DID of client
	t.Run("Register Did With Ledger", func(t *testing.T) {

		did, err := RegisterDidWithLedger(ac, seed)
		if err != nil {
			t.Errorf("Error occurred while registering did: %v\n", err)
			return
		}
		ac.SetPublicDid(did)

		did, err = RegisterDidWithLedger(cc, seed)
		if err != nil {
			t.Errorf("Error occurred while registering did: %v\n", err)
			return
		}
		cc.SetPublicDid(did)

		t.Logf("admin Public DID:  %v\n", ac.PublicDid())
		t.Logf("client Public DID: %v\n", cc.PublicDid())
	})
}

func TestAdminController_IssueCredential(t *testing.T) {

	ac := NewAdminController()
	cc := NewClientController()

	// Establish connection between client and admin
	t.Run("Establish connection", func(t *testing.T) {

		// create invitation
		inv, err := CreateInvitation(ac)
		if err != nil {
			t.Errorf("Error while creating invitation: %v\n", err)
			return
		}

		t.Logf("Invitation: %+v\n", inv)

		// receive invitation
		_, err = ReceiveInvitation(cc, inv)
		if err != nil {
			t.Errorf("Receive invitation failed: %v\n", err)
		}
	})

	time.Sleep(4 * time.Second)

	// Get connection ID of connection between admin and client (FOR ADMIN)
	t.Run("Get Connection Object", func(t *testing.T) {
		_, err := ac.GetConnection()
		if err != nil {
			t.Errorf("Failed to retrieve connection id: %v\n", err)
		}
		t.Logf("Connection:     %+v\n", ac.Connection)
	})

	// Register DID of client
	t.Run("Register Did With Ledger", func(t *testing.T) {

		did, err := RegisterDidWithLedger(ac, Seed())
		if err != nil {
			t.Errorf("Error occurred while registering did: %v\n", err)
			return
		}
		ac.SetPublicDid(did)
		t.Logf("Public DID : %v\n", did)
	})

	t.Run("Register Schema and Cred Def", func(t *testing.T) {

		schemaID, err := ac.RegisterSchema("schema")
		if err != nil {
			t.Errorf("Error occurred while registering schema: %v\n", err)
			return
		}
		t.Logf("Schema ID : %v\n", schemaID)

		credDefID, err := ac.RegisterCredentialDefinition(schemaID)
		if err != nil {
			t.Errorf("Error occurred while registering cred def: %v\n", err)
			return
		}
		t.Logf("CredDef ID : %v\n", credDefID)
	})

	t.Run("Issue Credential", func(t *testing.T) {

		err := ac.IssueCredential()
		if err != nil {
			t.Errorf("Error occurred while issuing credential: %v\n", err)
			return
		}
	})
}

func TestAdminController_VerifySignature(t *testing.T) {

	var signature string
	var err error
	var message = "Foo bar"
	var seed = Seed()

	// Create a client controller for signing
	cc := NewClientController()

	// Create an admin controller for verifying the signature
	ac := NewAdminController()

	// Register DID of client
	t.Run("Register Did With Ledger", func(t *testing.T) {

		_, err := RegisterDidWithLedger(ac, seed)
		if err != nil {
			t.Errorf("Error occurred while registering did: %v\n", err)
			return
		}

		_, err = RegisterDidWithLedger(cc, seed)
		if err != nil {
			t.Errorf("Error occurred while registering did: %v\n", err)
			return
		}

		t.Logf("admin Public DID:  %v\n", ac.PublicDid())
		t.Logf("client Public DID: %v\n", cc.PublicDid())
	})

	// Register DID of client
	t.Run("Create Signing Did", func(t *testing.T) {

		err = cc.CreateSigningDid()
		if err != nil {
			t.Errorf("Error occurred while registering did: %v\n", err)
			return
		}

		t.Logf("Signing DID: %v\n", cc.SigningDid)
		t.Logf("Signing VK: %v\n", cc.SigningVk)
	})

	// Establish connection between client and admin
	t.Run("Establish connection", func(t *testing.T) {

		// create invitation
		inv, err := CreateInvitation(ac)
		if err != nil {
			t.Errorf("Error while creating invitation: %v\n", err)
			return
		}

		t.Logf("Invitation: %+v\n", inv)

		// receive invitation
		_, err = ReceiveInvitation(cc, inv)
		if err != nil {
			t.Errorf("Receive invitation failed: %v\n", err)
		}
	})

	time.Sleep(2 * time.Second)

	// Get connection ID of connection between admin and client (FOR ADMIN)
	t.Run("Get Connection Object", func(t *testing.T) {
		_, err := ac.GetConnection()
		if err != nil {
			t.Errorf("Failed to retrieve connection id: %v\n", err)
		}
		t.Logf("Connection:     %+v\n", ac.Connection)
	})

	// Get connection ID of connection between admin and client (FOR CLIENT)
	t.Run("Get Connection Object", func(t *testing.T) {
		_, err := cc.GetConnection()
		if err != nil {
			t.Errorf("Failed to retrieve connection id: %v\n", err)
		}
		t.Logf("Connection:     %+v\n", cc.Connection)
	})

	// Sign a message with using the application signing DID
	t.Run("Sign Message", func(t *testing.T) {

		signature, err = cc.SignMessage(message)
		if err != nil {
			t.Errorf("Error occurred during signing: %v\n", err)
		}

		t.Logf("Signature: %s\n", signature)
	})

	time.Sleep(2 * time.Second)

	t.Run("Put Key to Ledger", func(t *testing.T) {

		r := <-PutKeyToLedger(cc, cc.SigningDid, cc.SigningVk)
		//if err != nil {
		//	t.Errorf("Failed to put key to ledger: %v\n", err)
		//}

		t.Logf("Successfully put key to ledger: %v\n", r)
	})

	// Verify signature
	t.Run("Verify Signature", func(t *testing.T) {

		t.Logf("Message: %v\n", message)
		t.Logf("Signature: %v\n", signature)
		t.Logf("Signing DID: %v\n", cc.SigningDid)
		t.Logf("Signing VK: %v\n", cc.SigningVk)

		verified, err := ac.VerifySignature(message, signature, cc.SigningDid, cc.SigningVk)
		if err != nil {
			t.Errorf("Error occurred while attempting to verify signature: %v\n", err)
		}

		t.Logf("Verified: %v\n", verified)
	})
}

func TestAdminController_RequireProof(t *testing.T) {

	ac := NewAdminController()
	cc := NewClientController()

	// Establish connection between client and admin
	t.Run("Establish connection", func(t *testing.T) {

		// create invitation
		inv, err := CreateInvitation(ac)
		if err != nil {
			t.Errorf("Error while creating invitation: %v\n", err)
			return
		}

		t.Logf("Invitation: %+v\n", inv)

		// receive invitation
		_, err = ReceiveInvitation(cc, inv)
		if err != nil {
			t.Errorf("Receive invitation failed: %v\n", err)
		}
	})

	time.Sleep(4 * time.Second)

	// Get connection ID of connection between admin and client (FOR ADMIN)
	t.Run("Get Connection Object", func(t *testing.T) {
		_, err := ac.GetConnection()
		if err != nil {
			t.Errorf("Failed to retrieve connection id: %v\n", err)
		}
		t.Logf("Connection:     %+v\n", ac.Connection)
	})

	// Register DID of client
	t.Run("Register Did With Ledger", func(t *testing.T) {

		_, err := RegisterDidWithLedger(ac, Seed())
		if err != nil {
			t.Errorf("Error occurred while registering did: %v\n", err)
			return
		}
		t.Logf("Public DID : %v\n", ac.PublicDid())
	})

	t.Run("Register Schema and Cred Def", func(t *testing.T) {

		schemaID, err := ac.RegisterSchema("schema")
		if err != nil {
			t.Errorf("Error occurred while registering schema: %v\n", err)
			return
		}
		t.Logf("Schema ID : %v\n", schemaID)

		credDefID, err := ac.RegisterCredentialDefinition(schemaID)
		if err != nil {
			t.Errorf("Error occurred while registering cred def: %v\n", err)
			return
		}
		t.Logf("CredDef ID : %v\n", credDefID)
	})

	t.Run("Issue Credential", func(t *testing.T) {

		err := ac.IssueCredential()
		if err != nil {
			t.Errorf("Error occurred while issuing credential: %v\n", err)
			return
		}
	})

	time.Sleep(2 * time.Second)

	t.Run("Request Proof", func(t *testing.T) {

		err := ac.RequireProof()
		if err != nil {
			t.Errorf("Error while trying to request proof: %v\n", err)
		}

	})

}