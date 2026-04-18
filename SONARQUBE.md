# SonarQube Configuration

This project uses SonarQube for code quality analysis.

## Setup

To enable SonarQube scanning, you need to configure the following GitHub/Gitea secrets:

### Secrets Required

- `SONAR_HOST_URL`: Your SonarQube server URL (e.g., `https://sonarqube.example.com`)
- `SONAR_LOGIN`: Your SonarQube authentication token

### Adding Secrets in Gitea

1. Go to your repository settings
2. Navigate to "Secrets and variables" → "Actions"
3. Add the following secrets:
   - `SONAR_HOST_URL`
   - `SONAR_LOGIN`

## SonarQube Properties

The `sonar-project.properties` file configures:

- **Project Key**: `epage-go`
- **Sources**: `src/` directory
- **Tests**: All `*_test.go` files in `src/`
- **Coverage**: Reports from `coverage.out`
- **Exclusions**: Test files excluded from main analysis

## Workflow Triggers

- **CI Workflow** (`.gitea/workflows/ci.yml`):
  - Runs on: Pull requests to `main` or `develop`
  - Tests, lints, and builds the application

- **SonarQube Workflow** (`.gitea/workflows/sonarqube.yml`):
  - Runs on: Pushes to `main` branch
  - Generates coverage reports and uploads to SonarQube

## Local SonarQube Scanning

To scan locally (requires SonarQube CLI):

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./src

# Run SonarQube scanner
sonar-scanner \
  -Dsonar.projectKey=epage-go \
  -Dsonar.sources=src \
  -Dsonar.host.url=https://sonarqube.example.com \
  -Dsonar.login=your_token
```

## Coverage Reports

Test coverage is automatically collected and can be viewed:

```bash
# Run tests with coverage
go test -v -coverprofile=coverage.out ./src

# View coverage report
go tool cover -html=coverage.out
```
