package main

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"strings"
)

type Config struct {
	profile     string
	user        string
	fingerprint string
	key_file    string
	tenancy     string
	region      string
}

const (
	userValidation    = "ocid1.user.oc1.."
	tenancyValidation = "ocid1.tenancy.oc1.."
)

func getHomeDir() (result string, err error) {
	u, err := user.Current()
	if err != nil {
		err = fmt.Errorf("can not get HomeDirectory due to: %s", err.Error())
		return "", err
	}
	result = u.HomeDir
	return result, err
}

func confirmAddProfile(s *bufio.Scanner) {
	fmt.Print("config file already exists. add a new profile? [Y/n]: ")
	for {
		s.Scan()
		input := s.Text()
		switch input {
		case "y", "Y":
			return
		case "n":
			fmt.Println("bye")
			os.Exit(0)
		default:
			fmt.Println("please enter Y or n")
			continue
		}
	}
}

func checkConfigExists(homeDir string, s *bufio.Scanner) (configExists bool, filePath string, err error) {
	filePath = homeDir + "/.oci/config"
	if f, err := os.Stat(filePath); os.IsNotExist(err) || f.IsDir() {
		configExists = false
		return configExists, filePath, nil
	} else if err == nil {
		configExists = true
		return configExists, filePath, nil
	} else {
		if err != nil {
			err = fmt.Errorf("can not create Configuration file due to: %s", err.Error())
			return false, "", err
		}
	}
	return
}

func scanField(s *bufio.Scanner, c *string, message string, validation string, errorMessage string) {
	fmt.Print(message)
	for {
		s.Scan()
		input := s.Text()
		if strings.HasPrefix(input, validation) {
			*c = input
			break
		}
		fmt.Println(errorMessage)
		fmt.Print(message)
	}
}

func receiveConfigValue(s *bufio.Scanner, c *Config) {
	scanField(s, &c.user, "Enter a user OCID: ", userValidation, "Error: Invalid OCID format. ")
	scanField(s, &c.fingerprint, "Enter a Fingerprint: ", "", "")
	scanField(s, &c.key_file, "Enter a directory path to Private key: ", "", "")
	scanField(s, &c.tenancy, "Enter a tenancy OCID: ", tenancyValidation, "Error: Invalid OCID format. ")
	scanField(s, &c.region, "Enter a region  ", "", "")
}

func configToByte(config *Config) []byte {
	str := fmt.Sprintf("\n[%s]\nuser=%s\nfingerprint=%s\nkey_file=%s\ntenancy=%s\nregion=%s\n",
		config.profile,
		config.user,
		config.fingerprint,
		config.key_file,
		config.tenancy,
		config.region,
	)
	return []byte(str)
}

func createNewConfig(byteData []byte, filePath string, c *Config) error {
	f, err := os.Create(filePath)
	if err != nil {
		err = fmt.Errorf("can not open config file due to: %s", err.Error())
		return err
	}
	_, err = f.Write(byteData)
	if err != nil {
		err = fmt.Errorf("can not write to config file due to: %s", err.Error())
		return err
	}
	fmt.Printf("Config written to %s\n", filePath)
	defer f.Close()
	return nil
}

func addNewProfile(byteData []byte, filePath string, c *Config) error {
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		err = fmt.Errorf("can not write to config file due to: %s", err.Error())
		return err
	}
	_, err = f.Write(byteData)
	if err != nil {
		err = fmt.Errorf("can not write to config file due to: %s", err.Error())
		return err
	}
	fmt.Printf("Config written to %s\n", filePath)
	defer f.Close()
	return nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	config := &Config{}

	homeDir, err := getHomeDir()
	if err != nil {
		panic(err)
	}

	configExists, filePath, err := checkConfigExists(homeDir, scanner)
	if err != nil {
		panic(err)
	}

	if configExists {
		confirmAddProfile(scanner)
		scanField(scanner, &config.profile, "enter profile name : ", "", "")
	} else {
		config.profile = "DEFAULT"
	}

	receiveConfigValue(scanner, config)

	configByteData := configToByte(config)

	if configExists {
		err = addNewProfile(configByteData, filePath, config)
		if err != nil {
			panic(err)
		}
	} else {
		err = createNewConfig(configByteData, filePath, config)
		if err != nil {
			panic(err)
		}
	}
}
