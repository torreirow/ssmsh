package commands

import (
	"github.com/abiosoft/ishell"
)

const cpUsage string = `
cp usage: cp [-rR] src dest
Copy a parameter from src to dest.
  -r Copy parameters recursively
`

func cp(c *ishell.Context) {
	paths, recurse := checkRecursion(c.Args)
	if len(paths) != 2 {
		shell.Println("Expected src and dst")
		shell.Println(cpUsage)
		return
	}
	dst := parsePath(paths[1])
	err := ps.Copy(parsePath(paths[0]), dst, recurse)
	if err != nil {
		shell.Println(err)
	} else {
		// Invalidate cache for destination path
		invalidatePathAndParents(dst.Name, dst.Region)
	}
}
