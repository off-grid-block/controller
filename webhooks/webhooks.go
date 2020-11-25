package webhooks

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/off-grid-block/deon-controller-go"
	"log"
	"net/http"
	"time"
)

const (
	adminWebhookUrl = "localhost:8022"
)

type AgentServer struct {
	prevState string
	responseChannel chan ChannelResponse
}

func ListenWebhooks() {

	as := AgentServer{}

	r := mux.NewRouter()
	webhooks := r.PathPrefix("/topic").Subrouter()

	// register webhooks
	webhooks.HandleFunc("/connections/", handleConnections).Methods("POST")
	webhooks.HandleFunc("/basicmessages/", handleMessages).Methods("POST")
	webhooks.HandleFunc("/issue_credential/", handleIssueCredential).Methods("POST")
	webhooks.HandleFunc("/present_proof/", as.handlePresentProof).Methods("POST")

	srv := &http.Server{
		Handler: r,
		Addr: adminWebhookUrl,
		WriteTimeout: 15 * time.Second,
		ReadTimeout: 15 * time.Second,
	}

	fmt.Printf("Listening for webhooks on %v...\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	fmt.Println("INSIDE HANDLE CONNECTIONS WEBHOOK")

	var i interface{}
	err := json.NewDecoder(r.Body).Decode(&i)
	if err != nil {
		fmt.Println("Error")
	}

	fmt.Printf("%v\n\n", i)
}

func handleMessages(w http.ResponseWriter, r *http.Request) {
	fmt.Println("INSIDE HANDLE MESSAGES WEBHOOK")

	var i interface{}
	err := json.NewDecoder(r.Body).Decode(&i)
	if err != nil {
		fmt.Println("Error")
	}

	fmt.Printf("%v\n\n", i)
}

type IssueCredentialMessage struct {
	State string `json:"state"`
	CredExID string `json:"credential_exchange_id"`
}

type IssueCredentialRequest struct {
	Comment     string                        `json:"comment"`
	CredPreview controller.CredentialProposal `json:"credential_proposal"`
}

func handleIssueCredential(w http.ResponseWriter, r *http.Request) {

	fmt.Println("INSIDE ADMIN ISSUE CREDENTIAL WEBHOOK")

	var i interface{}
	err := json.NewDecoder(r.Body).Decode(&i)
	if err != nil {
		fmt.Println("Error")
	}

	fmt.Printf("%v\n\n", i)

	//var message IssueCredentialMessage
	//err := json.NewDecoder(r.Body).Decode(&message)
	//if err != nil {
	//	fmt.Printf("Error while decoding issue cred message: %v\n", err)
	//	return
	//}
	//
	//credentialExchangeID := message.CredExID
	//
	//if message.State == "request_received" {
	//
	//	fmt.Println("inside request_received")
	//
	//	attrs := []CredentialAttribute{
	//		{"name": "app_name", "value": "voter"},
	//		{"name": "app_id", "value": "101"},
	//	}
	//
	//	credPreview := CredentialProposal{
	//		Type:      	"https://didcomm.org/issue-credential/1.0/credential-preview",
	//		Attributes: attrs,
	//	}
	//
	//	payload := IssueCredentialRequest{
	//		Comment: "comment",
	//		CredPreview: credPreview,
	//	}
	//
	//	_, err := common.SendRequest_POST(
	//		adminUrl,
	//		"/issue-credential/records/" + credentialExchangeID + "/issue",
	//		payload,
	//	)
	//	if err != nil {
	//		fmt.Printf("Error while decoding issue cred message: %v\n", err)
	//		return
	//	}
	//
	//}

}

type ProofMessage struct {
	PresentationExchangeID string `json:"presentation_exchange_id"`
	//PresReq PresentationRequest
	//Pres Presentation
}

type VerifyProofResponse struct {
	Verified string `json:"verified"`
}

type ChannelResponse struct {
	Verified bool
}

func (as *AgentServer) handlePresentProof(w http.ResponseWriter, r *http.Request) {

	fmt.Println("INSIDE PRESENT PROOF WEBHOOK - ADMIN")

	var i interface{}
	err := json.NewDecoder(r.Body).Decode(&i)
	if err != nil {
		fmt.Println("Error")
	}

	fmt.Printf("%v\n\n", i)

	//// decode proof details provided by client
	//var proofMessage ProofMessage
	//err := json.NewDecoder(r.Body).Decode(&proofMessage)
	//if err != nil {
	//	fmt.Fprintf(os.Stderr,"Error while decoding proof message: %v\n", err)
	//	//http.Error(w, "Failed to decode proof message", 500)
	//}
	//
	//// send proof details to admin agent for verification
	//resp, err := common.SendRequest_POST(
	//	adminUrl,
	//	"/present-proof/records/" + proofMessage.PresentationExchangeID + "/verify-presentation", nil)
	//if err != nil {
	//	fmt.Fprintf(os.Stderr,"Error while sending post request: %v\n", err)
	//	//http.Error(w, "Error while sending post request", 500)
	//}
	//defer resp.Body.Close()
	//
	//// decode verification response, retrieve relevant attributes
	//var proofResult VerifyProofResponse
	//err = json.NewDecoder(resp.Body).Decode(&proofResult)
	//if err != nil {
	//	fmt.Fprintf(os.Stderr,"Error while decoding proof response: %v\n", err)
	//	//http.Error(w, "Error while decoding proof response", 500)
	//}
	//
	//if proofResult.Verified != "true" {
	//	as.responseChannel <- ChannelResponse{Verified: false}
	//	fmt.Fprintf(os.Stdout,"Proof is false: %+v\n", proofResult)
	//	//http.Error(w, "Failed to verify proof", 500)
	//} else {
	//
	//	//presReq := proofMessage.PresentationRequest
	//	//pres := proofMessage.Presentation
	//	//
	//	//var attrs map[string]string
	//	//
	//	//for ref, attrSpec := range presReq.RequestedAttributes {
	//	//	attrs[attrSpec.Name] = pres.RequestedProof.RevealedAttrs
	//	//}
	//
	//	as.responseChannel <- ChannelResponse{Verified: true}
	//}
}