package main

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/lunalubra/snippetbox/internal/assert"
)

func TestPing(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	res := ts.get(t, "/ping")
	assert.Equal(t, res.status, http.StatusOK)
	assert.Equal(t, res.body, "OK")
}

func TestSnippetView(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name       string
		urlPath    string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "Valid ID",
			urlPath:    "/snippet/view/1",
			wantStatus: http.StatusOK,
			wantBody:   "An old silent pond...",
		},
		{
			name:       "Non-existent ID",
			urlPath:    "/snippet/view/2",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "Negative ID",
			urlPath:    "/snippet/view/-1",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "Decimal ID",
			urlPath:    "/snippet/view/1.23",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "String ID",
			urlPath:    "/snippet/view/foo",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "Empty ID",
			urlPath:    "/wnippet/view/",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts.resetClientCookieJar(t)

			res := ts.get(t, tt.urlPath)
			assert.Equal(t, res.status, tt.wantStatus)
			assert.True(t, strings.Contains(res.body, tt.wantBody))
		})
	}
}

func TestSnippetViewExtendForm(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantForm bool
	}{
		{
			name:     "Snippet close to expiry",
			urlPath:  "/snippet/view/1",
			wantForm: true,
		},
		{
			name:     "Long lived snippet",
			urlPath:  "/snippet/view/3",
			wantForm: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts.resetClientCookieJar(t)

			res := ts.get(t, tt.urlPath)
			assert.Equal(t, res.status, http.StatusOK)
			assert.Equal(t, strings.Contains(res.body, "/snippet/extend/"), tt.wantForm)
		})
	}
}

func TestSnippetExtendPost(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name              string
		urlPath           string
		authenticate      bool
		useValidCSRFToken bool
		wantStatus        int
		wantLocation      string
	}{
		{
			name:              "Extendable snippet",
			urlPath:           "/snippet/extend/1",
			authenticate:      true,
			useValidCSRFToken: true,
			wantStatus:        http.StatusSeeOther,
			wantLocation:      "/snippet/view/1",
		},
		{
			name:              "Not extendable snippet",
			urlPath:           "/snippet/extend/3",
			authenticate:      true,
			useValidCSRFToken: true,
			wantStatus:        http.StatusSeeOther,
			wantLocation:      "/snippet/view/3",
		},
		{
			name:              "String ID",
			urlPath:           "/snippet/extend/foo",
			authenticate:      true,
			useValidCSRFToken: true,
			wantStatus:        http.StatusNotFound,
		},
		{
			name:              "Negative ID",
			urlPath:           "/snippet/extend/-1",
			authenticate:      true,
			useValidCSRFToken: true,
			wantStatus:        http.StatusNotFound,
		},
		{
			name:              "Invalid CSRF token",
			urlPath:           "/snippet/extend/1",
			authenticate:      true,
			useValidCSRFToken: false,
			wantStatus:        http.StatusBadRequest,
		},
		{
			name:              "Unauthenticated",
			urlPath:           "/snippet/extend/1",
			authenticate:      false,
			useValidCSRFToken: true,
			wantStatus:        http.StatusSeeOther,
			wantLocation:      "/user/login",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts.resetClientCookieJar(t)

			if tt.authenticate {
				ts.login(t)
			}

			res := ts.get(t, "/snippet/view/1")

			form := url.Values{}
			if tt.useValidCSRFToken {
				form.Add("csrf_token", extractCSRFToken(t, res.body))
			}

			res = ts.postForm(t, tt.urlPath, form)

			assert.Equal(t, res.status, tt.wantStatus)
			if tt.wantLocation != "" {
				assert.Equal(t, res.headers.Get("Location"), tt.wantLocation)
			}
		})
	}
}

func TestUserSignup(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	const (
		validName     = "Bob"
		validPassword = "validPa$$word"
		validEmail    = "bob@example.com"
		formTag       = "<form action='/user/signup' method='POST' novalidate>"
	)

	tests := []struct {
		name              string
		userName          string
		userEmail         string
		userPassword      string
		useValidCSRFToken bool
		wantStatus        int
		wantFormTag       string
	}{
		{
			name:              "Valid submission",
			userName:          validName,
			userEmail:         validEmail,
			userPassword:      validPassword,
			useValidCSRFToken: true,
			wantStatus:        http.StatusSeeOther,
		},
		{
			name:              "Invalid CSRF Token",
			userName:          validName,
			userEmail:         validEmail,
			userPassword:      validPassword,
			useValidCSRFToken: false,
			wantStatus:        http.StatusBadRequest,
		},
		{
			name:              "Empty name",
			userName:          "",
			userEmail:         validEmail,
			userPassword:      validPassword,
			useValidCSRFToken: true,
			wantStatus:        http.StatusUnprocessableEntity,
			wantFormTag:       formTag,
		},
		{
			name:              "Empty email",
			userName:          validName,
			userEmail:         "",
			userPassword:      validPassword,
			useValidCSRFToken: true,
			wantStatus:        http.StatusUnprocessableEntity,
			wantFormTag:       formTag,
		},
		{
			name:              "Empty password",
			userName:          validName,
			userEmail:         validEmail,
			userPassword:      "",
			useValidCSRFToken: true,
			wantStatus:        http.StatusUnprocessableEntity,
			wantFormTag:       formTag,
		},
		{
			name:              "Invalid email",
			userName:          validName,
			userEmail:         "bob@example.",
			userPassword:      validPassword,
			useValidCSRFToken: true,
			wantStatus:        http.StatusUnprocessableEntity,
			wantFormTag:       formTag,
		},
		{
			name:              "Short password",
			userName:          validName,
			userEmail:         validEmail,
			userPassword:      "pa$$",
			useValidCSRFToken: true,
			wantStatus:        http.StatusUnprocessableEntity,
			wantFormTag:       formTag,
		},
		{
			name:              "Duplicate email",
			userName:          validName,
			userEmail:         "dupe@example.com",
			userPassword:      validPassword,
			useValidCSRFToken: true,
			wantStatus:        http.StatusUnprocessableEntity,
			wantFormTag:       formTag,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts.resetClientCookieJar(t)

			res := ts.get(t, "/user/signup")

			form := url.Values{}
			form.Add("name", tt.userName)
			form.Add("email", tt.userEmail)
			form.Add("password", tt.userPassword)
			if tt.useValidCSRFToken {
				form.Add("csrf_token", extractCSRFToken(t, res.body))
			}

			res = ts.postForm(t, "/user/signup", form)

			assert.Equal(t, res.status, tt.wantStatus)
			assert.True(t, strings.Contains(res.body, tt.wantFormTag))
		})
	}
}
