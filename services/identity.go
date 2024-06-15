package services

type Identity struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
}

func (i *Identity) GetIdentity() {

}
