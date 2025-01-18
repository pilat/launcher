# Launcher

**Launcher** is an extremely lightweight container designed to serve a simple, customizable web page with a list of links to other services. It's particularly useful as an entry point for cluster resources or as a home lab landing page.

This project offers seamless configuration, allowing users to define the appearance and behavior of the web page and customize the links presented to users.

---

## Configuration File Format

The Launcher uses a JSON file to define the settings and structure of the web page. Below is a sample configuration file:

### Example Configuration File (`config.json`)

```json
{
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

Root level parameters:
- `title` - Title of the page
- `background_color` - Hex code for the page's background color.
- `background_interval` - Time interval (in seconds) for changing the background images.
- `backgrounds` - Array of URLs for background images.
- `main_links` - List of primary links with url, image, label, and description.
- `additional_links` - List of additional links structured similarly to main_links.

Links level parameters:
- `url` - URL of the link
- `image` - URL of the image
- `label` - Label of the link
- `description` - Description of the link

## Usage

### Running from Docker

```bash
docker run -d \
    --name launcher \
    -p 8080:8080 \
    -v /path/to/config.json:/config.json \
    ghcr.io/pilat/launcher:latest
```

> **Note**: Replace `/path/to/config.json` with the path to your configuration file.

### Running as part of docker compose

```yaml
services:
  launcher:
    image: ghcr.io/pilat/launcher:latest
    ports:
      - 8080:8080
    volumes:
      - /path/to/config.json:/config.json
```

> **Note**: Replace `/path/to/config.json` with the path to your configuration file.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
