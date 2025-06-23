package usecases

import (
	"context"
	"fmt"
	"log"

	"github.com/thiagohmm/integracaoThothConsumer/internal/domain/entities"
	"go.opentelemetry.io/otel"
)

// StatusUseCase encapsula o repositório para atualização de status
type StatusUseCase struct {
	Repo entities.StatusRepository
}

// NewStatusUseCase retorna uma nova instância de StatusUseCase
func NewStatusUseCase(repo entities.StatusRepository) *StatusUseCase {
	return &StatusUseCase{Repo: repo}
}

// UpdateStatusProcesso atualiza o status do processo usando o repositório
func (uc *StatusUseCase) UpdateStatusProcesso(ctx context.Context, uuid string, novoStatus string) error {
	ctx, span := otel.Tracer("StatusUseCase").Start(ctx, "UpdateStatusProcesso")
	defer span.End()

	err := uc.Repo.UpdateStatusProcesso(ctx, uuid, novoStatus)
	if err != nil {
		return fmt.Errorf("error updating status: %w", err)
	}

	log.Printf("Status atualizado com sucesso para UUID: %s", uuid)
	return nil
}
