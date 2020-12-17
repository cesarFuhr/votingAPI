package adapters

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

const validadeServiceURL = "https://user-info.herokuapp.com/users/"
const ableString = "ABLE_TO_VOTE"

type jsonValidateRes struct {
	status string
}

// DocValidator abstraction to encapsulate the document validation
type DocValidator struct{}

// ValidateDocument returns if document is able to vote
func (v DocValidator) ValidateDocument(document string) (bool, error) {
	res, err := http.Get(validadeServiceURL + document)
	if err != nil {
		return false, err
	}

	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return false, err
	}

	var j jsonValidateRes
	json.Unmarshal(resBytes, &j)

	return strings.Contains(j.status, ableString), nil
}
