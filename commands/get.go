package commands

import (
	"github.com/abiosoft/ishell"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/bwhaley/ssmsh/parameterstore"
)

const getUsage string = `
get usage: get parameter ...
Get one or more parameters.
`

// Get parameters
func get(c *ishell.Context) {
	if len(c.Args) >= 1 {
		var params []parameterstore.ParameterPath
		for _, p := range c.Args {
			params = append(params, parsePath(p))
		}
		paramsByRegion := groupByRegion(params)
		for region, params := range paramsByRegion {
			resp, err := ps.Get(params, region)
			if err != nil {
				shell.Println("Error: ", err)
			} else {
				if len(resp) >= 1 {
					for i, param := range resp {
						if aws.StringValue(param.Type) == "SecureString" {
							if !ps.Decrypt {
								resp[i].Value = aws.String("<sensitive>")
							} else if aws.StringValue(param.Value) == "<encrypted>" {
								shell.Println("Warning:", aws.StringValue(param.Name), "returned <encrypted>. Verify that your IAM role has kms:Decrypt permission for the parameter's KMS key.")
							}
						}
					}
					printResult(resp)
				}
			}
		}
	} else {
		shell.Println(getUsage)
	}
}
