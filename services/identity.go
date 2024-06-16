package services

import (
	"database/sql"
	"errors"
	"log"
	"regexp"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
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

type IdenityDetails struct {
	Id          int      `json:"primaryContactId" db:"-"`
	Email       []string `json:"emails" db:"email"`
	PhoneNumber []string `json:"phoneNumbers" db:"phone_number"`
	SecondaryId []int    `json:"secondaryContactIds" db:"id"`
}

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

func generateResponse(pi *IdentityPrimary) (*IdenityDetails, error) {
	args := database.GetContactDetailsByIdArgs(pi.Id)
	row, err := db.Query(ctx, database.GetContactDetailsById, args)
	if err != nil {
		log.Printf("error fetching response details from db: %v\n", err)
		return nil, err
	}
	defer row.Close()

	details, err := pgx.CollectOneRow(row, pgx.RowToStructByName[IdenityDetails])
	if err != nil {
		log.Printf("error reding row details: %v\n", err)
		return nil, err
	}

	// adding primary details
	details.Id = pi.Id
	if pi.Email.Valid {
		primaryEmail := pi.Email.String
		details.Email = append(details.Email, "")
		copy(details.Email[1:], details.Email)

		details.Email[0] = primaryEmail
	}

	if pi.PhoneNumber.Valid {
		primaryPhoneNumber := pi.PhoneNumber.String
		details.PhoneNumber = append(details.PhoneNumber, "")
		copy(details.PhoneNumber[1:], details.PhoneNumber)

		details.PhoneNumber[0] = primaryPhoneNumber
	}
	return &details, nil
}

func (i *Identity) GetIdentity() (*IdenityDetails, error) {
	primaryContact, err := getPrimary(i)
	if err != nil {
		return nil, err
	}

	// create primary contact
	if primaryContact == nil {
		args := database.CreatePrimaryContactArgs(i.Email, i.PhoneNumber)
		_, err := db.Exec(ctx, database.CreatePrimaryContact, args)
		if err != nil {
			log.Printf("Error creating primary contact: %v\n", err)
			return nil, err
		}

		setValuesInRedis(i)
		// generate response
		res, err := generateResponse(primaryContact)
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	// create secondary contact if new value
	_, emailErr := redisClient.Get(ctx, i.Email).Result()
	_, phoneErr := redisClient.Get(ctx, i.PhoneNumber).Result()

	if emailErr != nil || phoneErr != nil {
		if errors.Is(emailErr, redis.Nil) || errors.Is(phoneErr, redis.Nil) {
			// one is new value
			args := database.CreateSecondaryContactArgs(primaryContact.Id, i.Email, i.PhoneNumber)
			_, err := db.Exec(ctx, database.CreateSecondaryContact, args)
			if err != nil {
				log.Printf("error creating secondary contact: %v\n", err)
				return nil, err
			}
			setValuesInRedis(i)
			// generate response
			res, err := generateResponse(primaryContact)
			if err != nil {
				return nil, err
			}
			return res, nil
		}

		log.Printf("error getting data from redis:\n emailErr: %v\nphoneErr:%v", emailErr, phoneErr)
		return nil, err
	}

	// generate response
	log.Println("primary id", primaryContact.Id)
	res, err := generateResponse(primaryContact)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func setValuesInRedis(i *Identity) {
	if i.Email != "" {
		err := redisClient.Set(ctx, i.Email, "1", 0).Err()
		if err != nil {
			log.Printf("error adding email in redis: %v\n", err)

		}

	}
	if i.PhoneNumber != "" {
		err := redisClient.Set(ctx, i.PhoneNumber, "1", 0).Err()
		if err != nil {
			log.Printf("error adding phoneNumber in redis: %v\n", err)

		}
	}
}

func getPrimary(i *Identity) (*IdentityPrimary, error) {
	email := i.Email
	phoneNumber := i.PhoneNumber

	args := database.GetPrimaryDetailsArgs(email, phoneNumber)
	row, err := db.Query(ctx, database.GetPrimaryDetailsBySecondary, args)
	if err != nil {
		log.Printf("Error querying db to get primary details by secondary: %v\n", err)
		return nil, err
	}
	defer row.Close()

	contact, secErr := pgx.CollectOneRow(row, pgx.RowToStructByName[IdentityPrimary])
	if secErr != nil && !errors.Is(secErr, pgx.ErrNoRows) {
		log.Printf("error collecting rows for get primary details by secondary: %v\n", err)
		return nil, secErr
	}

	gotPrimary := true
	if secErr != nil && errors.Is(secErr, pgx.ErrNoRows) {
		gotPrimary = false
	}

	// resolve primary
	args = database.GetPrimaryDetailsArgs(email, phoneNumber)
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

	if gotPrimary {
		for _, c := range contacts {
			if c.Id == contact.Id {
				continue
			}

			err := resolvePrimary(contact.Id, c.Id)
			if err != nil {
				return nil, err
			}
		}

		return &contact, nil
	}

	lenContact := len(contacts)

	if lenContact == 0 {
		return nil, nil
	}

	if lenContact == 1 {
		return &contacts[0], nil
	}

	// resolve primary
	record1 := contacts[0]
	record2 := contacts[1]

	if record1.Email.Valid && record1.Email.String == email {
		err := resolvePrimary(record1.Id, record2.Id)
		return &record1, err
	}

	err = resolvePrimary(record2.Id, record1.Id)
	return &record2, err
}

func resolvePrimary(primary, secondary int) error {
	args := database.ResolvePrimaryConflictArgs(primary, secondary)
	_, err := db.Exec(ctx, database.ResolvePrimaryConflict, args)
	if err != nil {
		log.Printf("Error resolving primary contact conflict: %v", err)
	}

	return err
}
