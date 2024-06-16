package database

import "github.com/jackc/pgx/v5"

const GetPrimaryDetails = `
	SELECT id, email, phone_number, linked_id
	FROM contact 
	WHERE (email = @email OR phone_number = @phoneNumber) 
	AND link_precedence = 'primary'
`

func GetPrimaryDetailsArgs(email, phoneNumber string) pgx.NamedArgs {
	return pgx.NamedArgs{
		"email":       email,
		"phoneNumber": phoneNumber,
	}
}

// create primary contact
const CreatePrimaryContact = `
	INSERT INTO contact (email, phone_number, link_precedence) 
	VALUES (COALESCE(@email, NULL), COALESCE(@phoneNumber, NULL), 'primary')
	RETURNING id
`

func CreatePrimaryContactArgs(email, phoneNumber string) pgx.NamedArgs {
	args := pgx.NamedArgs{}
	if email != "" {
		args["email"] = email
	}
	if phoneNumber != "" {
		args["phoneNumber"] = phoneNumber
	}

	return args
}
