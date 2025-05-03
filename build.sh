#!/bin/bash

# Build the WebAssembly version
GOOS=js GOARCH=wasm go build -o main.wasm

# Copy the wasm_exec.js file from the correct location
cp "/usr/lib/go/lib/wasm/wasm_exec.js" .

# Create a simple server script
cat > server.py << 'EOF'
from http.server import HTTPServer, SimpleHTTPRequestHandler
import sys
import socket

class CORSRequestHandler(SimpleHTTPRequestHandler):
    def end_headers(self):
        self.send_header('Access-Control-Allow-Origin', '*')
        self.send_header('Access-Control-Allow-Methods', 'GET')
        self.send_header('Cache-Control', 'no-store, no-cache, must-revalidate')
        return super().end_headers()

def find_free_port():
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
        s.bind(('', 0))
        return s.getsockname()[1]

if __name__ == '__main__':
    port = find_free_port()
    if len(sys.argv) > 1:
        port = int(sys.argv[1])
    server_address = ('', port)
    httpd = HTTPServer(server_address, CORSRequestHandler)
    print(f"Starting server on port {port}")
    print(f"Open your browser to: http://localhost:{port}")
    httpd.serve_forever()
EOF

echo "Build complete! Run the server with:"
echo "python3 server.py [port]"
echo "If no port is specified, a free port will be chosen automatically" 