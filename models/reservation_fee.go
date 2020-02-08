package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
)

// ReservationFee model
type ReservationFee struct {
	ID                   int64       `json:"ID"`
	CurrentBeneficiaryID int64       `json:"CurrentBeneficiaryID"`
	CurrentBeneficiary   string      `json:"CurrentBeneficiary"`
	FirstBeneficiaryID   NullInt64   `json:"FirstBeneficiaryID"`
	FirstBeneficiary     NullString  `json:"FirstBeneficiary"`
	CityCode             int64       `json:"CityCode"`
	City                 string      `json:"City"`
	AddressNumber        NullString  `json:"AddressNumber"`
	AddressStreet        NullString  `json:"AddressStreet"`
	RPLS                 NullString  `json:"RPLS"`
	Convention           NullString  `json:"Convention"`
	ConventionTypeID     NullInt64   `json:"ConventionTypeID"`
	ConventionType       NullString  `json:"ConventionType"`
	Count                int64       `json:"Count"`
	TransferDate         NullTime    `json:"TransferDate"`
	CommentID            NullInt64   `json:"CommentID"`
	Comment              NullString  `json:"Comment"`
	TransferID           NullInt64   `json:"TransferID"`
	Transfer             NullString  `json:"Transfer"`
	PMR                  bool        `json:"PMR"`
	ConventionDate       NullTime    `json:"ConventionDate"`
	EliseRef             NullString  `json:"EliseRef"`
	Area                 NullFloat64 `json:"Area"`
	EndYear              NullInt64   `json:"EndYear"`
	Loan                 NullFloat64 `json:"Loan"`
	Charges              NullFloat64 `json:"Charges"`
}

// ReservationFees embeddes an array of ReservationFee for json export and dedicated
// queries
type ReservationFees struct {
	Lines []ReservationFee `json:"ReservationFee"`
}

// ReservationFeeLine is used to decode one line of the batch for the upsert
// query of the reservation fee
type ReservationFeeLine struct {
	CurrentBeneficiary string      `json:"CurrentBeneficiary"`
	FirstBeneficiary   NullString  `json:"FirstBeneficiary"`
	City               string      `json:"City"`
	AddressNumber      NullString  `json:"AddressNumber"`
	AddressStreet      string      `json:"AddressStreet"`
	Convention         NullString  `json:"Convention"`
	Typology           NullString  `json:"Typology"`
	RPLS               NullString  `json:"RPLS"`
	ConventionType     NullString  `json:"ConventionType"`
	Count              int         `json:"Count"`
	Transfer           NullString  `json:"Transfer"`
	TransferDate       NullInt64   `json:"TransferDate"`
	PMR                bool        `json:"PMR"`
	Comment            NullString  `json:"Comment"`
	ConventionDate     NullInt64   `json:"ConventionDate"`
	Area               NullFloat64 `json:"Area"`
	EndYear            NullInt64   `json:"EndYear"`
	Loan               NullFloat64 `json:"Loan"`
	Charges            NullInt64   `json:"Charges"`
}

// ReservationFeeBatch embeddes an array of ReservationFeeLines for the import
// batch query
type ReservationFeeBatch struct {
	Lines []ReservationFeeLine `json:"ReservationFee"`
}

// PaginatedReservationFees embeddes an array of ReservationFees for json export
// with paginated informations
type PaginatedReservationFees struct {
	ReservationFees []ReservationFee `json:"ReservationFee"`
	Page            int64            `json:"Page"`
	ItemsCount      int64            `json:"ItemsCount"`
}

// Valid checks if fields complies with database constraints
func (r *ReservationFee) Valid() error {
	if r.CurrentBeneficiaryID == 0 {
		return fmt.Errorf("CurrentBeneficiaryID null")
	}
	if r.CityCode == 0 {
		return fmt.Errorf("CityCode null")
	}
	if r.Count == 0 {
		return fmt.Errorf("Count null")
	}
	return nil
}

