package cmd

import (
	"fmt"
	"kafkactl/kafka"
	"log"
	"os"
	"strings"

	"encoding/json"

	"bufio"
	"bytes"

	"github.com/IBM/sarama"
	"github.com/spf13/cobra"
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"
)

// listcmd to list all topic in kafka cluster

var listacl = &cobra.Command{
	Use:   "listacl",
	Short: "List acl",
	Run: func(cmd *cobra.Command, args []string) {
		ListPermission()
	},
}

var checkgrant = &cobra.Command{
	Use:   "checkgrant",
	Short: "Check grant",
	Run: func(cmd *cobra.Command, args []string) {
		CheckGrant()
	},
}

var check = &cobra.Command{
	Use:   "check",
	Short: "check file permission valid",
	Run: func(cmd *cobra.Command, args []string) {
		Checktemplate()
	},
}

func init() {

	rootCmd.AddCommand(listacl)
	rootCmd.AddCommand(checkgrant)
	rootCmd.AddCommand(check)
}

// List all topic in kafka cluster

func ListPermission() {
	client, err := kafka.ConnectKafka(bootstrapServer, saslEnabled, saslUser, saslPassword, saslMechanism, saslProtocol)
	if err != nil {
		log.Fatalf("Failed to client connect to kafka: %v", err)
	}
	admin, err := sarama.NewClusterAdminFromClient(client)
	if err != nil {
		log.Fatalf("Failed to create admin client: %v", err)
	}

	defer admin.Close()

	resourceName := "test4"
	Principal := "User:test3"
	temp := sarama.AclFilter{
		ResourceType:              sarama.AclResourceTopic,
		ResourceName:              &resourceName,
		Principal:                 &Principal,
		ResourcePatternTypeFilter: sarama.AclPatternLiteral,
		Operation:                 sarama.AclOperationAny,
		PermissionType:            sarama.AclPermissionAny,
	}
	fmt.Println(temp)
	acls, err := admin.ListAcls(temp)
	if err != nil {
		log.Fatalf("Failed to list acls ResourceName %s Principal %s: %v", resourceName, Principal, err)
	}
	for _, acl := range acls {
		for i := 0; i < len(acl.Acls); i++ {
			fmt.Println("--------------------------")
			fmt.Println("Permission Type: ", acl.Acls[i].PermissionType.String())
			fmt.Println("User: ", strings.Split(acl.Acls[i].Principal, ":")[1])
			fmt.Println("Host: ", acl.Acls[i].Host)
			fmt.Println("Operation: ", acl.Acls[i].Operation.String())
		}

	}
}

func CheckGrant() {
	// struct containing acl for kafka
	tempResource := sarama.Resource{
		ResourceName:        "test9",
		ResourceType:        sarama.AclResourceCluster,
		ResourcePatternType: sarama.AclPatternLiteral,
	}

	// Creating struct final contain acl for kafka
	Acl := sarama.Acl{
		PermissionType: sarama.AclPermissionAllow,
		Host:           "*",
		Principal:      "User:test9",
		Operation:      sarama.AclOperationRead,
	}
	fmt.Println(tempResource)
	fmt.Println(Acl)
	// Connect to kafka
	adminclient, err := kafka.AdminKafka(bootstrapServer, saslEnabled, saslUser, saslPassword, saslMechanism, saslProtocol)
	if err != nil {
		log.Fatalf("Failed to create admin client: %v", err)
	}
	defer adminclient.Close()
	err = adminclient.CreateACL(tempResource, Acl)
	if err != nil {
		fmt.Println("error:", err)
	}
}

func Checktemplate() {
	file, err := os.Open(file)
	if err != nil {
		log.Fatalf("Error opening yaml file: %v", err)
	}
	defer file.Close()

	// Create a buffer to hold the  contents of the file
	reader := bufio.NewReader(file)
	var buffer bytes.Buffer
	buf := make([]byte, 1024)
	//read file file in chunks
	for {
		n, err := reader.Read(buf)
		if err != nil {
			if err.Error() == "EOF" {
				break
			} else {
				log.Fatalf("Error reading yaml file: %v", err)
			}
		}
		buffer.Write(buf[:n])
	}
	// parse the yaml file
	var yamldata interface{}
	err = yaml.Unmarshal(buffer.Bytes(), &yamldata)
	if err != nil {
		log.Fatalf("Error parsing yaml file: %v", err)
	}
	//convert the yaml to json
	jsondata, err := json.Marshal(yamldata)
	if err != nil {
		log.Fatalf("Error converting yaml to json: %v", err)
	}
	schema := gojsonschema.NewStringLoader(`{
			"type": "object",
			"properties": {
				"rbac": {
					"type": "object",
					"properties": {
						"users": {
							"type": "array",
							"items": {
								"type": "object",
								"properties": {
									"name": {"type": "string"},
									"permissions": {
										"type": "array",
										"items": {
											"type": "object",
											"properties": {
												"topic": {
													"type": "array",
													"items": {
														"type": "object",
														"properties": {
															"name": {"type": "string"},
															"patternType": {
																"type": "string",
																"enum": ["literal", "prefixed", "any", "match"]
															},
															"permission": {
																"type": "array",
																"items": {
																	"type": "string",
																	"enum": ["read", "write", "alter", "alterconfigs", "create", "delete", "describe", "describeconfigs"]
																}
															}
														},
														"required": ["name", "patternType", "permission"]
													}
												},
												"group": {
													"type": "array",
													"items": {
														"type": "object",
														"properties": {
															"name": {"type": "string"},
															"permission": {
																"type": "array",
																"items": {
																	"type": "string",
																	"enum": ["delete", "read", "describe"]
																}
															}
														},
														"required": ["name", "permission"]
													}
												},
												"cluster": {
													"type": "array",
													"items": {
														"type": "object",
														"properties": {
															"name": {"type": "string"},
															"permission": {
																"type": "array",
																"items": {
																	"type": "string",
																	"enum": ["alter", "alterconfigs", "clusteraction", "create", "describe", "describeconfigs"]
																}
															}
														},
														"required": ["name", "permission", "patternType"]
													}
												},
												"transaction": {
													"type": "array",
													"items": {
														"type": "object",
														"properties": {
															"name": {"type": "string"},
															"permission": {
																"type": "array",
																"items": {
																	"type": "string",
																	"enum": ["describe", "write"]
																}
															}
														},
														"required": ["name", "permission", "patternType"]
													}
												}
											},
											"oneOf": [
												{
													"required": ["topic"]
												},
												{
													"required": ["group"]
												},
												{
													"required": ["cluster"]
												},
												{
													"required": ["transaction"]
												}
											]
										}
									}
								},
								"required": ["name", "permissions"]
							}
						}
					},
					"required": ["users"]
				}
			},
			"required": ["rbac"]
		}`)
	documentLoader := gojsonschema.NewBytesLoader(jsondata)
	//validate the schema
	result, err := gojsonschema.Validate(schema, documentLoader)
	if err != nil {
		log.Fatalf("Error while validating: %s", err)
	}
	if result.Valid() {
		fmt.Println("yaml file is valid")
	} else {
		fmt.Println("Error: yaml file is not valid")
		fmt.Println("Detail:")
		for _, desc := range result.Errors() {
			fmt.Printf("%s| %v\n", desc, desc.Context().String())
		}
	}

}
