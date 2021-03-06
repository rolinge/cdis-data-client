package tests

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/uc-cdis/gen3-client/gen3-client/jwt"
	"github.com/uc-cdis/gen3-client/gen3-client/mocks"
)

func Requesting(*http.Response) *http.Response {
	return &http.Response{}
}

func TestDoRequestWithSignedHeaderNoProfile(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockConfig := mocks.NewMockConfigureInterface(mockCtrl)
	testFunction := &jwt.Functions{Config: mockConfig}

	cred := jwt.Credential{KeyId: "", APIKey: "", AccessKey: "", APIEndpoint: ""}

	mockConfig.EXPECT().ParseConfig(gomock.Any()).Return(cred).Times(1)

	_, err := testFunction.DoRequestWithSignedHeader("default", "not_json", "/user/data/download/test_uuid", "", nil)

	if err == nil {
		t.Fail()
	}
}

func TestDoRequestWithSignedHeaderGoodToken(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockConfig := mocks.NewMockConfigureInterface(mockCtrl)
	mockRequest := mocks.NewMockRequestInterface(mockCtrl)
	testFunction := &jwt.Functions{Config: mockConfig, Request: mockRequest}

	cred := jwt.Credential{KeyId: "", APIKey: "fake_api_key", AccessKey: "non_exprired_token", APIEndpoint: "http://www.test.com"}
	mockedResp := &http.Response{
		Body:       ioutil.NopCloser(bytes.NewBufferString("{\"url\": \"http://www.test.com/user/data/download/test_uuid\"}")),
		StatusCode: 200,
	}

	mockConfig.EXPECT().ParseConfig("default").Return(cred).Times(1)
	mockRequest.EXPECT().MakeARequest("GET", "http://www.test.com/user/data/download/test_uuid", "non_exprired_token", "", gomock.Any(), gomock.Any()).Return(mockedResp, nil).Times(1)

	_, err := testFunction.DoRequestWithSignedHeader("default", "", "/user/data/download/test_uuid", "", nil)

	if err != nil {
		t.Fail()
	}
}

func TestDoRequestWithSignedHeaderCreateNewToken(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockConfig := mocks.NewMockConfigureInterface(mockCtrl)
	mockRequest := mocks.NewMockRequestInterface(mockCtrl)
	testFunction := &jwt.Functions{Config: mockConfig, Request: mockRequest}

	cred := jwt.Credential{KeyId: "", APIKey: "fake_api_key", AccessKey: "", APIEndpoint: "http://www.test.com"}
	mockedResp := &http.Response{
		Body:       ioutil.NopCloser(bytes.NewBufferString("{\"url\": \"www.test.com/user/data/download/\"}")),
		StatusCode: 200,
	}

	mockConfig.EXPECT().ParseConfig("default").Return(cred).Times(1)
	mockConfig.EXPECT().ReadFile(gomock.Any(), gomock.Any()).Times(1)
	mockConfig.EXPECT().UpdateConfigFile(cred, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1)

	mockRequest.EXPECT().RequestNewAccessKey("http://www.test.com/user/credentials/api/access_token", &cred).Return(nil).Times(1)
	mockRequest.EXPECT().MakeARequest(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockedResp, nil).Times(1)

	_, err := testFunction.DoRequestWithSignedHeader("default", "", "/user/data/download/test_uuid", "", nil)

	if err != nil {
		t.Fail()
	}
}

func TestDoRequestWithSignedHeaderRefreshToken(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockConfig := mocks.NewMockConfigureInterface(mockCtrl)
	mockRequest := mocks.NewMockRequestInterface(mockCtrl)
	testFunction := &jwt.Functions{Config: mockConfig, Request: mockRequest}

	cred := jwt.Credential{KeyId: "", APIKey: "fake_api_key", AccessKey: "expired_token", APIEndpoint: "http://www.test.com"}
	mockedResp := &http.Response{
		Body:       ioutil.NopCloser(bytes.NewBufferString("{\"url\": \"www.test.com/user/data/download/\"}")),
		StatusCode: 401,
	}

	mockConfig.EXPECT().ParseConfig("default").Return(cred).Times(1)
	mockConfig.EXPECT().ReadFile(gomock.Any(), gomock.Any()).Times(1)
	mockConfig.EXPECT().UpdateConfigFile(cred, gomock.Any(), "http://www.test.com", gomock.Any(), gomock.Any(), gomock.Any(), "default").Times(1)

	mockRequest.EXPECT().RequestNewAccessKey("http://www.test.com/user/credentials/api/access_token", &cred).Return(nil).Times(1)
	mockRequest.EXPECT().MakeARequest(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockedResp, nil).Times(2)

	_, err := testFunction.DoRequestWithSignedHeader("default", "", "/user/data/download/test_uuid", "", nil)

	if err != nil && !strings.Contains(err.Error(), "401") {
		t.Fail()
	}

}

