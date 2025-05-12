package repo

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/yourusername/atprogo/pkg/lexicon"
)

// Commit represents a commit in a repository
type Commit struct {
	ID        string    `json:"id"`
	Prev      string    `json:"prev,omitempty"`
	Data      []byte    `json:"data"`
	Signature []byte    `json:"sig,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}

// Repository represents a user's repository
type Repository struct {
	DID       string            `json:"did"`
	Head      string            `json:"head"`
	Commits   map[string]Commit `json:"commits"`
	Documents map[string][]byte `json:"documents"`
}

// Store is an interface for repository storage
type Store interface {
	GetRepository(ctx context.Context, did string) (*Repository, error)
	SaveRepository(ctx context.Context, repo *Repository) error
	GetCommit(ctx context.Context, did string, commitID string) (*Commit, error)
	SaveCommit(ctx context.Context, did string, commit *Commit) error
}

// NewRepository creates a new repository
func NewRepository(did string) *Repository {
	return &Repository{
		DID:       did,
		Commits:   make(map[string]Commit),
		Documents: make(map[string][]byte),
	}
}

// CreateCommit creates a new commit in the repository
func (r *Repository) CreateCommit(doc *lexicon.Document) (*Commit, error) {
	// Marshal document to JSON
	data, err := lexicon.MarshalDocument(doc)
	if err != nil {
		return nil, err
	}

	// Create commit
	commit := &Commit{
		Prev:      r.Head,
		Data:      data,
		CreatedAt: time.Now(),
	}

	// Generate commit ID
	hash := sha256.Sum256(data)
	commit.ID = hex.EncodeToString(hash[:])

	// Update repository
	r.Commits[commit.ID] = *commit
	r.Head = commit.ID
	r.Documents[doc.ID] = data

	return commit, nil
}

// GetDocument gets a document from the repository
func (r *Repository) GetDocument(id string) (*lexicon.Document, error) {
	data, ok := r.Documents[id]
	if !ok {
		return nil, fmt.Errorf("document not found: %s", id)
	}

	return lexicon.UnmarshalDocument(data)
}

// MemoryStore is an in-memory implementation of Store
type MemoryStore struct {
	repositories map[string]*Repository
}

// NewMemoryStore creates a new memory store
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		repositories: make(map[string]*Repository),
	}
}

// GetRepository gets a repository from the store
func (s *MemoryStore) GetRepository(ctx context.Context, did string) (*Repository, error) {
	repo, ok := s.repositories[did]
	if !ok {
		return nil, fmt.Errorf("repository not found: %s", did)
	}
	return repo, nil
}

// SaveRepository saves a repository to the store
func (s *MemoryStore) SaveRepository(ctx context.Context, repo *Repository) error {
	s.repositories[repo.DID] = repo
	return nil
}

// GetCommit gets a commit from the store
func (s *MemoryStore) GetCommit(ctx context.Context, did string, commitID string) (*Commit, error) {
	repo, err := s.GetRepository(ctx, did)
	if err != nil {
		return nil, err
	}

	commit, ok := repo.Commits[commitID]
	if !ok {
		return nil, fmt.Errorf("commit not found: %s", commitID)
	}

	return &commit, nil
}

// SaveCommit saves a commit to the store
func (s *MemoryStore) SaveCommit(ctx context.Context, did string, commit *Commit) error {
	repo, err := s.GetRepository(ctx, did)
	if err != nil {
		return err
	}

	repo.Commits[commit.ID] = *commit
	return nil
}
