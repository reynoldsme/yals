# Yet Another Link Shortener (YALS)

YALS is a simple link shortening service that accepts any validly formatted URI and returns a short (10 character case sensitive alphanumeric identifier) link which may then be used in place of the original link for public or private distribution.

## Requirements

Note: This application is likely to be widely portable across any OS or hardware architectures supported by the official Go compiler, but has only been tested on Linux / AMD64.

### Container Build Requirements

A container solution such as Docker or Podman compatible with OCI specification v1.0 or greater.

### Native Build Requirements

A functional installation of the official reference implementation Go compiler of version 1.19 or greater.

## Building

Clone the repository to the directory.

`cd` to the project directory.

### Container Build

Run:

```
docker build . -t yals
docker run -p 127.0.0.1:8086:8086 yals
```

### Native Build

Run:

```
go build
./yals
```

## Usage

### Request a Shortened URL

Note: As a URL cannot contain another url unescaped, the API caller is responsible for base64 encoding the URL to be shortened. You may ask yourself, why are we not doing the typical percent-encoding / URI encoding per [RFC 3986](https://www.rfc-editor.org/rfc/rfc3986.html#page-12)? Well! https://github.com/golang/go/issues/21955

To shorten the URL `https://www.example.com/coolthing?tacos=yes`:

`curl "http://127.0.0.1:8086/api/v1/shorten/aHR0cHM6Ly93d3cuZXhhbXBsZS5jb20vY29vbHRoaW5nP3RhY29zPXllcw=="`

Response:


HTTP response code: `200`

HTTP response body:

```json
{
    "url": "http://127.0.0.1:8086/aGVsbG8hCg",
    "error": ""
}
```

### Lookup a Shortened URL

To look up the identifier `aGVsbG8hCg`:

`curl "http://127.0.0.1:8086/api/v1/lookup/aGVsbG8hCg"`

Response:

HTTP response code: `200`

HTTP response body:

```json
{
    "url": "https://www.example.com/",
    "error": ""
}

```

### Using a Shortened URL

`curl "http://127.0.0.1:8086/aGVsbG8hCg"`

Response:

HTTP response code: `302` "Moved temporarily"
HTTP Location header: `https://www.example.com/`

## Other notes

The application defaults to `http://127.0.0.1:8086` for the bind address which is reflected in the generated short links.
