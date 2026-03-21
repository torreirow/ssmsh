package commands

import (
	"github.com/abiosoft/ishell"
)

const mvUsage string = `
mv usage: mv src dst
Move parameter from src to dst.
`

func mv(c *ishell.Context) {
	src := parsePath(c.Args[0])
	dst := parsePath(c.Args[1])
	err := ps.Move(src, dst)
	if err != nil {
		shell.Println(err)
	} else {
		// Invalidate cache for both source and destination paths
		invalidatePathAndParents(src.Name, src.Region)
		invalidatePathAndParents(dst.Name, dst.Region)
	}
}
