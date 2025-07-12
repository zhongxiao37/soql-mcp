# SOQL MCP Server

A Model Context Protocol (MCP) server that provides tools for querying Salesforce using SOQL (Salesforce Object Query Language).

## Features

- **SOQL Query Tool**: Execute SOQL queries against Salesforce and get results in JSON or table format
- **Debug Tool**: Get server configuration information
- **Hello Tool**: Simple greeting tool for testing
- **Terms Resource**: Access to terms and definitions

## Tools

### soql_query

Execute SOQL queries against Salesforce.

**Parameters:**

- `query` (required): The SOQL query to execute
- `format` (optional): Output format - 'json' or 'table' (default: 'json')

**Example queries:**

```soql
SELECT Id, Name FROM Account LIMIT 10
SELECT Id, Name, Email FROM Contact WHERE Email != NULL LIMIT 5
SELECT Name, StageName, Amount FROM Opportunity WHERE StageName = 'Closed Won'
```

### debug

Returns server configuration information.

### hello_world

Simple greeting tool for testing.

**Parameters:**

- `name` (required): Name of the person to greet

## Configuration

Set the following environment variables to configure the server:

### Server Configuration

- `MCP_SERVER_NAME`: Server name (default: "Demo 🚀")
- `MCP_SERVER_VERSION`: Server version (default: "1.0.0")
- `MCP_RESOURCE_PATH`: Path to resource files (required)
- `MCP_DEBUG`: Enable debug mode (default: false)
- `MCP_LOG_LEVEL`: Log level (default: "info")

### Salesforce Configuration

- `SALESFORCE_URL`: Salesforce login URL (default: "https://login.salesforce.com")
- `SALESFORCE_CLIENT_ID`: Connected App Client ID (required)
- `SALESFORCE_CLIENT_SECRET`: Connected App Client Secret (required)
- `SALESFORCE_USERNAME`: Salesforce username (required)
- `SALESFORCE_PASSWORD`: Salesforce password (required)
- `SALESFORCE_SECURITY_TOKEN`: Salesforce security token (optional, but usually required)

## Salesforce Setup

To use the SOQL query tool, you need to:

1. **Create a Connected App in Salesforce:**

   - Go to Setup → App Manager → New Connected App
   - Enable OAuth Settings
   - Add required OAuth scopes (at minimum: "Access and manage your data (api)")
   - Set callback URL (can be a placeholder like `https://localhost`)

2. **Get your Security Token:**

   - Go to Settings → Reset My Security Token
   - The token will be sent to your email

3. **Set Environment Variables:**
   ```bash
   export SALESFORCE_CLIENT_ID="your_client_id"
   export SALESFORCE_CLIENT_SECRET="your_client_secret"
   export SALESFORCE_USERNAME="your_username"
   export SALESFORCE_PASSWORD="your_password"
   export SALESFORCE_SECURITY_TOKEN="your_security_token"
   export MCP_RESOURCE_PATH="/path/to/your/terms.json"
   ```

## Usage

1. Build the server:

   ```bash
   go build -o soql-mcp
   ```

2. Run the server:

   ```bash
   ./soql-mcp
   ```

3. Use with an MCP client to execute SOQL queries:
   ```json
   {
     "tool": "soql_query",
     "arguments": {
       "query": "SELECT Id, Name FROM Account LIMIT 5",
       "format": "table"
     }
   }
   ```

## Output Formats

### JSON Format

Returns the raw Salesforce API response in JSON format, including metadata.

### Table Format

Returns a human-readable table format with:

- Total record count
- Records returned count
- Each record displayed as key-value pairs

## Error Handling

The tool provides detailed error messages for:

- Missing or invalid configuration
- Salesforce authentication failures
- SOQL syntax errors
- Network connectivity issues

## Development

To extend this server:

1. Add new tools in `pkg/tools/`
2. Register tools in `main.go`
3. Update configuration in `pkg/config.go` if needed
4. Add resources in `pkg/resources/` if needed

## License

This project is licensed under the MIT License.
