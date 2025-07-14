package dto

type Policy struct {
	Statement []struct {
		Effect    string                 `json:"Effect"`
		Principal map[string]interface{} `json:"Principal"`
		Action    interface{}            `json:"Action"`
		Resource  interface{}            `json:"Resource"`
	} `json:"Statement"`
}
