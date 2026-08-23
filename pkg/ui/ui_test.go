package ui

import (
	"reflect"
	"strings"
	"testing"
	"time"
	"ttm/pkg/config"
	"ttm/pkg/fs"
	"ttm/pkg/models"
	"ttm/pkg/store"
	sqlite "ttm/pkg/store/sqlite"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

func TestParseCommand(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{input: `/list`, want: []string{"list"}},
		{input: `/add "Plan release" "Prepare release notes"`, want: []string{"add", "Plan release", "Prepare release notes"}},
		{input: `/add 'Plan release'`, want: []string{"add", "Plan release"}},
	}

	for _, test := range tests {
		got, err := parseCommand(test.input)
		if err != nil {
			t.Fatalf("parseCommand(%q) returned an error: %v", test.input, err)
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("parseCommand(%q) = %#v, want %#v", test.input, got, test.want)
		}
	}
}

func TestParseCommandRejectsInvalidInput(t *testing.T) {
	for _, input := range []string{"add task", `/add "unfinished`} {
		if _, err := parseCommand(input); err == nil {
			t.Errorf("parseCommand(%q) returned nil error", input)
		}
	}
}

func TestUpdateSuggestions(t *testing.T) {
	m := newModel(nil, nil)
	m.input.SetValue("/")
	m.updateSuggestions()
	if len(m.suggestions) != len(commands) {
		t.Fatalf("suggestions = %d, want %d", len(m.suggestions), len(commands))
	}

	m.input.SetValue("/tag")
	m.updateSuggestions()
	if len(m.suggestions) != 2 || m.suggestions[0].name != "tag" || m.suggestions[1].name != "tags" {
		t.Fatalf("suggestions = %#v, want tag and tags", m.suggestions)
	}
}

func TestStartCommandRequiresTaskID(t *testing.T) {
	m := newModel(nil, nil)
	m.input.SetValue("/start")
	m.execute()

	if m.content != "Usage: /start <task_id>" {
		t.Errorf("start content = %q", m.content)
	}
}

func TestSummaryDays(t *testing.T) {
	tests := []struct {
		args    []string
		want    int
		wantErr string
	}{
		{want: 7},
		{args: []string{"14"}, want: 14},
		{args: []string{"0"}, wantErr: "positive whole number"},
		{args: []string{"week"}, wantErr: "positive whole number"},
		{args: []string{"7", "14"}, wantErr: "usage:"},
	}

	for _, test := range tests {
		got, err := summaryDays(test.args)
		if test.wantErr != "" {
			if err == nil || !strings.Contains(err.Error(), test.wantErr) {
				t.Errorf("summaryDays(%v) error = %v, want %q", test.args, err, test.wantErr)
			}
			continue
		}
		if err != nil || got != test.want {
			t.Errorf("summaryDays(%v) = %d, %v; want %d, nil", test.args, got, err, test.want)
		}
	}
}

func TestRenderSummaryGroupsSessionsByDayAndTotalsTasks(t *testing.T) {
	now := time.Date(2026, time.August, 22, 15, 0, 0, 0, time.UTC)
	plan := models.Task{ID: 1, Title: "Plan release"}
	tests := models.Task{ID: 2, Title: "Write tests"}
	rendered := renderSummary([]summarySession{
		{session: models.Session{TaskId: 1, StartTime: now.Add(-2 * time.Hour), EndTime: now.Add(-75 * time.Minute)}, task: plan},
		{session: models.Session{TaskId: 2, StartTime: now.AddDate(0, 0, -1), EndTime: now.AddDate(0, 0, -1).Add(45 * time.Minute)}, task: tests},
		{session: models.Session{TaskId: 1, StartTime: now.AddDate(0, 0, -1).Add(time.Hour), EndTime: now.AddDate(0, 0, -1).Add(95 * time.Minute)}, task: plan},
	}, []summaryTask{
		{task: plan, duration: 80 * time.Minute},
		{task: tests, duration: 45 * time.Minute},
	}, 120)

	for _, want := range []string{
		"Sessions",
		"Tasks",
		"Saturday, August 22, 2026",
		"Friday, August 21, 2026",
		"Plan release",
		"Write tests",
		"01:20",
		"00:45",
	} {
		if !strings.Contains(rendered, want) {
			t.Errorf("renderSummary() = %q, want it to contain %q", rendered, want)
		}
	}
	if strings.Index(rendered, "Saturday, August 22, 2026") > strings.Index(rendered, "Friday, August 21, 2026") {
		t.Errorf("renderSummary() does not order days newest first: %q", rendered)
	}
}

