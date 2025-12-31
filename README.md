# Classical Music MCP Server

> **⚠️ Demo Project**
>
> This is **not in use** by [Bachstreet](https://discord.com/servers/bachstreet-707021383288488047).
> Created for **Moonshot.ai** technical evaluation only.

---

An MCP server connecting AI assistants to IMSLP's 600,000+ public domain classical music scores.

## What It Does

Enables AI assistants to answer questions like:
- *"Where can I find Bach's Prelude in C major?"* → Returns PDF download links
- *"What key is Mozart K545 in?"* → Returns C major, instrumentation, opus info
- *"I need the Moonlight Sonata sheet music"* → Searches, identifies, provides free PDFs


## Available Tools

### `search_work`
Search IMSLP by title, composer, or catalog number.
- **Input:** `query` (string), `limit` (number, optional)
- **Returns:** Matching works with titles, composers, URLs

### `get_work_details`
Get metadata for a specific work.
- **Input:** `page_title` (string, from search results)
- **Returns:** Instrumentation, key, opus number, composer info

### `get_score_links`
Get PDF download links.
- **Input:** `page_title` (string, from search results)
- **Returns:** Direct download URLs, file sizes, edition info

**Project Structure:**
```
client/     # IMSLP API client
models/     # Data structures
tools/      # MCP tool implementations
main.go     # Server setup
```

## Future Enhancements

- Filter by instrumentation, period, difficulty
- Composer biography tool
- Caching for performance
- Link to public domain recordings