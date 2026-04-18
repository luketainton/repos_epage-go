# ePage (Go)

## Description

Send me an ePage via Pushover. This is a Go recreation of the original Python Flask application.

A simple web application that accepts form submissions and sends them via the Pushover API.

## Features

- Simple web form for name, email, and message
- Integration with Pushover API for push notifications
- Clean Bootstrap UI
- Easy to deploy with Docker

## Requirements

- Go 1.23 or higher
- Pushover API token and user key

## Environment Variables

Set the following environment variables:

- `PUSHOVER_API_TOKEN`: Your Pushover API token
- `PUSHOVER_USER_KEY`: Your Pushover user key
- `PORT`: Port to run the server on (default: 5000)
- `CSRF_KEY`: CSRF key for production (optional, generated automatically for development)

## Running Locally

1. Install dependencies:
```bash
go mod download
go mod tidy
```

2. Create a `.env` file (copy from `.env.example`):
```bash
cp .env.example .env
```

3. Edit `.env` with your actual Pushover credentials:
```
PUSHOVER_API_TOKEN=your_actual_token
PUSHOVER_USER_KEY=your_actual_key
PORT=5000
```

4. Run the application:
```bash
go run ./src
```

The application will automatically load the `.env` file and listen on the configured port (default: 5000).
The server will be available at `http://localhost:5000`

## Building

```bash
go build -o epage
./epage
```

## Docker

Build the Docker image:
```bash
docker build -t epage .
```

Run the container:
```bash
docker run -e PUSHOVER_API_TOKEN=your_token -e PUSHOVER_USER_KEY=your_key -p 5000:5000 epage
```

## API Integration

The application sends form data to the Pushover API with the following message format:

```
Name: <name>
Email: <email>

Message: <message>
```

With priority 1 (High) and sound "cosmic".