func TestSummaryCommandCreatesChildModel(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	st := store.NewStore(sqlite.NewStore())
	if err := st.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	m := newModel(nil, st)
	m.input.SetValue("/summary 14")
	m.execute()

	summary, ok := m.active.(summaryModel)
	if !ok {
		t.Fatalf("active model = %T, want summaryModel", m.active)
	}
	if summary.days != 14 {
		t.Errorf("summary days = %d, want 14", summary.days)
	}
}

func TestSearchCommandCreatesChildModel(t *testing.T) {
	m := newModel(nil, nil)
	m.input.SetValue("/search")
	m.execute()

	if _, ok := m.active.(searchModel); !ok {
		t.Errorf("active model = %T, want searchModel", m.active)
	}
}

func TestSearchModelUpdatesMatchingTasks(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	st := store.NewStore(sqlite.NewStore())
	if err := st.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	if err := st.InsertTask(models.Task{Title: "Plan release", Status: models.StatusOpen}); err != nil {
		t.Fatalf("InsertTask() error = %v", err)
	}
	if err := st.InsertTask(models.Task{Title: "Write tests", Status: models.StatusOpen}); err != nil {
		t.Fatalf("InsertTask() error = %v", err)
	}

	m := newSearchModel(st, 80)
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("release")})
	m = updated.(searchModel)
	if !strings.Contains(m.content, "Plan release") || strings.Contains(m.content, "Write tests") {
		t.Errorf("search content = %q", m.content)
	}

	taskID, err := fs.GetTaskIDFromTempID(1)
	if err != nil {
		t.Fatalf("GetTaskIDFromTempID() error = %v", err)
	}
	if taskID != 1 {
		t.Errorf("mapped task ID = %d, want 1", taskID)
	}
}

func TestSearchModelClosesToRoot(t *testing.T) {
	m := searchModel{content: "Search results"}
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	msg := cmd()

	if _, ok := msg.(searchClosedMsg); !ok {
		t.Fatalf("close message = %T, want searchClosedMsg", msg)
	}
}

func TestStartSessionRejectsInvalidTaskID(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	content := startSession(nil, "not-a-task-id", time.Time{})
	if !strings.Contains(content, "Error starting session:") {
		t.Errorf("startSession() = %q, want an invalid ID error", content)
	}
}

func TestStartSessionCreatesSessionForMappedTask(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	st := store.NewStore(sqlite.NewStore())
	if err := st.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	if err := st.InsertTask(models.Task{Title: "Write tests"}); err != nil {
		t.Fatalf("InsertTask() error = %v", err)
	}
	if err := fs.UpdateIDMapFile([]models.Task{{ID: 1, ListID: 1}}); err != nil {
		t.Fatalf("UpdateIDMapFile() error = %v", err)
	}

	startedAt := time.Date(2026, time.August, 14, 23, 55, 0, 0, time.UTC)
	content := startSession(st, "1", startedAt)
	if !strings.Contains(content, "Session started.") {
		t.Fatalf("startSession() = %q", content)
	}

	session, err := fs.ReadSessionFile()
	if err != nil {
		t.Fatalf("ReadSessionFile() error = %v", err)
	}
	if session.TaskID != 1 || !session.StartTime.Equal(startedAt) {
		t.Errorf("session = %#v, want task 1 at %v", session, startedAt)
	}
}

