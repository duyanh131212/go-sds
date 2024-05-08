package core

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

//Master node
type MasterNode struct {
	api *gin.Engine 	//api server
	ln net.Listener 	//listener
	svr *grpc.Server 	//server
	nodeSvr *NodeServiceGRPCServer //node service
}

func (n *MasterNode) Init() (err error) {
	//Create new listener
	n.ln, err = net.Listen("tcp",":50001")
	if err != nil{
		return err
	}

	//Create new gRPC server
	n.svr = grpc.NewServer()

	//Create new Node gRPC server
	n.nodeSvr = GetNodeServiceGRPCServer()

	//Register service
	RegisterNodeServiceServer(n.svr, n.nodeSvr)

	//setup api using gin
	n.api = gin.Default()
	n.api.POST("/tasks", func(c *gin.Context){
		//parse payload
		var payload struct {
			Cmd string `json:"cmd"`
		}
		if err := c.ShouldBindJSON(&payload); err != nil{
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		n.nodeSvr.CmdChanel <-payload.Cmd
		c.AbortWithStatus(http.StatusOK)
	})

	return nil

}

func (n *MasterNode) Start() {
	//Start gRPC server
	go n.svr.Serve(n.ln)

	//Start api server
	_ = n.api.Run(":9000")

	//wait for exit
	n.svr.Stop()
}

var node *MasterNode

func GetMasterNode() *MasterNode {
	if node == nil{
		node = &MasterNode{}

		if err:=node.Init(); err != nil{
			panic(err)
		}
	}

	return node
}