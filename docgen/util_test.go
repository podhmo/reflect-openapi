package docgen

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_toHtmlID(t *testing.T) {
	type args struct {
		operationID string
		method      string
		path        string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"simple", args{"main.FindPets", "GET", "/pets"}, "mainfindpets-get-pets"},
		{"patvar", args{"main.DeletePet", "DELETE", "/pets/{id}"}, "maindeletepet-delete-petsid"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toHtmlID(tt.args.operationID, tt.args.method, tt.args.path)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("toHtmlID() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_toDocumentInfo(t *testing.T) {
	type args struct {
		summary     string
		description string
	}
	tests := []struct {
		name   string
		args   args
		wantDi DocumentInfo
	}{
		{"directly", args{"xxx", "yyy"}, DocumentInfo{Summary: "xxx", Description: "yyy"}},
		{"python-docstring-like", args{"xxx", "yyy\n\nzzz"}, DocumentInfo{Summary: "yyy", Description: "zzz"}},
		{"decompose-description", args{"", "xxx\nyyy\nzzz"}, DocumentInfo{Summary: "xxx", Description: "yyy\nzzz"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDi := toDocumentInfo(tt.args.summary, tt.args.description)
			if diff := cmp.Diff(tt.wantDi, gotDi); diff != "" {
				t.Errorf("toDocumentInfo() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
