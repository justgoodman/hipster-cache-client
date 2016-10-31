package tcp

import (
	"fmt"
	"net"
	"unicode/utf8"
	"sync"

	"hipster-cache-client/common"
)

type TCPClient struct {
	serverAddress string
	serverPort    int
	clientPort    int
	logger        common.ILogger
	conn          *net.TCPConn
	sendMessageMutex sync.Mutex
}

func NewTCPClient(clientPort int, serverAddress string, serverPort int, logger common.ILogger) *TCPClient {
	return &TCPClient{clientPort: clientPort, serverAddress: serverAddress, serverPort: serverPort, logger: logger}
}

func (c *TCPClient) InitConnection() error {
	serverTCPAddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", c.serverAddress, c.serverPort))
	if err != nil {
		return err
	}

	localTCPAddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf(":%d", c.clientPort))

	c.conn, err = net.DialTCP("tcp", localTCPAddr, serverTCPAddr)
	if err != nil {
		return err
	}
	return nil
}

func (c *TCPClient) SendMessage(message string) (string, error) {
	c.sendMessageMutex.Lock()
	defer c.sendMessageMutex.Unlock()
	var buf [512]byte
	_, err := c.conn.Write([]byte(message))
	if err != nil {
		c.logger.Errorf(`Error send message "%s", error "%s"`, message, err.Error())
		c.conn.Close()
		c.InitConnection()
		return "", err
	}

	n, err := c.conn.Read(buf[0:])
	if err != nil {
		c.logger.Errorf(`Read message error: "%s"`, err.Error())
		c.conn.Close()
		c.InitConnection()
	}

	fmt.Println(string(buf[0:n]))
	//	response, err := ioutil.ReadAll(c.conn)
/*
	if err != nil {
		c.logger.Errorf(`Error response for message "%s", error "%s"`, message, err.Error())
		return "", err
	}
*/
	//	fmt.Printf(string(response))
	fmt.Printf(string(buf[0:n]))
	// return string(buf[0:n]), nil
	return c.parseResponse(string(buf[0:n]))
}

// If response in quotes this is error
func (c *TCPClient) parseResponse(response string) (string,error) {
	var (
		firstCharacter, lastCharacter rune
		size int
	)

	firstCharacter,size = utf8.DecodeRuneInString(response)
	if string(firstCharacter) != "\"" {
		return "", fmt.Errorf(response)
	}
	result := response[size:]
	lastCharacter, size = utf8.DecodeLastRuneInString(result)
	if string(lastCharacter) != "\"" {
		return "", fmt.Errorf(response)
	}
	result = result[:(len(result)-size)]
	return result, nil
}