// getOuterFields fetches fields that belong to other tables that reservation_fee
// in order to have a common part between create and update functions
func (r *ReservationFee) getOuterFields(db *sql.DB) error {
	return db.QueryRow(`SELECT b1.name,b2.name,hc.name,ht.name,ho.name
	FROM reservation_fee rf 
	JOIN beneficiary b1 ON rf.current_beneficiary_id=b1.id
	LEFT OUTER JOIN beneficiary b2 ON first_beneficiary_id=b2.id
	LEFT OUTER JOIN housing_convention hc ON hc.id=rf.convention_type_id
	LEFT OUTER JOIN housing_transfer ht ON ht.id=rf.transfer_id
	LEFT OUTER JOIN housing_comment ho ON ho.id=rf.comment_id
	WHERE rf.id=$1`, r.ID).Scan(&r.CurrentBeneficiary, &r.FirstBeneficiary,
		&r.ConventionType, &r.Transfer, &r.Comment)
}

// Create insert a new ReservationFee into database
func (r *ReservationFee) Create(db *sql.DB) error {
	if err := db.QueryRow(`INSERT into reservation_fee (current_beneficiary_id,
		first_beneficiary_id,city_code,address_number,address_street,rpls,
		convention,convention_type_id,count,transfer_date,transfer_id,pmr,comment_id,
		convention_date,elise_ref,area,end_year,loan,charges)
		VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19)
		RETURNING ID`, r.CurrentBeneficiaryID, r.FirstBeneficiaryID, r.CityCode,
		r.AddressNumber, r.AddressStreet, r.RPLS, r.Convention, r.ConventionTypeID, r.Count,
		r.TransferDate, r.TransferID, r.PMR, r.CommentID, r.ConventionDate, r.EliseRef,
		r.Area, r.EndYear, r.Loan, r.Charges).Scan(&r.ID); err != nil {
		return fmt.Errorf("insert %v", err)
	}
	if err := r.getOuterFields(db); err != nil {
		return fmt.Errorf("select %v", err)
	}
	return nil
}

