package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
)

// PaymentDemandLine is used to decode one line of payments demands batch
type PaymentDemandLine struct {
	IrisCode        string     `json:"IrisCode"`
	IrisName        string     `json:"IrisName"`
	CommitmentDate  int64      `json:"CommitmentDate"`
	BeneficiaryCode int64      `json:"BeneficiaryCode"`
	DemandNumber    int64      `json:"DemandNumber"`
	DemandDate      int64      `json:"DemandDate"`
	ReceiptDate     int64      `json:"ReceiptDate"`
	DemandValue     int64      `json:"DemandValue"`
	CsfDate         NullInt64  `json:"CsfDate"`
	CsfComment      NullString `json:"CsfComment"`
	DemandStatus    string     `json:"DemandStatus"`
	StatusComment   NullString `json:"StatusComment"`
}

// PaymentDemandBatch embeddes an array of PaymentDemandLine for dedicated query
type PaymentDemandBatch struct {
	Lines      []PaymentDemandLine `json:"PaymentDemand"`
	ImportDate time.Time           `json:"ImportDate"`
}

// PaymentDemand model
type PaymentDemand struct {
	ID              int64      `json:"ID"`
	ImportDate      time.Time  `json:"ImportDate"`
	IrisCode        string     `json:"IrisCode"`
	IrisName        string     `json:"IrisName"`
	BeneficiaryID   int64      `json:"BeneficiaryCode"`
	Beneficiary     string     `json:"Beneficiary"`
	DemandNumber    int64      `json:"DemandNumber"`
	DemandDate      time.Time  `json:"DemandDate"`
	ReceiptDate     time.Time  `json:"ReceiptDate"`
	DemandValue     int64      `json:"DemandValue"`
	CsfDate         NullTime   `json:"CsfDate"`
	CsfComment      NullString `json:"CsfComment"`
	DemandStatus    string     `json:"DemandStatus"`
	StatusComment   NullString `json:"StatusComment"`
	Excluded        NullBool   `json:"Excluded"`
	ExcludedComment NullString `json:"ExcludedComment"`
	ProcessedDate   NullTime   `json:"ProcessedDate"`
}

// PaymentDemands embeddes an array of PaymentDemand for json export and dedicated
// queries
type PaymentDemands struct {
	Lines []PaymentDemand `json:"PaymentDemand"`
}

// PaymentDemandCount model
type PaymentDemandCount struct {
	Date         time.Time `json:"Date"`
	UnProcessed  int64     `json:"Unprocessed"`
	UnControlled int64     `json:"Uncontrolled"`
}

// PaymentDemandCounts embeddes an array of PaymentDemandCount for json export
// and the dedicated query
type PaymentDemandCounts struct {
	Lines []PaymentDemandCount `json:"PaymentDemandCount"`
}

