package domain

type UniqueConstraintDatabaseError struct {
	Message string
}

func(ucDE *UniqueConstraintDatabaseError) Error() string {
	return ucDE.Message
}


type UnmappedDatabaseError struct {
	Message string
}

func(ucDE *UnmappedDatabaseError) Error() string {
	return ucDE.Message
}