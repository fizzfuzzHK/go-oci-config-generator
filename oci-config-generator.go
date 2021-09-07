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
				// ファイルの中身を読み取り、変数に追加する
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

func scanField(s *bufio.Scanner, c *string, message string, validation string) {
	fmt.Print(message)
	for {
		s.Scan()
		input := s.Text()
		if strings.HasPrefix(input, validation) {
			*c = input
			break
		}
		fmt.Print("Please Enter Appropriate Value!! : ")
	}
}

// func generateConfigFile(dir string) (err error) {
// 	d := dir + "/fuga.txt"
// 	fp, err := os.OpenFile(d, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
// 	if err != nil {
// 		err = fmt.Errorf("can not create Configuration file due to: %s", err.Error())
// 		return err
// 	}
// 	defer fp.Close()
// 	return nil
// }

func runScanner(s *bufio.Scanner, c *Config) {

	if !isNewConfig {
		scanField(s, &c.profile, "enter profile name : ", "")
	} else {
		c.profile = "DEFAULT"
	}
	scanField(s, &c.user, "enter user OCID : ", userValidation)
	scanField(s, &c.fingerprint, "enter fingerprint : ", "")
	scanField(s, &c.key_file, "enter path to private_key : ", "")
	scanField(s, &c.tenancy, "enter tenancy OCID : ", tenancyValidation)
	scanField(s, &c.region, "enter region : ", regionValidation)
}

// func createInteractiveCLI(s *bufio.Scanner, c *Config) {

// 	fmt.Println("enter profile name : ")
// 	for s.Scan() {
// 		input := s.Text()
// 		if strings.HasPrefix(input, "") {
// 			c.profile = input
// 			break
// 		}
// 		fmt.Println("error : invalid character")
// 		fmt.Println("enter fingerprint : ")
// 	}

// 	fmt.Println("enter user OCID : ")
// 	for s.Scan() {
// 		input := s.Text()
// 		if strings.HasPrefix(input, "ocid.") {
// 			c.field["user"] = input
// 			break
// 		}
// 		fmt.Println("error : invalid character")
// 		fmt.Println("enter user OCID : ")
// 	}

// 	fmt.Println("enter fingerprint : ")
// 	for s.Scan() {
// 		input := s.Text()
// 		if strings.HasPrefix(input, "") {
// 			c.field["fingerprint"] = input
// 			break
// 		}
// 		fmt.Println("error : invalid character")
// 		fmt.Println("enter fingerprint : ")

// 	}

// 	fmt.Println("enter path to private_key : ")
// 	for s.Scan() {
// 		input := s.Text()
// 		if strings.HasPrefix(input, "") {
// 			c.field["key_file"] = input
// 			break
// 		}
// 		fmt.Println("error : invalid character")
// 	}

// 	fmt.Println("enter tenancy OCID : ")
// 	for s.Scan() {
// 		input := s.Text()
// 		if strings.HasPrefix(input, "ocid.") {
// 			c.field["tenancy"] = input
// 			break
// 		}
// 		fmt.Println("error : invalid character")
// 	}

// 	fmt.Println("enter region : ")
// 	for s.Scan() {
// 		input := s.Text()
// 		if strings.HasPrefix(input, "") {
// 			c.field["region"] = input
// 			break
// 		}
// 		fmt.Println("error : invalid character")
// 	}
// }
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
	defer f.Close()
	return nil
}

// func writeToConfig(filePath string, c *Config) error {
// 	if !isNewConfig {
// 		value := fmt.Sprintf("\n[%s]\nuser=%s\nfingerprint=%s\nkey_file=%s\ntenancy=%s\nregion=%s\n",
// 			c.profile,
// 			c.user,
// 			c.fingerprint,
// 			c.key_file,
// 			c.tenancy,
// 			c.region,
// 		)
// 		data := []byte(value)
// 		f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0666)
// 		if err != nil {
// 			err = fmt.Errorf("can not open config file due to: %s", err.Error())
// 			return err
// 		}
// 		_, err = f.Write(data)
// 		if err != nil {
// 			err = fmt.Errorf("can not write to config file due to: %s", err.Error())
// 			return err
// 		}
// 		defer f.Close()
// 	} else {
// 		value := fmt.Sprintf("\n[%s]\nuser=%s\nfingerprint=%s\nkey_file=%s\ntenancy=%s\nregion=%s\n",
// 			c.profile,
// 			c.user,
// 			c.fingerprint,
// 			c.key_file,
// 			c.tenancy,
// 			c.region,
// 		)
// 		data := []byte(value)
// 		f, err := os.Create(filePath)
// 		if err != nil {
// 			err = fmt.Errorf("can not write to config file due to: %s", err.Error())
// 			return err
// 		}
// 		_, err = f.Write(data)
// 		if err != nil {
// 			err = fmt.Errorf("can not write to config file due to: %s", err.Error())
// 			return err
// 		}
// 		defer f.Close()
// 	}
// 	return nil
// }

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
