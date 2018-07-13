# cfddns

A dynamic DNS updater for Cloudflare

## Installation

- `go build`
- Create an A record in Cloudflare for your subdomain, point it anywhere

## Usage

### `cfddns -key=<key> -email=<email> -domain=test.dev -subdomain=remote.test.dev`

#### where:

**key:** Cloudflare API Key

**email:** Cloudflare E-Mail address

**domain:** The root domain (zone) in Cloudflare

**subdomain:** The name for the DNS record in Cloudflare

#### optionally:

**httpbin:** In the format `httpbin=https://httpbin.org/ip`, this should point to your own running instance of httpbin.org's docker package.

Or, a custom API that returns:
```json
{
    "origin": "0.0.0.0"
}
```

set that up on in Cron or Windows Task Scheduler