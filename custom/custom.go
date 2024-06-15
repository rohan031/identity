package custom

type MalformedRequest struct {
	Status  int
	Message string
}

func (mr *MalformedRequest) Error() string {
	return mr.Message
}
