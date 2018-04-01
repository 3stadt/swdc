package cmd

import (
	"fmt"
	"github.com/jawher/mow.cli"
)

type startOptions struct {
	testing *bool
}

func Start(cmd *cli.Cmd) {
	so := &startOptions{
		testing: cmd.BoolOpt("t testing", false, "use testing env including xdebug"),
	}
	cmd.Action = so.startAction

}

func (so *startOptions) startAction() {
	fmt.Println("...start function called")
}
