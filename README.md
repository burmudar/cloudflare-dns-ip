# Cloudflare DNS IP Updater
Small CLI tool that updates the A record with the public facing ip for a particular DNS record

## Usage:

### Help
```
Using the zone id the DNS record is retrieved and the content is updated to the latest public ip

Usage:
  cloudfare-dns [command]

Available Commands:
  delete       delete the DNS record with <dns-record-name>
  help         Help about any command
  list-records list DNS records present in zone <zoneName>
  update       Update a type A DNS record found in the given <zoneId> with the public IP
  version      displays version information

Flags:
  -h, --help           help for cloudfare-dns
  -t, --token string   Cloudflare API token file

Use "cloudfare-dns [command] --help" for more information about a command.
```
### The Token file
This tool uses an API Token to communicate with the Cloudflare API. The tool reads the API Token from a given file and the file
should have either 400 or 600 permissions set. Otherwise the tool will not use the token file and fail.

#### Example
Create a file with the API token
```
echo "MY TOKEN" >> token
```
Set the permissions
```
chmod 0600 token
```

### Running the tool
All commands require the token file and it is therefore a global flag. Meaning it's a required flag that should be provided before a command as specified, example:
```
cloudflare-dns -t token command
```

Each command has it's own required and optional commands which can be seen by executing the command together with the `-h` flag as per the following example:
```
cloudflare-dns -t token delete -h
Delete the DNS record with <dns-record-name> that is inside zone with <zone-name>

Usage:
  cloudfare-dns delete [flags]

Flags:
  -r, --dns-record-name strings   Name of the DNS record
  -h, --help                      help for delete
  -z, --zone-name string          Name of the Zone the DNS record resides in

Global Flags:
  -t, --token string   Cloudflare API token file
```
### Update and missing records
The update command allows one to update multiple records with either a manual IP specified on the command line, or the tool with dynamically discover the public ip. The update command has the following help output
```
Using the zone id the DNS record is retrieved and the content is updated to the latest public ip

``Using the zone id the DNS record is retrieved and the content is updated to the latest public ip

Usage:
  cloudfare-dns update [flags]

Flags:
  -r, --dns-record-names strings   Name of one or more DNS records. If more than one record is specified separated them with a comma
  -h, --help                       help for update
      --ip string                  Set the content of the dns record to this ip
      --ttl int                    TTL (in seconds) to set on the DNS record (default 3600)
  -z, --zone-name string           Name of the Zone the DNS record resides in

Global Flags:
  -t, --token string   Cloudflare API token file
```
Notice that the flag `-r,--dns-record-name` accepts `strings`, which means it accepts multiple values separated by a `,`. This allows one to update multiple subdomains as per the following example:
```
cloudflare-dns -t token update -r test,media,files -z burmudar.dev
Locating DNS record: test.burmudar.dev ...NOT FOUND
--- Creating DNS Record ---
ID........:
ZoneID....: d8e42ceecc223f1db02246f959ca741f
Name......: test.burmudar.dev
Type......: A
Content...: 169.0.54.153
Proxied...: false
Priority..: 10
TTL.......: 3600
Locating DNS record: media.burmudar.dev ...FOUND
--- Current DNS Record ---
ID........: adda54f06cd2a92ac8e51f1271058d01
ZoneID....: d8e42ceecc223f1db02246f959ca741f
ZoneName..: burmudar.dev
Name......: media.burmudar.dev
Type......: A
Content...: 169.0.54.153
Proxiable.: true
Proxied...: false
TTL.......: 300
Locked....: false
Created...: <nil>
Modified..: <nil>
Meta
AutoAdd.......: false
ManagedByApps.: false
ManagedByArgo.: false
Source........: primary
DNS Record [A media.burmudar.dev] content already contains: 169.0.54.153
Locating DNS record: files.burmudar.dev ...FOUND
--- Current DNS Record ---
ID........: 0391aa76c0335c6ea9b28669b00226b3
ZoneID....: d8e42ceecc223f1db02246f959ca741f
ZoneName..: burmudar.dev
Name......: files.burmudar.dev
Type......: A
Content...: 169.0.54.153
Proxiable.: true
Proxied...: false
TTL.......: 3600
Locked....: false
Created...: <nil>
Modified..: <nil>
Meta
AutoAdd.......: false
ManagedByApps.: false
ManagedByArgo.: false
Source........: primary
DNS Record [A files.burmudar.dev] content already contains: 169.0.54.153
```
As can be seen by the previous example, by specifying `-r test,media,files` the tool expanded the records to be fully qualified and then updated each dns record.
Another IMPORTANT thing to note: When a DNS record is NOT FOUND, the update command will CREATE the dns record! This happend in the previous example with `test`. Currently this behaviour CANNOT BE TURNED OFF but will be optional in the next release!
