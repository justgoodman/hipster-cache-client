package hipsterCacheClient

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"

	"hipster-cache-client/common"
	"hipster-cache-client/tcp"
)

type HipsterCacheClient struct {
	logger             common.ILogger
	proxyClient        *tcp.ProxyClient
	serversClient      map[string]*tcp.TCPClient
	serversClientMutex sync.RWMutex
}

func NewHipsterCacheClient(proxyServerAddress string, proxyServerPort int, logger common.ILogger) *HipsterCacheClient {
	return &HipsterCacheClient{
		proxyClient:   tcp.NewProxyClient(proxyServerAddress, proxyServerPort, logger),
		logger:        logger,
		serversClient: make(map[string]*tcp.TCPClient),
	}
}

func (c *HipsterCacheClient) Init() error {
	return c.proxyClient.InitConnection()
}

func (c *HipsterCacheClient) Set(key, value string) error {
	command := fmt.Sprintf("SET %s %s", key, value)
	_, err := c.sendCommand(key, command)
	return err
}

func (c *HipsterCacheClient) Get(key string) (string, error) {
	command := fmt.Sprintf("GET %s", key)
	return c.sendCommand(key, command)
}

func (c *HipsterCacheClient) LPush(key, value string) error {
	command := fmt.Sprintf("LPUSH %s %s", key, value)
	_, err := c.sendCommand(key, command)
	return err
}

func (c *HipsterCacheClient) LSet(key string, index int, value string) error {
	command := fmt.Sprintf("LSET %s %d %s", key, index, value)
	_, err := c.sendCommand(key, command)
	return err
}

func (c *HipsterCacheClient) LRange(key string, indexStart, indexEnd int) ([]string, error) {
	var (
		size int
	)

	command := fmt.Sprintf("LRANGE %s %d %d", key, indexStart, indexEnd)
	result, err := c.sendCommand(key, command)
	if err != nil {
		return []string{}, err
	}

	// Remove quotes
	values := strings.Split(result, "\n")
	lenValues := len(values)
	res := make([]string, lenValues, lenValues)
	for i, value := range values {
		if i != lenValues-1 {
			_, size = utf8.DecodeLastRuneInString(value)
			value = value[:(len(value) - size)]
		}
		if i != 0 {
			_, size = utf8.DecodeRuneInString(value)
			value = value[size:]
		}
		res[i] = value
	}
	return res, err
}

func (c *HipsterCacheClient) LLen(key string) (int, error) {
	var (
		result string
		err    error
		lenght int
	)
	command := fmt.Sprintf("LLEN %s", key)
	result, err = c.sendCommand(key, command)
	if err != nil {
		return -1, err
	}
	lenght, err = strconv.Atoi(result)
	if err != nil {
		return -1, err
	}
	return lenght, nil
}

func (c *HipsterCacheClient) DSet(key, field, value string) error {
	command := fmt.Sprintf("DSET %s %s %s", key, field, value)
	_, err := c.sendCommand(key, command)
	return err
}

func (c *HipsterCacheClient) DGet(key, field string) (string, error) {
	command := fmt.Sprintf("DGET %s %s", key, field)
	return c.sendCommand(key, command)
}

/*
func (c *HipsterCacheClient) DGetAll(key string) ([]string, error) {
	command := fmt.Sprintf("DGETALL %s", key)
	result,err := c.sendCommand(key, command)
	return strings.Split(result,"\n"), err
}
*/

func (c *HipsterCacheClient) getServerClient(cacheServerAddress string, cacheServerPort int) (*tcp.TCPClient, error) {
	cacheServerKey := fmt.Sprintf("%s:%d", cacheServerAddress, cacheServerPort)
	c.serversClientMutex.RLock()
	cacheServerClient, ok := c.serversClient[cacheServerKey]
	c.serversClientMutex.RUnlock()
	if !ok {
		c.serversClientMutex.Lock()
		cacheServerClient = tcp.NewTCPClient(cacheServerAddress, cacheServerPort, c.logger)
		c.serversClient[cacheServerKey] = cacheServerClient
		c.serversClientMutex.Unlock()
		if err := cacheServerClient.InitConnection(); err != nil {
			return nil, err
		}
	}
	return cacheServerClient, nil
}
func (c *HipsterCacheClient) sendCommand(key string, command string) (string, error) {
	var (
		cacheServerClient *tcp.TCPClient
	)
	cacheServerAddress, cacheServerPort, err := c.proxyClient.GetShardAddress(key)
	if err != nil {
		return "", err
	}
	cacheServerClient, err = c.getServerClient(cacheServerAddress, cacheServerPort)
	if err != nil {
		return "", err
	}
	fmt.Printf("\n CacheServerClient: %#v", cacheServerClient)
	return cacheServerClient.SendMessage(command)
}
