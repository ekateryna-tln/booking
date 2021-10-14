package forms

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("post", "/whatever", nil)
	form := New(r.PostForm)
	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_CheckRequiredFields(t *testing.T) {
	r := httptest.NewRequest("post", "/whatever", nil)
	form := New(r.PostForm)
	form.CheckRequiredFields("a", "b", "c")
	isValid := form.Valid()
	if isValid {
		t.Error("form shows valid when required fields missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "b")
	postedData.Add("c", "c")

	form = New(postedData)
	form.CheckRequiredFields("a", "b", "c")
	isValid = form.Valid()
	if !isValid {
		t.Error("form shows does not have required fields when it does")
	}
}

func TestForm_IsEmail(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)
	form.IsEmail("email")
	isValid := form.Valid()
	if isValid {
		t.Error("form shows valid email for non-existent field")
	}

	postedData = url.Values{}
	postedData.Add("email", "invalid_email")
	form = New(postedData)
	form.IsEmail("email")
	isValid = form.Valid()
	if isValid {
		t.Error("form shows valid email when it does not")
	}

	postedData = url.Values{}
	postedData.Add("email", "valid_email@test.test")
	form = New(postedData)
	form.IsEmail("email")
	isValid = form.Valid()
	if !isValid {
		t.Error("form shows invalid email when it does")
	}
}

func TestForm_HasField(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)
	has := form.HasField("non_existed_field")
	if has {
		t.Error("form shows has field when it does not")
	}

	postedData.Add("existed_field", "existed_field_value")
	form = New(postedData)
	has = form.HasField("existed_field")
	if !has {
		t.Error("form shows does not have field when it should")
	}
}

func TestForm_MinLength(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)
	form.MinLength("some_field", 100)
	isValid := form.Valid()
	if isValid {
		t.Error("form shows minlength of non-existent field")
	}

	isError := form.Errors.Get("some_field")
	if isError == "" {
		t.Error("should have an error but did not get one")
	}

	postedData.Add("some_field", "some_value")
	form = New(postedData)
	form.MinLength("some_field", 100)
	isValid = form.Valid()
	if isValid {
		t.Error("form shows minlength of 100 met when data is shorter")
	}

	postedData.Add("some_other_field", "some_other_value")
	form = New(postedData)
	form.MinLength("some_other_field", 3)
	isValid = form.Valid()
	if !isValid {
		t.Error("form shows minlength of 3 is not met when it is")
	}

	isError = form.Errors.Get("some_other_field")
	if isError != "" {
		t.Error("should not have an error but got one")
	}
}
