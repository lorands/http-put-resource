package out

import (
	"github.com/lorands/http-resource"
)

type Request struct {
	Source resource.Source `json:"source"`
	Params Params          `json:"params"`
}

type Params struct {
	From   string `json:"from"`
	FromRe string `json:"from-re-filter"`
	To     string `json:"to"`
}

type Response struct {
	Version  resource.Version        `json:"version"`
	Metadata []resource.MetadataPair `json:"metadata"`
}
