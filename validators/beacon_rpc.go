package validators

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type BeaconApi struct {
	Network  string `json:"network" yaml:"network"`
	Endpoint string `json:"endpoint" yaml:"endpoint"`
}

func (rpc BeaconApi) GetValidatorData(validator Validator) (data ValidatorData, err error) {
	var validatorID string
	if validator.Index != "" {
		validatorID = validator.Index
	} else {
		validatorID = validator.Pubkey
	}

	if validatorID == "" {
		err = fmt.Errorf("validator must define index or pubkey")
		return
	}

	uri := strings.Join([]string{rpc.Endpoint, "eth/v1/beacon/states/head/validators", validatorID}, "/")

	resp, err := http.Get(uri)
	if err != nil {
		return
	}

	rawData, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if json.Unmarshal(rawData, &data); err != nil {
		return
	}

	return
}
