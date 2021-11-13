package handlers

import (
	"github.com/DumanYessengali/ardanlabWebService/business/auth"
	"github.com/DumanYessengali/ardanlabWebService/business/data/book"
)

type bookGroup struct {
	book book.Book
	auth *auth.Auth
}
