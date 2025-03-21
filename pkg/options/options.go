package options

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type Options struct {
	MachineID     string
	MachineFolder string

	ClientOptions    *ClientOptions
	ProjectID        string
	AvailabilityZone string
	Flavor           string
	DiskSize         int64
}

type ClientOptions struct {
	Region                string
	PrivateKeyPath        string
	ServiceAccountKeyPath string
}

func FromEnv(skipMachine bool) (*Options, error) {
	retOptions := &Options{
		ClientOptions: &ClientOptions{},
	}
	var err error
	if !skipMachine {
		retOptions.MachineID, err = fromEnvOrError("MACHINE_ID")
		if err != nil {
			return nil, err
		}

		retOptions.MachineFolder, err = fromEnvOrError("MACHINE_FOLDER")
		if err != nil {
			return nil, err
		}
	}

	retOptions.ProjectID, err = fromEnvOrError("STACKIT_PROJECT_ID")
	if err != nil {
		return nil, err
	}
	retOptions.Flavor, err = fromEnvOrError("STACKIT_FLAVOR")
	if err != nil {
		return nil, err
	}
	retOptions.AvailabilityZone, err = fromEnvOrError("STACKIT_AVAILABILITY_ZONE")
	if err != nil {
		return nil, err
	}
	diskSize, err := fromEnvOrError("STACKIT_DISK_SIZE")
	if err != nil {
		return nil, err
	}
	retOptions.DiskSize, err = extractInt64(diskSize)
	if err != nil {
		return nil, err
	}
	retOptions.ClientOptions.Region, err = fromEnvOrError("STACKIT_REGION")
	if err != nil {
		return nil, err
	}
	retOptions.ClientOptions.PrivateKeyPath, err = fromEnvOrError("STACKIT_PRIVATE_KEY_PATH")
	if err != nil {
		return nil, err
	}
	retOptions.ClientOptions.ServiceAccountKeyPath, err = fromEnvOrError("STACKIT_SERVICE_ACCOUNT_KEY_PATH")
	if err != nil {
		return nil, err
	}

	return retOptions, nil
}

func fromEnvOrError(name string) (string, error) {
	val := os.Getenv(name)
	if val == "" {
		return "", fmt.Errorf("couldn't find option %s in environment, please make sure %s is defined", name, name)
	}

	return val, nil
}

func extractInt64(s string) (int64, error) {
	re := regexp.MustCompile(`\d+`)
	match := re.FindString(s)

	if match == "" {
		return 0, fmt.Errorf("couldn't find number in %s", s)
	}

	result, err := strconv.ParseInt(match, 10, 64)
	if err != nil {
		return 0, err
	}
	return result, nil
}
