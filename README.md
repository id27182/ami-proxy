# ami-proxy
ami-proxy - is a service, which acts similar to a http proxy. It listens for incoming http requests, and redirects those requests to host, specified in configuration (destination host). ami-proxy attachs an auth header with token, which was received from Azure vm managed service identity endpoint (more about azure managed service identity: https://docs.microsoft.com/en-us/azure/active-directory/managed-identities-azure-resources/qs-configure-cli-windows-vm).

ami-proxy is useful for situations, when there is a need to run legacy software, initialy designed for on premise, in the cloud (in Azure). It allow to add cloud-native auth and security layer to communication between 2 or more services, if there is no way to rewrite any of them. 

There was a real use case, which inspired me to write this service: my company had a client-server app, which was designed for on premise, and I was need to move this app to an Azure. I was need to run a server part on a vm with a public IP. Clients was deployed to a lot of azure VMs within same azure subscribtion by deploy automation tool. There was no way to restrict an access to server vm by firevall rules or implement any kind of vnet peering between clients vnets and server vnet, so I've decided to put an Azure Web App, which will act lice a reverse proxy in front of my server part of app (configs are avaliable in ami-reverse-proxy-config folder). I've configured an Azure AD auth on my reverse proxy app and the next step was to find a way, how I'll add auth header to clients requests, and this is where ami-proxy come in the play.

### configuration
ami-proxy can read configuration from environment variables or from configuration file. 
Following config keys are required, in case if service is configured with environment variables: 

AMIPROXY_PROXY_DEST_PROTOCOL - protocol, which should be used in requests, to destination service (http or https)

AMIPROXY_PROXY_DEST_HOST - FQDN or ip of destination service

AMIPROXY_PROXY_DEST_PORT - port, which is listened by destination service

AMIPROXY_PROXY_DEST_RESOURCE - resource id or resource url of destination service (used when getting auth token)

AMIPROXY_PROXY_BIND_PORT - port, on which proxy will listen for incoming http requests

