# goku.dev

A lightweight, personal profile website built with Go. This application serves a markdown-based profile page with a clean, responsive design.

## Features

- **Markdown-driven content**: Edit your profile in [profile.md](profile.md) using simple markdown syntax
- **Go-powered**: Fast, efficient web server with minimal dependencies
- **Responsive design**: Clean, modern UI that works on all devices
- **Static asset support**: Serves CSS, images, and other static files
- **Single binary deployment**: Compile to a single executable for easy deployment

## Prerequisites

- Go 1.24.2 or higher
- Git (for cloning the repository)

## Quick Start

### Development

1. **Clone the repository**
   ```bash
   git clone https://github.com/goku-devv/goku.dev.git
   cd goku.dev
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Run the server**
   ```bash
   go run main.go
   ```

4. **Visit your profile**

   Open your browser and navigate to `http://localhost:8080`

### Customizing Your Profile

Edit [profile.md](profile.md) to update your profile content. The server will read this file on each request, so you can modify it while the server is running.

## Project Structure

```
goku.dev/
├── main.go              # Main application code
├── profile.md           # Your profile content (markdown)
├── templates/
│   └── layout.html      # HTML template
├── static/
│   ├── style.css        # Stylesheet
│   └── goku.jpg         # Profile image
├── go.mod               # Go module definition
└── go.sum               # Go module checksums
```

## Building for Production

### Build for Current Platform

```bash
go build -o profile main.go
```

### Cross-Platform Builds

**Linux (64-bit)**
```bash
GOOS=linux GOARCH=amd64 go build -o profile main.go
```

**Windows (64-bit)**
```bash
GOOS=windows GOARCH=amd64 go build -o profile.exe main.go
```

**macOS (Intel)**
```bash
GOOS=darwin GOARCH=amd64 go build -o profile main.go
```

**macOS (Apple Silicon)**
```bash
GOOS=darwin GOARCH=arm64 go build -o profile main.go
```

## Deployment

### Linux Server Deployment

1. **Build the Linux binary**
   ```bash
   GOOS=linux GOARCH=amd64 go build -o profile main.go
   ```

2. **Transfer files to your server**
   ```bash
   scp profile user@server:/path/to/app/
   scp profile.md user@server:/path/to/app/
   scp -r templates/ user@server:/path/to/app/
   scp -r static/ user@server:/path/to/app/
   ```

3. **Run on the server**
   ```bash
   chmod +x profile
   ./profile
   ```

### Running as a Service (systemd)

Create a systemd service file at `/etc/systemd/system/goku-profile.service`:

```ini
[Unit]
Description=Goku Profile Website
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/path/to/app
ExecStart=/path/to/app/profile
Restart=always

[Install]
WantedBy=multi-user.target
```

Enable and start the service:
```bash
sudo systemctl enable goku-profile
sudo systemctl start goku-profile
```

### Reverse Proxy with Nginx

```nginx
server {
    listen 80;
    server_name goku.dev;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

## Dependencies

- [gomarkdown/markdown](https://github.com/gomarkdown/markdown) - Markdown parsing and rendering

## Configuration

The server runs on port **8080** by default. To change the port, modify line 77 in [main.go](main.go):

```go
http.ListenAndServe(":8080", nil)  // Change port here
```

## Development Tips

- The server logs all requests and errors to stdout
- External links in markdown automatically open in new tabs
- Markdown headings automatically get IDs for anchor linking
- Static files are served from the `/static/` URL path

## License

MIT License - feel free to use this for your own profile website!

## Author

**Goku**
- Email: hi.im@goku.dev
- X: [@goku_dev](https://x.com/goku_dev)
- GitHub: [@goku-devv](https://github.com/goku-devv)

---

Built with ❤️ using Go
