package models

import (
	"net"
	"time"
)

type Election struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	StartDate   string        `json:"start_date"` // puede ser string o time.Time según cómo quieras manejarlo
	EndDate     string        `json:"end_date"`
	Encrypted   bool          `json:"encrypted"`
	Anonymous   bool          `json:"anonymous"`
	Candidates  []Candidate   `json:"candidates"`  // si tienes un struct Candidate, úsalo
	Voters      []Voter `json:"voters"`      // si tienes un struct Voter, úsalo
	Authorities []ElectionAuthority `json:"authorities"` // si tienes struct Authority, úsalo
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
}

type Candidate struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	PhotoURL       string `json:"photo_url"`
	CandidateOrder int    `json:"candidate_order"`
}

type Voter struct {
	Token      string `json:"token"`
}

type Ballot struct {
	Uuid                    string   `json:"Uuid"`
	Vote                    []string `json:"vote"`
	VotingDeviceFingerprint string   `json:"fingerprint"`
	IPAddress               net.IP   `json:"ip_address"`
	Status                  string   `json:"status"`
}

