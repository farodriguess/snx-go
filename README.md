# snxgo

snxgo is a utility to connect to Checkpoint VPN without needing a browser. This project is based on Ralf Schlatterbeck's script: https://github.com/schlatterbeck/snxvpn

## Installing Pre compiled
Download the compiled binary from [Release Pages](https://github.com/francisco-anderson/snx-go/releases), 

```sh
# Dowload latest version for linux
$ wget -qO snxgo https://github.com/francisco-anderson/snx-go/releases/latest/download/snxgo-linux-amd64

# compare binary checksum output against version in release notes page (https://github.com/francisco-anderson/snx-go/releases)
# if checksum do not match, binary was not successfully downloaded
$ sha256sum snxgo

# make binary as executable
$ chmod +x snxgo

# move binary to binaries directory
$ sudo mv snxgo /usr/local/bin/snxgo

```

**Note for Mac users: On mac you need to go to security and privacy settings and allow this app to run. This is required due to the application not being signed.**

## Building from Sources

### Prerequisites
- go
- make

To build the project use the make command

```sh
make
```

To install the generated binary use the command:
```sh
sudo make install
```
or alternatively use `PREFIX` to define an alternative installation directory:
```sh
sudo make install PREFIX=/usr
```

## Using

To connect VPN:
```sh
snxgo --host somehost.test --realm ssl_vpn --user user
```

### Configuration Parameters

| Parâmetro           | Tipo        | Decrição                                |
|:--------------------|:-----------:|:----------------------------------------|
|    --host           |   required  | Connection host for VPN                 |
|    --user           |   required  | User for VPN connection                 |
|    --password       |   optional¹ | VPN connection password                 |
|    --realm          |   required  | Realm for VPN connection.               |
|    --skip-security  |   optional  | Disable SSL certificate checking        |
|    --debug          |   optional  | Enable debug logging                    |
|    --version        |   optional  | Display version and build date          |
|    --snx-path       |   optional  | Alternative SNX Path                    |

¹ - Will be asked if not informed

For more parameters and configuration using the parameter `--help`: