package pds

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository represents a user's repository
type Repository struct {
	DID       string    `json:"did"`
	Head      string    `json:"head"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Commit represents a commit in a repository
type Commit struct {
	ID           string    `json:"id"`
	RepositoryDID string    `json:"repositoryDid"`
	Prev         string    `json:"prev,omitempty"`
	Data         []byte    `json:"data"`
	Signature    []byte    `json:"signature,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
}

// Document represents a document in a repository
type Document struct {
	ID           string          `json:"id"`
	RepositoryDID string          `json:"repositoryDid"`
	Type         string          `json:"type"`
	Value        json.RawMessage `json:"value"`
	CreatedAt    time.Time       `json:"createdAt"`
	UpdatedAt    time.Time       `json:"updatedAt"`
}

// RepositoryRepository handles repository data access
type RepositoryRepository struct {
	db *pgxpool.Pool
}

// NewRepositoryRepository creates a new repository repository
func NewRepositoryRepository(db *pgxpool.Pool) *RepositoryRepository {
	return &RepositoryRepository{db: db}
}

// CreateRepository creates a new repository
func (r *RepositoryRepository) CreateRepository(ctx context.Context, repo *Repository) error {
	query := `
		INSERT INTO repositories (did, head, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
	`
	now := time.Now()
	_, err := r.db.Exec(ctx, query,
		repo.DID,
		repo.Head,
		now,
		now,
	)
	if err != nil {
		return fmt.Errorf("failed to create repository: %w", err)
	}

	repo.CreatedAt = now
	repo.UpdatedAt = now

	return nil
}

// GetRepository gets a repository by DID
func (r *RepositoryRepository) GetRepository(ctx context.Context, did string) (*Repository, error) {
	query := `
		SELECT did, head, created_at, updated_at
		FROM repositories
		WHERE did = $1
	`
	var repo Repository
	err := r.db.QueryRow(ctx, query, did).Scan(
		&repo.DID,
		&repo.Head,
		&repo.CreatedAt,
		&repo.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get repository: %w", err)
	}
	return &repo, nil
}

// UpdateRepository updates a repository
func (r *RepositoryRepository) UpdateRepository(ctx context.Context, repo *Repository) error {
	query := `
		UPDATE repositories
		SET head = $1, updated_at = $2
		WHERE did = $3
	`
	now := time.Now()
	_, err := r.db.Exec(ctx, query,
		repo.Head,
		now,
		repo.DID,
	)
	if err != nil {
		return fmt.Errorf("failed to update repository: %w", err)
	}

	repo.UpdatedAt = now

	return nil
}

