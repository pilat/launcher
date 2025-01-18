# Launcher

**Launcher** is a lightweight container designed to serve a simple, easy-to-configure web page with a list of links to other services. It’s particularly useful as an entry point for cluster resources or as a home lab landing page.

## Key Features

- **Simple Configuration**: Minimal effort required to set up and manage the page using a straightforward JSON configuration file.
- **Dynamic Backgrounds**: Add rotating background images with configurable intervals.
- **Basic Link Organization**: Easily structure links into main and additional sections with labels, descriptions, and optional icons.
- **Lightweight and Efficient**: Minimal resource usage, ideal for running in containers.

---

## Configuration File Format

The Launcher uses a JSON file for settings and web page structure. Below is a detailed explanation of the configuration format:

### Example Configuration File (`config.json`)

```json
{
  "title": "My Launcher Page",
  "background_color": "#6495ED",
  "background_interval": 10,
  "backgrounds": [
    "https://example.com/image1.jpg",
    "https://example.com/image2.jpg",
    "https://example.com/image3.jpg"
  ],
  "main_links": [
    {
      "url": "https://example.com/service1",
      "image": "https://example.com/icon1.png",
      "label": "Service 1",
      "description": "Description for service 1"
    },
    {
      "url": "https://example.com/service2",
      "image": "https://example.com/icon2.png",
      "label": "Service 2",
      "description": "Description for service 2"
    }
  ],
  "additional_links": [
    {
      "url": "https://example.com/tool1",
      "image": "https://example.com/icon3.png",
      "label": "Tool 1",
      "description": "Description for tool 1"
    },
    {
      "url": "https://example.com/tool2",
      "image": "https://example.com/icon4.png",
      "label": "Tool 2",
      "description": "Description for tool 2"
    }
  ]
}
```

### Configuration Parameters

| Parameter              | Type                 | Description                                                                 |
|------------------------|----------------------|-----------------------------------------------------------------------------|
| **`title`**            | string               | Title of the web page.                                                     |
| **`background_color`** | string (hex)         | Background color in hex format (e.g., `#FFFFFF`).                          |
| **`background_interval`** | integer           | Time interval (in seconds) for changing background images.                 |
| **`backgrounds`**      | array of strings     | URLs of background images.                                                 |
| **`main_links`**       | array of objects     | List of primary links. Each object supports:                               |
|                        |                      | • **`url`**: The URL the link points to.                                   |
|                        |                      | • **`image`**: URL of the link’s image/icon. Can also use `file:///` for local files. |
|                        |                      | • **`label`**: Display text for the link.                                  |
|                        |                      | • **`description`**: A short description of the link.                      |
| **`additional_links`** | array of objects     | Secondary links structured similarly to `main_links`.                      |

---

## Running the Application

You can run the Launcher using Docker or Docker Compose.

### Running from Docker

```bash
docker run -d \
    --name launcher \
    -p 8080:8080 \
    -v /path/to/config.json:/config.json \
    ghcr.io/pilat/launcher:latest
```

> **Note**: Replace `/path/to/config.json` with the absolute path to your configuration file.

### Running with Docker Compose

```yaml
services:
  launcher:
    image: ghcr.io/pilat/launcher:latest
    pull_policy: always
    restart: always
    volumes:
      - ./config.json:/config.json:ro
```

> **Note**: Replace `./config.json` with the path to your configuration file.

---

## Advanced Topics

### Environment Variables

Launcher supports the following environment variables:

- **`LISTEN_ADDRESS`** *(optional)*: Address and port the server listens on (default: `0.0.0.0:8080`).
- **`CONFIG_FILEPATH`** *(optional)*: Path to the configuration file. If not provided, Launcher will attempt to locate `config.json` in the current directory.

### Caching

Launcher uses a caching mechanism for background images and icons:
- Cached files are stored in `/tmp/launcher-cache`.
- Only images specified in the configuration file are allowed and cached.

### Using Static Files

Launcher can serve static assets (e.g., local images) from a specified directory using the `file:///` scheme in the configuration file. To enable this:

1. Mount the directory containing your static files into the container (e.g., `/static`).
2. Reference local files in the configuration file using the `file:///` scheme (e.g., `"file:///static/example.png"`).

This is particularly useful for including locally hosted icons or other resources.

---

## Development and Debugging

To build and run the application locally:

1. Clone the repository:
   ```bash
   git clone https://github.com/your-repo/launcher.git
   cd launcher
   ```

2. Build the application:
   ```bash
   go build -o launcher .
   ```

3. Run the application:
   ```bash
   ./launcher
   ```

4. Access the Launcher at [http://localhost:8080](http://localhost:8080).

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository.
2. Create a feature branch.
3. Commit your changes with descriptive messages.
4. Submit a pull request.

For bug reports and feature requests, please open an issue in the repository.