// Update set excluded fields in the database
func (p *PaymentDemand) Update(db *sql.DB) error {
	res, err := db.Exec(`UPDATE payment_demands SET excluded=$1,excluded_comment=$2
	WHERE id=$3`, p.Excluded, p.ExcludedComment, p.ID)
	if err != nil {
		return fmt.Errorf("update %v", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected %v", err)
	}
	if count == 0 {
		return fmt.Errorf("demande de paiement introuvable")
	}
	return nil
}

// GetAll fetches all payment demand from database
func (p *PaymentDemands) GetAll(db *sql.DB) error {
	rows, err := db.Query(`SELECT p.id,p.import_date,p.iris_code,p.iris_name,
	p.beneficiary_id,b.name,p.demand_number,p.demand_date,p.receipt_date,
	p.demand_value,p.csf_date,p.csf_comment,p.demand_status,p.status_comment,
	p.excluded,p.excluded_comment,p.processed_date
	FROM payment_demands p
	JOIN beneficiary b on  b.id=p.beneficiary_id`)
	if err != nil {
		return fmt.Errorf("select %v", err)
	}
	var l PaymentDemand
	for rows.Next() {
		if err = rows.Scan(&l.ID, &l.ImportDate, &l.IrisCode, &l.IrisName,
			&l.BeneficiaryID, &l.Beneficiary, &l.DemandNumber, &l.DemandDate,
			&l.ReceiptDate, &l.DemandValue, &l.CsfDate, &l.CsfComment, &l.DemandStatus,
			&l.StatusComment, &l.Excluded, &l.ExcludedComment, &l.ProcessedDate); err != nil {
			return fmt.Errorf("scan %v", err)
		}
		p.Lines = append(p.Lines, l)
	}
	err = rows.Err()
	if err != nil {
		return fmt.Errorf("rows err %v", err)
	}
	if len(p.Lines) == 0 {
		p.Lines = []PaymentDemand{}
	}
	return nil
}

// Validate checks if a payment batch has correct fields
func (p *PaymentDemandBatch) Validate() error {
	for i, l := range p.Lines {
		if l.IrisCode == "" {
			return fmt.Errorf("ligne %d IrisCode vide", i+1)
		}
		if l.IrisName == "" {
			return fmt.Errorf("ligne %d IrisName vide", i+1)
		}
		if int64(l.CommitmentDate) == 0 {
			return fmt.Errorf("ligne %d CommitmentDate vide", i+1)
		}
		if l.BeneficiaryCode == 0 {
			return fmt.Errorf("ligne %d BeneficiaryCode vide", i+1)
		}
		if l.DemandNumber == 0 {
			return fmt.Errorf("ligne %d DemandNumber vide", i+1)
		}
		if int64(l.DemandDate) == 0 {
			return fmt.Errorf("ligne %d DemandDate vide", i+1)
		}
		if int64(l.ReceiptDate) == 0 {
			return fmt.Errorf("ligne %d ReceiptDate vide", i+1)
		}
	}
	if p.ImportDate.IsZero() {
		return fmt.Errorf("date d'import non définie")
	}
	return nil
}

// excel2Time convert a int64 corresponding to an Excel integer date to time.Time
func excel2Time(i int64) time.Time {
	return time.Date(1899, 12, 30, 0, 0, 0, 0, time.UTC).
		Add(time.Duration(i*24) * time.Hour)
}

// Save import a batch of PaymentDemandLine and update the database accordingly.
// The batch must be valid i.e. the Validate() function should be called before
// using Save().
// The import process uses a temporary table to store the batch. This batch is
// first modified using a view to select the last beneficiary in case of
// duplicated lines due to the query the generates the batch. Only lines
// whose tuples of (iris_code,beneficiary_code,demand_number) are not already
// in the payment demands tables are added. The ImportDate field of the
// batch is used to fill the import_date of the newly inserted lines.
// For the existing lines, the csf_date, csf_comment,demand_status,status_comment
// and demand_value are updated.
// The null process_date are updated when the corresponding row in the database
// is missing in the batch.
func (p *PaymentDemandBatch) Save(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("tx begin %v", err)
	}

	if _, err := tx.Exec("DELETE from temp_payment_demands"); err != nil {
		tx.Rollback()
		return fmt.Errorf("delete %v", err)
	}

	stmt, err := tx.Prepare(pq.CopyIn("temp_payment_demands", "iris_code",
		"iris_name", "commitment_date", "beneficiary_code", "demand_number",
		"demand_date", "receipt_date", "demand_value", "csf_date", "csf_comment",
		"demand_status", "status_comment"))
	if err != nil {
		return fmt.Errorf("prepare stmt %v", err)
	}
	defer stmt.Close()
	for _, r := range p.Lines {
		if _, err = stmt.Exec(r.IrisCode, r.IrisName, excel2Time(r.CommitmentDate),
			r.BeneficiaryCode, r.DemandNumber, excel2Time(r.DemandDate),
			excel2Time(r.ReceiptDate), r.DemandValue, nullExcel2NullTime(r.CsfDate), r.CsfComment,
			r.DemandStatus, r.StatusComment); err != nil {
			tx.Rollback()
			return fmt.Errorf("insertion de %+v  %v", r, err)
		}
	}
	if _, err = stmt.Exec(); err != nil {
		tx.Rollback()
		return fmt.Errorf("statement exec flush %v", err)
	}
	type query struct {
		Query string
		Args  []interface{}
	}
	queries := []query{
		{Query: `INSERT INTO payment_demands (import_date,iris_code,
			iris_name,beneficiary_id,demand_number,demand_date,receipt_date,demand_value,
			csf_date,csf_comment,demand_status,status_comment,excluded,excluded_comment,
			processed_date)
		SELECT $1,t.iris_code,t.iris_name,b.id,t.demand_number,t.demand_date,
			t.receipt_date,t.demand_value,t.csf_date,t.csf_comment,t.demand_status,
			t.status_comment,FALSE,NULL::text,NULL::date
		FROM imported_payment_demands t
		JOIN beneficiary b ON b.code=t.beneficiary_code
		WHERE (t.iris_code,t.beneficiary_code,t.demand_number) NOT IN 
		(SELECT iris_code,beneficiary_code,demand_number FROM payment_demands)`,
			Args: []interface{}{p.ImportDate}},
		{
			Query: `UPDATE payment_demands SET csf_date=t.csf_date,csf_comment=t.csf_comment,
			demand_status=t.demand_status,status_comment=t.status_comment,
			demand_value=t.demand_value
			FROM (SELECT t.*,b.id AS beneficiary_id FROM imported_payment_demands t
				JOIN beneficiary b ON t.beneficiary_code=b.code) t
			WHERE (payment_demands.iris_code=t.iris_code AND
			payment_demands.beneficiary_id=t.beneficiary_id AND
			payment_demands.demand_number=t.demand_number)`,
			Args: []interface{}{}},
		{
			Query: `UPDATE payment_demands SET processed_date=$1
			WHERE (iris_code,beneficiary_id,demand_number) NOT IN 	
				(SELECT t.iris_code,b.id,t.demand_number FROM imported_payment_demands t
					JOIN beneficiary b ON t.beneficiary_code=b.code)
				AND processed_date IS NULL`,
			Args: []interface{}{p.ImportDate}},
		{
			Query: `DELETE from temp_payment_demands`,
			Args:  []interface{}{}},
	}
	for i, q := range queries {
		if _, err := tx.Exec(q.Query, q.Args...); err != nil {
			tx.Rollback()
			return fmt.Errorf("requête %d %v", i+1, err)
		}
	}
	return tx.Commit()
}