func TestEndSessionSavesAndRemovesActiveSession(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	st := store.NewStore(sqlite.NewStore())
	if err := st.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	if err := st.InsertTask(models.Task{Title: "Write tests"}); err != nil {
		t.Fatalf("InsertTask() error = %v", err)
	}

	startedAt := time.Date(2026, time.August, 14, 23, 0, 0, 0, time.UTC)
	endedAt := startedAt.Add(25 * time.Minute)
	if err := fs.CreateSessionFile(1, startedAt); err != nil {
		t.Fatalf("CreateSessionFile() error = %v", err)
	}

	content := endSession(st, endedAt)
	if !strings.Contains(content, "Session ended.") {
		t.Fatalf("endSession() = %q", content)
	}
	if fs.SessionFileExists() {
		t.Error("session file still exists after ending the session")
	}

	sessions, err := st.GetSessionsByTaskID(1)
	if err != nil {
		t.Fatalf("GetSessionsByTaskID() error = %v", err)
	}
	if len(sessions) != 1 || !sessions[0].StartTime.Equal(startedAt) || !sessions[0].EndTime.Equal(endedAt) {
		t.Errorf("sessions = %#v, want one session from %v to %v", sessions, startedAt, endedAt)
	}
}

func TestCancelSessionRemovesActiveSessionWithoutSaving(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	if err := fs.CreateSessionFile(1, time.Now()); err != nil {
		t.Fatalf("CreateSessionFile() error = %v", err)
	}

	if content := cancelSession(); content != "Session cancelled." {
		t.Errorf("cancelSession() = %q", content)
	}
	if fs.SessionFileExists() {
		t.Error("session file still exists after cancelling the session")
	}
}

func TestAddFormWalksThroughFields(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	st := store.NewStore(sqlite.NewStore())
	if err := st.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	m := newModel(&config.Config{
		AddFlags: config.ConfigDefaultFlags{
			Category: string(models.CategoryTask),
			Priority: string(models.PriorityHigh),
			Status:   string(models.StatusOpen),
		},
	}, st)
	m.input.SetValue("/add")
	m.execute()

	add, ok := m.active.(addModel)
	if !ok {
		t.Fatalf("active model = %T, want addModel", m.active)
	}
	if add.step != addStepTitle || add.input.Placeholder != "Enter title..." {
		t.Fatalf("add form did not start at title: step=%d placeholder=%q", add.step, add.input.Placeholder)
	}

	add.input.SetValue("Plan release")
	updated, _ := add.submitAddField()
	add = updated.(addModel)
	if add.step != addStepDescription || add.draft.Title != "Plan release" {
		t.Fatalf("title was not saved: step=%d title=%q", add.step, add.draft.Title)
	}

	add.input.SetValue("Prepare notes")
	updated, _ = add.submitAddField()
	add = updated.(addModel)
	if add.step != addStepPriority || add.draft.Description != "Prepare notes" {
		t.Fatalf("description was not saved: step=%d description=%q", add.step, add.draft.Description)
	}

	add.input.SetValue("medium")
	updated, _ = add.submitAddField()
	add = updated.(addModel)
	if add.step != addStepTags || add.draft.Priority != models.PriorityMedium {
		t.Fatalf("priority was not saved: step=%d priority=%q", add.step, add.draft.Priority)
	}
}

func TestAddModelSuggestsCommaSeparatedTags(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	st := store.NewStore(sqlite.NewStore())
	if err := st.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	if err := st.InsertTask(models.Task{Title: "Source task", Tags: []string{"urgent", "work"}}); err != nil {
		t.Fatalf("InsertTask() error = %v", err)
	}

	m := newAddModel(&config.Config{}, st, 80)
	m.nextAddStep(addStepTags, "Enter tags (comma-separated)...", "Tags")
	m.loadTags()
	m.input.SetValue("work, ur")
	m.updateTagSuggestions()

	if !reflect.DeepEqual(m.suggestions, []string{"urgent"}) {
		t.Fatalf("suggestions = %#v, want urgent", m.suggestions)
	}

	m.input.SetValue(completedTagValue(m.input.Value(), m.suggestions[m.selected]))
	if m.input.Value() != "work, urgent" {
		t.Errorf("input = %q, want work, urgent", m.input.Value())
	}
}

func TestParseTags(t *testing.T) {
	got := parseTags("work, urgent, ,planning")
	want := []string{"work", "urgent", "planning"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("parseTags() = %#v, want %#v", got, want)
	}
}

func TestTasksCommandRejectsArguments(t *testing.T) {
	m := newModel(nil, nil)
	m.input.SetValue("/tasks one two")
	m.execute()

	if m.active != nil {
		t.Fatalf("active model = %T, want nil", m.active)
	}
	if m.content != "Usage: /tasks" {
		t.Errorf("content = %q, want usage message", m.content)
	}
}

