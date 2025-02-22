#!/usr/bin/env python3
import os
import sys
import signal
import subprocess
import time
import psutil
import argparse

FLASK_PORT = 5000
TYPING_TUTOR_PORT = 8080
REACT_PORT = 5173

def is_port_in_use(port):
    """Check if a port is in use."""
    for conn in psutil.net_connections():
        if conn.laddr.port == port:
            return True
    return False

def kill_process_on_port(port):
    """Kill any process running on the specified port."""
    for conn in psutil.net_connections():
        if conn.laddr.port == port:
            try:
                process = psutil.Process(conn.pid)
                process.terminate()
                process.wait()
            except (psutil.NoSuchProcess, psutil.AccessDenied):
                pass

def start_flask_server():
    """Start the Flask backend server."""
    if is_port_in_use(FLASK_PORT):
        print(f"Port {FLASK_PORT} is already in use. Stopping existing process...")
        kill_process_on_port(FLASK_PORT)
        time.sleep(1)

    os.chdir('lang-portal/backend-flask')
    flask_process = subprocess.Popen(
        ['python3', '-m', 'flask', 'run'],
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE
    )
    print(f"Flask server started on port {FLASK_PORT}")
    os.chdir('../..')
    return flask_process

def start_typing_tutor():
    """Start the typing tutor frontend server."""
    if is_port_in_use(TYPING_TUTOR_PORT):
        print(f"Port {TYPING_TUTOR_PORT} is already in use. Stopping existing process...")
        kill_process_on_port(TYPING_TUTOR_PORT)
        time.sleep(1)

    os.chdir('typing-tutor')
    typing_tutor_process = subprocess.Popen(
        ['npx', 'http-server', './public', '-p', str(TYPING_TUTOR_PORT)],
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE
    )
    print(f"Typing tutor server started on port {TYPING_TUTOR_PORT}")
    os.chdir('..')
    return typing_tutor_process

def start_react_server():
    """Start the React frontend server."""
    if is_port_in_use(REACT_PORT):
        print(f"Port {REACT_PORT} is already in use. Stopping existing process...")
        kill_process_on_port(REACT_PORT)
        time.sleep(1)

    os.chdir('lang-portal/frontend-react')
    react_process = subprocess.Popen(
        ['npm', 'run', 'dev'],
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE
    )
    print(f"React server started on port {REACT_PORT}")
    os.chdir('../..')
    return react_process

def init_database():
    """Initialize the database with fresh data."""
    os.chdir('lang-portal/backend-flask')
    subprocess.run(['python3', 'init_db.py'], check=True)
    os.chdir('../..')

def main():
    print("Starting servers...")
    
    # Initialize the database with fresh data
    init_database()
    
    # Start all servers
    flask_process = start_flask_server()
    typing_tutor_process = start_typing_tutor()
    react_process = start_react_server()
    
    print("\nServers are running!")
    print("Access the typing tutor at:")
    print("  Adjectives: http://localhost:8080/public/index.html?group_id=1")
    print("  Verbs: http://localhost:8080/public/index.html?group_id=2")
    print("\nAccess the language portal at:")
    print("  http://localhost:5173")
    
    def signal_handler(signum, frame):
        print("\nStopping all servers...")
        flask_process.terminate()
        typing_tutor_process.terminate()
        react_process.terminate()
        sys.exit(0)
    
    signal.signal(signal.SIGINT, signal_handler)
    
    try:
        flask_process.wait()
        typing_tutor_process.wait()
        react_process.wait()
    except KeyboardInterrupt:
        signal_handler(None, None)

if __name__ == '__main__':
    main()
