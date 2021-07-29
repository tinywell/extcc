// Copyright the Hyperledger Fabric contributors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package cmds

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

// Post ...
func Post(url string, body []byte) ([]byte, error) {

	buf := bytes.NewReader(body)
	res, err := http.Post(url, "application/json", buf)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return resBody, nil
}
