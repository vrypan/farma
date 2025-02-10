package fctools

import (
	"testing"
)

func Test_UserJson(t *testing.T) {
	u := NewUser().FromFid(280).FetchUserData(nil, nil)
	s, err := u.Json("", false, false)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("\n" + string(s))
}

func Test_UserString(t *testing.T) {
	u := NewUser().FromFname(nil, "vrypan").FetchUserData(nil, nil)
	t.Logf("\n%s", u)
}