// Update changes a reservation fee in the database using the fields
func (r *ReservationFee) Update(db *sql.DB) error {
	res, err := db.Exec(`UPDATE reservation_fee SET current_beneficiary_id=$1,
		first_beneficiary_id=$2,city_code=$3,address_number=$4,address_street=$5,
		rpls=$6,convention=$7,convention_type_id=$8,count=$9,transfer_date=$10,
		transfer_id=$11,comment_id=$12,convention_date=$13,elise_ref=$14,area=$15,
		end_year=$16,loan=$17,charges=$18,pmr=$19 WHERE ID=$20`, r.CurrentBeneficiaryID,
		r.FirstBeneficiaryID, r.CityCode, r.AddressNumber, r.AddressStreet, r.RPLS,
		r.Convention, r.ConventionTypeID, r.Count, r.TransferDate, r.TransferID,
		r.CommentID, r.ConventionDate, r.EliseRef, r.Area, r.EndYear, r.Loan,
		r.Charges, r.PMR, r.ID)
	if err != nil {
		return fmt.Errorf("update %v", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("count %v", err)
	}
	if count != 1 {
		return fmt.Errorf("réservation non trouvée")
	}
	if err := r.getOuterFields(db); err != nil {
		return fmt.Errorf("select %v", err)
	}
	return nil
}

// Delete removes a reservation fee from database
func (r *ReservationFee) Delete(db *sql.DB) error {
	res, err := db.Exec(`DELETE FROM reservation_fee WHERE id=$1`, r.ID)
	if err != nil {
		return fmt.Errorf("delete %v", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("count %v", err)
	}
	if count != 1 {
		return fmt.Errorf("réservation non trouvée")
	}
	return nil
}

// Save import a batch of reservation fee, updating the housing transfer, housing
// convention, housing typology, housing comment and convention type tables
func (r *ReservationFeeBatch) Save(db *sql.DB) error {
	for i, l := range r.Lines {
		if l.CurrentBeneficiary == "" {
			return fmt.Errorf("line %d, CurrentBeneficiary empty", i+1)
		}
		if l.City == "" {
			return fmt.Errorf("line %d, City empty", i+1)
		}
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(pq.CopyIn("temp_reservation_fee", "current_beneficiary",
		"first_beneficiary", "city", "address_number", "address_street", "convention",
		"typology", "rpls", "convention_type", "count", "transfer", "transfer_date",
		"pmr", "comment", "convention_date", "area", "end_year", "loan", "charges"))
	if err != nil {
		return err
	}
	defer stmt.Close()
	var (
		transferDate, conventionDate NullTime
		b                            = time.Date(1899, 12, 30, 0, 0, 0, 0, time.UTC)
	)
	for _, r := range r.Lines {
		transferDate.Valid = r.TransferDate.Valid
		conventionDate.Valid = r.ConventionDate.Valid
		if r.TransferDate.Valid {
			transferDate.Time = b.Add(time.Duration(r.TransferDate.Int64*24) * time.Hour)
		}
		if r.ConventionDate.Valid {
			conventionDate.Time = b.Add(time.Duration(r.ConventionDate.Int64*24) * time.Hour)
		}
		if _, err = stmt.Exec(r.CurrentBeneficiary, r.FirstBeneficiary, r.City,
			r.AddressNumber, r.AddressStreet, r.Convention, r.Typology, r.RPLS,
			r.ConventionType, r.Count, r.Transfer, transferDate, r.PMR, r.Comment,
			conventionDate, r.Area, r.EndYear, r.Loan, r.Charges); err != nil {
			tx.Rollback()
			return err
		}
	}
	if _, err = stmt.Exec(); err != nil {
		tx.Rollback()
		return fmt.Errorf("statement flush exec %v", err)
	}
	queries := []string{
		`INSERT INTO housing_typology(name)
			SELECT DISTINCT typology FROM temp_reservation_fee WHERE typology NOTNULL
			ON CONFLICT DO NOTHING`, // 0
		`INSERT INTO housing_transfer(name) 
			SELECT DISTINCT transfer FROM temp_reservation_fee WHERE transfer NOTNULL
			ON CONFLICT DO NOTHING`, // 1
		`INSERT INTO convention_type(name)
			SELECT DISTINCT convention_type FROM temp_reservation_fee 
			WHERE convention_type NOTNULL
			ON CONFLICT DO NOTHING`, // 2
		`INSERT INTO housing_comment(name)
			SELECT DISTINCT comment FROM temp_reservation_fee WHERE comment NOTNULL
			ON CONFLICT DO NOTHING`, // 3
		`INSERT INTO reservation_fee (current_beneficiary_id, first_beneficiary_id,
				city_code,address_number,address_street,rpls,convention,convention_type_id,
				count,transfer_date,transfer_id,pmr,comment_id,convention_date,elise_ref,
				area,end_year,loan,charges)
			SELECT b1.id,b2.id,c.insee_code,rf.address_number,rf.address_street,rf.rpls,
				rf.convention,ct.id,rf.count,rf.transfer_date,ht.id,rf.pmr,hc.id,
				rf.convention_date,NULL,rf.area,rf.end_year,rf.loan,rf.charges
			FROM temp_reservation_fee rf
			JOIN beneficiary b1 ON b1.name=rf.current_beneficiary
			LEFT JOIN beneficiary b2 ON b2.name=rf.first_beneficiary
			JOIN city c ON rf.city=c.name
			LEFT JOIN convention_type ct ON ct.name=rf.convention_type
			LEFT JOIN housing_transfer ht ON ht.name=rf.transfer
			LEFT JOIN housing_comment hc ON hc.name=rf.comment`,
		`DELETE FROM temp_reservation_fee`, // 4
	}
	for i, q := range queries {
		_, err = tx.Exec(q)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("requête %d : %s", i, err.Error())
		}
	}
	tx.Commit()
	return nil
}

// Get fetches all paginated reservation fees from database that match the
// paginated query
func (p *PaginatedReservationFees) Get(db *sql.DB, q *PaginatedQuery) error {
	var count int64
	if err := db.QueryRow(`SELECT count(1) FROM reservation_fee rf
	JOIN beneficiary b1 ON b1.id=rf.current_beneficiary_id
	LEFT JOIN beneficiary b2 ON b2.id=rf.first_beneficiary_id
	JOIN city c ON rf.city_code=c.insee_code
	LEFT JOIN convention_type ct ON rf.convention_type_id=ct.id
	LEFT JOIN housing_comment cmt ON cmt.id=rf.comment_id
	LEFT JOIN housing_transfer ht ON ht.id=rf.transfer_id
	WHERE (b1.name ILIKE $1 OR b2.name ILIKE $1 OR c.name ILIKE $1 OR 
		rf.convention ILIKE $1 OR cmt.name ILIKE $1 OR rf.address_street ILIKE $1 OR
		rf.address_number ILIKE $1 OR ht.name ILIKE $1 OR rf.elise_ref ILIKE $1)`,
		"%"+q.Search+"%").Scan(&count); err != nil {
		return fmt.Errorf("count query failed %v", err)
	}
	offset, newPage := GetPaginateParams(q.Page, count)

	rows, err := db.Query(`SELECT rf.id,rf.current_beneficiary_id,b1.name,
		rf.first_beneficiary_id,b2.name,rf.city_code,c.name,rf.address_number,
		rf.address_street,rf.rpls,rf.convention,rf.convention_type_id,ct.name,
		rf.count,rf.transfer_date,rf.comment_id,cmt.name,rf.transfer_id,ht.name,
		rf.pmr,rf.convention_date,rf.elise_ref,rf.area,rf.end_year,rf.loan,rf.charges
	FROM reservation_fee rf
	JOIN beneficiary b1 ON b1.id=rf.current_beneficiary_id
	LEFT JOIN beneficiary b2 ON b2.id=rf.first_beneficiary_id
	JOIN city c ON rf.city_code=c.insee_code
	LEFT JOIN convention_type ct ON rf.convention_type_id=ct.id
	LEFT JOIN housing_comment cmt ON cmt.id=rf.comment_id
	LEFT JOIN housing_transfer ht ON ht.id=rf.transfer_id
	WHERE (b1.name ILIKE $1 OR b2.name ILIKE $1 OR c.name ILIKE $1 OR 
		rf.convention ILIKE $1 OR cmt.name ILIKE $1 OR rf.address_street ILIKE $1 OR
		rf.address_number ILIKE $1 OR ht.name ILIKE $1 OR rf.elise_ref ILIKE $1)
	ORDER BY 1 LIMIT $2 OFFSET $3`, "%"+q.Search+"%", PageSize, offset)
	if err != nil {
		return err
	}
	var row ReservationFee
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&row.ID, &row.CurrentBeneficiaryID, &row.CurrentBeneficiary,
			&row.FirstBeneficiaryID, &row.FirstBeneficiary, &row.CityCode, &row.City,
			&row.AddressNumber, &row.AddressStreet, &row.RPLS, &row.Convention,
			&row.ConventionTypeID, &row.ConventionType, &row.Count, &row.TransferDate,
			&row.CommentID, &row.Comment, &row.TransferID, &row.Transfer, &row.PMR,
			&row.ConventionDate, &row.EliseRef, &row.Area, &row.EndYear, &row.Loan,
			&row.Charges); err != nil {
			return err
		}
		p.ReservationFees = append(p.ReservationFees, row)
	}
	err = rows.Err()
	if len(p.ReservationFees) == 0 {
		p.ReservationFees = []ReservationFee{}
	}
	p.Page = newPage
	p.ItemsCount = count
	return err
}
