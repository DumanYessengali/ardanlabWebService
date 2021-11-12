package tests

import (
	"bytes"
	"encoding/json"
	"github.com/DumanYessengali/ardanlabWebService/app/sales-api/handlers"
	"github.com/DumanYessengali/ardanlabWebService/business/auth"
	"github.com/DumanYessengali/ardanlabWebService/business/data/user"
	"github.com/DumanYessengali/ardanlabWebService/business/tests"
	"github.com/google/go-cmp/cmp"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

type UserTests struct {
	app        http.Handler
	kid        string
	userToken  string
	adminToken string
}

func TestUsers(t *testing.T) {
	test := tests.NewIntegration(t)
	t.Cleanup(test.Teardown)

	shutdown := make(chan os.Signal, 1)
	tests := UserTests{
		app:        handlers.API("develop", shutdown, test.Log, test.Auth, test.DB),
		kid:        test.KID,
		userToken:  test.Token(test.KID, "user@example.com", "gophers"),
		adminToken: test.Token(test.KID, "admin@example.com", "gophers"),
	}

	t.Run("crudUsers", tests.crudUser)
}

func (ut *UserTests) crudUser(t *testing.T) {
	nu := ut.postUser201(t)
	defer ut.deleteUser204(t, nu.ID)

	ut.getUser200(t, nu.ID)
	ut.putUser204(t, nu.ID)
	ut.putUser403(t, nu.ID)
}

func (ut *UserTests) postUser201(t *testing.T) user.Info {
	nu := user.NewUser{
		Name:            "Duman Yessengeli",
		Email:           "dumanyessengali@gmail.com",
		Roles:           []string{auth.RoleAdmin},
		Password:        "gophers",
		PasswordConfirm: "gophers",
	}

	body, err := json.Marshal(&nu)
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+ut.adminToken)
	ut.app.ServeHTTP(w, r)

	var got user.Info
	t.Log("Given the need to create a new user with the users endpoint.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using the declared user value.", testID)
		{
			if w.Code != http.StatusCreated {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 201 for the response : %v", tests.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 201 for the response.", tests.Success, testID)

			if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to unmarshal the response : %v", tests.Failed, testID, err)
			}

			exp := got
			exp.Name = "Duman Yessengeli"
			exp.Email = "dumanyessengali@gmail.com"
			exp.Roles = []string{auth.RoleAdmin}

			if diff := cmp.Diff(got, exp); diff != "" {
				t.Fatalf("\t%s\tTest %d:\tShould get the expected result. Diff:\n%s", tests.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould get the expected result.", tests.Success, testID)
		}
	}

	return got
}
func (ut *UserTests) putUser204(t *testing.T, id string) {
	body := `{"name": "Jacob Walker"}`

	r := httptest.NewRequest(http.MethodPut, "/users/"+id, strings.NewReader(body))
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+ut.adminToken)
	ut.app.ServeHTTP(w, r)

	t.Log("Given the need to update a user with the users endpoint.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using the modified user value.", testID)
		{
			if w.Code != http.StatusNoContent {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 204 for the response : %v", tests.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 204 for the response.", tests.Success, testID)

			r = httptest.NewRequest(http.MethodGet, "/users/"+id, nil)
			w = httptest.NewRecorder()

			r.Header.Set("Authorization", "Bearer "+ut.adminToken)
			ut.app.ServeHTTP(w, r)

			if w.Code != http.StatusOK {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 200 for the retrieve : %v", tests.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 200 for the retrieve.", tests.Success, testID)

			var ru user.Info
			if err := json.NewDecoder(w.Body).Decode(&ru); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to unmarshal the response : %v", tests.Failed, testID, err)
			}

			if ru.Name != "Jacob Walker" {
				t.Fatalf("\t%s\tTest %d:\tShould see an updated Name : got %q want %q", tests.Failed, testID, ru.Name, "Jacob Walker")
			}
			t.Logf("\t%s\tTest %d:\tShould see an updated Name.", tests.Success, testID)

			if ru.Email != "dumanyessengali@gmail.com" {
				t.Fatalf("\t%s\tTest %d:\tShould not affect other fields like Email : got %q want %q", tests.Failed, testID, ru.Email, "bill@ardanlabs.com")
			}
			t.Logf("\t%s\tTest %d:\tShould not affect other fields like Email.", tests.Success, testID)
		}
	}
}
func (ut *UserTests) deleteUser204(t *testing.T, id string) {
	r := httptest.NewRequest(http.MethodDelete, "/users/"+id, nil)
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+ut.adminToken)
	ut.app.ServeHTTP(w, r)

	t.Log("Given the need to validate deleting a user that does exist.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using the new user %s.", testID, id)
		{
			if w.Code != http.StatusNoContent {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 204 for the response : %v", tests.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 204 for the response.", tests.Success, testID)
		}
	}
}

