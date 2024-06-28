package trace

import "github.com/google/uuid"

func ID() uuid.UUID {
	return uuid.New()
}

func String() string {
	return ID().String()
}
