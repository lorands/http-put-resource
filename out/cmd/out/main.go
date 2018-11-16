package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"

	// "io/ioutil"

	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"time"

	resource "github.com/lorands/http-put-resource"
	"github.com/lorands/http-put-resource/out"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal(fmt.Sprintf("usage: %v <sources directory>", os.Args[0]))
		os.Exit(1)
	}

	sourceDir := os.Args[1]

	say("Input directory set.", sourceDir)

	var request out.Request
	inputRequest(&request)

	say("Request params set: %v\n", request)

	toPathPrefix := processTemplatedTo(request.Params.To)

	say("Target output URL: %v\n", request.Source.URL+"/"+toPathPrefix)

	httpPut := prepareHTTPPut(request.Source)

	var re *regexp.Regexp

	if request.Params.FromRe != "" {
		re = regexp.MustCompile(request.Params.FromRe)
	}

	workFolder := sourceDir + "/" + request.Params.From

	say("workfolder is %v\n", workFolder)

	filepath.Walk(workFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath := path[len(workFolder):]
		if re != nil {
			if re.MatchString(relPath) {
				say("Matched path  %v\n", relPath)
				if !info.IsDir() {
					httpPut(path, toPathPrefix+"/"+path[len(workFolder):])
				}
			}
		} else {
			if !info.IsDir() {
				httpPut(path, toPathPrefix+"/"+path[len(workFolder):])
			}
		}
		return nil
	})

	//TODO put url to metadata
	var metadata []resource.MetadataPair
	timestamp := time.Now()
	version := resource.Version{
		Timestamp: timestamp,
	}
	//output to stdout...
	response := out.Response{
		Version:  version,
		Metadata: metadata,
	}

	outputResponse(response)

}

//process path from env variables
func processTemplatedTo(tmpl string) string {
	envMap, _ := envToMap()
	t := template.Must(template.New("tmpl").Parse(tmpl))
	var b bytes.Buffer
	t.Execute(&b, envMap)
	return b.String()
}

func envToMap() (map[string]string, error) {
	envMap := make(map[string]string)
	var err error

	for _, v := range os.Environ() {
		split_v := strings.Split(v, "=")
		envMap[split_v[0]] = split_v[1]
	}

	return envMap, err
}

func prepareHTTPPut(src resource.Source) func(path string, to string) error {

	client := &http.Client{}

	f := func(path string, to string) error {
		say("To PUT file: %v\n", path)

		var reader io.Reader

		file, err := os.Open(path)
		defer file.Close()
		reader = bufio.NewReader(file)
		req, err := http.NewRequest("PUT", src.URL+"/"+to, reader)
		if err != nil {
			return err
		}
		req.SetBasicAuth(src.Username, src.Password)

		resp, err := client.Do(req)

		if err != nil {
			fatal(fmt.Sprintf("Error response from http. %v\n", resp), err)
		}

		return err
	}

	return f
}

func fatal(message string, err error) {
	fmt.Fprintf(os.Stderr, "error %s: %s\n", message, err)
	os.Exit(1)
}

func inputRequest(request *out.Request) {
	if err := json.NewDecoder(os.Stdin).Decode(request); err != nil {
		log.Fatal("reading request from stdin", err)
	}
}

func outputResponse(response out.Response) {
	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		log.Fatal("writing response to stdout", err)
	}
}

func say(message string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, message, args...)
}
