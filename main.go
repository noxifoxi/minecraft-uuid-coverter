package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	// current directory or the first argument
	root := "."
	if len(os.Args) > 1 {
		root = os.Args[1]
	}

	// just walk through all files and alter them
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) == ".json" {
			alterFile(path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func alterFile(file string) {
	re := regexp.MustCompile(`SkullOwner:\{Id:\\"([0-9a-f-]{36})\\"`)

	// don't pass a wrong formatted file
	json, err := ioutil.ReadFile(file)

	if err != nil {
		log.Fatal(err)
	}

	// skip file if no match was found
	if len(re.FindSubmatch(json)) < 2 {
		return
	}

	intArray, err := convertUUID(string(re.FindSubmatch(json)[1]))
	if err != nil {
		panic(err)
	}

	// replace old UUID with int_array format
	json = re.ReplaceAll(json, []byte("SkullOwner:{Id:"+parseUUIDArray(intArray)))

	fmt.Print(file + " updated.")

	// your original file is now rip
	ioutil.WriteFile(file, json, 0666)
}

func parseUUIDArray(uuid [4]int32) string {
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
		n, err := strconv.ParseInt(slice, 16, 64)
		if err != nil {
			// fuck you
			return uuidIntArray, err
		}

		// convert int64 to int32 because stupid
		uuidIntArray[i] = int32(n)
	}

	return uuidIntArray, nil
}
