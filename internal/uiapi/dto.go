package uiapi

import "aiguard/internal/model"

type StartReviewRequest struct {
	MRURL         string `json:"mrUrl"`
	RepoURL       string `json:"repoUrl"`
	LocalRepoPath string `json:"localRepoPath"`
	SourceBranch  string `json:"sourceBranch"`
	TargetBranch  string `json:"targetBranch"`
	ConfigPath    string `json:"configPath"`
	WorkspaceDir  string `json:"workspaceDir"`
}

type RuntimeContextRequest struct {
	ConfigPath   string `json:"configPath"`
	WorkspaceDir string `json:"workspaceDir"`
}

type PrepareRepositoryResponse struct {
	RepoURL           string   `json:"repoUrl"`
	SourceBranch      string   `json:"sourceBranch"`
	TargetBranch      string   `json:"targetBranch"`
	AvailableBranches []string `json:"availableBranches"`
	LogPath           string   `json:"logPath"`
	Message           string   `json:"message"`
}

type RepositorySuggestion struct {
	RepoURL       string   `json:"repoUrl"`
	Candidates    []string `json:"candidates"`
	ResolvedByMR  bool     `json:"resolvedByMr"`
	ManualRepoURL string   `json:"manualRepoUrl"`
	Message       string   `json:"message"`
}

type LogState struct {
	LogPath   string `json:"logPath"`
	Content   string `json:"content"`
	UpdatedAt string `json:"updatedAt"`
}

type CacheClearResult struct {
	WorkspaceDir string   `json:"workspaceDir"`
	LogPath      string   `json:"logPath"`
	Removed      []string `json:"removed"`
	Message      string   `json:"message"`
}

type ProgressEvent struct {
	TaskID  string `json:"taskId"`
	Stage   string `json:"stage"`
	Percent int    `json:"percent"`
	Message string `json:"message"`
	High    int    `json:"high"`
	Medium  int    `json:"medium"`
	Low     int    `json:"low"`
}

type ReviewDoneEvent struct {
	TaskID       string       `json:"taskId"`
	ReportDir    string       `json:"reportDir"`
	HTMLPath     string       `json:"htmlPath"`
	MarkdownPath string       `json:"markdownPath"`
	JSONPath     string       `json:"jsonPath"`
	Report       model.Report `json:"report"`
}

type HistoryItem struct {
	TaskID       string        `json:"taskId"`
	Title        string        `json:"title"`
	RepoURL      string        `json:"repoUrl"`
	SourceRef    string        `json:"sourceRef"`
	TargetRef    string        `json:"targetRef"`
	CreatedAt    string        `json:"createdAt"`
	ReportDir    string        `json:"reportDir"`
	HTMLPath     string        `json:"htmlPath"`
	MarkdownPath string        `json:"markdownPath"`
	JSONPath     string        `json:"jsonPath"`
	TotalIssues  int           `json:"totalIssues"`
	Summary      model.Summary `json:"summary"`
}

type OpenReportRequest struct {
	HTMLPath  string `json:"htmlPath"`
	ReportDir string `json:"reportDir"`
}

type OpenPathResult struct {
	Path    string `json:"path"`
	Mode    string `json:"mode"`
	Message string `json:"message"`
}
