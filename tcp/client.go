package tcp

import (
	"fmt"
	"net"
	"sync"
	"unicode/utf8"

	"hipster-cache-client/common"
)

type TCPClient struct {
	serverAddress    string
	serverPort       int
	logger           common.ILogger
	conn             net.Conn
	sendMessageMutex sync.Mutex
}

func NewTCPClient(serverAddress string, serverPort int, logger common.ILogger) *TCPClient {
	return &TCPClient{serverAddress: serverAddress, serverPort: serverPort, logger: logger}
}

func (c *TCPClient) InitConnection() error {
	var err error
	for i := 0; i < 10; i++ {
		c.conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", c.serverAddress, c.serverPort))
		if err != nil {
			fmt.Printf("Error connection try %d error %s", i+1, err.Error())
		} else {
			break
		}
	}
	return err
}

func (c *TCPClient) SendMessage(message string) (string, error) {
	var err error
	fmt.Printf("\n Send Message : %s", message)
	c.sendMessageMutex.Lock()
	defer c.sendMessageMutex.Unlock()
	var buf [512]byte
	if c.conn == nil {
		fmt.Printf("Error conn is nill")
		err = c.InitConnection()
		if err != nil {
			return "", err
		}
	}
	_, err = c.conn.Write([]byte(message))
	if err != nil {
		c.logger.Errorf(`Error send message "%s", error "%s"`, message, err.Error())
		c.conn.Close()
		return "", err
	}

	n, err := c.conn.Read(buf[0:])
	if err != nil {
		c.logger.Errorf(`Read message error: "%s"`, err.Error())
		c.conn.Close()
	}
	fmt.Printf(string(buf[0:n]))
	return c.parseResponse(string(buf[0:n]))
}

// If response doesn't include quotes this is error
func (c *TCPClient) parseResponse(response string) (string, error) {
	var (
		firstCharacter, lastCharacter rune
		size                          int
	)
	// Remove \n -> lastCharacter
	lastCharacter, size = utf8.DecodeLastRuneInString(response)
	if string(lastCharacter) != "\n" {
		return "", fmt.Errorf("Last Character in not \\n %s", response)
	}

	response = response[:(len(response) - size)]

	if response == "OK" {
		return "", nil
	}

	firstCharacter, size = utf8.DecodeRuneInString(response)
	if string(firstCharacter) != "\"" {
		fmt.Printf("Error FirstCharacted `%s` in string `%s`", string(firstCharacter), response)
		return "", fmt.Errorf(response)
	}
	result := response[size:]
	lastCharacter, size = utf8.DecodeLastRuneInString(result)
	if string(lastCharacter) != "\"" {
		fmt.Printf("Error LastCharacted `%s` in string `%s`", string(lastCharacter), response)
		return "", fmt.Errorf(response)
	}
	result = result[:(len(result) - size)]
	return result, nil
}
