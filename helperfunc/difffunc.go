package helperfunc

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strings"
)

func Help() {
	fmt.Print(`Simple Storage Service.

**Usage:**
    triple-s [-port <N>] [-dir <S>]  
    triple-s --help

**Options:**
- --help     Show this screen.
- --port N   Port number
- --dir S    Path to the directory`, "\n")
}

func AllFlags() (int, string) {
	helpFlag := flag.Bool("help", false, "help")
	portFlag := flag.Int("port", 8080, "port")
	dirFlag := flag.String("dir", "data", "dir")
	flag.Usage = Help
	flag.Parse()

	if *helpFlag == true {
		Help()
		os.Exit(0)
	}

	if validPort := isValidPort(*portFlag); !validPort {
		fmt.Fprintf(os.Stderr, "This port isn't a valid\n")
		os.Exit(1)
	}
	if validName := IsValidName(*dirFlag); !validName {
		fmt.Fprintf(os.Stderr, "This name isn't a valid\n")
		os.Exit(1)
	}

	return *portFlag, *dirFlag
}

func isValidPort(portNum int) bool {
	if portNum < 1 || portNum > 65535 {
		return false
	}

	return true
}

func IsValidName(name string) bool {
	re := regexp.MustCompile("^[a-z0-9-\\.]+$")

	if strings.Contains(name, "..") || strings.Contains(name, "--") || strings.Contains(name, "-.") || strings.Contains(name, ".-") {
		return false
	}

	if net.ParseIP(name) != nil {
		fmt.Fprintf(os.Stderr, "It's ip address: %s\n", name)
		os.Exit(1)
	}

	return re.MatchString(name)
}

func CreateDir(dirFlag string) {
	_, err := os.Stat(dirFlag)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Папка не существует. Создаем!")
		}
	} else {
		_ = os.RemoveAll(dirFlag)
	}

	_ = os.Mkdir(dirFlag, 0o766)

	_, err1 := os.Create(dirFlag + "/buckets.csv")
	if err1 != nil {
		log.Fatal(err1)
	}
}

func WriteBytesToFile(filename, bucket, dir string, data []byte) error {
	file, err := os.Create(dir + "/" + bucket + "/" + filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func checkIPAddress(ip string) {
	if net.ParseIP(ip) == nil {
		fmt.Printf("IP Address: %s - Invalid\n", ip)
	} else {
		fmt.Printf("IP Address: %s - Valid\n", ip)
	}
}
