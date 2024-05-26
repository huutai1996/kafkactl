package config

type Config struct {
	RBAC RBAC `yaml:"rbac"`
}

type RBAC struct {
	Users []User `yaml:"users"`
}

type User struct {
	Name        string        `yaml:"name"`
	Permissions []Permissions `yaml:"permissions"`
}

type Permissions struct {
	Topic       []TopicDetail       `yaml:"topic,omitempty"`
	Group       []GroupDetail       `yaml:"group,omitempty"`
	Cluster     []ClusterDetail     `yaml:"cluster,omitempty"`
	Transaction []TransactionDetail `yaml:"transaction,omitempty"`
}

type TopicDetail struct {
	Name        string   `yaml:"name"`
	PatternType string   `yaml:"patternType"`
	Permission  []string `yaml:"permission"`
}

type GroupDetail struct {
	Name       string   `yaml:"name"`
	Permission []string `yaml:"permission"`
}

type ClusterDetail struct {
	Name       string   `yaml:"name"`
	Permission []string `yaml:"permission"`
}

type TransactionDetail struct {
	Name        string   `yaml:"name"`
	Permission  []string `yaml:"permission"`
	PatternType string   `yaml:"patternType"`
}

var ResourceTypeMap = map[string]int{
	"any":             1,
	"topic":           2,
	"group":           3,
	"cluster":         4,
	"transaction":     5,
	"delegationtoken": 6,
}

var ResourcePatternTypeMap = map[string]int{
	"any":      1,
	"match":    2,
	"literal":  3,
	"prefixed": 4,
}

var PermissionMap = map[string]int{
	"any":             1,
	"all":             2,
	"read":            3,
	"write":           4,
	"create":          5,
	"delete":          6,
	"alter":           7,
	"describe":        8,
	"clusteraction":   9,
	"describeconfigs": 10,
	"alterconfigs":    11,
	"idempotentwrite": 12,
}
