package cli

import (
	"errors"
)

func SplitArgs(args []string) ([]string, []string, error) {
	var (
		i         int
		arg       string
		selfArgs  []string
		childArgs []string
	)
	found := false
	for i, arg = range args {
		if arg == "--" {
			found = true
			if i+1 == len(args) {
				return nil, nil, errors.New("invalid argument")
			}

			if len(args[i+1:]) == 0 {
				return nil, nil, errors.New("invalid argument")
			}

			selfArgs = args[1:i]
			childArgs = args[i+1:]
			break
		}
	}

	if !found {
		return nil, nil, errors.New("invalid argument")
	}
	return selfArgs, childArgs, nil

}
