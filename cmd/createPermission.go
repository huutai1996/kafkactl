package cmd

import (
	"kafkactl/config"
	"kafkactl/kafka"
	"log"

	"github.com/IBM/sarama"
)

// function grant permission
func CreatePermission(user string, resourceType string, resouceName string, resourcePatternType string, permission []string) error {
	var arrayPermissions []*sarama.ResourceAcls
	// struct containing acl for kafka
	tempResource := sarama.Resource{
		ResourceName: resouceName,
	}
	tempResourceAcl := sarama.ResourceAcls{}
	// map value in yaml to lib sarama
	tempResource.ResourceType = sarama.AclResourceType(config.ResourceTypeMap[resourceType])
	tempResource.ResourcePatternType = sarama.AclResourcePatternType(config.ResourcePatternTypeMap[resourceType])
	tempResourceAcl.Resource = tempResource
	// Creating struct final contain acl for kafka
	for _, r := range permission {
		temAcl1 := &sarama.Acl{
			PermissionType: sarama.AclPermissionAllow,
			Host:           "*",
			Principal:      "User:" + user,
		}
		temAcl1.Operation = sarama.AclOperation(config.PermissionMap[r])
		tempResourceAcl.Acls = append(tempResourceAcl.Acls, temAcl1)
	}
	arrayPermissions = append(arrayPermissions, &tempResourceAcl)
	// Connect to kafka
	client, err := kafka.ConnectKafka(bootstrapServer, saslEnabled, saslUser, saslPassword, saslMechanism, saslProtocol)
	if err != nil {
		log.Fatalf("Failed to connect to kafka: %v", err)
	}
	adminclient, err := sarama.NewClusterAdminFromClient(client)
	if err != nil {
		log.Fatalf("Failed to create admin client: %v", err)
	}
	defer adminclient.Close()
	err = adminclient.CreateACLs(arrayPermissions)
	if err != nil {
		return err
	}

	return nil
}