func (ut *UserTests) getUser200(t *testing.T, id string) {
	r := httptest.NewRequest(http.MethodGet, "/users/"+id, nil)
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+ut.adminToken)
	ut.app.ServeHTTP(w, r)

	t.Log("Given the need to validate getting a user that exsits.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using the new user %s.", testID, id)
		{
			if w.Code != http.StatusOK {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 200 for the response : %v", tests.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 200 for the response.", tests.Success, testID)

			var got user.Info
			if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to unmarshal the response : %v", tests.Failed, testID, err)
			}

			exp := got
			exp.ID = id
			exp.Name = "Duman Yessengeli"
			exp.Email = "dumanyessengali@gmail.com"
			exp.Roles = []string{auth.RoleAdmin}

			if diff := cmp.Diff(got, exp); diff != "" {
				t.Fatalf("\t%s\tTest %d:\tShould get the expected result. Diff:\n%s", tests.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould get the expected result.", tests.Success, testID)
		}
	}
}
func (ut *UserTests) getUser403(t *testing.T) {
	t.Log("Given the need to validate regular users can't fetch other users.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen fetching the admin user as a regular data.", testID)
		{
			r := httptest.NewRequest(http.MethodGet, "/users/"+tests.AdminID, nil)
			w := httptest.NewRecorder()

			r.Header.Set("Authorization", "Bearer "+ut.userToken)
			ut.app.ServeHTTP(w, r)

			if w.Code != http.StatusForbidden {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 403 for the response : %v", tests.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 403 for the response.", tests.Success, testID)

			recv := w.Body.String()
			resp := `{"error":"attempted action is not allowed"}`
			if resp != recv {
				t.Log("Got :", recv)
				t.Log("Want:", resp)
				t.Fatalf("\t%s\tTest %d:\tShould get the expected result.", tests.Failed, testID)
			}
			t.Logf("\t%s\tTest %d:\tShould get the expected result.", tests.Success, testID)
		}

		testID = 1
		t.Logf("\tTest %d:\tWhen fetching the user as themselves.", testID)
		{

			r := httptest.NewRequest(http.MethodGet, "/v1/users/"+tests.UserID, nil)
			w := httptest.NewRecorder()

			r.Header.Set("Authorization", "Bearer "+ut.userToken)
			ut.app.ServeHTTP(w, r)

			if w.Code != http.StatusOK {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 200 for the response : %v", tests.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 200 for the response.", tests.Success, testID)
		}
	}
}
func (ut *UserTests) putUser403(t *testing.T, id string) {
	body := `{"name": "Anna Walker"}`

	r := httptest.NewRequest(http.MethodPut, "/users/"+id, strings.NewReader(body))
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+ut.userToken)
	ut.app.ServeHTTP(w, r)

	t.Log("Given the need to update a user with the users endpoint.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen a non-admin user makes a request", testID)
		{
			if w.Code != http.StatusForbidden {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 403 for the response : %v", tests.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 403 for the response.", tests.Success, testID)
		}
	}
}
