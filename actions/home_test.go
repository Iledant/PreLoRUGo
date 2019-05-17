package actions

import (
	"net/http"
	"strings"
	"testing"
)

// testHome is the entry point for testing all home requests
func testHome(t *testing.T, c *TestContext) {
	t.Run("Home", func(t *testing.T) {
		testGetHome(t, c)
	})
}

// testGetHome checks if route is user protected and datas correctly sent back
func testGetHome(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: "",
			RespContains: []string{`Token absent`},
			Count:        1,
			StatusCode:   http.StatusInternalServerError}, // 0 : token empty
		{Token: c.Config.Users.User.Token,
			RespContains: []string{`"Commitment":`, `"Payment":`},
			Count:        0,
			StatusCode:   http.StatusOK}, // 1 : ok
	}
	for i, tc := range tcc {
		response := c.E.GET("/api/home").
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
		body := string(response.Content)
		for _, r := range tc.RespContains {
			if !strings.Contains(body, r) {
				t.Errorf("GetHomes[%d]\n  ->attendu %s\n  ->reçu: %s", i, r, body)
			}
		}
		status := response.Raw().StatusCode
		if status != tc.StatusCode {
			t.Errorf("GetHomes[%d]  ->status attendu %d  ->reçu: %d", i, tc.StatusCode, status)
		}
		if status == http.StatusOK {
			count := strings.Count(body, `"Month"`)
			if count != tc.Count {
				t.Errorf("GetHomes[%d]  ->nombre attendu %d  ->reçu: %d", i, tc.Count, count)
			}
		}
	}
}
