package googledocs

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
	"ttm/pkg/config"
	"ttm/pkg/models"

	docs "google.golang.org/api/docs/v1"
	"google.golang.org/api/option"
)

type Store struct {
	config  config.GoogleDocsConfig
	service *docs.Service
	mu      sync.Mutex
}

type documentData struct {
	Tasks    []models.Task    `json:"tasks"`
	Sessions []models.Session `json:"sessions"`
	Notes    []models.Note    `json:"notes"`
}

func NewStore(cfg config.GoogleDocsConfig) *Store {
	return &Store{config: cfg}
}

func (s *Store) Init() error {
	if s.config.DocumentID == "" {
		return fmt.Errorf("google docs storage requires storage.googleDocs.documentId")
	}
	if s.config.CredentialsFile == "" {
		return fmt.Errorf("google docs storage requires storage.googleDocs.credentialsFile")
	}

	service, err := docs.NewService(
		context.Background(),
		option.WithCredentialsFile(s.config.CredentialsFile),
		option.WithScopes(docs.DocumentsScope),
	)
	if err != nil {
		return fmt.Errorf("create Google Docs client: %w", err)
	}
	s.service = service

	_, err = s.read()
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) InsertTask(task models.Task) (models.Task, error) {
	err := s.modify(func(data *documentData) error {
		task.ID = nextTaskID(data.Tasks)
		task.CreatedAt = time.Now()
		task.UpdatedAt = task.CreatedAt
		data.Tasks = append(data.Tasks, task)
		return nil
	})

	if err != nil {
		return models.Task{}, err
	}

	return task, nil
}

func (s *Store) GetTaskByID(taskID int64) (models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := s.read()
	if err != nil {
		return models.Task{}, err
	}
	for _, task := range data.Tasks {
		if task.ID == taskID {
			return task, nil
		}
	}
	return models.Task{}, fmt.Errorf("task %d not found", taskID)
}

func (s *Store) ListTasks(search string, category models.Category, status models.Status, priority models.Priority) ([]models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := s.read()
	if err != nil {
		return nil, err
	}

	tasks := make([]models.Task, 0, len(data.Tasks))
	for _, task := range data.Tasks {
		if (search == "" || strings.Contains(task.Title, search)) &&
			(search == "" || strings.Contains(task.Description, search)) &&
			(category == "" || task.Category == category) &&
			(status == "" || task.Status == status) &&
			(priority == "" || task.Priority == priority) {
			tasks = append(tasks, task)
		}
	}
	return tasks, nil
}

func (s *Store) SearchTasks(search models.TaskSearch) ([]models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := s.read()
	if err != nil {
		return nil, err
	}

	tasks := make([]models.Task, 0, len(data.Tasks))
	for _, task := range data.Tasks {
		if taskMatchesSearch(task, search) {
			tasks = append(tasks, task)
		}
	}
	return tasks, nil
}

func taskMatchesSearch(task models.Task, search models.TaskSearch) bool {
	contains := func(value, search string) bool {
		return search == "" || strings.Contains(strings.ToLower(value), strings.ToLower(search))
	}
	hasTag := func(search string) bool {
		for _, tag := range task.Tags {
			if contains(tag, search) {
				return true
			}
		}
		return false
	}
	if search.General != "" &&
		!contains(task.Title, search.General) &&
		!contains(task.Description, search.General) &&
		!contains(string(task.Category), search.General) &&
		!contains(string(task.Priority), search.General) &&
		!contains(string(task.Status), search.General) &&
		!hasTag(search.General) {
		return false
	}
	if !contains(task.Title, search.Title) ||
		!contains(task.Description, search.Description) ||
		!contains(string(task.Category), search.Category) ||
		!contains(string(task.Priority), search.Priority) ||
		!contains(string(task.Status), search.Status) {
		return false
	}
	for _, tag := range search.Tags {
		if !hasTag(tag) {
			return false
		}
	}
	return true
}

func (s *Store) UpdateTitle(taskID int64, title string) error {
	return s.updateTask(taskID, func(task *models.Task) { task.Title = title })
}

