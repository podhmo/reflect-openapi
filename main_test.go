package reflectopenapi_test

import (
	"log"
	"testing"
)

func TestMain(m *testing.M) {
	log.SetFlags(0)
	log.SetPrefix("# ")
	m.Run()
}
