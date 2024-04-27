package model
type Book struct {
	Id       int     `json:"id" db:"id"`
	Title    string  `json:"title" db:"title"`
	AuthorId int     `json:"author_id" db:"author_id"`
	Date     string  `json:"date" db:"date"`
	Rating   float32 `json:"rating" db:"rating"`
}

type Author struct {
	Name      string `json:"name" db:"name"`
	Surname   string `json:"surname" db:"surname"`
	BirthDate string `json:"birth_date" db:"birth_date"`
}