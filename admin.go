package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"os"
)

const (
	adminUrl = "http://localhost:8021"
)

type AdminController struct {
	alias      string
	did        string // public did
	SchemaID   string
	CredDefID  string
	Connection Connection
	agentUrl   string
}

func NewAdminController() (*AdminController, error) {

	ac := &AdminController{
		alias: "admin",
		agentUrl: adminUrl,
	}

	_, err := RegisterDidWithLedger(ac, Seed())
	if err != nil {
		return nil, fmt.Errorf("Failed initialization of new client controller: %v\n", err)
	}

	return ac, nil
}

func (ac *AdminController) Alias() string {
	return ac.alias
}

func (ac *AdminController) AgentUrl() string {
	return ac.agentUrl
}

func (ac *AdminController) PublicDid() string {
	return ac.did
}

func (ac *AdminController) SetPublicDid(did string) {
	ac.did = did
}

func (ac *AdminController) ConnectionDid() string {
	return ac.Connection.MyDID
}

type VerifySignatureRequest struct {
	Message string `json:"message"`
	Signature string `json:"signature"`
	MyDid string `json:"my_did"`
	TheirDid string `json:"their_did"`
	SigningDid string `json:"signing_did"`
	SigningVk string `json:"signing_vk"`
}

type VerifySignatureResponse struct {
	Status string `json:"status"`
}

type RegisterSchemaRequest struct {
	Name string `json:"schema_name"`
	Version string `json:"schema_version"`
	Attributes []string `json:"attributes"`
}

type RegisterSchemaResponse struct {
	SchemaID string `json:"schema_id"`
}

// Returns schema ID needed for registering cred def
func (ac *AdminController) RegisterSchema(name string) (string, error) {

	payload := RegisterSchemaRequest{
		Name: name,
		Version: "1.0",
		Attributes: []string{"app_name", "app_id"},
	}

	resp, err := SendRequest_POST(adminUrl, "/schemas", payload)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var schemaResp RegisterSchemaResponse
	err = json.NewDecoder(resp.Body).Decode(&schemaResp)
	if err != nil {
		return "", err
	}

	// just return schema ID; it's all we need to register cred def
	ac.SchemaID = schemaResp.SchemaID
	return schemaResp.SchemaID, nil
}

type RegisterCredDefRequest struct {
	Tag string `json:"tag"`
	SchemaID string `json:"schema_id"`
	SupportRevocation bool `json:"support_revocation"`
	//RevocationRegistrySize string `json:"revocation_registry_size"`
}

type RegisterCredDefResponse struct {
	CredDefID string `json:"credential_definition_id"`
}

// Returns credential definition ID (string)
func (ac *AdminController) RegisterCredentialDefinition(schemaID string) (string, error) {

	payload := RegisterCredDefRequest{
		Tag: "default",
		SchemaID: schemaID,
		SupportRevocation: false,
		//RevocationRegistrySize: "1000",
	}

	// send request for registering schema to agent
	resp, err := SendRequest_POST(adminUrl, "/credential-definitions", payload)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var credDefResp RegisterCredDefResponse
	err = json.NewDecoder(resp.Body).Decode(&credDefResp)
	if err != nil {
		return "", err
	}

	ac.CredDefID = credDefResp.CredDefID
	return credDefResp.CredDefID, nil
}

// Get connection ID of connection with client agent
func (ac *AdminController) GetConnection() (GetConnectionResponse, error) {

	var getConnResp GetConnectionResponse

	resp, err := SendRequest_GET(
		adminUrl,
		"/connections",
		nil,
	)
	if err != nil {
		return getConnResp, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&getConnResp)
	if err != nil {
		return getConnResp, err
	}

	if len(getConnResp.Results) == 0 {
		return getConnResp, errors.New("no connections found")
	}

	// save the connection ID
	ac.Connection = getConnResp.Results[0]
	return getConnResp, nil
}

