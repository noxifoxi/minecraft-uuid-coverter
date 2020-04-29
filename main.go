package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// UUIDRegex is the regex to find UUIDs
const UUIDRegex = `([0-9a-f]{8}-([0-9a-f]{4}-){3}[0-9a-f]{12})`

// TrUUIDRegex is the simple regex to find trimmed UUIDs
const TrUUIDRegex = `([0-9a-f]{32})`

func main() {
	// yey flags
	uuidP := flag.String("uuid", "", "UUID to convert")
	helpP := flag.Bool("help", false, "help?")
	fileP := flag.String("file", "", "Specify a file to change")
	dirP := flag.String("dir", ".", "Directory to scan for files")
	extP := flag.String("ext", "", "Only scan these extensions")
	recursiveP := flag.Bool("r", false, "Scan directory recursive")
	simulateP := flag.Bool("simulate", false, "Don't change files, only simulate the process")
	flag.Parse()

	if *helpP {
		println("this unholy place has no help")
		return
	}

	// do have to deal with arguments?
	if len(os.Args) > 1 {
		var r *regexp.Regexp

		var uuid *string = &os.Args[1]
		if *uuidP != "" {
			uuid = uuidP
		}

		if matched, _ := regexp.MatchString(UUIDRegex, *uuid); matched {
			r = regexp.MustCompile(UUIDRegex) // UUID
		} else if matched, _ := regexp.MatchString(TrUUIDRegex, *uuid); matched {
			r = regexp.MustCompile(TrUUIDRegex) // Trimmed UUID
		}

		if r != nil {
			// return the converted UUID
			intArray, _ := convertUUID(r.FindStringSubmatch(*uuid)[0])
			fmt.Println(stringifyArray(intArray))
			return
		}
	}

	if *fileP != "" {
		alterFile(*fileP, *simulateP)
	} else {
		// just walk through all files and alter them
		err := filepath.Walk(*dirP, func(path string, info os.FileInfo, err error) error {
			if info.Name() == ".git" || (!*recursiveP && path != "." && info.IsDir()) {
				return filepath.SkipDir
			}

			if info.IsDir() || filepath.Ext(path) == ".exe" {
				return nil
			}

			if *extP == "" || filepath.Ext(path) == *extP {
				alterFile(path, *simulateP)
			}
			return nil
		})
		if err != nil {
			panic(err)
		}
	}
}

func alterFile(file string, simulate bool) {
	re := regexp.MustCompile(UUIDRegex)

	rFile, err := ioutil.ReadFile(file)

	if err != nil {
		log.Fatal(err)
		return
	}

	// skip file if no match was found
	if len(re.FindSubmatch(rFile)) < 2 {
		fmt.Println(file + " skipped.")
		return
	}

	// replace UUID with int_array format
	rFile = re.ReplaceAllFunc(rFile, func(match []byte) []byte {
		intArray, _ := convertUUID(string(re.FindSubmatch(match)[1]))
		return []byte(stringifyArray(intArray))
	})

	fmt.Println(file + " updated.")

	if !simulate {
		// your original file is now rip
		ioutil.WriteFile(file, rFile, 0666)
	}
}

func stringifyArray(uuid [4]int32) string {
	parsed := "[I;"

	for i, n := range uuid {
		parsed += fmt.Sprintf("%d", n)
		if i < len(uuid)-1 {
			parsed += ","
		}
	}

	return parsed + "]"
}

func convertUUID(uuid string) ([4]int32, error) {
	// remove "-"
	uuid = strings.ReplaceAll(uuid, "-", "")

	var uuidIntArray [4]int32

	// for real, why?
	if len(uuid) != 32 {
		return uuidIntArray, errors.New(uuid + " is not a valid UUID")
	}

	for i := 0; i < 4; i++ {
		slice := uuid[i*8 : i*8+8]

		// convert hex to int64 because stupid
		n, _ := strconv.ParseInt(slice, 16, 64)

		// convert int64 to int32 because stupid
		uuidIntArray[i] = int32(n)
	}

	return uuidIntArray, nil
}