func (s *Store) UpdateDescription(taskID int64, description string) error {
	return s.updateTask(taskID, func(task *models.Task) { task.Description = description })
}

func (s *Store) UpdateCategory(taskID int64, category models.Category) error {
	return s.updateTask(taskID, func(task *models.Task) { task.Category = category })
}

func (s *Store) UpdatePriority(taskID int64, priority models.Priority) error {
	return s.updateTask(taskID, func(task *models.Task) { task.Priority = priority })
}

func (s *Store) UpdateStatus(taskID int64, status models.Status) error {
	return s.updateTask(taskID, func(task *models.Task) { task.Status = status })
}

func (s *Store) UpdateTags(taskID int64, tags []string) error {
	return s.updateTask(taskID, func(task *models.Task) { task.Tags = tags })
}

func (s *Store) UpdateOpenedAt(taskID int64, openedAt time.Time) error {
	return s.updateTask(taskID, func(task *models.Task) { task.OpenedAt = openedAt })
}

func (s *Store) UpdateClosedAt(taskID int64, closedAt time.Time) error {
	return s.updateTask(taskID, func(task *models.Task) { task.ClosedAt = closedAt })
}

func (s *Store) AddSession(session models.Session) error {
	return s.modify(func(data *documentData) error {
		if !hasTask(data.Tasks, session.TaskId) {
			return fmt.Errorf("task %d not found", session.TaskId)
		}
		session.ID = nextSessionID(data.Sessions)
		data.Sessions = append(data.Sessions, session)
		return nil
	})
}

func (s *Store) GetSessionsByTaskID(taskID int) ([]models.Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := s.read()
	if err != nil {
		return nil, err
	}
	sessions := []models.Session{}
	for _, session := range data.Sessions {
		if session.TaskId == int64(taskID) {
			sessions = append(sessions, session)
		}
	}
	return sessions, nil
}

func (s *Store) GetSessionsByTimeRange(startTime time.Time, endTime time.Time) ([]models.Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := s.read()
	if err != nil {
		return nil, err
	}
	sessions := []models.Session{}
	for _, session := range data.Sessions {
		if !session.StartTime.Before(startTime) && !session.EndTime.After(endTime) {
			sessions = append(sessions, session)
		}
	}
	return sessions, nil
}

func (s *Store) InsertTags(taskID int64, tags []string) error {
	return s.updateTask(taskID, func(task *models.Task) {
		task.Tags = append(task.Tags, tags...)
	})
}

func (s *Store) GetTagsByTaskID(taskID int64) ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := s.read()
	if err != nil {
		return nil, err
	}
	for _, task := range data.Tasks {
		if task.ID == taskID {
			return task.Tags, nil
		}
	}
	return nil, fmt.Errorf("task %d not found", taskID)
}

func (s *Store) ListTagCounts() ([]models.TagCount, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := s.read()
	if err != nil {
		return nil, err
	}

	counts := make(map[string]int)
	for _, task := range data.Tasks {
		taskTags := make(map[string]struct{})
		for _, tag := range task.Tags {
			taskTags[tag] = struct{}{}
		}
		for tag := range taskTags {
			if tag != "" {
				counts[tag]++
			}
		}
	}

	tags := make([]models.TagCount, 0, len(counts))
	for tag, count := range counts {
		tags = append(tags, models.TagCount{Tag: tag, Count: count})
	}
	sort.Slice(tags, func(i, j int) bool {
		left, right := strings.ToLower(tags[i].Tag), strings.ToLower(tags[j].Tag)
		if left == right {
			return tags[i].Tag < tags[j].Tag
		}
		return left < right
	})
	return tags, nil
}

func (s *Store) InsertNote(note models.Note) (models.Note, error) {
	err := s.modify(func(data *documentData) error {
		if !hasTask(data.Tasks, note.TaskID) {
			return fmt.Errorf("task %d not found", note.TaskID)
		}
		note.ID = nextNoteID(data.Notes)
		note.CreatedAt = time.Now()
		data.Notes = append(data.Notes, note)
		return nil
	})

	if err != nil {
		return models.Note{}, err
	}

	return note, nil
}

