package ubot

import (
	"reflect"
	"testing"
)

func TestVAsInteger(t *testing.T) {
	var (
		err   error
		value int64
	)
	v := V{5.}
	if value, err = v.AsInteger(); err != nil {
		t.Fail()
	}
	if value != 5 {
		t.Fatalf("value != 5, %v", value)
	}
}

func TestV_AsFloat(t *testing.T) {
	var (
		err   error
		value float64
	)
	v := V{5.}
	if value, err = v.AsFloat(); err != nil {
		t.Fail()
	}
	if value != 5 {
		t.Fatalf("value != 5, %v", value)
	}
}

func TestV_AsString(t *testing.T) {
	var (
		err   error
		value string
	)
	v := V{"value"}
	if value, err = v.AsString(); err != nil {
		t.Fail()
	}
	if value != "value" {
		t.Fatalf("value != 5, %v", value)
	}
}

func TestV_AsBool(t *testing.T) {
	var (
		err   error
		value bool
	)
	v := V{true}
	if value, err = v.AsBool(); err != nil {
		t.Fail()
	}
	if value != true {
		t.Fatalf("value != 5, %v", value)
	}
}

func TestV_AsObject(t *testing.T) {
	var (
		err error
	)
	o := map[string]interface{}{
		"value": 1.,
	}
	v := V{o}
	if _, err = v.AsObject(); err != nil {
		t.Fatal(err)
	}
}

func TestV_AsArray(t *testing.T) {
	var (
		err error
	)
	o := []interface{}{1, 2, 3}
	v := V{o}
	if _, err = v.AsArray(); err != nil {
		t.Fatal(err)
	}
}

func TestO_Get(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name       string
		o          *O
		args       args
		wantResult V
		wantErr    bool
	}{

		{
			name:       "single prop",
			o:          &O{"value": 1},
			args:       args{"value"},
			wantResult: V{1},
			wantErr:    false,
		},
		{
			name:       "single prop err",
			o:          &O{"value1": 0},
			args:       args{"value"},
			wantResult: V{nil},
			wantErr:    true,
		},
		{
			name: "nested prop",
			o: &O{
				"value": map[string]interface{}{
					"inner": 1,
				},
			},
			args:       args{"value.inner"},
			wantResult: V{1},
			wantErr:    false,
		},
		{
			name: "nested prop err",
			o: &O{
				"value": map[string]interface{}{
					"innerNo": 1,
				},
			},
			args:       args{"value.inner"},
			wantResult: V{nil},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := tt.o.Get(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("O.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("O.Get() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestO_GetObject(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name       string
		o          *O
		args       args
		wantResult O
		wantErr    bool
	}{
		{
			name: "object ok",
			o: &O{
				"inner": map[string]interface{}{},
			},
			args:       args{"inner"},
			wantResult: O{},
			wantErr:    false,
		},
		{
			name: "object ko",
			o: &O{
				"inner": 1,
			},
			args:       args{"inner"},
			wantResult: nil,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := tt.o.GetObject(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("O.GetObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("O.GetObject() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestO_GetString(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name       string
		o          *O
		args       args
		wantResult string
		wantErr    bool
	}{
		{
			name: "has string",
			o: &O{
				"path": "value",
			},
			args:       args{"path"},
			wantResult: "value",
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := tt.o.GetString(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("O.GetString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResult != tt.wantResult {
				t.Errorf("O.GetString() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