// CreateCommit creates a new commit
func (r *RepositoryRepository) CreateCommit(ctx context.Context, commit *Commit) error {
	query := `
		INSERT INTO commits (id, repository_did, prev, data, signature, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	now := time.Now()
	_, err := r.db.Exec(ctx, query,
		commit.ID,
		commit.RepositoryDID,
		commit.Prev,
		commit.Data,
		commit.Signature,
		now,
	)
	if err != nil {
		return fmt.Errorf("failed to create commit: %w", err)
	}

	commit.CreatedAt = now

	return nil
}

// GetCommit gets a commit by ID
func (r *RepositoryRepository) GetCommit(ctx context.Context, id string) (*Commit, error) {
	query := `
		SELECT id, repository_did, prev, data, signature, created_at
		FROM commits
		WHERE id = $1
	`
	var commit Commit
	err := r.db.QueryRow(ctx, query, id).Scan(
		&commit.ID,
		&commit.RepositoryDID,
		&commit.Prev,
		&commit.Data,
		&commit.Signature,
		&commit.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get commit: %w", err)
	}
	return &commit, nil
}

// CreateDocument creates a new document
func (r *RepositoryRepository) CreateDocument(ctx context.Context, doc *Document) error {
	query := `
		INSERT INTO documents (id, repository_did, type, value, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	now := time.Now()
	_, err := r.db.Exec(ctx, query,
		doc.ID,
		doc.RepositoryDID,
		doc.Type,
		doc.Value,
		now,
		now,
	)
	if err != nil {
		return fmt.Errorf("failed to create document: %w", err)
	}

	doc.CreatedAt = now
	doc.UpdatedAt = now

	return nil
}

// GetDocument gets a document by ID
func (r *RepositoryRepository) GetDocument(ctx context.Context, repositoryDID, id string) (*Document, error) {
	query := `
		SELECT id, repository_did, type, value, created_at, updated_at
		FROM documents
		WHERE repository_did = $1 AND id = $2
	`
	var doc Document
	err := r.db.QueryRow(ctx, query, repositoryDID, id).Scan(
		&doc.ID,
		&doc.RepositoryDID,
		&doc.Type,
		&doc.Value,
		&doc.CreatedAt,
		&doc.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get document: %w", err)
	}
	return &doc, nil
}

// GetDocumentsByType gets documents by type
func (r *RepositoryRepository) GetDocumentsByType(ctx context.Context, repositoryDID, docType string) ([]*Document, error) {
	query := `
		SELECT id, repository_did, type, value, created_at, updated_at
		FROM documents
		WHERE repository_did = $1 AND type = $2
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, repositoryDID, docType)
	if err != nil {
		return nil, fmt.Errorf("failed to get documents: %w", err)
	}
	defer rows.Close()

	var docs []*Document
	for rows.Next() {
		var doc Document
		err := rows.Scan(
			&doc.ID,
			&doc.RepositoryDID,
			&doc.Type,
			&doc.Value,
			&doc.CreatedAt,
			&doc.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan document: %w", err)
		}
		docs = append(docs, &doc)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating documents: %w", err)
	}
	return docs, nil
}

// CreatePost creates a new post document
func (r *RepositoryRepository) CreatePost(ctx context.Context, did string, content string) (*Document, error) {
	// Get repository
	repo, err := r.GetRepository(ctx, did)
	if err != nil {
		// Create repository if it doesn't exist
		repo = &Repository{
			DID:  did,
			Head: "",
		}
		if err := r.CreateRepository(ctx, repo); err != nil {
			return nil, fmt.Errorf("failed to create repository: %w", err)
		}
	}

	// Create post value
	postValue := map[string]interface{}{
		"text":      content,
		"createdAt": time.Now().Format(time.RFC3339),
	}
	postValueJSON, err := json.Marshal(postValue)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal post value: %w", err)
	}

	// Create document ID
	postID := fmt.Sprintf("post/%s", time.Now().Format("20060102150405.000"))

	// Create document
	doc := &Document{
		ID:           postID,
		RepositoryDID: did,
		Type:         "app.bsky.feed.post",
		Value:        postValueJSON,
	}
	if err := r.CreateDocument(ctx, doc); err != nil {
		return nil, fmt.Errorf("failed to create document: %w", err)
	}

	// Create commit data
	commitData := map[string]interface{}{
		"op":    "create",
		"path":  postID,
		"value": postValue,
	}
	commitDataJSON, err := json.Marshal(commitData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal commit data: %w", err)
	}

	// Create commit ID
	hash := sha256.Sum256(commitDataJSON)
	commitID := hex.EncodeToString(hash[:])

	// Create commit
	commit := &Commit{
		ID:           commitID,
		RepositoryDID: did,
		Prev:         repo.Head,
		Data:         commitDataJSON,
	}
	if err := r.CreateCommit(ctx, commit); err != nil {
		return nil, fmt.Errorf("failed to create commit: %w", err)
	}

	// Update repository head
	repo.Head = commitID
	if err := r.UpdateRepository(ctx, repo); err != nil {
		return nil, fmt.Errorf("failed to update repository: %w", err)
	}

	return doc, nil
}
