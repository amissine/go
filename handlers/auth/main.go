package auth

import (
//
)

// Strategy defines strategy for handling auth requests.
// Additionally for all methods we can add memo preimage to params
// list for ease of use or provide helper methods to extract memo
// from AuthData.
type Strategy interface {
	// SanctionsCheck performs AML sanctions check of the sender.
	SanctionsCheck(data AuthData, response *AuthResponse) error
	// GetUserData check if user data is required and if so decides
	// whether to allow access to customer data or not.
	GetUserData(data AuthData, response *AuthResponse) error
	// PersistTransaction save authorized transaction to persistent storage so
	// memo preimage can be fetched when a transaction is sent.
	PersistTransaction(data AuthData) error
}

type Handler struct {
	Strategy Strategy
}

// AuthStatus represents auth status returned by Auth Server
type AuthStatus string

const (
	// AuthStatusOk is returned when authentication was successful
	AuthStatusOk AuthStatus = "ok"
	// AuthStatusPending is returned when authentication is pending
	AuthStatusPending AuthStatus = "pending"
	// AuthStatusDenied is returned when authentication was denied
	AuthStatusDenied AuthStatus = "denied"
)

// AuthRequest represents auth request sent to compliance server
type AuthRequest struct {
	// Stringified AuthData JSON object
	Data string `name:"data" required:""`
	// Signature of sending FI
	Signature string `name:"sig" required:""`
}

// AuthData represents how AuthRequest.Data field looks like. It is Marshalled because of the attached signature.
type AuthData struct {
	// The stellar address of the customer that is initiating the send.
	Sender string `json:"sender"`
	// If the caller needs the recipient's AML info in order to send the payment.
	NeedInfo bool `json:"need_info"`
	// The transaction that the sender would like to send in XDR format. This transaction is unsigned.
	Tx string `json:"tx"`
	// The full text of the memo the hash of this memo is included in the transaction.
	Memo string `json:"memo"`
}

// AuthResponse represents response sent by auth server
type AuthResponse struct {
	// If this FI is willing to share AML information or not. {ok, denied, pending}
	InfoStatus AuthStatus `json:"info_status"`
	// If this FI is willing to accept this transaction. {ok, denied, pending}
	TxStatus AuthStatus `json:"tx_status"`
	// (only present if info_status is ok) JSON of the recipient's AML information. in the Stellar memo convention
	DestInfo string `json:"dest_info,omitempty"`
	// (only present if info_status or tx_status is pending) Estimated number of seconds till the sender can check back for a change in status. The sender should just resubmit this request after the given number of seconds.
	Pending int `json:"pending,omitempty"`
}
