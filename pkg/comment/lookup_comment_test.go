package comment

import (
	"testing"
)

func TestIt(t *testing.T) {
	l := NewLookup()

	cases := []struct {
		Msg     string
		Input   interface{}
		Comment string
	}{
		{
			Msg:     "func.0",
			Input:   ExtractText,
			Comment: "ExtractText extract full text of comment-group",
		},
		{
			Msg:     "func.1",
			Input:   NewLookup,
			Comment: "NewLookup is the factory function creating Lookup",
		},
		// methods is not supported yet.
		{
			Msg:     "method.0",
			Input:   l.LookupCommentTextFromFunc, // github.com/podhmo/reflect-openapi/pkg/comment.(*Lookup).LookupCommentTextFromFunc-fm
			Comment: "",
		},
	}
	for _, c := range cases {
		c := c
		t.Run(c.Msg, func(t *testing.T) {
			got, err := l.LookupCommentTextFromFunc(c.Input)
			if err != nil {
				t.Fatalf("!! %+v", err)
			}

			if want := c.Comment; want != got {
				t.Errorf("want:\n\t%q\nbut got:\n\t%q\n", want, got)
			}
		})
	}
}
