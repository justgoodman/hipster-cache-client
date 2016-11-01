package tcp

import (
	"fmt"
	"strconv"
	"strings"

	"hipster-cache-client/common"
)

type ProxyClient struct {
	tcpClient *TCPClient
}

func NewProxyClient(serverAddress string, serverPort int, logger common.ILogger) *ProxyClient {
	tcpClient := &TCPClient{serverAddress: serverAddress, serverPort: serverPort, logger: logger}
	return &ProxyClient{tcpClient: tcpClient}
}

func (c *ProxyClient) InitConnection() error {
	return c.tcpClient.InitConnection()
}

func (c *ProxyClient) GetShardAddress(key string) (address string, port int, err error) {
	var result string
	command := fmt.Sprintf("GET_SHARD %s", key)
	result, err = c.tcpClient.SendMessage(command)
	if err != nil {
		return
	}
	parameters := strings.Split(result, ":")
	if len(parameters) != 2 {
		err = fmt.Errorf(`Error: Incorrect get shard response: "%s"`, result)
		return
	}
	address = parameters[0]
	portString := parameters[1]
	port, err = strconv.Atoi(portString)
	return
}
