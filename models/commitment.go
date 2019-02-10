package models

import "time"

// Commitment model whose fields are directly linked to
// the IRIS dashboard PreLoRU
type Commitment struct {
	ID              int64     `json:"id"`
	ActionID        int64     `json:"action_id"`
	Date            time.Time `json:"date"`
	Value           int64     `json:"value"`
	Name            string    `json:"name"`
	ReportID        int64     `json:"report"`
	LapseDate       time.Time `json:"lapse_date"`
	IrisCode        string    `json:"iris_code"`
	CoriolisYear    int64     `json:"coriolis_year"`
	BeneficiaryID   int64     `json:"beneficiary_id"`
	CoriolisEgtCode string    `json:"coriolis_egt_code"`
	CoriolisEgtNum  string    `json:"coriolis_egt_num"`
	CoriolisEgtLine string    `json:"coriolis_egt_line"`
}

// ImportedCommitment is used for one row of imported commitments
// that are inserted into the temporary table before beeing used
// to create commitment table entries
type ImportedCommitment struct {
	Report          string    `json:"report"`
	Action          string    `json:"action"`
	IrisCode        string    `json:"iris_code"`
	CoriolisYear    int       `json:"coriolis_year"`
	CoriolisEgtCode string    `json:"coriolis_egt_code"`
	CoriolisEgtNum  string    `json:"coriolis_egt_num"`
	CoriolisEgtLine string    `json:"coriolis_egt_line"`
	Name            string    `json:"name"`
	Beneficiary     string    `json:"beneficiary"`
	BeneficiaryCode int       `json:"beneficiary_code"`
	Date            time.Time `json:"date"`
	Value           int64     `json:"value"`
	LapseDate       time.Time `json:"lapse_date"`
}
