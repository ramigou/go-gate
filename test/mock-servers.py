#!/usr/bin/env python3

import json
import threading
import time
from http.server import HTTPServer, BaseHTTPRequestHandler
from urllib.parse import urlparse, parse_qs

class MockServerHandler(BaseHTTPRequestHandler):
    def __init__(self, server_name, *args, **kwargs):
        self.server_name = server_name
        super().__init__(*args, **kwargs)

    def do_GET(self):
        self.send_response(200)
        self.send_header('Content-type', 'application/json')
        self.end_headers()
        
        response_data = {
            'server': self.server_name,
            'method': self.command,
            'path': self.path,
            'headers': dict(self.headers),
            'client_address': self.client_address[0],
            'timestamp': time.time()
        }
        
        if self.path == '/health':
            response_data['status'] = 'healthy'
        
        self.wfile.write(json.dumps(response_data, indent=2).encode())
        
        print(f"[{self.server_name}] {self.command} {self.path} from {self.client_address[0]}")

    def do_POST(self):
        content_length = int(self.headers.get('Content-Length', 0))
        post_data = self.rfile.read(content_length)
        
        self.send_response(200)
        self.send_header('Content-type', 'application/json')
        self.end_headers()
        
        response_data = {
            'server': self.server_name,
            'method': self.command,
            'path': self.path,
            'headers': dict(self.headers),
            'client_address': self.client_address[0],
            'body': post_data.decode('utf-8', errors='ignore'),
            'timestamp': time.time()
        }
        
        self.wfile.write(json.dumps(response_data, indent=2).encode())
        
        print(f"[{self.server_name}] {self.command} {self.path} from {self.client_address[0]}")

    def log_message(self, format, *args):
        pass

def create_handler(server_name):
    def handler_class(*args, **kwargs):
        return MockServerHandler(server_name, *args, **kwargs)
    return handler_class

def start_server(port, server_name):
    handler_class = create_handler(server_name)
    httpd = HTTPServer(('localhost', port), handler_class)
    print(f"Starting {server_name} on port {port}")
    httpd.serve_forever()

def main():
    servers = [
        (3001, "API Server 1"),
        (3002, "API Server 2"),
        (4000, "Web Server")
    ]
    
    print("Starting mock backend servers...")
    print("Press Ctrl+C to stop all servers")
    
    threads = []
    for port, name in servers:
        thread = threading.Thread(target=start_server, args=(port, name))
        thread.daemon = True
        thread.start()
        threads.append(thread)
    
    try:
        while True:
            time.sleep(1)
    except KeyboardInterrupt:
        print("\nShutting down servers...")

if __name__ == "__main__":
    main()