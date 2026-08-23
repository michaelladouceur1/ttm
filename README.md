## Usage

`ttm [command] [subcommands]`

### Task

`ttm task [subcommands]`

#### Main

- `ttm view`: Serve webpage displaying tasks and sessions in user-friendly UI

#### Subcommands

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

### Session

`ttm session [subcommands]`

#### Subcommands

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

### Terminal UI

Run `ttm` without arguments to open the full-screen terminal UI. Type `/` in
the command input to browse commands, use the arrow keys or Tab to select one,
and press Enter to run it.

- `/add`: add a task by entering its title, description, priority, and tags
- `/tasks ["query"]`: list open tasks, optionally filtered by a query
- `/tags`: list all tags and their task counts
- `/search`: search task text or use `*field:value` filters
- `/summary [TIME_PERIOD_IN_DAYS]`: show sessions by day and task totals (defaults to 7 days)
- `/note [task_id]`: add a note to a task, or the active session's task when omitted
- `/start <task_id>`: start a session for a listed task
- `/end`: end and save the active session
- `/cancel`: discard the active session

### Windows installation

Windows releases include `ttm_setup.exe`. Run the installer to install TTM in
`Program Files\TTM`, add Start menu shortcuts, and register an uninstaller.
To create the installer from source on Windows, install [NSIS](https://nsis.sourceforge.io/)
and run:

```sh
make package-windows
```

The installer is written to `dist/ttm_setup.exe`.

### Shorthand

- `ttm add`: Add a task
- `ttm list`: List tasks
- `ttm search <query>`: Search task fields. Use field filters such as `"$tags:work,urgent $title:Task 1"`.
- `ttm close`: Close task
- `ttm start`: Start session
- `ttm pause`: Pause session
- `ttm end`: End session
- `ttm cancel`: Cancel session

### Storage configuration

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

### Logging themes

Set `logging.theme` in `~/.ttm/config.yaml` to choose the output format. The
default `classic` theme uses colored tables and summary trees. `compact` is a
color-free, line-oriented format for narrow terminals and scripting.

```yaml
logging:
  theme: compact
```

#### TODO

- [ ] TUI commands
  - [x] /add (started; matching ability to CLI)
  - [x] /tag - add tag to task
  - [x] /note - add note to task
  - [x] /start - start session for task
  - [ ] /pause - pause session for task
  - [x] /end - end session for task
  - [x] /cancel - cancel session for task
  - [x] /tasks (started; matches ability to CLI but needs to only list with filter for priority, status, tags, etc.)
  - [x] /summary - summary of sessions and tasks worked on in a given time period
  - [x] /tags - list all tags and their counts
  - [x] /search - search task fields
  - [ ] /update - update task fields
  - [ ] /detail - shows task details with notes and sessions
  - [x] /close - close task
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
