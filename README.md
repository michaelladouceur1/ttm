## Usage

`ttm [command] [subcommands]`

## Task

`ttm task [subcommands]`

### Main

- `ttm view`: Serve webpage displaying tasks and sessions in user-friendly UI

### Subcommands

- `add` : Add a task
  - `ttm task add [title] [description?]`
- `list` : List a task
  - `ttm task list`
- `update` : Update a task
  - `ttm task update [taskId]`
- `summary` : Summarize tasks for given time period
  - `ttm task summary`
- `close` : Close a task
  - `ttm task close [taskId]`

## Session

`ttm session [subcommands]`

### Subcommands

- `start` : Start a session
  - `ttm session start [taskId]`
- `end` : End current session
  - `ttm session end`
- `cancel` : Cancel current session (no save)
  - `ttm session cancel`
- `info` : Get info about current session
  - `ttm session info`
- `summary` : Summarize sessions for given time period
  - `ttm session summary`

## Terminal UI

Run `ttm` without arguments to open the full-screen terminal UI. Type `/` in
the command input to browse commands, use the arrow keys or Tab to select one,
and press Enter to run it.

- `/add`: add a task by entering its title, description, priority, and tags
- `/tasks ["query"]`: list open tasks, optionally filtered by a query
- `/tags`: list all tags and their task counts
- `/search`: search task text or use `*field:value` filters
- `/summary [TIME_PERIOD_IN_DAYS]`: show sessions by day and task totals (defaults to 7 days)
- `/note [task_id]`: add a note to a task, or the active session's task when omitted
- `/notes [task_id]`: list a task's note contents, or the active session's task notes when omitted
- `/start <task_id>`: start a session for a listed task
- `/open <task_id>...`: mark listed tasks as open
- `/standby <task_id>...`: put listed tasks on standby
- `/end`: end and save the active session
- `/cancel`: discard the active session

## Windows installation

