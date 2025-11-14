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

### TUI Mode

#### Command List

- Add task
- Start session
- List tasks

### Shorthand

- `ttm add`: Add a task
- `ttm list`: List tasks
- `ttm close`: Close task
- `ttm start`: Start session
- `ttm pause`: Pause session
- `ttm end`: End session
- `ttm cancel`: Cancel session

#### TODO

- [ ] TUI mode
- [ ] Notes support per task/session
- [ ] User support for multi-user environments
- [ ] Time-based reports viewing tasks and sessions
- [ ] Integration with project management tools (e.g., Jira, Trello)
