# phpVMS ACARS Client

A Go application that replaces the existing Python PySide client for phpVMS. This application consumes UDP JSON payloads sent by the FlyWithLua script in X-Plane 12 and performs the complete PIREP workflow against a phpVMS backend.

## Features

- Listens for UDP JSON payloads from X-Plane 12
- Implements the complete phpVMS API client
- Manages the PIREP workflow (prefile, updates, file, cancel)
- Optional status interface via HTTP (disabled by default)
- Interactive Terminal User Interface (TUI) for monitoring and control
- Configurable via environment variables

## Requirements

- Go 1.24 or later
- X-Plane 12 with FlyWithLua script sending UDP JSON payloads
- phpVMS backend

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/phpvms-xplane.git
   cd phpvms-xplane
   ```

2. Build the application:
   ```
   make build
   ```

3. Run the application:
   ```
   make run
   ```

## Configuration

The application can be configured using environment variables:

| Environment Variable | Description | Default |
|----------------------|-------------|---------|
| PHPVMS_BASE_URL | Base URL of the phpVMS API (required) | |
| PHPVMS_API_KEY | API key for phpVMS authentication (required) | |
| UDP_BIND_HOST | Host to bind the UDP listener to | 0.0.0.0 |
| UDP_BIND_PORT | Port to bind the UDP listener to | 47777 |
| TUI_ENABLED | Enable Terminal User Interface | true |
| LOG_LEVEL | Log level (debug, info, warn, error) | info |

### Using Environment Variables

Example:
```
PHPVMS_BASE_URL=https://your-phpvms-instance.com \
PHPVMS_API_KEY=your-api-key \
UDP_BIND_HOST=0.0.0.0 \
UDP_BIND_PORT=47777 \
TUI_ENABLED=true \
LOG_LEVEL=info \
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
```

The application will automatically load the `.env` file if it exists. You can also specify a custom path to the `.env` file using the `-config` flag:

```
./build/phpvmsd -config /path/to/your/.env
```

For development, you can use the `dev` target in the Makefile:
```
make dev
```

## UDP Message Format

The application expects UDP messages in the following JSON format:

```json
{
  "status": "ENR",
  "position": {
    "lat": 40.6398,
    "lon": -73.7789,
    "altitude_msl": 12000.0,
    "altitude_agl": 500.0,
    "gs": 320.0,
    "sim_time": 1724167500,
    "distance": 217.4,
    "heading": 255.0,
    "ias": 250.0,
    "vs": -500.0
  },
  "fuel": 4520.0,
  "flight_time": 83.0,
  "events": [{"log":"Passing 10,000 ft","sim_time":1724167400}]
}
```

Field descriptions:
- `status`: PirepStatus code (e.g., INI, BST, TXI, TOF, ENR, ARR, PSD)
- `position.lat`: Latitude in degrees
- `position.lon`: Longitude in degrees
- `position.altitude_msl`: Altitude above mean sea level in feet
- `position.altitude_agl`: Altitude above ground level in feet
- `position.gs`: Ground speed in knots
- `position.sim_time`: Simulation time in Unix epoch seconds
- `position.distance`: Distance in nautical miles
- `position.heading`: Heading in degrees
- `position.ias`: Indicated airspeed in knots
- `position.vs`: Vertical speed in feet per minute
- `fuel`: Remaining fuel in kilograms
- `flight_time`: Flight time in minutes since block off
- `events`: Optional array of events with log messages and timestamps

## Usage

1. Start the application with the required configuration.
2. Configure X-Plane 12 with FlyWithLua to send UDP JSON payloads to the host and port where the application is running.
3. Start a flight in X-Plane 12.
4. The application will receive the UDP payloads and update the PIREP in phpVMS.
5. Interact with the application through the Terminal User Interface (if TUI_ENABLED is true).

### Terminal User Interface (TUI)

The application includes an interactive Terminal User Interface (TUI) that allows you to:

- Monitor UDP metrics in real-time (packets received, last sender, last status, etc.)
- Start a new flight by entering flight details (airline ID, flight number, aircraft ID, departure, arrival)
- File a completed flight
- Cancel an active flight
- Reset the active PIREP

#### TUI Keyboard Shortcuts

- `?`: Toggle help view
- `q` or `Ctrl+C`: Quit the application
- `s`: Start a new flight (when in flight input mode)
- `f`: File the active flight
- `c`: Cancel the active flight
- `r`: Reset the active PIREP
- `Tab`: Navigate to the next input field
- `Shift+Tab`: Navigate to the previous input field

You can disable the TUI by setting `TUI_ENABLED=false` in your environment or .env file, or by using the `-tui=false` command-line flag.

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

This project is licensed under the MIT License - see the LICENSE file for details.
