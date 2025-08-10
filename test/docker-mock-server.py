#!/usr/bin/env python3

import json
import os
import time
from http.server import HTTPServer, BaseHTTPRequestHandler

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

def main():
    # Get server info from environment variables
    server_name = os.environ.get('SERVER_NAME', 'Mock Server')
    server_port = int(os.environ.get('SERVER_PORT', '8000'))
    
    handler_class = create_handler(server_name)
    httpd = HTTPServer(('0.0.0.0', server_port), handler_class)
    
    print(f"Starting {server_name} on port {server_port}")
    httpd.serve_forever()

if __name__ == "__main__":
    main()