// GetAll fetches the count of the unprocessed or uncontrolled payment demands
// i.e. the number of the difference between the count of newly arrived demands
// and the count of controlled or processed demands for the 30 last days
func (p *PaymentDemandCounts) GetAll(db *sql.DB) error {
	rows, err := db.Query(`WITH t AS (SELECT CURRENT_DATE - generate_series(0,30) d ORDER BY 1),
  arrived AS (SELECT t.d,count(1) nb FROM payment_demands p, t
   WHERE p.receipt_date>t.d-30 AND p.receipt_date <= t.d GROUP BY 1),
  processed AS (SELECT t.d,count(1) nb FROM payment_demands p, t
   WHERE p.processed_date>t.d-30 AND p.processed_date <= t.d GROUP BY 1),
  controlled AS (SELECT t.d,count(1) nb FROM payment_demands p, t
   WHERE p.csf_date>t.d-30 AND p.csf_date <= t.d GROUP BY 1)
SELECT t.d, COALESCE(a.nb,0)-COALESCE(c.nb,0), COALESCE(a.nb,0)-COALESCE(p.nb,0)
  FROM t 
  LEFT JOIN arrived a ON a.d=t.d
  LEFT JOIN processed p ON p.d=t.d
	LEFT JOIN controlled c ON c.d=t.d
	ORDER BY 1`)
	if err != nil {
		return fmt.Errorf("select %v", err)
	}
	var l PaymentDemandCount
	for rows.Next() {
		if err = rows.Scan(&l.Date, &l.UnControlled, &l.UnProcessed); err != nil {
			return fmt.Errorf("scan %v", err)
		}
		p.Lines = append(p.Lines, l)
	}
	err = rows.Err()
	if err != nil {
		return fmt.Errorf("rows err %v", err)
	}
	if len(p.Lines) == 0 {
		p.Lines = []PaymentDemandCount{}
	}
	return nil
}
