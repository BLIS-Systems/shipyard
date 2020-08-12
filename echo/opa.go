package echo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo"
)

type opaResponse struct {
	DecisionID string    `json:"decision_id"`
	Result     opaResult `json:"result"`
}

type opaResult struct {
	Allow bool `json:"allow"`
}

type opaRequest struct {
	Input requestData `json:"input"`
}

type requestData struct {
	User   string `json:"user"`
	Path   string `json:"path"`
	Method string `json:"method"`
}

func (s *Service) authorize(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Prep the request
		oReq, err := json.Marshal(
			opaRequest{
				Input: requestData{
					User:   c.Request().Context().Value("user").(string),
					Path:   c.Path(),
					Method: c.Request().Method,
				},
			},
		)
		if err != nil {
			s.logger.Errorf("Error trying to create request to OPA: %s\n", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Error while contacting authorization server")
		}

		req, err := http.NewRequest(
			"POST",
			fmt.Sprintf("%s/v1/data/shipyard", s.opaAddress),
			bytes.NewReader(oReq),
		)
		if err != nil {
			s.logger.Errorf("Error trying to create request to OPA: %s\n", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Error while contacting authorization server")
		}

		req.Header.Set("authorization", fmt.Sprintf("Bearer %s", s.opaToken))

		// Make the request
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			s.logger.Errorf("Error while making request to OPA: %s", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Error while contacting authorization server")
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			s.logger.Errorf("Got back non 200 status from OPA: %d", res.StatusCode)
			return echo.NewHTTPError(http.StatusInternalServerError, "Error while contacting authorization server")
		}

		// Read the body
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			s.logger.Errorf("Unable to read body from OPA: %s", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Error while contacting authorization server")
		}

		// Unmarshal the body
		oRes := opaResponse{}
		err = json.Unmarshal(body, &oRes)
		if err != nil {
			s.logger.Errorf("Unable to parse body from OPA: %s", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Error while contacting authorization server")
		}

		// If OPA approved then allow the request, else reject with a 403
		if oRes.Result.Allow {
			return next(c)
		} else {
			return echo.NewHTTPError(http.StatusForbidden, "Unauthorized")
		}

	}
}
