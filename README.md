# SOQL MCP Server

A Model Context Protocol (MCP) server that provides tools for querying Salesforce using SOQL (Salesforce Object Query Language).

## Features

- **SOQL Query Tool**: Execute SOQL queries against Salesforce and get results in JSON or table format

![SOQL MCP Server Demo](assets/images/soql-mcp.gif)

## Tools

### soql_query

Execute SOQL queries against Salesforce.

**Parameters:**

- `query` (required): The SOQL query to execute

**Example queries:**

```soql
SELECT Id, Name FROM Account LIMIT 10
SELECT Id, Name, Email FROM Contact WHERE Email != NULL LIMIT 5
SELECT Name, StageName, Amount FROM Opportunity WHERE StageName = 'Closed Won'
```

### describe

Describe Salesforce objects to get their metadata, fields, and properties.

**Parameters:**

- `object` (required): The Salesforce object name to describe (e.g., Account, Contact, Opportunity)
- `format` (optional): Output format: 'json' or 'table' (default: table)

**Example usage:**

```
object: Account
format: table
```

```
object: Contact
format: json
```

### debug

Return server configuration information for troubleshooting purposes.

**Parameters:**

No parameters required.

**Usage:**

This tool displays current server configuration including server name, version, resource path, debug mode status, and log level.

## Configuration

Set the following environment variables to configure the server:

### Server Configuration

```json
{
  "mcpServers": {
    "soql-mcp": {
      "command": "soql-mcp",
      "type": "stdio",
      "env": {
        "MCP_SERVER_NAME": "SOQL MCP Server",
        "MCP_SERVER_VERSION": "1.0.0",
        "MCP_RESOURCE_PATH": "/soql-mcp/terms.json",
        "MCP_DEBUG": "true",
        "MCP_LOG_LEVEL": "debug",
        "SALESFORCE_URL": "https://login.salesforce.com",
        "SALESFORCE_CLIENT_ID": "XXX",
        "SALESFORCE_CLIENT_SECRET": "XXX",
        "SALESFORCE_USERNAME": "XXX",
        "SALESFORCE_PASSWORD": "XXX",
        "SALESFORCE_SECURITY_TOKEN": "XXX"
      }
    }
  }
}
```

## Build

```bash
go build -o soql-mcp
```
