package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/go_assignment_3/dtos"
	"github.com/go_assignment_3/model"
)

type handler struct {
	db *sqlx.DB
	rdb *redis.Client
}


func (h *handler) createBook(c *gin.Context) {
	var input model.Book
	if err := c.BindJSON(&input); err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	query := fmt.Sprintf("insert into %s (author_id, title, date) values($1, $2, $3) returning id", BookTable)
	id, err := h.db.Exec(query, input.AuthorId, input.Title, input.Date)

	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return 
	}

	c.JSON(http.StatusOK, map[string]string{
		"message:": fmt.Sprintf("Added new book with id: %d", id),
	})

}

func (h *handler) getBookById(c *gin.Context) {
	id, err := ValidateId(c)
	
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	book, err := h.getBookFromRedis(c, id)

	if err == nil {
		c.JSON(http.StatusOK, book)
		return
	}

	book, err = h.getBookFromDB(id)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.cacheBook(c, id, book)
	c.JSON(http.StatusOK, book)
}

func (h *handler) getBookFromDB(id int) (dtos.Book, error){
	var book dtos.Book
	query := fmt.Sprintf(`select b.title, a.name, a.surname, b.rating, b.data from %s b
				join %s a on a.id = b.author_id
				where b.id = $1`, BookTable, AuthorTable)
	
	err := h.db.Get(&book, query, id)
	return book, err
}

func (h *handler) getBookFromRedis(c *gin.Context, id int) (dtos.Book, error){
	var book dtos.Book
	JSONbook, err := h.rdb.Get(c, strconv.Itoa(id)).Result()
	
	if err != nil {
		return dtos.Book{}, err
	}

	if err := json.Unmarshal([]byte(JSONbook), &book); err != nil {
		return dtos.Book{}, err
	}

	return book, nil
}

func (h *handler) cacheBook(c *gin.Context, id int, book dtos.Book) error {
	marshalledBook, err := json.Marshal(book)
	if err != nil {
		return err
	}
	if err := h.rdb.Set(c, strconv.Itoa(id), marshalledBook, time.Hour).Err(); err != nil {
		return err
	}
	return nil	
}

func (h *handler) RouterBuilder() *gin.Engine {
	router := gin.New()

	router.POST("/", h.createBook)
	router.GET("/:id", h.getBookById)

	return router
}

func HandlerBuilder(db *sqlx.DB, rdb *redis.Client) *handler{
	return &handler{
		db: db,
		rdb: rdb,
	}
}