func TestTasksCommandCreatesChildModel(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	st := store.NewStore(sqlite.NewStore())
	if err := st.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	if err := st.InsertTask(models.Task{Title: "Write tests", Category: models.CategoryTask, Priority: models.PriorityLow, Status: models.StatusOpen}); err != nil {
		t.Fatalf("InsertTask() error = %v", err)
	}

	cfg := &config.Config{ListFlags: config.ConfigListFlags{Status: []string{string(models.StatusOpen)}}}
	m := newModel(cfg, st)
	m.input.SetValue("/tasks")
	m.execute()

	list, ok := m.active.(tasksModel)
	if !ok {
		t.Fatalf("active model = %T, want tasksModel", m.active)
	}
	if !strings.Contains(list.content, "Write tests") {
		t.Errorf("tasks content = %q, want it to contain the task title", list.content)
	}
}

func TestTasksCommandTogglingFilterUpdatesTasks(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	st := store.NewStore(sqlite.NewStore())
	if err := st.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	if err := st.InsertTask(models.Task{Title: "Open task", Category: models.CategoryTask, Priority: models.PriorityLow, Status: models.StatusOpen}); err != nil {
		t.Fatalf("InsertTask() error = %v", err)
	}
	if err := st.InsertTask(models.Task{Title: "Closed task", Category: models.CategoryTask, Priority: models.PriorityLow, Status: models.StatusClosed}); err != nil {
		t.Fatalf("InsertTask() error = %v", err)
	}

	cfg := &config.Config{ListFlags: config.ConfigListFlags{Status: []string{string(models.StatusOpen)}}}
	m := newTasksModel(cfg, st, 80)
	if !strings.Contains(m.content, "Open task") || strings.Contains(m.content, "Closed task") {
		t.Fatalf("default filtered content = %q, want only open task", m.content)
	}

	// Move the cursor to the "Closed" status option and toggle it on.
	m.cursor = len(categoryFilterOptions) + 1
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeySpace})
	m = updated.(tasksModel)

	if !strings.Contains(m.content, "Open task") || !strings.Contains(m.content, "Closed task") {
		t.Fatalf("updated content = %q, want both tasks after toggling status filter", m.content)
	}
}

func TestTaskFiltersRenderOptionsHorizontally(t *testing.T) {
	m := tasksModel{
		width: 80,
		groups: []filterGroup{
			newFilterGroup("Category", categoryFilterOptions),
			newFilterGroup("Status", statusFilterOptions, string(models.StatusOpen)),
			newFilterGroup("Priority", priorityFilterOptions),
		},
	}

	rendered := m.renderFilters()
	for _, want := range []string{
		"Category: > [ ] Task   [ ] Meeting",
		"Status:     [x] Open   [ ] Closed",
		"Priority:   [ ] Low   [ ] Medium   [ ] High",
	} {
		if !strings.Contains(rendered, want) {
			t.Errorf("renderFilters() = %q, want it to contain %q", rendered, want)
		}
	}
}

func TestTaskFiltersWrapAtTerminalWidth(t *testing.T) {
	m := tasksModel{
		width: 32,
		groups: []filterGroup{
			newFilterGroup("Category", categoryFilterOptions),
		},
	}

	for _, line := range strings.Split(strings.TrimSpace(m.renderFilters()), "\n") {
		if len(line) > 28 {
			t.Errorf("filter line %q exceeds available width", line)
		}
	}
}

func TestRenderTagCounts(t *testing.T) {
	rendered := renderTagCounts([]models.TagCount{
		{Tag: "planning", Count: 2},
		{Tag: "work", Count: 1},
	})
	if rendered != "Tags\n\n(2) planning\n(1) work" {
		t.Errorf("renderTagCounts() = %q", rendered)
	}
}

func TestTagsCommandRejectsArguments(t *testing.T) {
	m := newModel(nil, nil)
	m.input.SetValue("/tags extra")
	m.execute()

	if m.content != "Usage: /tags" {
		t.Errorf("tags content = %q", m.content)
	}
}

func TestTagCommandCreatesChildModel(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	st := store.NewStore(sqlite.NewStore())
	if err := st.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	m := newModel(nil, st)
	m.input.SetValue("/tag 1")
	m.execute()

	if _, ok := m.active.(tagModel); !ok {
		t.Errorf("active model = %T, want tagModel", m.active)
	}
}

