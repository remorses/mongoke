package fakedata

import (
	"testing"

	"github.com/remorses/mongoke/src/testutil"
)

func TestGenerate(t *testing.T) {
	fakeData, err := NewFakeData(testutil.UserSchema)
	if err != nil {
		t.Error(err)
	}
	x, err := fakeData.Generate("User")
	if err != nil {
		t.Error(err)
	}
	print(testutil.Pretty(x))
}
