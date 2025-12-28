#include <iostream>
#include <fstream>
#include <sstream>
#include <string>
#include <thread>
#include <chrono>

// GRPC INCLUDES
#include <grpcpp/grpcpp.h>
#include "sentinel.grpc.pb.h"

using grpc::Channel;
using grpc::ClientContext;
using grpc::Status;
using grpc::ClientWriter;
using sentinel::MetricsService;
using sentinel::MetricData;
using sentinel::Empty;

// --- CPU LOGIC (The same as before) ---
struct CpuStats {
    long long user, nice, system, idle, iowait, irq, softirq, steal;
};

CpuStats readCpuStats() {
    std::ifstream file("/proc/stat");
    std::string line, label;
    std::getline(file, line);
    std::stringstream ss(line);
    ss >> label;
    CpuStats stats;
    ss >> stats.user >> stats.nice >> stats.system >> stats.idle 
       >> stats.iowait >> stats.irq >> stats.softirq >> stats.steal;
    return stats;
}

long long getActiveTime(const CpuStats& s) {
    return s.user + s.nice + s.system + s.irq + s.softirq + s.steal;
}

long long getIdleTime(const CpuStats& s) {
    return s.idle + s.iowait;
}
// --------------------------------------

// --- NETWORKING LOGIC (The New Part) ---
class SentinelClient {
public:
    SentinelClient(std::shared_ptr<Channel> channel)
        : stub_(MetricsService::NewStub(channel)) {}

    // The main loop that pushes data
    void StreamMetrics() {
        ClientContext context;
        Empty response;

        // 1. Open the stream
        std::unique_ptr<ClientWriter<MetricData>> writer(
            stub_->ReportMetrics(&context, &response));

        std::cout << "Connected to Server. Streaming data..." << std::endl;

        CpuStats prev = readCpuStats();

        // Infinite Loop
        while (true) {
            std::this_thread::sleep_for(std::chrono::seconds(1));
            CpuStats current = readCpuStats();

            long long totalDelta = (getActiveTime(current) + getIdleTime(current)) - 
                                   (getActiveTime(prev) + getIdleTime(prev));
            long long idleDelta = getIdleTime(current) - getIdleTime(prev);

            double usage = 0.0;
            if (totalDelta > 0) {
                usage = (double)(totalDelta - idleDelta) / totalDelta * 100.0;
            }

            // 2. Pack the data into the Protobuf message
            MetricData data;
            data.set_agent_id("Agent-007");
            data.set_cpu_usage(usage);
            data.set_timestamp(time(NULL));

            // 3. Send it!
            if (!writer->Write(data)) {
                // If Write returns false, the stream is broken (server died)
                std::cout << "Stream broken. Exiting." << std::endl;
                break;
            }

            printf("Sent: %.2f%%   \r", usage);
            fflush(stdout);

            prev = current;
        }

        // Close stream cleanly
        writer->WritesDone();
        Status status = writer->Finish();
    }

private:
    std::unique_ptr<MetricsService::Stub> stub_;
};

int main() {
    // Connect to localhost:50051 (The standard gRPC port)
    // We use "Insecure" because we aren't setting up SSL certificates yet
    SentinelClient client(grpc::CreateChannel(
        "localhost:50051", grpc::InsecureChannelCredentials()));

    client.StreamMetrics();

    return 0;
}