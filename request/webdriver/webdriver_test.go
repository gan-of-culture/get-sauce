package webdriver

import (
	"testing"

	"github.com/gan-of-culture/get-sauce/test"
)

func TestSolveChallenge(t *testing.T) {
	t.Run("Default test", func(t *testing.T) {

		wd, err := New()
		test.CheckError(t, err)

		cookies, err := wd.SolveChallenge("https://hentaibar.com/")
		test.CheckError(t, err)

		if cookies == nil {
			t.Errorf("Got: %v - Want: %v", cookies, "cookies")
		}
	})
}
