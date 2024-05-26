package cmd

import (
	"fmt"
	"kafkactl/config"
	"log"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// grant permission from file
var grant = &cobra.Command{
	Use:   "grant",
	Short: "Read permissions from the file",
	Run: func(cmd *cobra.Command, args []string) {
		fileCmd()
	},
}

func init() {
	rootCmd.AddCommand(grant)
}

// func to read file yaml
func fileCmd() {
	if file == "" {
		log.Fatalf("File not specified")
	}
	file, err := os.Open(file)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	// Decode file yaml
	var filestruct config.Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&filestruct); err != nil {
		log.Fatalf("Errror decoding yaml: %v", err)
	}
	for _, user := range filestruct.RBAC.Users {
		for _, permission := range user.Permissions {
			log.Printf("Start processing user: %s\n", user.Name)
			// Grant permission for topic with user
			for _, topic := range permission.Topic {
				// check permission of user on resource
				acls, err := CheckPermission(user.Name, "topic", topic.Name, topic.PatternType)
				if err != nil {
					log.Printf("Failed to check permission for user: %s to topic %s PatterType: %s | error detail: %v\n", user.Name, topic.Name, topic.PatternType, err)
				}
				addACLs, removeACLs := CompareAcls(topic.Permission, acls)
				//grant permission
				if len(addACLs) > 0 {
					err := CreatePermission(user.Name, "topic", topic.Name, topic.PatternType, addACLs)
					if err != nil {
						log.Printf("Failed to grant permission for user: %s to topic %s with permissions %v and pattern type: %s | error detail:  %v\n", user.Name, topic.Name, addACLs, topic.PatternType, err)
					} else {
						log.Printf("Succeed to grant permission for user: %s to topic %s with permissions %v and pattern type: %s\n", user.Name, topic.Name, addACLs, topic.PatternType)
					}
				}
				//remove permission
				if len(removeACLs) > 0 {
					RemovePermission(user.Name, "topic", topic.Name, topic.PatternType, removeACLs)
				}
			}
			//Grant permission for group with user
			for _, group := range permission.Group {
				// Check permission group
				acls, err := CheckPermission(user.Name, "group", group.Name, "literal")
				if err != nil {
					log.Printf("Failed to check permission for user: %s to group %s PatterType: literal | error detail: %v\n", user.Name, group.Name, err)
				}
				addACLs, removeACLs := CompareAcls(group.Permission, acls)
				//grant permission
				if len(addACLs) > 0 {
					err := CreatePermission(user.Name, "group", group.Name, "literal", addACLs)
					if err != nil {
						log.Printf("Failed to grant permission for user: %s to group %s with permissions %v and pattern type: literal | error detail:  %v\n", user.Name, group.Name, addACLs, err)
					} else {
						log.Printf("Succeed to grant permission for user: %s to group %s with permissions %v and pattern type: literal\n", user.Name, group.Name, addACLs)
					}
				}
				//remove permission
				if len(removeACLs) > 0 {
					RemovePermission(user.Name, "group", group.Name, "literal", removeACLs)
				}

			}
			//Grant permission for cluster with user
			for _, cluster := range permission.Cluster {
				// Check permission cluster
				acls, err := CheckPermission(user.Name, "cluster", cluster.Name, "literal")
				if err != nil {
					log.Printf("Failed to check permission for user: %s to cluster %s PatterType: literal | error detail: %v\n", user.Name, cluster.Name, err)
				}
				addACLs, removeACLs := CompareAcls(cluster.Permission, acls)
				fmt.Println(cluster.Permission)
				fmt.Println(acls)
				fmt.Println(addACLs)
				fmt.Println(removeACLs)
				//grant permission
				if len(addACLs) > 0 {
					err := CreatePermission(user.Name, "cluster", cluster.Name, "literal", addACLs)
					if err != nil {
						log.Printf("Failed to grant permission for user: %s to cluster %s with permissions %v and pattern type: literal | error detail:  %v\n", user.Name, cluster.Name, addACLs, err)
					} else {
						log.Printf("Succeed to grant permission for user: %s to cluster %s with permissions %v and pattern type: literal\n", user.Name, cluster.Name, addACLs)
					}
				}
				//remove permission
				if len(removeACLs) > 0 {
					RemovePermission(user.Name, "cluster", cluster.Name, "literal", removeACLs)
				}

			}
			for _, transaction := range permission.Transaction {
				// Check permission transaction
				acls, err := CheckPermission(user.Name, "transaction", transaction.Name, transaction.PatternType)
				if err != nil {
					log.Printf("Failed to check permission for user: %s to transaction %s PatterType: %s | error detail: %v\n", user.Name, transaction.Name, transaction.PatternType, err)
				}
				addACLs, removeACLs := CompareAcls(transaction.Permission, acls)
				//grant permission
				if len(addACLs) > 0 {
					err := CreatePermission(user.Name, "transaction", transaction.Name, transaction.PatternType, addACLs)
					if err != nil {
						log.Printf("Failed to grant permission for user: %s to transaction %s with permissions %v and pattern type: %s | error detail:  %v\n", user.Name, transaction.Name, addACLs, transaction.PatternType, err)
					} else {
						log.Printf("Succeed to grant permission for user: %s to transaction %s with permissions %v and pattern type: %s\n", user.Name, transaction.Name, addACLs, transaction.PatternType)
					}
				}
				//remove permission
				if len(removeACLs) > 0 {
					RemovePermission(user.Name, "transaction", transaction.Name, transaction.PatternType, removeACLs)
				}

			}

		}
	}
}
