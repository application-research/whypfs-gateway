// It creates a new Echo instance, adds some middleware, creates a new WhyPFS node, creates a new GatewayHandler, and then
// adds a route to the Echo instance
package main

import (
	"context"
	"fmt"
	whypfs "github.com/application-research/whypfs-core"
	cid2 "github.com/ipfs/go-cid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"io"
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
		Ctx:       ctx,
		Datastore: whypfs.NewInMemoryDatastore(),
	})

	whypfsPeer.BootstrapPeers(whypfs.DefaultBootstrapPeers())

	node = whypfsPeer

	gw = gateway.NewGatewayHandler(node)
	if err != nil {
		panic(err)
	}

	// Routes
	e.GET("/gw/:path", GatewayResolverCheckHandler)
	e.GET("/gw/dir/:path", GatewayDirResolverCheckHandler)
	e.GET("/gw/file/:path", GatewayFileResolverCheckHandler)

	// Upload for testing
	e.POST("/upload", func(c echo.Context) error {
		file, err := c.FormFile("file")
		if err != nil {
			return err
		}
		src, err := file.Open()
		if err != nil {
			return err
		}

		addNode, err := node.AddPinFile(c.Request().Context(), src, nil)
		if err != nil {
			return err
		}
		c.Response().Write([]byte(addNode.Cid().String()))
		return nil
	})

	// Start server
	e.Logger.Fatal(e.Start("0.0.0.0:1313"))
}

func GatewayDirResolverCheckHandler(c echo.Context) error {
	p := c.Param("path")
	req := c.Request().Clone(c.Request().Context())
	req.URL.Path = p

	fmt.Printf("Request path: " + p)
	cid, err := cid2.Decode(p)

	if err != nil {
		panic(err)
	}
	//	 check if file or dir.

	rscDir, err := node.GetDirectoryWithCid(c.Request().Context(), cid)
	if err != nil {
		panic(err)
	}

	rscDir.GetNode()

	c.Response().Write([]byte("nice dir"))
	return nil
}

func GatewayFileResolverCheckHandler(c echo.Context) error {
	p := c.Param("path")
	req := c.Request().Clone(c.Request().Context())
	req.URL.Path = p

	fmt.Printf("Request path: " + p)
	cid, err := cid2.Decode(p)

	if err != nil {
		panic(err)
	}
	//	 check if file or dir.
	rsc, err := node.GetFile(c.Request().Context(), cid)
	if err != nil {
		panic(err)
	}

	content, err := io.ReadAll(rsc)

	c.Response().Write(content)
	return nil
}

// It takes a request, and forwards it to the gateway
func GatewayResolverCheckHandler(c echo.Context) error {
	p := c.Param("path")
	req := c.Request().Clone(c.Request().Context())
	req.URL.Path = p

	fmt.Printf("Request path: " + p)
	cid, err := cid2.Decode(p)

	if err != nil {
		panic(err)
	}

	rscDir, err := node.GetDirectoryWithCid(c.Request().Context(), cid)
	if err != nil {
		// 	check if file.
		rscFile, err := node.GetFile(c.Request().Context(), cid)
		content, err := io.ReadAll(rscFile)
		c.Response().Write(content)
		return nil
		if err != nil {
			panic(err)
		}
	}

	rscDir.GetNode()

	c.Response().Write([]byte("nice dir"))
	return nil
}
