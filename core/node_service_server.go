package core

import "context"

type NodeServiceGRPCServer struct {
	UnimplementedNodeServiceServer
	//channel to receive command
	CmdChanel chan string
}

func (n NodeServiceGRPCServer) ReportStatus(ctx context.Context, request *Request) (*Response, error) {
	return &Response{Data:"ok"}, nil
}

func (n NodeServiceGRPCServer) AssignTask (request *Request, server NodeService_AssignTaskServer) error {
	for {
		select {
		case cmd:= <- n.CmdChanel:
			if err:=server.Send(&Response{Data:cmd}); err != nil {
				return err
			}
		}
	}
}

var server *NodeServiceGRPCServer

func GetNodeServiceGRPCServer() *NodeServiceGRPCServer {
	if server == nil{
		server = &NodeServiceGRPCServer{
			CmdChanel: make(chan string),
		}
	}
	return server
}