func TestTagModelSuggestsAndSavesTags(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	st := store.NewStore(sqlite.NewStore())
	if err := st.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	if err := st.InsertTask(models.Task{Title: "Tagged task"}); err != nil {
		t.Fatalf("InsertTask() error = %v", err)
	}
	if err := st.InsertTask(models.Task{Title: "Source task", Tags: []string{"urgent", "work"}}); err != nil {
		t.Fatalf("InsertTask() error = %v", err)
	}
	if err := fs.UpdateIDMapFile([]models.Task{{ID: 1, ListID: 1}, {ID: 2, ListID: 2}}); err != nil {
		t.Fatalf("UpdateIDMapFile() error = %v", err)
	}

	m := newTagModel(st, 80, "1")
	m.input.SetValue("ur")
	m.updateSuggestions()
	if !reflect.DeepEqual(m.suggestions, []string{"urgent"}) {
		t.Fatalf("suggestions = %#v, want urgent", m.suggestions)
	}

	m.completeSuggestion()
	if m.input.Value() != "urgent" {
		t.Fatalf("input = %q, want urgent", m.input.Value())
	}
	m.input.SetValue("urgent, new-tag")
	if !m.saveTags() {
		t.Fatal("saveTags() returned false")
	}

	task, err := st.GetTaskByID(1)
	if err != nil {
		t.Fatalf("GetTaskByID() error = %v", err)
	}
	if !reflect.DeepEqual(task.Tags, []string{"urgent", "new-tag"}) {
		t.Errorf("tags = %#v, want urgent and new-tag", task.Tags)
	}
}

func TestListModelClosesToRoot(t *testing.T) {
	m := tasksModel{content: "No tasks found."}
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	msg := cmd()

	if _, ok := msg.(tasksClosedMsg); !ok {
		t.Fatalf("close message = %T, want tasksClosedMsg", msg)
	}
}

func TestRootViewportScrollsCommandContent(t *testing.T) {
	m := newModel(nil, nil)
	m.viewport = viewport.New(20, 1)
	m.active = detailModel{content: "first\nsecond"}
	m.setViewportContent()

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = updated.(model)
	if m.viewport.YOffset != 1 {
		t.Errorf("viewport offset = %d, want 1", m.viewport.YOffset)
	}
}

func TestRenderNotesShowsCreationTimeAndContent(t *testing.T) {
	notes := []models.Note{
		{
			Content:   "Follow up with the design team.",
			CreatedAt: time.Date(2026, time.August, 14, 20, 1, 0, 0, time.UTC),
		},
	}

	rendered := (&detailModel{}).renderNotes(notes, 80)
	for _, want := range []string{"Notes", "Created At", "2026-08-14 20:01", "Follow up with the design team."} {
		if !strings.Contains(rendered, want) {
			t.Errorf("renderNotes() = %q, want it to contain %q", rendered, want)
		}
	}
}

func TestRenderSessionsShowsTotalAndBreakdown(t *testing.T) {
	startedAt := time.Date(2026, time.August, 14, 20, 1, 0, 0, time.UTC)
	task := models.Task{
		Duration: time.Time{}.Add(95 * time.Minute),
		Sessions: []models.Session{
			{StartTime: startedAt.Add(65 * time.Minute), EndTime: startedAt.Add(95 * time.Minute)},
			{StartTime: startedAt, EndTime: startedAt.Add(65 * time.Minute)},
		},
	}

	rendered := (&detailModel{}).renderSessions(task, 38)
	for _, want := range []string{
		"Sessions",
		"Total time spent: 01:35",
		"Started At",
		"Duration",
		"2026-08-14 20:01",
		"01:05",
		"2026-08-14 21:06",
		"00:30",
	} {
		if !strings.Contains(rendered, want) {
			t.Errorf("renderSessions() = %q, want it to contain %q", rendered, want)
		}
	}
}

func TestDetailWidthHasMinimum(t *testing.T) {
	if got := detailWidth(80); got != detailMinWidth {
		t.Errorf("detailWidth(80) = %d, want %d", got, detailMinWidth)
	}
	if got := detailWidth(160); got != 160 {
		t.Errorf("detailWidth(160) = %d, want 160", got)
	}
}
