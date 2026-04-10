package commands

import (
	"github.com/abiosoft/ishell"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/torreirow/parsh/parameterstore"
)

const getUsage string = `
get usage: get [-d|--decrypt] parameter ...
Get one or more parameters.
  -d, --decrypt   Decrypt SecureString parameters (overrides global decrypt setting)
`

// Get parameters
func get(c *ishell.Context) {
	if len(c.Args) >= 1 {
		// Parse flags and paths
		var params []parameterstore.ParameterPath
		decryptOverride := false
		hasDecryptFlag := false

		for _, arg := range c.Args {
			if arg == "-d" || arg == "--decrypt" {
				decryptOverride = true
				hasDecryptFlag = true
			} else {
				params = append(params, parsePath(arg))
			}
		}

		if len(params) == 0 {
			shell.Println("Error: No parameters specified")
			shell.Println(getUsage)
			return
		}

		// Determine decrypt setting: flag overrides global setting
		shouldDecrypt := ps.Decrypt
		if hasDecryptFlag {
			shouldDecrypt = decryptOverride
		}

		paramsByRegion := groupByRegion(params)
		for region, params := range paramsByRegion {
			resp, err := ps.Get(params, region, shouldDecrypt)
			if err != nil {
				shell.Println("Error: ", err)
			} else {
				if len(resp) >= 1 {
					for i, param := range resp {
						if aws.StringValue(param.Type) == "SecureString" {
							if !shouldDecrypt {
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
