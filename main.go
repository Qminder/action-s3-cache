package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	action := Action{
		Action:    os.Getenv("ACTION"),
		Bucket:    os.Getenv("BUCKET"),
		S3Class:   os.Getenv("S3_CLASS"),
		Key:       fmt.Sprintf("%s.zip", os.Getenv("KEY")),
		Artifacts: strings.Split(strings.TrimSpace(os.Getenv("ARTIFACTS")), "\n"),
	}

	switch act := action.Action; act {
	case PutAction:
		if len(action.Artifacts[0]) <= 0 {
			log.Fatal("No artifacts patterns provided")
		}

		log.Printf("Storing cached object. Key: %s", action.Key)
		if err := Zip(action.Key, action.Artifacts); err != nil {
			log.Fatal(err)
		}

		if err := PutObject(action.Key, action.Bucket, action.S3Class); err != nil {
			log.Fatal(err)
		}
	case GetAction:
		log.Printf("Trying to restore cache. Key: %s", action.Key)
		exists, err := ObjectExists(action.Key, action.Bucket)
		if err != nil {
			log.Fatalf("Failed to check if objext exist in S3: %s", err)
		}

		// Get and unzip if object exists
		if exists {
			log.Println("Cache hit. Downloading.")
			if err := GetObject(action.Key, action.Bucket); err != nil {
				log.Print(err)
				return
			}

			if err := Unzip(action.Key); err != nil {
				log.Print(err)
			}
		} else {
			log.Println("Cache miss.")
		}

		err = setOutput("cache-hit", strconv.FormatBool(exists))
		if err != nil {
			log.Fatalf("Failed to set output: %s", err)
		}

	case DeleteAction:
		if err := DeleteObject(action.Key, action.Bucket); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("Action \"%s\" is not allowed. Valid options are: [%s, %s, %s]", act, PutAction, DeleteAction, GetAction)
	}
}

func setOutput(name, value string) error {
	file := os.Getenv("GITHUB_OUTPUT")
	if file == "" {
		return errors.New("GITHUB_OUTPUT env variable not specified")
	}

	return appendToFile(file, fmt.Sprintf("%s=%s\n", name, value))
}

func appendToFile(file, content string) error {
	f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	_, err = f.WriteString(content)
	closeErr := f.Close()
	if err != nil {
		return err
	}

	return closeErr
}
