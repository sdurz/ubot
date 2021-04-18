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
)

type errEncoding string

func (e errEncoding) Error() string {
	return string(e)
}

const errComplexType = errEncoding("type is complex")

func decodeJsonResponse(bytes []byte) (result interface{}, err error) {
	var reply UReply
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

type ApiClient interface {
	GetBytes(URl string) (result []byte, err error)
	PostBytes(URL string, data interface{}) (result []byte, err error)
	GetJson(URL string) (result interface{}, err error)
	PostJson(URL string, request interface{}) (result interface{}, err error)
	PostMultipart(URL string, request O) (result interface{}, err error)
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

	result, err = ioutil.ReadAll(resp.Body)
	if err == nil {
		log.Println(string(result))
	}
	return
}

func (h *httpApiClient) PostBytes(URL string, data interface{}) (result []byte, err error) {
	reqBody, err := json.Marshal(data)
	if err != nil {
		return
	}

	log.Println("posting request: ", string(reqBody))
	resp, err := h.httpClient.Post(URL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()

	result, err = ioutil.ReadAll(resp.Body)
	if err == nil {
		log.Println("got response: ", string(result))
	}
	return
}

func (h *httpApiClient) GetJson(URL string) (result interface{}, err error) {
	bytes, err := h.GetBytes(URL)
	if err != nil {
		return
	}
	result, err = decodeJsonResponse(bytes)
	return
}

func (h *httpApiClient) PostJson(URL string, request interface{}) (result interface{}, err error) {
	bytes, err := h.PostBytes(URL, request)
	if err != nil {
		return
	}
	result, err = decodeJsonResponse(bytes)
	return
}

// Multipart POSTing (for sending objects)
func (h *httpApiClient) PostMultipart(URL string, request O) (result interface{}, err error) {
	// Prepare a form that you will submit to that URL.
	var contentType string
	var buffer *bytes.Buffer
	if contentType, buffer, err = prepareMultipart(request); err != nil {
		return
	}
	log.Println(string(buffer.Bytes()))

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
	switch value.(type) {
	case O, A:
		var dataBytes []byte
		if dataBytes, err = json.Marshal(value); err != nil {
			return
		}
		if fw, err = writer.CreatePart(textproto.MIMEHeader{
			"Content-Type":        {"application/json"},
			"Content-Disposition": {fmt.Sprintf("form-data; name=\"%v\"", name)},
		}); err != nil {
			return
		}
		fr = bytes.NewBuffer(dataBytes)
	default:
		log.Fatalf("Unsupported type %v", value)
	}
	return
}
