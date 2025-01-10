package entities

type VectorType struct {
	ID        int
	Name      string
	Protected bool
}

type FreeVector struct {
	ID       int
	TypeID   int
	TypeName string
	PicLink  string
	Question string
	UserID   int64
	Answers  []VectorAnswer
}

type VectorAnswer struct {
	ID         int
	QuestionID int
	Answer     string
}

type VectorScore struct {
	UID        int64
	VectorType int
	Score      int
}
