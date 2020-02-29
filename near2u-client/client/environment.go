package client

type Environment struct {
	Name      string                 `json:name`
	SensorMap map[string]interface{} `json:sensors`
	LastModified
}