package database

import "github.com/jackc/pgx/v5"

const GetPrimaryDetails = `
	SELECT id, email, phone_number, linked_id
	FROM contact 
	WHERE (email = @email OR phone_number = @phoneNumber) 
	AND link_precedence = @link
`

func GetPrimaryDetailsArgs(email, phoneNumber, link string) pgx.NamedArgs {
	return pgx.NamedArgs{
		"email":       email,
		"phoneNumber": phoneNumber,
		"link":        link,
	}
}

// get contact by id
const GetContactById = `
	SELECT id, email, phone_number, linked_id
	FROM contact 
	WHERE id=@id
`

func GetContactByIdArgs(id int32) pgx.NamedArgs {
	return pgx.NamedArgs{
		"id": id,
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
