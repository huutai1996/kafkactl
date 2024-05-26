package cmd

import (
	"fmt"
	"kafkactl/config"
	"kafkactl/kafka"
	"log"
	"strings"

	"github.com/IBM/sarama"
)

//function check permission

func CheckPermission(user string, resourceType string, resourceName string, resourcePatternType string) ([]string, error) {

	// Acl filter
	usertemp := "User:" + user
	temp := sarama.AclFilter{
		ResourceName:   &resourceName,
		Principal:      &usertemp,
		Operation:      sarama.AclOperationAny,
		PermissionType: sarama.AclPermissionAny,
	}

	temp.ResourceType = sarama.AclResourceType(config.ResourceTypeMap[resourceType])

	temp.ResourcePatternTypeFilter = sarama.AclResourcePatternType(config.ResourcePatternTypeMap[resourcePatternType])
	// connect to kafka
	adminclient, err := kafka.AdminKafka(bootstrapServer, saslEnabled, saslUser, saslPassword, saslMechanism, saslProtocol)
	if err != nil {
		log.Fatalf("Failed to create admin client: %v", err)
	}
	defer adminclient.Close()
	acls, err := adminclient.ListAcls(temp)
	if err != nil {
		return nil, err
	}
	fmt.Println(acls)

	// process acl
	result := []string{}
	for _, acl := range acls {
		for i := 0; i < len(acl.Acls); i++ {
			result = append(result, acl.Acls[i].Operation.String())
		}
	}
	return result, nil
}

// func to compare acl of user in file and acl of user in kafka
// Output function is addACl and removeACl
func CompareAcls(sourceACL, destACL []string) (addACL, removeACL []string) {
	if len(sourceACL) == 0 && len(destACL) == 0 {
		return []string{}, []string{}
	}

	destACLToLower := make([]string, len(destACL))
	for i := range destACL {
		destACLToLower[i] = strings.ToLower(destACL[i])
	}
	sourceACLToLower := make([]string, len(sourceACL))
	for i := range sourceACL {
		sourceACLToLower[i] = strings.ToLower(sourceACL[i])
	}
	if len(sourceACLToLower) == 0 {
		return []string{}, destACLToLower
	}
	if len(destACLToLower) == 0 {
		return sourceACL, []string{}
	}
	sourceMap := make(map[string]struct{})
	destMap := make(map[string]struct{})

	for _, a := range sourceACLToLower {
		sourceMap[a] = struct{}{}
	}

	for _, b := range destACLToLower {
		destMap[b] = struct{}{}
	}

	for a := range sourceMap {
		if _, ok := destMap[a]; !ok {
			addACL = append(addACL, a)
		}
	}

	for b := range destMap {
		if _, ok := sourceMap[b]; !ok {
			removeACL = append(removeACL, b)
		}
	}
	return addACL, removeACL
}
