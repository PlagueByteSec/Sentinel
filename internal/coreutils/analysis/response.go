package analysis

import (
	"io"
	"net/http"
	"strings"

	pools "github.com/PlagueByteSec/sdakit-project/v2/internal/datapools"
	"github.com/PlagueByteSec/sdakit-project/v2/internal/logging"
	"github.com/PlagueByteSec/sdakit-project/v2/internal/shared"
	"github.com/PlagueByteSec/sdakit-project/v2/pkg"
	"github.com/fhAnso/astkit"
)

func (check *SubdomainCheck) getResponse(url string) *http.Response {
	response := check.AnalysisSendRequest(AnalysisRequestConfig{Method: "GET", URL: url, Header: "", Value: ""})
	return pkg.Tern(response == nil, nil, response)
}

func (check *SubdomainCheck) responseGetBody(response *http.Response) []byte {
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		logging.GLogger.Log(err.Error())
		return nil
	}
	return responseBody
}

func (check *SubdomainCheck) checkPage(pageType string, pageInvestigate func(string, *http.Response) bool, successMessage string) {
	url := astkit.MakeUrl(astkit.HTTP(astkit.Basic), check.Subdomain)
	response := check.getResponse(url)
	if response == nil {
		return
	}
	if ok := pageInvestigate(url, response); ok {
		check.ConsoleOutput <- successMessage
		if pageType == "login" {
			pools.ManagePool(pools.PoolAction(pools.PoolAppend), check.Subdomain, &shared.GPoolBase.PoolLoginSubdomains)
		}
	}
}

func checkPageLogin(responseBody string) bool {
	if len(responseBody) != 0 {
		for idx := 0; idx < len(loginIndicators); idx++ {
			if strings.Contains(responseBody, loginIndicators[idx]) {
				return true
			}
		}
	}
	return false
}

func (check *SubdomainCheck) isLoginPage(url string, response *http.Response) bool {
	return checkPageLogin(string(check.responseGetBody(response)))
}
