# Resources

The server exposes static and templated MCP resources to help clients inspect capabilities before making tool calls.

## Static Resources

| Name | URI | Purpose |
| --- | --- | --- |
| `server_overview` | `scholar://server/overview` | Server summary, version, transport, and enabled capabilities. |
| `tool_catalog` | `scholar://server/tools` | Tool names, arguments, and return shapes. |
| `runtime_config` | `scholar://server/config` | Runtime config values such as timeout and rate limit. |
| `limitations` | `scholar://server/limitations` | Operational boundaries and caveats. |

## Resource Template

| Name | URI Template | Purpose |
| --- | --- | --- |
| `search_guide` | `scholar://search-guide/{topic}` | Topic-specific guidance for how to search and interpret Scholar results. |

## Example

Example resource read:

```text
scholar://search-guide/graph%20neural%20networks
```
