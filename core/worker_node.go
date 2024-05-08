package core

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
)

type WorkerNode struct {
	conn *grpc.ClientConn  // grpc client connection
	c    NodeServiceClient // grpc client
}

func (n *WorkerNode) Init() (err error) {
	// connect to master node
	n.conn, err = grpc.Dial("localhost:50001", grpc.WithInsecure())
	if err != nil {
		return err
	}

	// grpc client
	n.c = NewNodeServiceClient(n.conn)

	return nil
}


func (n *WorkerNode) Start(){

	// log
	fmt.Println("worker node started")

	// report status
	_, _ = n.c.ReportStatus(context.Background(), &Request{})

	// assign task
	stream, _ := n.c.AssignTask(context.Background(), &Request{})

	for{
		res, err := stream.Recv()
		if err!=nil{
			return
		}

		fmt.Println("receive command: ", res.Data)
		

		// Remove this part for security issue
		// execute command
        //parts := strings.Split(res.Data, " ")
        // if err := exec.Command(parts[0], parts[1:]...).Run(); err != nil {
        //     fmt.Println(err)
        // }
	}

}

var workerNode *WorkerNode

func GetWorkerNode() *WorkerNode {
    if workerNode == nil {
        // node
        workerNode = &WorkerNode{}

        // initialize node
        if err := workerNode.Init(); err != nil {
            panic(err)
        }
    }

    return workerNode
}
