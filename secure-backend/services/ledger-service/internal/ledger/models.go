package ledger

import "time"

type Transaction struct {
	ID     string `json:"id"`
	Amount int    `json:"amount"`
	From   string `json:"from"`
	To     string `json:"to"`
}

type TransactionResponse struct {
	Status    string    `json:"status"`
	TxID      string    `json:"tx_id"`
	Timestamp time.Time `json:"timestamp"`
	Ledger    string    `json:"ledger"`
}

type HealthResponse struct {
	Status        string `json:"status"`
	Service       string `json:"service"`
	Version       string `json:"version"`
	TotalTxCount  int    `json:"total_tx_count"`
}