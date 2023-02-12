package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "a")
	postedData.Add("c", "a")

	r, _ = http.NewRequest("POST", "/whatever", nil)

}

func TestForm_Has(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)
	has := form.Has("whatever")
	if has {
		t.Error("form shows has field when it does not")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")

	form = New(postedData)

	has = form.Has("a")

	if !has {
		t.Error("shows form does not have field when it should")
	}
}

func TestForm_MinLength(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.MinLength("x", 10)
	if form.Valid() {
		t.Error("form shows min length for non-existent field")
	}

	isError := form.Errors.Get("x")
	if isError == "" {
		t.Error("should have an error but did not get one")
	}

	postedValues := url.Values{}
	postedValues.Add("some_field", "some value")
	form = New(postedValues)

	form.MinLength("some_field", 100)
	if form.Valid() {
		t.Error("shows minlength of 100 met when data is shorter")
	}

	postedValues = url.Values{}
	postedValues.Add("another_field", "abc123")
	form = New(postedValues)

	form.MinLength("another_field", 1)
	if !form.Valid() {
		t.Error("shows minlength of 1 is not met when it is")
	}
	isError = form.Errors.Get("another_field")
	if isError != "" {
		t.Error("should not have an error but got one")
	}

	// hasMinLength := form.MinLength("whatever", 3, r)
	// if hasMinLength {
	// 	t.Error("no field but still minlength satisfied")

	// }
	// postedData := url.Values{}
	// postedData.Add("a", "a")

	// form = New(postedData)

	// hasMinLength = form.MinLength("a", 1, r)

	// if !hasMinLength {
	// 	t.Error("had minlength satisfied still showed false")
	// }

}

func TestForm_IsEmail(t *testing.T) {
	postedValues := url.Values{}
	form := New(postedValues)
	form.IsEmail("x")
	if form.Valid() {
		t.Error("form shows valid email when it's non existent field")
	}
	postedValues = url.Values{}
	postedValues.Add("email", "x")
	form = New(postedValues)
	form.IsEmail("email")
	if form.Valid() {
		t.Error("shows valid form even though email is invalid")
	}

	postedValues = url.Values{}
	postedValues.Add("email", "a@gmail.com")
	form = New(postedValues)
	form.IsEmail("email")
	if !form.Valid() {
		t.Error("shows invalid form although email is valid")
	}

}
