// It creates a new Echo instance, adds some middleware, creates a new WhyPFS node, creates a new GatewayHandler, and then
// adds a route to the Echo instance
package main

import (
	"context"
	"fmt"
	whypfs "github.com/application-research/whypfs-core"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "net/http"
	"os"
	"os/signal"
	"syscall"
	"whypfs-gateway/gateway"
)

var (
	node *whypfs.Node
	gw   *gateway.GatewayHandler
	// OsSignal signal used to shutdown
	OsSignal chan os.Signal
)

func main() {
	OsSignal = make(chan os.Signal, 1)
	GatewayRoutersConfig()
	LoopForever()
}

// LoopForever on signal processing
func LoopForever() {
	fmt.Printf("Entering infinite loop\n")

	signal.Notify(OsSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)
	_ = <-OsSignal

	fmt.Printf("Exiting infinite loop received OsSignal\n")
}

func GatewayRoutersConfig() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	whypfsPeer, err := whypfs.NewNode(whypfs.NewNodeParams{
		Ctx: ctx,
	})

	whypfsPeer.BootstrapPeers(whypfs.DefaultBootstrapPeers())
	//samplePeer, err := multiaddr.NewMultiaddr("/ip4/145.40.93.107/tcp/6745/p2p/12D3KooWBUriTeu6YoJsuSg5gCvEecim9xKdZaW8fE54LUzmWJn7")
	//pinfo2, err := peer.AddrInfoFromP2pAddr(samplePeer)
	//whypfsPeer.BootstrapPeers([]peer.AddrInfo{*pinfo2})
	//
	//content := []byte("letsrebuildtolearnnewthings!")
	//buf := bytes.NewReader(content)
	//// bafybeiawc5enlmxtwdbnts3mragh5eyhl3wn5qekvimw72igdj45lixbo4
	//n, err := whypfsPeer.AddPinFile(context.Background(), buf, nil) // default configurations
	//fmt.Println(n.Cid().String())
	node = whypfsPeer

	gw = gateway.NewGatewayHandler(node)
	if err != nil {
		panic(err)
	}

	// Routes
	e.GET("/gw/:path", GatewayResolverCheckHandler)

	// Start server
	e.Logger.Fatal(e.Start("0.0.0.0:1313"))
}

// It takes a request, and forwards it to the gateway
func GatewayResolverCheckHandler(c echo.Context) error {
	p := c.Param("path")
	req := c.Request().Clone(c.Request().Context())
	req.URL.Path = p

	gw.ServeHTTP(c.Response(), req)
	return nil
}
