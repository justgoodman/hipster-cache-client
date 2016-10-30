package HipsterCacheClient

import (
	"sync"
)

type HipsterCacheClient struct {
	proxyServerAddress string
	proxyClient *tcp.TCPClient
	serversClient map[string]*tcp.TCPClient
	serversClientMutex sync.RWMutex
}

func NewHipsterCacheClient(proxyServerAddress string) *HipsterCacheClient {
	return &HipsterCacheClient{proxyServerAddress: proxyServerAddress, serversClient: make(map[string]*tcp.TCPClient)}
}

func (c *HipsterCacheClient) Set(key,value string) error {
	command := fmt.Sprintf("SET %s %s", key, value)
	_, err := c.sendCommand(key, command)
	return err
}

func (c *HipsterCacheClient) Get(key string) string,error {
	command := fmt.Sprintf("GET %s", key, value)
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
}

func (c *HipsterCacheClient) LRange(key string, indexStart,indexEnd int) ([]string,error) {
	command := fmt.Sprintf("LRANGE %s %d %d", key, indexStart, indexEnd)
	result, err := c.sendCommand(key,command)
	return strings.Split(result,"\n"), err
}

func (c *HipsterCacheClient) LLen(key string) (int,error) {
	command := fmt.Sprintf("LLEN %s", key)
	result,err := c.sendCommand(key, command)
	return int(result), err
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

func (c *HipsterCacheClient) DGetAll(key) ([]string, error) {
	command := fmt.Sprintf("DGETALL %s", key)
	result,err := c.sendCommand(key, command)
	return strings.Split(result,"\n"), err
}

func (c *HipsterCacheClient) getServerClient(cacheServerAddress string) *tcp.TCPClient {
	c.serversClientMutex.RLock()
	cacheServerClient, ok := s.serversClient[cacheServerAddress]
	c.serverClientMutex.RUnlock()
	if !ok {
		c.serversClientMutes.Lock()
		cacheServerClient = NewTCPClient(cacheServerAddress)
		c.serversClient[cacheServerAddress] = cacheServerClient
		c.serversClientMutex.Unlock()
	}
	return cacheServerClient
}
func (c *HipsterCacheClient) sendCommand(key string,command string) (string,error) {
	cacheServerAddress := c.proxyClient.getCacheServerAddress(key)
	cacheServerClient := c.getServerClient(cacheServerAddress)
	return cacheServerClient.sendMessage(command)
}
