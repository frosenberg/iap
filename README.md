# IAP

_A simplified GCP IAP like authentication proxy_

The concept is simple, with organization that need to handle more and more employees working remote
and more and more devices, protecting internal services and resources can't be effectively done
with tools like private networks, firewalls and VPN. IAP (Identity Aware Proxy) will take in all
requests to services, authentication the current user, make sure it has access to the requested
services/resource, and only then proxy to the backing service adding headers like `X-IAP-User`
for the service to know which user is making the request.
