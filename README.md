# asset-storage

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

    spa-server [-gzip] [-base <directory>] [-spa <http-root-directory of spa-app>]

### rest-server

A HTTP Server which provides a REST-API to access the storage. This server is included in `spa-server`.

*Note: This is work in progress, important features like TLS and authentication are missing at the moment*

    rest-server [-gzip] [-base <directory>]

### metadata-db-create

Update meta-data-database by reading all meta-data files and writing contents to database.

Not required if database is intact, because `add` also updates the database.

    metadata-db-create [-base <directory>]

## App Commandline args


Parameter     | Description
--------------| -----------
base          | Asset-storage base dir containing all data (file, meta-data, database). Default is `$HOME/asset-storage`
gzip          | Use gzip to compress data. **Important:** Cannot be mixed, use always or never for one storage.
maxmem        | Max size in bytes when reading files while adding to storage. If a file is larger, it will not be read into memory and a temp-file will be used
spa           | HttpRoot-Directory which contains the SPA-files (HTML, JS, etc)