func TestCheckPrivilegesNoProfile(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockConfig := mocks.NewMockConfigureInterface(mockCtrl)
	testFunction := &jwt.Functions{Config: mockConfig}

	cred := jwt.Credential{KeyId: "", APIKey: "", AccessKey: "", APIEndpoint: ""}

	mockConfig.EXPECT().ParseConfig(gomock.Any()).Return(cred).Times(1)

	_, _, err := testFunction.CheckPrivileges("default", "")

	if err == nil {
		t.Errorf("Expected an error on missing credentials in configuration, but not received")
	}
}

func TestCheckPrivilegesNoAccess(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockConfig := mocks.NewMockConfigureInterface(mockCtrl)
	mockRequest := mocks.NewMockRequestInterface(mockCtrl)
	testFunction := &jwt.Functions{Config: mockConfig, Request: mockRequest}

	cred := jwt.Credential{KeyId: "", APIKey: "fake_api_key", AccessKey: "non_exprired_token", APIEndpoint: "http://www.test.com"}
	mockedResp := &http.Response{
		Body:       ioutil.NopCloser(bytes.NewBufferString("{\"project_access\": {}}")),
		StatusCode: 200,
	}

	mockConfig.EXPECT().ParseConfig("default").Return(cred).Times(1)
	mockRequest.EXPECT().MakeARequest("GET", "http://www.test.com/user/user", "non_exprired_token", "", gomock.Any(), gomock.Any()).Return(mockedResp, nil).Times(1)

	_, receivedAccess, err := testFunction.CheckPrivileges("default", "")

	expectedAccess := make(map[string]interface{})

	if err != nil {
		t.Errorf("Expected no errors, received an error \"%v\"", err)
	} else if !reflect.DeepEqual(receivedAccess, expectedAccess) {
		t.Errorf("Expected no user access, received %v", receivedAccess)
	}
}

func TestCheckPrivilegesGrantedAccess(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockConfig := mocks.NewMockConfigureInterface(mockCtrl)
	mockRequest := mocks.NewMockRequestInterface(mockCtrl)
	testFunction := &jwt.Functions{Config: mockConfig, Request: mockRequest}

	cred := jwt.Credential{KeyId: "", APIKey: "fake_api_key", AccessKey: "non_exprired_token", APIEndpoint: "http://www.test.com"}

	grantedAccessJson := "{\"project_access\": " +
		"{\"test_project\": {" +
		"\"0\": \"read\"," +
		"\"1\": \"create\"," +
		"\"2\": \"read-storage\"," +
		"\"3\": \"update\"," +
		"\"4\": \"delete\"}" +
		"}}"

	mockedResp := &http.Response{
		Body:       ioutil.NopCloser(bytes.NewBufferString(grantedAccessJson)),
		StatusCode: 200,
	}

	mockConfig.EXPECT().ParseConfig("default").Return(cred).Times(1)
	mockRequest.EXPECT().MakeARequest("GET", "http://www.test.com/user/user", "non_exprired_token", "", gomock.Any(), gomock.Any()).Return(mockedResp, nil).Times(1)

	_, expectedAccess, err := testFunction.CheckPrivileges("default", "")

	receivedAccess := make(map[string]interface{})
	receivedAccess["test_project"] = map[string]interface{}{
		"0": "read",
		"1": "create",
		"2": "read-storage",
		"3": "update",
		"4": "delete"}

	if err != nil {
		t.Errorf("Expected no errors, received an error \"%v\"", err)
	} else if !reflect.DeepEqual(expectedAccess, receivedAccess) {
		t.Errorf(`Expected user access and received user access are note the same.
        Expected: %v
        Received: %v`, expectedAccess, receivedAccess)
	}
}
