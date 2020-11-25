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
	clientWebhookUrl = "localhost:8032"
)

type AgentServer struct {
	prevState string
	responseChannel chan ChannelResponse
}

type ChannelResponse struct {
	Verified bool
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
		Addr: clientWebhookUrl,
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
	state                  string `json:"state"`
	credentialExchangeID   string `json:"credential_exchange_id"`
	credentialID           string
	metadata               controller.CredentialRequestMetadata
	credentialDefinitionID string
	schemaID               string
}

func handleIssueCredential(w http.ResponseWriter, r *http.Request) {

	fmt.Println("INSIDE CLIENT ISSUE CREDENTIAL WEBHOOK")

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
	//if message.state == "offer_received" {
	//
	//	fmt.Println("inside offer_received")
	//
	//	_, err := common.SendRequest_POST(
	//		clientUrl,
	//		"/issue-credential/records/" + message.credentialExchangeID + "/send-request",
	//		nil,
	//	)
	//	if err != nil {
	//		fmt.Printf("Error while sending issue cred post request: %v\n", err)
	//		return
	//	}
	//
	//} else if message.state == "credential_acked" {
	//	fmt.Printf("Stored credential\n", err)
	//}
}

type PresentProofMessage struct {
	State                  string                      `json:"state"`
	PresentationExchangeID string                      `json:"presentation_exchange_id"`
	PresentationRequest    controller.IndyProofRequest `json:"presentation_request"`
}

func (as *AgentServer) handlePresentProof(w http.ResponseWriter, r *http.Request) {

	fmt.Println("INSIDE PRESENT PROOF WEBHOOK - CLIENT")

	var i interface{}
	err := json.NewDecoder(r.Body).Decode(&i)
	if err != nil {
		fmt.Println("Error")
	}

	fmt.Printf("%v\n\n", i)

	//var message PresentProofMessage
	//err := json.NewDecoder(r.Body).Decode(&message)
	//if err != nil {
	//	fmt.Printf("Error while decoding proof request: %v\n", err)
	//	return
	//}
	//
	//fmt.Printf("%+v\n", message)
	//
	////fmt.Printf("%+v\n", message)
	//if message.State == "request_received" {
	//
	//	resp, err := common.SendRequest_GET(
	//		clientUrl,
	//		"/present-proof/records/" + message.PresentationExchangeID + "/credentials",
	//		nil,
	//	)
	//	if err != nil {
	//		fmt.Printf("Error while querying for credentials: %v\n", err)
	//		return
	//	}
	//	defer resp.Body.Close()
	//
	//	var cred interface{}
	//	err = json.NewDecoder(resp.Body).Decode(&cred)
	//	if err != nil {
	//		fmt.Printf("Error while decoding credentials: %v\n", err)
	//		return
	//	}
	//
	//	fmt.Printf("%+v\n", cred)
	//}
}