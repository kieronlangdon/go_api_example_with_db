package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/go-resty/resty"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

//Maintestsuite a
type Maintestsuite struct {
	suite.Suite
	APIClient *resty.Client
}

//SetupTest a
func (suite *Maintestsuite) SetupTest() {
	suite.APIClient = resty.New()
}

//TestGetBooksStatusCodeShouldEqual200 a
func (suite *Maintestsuite) TestGetBooksStatusCodeShouldEqual200() {

	resp, _ := suite.APIClient.R().Get("http://localhost:8080/api/books")

	assert.Equal(suite.T(), 200, resp.StatusCode())
	assert.Equal(suite.T(), "application/json", resp.Header().Get("Content-Type"))
}

//TestGetBooksValuesShouldBeEqual
func (suite *Maintestsuite) TestGetBooksValuesShouldBeEqual() {

	resp, _ := suite.APIClient.R().
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetResult(&BookResponse{}).
		ForceContentType("application/json").
		Get("http://localhost:8001/api/books")

	var myResponse []BookResponse
	err := json.Unmarshal(resp.Body(), &myResponse)
	if err != nil {
		fmt.Println(err)
		return
	}
	assert.Equal(suite.T(), "1", myResponse[0].ID)
	assert.Equal(suite.T(), "Book one", myResponse[0].Title)

}

//TestMaintestsuite test method
func TestMaintestsuite(t *testing.T) {
	suite.Run(t, new(Maintestsuite))
}