Windows releases include `ttm_setup.exe`. Run the installer to install TTM in
`Program Files\TTM`, add Start menu shortcuts, and register an uninstaller.
To create the installer from source on Windows, install [NSIS](https://nsis.sourceforge.io/)
and run:

```sh
make package-windows
```

The installer is written to `dist/ttm_setup.exe`.

## Shorthand

- `ttm add`: Add a task
- `ttm list`: List tasks
- `ttm search <query>`: Search task fields. Use field filters such as `"$tags:work,urgent $title:Task 1"`.
- `ttm close`: Close task
- `ttm start`: Start session
- `ttm pause`: Pause session
- `ttm end`: End session
- `ttm cancel`: Cancel session

## Storage configuration

TTM uses SQLite by default. To store task and session data in a Google Doc instead,
set the following values in `~/.ttm/config.yaml`:

```yaml
storage:
  type: google-docs
  googleDocs:
    documentId: "your-google-document-id"
    credentialsFile: "/absolute/path/to/service-account.json"
```

Create a Google Cloud service account, enable the Google Docs API, download its
JSON credential file, and share the target Google Doc with the service account's
email address as an editor. The document content is managed as JSON by TTM; do not
edit that JSON manually. Credentials remain in the credential file rather than
being stored in the document or committed to source control.

To use PostgreSQL, configure its connection settings:

```yaml
storage:
  type: postgres
  postgres:
    host: localhost
    port: 5432
    user: ttm
    dbname: ttmdb
    sslmode: disable
    passwordEnv: TTM_POSTGRES_PASSWORD
```

`passwordEnv` names an environment variable containing the password, so the
password does not need to be stored in `config.yaml`. When `passwordEnv` and
`password` are omitted, the PostgreSQL driver also supports `PGPASSWORD` and a
matching `~/.pgpass` (or `PGPASSFILE`) entry. `password` is available only for
setups where storing it in the plaintext configuration file is acceptable.

## Logging themes

Set `logging.theme` in `~/.ttm/config.yaml` to choose the output format. The
default `classic` theme uses colored tables and summary trees. `compact` is a
color-free, line-oriented format for narrow terminals and scripting.

```yaml
logging:
  theme: compact
```

## PostgreSQL over Tailscale setup

Use this when you want TTM on one machine to use PostgreSQL on another machine over a private Tailscale network.

### 1. Install and connect both machines to Tailscale
- Install Tailscale on:
  - your TTM client machine
  - your PostgreSQL server machine
- Ensure both are in the same tailnet and can ping each other over Tailscale.
- Note the PostgreSQL server Tailscale address (example: 100.88.77.66).

### 2. Configure PostgreSQL to listen for remote connections
On the PostgreSQL server, edit postgresql.conf:

    listen_addresses = '127.0.0.1,100.88.77.66'
    port = 5432

You can use * for listen_addresses, but binding only loopback + Tailscale IP is tighter.

Restart PostgreSQL and confirm it is listening:

    sudo systemctl restart postgresql
    ss -lntp | rg 5432

You should see 5432 bound to the Tailscale address (or 0.0.0.0 if using *), not only 127.0.0.1.

### 3. Create app user and database
TTM initializes tables, but the database itself must already exist.

    sudo -u postgres psql

    CREATE ROLE ttm_app LOGIN PASSWORD 'replace-with-strong-password';
    CREATE DATABASE ttmdb OWNER ttm_app;
    GRANT ALL PRIVILEGES ON DATABASE ttmdb TO ttm_app;

    \q

### 4. Allow Tailscale clients in pg_hba.conf
Add one of these rules on the PostgreSQL server.

Tighter (single client):

    hostssl  ttmdb  ttm_app  100.77.66.55/32  scram-sha-256

Broader (all Tailscale IPv4 nodes):

    hostssl  ttmdb  ttm_app  100.64.0.0/10    scram-sha-256

If you enable PostgreSQL TLS, use hostssl instead of hostnossl (or always use hostssl as shown above).

Reload PostgreSQL:

    sudo systemctl reload postgresql

### 5. Allow firewall access on the DB server
Allow inbound TCP 5432 from:
- the client Tailscale IP (/32), or
- the Tailscale CGNAT range 100.64.0.0/10

Do not expose 5432 to the public internet.

### 6. Configure TTM
In ~/.ttm/config.yaml:

    storage:
      type: postgres
      postgres:
        host: "100.88.77.66"
        port: 5432
        user: "ttm_app"
        dbname: "ttmdb"
        sslmode: require
        passwordEnv: "TTM_POSTGRES_PASSWORD"

Set password in environment (recommended over plaintext config):

    export TTM_POSTGRES_PASSWORD='replace-with-strong-password'

### 7. Validate connectivity before running TTM
From the client machine:

    pg_isready -h 100.88.77.66 -p 5432 -U ttm_app -d ttmdb
    psql -h 100.88.77.66 -U ttm_app -d ttmdb -c "select 1;"

If these work, TTM should connect and initialize schema objects automatically.

### Common errors and fixes

- Connection refused
  - PostgreSQL not listening on Tailscale interface, or firewall is blocking 5432.

- No pg_hba.conf entry for host ... no encryption
  - Add matching hostnossl rule for client IP/range, or switch to hostssl + sslmode require.

- Database does not exist
  - Create the database first (TTM does not create the DB itself).

### Security recommendations

1. Use a dedicated least-privilege DB user (not postgres superuser).
2. Prefer passwordEnv instead of storing password in config file.
3. Rotate credentials if they were ever exposed.
4. Restrict pg_hba and firewall rules to least privilege.
5. Consider TLS (sslmode require) in addition to Tailscale encryption.

## TODO

- [ ] TUI commands
  - [x] /add (started; matching ability to CLI)
  - [x] /tag - add tag to task
  - [x] /note - add note to task
  - [x] /notes - list notes for task (or current session's task if no task specified)
  - [x] /start - start session for task
  - [ ] /pause - pause session for task
  - [x] /end - end session for task
  - [x] /cancel - cancel session for task
  - [x] /tasks (started; matches ability to CLI but needs to only list with filter for priority, status, tags, etc.)
  - [x] /summary - summary of sessions and tasks worked on in a given time period
  - [x] /tags - list all tags and their counts
  - [x] /search - search task fields
  - [ ] /update - update task fields
  - [x] /detail - shows task details with notes and sessions
  - [x] /close - close task
  - [x] /open - open task
  - [x] /standby - put task on standby
  - [ ] /sessions - list sessions for given time period
  - [ ] /users - list users and their task/session counts
  - [ ] /add-user - add user to multi-user environment
  - [ ] /assign - assign task to user
  - [ ] /csv - export tasks and sessions to CSV
- [ ] Commands legend and help
- [x] Task notes
- [ ] User support for multi-user environments
- [ ] Command history (e.g., up arrow to repeat last command)
- [ ] Focus (Analogous to sprint) for tasks and sessions
- [ ] Time-based reports viewing tasks and sessions
- [ ] Integration with copilot cli; send tasks to copilot agent
- [ ] Integration with project management tools (e.g., Jira, Trello)
