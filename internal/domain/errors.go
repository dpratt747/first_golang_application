package domain

type UniqueConstraintDatabaseError struct {
	Message string
}

func (ucDE *UniqueConstraintDatabaseError) Error() string {
	return ucDE.Message
}

type UnmappedDatabaseError struct {
	Message string
}

func (ucDE *UnmappedDatabaseError) Error() string {
	return ucDE.Message
}

type DatabaseTransactionError struct {
	Message string
}

func (ucDE *DatabaseTransactionError) Error() string {
	return ucDE.Message
}

type UserNotFoundError struct {
	Message string
}

func (ucDE *UserNotFoundError) Error() string {
	return ucDE.Message
}
