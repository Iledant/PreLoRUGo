package models

import "database/sql"

// BeneficiaryPayment is used to decode one line of the query dedicated to
// payment per month calculus of a beneficiary
type BeneficiaryPayment struct {
	Year  int64   `json:"Year"`
	Month int64   `json:"Month"`
	Value float64 `json:"Value"`
}

// BeneficiaryPayments embeddes an array of BeneficiaryPayment for json export
// to encapsulate all datas of the payment per month query
type BeneficiaryPayments struct {
	Lines []BeneficiaryPayment `json:"BeneficiaryPayment"`
}

// GetAll fetches the payments per month and year for a given beneficiary
func (b *BeneficiaryPayments) GetAll(ID int64, db *sql.DB) error {
	rows, err := db.Query(`select p.year,m,0.01*sum(p.value)::double precision
  from generate_series(1,12) m
  join payment p on m=extract(month from p.creation_date)
  join commitment c on p.commitment_id=c.id   
  where c.beneficiary_id=$1
  group by 1,2 order by 1,2;`, ID)
	if err != nil {
		return err
	}
	var r BeneficiaryPayment
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&r.Year, &r.Month, &r.Value); err != nil {
			return err
		}
		b.Lines = append(b.Lines, r)
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	if len(b.Lines) == 0 {
		b.Lines = []BeneficiaryPayment{}
	}
	return nil
}
