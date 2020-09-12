package agent

import (
	"fmt"
	"google.golang.org/grpc"

	pb "hurracloud.io/jawhar/internal/agent/proto"
)

var (
	conn   *grpc.ClientConn
	Client pb.HurraAgentClient
)

func Connect(agentHost string, agentPort int) {

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", agentHost, agentPort), opts...)
	if err != nil {
		panic("failed to connect to agent")
	}

	Client = pb.NewHurraAgentClient(conn)
}

func Close() {
	conn.Close()
}
