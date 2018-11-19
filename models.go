package resource

import "time"

type Source struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
	Verbose  bool   `json:"verbose,omitempty"`
}

type Version struct {
	Timestamp time.Time `json:"timestamp,omitempty"`
}

type MetadataPair struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
