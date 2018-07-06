package testutils

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/stretchr/testify/suite"
)

func Execute(suite *suite.Suite, cmd *exec.Cmd) {
	w := bytes.NewBuffer(nil)
	cmd.Stderr = w
	err := cmd.Run()
	suite.NoError(err)
	fmt.Printf("Stderr: %s\n", string(w.Bytes()))
}
