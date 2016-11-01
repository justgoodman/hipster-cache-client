package test

import(
	"os"
	"fmt"
	"strconv"

	"hipster-cache-client/common"
	"hipster-cache-client"

	"github.com/juju/loggo"
)

type clientDSL struct {
	logger common.ILogger
}

func NewClientDSL() *clientDSL {
	return &clientDSL{logger: loggo.GetLogger(""),}
}

func (dsl *clientDSL) WithLogger(logger common.ILogger) *clientDSL {
	dsl.logger = logger
	return dsl
}

func (dsl *clientDSL) Do() *hipsterCacheClient.HipsterCacheClient {
	proxyServerAddress := os.Getenv("PROXY_ADDRESS")
	if proxyServerAddress == "" {
		panic(fmt.Sprintf(`Can't read PROXY_ADDRESS for integration test`))
	}

	proxyServerPortString := os.Getenv("PROXY_PORT")
	if proxyServerPortString == "" {
		panic(fmt.Sprintf(`Can't read PROXY_PORT for integration test`))
	}
	proxyServerPort, err := strconv.Atoi(proxyServerPortString)
	if err != nil {
		panic(fmt.Sprintf(`PROXY_PORT is not the int type`))
	}


	client := hipsterCacheClient.NewHipsterCacheClient(proxyServerAddress,proxyServerPort,dsl.logger)

	err = client.Init()

	if err != nil {
		panic(fmt.Sprintf(`Can't init connection for integration test`))
	}
	return client
}