func (s *Store) GetNotesByTaskID(taskID int64) ([]models.Note, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := s.read()
	if err != nil {
		return nil, err
	}
	notes := []models.Note{}
	for _, note := range data.Notes {
		if note.TaskID == taskID {
			notes = append(notes, note)
		}
	}
	return notes, nil
}

func (s *Store) updateTask(taskID int64, update func(*models.Task)) error {
	return s.modify(func(data *documentData) error {
		for i := range data.Tasks {
			if data.Tasks[i].ID == taskID {
				update(&data.Tasks[i])
				data.Tasks[i].UpdatedAt = time.Now()
				return nil
			}
		}
		return fmt.Errorf("task %d not found", taskID)
	})
}

func (s *Store) modify(modifier func(*documentData) error) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := s.read()
	if err != nil {
		return err
	}
	if err := modifier(&data); err != nil {
		return err
	}
	return s.write(data)
}

func (s *Store) read() (documentData, error) {
	document, err := s.service.Documents.Get(s.config.DocumentID).Do()
	if err != nil {
		return documentData{}, fmt.Errorf("read Google Doc %q: %w", s.config.DocumentID, err)
	}

	text := strings.TrimSpace(documentText(document))
	if text == "" {
		return documentData{Tasks: []models.Task{}, Sessions: []models.Session{}}, nil
	}

	var data documentData
	if err := json.Unmarshal([]byte(text), &data); err != nil {
		return documentData{}, fmt.Errorf("decode TTM data in Google Doc %q: %w", s.config.DocumentID, err)
	}
	if data.Tasks == nil {
		data.Tasks = []models.Task{}
	}
	if data.Sessions == nil {
		data.Sessions = []models.Session{}
	}
	return data, nil
}

func (s *Store) write(data documentData) error {
	text, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("encode TTM data: %w", err)
	}

	document, err := s.service.Documents.Get(s.config.DocumentID).Do()
	if err != nil {
		return fmt.Errorf("read Google Doc %q before writing: %w", s.config.DocumentID, err)
	}

	requests := make([]*docs.Request, 0, 2)
	if endIndex := documentEndIndex(document); endIndex > 1 {
		requests = append(requests, &docs.Request{
			DeleteContentRange: &docs.DeleteContentRangeRequest{
				Range: &docs.Range{StartIndex: 1, EndIndex: endIndex},
			},
		})
	}
	requests = append(requests, &docs.Request{
		InsertText: &docs.InsertTextRequest{
			Location: &docs.Location{Index: 1},
			Text:     string(text),
		},
	})

	_, err = s.service.Documents.BatchUpdate(s.config.DocumentID, &docs.BatchUpdateDocumentRequest{
		Requests: requests,
	}).Do()
	if err != nil {
		return fmt.Errorf("write TTM data to Google Doc %q: %w", s.config.DocumentID, err)
	}
	return nil
}

func documentText(document *docs.Document) string {
	var text strings.Builder
	for _, element := range document.Body.Content {
		if element.Paragraph == nil {
			continue
		}
		for _, paragraphElement := range element.Paragraph.Elements {
			if paragraphElement.TextRun != nil {
				text.WriteString(paragraphElement.TextRun.Content)
			}
		}
	}
	return text.String()
}

func documentEndIndex(document *docs.Document) int64 {
	var endIndex int64 = 1
	for _, element := range document.Body.Content {
		if element.EndIndex > endIndex {
			endIndex = element.EndIndex
		}
	}
	return endIndex - 1
}

func nextTaskID(tasks []models.Task) int64 {
	var maxID int64
	for _, task := range tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}
	return maxID + 1
}

func nextSessionID(sessions []models.Session) int64 {
	var maxID int64
	for _, session := range sessions {
		if session.ID > maxID {
			maxID = session.ID
		}
	}
	return maxID + 1
}

func hasTask(tasks []models.Task, taskID int64) bool {
	for _, task := range tasks {
		if task.ID == taskID {
			return true
		}
	}
	return false
}

func nextNoteID(notes []models.Note) int64 {
	var maxID int64
	for _, note := range notes {
		if note.ID > maxID {
			maxID = note.ID
		}
	}
	return maxID + 1
}
