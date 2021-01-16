# Cloudflare DNS IP Updater
Small CLI tool that updates the A record with the public facing ip for a particular DNS record

## Usage:

### Help
```
Using the zone id the DNS record is retrieved and the content is updated to the latest public ip

Usage:
  cloudfare-dns update [flags]

Flags:
  -r, --dns-record-name string   Name of the DNS record where the A record needs to be updated
  -h, --help                     help for update
      --ip string                Set the content of the dns record to this ip
  -t, --token string             Cloudflare API token
      --ttl int                  TTL (in seconds) to set on the DNS record (default 3600)
  -z, --zone-name string         Name of the Zone the DNS record resides in
  ```

  ### Update DNS Record 'media.burmudar.dev' in Zone 'burmudar.dev' with a static ip 127.0.0.1
  ```
  cloudflare-dns update -t $(cat TOKEN) -z burmudar.dev -r media.burmudar.dev --ip 127.0.0.1  
  Listing zones ...2 listed zones
  Locating zone: burmudar.dev ...FOUND!
  Listing DNS Records for zone 'burmudar.dev' using id 'd8e42ceecc223f1db02246f959ca741f' ...3 listed dns records
  Locatin DNS Record: media.burmudar.dev ...FOUND!
  Manually setting ip ...127.0.0.1
  Updating DNS [A media.burmudar.dev] record content with ip: 127.0.0.1
  {"result":{"id":"adda54f06cd2a","zone_id":"d8e42ceecc223f","zone_name":"burmudar.dev","name":"media.burmudar.dev","type":"A","content":"127.0.0.1","proxiable":false,"proxied":false,"ttl":3600,"locked":false,"meta":{"auto_added":false,"managed_by_apps":false,"managed_by_argo_tunnel":false,"source":"primary"},"created_on":"2021-01-16T20:22:27.655436Z","modified_on":"2021-01-16T20:22:27.655436Z"},"success":true,"errors":[],"messages":[]}
  Updated!
  ```

  ### Update DNS Record 'media.burmudar.dev' in Zone 'burmudar.dev' with our public facing ip
  ```
  cloudflare-dns update -t $(cat TOKEN) -z burmudar.dev -r media.burmudar.dev
  Listing zones ...2 listed zones
  Locating zone: burmudar.dev ...FOUND!
  Listing DNS Records for zone 'burmudar.dev' using id 'd8e42ceecc223f' ...3 listed dns records
  Locatin DNS Record: media.burmudar.dev ...FOUND!
  Discovering public ip ...169.0.60.42
  {"result":{"id":"adda54f06cd2a92ac8e","zone_id":"d8e42ceecc223f1db0","zone_name":"burmudar.dev","name":"media.burmudar.dev","type":"A","content":"169.0.0.1","proxiable":true,"proxied":false,"ttl":3600,"locked":false,"meta":{"auto_added":false,"managed_by_apps":false,"managed_by_argo_tunnel":false,"source":"primary"},"created_on":"2021-01-16T20:22:39.996491Z","modified_on":"2021-01-16T20:22:39.996491Z"},"success":true,"errors":[],"messages":[]}
  DNS [A media.burmudar.dev] content already contains: 169.0.0.1
  ```
