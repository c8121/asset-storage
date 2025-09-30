# asset-storage

Basically a DAM-System (Digital Asset Management). Includes a Webserver providing a SPA (Single Page Application) 
to browse the storage content.

## Goals

**DAM**: Manage Pictures, Documents...

**Portable**: Runs on Windows, Linux, Mac. No installation required, just copy. Can be copied to a external HD or Stick
which then can be used to store assets from different devices.

**Server included**: Can be used on a NAS for example.

**Deduplication**: Same file is stored automatically only once, no matter how often it was added.

**Resilience**: Files are stored in directories, one per time period, to enable quick backups even on large storages. Meta-Data stored separately in JSON-Format. Database can be recreated from Meta-Data and vice versa.

## Features

Store files into an archive directory

- Deduplicated (same file only once, determined by content-hash)
- Organized in subdirectories, one per time period. Old directories will no be touched again to enable incremental backups.
- Meta-data is stored separately and updated when new files are added or existing files with different origin/owner
- A database is created to be able to find/browse data

## Apps

### add

To add new files to the archive

    add [-gzip] [-maxmem <bytes>] [-base <directory>] [-r] <file or directory>

### spa-server

A HTTP Server which provides a Single-Page-Application to browse the storage

*Note: This is work in progress, important features like TLS and authentication are missing at the moment*

    spa-server [-gzip] [-xor <key>] [-base <directory>] [-spa <http-root-directory of spa-app>]

### rest-server

A HTTP Server which provides a REST-API to access the storage. This server is included in `spa-server`.

*Note: This is work in progress, important features like TLS and authentication are missing at the moment*

    rest-server [-gzip] [-xor <key>] [-base <directory>]

### metadata-db-create

Update meta-data-database by reading all meta-data files and writing contents to database.

Not required if database is intact, because `add` also updates the database.

    metadata-db-create [-base <directory>]

## App Commandline args


| Parameter         | Description                                                                                                                                                                                                                                                                                                                                     |
|-------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| base <dir&gt;     | Asset-storage base dir containing all data (file, meta-data, database). Default is `$HOME/asset-storage`                                                                                                                                                                                                                                        |
| gzip              | Use gzip to compress data.<br/> **Important:** Cannot be mixed, use always or never for one storage.                                                                                                                                                                                                                                            |
| maxmem <bytes&gt; | Max size in bytes when reading files while adding to storage. If a file is larger, it will not be read into memory and a temp-file will be used                                                                                                                                                                                                 |
| skip-meta         | When adding files: Skip updating meta-data if file exists.
| spa <dir&gt;      | HttpRoot-Directory which contains the SPA-files (HTML, JS, etc)                                                                                                                                                                                                                                                                                 |
| xor <key&gt;      | Content will be XOR'ed to obfusicate. This is to avoid manual changes to files (when content is XOR'ed, files cannot be openend and modified directly from storage directory) <br/>**Important:** Cannot be mixed, use always with same key or never for one storage. <br/>**Important:** Use same key for all apps with same storage directory |

## Libraries used

I have used these libraries (many thanks to the authors)

- SQLite: https://pkg.go.dev/modernc.org/sqlite
- Web Services: https://pkg.go.dev/github.com/gin-gonic/gin
- Webp decoding: https://github.com/HugoSmits86/nativewebp
- Image scaling: https://pkg.go.dev/golang.org/x/image/draw
- MIME type detection: https://pkg.go.dev/github.com/gabriel-vasile/mimetype
- Web UI: https://vuejs.org/
- CSS: https://getbootstrap.com/docs/5.0/getting-started/introduction/
