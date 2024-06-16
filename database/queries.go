package database

import "github.com/jackc/pgx/v5"

const GetPrimaryDetails = `
	SELECT id, email, phone_number, linked_id
	FROM contact 
	WHERE (email = @email OR phone_number = @phoneNumber) 
	AND link_precedence = 'primary'
`

func GetPrimaryDetailsArgs(email, phoneNumber string) pgx.NamedArgs {
	args := pgx.NamedArgs{}
	if email != "" {
		args["email"] = email
	}
	if phoneNumber != "" {
		args["phoneNumber"] = phoneNumber
	}

	return args
}

// get primary Contact from secondary
const GetPrimaryDetailsBySecondary = `
WITH linkedId AS (
	SELECT linked_id AS id FROM contact
	WHERE (email = @email OR phone_number = @phoneNumber) 
	AND link_precedence = 'secondary'
)
SELECT c.id, c.email, c.phone_number, c.linked_id FROM contact c, linkedId l
WHERE c.id = l.id
`

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

// resolving primary conflict
const ResolvePrimaryConflict = `
UPDATE contact
SET link_precedence = 'secondary', linked_id = @pId
WHERE id = @sId
`

func ResolvePrimaryConflictArgs(pId, sId int) pgx.NamedArgs {
	return pgx.NamedArgs{
		"pId": pId,
		"sId": sId,
	}
}
