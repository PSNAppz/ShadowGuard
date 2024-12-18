# ShadowGuard

"ShadowGuard" is a robust and flexible security engine for web applications and APIs, providing a wide array of features to safeguard your services. Developed in Go, "ShadowGuard" acts as a protective layer, analyzing and filtering incoming HTTP/HTTPS requests to ensure that only legitimate traffic gets through.

## Key Features:
- Traffic Filtering: Detects and blocks common web-based attacks like SQL Injection, Cross-Site Scripting (XSS), Cross-Site Request Forgery (CSRF).
- Rate Limiting: Prevents Denial-of-Service (DoS) and brute-force attacks by tracking and limiting requests from individual IP addresses.
- IP Blacklisting/Whitelisting: Allows blocking or permitting requests from specific IP addresses.
- Content Security Policy Enforcement: Ensures server responses adhere to a specified content security policy.
- Anomaly Detection: Learns normal request patterns and alerts when anomalous patterns occur.
- Logging and Alerting: Logs all requests and responses, providing real-time alerts for identified threats.
- Integration with Threat Intelligence Platforms: Stays updated with the latest threats and vulnerabilities.
- API Schema Validation: Validates request payloads against predefined schemas.

## Modes of Operation:
"ShadowGuard" can operate in two modes: Active and Passive. In ![Active Mode](https://img.shields.io/badge/ACTIVE_MODE-FF0000), it actively filters incoming requests and blocks malicious traffic. In ![Passive Mode](https://img.shields.io/badge/PASSIVE_MODE-8A2BE2), it only logs requests and responses, and alerts when it detects a threat. Active mode is recommended for APIs which are critical to the functioning of the application.

## Architecture:
"ShadowGuard" employs a modular architecture, built around the Factory Design Pattern, allowing for robust scalability and reliability. Designed as a middleware, it integrates seamlessly with any application, independent of the backend technology stack. It employs a plugin-based system where each plugin handles a specific type of security task. This plugin-based approach provides an extensible framework where new plugins can be added to enhance the security features without disrupting the existing system.

The behavior of each plugin is driven by configuration, making "ShadowGuard" highly adaptable to the needs of different applications. This configurability enables fine-tuning of the security parameters to match the specific requirements and threat profiles of each application.

The architecture also facilitates both active and passive modes of operation, allowing the system to either block malicious traffic actively or to monitor and alert on potential threats passively. This flexibility of operation modes allows "ShadowGuard" to be tailored to the specific security posture of your application or API.

## Getting Started:
TODO: Instructions on how to setup "ShadowGuard", its dependencies, and how to get it running.

### Database
Run `build.sh` to setup install the necessary dependencies including Postgresql and configures the gorm database.

This script does a couple of things in this sequence.
1. Installs all of the necessary dependencies.
2. Starts the Postgres service.
3. Runs the necessary SQL commands to create the user and database needed for gorm.
4. Configures Postgres to allow md5 connections.

```shell
sudo bash build.sh
```

## How to Use:
The program can be ran in a multitde of ways.

### Go
```shell
go run cmd/main.go
```

### Shell
```shell
chmod +x run.sh
./run.sh
```

### Docker
```shell
docker build . -t shadow_guard 
docker run --network=host shadow_guard
```

## Unit Tests:

In order to run unit tests, you can use the shell script `run_tests.sh` in the root directory. The unit tests can be ran using convential Go commands.

## Documentation:
TODO: Link to full API documentation, or brief outline of main methods and how to use them.

## License:
TODO: Information on the licensing of "ShadowGuard".
