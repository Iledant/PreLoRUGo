package actions

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
)

// testBeneficiaryDatas is the entry point for testing all renew projet requests
func testBeneficiaryDatas(t *testing.T, c *TestContext) {
	t.Run("BeneficiaryDatas", func(t *testing.T) {
		testGetBeneficiaryDatas(t, c)
	})
}

// testGetBeneficiaryDatas checks if route is user protected and datas correctly sent back
func testGetBeneficiaryDatas(t *testing.T, c *TestContext) {
	tcc := []TestCase{
		{Token: "",
			Sent:         []byte(`Page=2&Year=2010&Search=savigny`),
			RespContains: []string{`Token absent`},
			Count:        1,
			ID:           3,
			StatusCode:   http.StatusInternalServerError}, // 0 : token empty
		{Token: c.Config.Users.User.Token,
			Sent:         []byte(`Page=2&Year=a&Search=savigny`),
			RespContains: []string{`Page de données bénéficiaire, décodage Year :`},
			Count:        1,
			ID:           3,
			StatusCode:   http.StatusInternalServerError}, // 1 : bad param query
		{Token: c.Config.Users.User.Token,
			Sent:         []byte(`Page=2&Year=2010&Search=savigny`),
			RespContains: []string{`"Datas":[],"Page":1,"ItemsCount":0`},
			Count:        0,
			ID:           0,
			StatusCode:   http.StatusOK}, // 2 : bad ID
		{Token: c.Config.Users.User.Token,
			Sent: []byte(`Page=2&Year=2010&Search=savigny`),
			//cSpell: disable
			RespContains: []string{`"Datas":[{"ID":3,"Date":"2015-04-13T00:00:00Z","Value":30000000,"Name":"91 - SAVIGNY SUR ORGE - AV DE LONGJUMEAU - 65 PLUS/PLAI","Available":30000000}],"Page":1,"ItemsCount":1`},
			//cSpell: enable
			ID:         3,
			Count:      1,
			StatusCode: http.StatusOK}, // 3 : ok
	}
	for i, tc := range tcc {
		response := c.E.GET("/api/beneficiary/"+strconv.Itoa(tc.ID)+"/datas").WithQueryString(string(tc.Sent)).
			WithHeader("Authorization", "Bearer "+tc.Token).Expect()
		body := string(response.Content)
		for _, r := range tc.RespContains {
			if !strings.Contains(body, r) {
				t.Errorf("GetBeneficiaryDatas[%d]\n  ->attendu %s\n  ->reçu: %s", i, r, body)
			}
		}
		status := response.Raw().StatusCode
		if status != tc.StatusCode {
			t.Errorf("GetBeneficiaryDatas[%d]  ->status attendu %d  ->reçu: %d", i, tc.StatusCode, status)
		}
		if status == http.StatusOK {
			count := strings.Count(body, `"ID"`)
			if count != tc.Count {
				t.Errorf("GetBeneficiaryDatas[%d]  ->nombre attendu %d  ->reçu: %d", i, tc.Count, count)
			}
		}
	}
}
