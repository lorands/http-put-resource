package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	resource "github.com/lorands/http-resource"
	"github.com/lorands/http-resource/out"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal(fmt.Sprintf("usage: %v <sources directory>", os.Args[0]))
		os.Exit(1)
	}

	var request out.Request
	inputRequest(&request)

	sourceDir := os.Args[1]

	toPathPrefix := processTemplatedTo(request.Params.To)

	fmt.Println("Target output URL: ", request.Source.URL + "/" + toPathPrefix)

	httpPut := prepareHTTPPut(request.Source)

	var re *regexp.Regexp

	if request.Params.FromRe != "" {
		re = regexp.MustCompile(request.Params.FromRe)
		// fmt.Println(re)
	}	

	workFolder := sourceDir + "/" + request.Params.From
	// readDirRecursively(sourceDir + "/" + request.Params.From)

	filepath.Walk(workFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath := path[len(workFolder):]
		if re != nil {
			if re.MatchString(relPath) {
				fmt.Println(relPath)
				if !info.IsDir() {
					httpPut(path, toPathPrefix+path[len(workFolder):])
				}
			}
		}
		return nil
	})
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
		fmt.Println("To PUT file: ", path)

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
			fmt.Println(resp)
		}

		return err
	}

	return f
}

func readDirRecursively(path string) {

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if f.IsDir() {
			readDirRecursively(path + "/" + f.Name())
		}
		fmt.Println(f.Name())
	}

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
