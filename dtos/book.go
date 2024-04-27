package dtos

type Book struct {
	Title         string  `json:"title" db:"title"`
	AuthorName    string  `json:"author_name" db:"author_name"`
	AuthorSurname string  `json:"author_surname" db:"author_surname"`
	Date          string  `json:"date" db:"date"`
	Rating        float32 `json:"rating" db:"rating"`
}