package execute

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/kasaikou/markflow/pkg/models"
)

func lookPath(name string) (string, error) {
	found, err := exec.LookPath(name)
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return "", models.NewCommandNotFoundError(name)
		} else {
			panic(err)
		}
	}
	return found, nil
}

func createCmd(execution models.Execution) (path string, args []string, err error) {

	path = execution.Path.String()
	args = execution.AdditionalArgs

	switch lang := execution.Lang.String(); lang {
	case "sh", "shell":
		if path == "" {
			if s, ok := os.LookupEnv("SHELL"); ok {
				path = s
			}

			if path, err = lookPath("sh"); err != nil {
				return "", nil, errors.Join(models.NewEnvironmentNotFoundError("SHELL"), err)
			}
		}

		args = append(args, "-c", execution.Script)
		return path, args, nil

	case "bash":
		if path == "" {
			if path, err = lookPath("bash"); err != nil {
				return "", nil, err
			}
		}

		args = append(args, "-c", execution.Script)
		return path, args, nil

	case "fish":
		if path == "" {
			if path, err = lookPath("bash"); err != nil {
				return "", nil, err
			}
		}

		args = append(args, "-c", execution.Script)
		return path, args, nil

	case "py", "python":
		if path == "" {
			if path, err = lookPath("python"); err != nil {
				return "", nil, err
			}
		}

		args = append(args, "-c", execution.Script)
		return path, args, nil

	case "js", "javascript":
		if path == "" {
			if path, err = lookPath("node"); err != nil {
				return "", nil, err
			}
		}

		args = append(args, "-c", execution.Script)
		return path, args, nil

	case "ts", "typescript":

		if path == "" {
			if path, err = lookPath("tsnode"); err != nil {
				return "", nil, err
			}
		}

		args = append(args, "-c", execution.Script)
		return path, args, nil

	default:
		panic(fmt.Sprintf("unknown lang '%s'", lang))
	}
}
