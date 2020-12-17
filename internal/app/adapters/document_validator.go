package adapters

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const validateServiceURL = "https://user-info.herokuapp.com/users/"
const ableString = "ABLE_TO_VOTE"

// DocValidator abstraction to encapsulate the document validation
type DocValidator struct{}

// ValidateDocument returns if document is able to vote
func (v DocValidator) ValidateDocument(document string) (bool, error) {
	res, err := http.Get(validateServiceURL + document)
	if err != nil {
		return false, err
	}

	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return false, err
	}

	j := map[string]interface{}{}
	json.Unmarshal(resBytes, &j)

	return j["status"] == ableString, nil
}
