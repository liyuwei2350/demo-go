package main

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

type Book struct {
	ID     int64  `json:"id"`
	Title  string `json:"title" binding:"required"`
	Author string `json:"author" binding:"required"`
}

type bookStore struct {
	mu     sync.RWMutex
	nextID int64
	books  map[int64]Book
}

type createBookRequest struct {
	Title  string `json:"title" binding:"required"`
	Author string `json:"author" binding:"required"`
}

type updateBookRequest struct {
	Title  string `json:"title" binding:"required"`
	Author string `json:"author" binding:"required"`
}

func newBookStore() *bookStore {
	return &bookStore{
		nextID: 1,
		books:  make(map[int64]Book),
	}
}

func registerBookRoutes(r *gin.Engine, store *bookStore) {
	books := r.Group("/books")
	{
		books.POST("", store.createBook)
		books.GET("", store.listBooks)
		books.GET("/:id", store.getBook)
		books.PUT("/:id", store.updateBook)
		books.DELETE("/:id", store.deleteBook)
	}
}

func (s *bookStore) createBook(c *gin.Context) {
	var req createBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	book := Book{
		ID:     s.nextID,
		Title:  req.Title,
		Author: req.Author,
	}
	s.books[book.ID] = book
	s.nextID++

	c.JSON(http.StatusCreated, book)
}

func (s *bookStore) listBooks(c *gin.Context) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	books := make([]Book, 0, len(s.books))
	for _, book := range s.books {
		books = append(books, book)
	}

	c.JSON(http.StatusOK, books)
}

func (s *bookStore) getBook(c *gin.Context) {
	id, ok := parseBookID(c)
	if !ok {
		return
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	book, exists := s.books[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	c.JSON(http.StatusOK, book)
}

func (s *bookStore) updateBook(c *gin.Context) {
	id, ok := parseBookID(c)
	if !ok {
		return
	}

	var req updateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	book, exists := s.books[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	book.Title = req.Title
	book.Author = req.Author
	s.books[id] = book

	c.JSON(http.StatusOK, book)
}

func (s *bookStore) deleteBook(c *gin.Context) {
	id, ok := parseBookID(c)
	if !ok {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.books[id]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	delete(s.books, id)
	c.Status(http.StatusNoContent)
}

func parseBookID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book id"})
		return 0, false
	}

	return id, true
}
