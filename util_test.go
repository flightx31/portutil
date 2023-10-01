package portutil

import (
	"fmt"
	"testing"
)

func TestUtil(t *testing.T) {
	l := L{}
	SetLogger(l)
	portConnection, err := FindOpenPort(3000, 30)
	fmt.Println(portConnection, err)
}
