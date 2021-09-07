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

// Enter a location for your config [/home/isucon/.oci/config]:
// Enter a user OCID:
// Error: Invalid OCID format. Instructions to find OCIDs: https://docs.cloud.oracle.com/Content/API/Concepts/apisigningkey.htm#Other

const (
	userValidation    = "ocid1.user.oc1.."
	tenancyValidation = "ocid1.tenancy.oc1.."
	regionValidation  = "ocid1.tenancy.oc1.."
)

var isNewConfig bool

func getHomeDir() (result string, err error) {
	u, err := user.Current()
	if err != nil {
		err = fmt.Errorf("can not get HomeDirectory due to: %s", err.Error())
		return "", err
	}
	result = u.HomeDir
	return result, err
}

func checkConfigExists(dir string, s *bufio.Scanner) (filePath string, err error) {
	filePath = dir + "/fuga.txt"
	if f, err := os.Stat(filePath); os.IsNotExist(err) || f.IsDir() {
		isNewConfig = true
		return filePath, nil
	} else if err == nil {
		fmt.Print("config file already exists. add a new profile? [Y/n]: ")
		for {
			s.Scan()
			input := s.Text()
			switch input {
			case "y", "Y":
				goto a
			case "n":
				fmt.Println("bye")
				os.Exit(0)
			default:
				fmt.Println("what?")
				continue
			}
		}
	a:
		if err != nil {
			err = fmt.Errorf("can not create Configuration file due to: %s", err.Error())
			return filePath, err
		}
		isNewConfig = false
		return filePath, nil
	} else {
		if err != nil {
			err = fmt.Errorf("can not create Configuration file due to: %s", err.Error())
			return filePath, err
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
		fmt.Print(errorMessage)
	}
}

func runScanner(s *bufio.Scanner, c *Config) {
	if !isNewConfig {
		scanField(s, &c.profile, "enter profile name : ", "", "")
	} else {
		c.profile = "DEFAULT"
	}
	scanField(s, &c.user, "Enter a user OCID: ", userValidation, "Error: Invalid OCID format. ")
	scanField(s, &c.fingerprint, "Enter a Fingerprint: ", "", "")
	scanField(s, &c.key_file, "Enter a directory path to Private key: ", "", "")
	scanField(s, &c.tenancy, "Enter a tenancy OCID: ", tenancyValidation, "Error: Invalid OCID format. ")
	scanField(s, &c.region, "Enter a region  ", regionValidation, "")
}

func createNewConfig(filePath string, c *Config) error {
	value := fmt.Sprintf("[%s]\nuser=%s\nfingerprint=%s\nkey_file=%s\ntenancy=%s\nregion=%s\n",
		c.profile,
		c.user,
		c.fingerprint,
		c.key_file,
		c.tenancy,
		c.region,
	)
	data := []byte(value)
	f, err := os.Create(filePath)
	if err != nil {
		err = fmt.Errorf("can not open config file due to: %s", err.Error())
		return err
	}
	_, err = f.Write(data)
	if err != nil {
		err = fmt.Errorf("can not write to config file due to: %s", err.Error())
		return err
	}
	fmt.Println("Config written to %s", filePath)
	defer f.Close()
	return nil
}

func addNewProfile(filePath string, c *Config) error {
	value := fmt.Sprintf("\n[%s]\nuser=%s\nfingerprint=%s\nkey_file=%s\ntenancy=%s\nregion=%s\n",
		c.profile,
		c.user,
		c.fingerprint,
		c.key_file,
		c.tenancy,
		c.region,
	)
	data := []byte(value)
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		err = fmt.Errorf("can not write to config file due to: %s", err.Error())
		return err
	}
	_, err = f.Write(data)
	if err != nil {
		err = fmt.Errorf("can not write to config file due to: %s", err.Error())
		return err
	}
	fmt.Println("Config written to %s", filePath)
	defer f.Close()
	return nil
}

func main() {
	dir, err := getHomeDir()
	if err != nil {
		panic(err)
	}

	s := bufio.NewScanner(os.Stdin)
	filePath, err := checkConfigExists(dir, s)
	if err != nil {
		panic(err)
	}

	c := &Config{}

	runScanner(s, c)

	if !isNewConfig {
		err = addNewProfile(filePath, c)
		if err != nil {
			panic(err)
		}
	} else {
		err = createNewConfig(filePath, c)
		if err != nil {
			panic(err)
		}
	}
}
