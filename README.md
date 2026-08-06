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

- `/add "title" ["description"]`: add a task
- `/list ["query"]`: list open tasks, optionally filtered by a query

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

- [ ] Notes support per task/session
- [ ] User support for multi-user environments
- [ ] Time-based reports viewing tasks and sessions
- [ ] Integration with project management tools (e.g., Jira, Trello)
