package main
import (
	"encoding/json" // <--- NEW
	"fmt"
	"io"
	"log"
	"net"
	"net/http"      // <--- NEW

	pb "sentinel-backend/gen"
	"google.golang.org/grpc"
    // (Keep your sqlite import in db.go, don't need it here)
)

type server struct {
	pb.UnimplementedMetricsServiceServer
}

func (s *server) ReportMetrics(stream pb.MetricsService_ReportMetricsServer) error {
	log.Println("New Agent Connected!")

	for {
		data, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.Empty{})
		}
		if err != nil {
			return err
		}

		// LOGIC UPDATE: Write to DB instead of just printing
		RecordMetric(data.AgentId, data.CpuUsage, data.Timestamp)
		
		// Optional: Still print so we know it's working
		fmt.Printf(">> [Stored] %s: %.2f%%\n", data.AgentId, data.CpuUsage)
	}
}
// HTTP Handler: Returns JSON data
func getMetricsHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Get data from DB
	metrics, err := GetRecentMetrics()
	if err != nil {
		http.Error(w, "Failed to fetch metrics", http.StatusInternalServerError)
		return
	}

	// 2. Set Headers (Allow CORS so Next.js can talk to us later)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 3. Encode to JSON
	json.NewEncoder(w).Encode(metrics)
}

func main() {
	
	// 1. Initialize Database
	InitDB()

	// --- START HTTP SERVER (in background) ---
	go func() {
		http.HandleFunc("/metrics", getMetricsHandler)
		log.Println("HTTP Server listening on :8080...")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()
	// ----------------------------------------

	// --- START GRPC SERVER (blocking) ---
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterMetricsServiceServer(s, &server{})

	log.Println("gRPC Server listening on :50051...")
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
