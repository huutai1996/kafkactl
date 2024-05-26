package kafka

import (
	"log"

	"crypto/sha256"
	"crypto/sha512"

	"github.com/IBM/sarama"
	"github.com/xdg-go/scram"
)

var (
	SHA512 scram.HashGeneratorFcn = sha512.New
	SHA256 scram.HashGeneratorFcn = sha256.New
)

type XDGSSCRAMClient struct {
	*scram.Client
	*scram.ClientConversation
	scram.HashGeneratorFcn
}

// Handshake Scram

func (x *XDGSSCRAMClient) Begin(userName, password, authid string) (err error) {
	x.Client, err = x.HashGeneratorFcn.NewClient(userName, password, authid)
	if err != nil {
		return err
	}
	x.ClientConversation = x.Client.NewConversation()
	return nil
}

func (x *XDGSSCRAMClient) Step(challenge string) (response string, err error) {
	response, err = x.ClientConversation.Step(challenge)
	return
}

func (x *XDGSSCRAMClient) Done() bool {
	return x.ClientConversation.Done()
}

// Create the connection connect to kafka
func ConnectKafka(bootstrapServer string, saslEnable bool, saslUser, saslPassword, saslMechanism, saslProtocol string) (sarama.Client, error) {
	saramaConfig := sarama.NewConfig()
	if saslEnable {
		saramaConfig.Net.SASL.Enable = true
		saramaConfig.Net.SASL.User = saslUser
		saramaConfig.Net.SASL.Password = saslPassword
		saramaConfig.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
		saramaConfig.Net.SASL.Handshake = true
		saramaConfig.Version = sarama.V3_6_0_0
		saramaConfig.Net.TLS.Enable = false
		switch saslMechanism {
		case "SCRAM-SHA-512":
			saramaConfig.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &XDGSSCRAMClient{HashGeneratorFcn: SHA512} }
		case "SCRAM-SHA-256":
			saramaConfig.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &XDGSSCRAMClient{HashGeneratorFcn: SHA256} }
		default:
			log.Fatal("Tool only supported SCRAM-SHA-512, SCRAM-SHA-256")
		}

	}

	client, err := sarama.NewClient([]string{bootstrapServer}, saramaConfig)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func AdminKafka(bootstrapServer string, saslEnable bool, saslUser, saslPassword, saslMechanism, saslProtocol string) (sarama.ClusterAdmin, error) {
	saramaConfig := sarama.NewConfig()
	if saslEnable {
		saramaConfig.Net.SASL.Enable = true
		saramaConfig.Net.SASL.User = saslUser
		saramaConfig.Net.SASL.Password = saslPassword
		saramaConfig.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
		saramaConfig.Net.SASL.Handshake = true
		saramaConfig.Version = sarama.V3_6_0_0
		saramaConfig.Net.TLS.Enable = false
		switch saslMechanism {
		case "SCRAM-SHA-512":
			saramaConfig.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &XDGSSCRAMClient{HashGeneratorFcn: SHA512} }
		case "SCRAM-SHA-256":
			saramaConfig.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &XDGSSCRAMClient{HashGeneratorFcn: SHA256} }
		default:
			log.Fatal("Tool only supported SCRAM-SHA-512, SCRAM-SHA-256")
		}

	}

	adminclient, err := sarama.NewClusterAdmin([]string{bootstrapServer}, saramaConfig)
	if err != nil {
		return nil, err
	}
	return adminclient, nil
}
