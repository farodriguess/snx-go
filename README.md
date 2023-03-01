# snxgo

snxgo is a utility to connect to Checkpoint VPN without needing a browser. This project is based on Ralf Schlatterbeck's script: https://github.com/schlatterbeck/snxvpn

## Installing
Download the compiled binary from [Release Pages](https://github.com/francisco-anderson/snx-go/releases) or compile the binary

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

¹ - Will be asked if not informed

For more parameters and configuration using the parameter `--help`: