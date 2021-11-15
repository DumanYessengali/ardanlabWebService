package handlers

import (
	"context"
	"fmt"
	"github.com/DumanYessengali/ardanlabWebService/business/auth"
	"github.com/DumanYessengali/ardanlabWebService/business/data/book_info"
	"github.com/DumanYessengali/ardanlabWebService/business/errors"
	"github.com/DumanYessengali/ardanlabWebService/foundation/web"
	errs "github.com/pkg/errors"
	"net/http"
	"strconv"
)

type bookInfoGroup struct {
	bookInfo book_info.BookInfo
	auth     *auth.Auth
}

func (bg bookInfoGroup) query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	params := web.Params(r)
	pageNumber, err := strconv.Atoi(params["page"])
	if err != nil {
		return web.NewRequestError(fmt.Errorf("invalid page format: %s", params["page"]), http.StatusBadRequest)
	}

	rowsPerPage, err := strconv.Atoi(params["rows"])
	if err != nil {
		return web.NewRequestError(fmt.Errorf("invalid rows format: %s", params["rows"]), http.StatusBadRequest)
	}

	books, err := bg.bookInfo.Query(ctx, v.TraceID, pageNumber, rowsPerPage)
	if err != nil {
		return errs.Wrap(err, "unable to query for books")
	}

	return web.Respond(ctx, w, books, http.StatusOK)
}

func (bg bookInfoGroup) queryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errs.New("claims missing from context")
	}

	params := web.Params(r)
	usr, err := bg.bookInfo.QueryByID(ctx, v.TraceID, claims, params["id"])

	if err != nil {
		if err != nil {
			switch err {
			case errors.ErrInvalidID:
				return web.NewRequestError(err, http.StatusBadRequest)
			case errors.ErrNotFound:
				return web.NewRequestError(err, http.StatusNotFound)
			case errors.ErrForbidden:
				return web.NewRequestError(err, http.StatusForbidden)
			default:
				return errs.Wrapf(err, "ID: %s", params["id"])
			}
		}
	}

	return web.Respond(ctx, w, usr, http.StatusOK)
}

func (bg bookInfoGroup) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	var nb book_info.NewBookInfo
	if err := web.Decode(r, &nb); err != nil {
		return errs.Wrapf(err, "unable to decode payload")
	}

	usr, err := bg.bookInfo.Create(ctx, v.TraceID, nb, v.Now)
	if err != nil {
		return errs.Wrapf(err, "Book: %+v", &usr)
	}

	return web.Respond(ctx, w, usr, http.StatusCreated)
}

func (bg bookInfoGroup) update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errs.New("claims missing from context")
	}

	var bpd book_info.UpdateBookInfo
	if err := web.Decode(r, &bpd); err != nil {
		return errs.Wrap(err, "unable to decode payload")
	}

	params := web.Params(r)
	err := bg.bookInfo.Update(ctx, v.TraceID, claims, params["id"], bpd, v.Now)
	if err != nil {
		switch err {
		case errors.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case errors.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case errors.ErrForbidden:
			return web.NewRequestError(err, http.StatusForbidden)
		default:
			return errs.Wrapf(err, "ID: %s Book: %+v", params["id"], &bpd)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

func (bg bookInfoGroup) delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	params := web.Params(r)
	err := bg.bookInfo.Delete(ctx, v.TraceID, params["id"])

	if err != nil {
		if err != nil {
			switch err {
			case errors.ErrInvalidID:
				return web.NewRequestError(err, http.StatusBadRequest)
			case errors.ErrNotFound:
				return web.NewRequestError(err, http.StatusNotFound)
			case errors.ErrForbidden:
				return web.NewRequestError(err, http.StatusForbidden)
			default:
				return errs.Wrapf(err, "ID: %s", params["id"])
			}
		}
	}
	return web.Respond(ctx, w, nil, http.StatusNoContent)
}
