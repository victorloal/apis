package models

import (
	"net"
	"time"

	"github.com/google/uuid"
)

type Election struct {
	Uuid        uuid.UUID `json:"uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Status      string    `json:"status"`
	IsActive    bool      `json:"is_active"`
}

type TallyResult struct {
	TotalVotes int       `json:"total_votes"`
	Results    []string  `json:"results"`
	Skey       []byte    `json:"s_key"`
	ComputedAt time.Time `json:"computed_at"`
}

type HomomorphicKey struct {
	PKey   []byte `json:"p_key"`
	Params []byte `json:"params"`
}

type AuditConfig struct {
	EnableBallotAudit bool `json:"enable_ballot_audit"`
	EnableAccessLogs  bool `json:"enable_access_logs"`
}

type ElectionAuthority struct {
	CC       uint   `json:"cc"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
	SKey     []byte `json:"s_key"`
	IsActive bool   `json:"is_active"`
}

type Candidate struct {
	Name           string `json:"name"`
	Description    string `json:"description"`
	PhotoURL       string `json:"photo_url"`
	CandidateOrder int    `json:"candidate_order"`
}

type Voter struct {
	Token      string `json:"token"`
	VoteStatus bool   `json:"vote_status"`
	IsActive   bool   `json:"is_active"`
}

type Ballot struct {
	Uuid                    string   `json:"Uuid"`
	Vote                    []string `json:"vote"`
	VotingDeviceFingerprint string   `json:"fingerprint"`
	IPAddress               net.IP   `json:"ip_address"`
	Status                  string   `json:"status"`
}

