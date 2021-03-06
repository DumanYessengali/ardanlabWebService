package user

import (
	"context"
	"database/sql"
	"github.com/DumanYessengali/ardanlabWebService/business/auth"
	"github.com/DumanYessengali/ardanlabWebService/business/errors"
	"github.com/DumanYessengali/ardanlabWebService/foundation/database"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	errs "github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

type User struct {
	log *log.Logger
	db  *sqlx.DB
}

func New(log *log.Logger, db *sqlx.DB) User {
	return User{
		log: log,
		db:  db,
	}
}

func (u User) Create(ctx context.Context, traceID string, nu NewUser, now time.Time) (Info, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		return Info{}, errs.Wrap(err, "generating password hash")
	}

	usr := Info{
		ID:           uuid.New().String(),
		Name:         nu.Name,
		Email:        nu.Email,
		Roles:        nu.Roles,
		PasswordHash: hash,
		DateCreated:  now.UTC(),
		DateUpdated:  now.UTC(),
	}

	const q = `INSERT INTO users (user_id, name, email, password_hash, roles, date_created, date_updated)
			   VAlUES ($1, $2, $3, $4, $5, $6, $7)`

	u.log.Printf("%s : %s query : %s", traceID, "user.Create",
		database.Log(q, usr.ID, usr.Name, usr.Email, usr.PasswordHash, usr.Roles, usr.DateCreated, usr.DateUpdated),
	)

	if _, err := u.db.ExecContext(ctx, q, usr.ID, usr.Name, usr.Email, usr.PasswordHash, usr.Roles, usr.DateCreated, usr.DateUpdated); err != nil {
		return Info{}, errs.Wrap(err, "inserting user")
	}
	return usr, nil
}

func (u User) Update(ctx context.Context, traceID string, claims auth.Claims, userID string, uu UpdateUser, now time.Time) error {
	usr, err := u.QueryByID(ctx, traceID, claims, userID)
	if err != nil {
		return err
	}

	if uu.Name != nil {
		usr.Name = *uu.Name
	}
	if uu.Email != nil {
		usr.Email = *uu.Email
	}
	if uu.Roles != nil {
		usr.Roles = uu.Roles
	}
	if uu.Password != nil {
		pw, err := bcrypt.GenerateFromPassword([]byte(*uu.Password), bcrypt.DefaultCost)
		if err != nil {
			return errs.Wrap(err, "generating password hash")
		}
		usr.PasswordHash = pw
	}
	usr.DateUpdated = now

	const q = `
	UPDATE
		users
	SET 
		"name" = $2,
		"email" = $3,
		"roles" = $4,
		"password_hash" = $5,
		"date_updated" = $6
	WHERE
		user_id = $1`

	u.log.Printf("%s: %s: %s", traceID, "user.Update",
		database.Log(q, usr.ID, usr.Name, usr.Email, usr.Roles, usr.PasswordHash, usr.DateCreated, usr.DateUpdated),
	)

	if _, err = u.db.ExecContext(ctx, q, userID, usr.Name, usr.Email, usr.Roles, usr.PasswordHash, usr.DateUpdated); err != nil {
		return errs.Wrap(err, "updating user")
	}

	return nil
}

func (u User) Delete(ctx context.Context, traceID string, userID string) error {
	if _, err := uuid.Parse(userID); err != nil {
		return errors.ErrInvalidID
	}
	const q = `DELETE FROM users where user_id = $1`

	u.log.Printf("%s : %s query : %s", traceID, "user.Delete",
		database.Log(q, userID),
	)

	if _, err := u.db.ExecContext(ctx, q, userID); err != nil {
		return errs.Wrapf(err, "deleting user %s", userID)
	}
	return nil
}

func (u User) Query(ctx context.Context, traceID string, pageNumber, rowsPerPage int) ([]Info, error) {
	const q = `SELECT *FROM users ORDER BY user_id OFFSET $1 ROWS FETCH NEXT $2 ROWS ONLY`
	offset := (pageNumber - 1) * rowsPerPage

	u.log.Printf("%s : %s query : %s", traceID, "user.Query", database.Log(q, offset, rowsPerPage))

	users := []Info{}

	if err := u.db.SelectContext(ctx, &users, q, offset, rowsPerPage); err != nil {
		return nil, errs.Wrap(err, "selecting users")
	}

	return users, nil
}

func (u User) QueryByID(ctx context.Context, traceID string, claims auth.Claims, userID string) (Info, error) {
	if _, err := uuid.Parse(userID); err != nil {
		return Info{}, errors.ErrInvalidID
	}
	if !claims.Authorized(auth.RoleAdmin) && claims.Subject != userID {
		return Info{}, errors.ErrForbidden
	}

	const q = `SELECT * FROM users WHERE user_id =$1`

	u.log.Printf("%s : %s query : %s", traceID, "user.QueryByID",
		database.Log(q, userID),
	)

	var usr Info
	if err := u.db.GetContext(ctx, &usr, q, userID); err != nil {
		if err == sql.ErrNoRows {
			return Info{}, errors.ErrNotFound
		}
		return Info{}, errs.Wrapf(err, "selecting user %q", userID)
	}

	return usr, nil
}

func (u User) QueryByEmail(ctx context.Context, traceID string, claims auth.Claims, email string) (Info, error) {
	const q = `SELECT * FROM users WHERE email =$1`

	u.log.Printf("%s : %s query : %s", traceID, "user.QueryByEmail",
		database.Log(q, email),
	)

	var usr Info
	if err := u.db.GetContext(ctx, &usr, q, email); err != nil {
		if err == sql.ErrNoRows {
			return Info{}, errors.ErrNotFound
		}
		return Info{}, errs.Wrapf(err, "selecting user %q", email)
	}

	if !claims.Authorized(auth.RoleAdmin) && claims.Subject != usr.ID {
		return Info{}, errors.ErrForbidden
	}

	return usr, nil
}

func (u User) Authenticate(ctx context.Context, traceID string, now time.Time, email, password string) (auth.Claims, error) {
	const q = `
	SELECT
		*
	FROM
		users
	WHERE
		email = $1`

	u.log.Printf("%s: %s: %s", traceID, "user.Authenticate",
		database.Log(q, email),
	)

	var usr Info
	if err := u.db.GetContext(ctx, &usr, q, email); err != nil {

		// Normally we would return ErrNotFound in this scenario but we do not want
		// to leak to an unauthenticated user which emails are in the system.
		if err == sql.ErrNoRows {
			return auth.Claims{}, errors.ErrAuthenticationFailure
		}

		return auth.Claims{}, errs.Wrap(err, "selecting single user")
	}

	// Compare the provided password with the saved hash. Use the bcrypt
	// comparison function so it is cryptographically secure.
	if err := bcrypt.CompareHashAndPassword(usr.PasswordHash, []byte(password)); err != nil {
		return auth.Claims{}, errors.ErrAuthenticationFailure
	}

	// If we are this far the request is valid. Create some claims for the user
	// and generate their token.
	claims := auth.Claims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    "service project",
			Subject:   usr.ID,
			Audience:  "students",
			ExpiresAt: now.Add(time.Hour).Unix(),
			IssuedAt:  now.Unix(),
		},
		Roles: usr.Roles,
	}

	return claims, nil
}
