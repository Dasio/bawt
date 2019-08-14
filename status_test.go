package bawt

import (
	"reflect"
	"testing"
)

func TestNewStatus(t *testing.T) {
	tests := []struct {
		name string
		want Status
	}{
		{"Test defaults", Status{DB: "Not ok", Chat: "Not ok", HTTP: "N/A"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStatus(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatus_Update(t *testing.T) {
	type fields struct {
		DB   string
		Chat string
		HTTP string
	}
	type args struct {
		comp  string
		value string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"Test component DB", fields{DB: "Not ok", Chat: "Not ok", HTTP: "N/A"}, args{"db", "ok"}, false},
		{"Test component Chat", fields{DB: "Not ok", Chat: "Not ok", HTTP: "N/A"}, args{"chat", "ok"}, false},
		{"Test component HTTP", fields{DB: "Not ok", Chat: "Not ok", HTTP: "N/A"}, args{"http", "ok"}, false},
		{"Test component wild", fields{DB: "Not ok", Chat: "Not ok", HTTP: "N/A"}, args{"wild", "ok"}, true},
		{"Test component blank", fields{DB: "Not ok", Chat: "Not ok", HTTP: "N/A"}, args{"", "ok"}, true},
		{"Test value N/A", fields{DB: "Not ok", Chat: "Not ok", HTTP: "N/A"}, args{"db", "N/A"}, false},
		{"Test value Ok", fields{DB: "Not ok", Chat: "Not ok", HTTP: "N/A"}, args{"db", "Ok"}, false},
		{"Test value Not ok", fields{DB: "Not ok", Chat: "Not ok", HTTP: "N/A"}, args{"db", "Not ok"}, false},
		{"Test value wild", fields{DB: "Not ok", Chat: "Not ok", HTTP: "N/A"}, args{"db", "wild"}, true},
		{"Test value blank", fields{DB: "Not ok", Chat: "Not ok", HTTP: "N/A"}, args{"db", ""}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Status{
				DB:   tt.fields.DB,
				Chat: tt.fields.Chat,
				HTTP: tt.fields.HTTP,
			}
			if err := s.Update(tt.args.comp, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Status.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
