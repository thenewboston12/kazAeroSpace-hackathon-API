package main

import (
	"database/sql"
	"errors"
	"fmt"
)

var DELETED_STATUS = "DELETED_STATUS"
var ACTIVE_STATUS = "ACTIVE_STATUS"

var ErroFailedToInsert = errors.New("failed to insert record")
var ErroFailedToDelete = errors.New("failed to delete record")
var ErroFailedToUpdate = errors.New("failed to update record")
var ErroFailedToGetAll = errors.New("failed to get all records")
var ErroFailedToGetById = errors.New("failed to get by record by id")
var ErrorRecordNotExists = errors.New("invalid id. record doesn't exists")

type Record struct {
	Id         int64   `json:"id"`
	Name       string  `json:"name"`
	Surname    string  `json:"surname"`
	Middlename string  `json:"middlename"`
	Iin        string  `json:"iin"`
	Dom        int     `json:"dom"`
	Kv         int     `json:"kv"`
	City       string  `json:"city"`
	Street     string  `json:"street"`
	CadastrNum string  `json:"cadastr_num"`
	Area       float64 `json:"area"`
	Latitude   float64 `json:"lat"`
	Longitude  float64 `json:"long"`
	Comment    string  `json:"comment"`
	Status     string  `json:"status"`
}

type RecordsRepository struct {
	Db *sql.DB
}

func NewRecordsRepo(Db *sql.DB) *RecordsRepository {
	return &RecordsRepository{
		Db: Db,
	}
}

// here we will write sql commands and make queries to db
func (repo *RecordsRepository) CreateRecord(r *Record) (*Record, error) {
	sqlStatement := `
		INSERT INTO records (name, surname, middle_name, iin, dom, kv, city, street,
			cadastr_num, area, lat, long, comment
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING  id, name, surname, middle_name, iin, dom, kv, city, street, cadastr_num, area, lat, long, comment, status
	`

	// execute statement to create entry in db
	args := []interface{}{r.Name, r.Surname, r.Middlename, r.Iin, r.Dom, r.Kv, r.City, r.Street, r.CadastrNum, r.Area, r.Latitude, r.Longitude, r.Comment}

	var rec Record

	err := repo.Db.QueryRow(sqlStatement, args...).Scan(&rec.Id, &rec.Name, &rec.Surname, &rec.Middlename, &rec.Iin, &rec.Dom, &rec.Kv, &rec.City, &rec.Street, &rec.CadastrNum, &rec.Area, &rec.Latitude, &rec.Longitude, &rec.Comment, &rec.Status)

	if err != nil {
		fmt.Println(err)
		return nil, ErroFailedToInsert
	}

	return &rec, nil
}

//TODO
func (repo *RecordsRepository) UpdateRecord(r *Record) (*Record, error) {
	sqlStatement := `
		UPDATE records
		SET name = $1, surname=$2, middle_name=$3, iin=$4, dom=$5, kv=$6, city=$7, street=$8, cadastr_num=$9, area=$10, lat=$11, long=$12, comment=$13
		WHERE id=$14
		RETURNING id, name, surname, middle_name, iin, dom, kv, city, street, cadastr_num, area, lat, long, comment, status
	`
	args := []interface{}{r.Name, r.Surname, r.Middlename, r.Iin, r.Dom, r.Kv, r.City, r.Street, r.CadastrNum, r.Area, r.Latitude, r.Longitude, r.Comment, r.Id}

	var u Record
	err := repo.Db.QueryRow(sqlStatement, args...).Scan(&u.Id, &u.Name, &u.Surname, &u.Middlename, &u.Iin, &u.Dom, &u.Kv, &u.City, &u.Street, &u.CadastrNum, &u.Area, &u.Latitude, &u.Longitude, &u.Comment, &u.Status)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, ErrorRecordNotExists
		}
		return nil, err
	}
	return &u, nil
}

func (repo *RecordsRepository) DeleteRecord(id int64) (*Record, error) {
	sqlStatement := `
		UPDATE records 
		SET status = $1
		WHERE id = $2
		RETURNING id, name, surname, middle_name, iin, dom, kv, city, street, cadastr_num, area, lat, long, comment, status
	`

	// update record in db
	args := []interface{}{DELETED_STATUS, id}

	var u Record
	err := repo.Db.QueryRow(sqlStatement, args...).Scan(&u.Id, &u.Name, &u.Surname, &u.Middlename, &u.Iin, &u.Dom, &u.Kv, &u.City, &u.Street, &u.CadastrNum, &u.Area, &u.Latitude, &u.Longitude, &u.Comment, &u.Status)

	if err != nil {
		fmt.Println(err)
		return nil, ErroFailedToDelete
	}

	return &u, nil
}

//TODO

func (repo *RecordsRepository) GetAllRecords() ([]*Record, error) {
	rows, err := repo.Db.Query("SELECT * FROM records WHERE status = $1", ACTIVE_STATUS)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*Record

	for rows.Next() {
		var rec Record
		if err := rows.Scan(&rec.Id, &rec.Name, &rec.Surname, &rec.Middlename, &rec.Iin, &rec.Dom, &rec.Kv, &rec.City, &rec.Street, &rec.CadastrNum, &rec.Area, &rec.Latitude, &rec.Longitude, &rec.Comment, &rec.Status); err != nil {
			return records, ErroFailedToGetAll
		}
		records = append(records, &rec)
	}

	if err = rows.Err(); err != nil {
		return records, ErroFailedToGetAll
	}
	return records, nil

}

func (repo *RecordsRepository) GetRecordById(id int64) (*Record, error) {
	var rec Record
	err := repo.Db.QueryRow("SELECT * FROM records WHERE id=$1", id).Scan(&rec.Id, &rec.Name, &rec.Surname, &rec.Middlename, &rec.Iin, &rec.Dom, &rec.Kv, &rec.City, &rec.Street, &rec.CadastrNum, &rec.Area, &rec.Latitude, &rec.Longitude, &rec.Comment, &rec.Status)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, ErrorRecordNotExists
		}
		return nil, err
	}
	return &rec, nil
}
