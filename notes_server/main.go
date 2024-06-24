package main

import (
	"context"
	"flag"
	"fmt"
	"grpc-tutorial/pb"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	port        = flag.Int("port", 50051, "The server port")
	gatewayPort = flag.Int("gateway-port", 8080, "The HTTP gateway port")
	storagePath = "testdata" 
	logger      *log.Logger
)

type notesServer struct {
	pb.UnimplementedNotesServer
}

func (s *notesServer) Save(ctx context.Context, n *pb.Note) (*pb.NoteSaveReply, error) {
	logger.Printf("Received a note to save: %v", n.Title)
	err := pb.SaveToDisk(n, storagePath)

	if err != nil {
		return &pb.NoteSaveReply{Saved: false}, err
	}

	return &pb.NoteSaveReply{Saved: true}, nil
}

func (s *notesServer) Load(ctx context.Context, search *pb.NoteSearch) (*pb.Note, error) {
	logger.Printf("Received a note to load: %v", search.Keyword)
	n, err := pb.LoadFromDisk(search.Keyword, storagePath)

	if err != nil {
		return &pb.Note{}, err
	}

	return n, nil
}

func runGRPCServer() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterNotesServer(s, &notesServer{})
	reflection.Register(s)
	logger.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		logger.Fatalf("failed to serve: %v", err)
	}
}

func runHTTPGateway(ctx context.Context) {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := pb.RegisterNotesHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%d", *port), opts)
	if err != nil {
		logger.Fatalf("failed to register gRPC gateway: %v", err)
	}

	logger.Printf("HTTP gateway listening on port %d", *gatewayPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *gatewayPort), mux); err != nil {
		logger.Fatalf("failed to serve HTTP gateway: %v", err)
	}
}

func main() {
	// Initialize the logger
	logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	flag.Parse()
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go runGRPCServer()
	runHTTPGateway(ctx)
}
