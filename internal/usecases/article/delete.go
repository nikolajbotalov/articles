package usecases

import (
	"context"
	"go.uber.org/zap"
)

func (uc *articleUseCase) DeleteArticle(ctx context.Context, id string) error {
	uc.logger.Info("Deleting article", zap.String("id", id))

	if err := uc.repo.Delete(ctx, id); err != nil {
		uc.logger.Error("Failed to delete article", zap.Error(err))
		return err
	}

	uc.logger.Info("Successfully deleted article", zap.String("id", id))
	return nil
}
