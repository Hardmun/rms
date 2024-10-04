**HOW TO INSTALL**

1. **Build an executable file using `go build` (CGO_ENABLED=1 must be enabled)**
   ```bash
   export CGO_ENABLED=1
   go build -o rms rms.go
   ```

2. **Install as a service**
   ```bash
   ./rms --install
   ```

3. **Check the status**
   ```bash
   systemctl status rms.service
   ```

   Example output:
   ```
   ● rms.service - rms
        Loaded: loaded (/etc/systemd/system/rms.service; enabled; vendor preset: enabled)
        Active: active (running) since Sun 2024-09-15 14:44:10 UTC; 17s ago
     Main PID: 6218 (rms)
        Tasks: 6 (limit: 23644)
        Memory: 5.0M
        CGroup: /system.slice/rms.service
             └─6218 /usr/local/rms/rms
   Sep 15 14:44:10 ubnt systemd[1]: Started rms.
   ```

   If the service is not running, execute the following commands:
    - a) Check if the service exists at `/etc/systemd/system/rms.service`
    - b) Run the following commands to start and enable the service:
      ```bash
      systemctl start rms.service
      systemctl enable rms.service
      ```

4. **Edit program settings**  
   Open the settings file for editing:
   ```bash
   nano /usr/local/rms/settings.json
   ```

  ```json
  {
  "serverPgsql": "Postgres server name in local domain",
  "portPgsql": "Postgres server port",
  "databasePgsql": "Postgres database name",
  "loginPgsql": "Postgres username",
  "passwordPgsql": "Postgres password",
  "httpServer": "Base URL, e.g., https://myshop.com",
  "httpPort": "Port if necessary, e.g., 8080. It's a local port",
  "urlPort": "Port for redirection. e.g. https://myshop.com:9090/exchange",
  "httpRedirectPrefix": "redirection resource. https://myshop.com/exchange redirect to our server 95.134.23.44:8080",
  "fileDirectory": "Path to files, e.g., /var/1C/storage",
  "requestBlockMin": "1440 (Request frequency in minutes, e.g., once per day)(number)",
  "connLimit": "Maximum simultaneous connections per user(number)",
  "userTtlMin": "Session duration on server in minutes(number)",
  "emailTech": {
    "emailTo": "Email for requests",
    "SMTPPort": "SMTP port, e.g., 587",
    "SMTPServer": "SMTP server name, e.g., smtp.gmail.com",
    "password": "Email password",
    "email": "Full email address, e.g., requestserver@gmail.com"
  },
  "tables": {
    "collection": {},
    "images": {},
    "details": {}
  }
}
```

**Note**: The `tables` field mapping is automatically populated using an external program via the API.

5. **SSL Configuration**  
   Create an `ssl` directory and place your `.key` and `.crt` files inside it, then restart the RMS service:
   ```bash
   mkdir /usr/local/rms/ssl
   systemctl restart rms.service
   ```

---