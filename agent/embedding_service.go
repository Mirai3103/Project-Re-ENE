package agent

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Mirai3103/Project-Re-ENE/config"
	"github.com/Mirai3103/Project-Re-ENE/embedding"
	"github.com/Mirai3103/Project-Re-ENE/store"
	"github.com/google/uuid"
)

type EmbeddingService struct {
	cfg    *config.Config
	logger *slog.Logger
	model  embedding.Model
	store  *store.Queries
}

func NewEmbeddingService(cfg *config.Config, logger *slog.Logger, model embedding.Model, store *store.Queries) *EmbeddingService {
	return &EmbeddingService{cfg: cfg, logger: logger, model: model, store: store}
}

func (s *EmbeddingService) EmbedText(ctx context.Context, text string) ([]float32, error) {
	return s.model.Get(ctx, text)
}

func (s *EmbeddingService) AddMemory(ctx context.Context, memory *store.Memory) error {
	vector, err := s.model.Get(ctx, *memory.Content)
	if err != nil {
		return err
	}
	// check if memory already exists
	similarMemories, err := s.store.SimilarMemories(ctx, vector, 10, 0.7)
	if err != nil {
		return err
	}
	if len(similarMemories) > 0 {
		return fmt.Errorf("memory already exists")
	}
	memory.Embedding = store.Float32ToBytes(vector)
	return s.store.CreateMemory(ctx, store.CreateMemoryParams{
		ID:          memory.ID,
		Content:     memory.Content,
		Importance:  memory.Importance,
		Confidence:  memory.Confidence,
		Tags:        memory.Tags,
		Embedding:   memory.Embedding,
		UserID:      memory.UserID,
		Source:      memory.Source,
		CharacterID: memory.CharacterID,
	})
}

func (s *EmbeddingService) AddUserFact(ctx context.Context, fact *store.UserFact) error {
	return s.store.AddUserFact(ctx, store.AddUserFactParams{
		ID:     uuid.New().String(),
		UserID: fact.UserID,
		Name:   fact.Name,
		Value:  fact.Value,
		Type:   fact.Type,
	})
}

func (s *EmbeddingService) AddCharacterFact(ctx context.Context, fact *store.CharacterFact) error {
	return s.store.AddCharacterFact(ctx, store.AddCharacterFactParams{
		ID:          uuid.New().String(),
		CharacterID: fact.CharacterID,
		Name:        fact.Name,
		Value:       fact.Value,
		Type:        fact.Type,
	})
}
