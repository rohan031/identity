package services

import (
	"database/sql"
	"log"
	"regexp"

	"github.com/jackc/pgx/v5"
	"github.com/rohan031/identity/database"
)

type Identity struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
}

type IdentityPrimary struct {
	Id          int            `db:"id"`
	Email       sql.NullString `db:"email"`
	PhoneNumber sql.NullString `db:"phone_number"`
	LinkedId    sql.NullInt32  `db:"linked_id"`
}

type IdentityId struct {
	Id int `db:"id"`
}

// func (i *Identity) GetIdentity() error {
// 	primaryId, _, _, err := getPrimary(i.Email, i.PhoneNumber)
// 	if err != nil {
// 		return err
// 	}

// 	if primaryId == -1 {
// 		id, err := createPrimary(i.Email, i.PhoneNumber)
// 		if err != nil {
// 			return err
// 		}

// 		log.Println(id)
// 		return nil
// 	}

// 	return nil
// }

// func createPrimary(email, phoneNumber string) (int, error) {
// 	args := database.CreatePrimaryContactArgs(email, phoneNumber)
// 	rows, err := db.Query(ctx, database.CreatePrimaryContact, args)
// 	if err != nil {
// 		log.Printf("Error creating primary contact: %v\n", err)
// 		return -1, err
// 	}
// 	defer rows.Close()

// 	contact, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[IdentityId])
// 	if err != nil {
// 		log.Printf("Error reading rows create primary: %v\n", err)
// 		return -1, err
// 	}

// 	return contact.Id, nil
// }

// func getPrimary(email, phoneNumber string) (int, bool, *Identity, error) {
// 	args := database.GetPrimaryDetailsArgs(email, phoneNumber)
// 	rows, err := db.Query(ctx, database.GetPrimaryDetails, args)
// 	if err != nil {
// 		log.Printf("Error getting data from db: %v\n", err)
// 		return 0, false, &Identity{}, err
// 	}
// 	defer rows.Close()

// 	contacts, err := pgx.CollectRows(rows, pgx.RowToStructByName[IdentityPrimary])
// 	if err != nil {
// 		log.Printf("Error reading rows get primary: %v", err)
// 		return 0, false, &Identity{}, err
// 	}

// 	contactsLength := len(contacts)

// 	if contactsLength == 0 {
// 		return -1, false, &Identity{}, nil
// 	}

// 	if contactsLength == 1 {
// 		email := ""
// 		if contacts[0].Email.Valid {
// 			email = contacts[0].Email.String
// 		}

// 		phoneNumber := ""
// 		if contacts[0].PhoneNumber.Valid {
// 			phoneNumber = contacts[0].PhoneNumber.String
// 		}

// 		return contacts[0].Id, true, &Identity{
// 			Email:       email,
// 			PhoneNumber: phoneNumber,
// 		}, nil
// 	}

// 	// resolve primary

// 	return 0, false, &Identity{
// 		Email:       email,
// 		PhoneNumber: phoneNumber,
// 	}, nil

// }

func validEmail(email string) bool {
	regexEmail := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	emailValid, _ := regexp.MatchString(regexEmail, email)

	return emailValid
}

func validPhone(phoneNumber string) bool {
	regexPhone := `^[0-9]+$`
	phoneValid, _ := regexp.MatchString(regexPhone, phoneNumber)

	return phoneValid
}

func (i *Identity) ValidateBody() bool {
	return !(i.Email == "" && i.PhoneNumber == "") &&
		(i.Email == "" || validEmail(i.Email)) &&
		(i.PhoneNumber == "" || validPhone(i.PhoneNumber))

}

func (i *Identity) GetIdentity() error {
	primaryContact, err := getPrimary(i)
	if err != nil {
		return err
	}

	// create primary contact
	if primaryContact == nil {
		args := database.CreatePrimaryContactArgs(i.Email, i.PhoneNumber)
		_, err := db.Exec(ctx, database.CreatePrimaryContact, args)
		if err != nil {
			log.Printf("Error creating primary contact: %v\n", err)
		}
		return err
	}

	log.Println("primary id", primaryContact.Id)
	return nil
}

func getPrimary(i *Identity) (*IdentityPrimary, error) {
	email := i.Email
	phoneNumber := i.PhoneNumber

	args := database.GetPrimaryDetailsArgs(email, phoneNumber, "primary")
	rows, err := db.Query(ctx, database.GetPrimaryDetails, args)
	if err != nil {
		log.Printf("Error querying db to get primary details: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	contacts, err := pgx.CollectRows(rows, pgx.RowToStructByName[IdentityPrimary])
	if err != nil {
		log.Printf("error collecting rows for get primary details: %v\n", err)
		return nil, err
	}

	lenContact := len(contacts)

	if lenContact == 0 {
		args := database.GetPrimaryDetailsArgs(email, phoneNumber, "secondary")
		rows, err := db.Query(ctx, database.GetPrimaryDetails, args)
		if err != nil {
			log.Printf("Error querying db to get primary details: %v\n", err)
			return nil, err
		}
		defer rows.Close()

		contacts, err := pgx.CollectRows(rows, pgx.RowToStructByName[IdentityPrimary])
		if err != nil {
			log.Printf("error collecting rows for get primary details: %v\n", err)
			return nil, err
		}

		if len(contacts) == 0 {
			return nil, nil
		}

		if !contacts[0].LinkedId.Valid {
			return nil, nil
		}

		args = database.GetContactByIdArgs(contacts[0].LinkedId.Int32)
		row, err := db.Query(ctx, database.GetContactById, args)
		if err != nil {
			log.Printf("Error querying db to get primary details by id: %v\n", err)
			return nil, err
		}
		defer row.Close()

		contact, err := pgx.CollectOneRow(row, pgx.RowToStructByName[IdentityPrimary])
		if err != nil {
			log.Printf("error collecting rows for get primary details by id: %v\n", err)
			return nil, err
		}

		return &contact, nil
	}

	if lenContact == 1 {
		return &contacts[0], nil
	}

	// resolve primary
	record1 := contacts[0]
	record2 := contacts[1]

	if record1.Email.Valid && record1.Email.String == email {
		err := resolvePrimary(record1, record2)
		return &record1, err
	}

	err = resolvePrimary(record2, record1)
	return &record2, err
}

func resolvePrimary(primary, secondary IdentityPrimary) error {
	args := database.ResolvePrimaryConflictArgs(primary.Id, secondary.Id)
	_, err := db.Exec(ctx, database.ResolvePrimaryConflict, args)
	if err != nil {
		log.Printf("Error resolving primary contact conflict: %v", err)
	}

	return err
}
