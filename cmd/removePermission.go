package cmd

import (
	"kafkactl/config"
	"kafkactl/kafka"
	"log"

	"github.com/IBM/sarama"
)

func RemovePermission(user, resourceType, resourceName, resourcePatternType string, permissions []string) {

	// Acl filter
	usertemp := "User:" + user
	ACLFilter := sarama.AclFilter{
		ResourceName:   &resourceName,
		Principal:      &usertemp,
		Operation:      sarama.AclOperationAny,
		PermissionType: sarama.AclPermissionAllow,
	}

	ACLFilter.ResourceType = sarama.AclResourceType(config.ResourceTypeMap[resourceType])

	ACLFilter.ResourcePatternTypeFilter = sarama.AclResourcePatternType(config.ResourcePatternTypeMap[resourcePatternType])

	// connect to kafka
	adminclient, err := kafka.AdminKafka(bootstrapServer, saslEnabled, saslUser, saslPassword, saslMechanism, saslProtocol)
	if err != nil {
		log.Fatalf("Failed to create admin client: %v", err)
	}
	defer adminclient.Close()
	// loop delete permission
	for _, p := range permissions {
		tempACLFilter := ACLFilter
		tempACLFilter.Operation = sarama.AclOperation(config.PermissionMap[p])

		_, err := adminclient.DeleteACL(tempACLFilter, true)
		if err != nil {
			log.Printf("Failed to delete acl user: %s, resoucetype %s, resourcename %s, permission: %s, error: %v ", user, resourceType, resourceName, p, err)
		}
		log.Printf("Succeed to delete acl user: %s, resoucetype %s, resourcename %s, permission: %s, error: %v ", user, resourceType, resourceName, p, err)
	}
}
