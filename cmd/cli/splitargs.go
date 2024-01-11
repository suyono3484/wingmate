package cli

import (
	"errors"
	"os"
)

func SplitArgs() ([]string, []string, error) {
	var (
		i         int
		arg       string
		selfArgs  []string
		childArgs []string
	)
	found := false
	for i, arg = range os.Args {
		if arg == "--" {
			found = true
			if i+1 == len(os.Args) {
				return nil, nil, errors.New("invalid argument")
			}

			if len(os.Args[i+1:]) == 0 {
				return nil, nil, errors.New("invalid argument")
			}

			selfArgs = os.Args[1:i]
			childArgs = os.Args[i+1:]
			break
		}

		if !found {
			return nil, nil, errors.New("invalid argument")
		}
	}
	return selfArgs, childArgs, nil

}
