package zahif

import (
	"fmt"
	"google.golang.org/grpc"

	pb "hurracloud.io/jawhar/internal/zahif/proto"
)

var (
	conn   *grpc.ClientConn
	Client pb.ZahifClient
)

func Connect(zahifHost string, zahifPort int) {

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", zahifHost, zahifPort), opts...)
	if err != nil {
		panic("failed to connect to zahif")
	}

	Client = pb.NewZahifClient(conn)
}

func Close() {
	conn.Close()
}
