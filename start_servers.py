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
    return flask_process

def start_typing_tutor():
    """Start the typing tutor frontend server."""
    if is_port_in_use(TYPING_TUTOR_PORT):
        print(f"Port {TYPING_TUTOR_PORT} is already in use. Stopping existing process...")
        kill_process_on_port(TYPING_TUTOR_PORT)
        time.sleep(1)

    os.chdir('../../typing-tutor')
    http_process = subprocess.Popen(
        ['python3', '-m', 'http.server', str(TYPING_TUTOR_PORT)],
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE
    )
    print(f"Typing tutor server started on port {TYPING_TUTOR_PORT}")
    return http_process

def init_database():
    """Initialize the database with fresh data."""
    print("Initializing database...")
    os.chdir('lang-portal/backend-flask')
    subprocess.run(['python3', 'init_db.py'])
    os.chdir('../..')
    print("Database initialized successfully")

def main():
    parser = argparse.ArgumentParser(description='Start servers for the language learning application')
    parser.add_argument('--init-db', action='store_true', help='Initialize the database before starting servers')
    args = parser.parse_args()

    # Store the original working directory
    original_dir = os.getcwd()

    try:
        if args.init_db:
            init_database()

        print("Starting servers...")
        flask_process = start_flask_server()
        typing_tutor_process = start_typing_tutor()

        print("\nServers are running!")
        print("Access the typing tutor at:")
        print(f"  Adjectives: http://localhost:{TYPING_TUTOR_PORT}/public/index.html?group_id=1")
        print(f"  Verbs: http://localhost:{TYPING_TUTOR_PORT}/public/index.html?group_id=2")
        print("\nPress Ctrl+C to stop all servers")

        # Wait for keyboard interrupt
        flask_process.wait()
        typing_tutor_process.wait()

    except KeyboardInterrupt:
        print("\nStopping servers...")
        flask_process.terminate()
        typing_tutor_process.terminate()
        flask_process.wait()
        typing_tutor_process.wait()
        print("Servers stopped")

    finally:
        # Restore original working directory
        os.chdir(original_dir)

if __name__ == '__main__':
    main()
