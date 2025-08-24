# phpvms-xplane (PXP)

While [phpVMS](https://www.phpvms.net) integrates directly with
[vmsACARS](https://docs.phpvms.net/acars/overview), phpvms-xplane (PXP
for short) offers an alternative for those who can't (or don't wish to)
run it, such as Windows users. It connects to X-Plane 12 via a provided
FlyWithLua script.

It's designed to be:
- fast
- small
- simple
- portable
- observable.

This is a work of love and not of profit; it doesn't come with promises,
but I will update it as long as I use it!

## Features

- Listens for UDP JSON payloads from X-Plane 12
- Implements the complete phpVMS API client
- Manages the PIREP workflow (prefile, updates, file, cancel)
- Interactive Terminal User Interface (TUI) for monitoring and control
- Configurable via environment variables

## Requirements

- X-Plane 12 with FlyWithLua script sending UDP JSON payloads
- phpVMS backend

## Usage

(something good planned for this section!)

## Configuration

The application can be configured using environment variables:

| Environment Variable | Description                                  | Default |
|----------------------|----------------------------------------------|---------|
| PHPVMS_BASE_URL      | Base URL of the phpVMS API (required)        |         |
| PHPVMS_API_KEY       | API key for phpVMS authentication (required) |         |
| UDP_BIND_HOST        | Host to bind the UDP listener to             | 0.0.0.0 |
| UDP_BIND_PORT        | Port to bind the UDP listener to             | 47777   |
| TUI_ENABLED          | Enable Terminal User Interface               | true    |
| LOG_LEVEL            | Log level (debug, info, warn, error)         | info    |
| SIMBRIEF_USER_ID     | The pilot's numeric SimBrief ID              |         |

### Using Environment Variables

Example:
```
PHPVMS_BASE_URL=https://your-phpvms-instance.com \
PHPVMS_API_KEY=your-api-key \
UDP_BIND_HOST=0.0.0.0 \
UDP_BIND_PORT=47777 \
TUI_ENABLED=true \
LOG_LEVEL=info \
SIMBRIEF_USER_ID=123456 \
./build/phpvmsd
```

### Using a .env File

You can also use a `.env` file to set environment variables. Create a file named `.env` in the same directory as the application with the following content:

```
PHPVMS_BASE_URL=https://your-phpvms-instance.com
PHPVMS_API_KEY=your-api-key
UDP_BIND_HOST=0.0.0.0
UDP_BIND_PORT=47777
TUI_ENABLED=true
LOG_LEVEL=info
SIMBRIEF_USER_ID=123456
```

The application will automatically load the `.env` file if it exists. You can also specify a custom path to the `.env` file using the `-config` flag:

```
./build/phpvmsd -config /path/to/your/.env
```

For development, you can use the `dev` target in the Makefile:
```
make dev
```

## Usage

1. Start the application with the required configuration.
2. Configure X-Plane 12 with FlyWithLua to send UDP JSON payloads to the host and port where the application is running.
3. Start a flight in X-Plane 12.
4. The application will receive the UDP payloads and update the PIREP in phpVMS.
5. Interact with the application through the Terminal User Interface (if TUI_ENABLED is true).

## Development

To run tests:
```
make test
```

To clean build artifacts:
```
make clean
```

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](./LICENSE) for details.
