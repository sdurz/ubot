package ubot

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

func TestNewBytesUploadFile(t *testing.T) {
	type args struct {
		fileName  string
		fileBytes []byte
	}
	tests := []struct {
		name       string
		args       args
		wantResult UploadFile
		wantErr    bool
	}{
		{
			name: "test1",
			args: args{
				fileName:  "afilename",
				fileBytes: []byte{1, 2, 3, 4},
			},
			wantResult: &bytesUploadFile{
				fileName:  "afilename",
				fileBytes: []byte{1, 2, 3, 4},
			},
			wantErr: false,
		},
		{
			name: "test2 err zero filename",
			args: args{
				fileName:  "",
				fileBytes: []byte{1, 2, 3, 4},
			},
			wantResult: nil,
			wantErr:    true,
		},
		{
			name: "test3 nil bytes",
			args: args{
				fileName:  "dummyval",
				fileBytes: nil,
			},
			wantResult: nil,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := NewBytesUploadFile(tt.args.fileName, tt.args.fileBytes)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBytesUploadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("NewBytesUploadFile() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestNewReaderUploadFile(t *testing.T) {
	type args struct {
		fileName string
		reader   io.Reader
	}
	var byteData []byte = []byte{1, 2, 3, 4}
	tests := []struct {
		name       string
		args       args
		wantResult UploadFile
		wantErr    bool
	}{
		{
			name: "test1",
			args: args{
				fileName: "afilename",
				reader:   bytes.NewReader(byteData),
			},
			wantResult: &readerUploadFile{
				fileName: "afilename",
				reader:   bytes.NewReader(byteData),
			},
			wantErr: false,
		},
		{
			name: "test2 zero filename",
			args: args{
				fileName: "",
				reader:   bytes.NewReader(byteData),
			},
			wantResult: nil,
			wantErr:    true,
		},
		{
			name: "test3 nil reader",
			args: args{
				fileName: "aname",
				reader:   nil,
			},
			wantResult: nil,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := NewReaderUploadFile(tt.args.fileName, tt.args.reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewReaderUploadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("NewReaderUploadFile() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
