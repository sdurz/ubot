package ubot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strconv"
	"strings"

	"github.com/sdurz/axon"
)

const errComplexType = errEncoding("type is complex")

type errEncoding string

func (e errEncoding) Error() string {
	return string(e)
}

type apiResponse struct {
	Ok          bool            `json:"ok"`
	ErrorCode   int64           `json:"error_code,omitempty"`
	Description string          `json:"description,omitempty"`
	Result      json.RawMessage `json:"result"`
}

func decodeJsonResponse(bytes []byte) (result interface{}, err error) {
	var reply apiResponse
	err = json.Unmarshal(bytes, &reply)
	if err != nil {
		return
	}

	if reply.Ok {
		log.Println(string(reply.Result))
		err = json.Unmarshal(reply.Result, &result)
	} else {
		err = errors.New(reply.Description)
	}
	return
}

// ApiClient serves as a mocking wrapper for http.Client
// there's no need to export this
type apiClient interface {
	GetBytes(URl string) (result []byte, err error)
	PostBytes(URL string, data interface{}) (result []byte, err error)
	GetJson(URL string) (result interface{}, err error)
	PostJson(URL string, request interface{}) (result interface{}, err error)
	PostMultipart(URL string, request axon.O) (result interface{}, err error)
}

type httpApiClient struct {
	httpClient http.Client
}

func (h *httpApiClient) GetBytes(URL string) (result []byte, err error) {
	log.Println(URL)
	resp, err := h.httpClient.Get(URL)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()

	if result, err = ioutil.ReadAll(resp.Body); err == nil {
		log.Println(string(result))
	}
	return
}

// PostBytes perform a low level post requesto to the API server.
func (h *httpApiClient) PostBytes(URL string, data interface{}) (result []byte, err error) {
	var (
		buffer []byte
		resp   *http.Response
	)
	if buffer, err = json.Marshal(data); err != nil {
		return
	}
	log.Println("posting request: ", string(buffer))
	if resp, err = h.httpClient.Post(URL, "application/json", bytes.NewBuffer(buffer)); err != nil {
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()

	if result, err = ioutil.ReadAll(resp.Body); err == nil {
		log.Println("got response: ", string(result))
	}
	return
}

// PostJson perform a get to the API server. Get JSON encoded response payload.
func (h *httpApiClient) GetJson(URL string) (result interface{}, err error) {
	var buffer []byte
	if buffer, err = h.GetBytes(URL); err != nil {
		return
	}
	result, err = decodeJsonResponse(buffer)
	return
}

// PostJson perform a JSON encoded post to the API server. Get JSON encoded response payload.
func (h *httpApiClient) PostJson(URL string, request interface{}) (result interface{}, err error) {
	var buffer []byte
	if buffer, err = h.PostBytes(URL, request); err != nil {
		return
	}
	result, err = decodeJsonResponse(buffer)
	return
}

// PostMultipart posts a multipart encoded request body to the API server
// Use this for sending files, photos, documents...
func (h *httpApiClient) PostMultipart(URL string, request axon.O) (result interface{}, err error) {
	// Prepare a form that you will submit to that URL.
	var contentType string
	var buffer *bytes.Buffer
	if contentType, buffer, err = prepareMultipart(request); err != nil {
		return
	}
	log.Println(buffer.String())

	// Submit the request
	var resp *http.Response
	if resp, err = http.Post(URL, contentType, buffer); err != nil {
		return
	}
	defer resp.Body.Close()

	// Check the response
	var bytes []byte
	if bytes, err = ioutil.ReadAll(resp.Body); err == nil {
		log.Println(string(bytes))
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", resp.Status)
	} else {
		result, err = decodeJsonResponse(bytes)
	}
	return
}

func prepareMultipart(request map[string]interface{}) (contentType string, buffer *bytes.Buffer, err error) {
	// Prepare a form that you will submit to that URL.
	buffer = &bytes.Buffer{}
	writer := multipart.NewWriter(buffer)
	// Don't forget to close the multipart writer, otherwise the request will miss the terminating boundary.
	defer writer.Close()

	for key, value := range request {
		var (
			fw io.Writer
			fr io.Reader
		)
		fw, fr, err = prepareSimpleValuePart(key, value, writer)
		if errors.Is(err, errComplexType) {
			fw, fr, err = prepareComplexValuePart(key, value, writer)
		}
		if _, err = io.Copy(fw, fr); err != nil {
			return
		}
	}
	contentType = writer.FormDataContentType()
	return
}

func prepareSimpleValuePart(name string, value interface{}, writer *multipart.Writer) (fw io.Writer, fr io.Reader, err error) {
	switch val := value.(type) {
	case string:
		if fw, err = writer.CreateFormField(name); err != nil {
			return
		}
		fr = strings.NewReader(val)
	case int64:
		if fw, err = writer.CreateFormField(name); err != nil {
			return
		}
		fr = strings.NewReader(fmt.Sprintf("%d", val))
	case float64:
		if fw, err = writer.CreateFormField(name); err != nil {
			return
		}
		fr = strings.NewReader(strconv.FormatFloat(val, 'f', 6, 64))
	case bool:
		if fw, err = writer.CreateFormField(name); err != nil {
			return
		}
		fr = strings.NewReader(fmt.Sprintf("%t", val))
	default:
		err = errComplexType
	}
	return
}

func prepareComplexValuePart(name string, value interface{}, writer *multipart.Writer) (fw io.Writer, fr io.Reader, err error) {
	switch unwrapped := value.(type) {
	case axon.O, axon.A:
		var dataBytes []byte
		if dataBytes, err = json.Marshal(unwrapped); err != nil {
			return
		}
		if fw, err = writer.CreatePart(textproto.MIMEHeader{
			"Content-Type":        {"application/json"},
			"Content-Disposition": {fmt.Sprintf("form-data; name=\"%v\"", name)},
		}); err != nil {
			return
		}
		fr = bytes.NewBuffer(dataBytes)
	case UploadFile:
		if fw, err = writer.CreatePart(textproto.MIMEHeader{
			"Content-Type":        {"application/octet-stream"},
			"Content-Disposition": {fmt.Sprintf("form-data; name=\"%v\"; filename=\"%v\"", name, unwrapped.GetName())},
		}); err != nil {
			return
		}
		fr = unwrapped.GetReader()
	default:
		log.Fatalf("Unsupported type %v", value)
	}
	return
}
