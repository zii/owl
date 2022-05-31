package db

import (
	"fmt"
	"testing"
)

func TestOrm(t *testing.T) {
	fmt.Println("escape:", literal(`cat"`))
	fmt.Println("join", joinStrings([]string{"a", "b", `c"`}))
	p := Prepare()
	p.And("a in (?)")
}
