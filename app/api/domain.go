package api

type createdID struct {
	ID any `json:"id"`
}

func newCreatedID(id any) createdID {
	return createdID{
		ID: id,
	}
}