type CredentialOfferRequest struct {
	ConnID string `json:"connection_id"`
	CredDefID string `json:"cred_def_id"`
	Comment string `json:"comment"`
	AutoRemove bool `json:"auto_remove"`
	//Trace bool `json:"trace"`
	CredProposal CredentialProposal `json:"credential_proposal"`
}

type CredentialProposal struct {
	Type string `json:"@type"`
	Attributes []map[string]interface{} `json:"attributes"`
}

//type CredentialAttribute map[string]interface{}

// issue credential to client agent
func (ac *AdminController) IssueCredential(appName string, appID string) error {

	attrs := []map[string]interface{}{
		{"name": "app_name", "value": appName},
		{"name": "app_id", "value": appID},
	}

	credProposal := CredentialProposal{
		Type:      	"https://didcomm.org/issue-credential/1.0/credential-preview",
		Attributes: attrs,
	}

	offerRequest := CredentialOfferRequest{
		ConnID:      ac.Connection.ConnectionID,
		CredDefID:   ac.CredDefID,
		Comment:     "Issue credential to client agent",
		AutoRemove:  false,
		//Trace:       false,
		CredProposal: credProposal,
	}

	fmt.Fprintf(os.Stdout, "%+v\n", offerRequest)

	_, err := SendRequest_POST(adminUrl, "/issue-credential/send", offerRequest)
	if err != nil {
		return err
	}

	return nil
}

// verify signature provided in transaction proposal
func (ac *AdminController) VerifySignature(message, signature, did, vk string) (bool, error) {

	payload := VerifySignatureRequest{
		Message:   	message,
		Signature: 	signature,
		MyDid: 		ac.Connection.MyDID,
		TheirDid:	ac.Connection.TheirDID,
		SigningDid:	did,
		SigningVk:	vk,
	}

	resp, err := SendRequest_POST(adminUrl, "/connections/verify-transaction", payload)
	if err != nil {
		return false, fmt.Errorf("Error occurred while sending post request: %v\n", err)
	}
	defer resp.Body.Close()

	var verifySigResp VerifySignatureResponse
	err = json.NewDecoder(resp.Body).Decode(&verifySigResp)
	if err != nil {
		return false, fmt.Errorf("Error occurred during json decoding: %v\n", err)
	}

	return verifySigResp.Status == "True", nil
}

type Attribute map[string]string
type reqAttribute map[string]Attribute

type Predicate map[string]interface{}
type reqPredicate map[string]Predicate

type RequireProofRequest struct {
	ConnectionID string           `json:"connection_id"`
	ProofRequest IndyProofRequest `json:"proof_request"`
}

// Request a proof from the client who submitted the transaction
func (ac *AdminController) RequireProof() error {

	nonce, _ := uuid.NewRandom()

	indyProofReq := IndyProofRequest{
		Name:    "simple_proof",
		Version: "1.0",
		Nonce:   nonce.String(),
		ReqAttr: map[string]map[string]interface{}{
			"0_app_name_uuid": {
				"name": "app_name",
				"restrictions": []map[string]string{{"issuer_did": ac.PublicDid()}},
			},
		},
		//ReqPred: map[string]map[string]interface{}{},
		ReqPred: map[string]map[string]interface{}{
			"0_app_id_GE_uuid": {
				"name": "app_id",
				"p_type": ">=",
				"p_value": 30,
				"restrictions": []map[string]string{{"issuer_did": ac.PublicDid()}},
			},
		},
	}

	payload := RequireProofRequest{
		ConnectionID: ac.Connection.ConnectionID,
		ProofRequest: indyProofReq,
	}

	fmt.Fprintf(os.Stdout, "payload: %+v\n", payload)

	_, err := SendRequest_POST(adminUrl, "/present-proof/send-request", payload)
	if err != nil {
		return fmt.Errorf("Failed to send post request to send-request: %v\n", err)
	}

	return nil
}


