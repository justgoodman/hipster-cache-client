package client

import (
	"sync"
	"fmt"
	"strings"
	"strconv"

	"hipster-cache-client/tcp"
	"hipster-cache-client/common"
)

type HipsterCacheClient struct {
	clientPort int
	logger common.ILogger
	proxyClient *tcp.ProxyClient
	serversClient map[string]*tcp.TCPClient
	serversClientMutex sync.RWMutex
}

func NewHipsterCacheClient(proxyServerAddress string,proxyServerPort,clientPort int,logger common.ILogger) *HipsterCacheClient {
	return &HipsterCacheClient{
		proxyClient: tcp.NewProxyClient(clientPort, proxyServerAddress, proxyServerPort, logger),
		clientPort: clientPort,
		logger: logger,
		serversClient: make(map[string]*tcp.TCPClient),
	}
}

func (c *HipsterCacheClient) Init() error {
	return c.proxyClient.InitConnection()
}

func (c *HipsterCacheClient) Set(key,value string) error {
	command := fmt.Sprintf("SET %s %s", key, value)
	_, err := c.sendCommand(key, command)
	return err
}

func (c *HipsterCacheClient) Get(key string) (string,error) {
	command := fmt.Sprintf("GET %s", key)
	return c.sendCommand(key, command)
}

func (c *HipsterCacheClient) LPush(key, value string) error {
	command := fmt.Sprintf("LPUSH %s %s", key, value)
	_, err := c.sendCommand(key, command)
	return err
}

func (c *HipsterCacheClient) LSet(key string, index int, value string) error {
	command := fmt.Sprintf("LSET %s %d %s",key, index, value)
	_, err := c.sendCommand(key, command)
	return err
}

func (c *HipsterCacheClient) LRange(key string, indexStart,indexEnd int) ([]string,error) {
	command := fmt.Sprintf("LRANGE %s %d %d", key, indexStart, indexEnd)
	result, err := c.sendCommand(key,command)
	return strings.Split(result,"\n"), err
}

func (c *HipsterCacheClient) LLen(key string) (int,error) {
	var (
		result string
		err error
		lenght int
	)
	command := fmt.Sprintf("LLEN %s", key)
	result,err = c.sendCommand(key, command)
	if err != nil {
		return -1, err
	}
	lenght,err = strconv.Atoi(result)
	if err != nil {
		return -1, err
	}
	return lenght, nil
}

func (c *HipsterCacheClient) DSet(key, field, value string) error {
	command := fmt.Sprintf("DSET %s %s %s", key,field,value)
	_, err := c.sendCommand(key, command)
	return err
}

func (c *HipsterCacheClient) DGet(key, field string) (string,error) {
	command := fmt.Sprintf("DGET %s %s", key, field)
	return c.sendCommand(key, command)
}

func (c *HipsterCacheClient) DGetAll(key string) ([]string, error) {
	command := fmt.Sprintf("DGETALL %s", key)
	result,err := c.sendCommand(key, command)
	return strings.Split(result,"\n"), err
}

func (c *HipsterCacheClient) getServerClient(cacheServerAddress string, cacheServerPort int) *tcp.TCPClient {
	cacheServerKey := fmt.Sprintf("%s:%d", cacheServerAddress,cacheServerPort)
	c.serversClientMutex.RLock()
	cacheServerClient, ok := c.serversClient[cacheServerKey]
	c.serversClientMutex.RUnlock()
	if !ok {
		c.serversClientMutex.Lock()
		cacheServerClient = tcp.NewTCPClient(c.clientPort, cacheServerAddress, cacheServerPort, c.logger)
		c.serversClient[cacheServerKey] = cacheServerClient
		c.serversClientMutex.Unlock()
		cacheServerClient.InitConnection()
	}
	return cacheServerClient
}
func (c *HipsterCacheClient) sendCommand(key string,command string) (string,error) {
	cacheServerAddress, cacheServerPort, err := c.proxyClient.GetShardAddress(key)
	if err != nil {
		return "",err
	}
	cacheServerClient := c.getServerClient(cacheServerAddress,cacheServerPort)
	return cacheServerClient.SendMessage(command)
}
