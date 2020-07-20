# Weather Proxy

A web application to proxy Weather Underground updates from weather stations to PostgreSQL.

## Installation

- Install PostgreSQL, create a database and user
- Grab the latest version from releases or build Weather Proxy via `make build`
- Generate a configuration for Weather Proxy via `./weatherproxy -g`
- Edit the configuration with the correct values
- Run Weather Proxy on port 80 (or use NGINX/Apache to proxy to it)
- Point the hostname `rtupdate.wunderground.com` at the proxy for your local network
- Enjoy having your own weather data!

For systemd service example, see [systemd/weatherproxy.service](systemd/weatherproxy.service).
