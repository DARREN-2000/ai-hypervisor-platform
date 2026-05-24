package vmmanager

import (
	"context"
	"time"

	"github.com/DARREN-2000/ai-hypervisor-platform/internal/models"
	apierrors "github.com/DARREN-2000/ai-hypervisor-platform/pkg/errors"
)

// CreateTemplate stores a new VM template.
func (s *Service) CreateTemplate(ctx context.Context, template *models.VMTemplate) error {
	if s.templateRepo == nil {
		return apierrors.InternalError("template repository is not configured")
	}
	if template == nil {
		return apierrors.ValidationError("template is required")
	}
	if template.ID == "" {
		template.ID = models.NewID()
	}
	if template.Name == "" {
		return apierrors.ValidationError("template name is required")
	}

	now := time.Now().UTC()
	if template.CreatedAt.IsZero() {
		template.CreatedAt = now
	}
	template.UpdatedAt = now

	if err := s.templateRepo.Create(ctx, template); err != nil {
		return apierrors.InternalError("failed to create template").WithCause(err)
	}
	return nil
}

// UpdateTemplate updates an existing VM template.
func (s *Service) UpdateTemplate(ctx context.Context, template *models.VMTemplate) error {
	if s.templateRepo == nil {
		return apierrors.InternalError("template repository is not configured")
	}
	if template == nil || template.ID == "" {
		return apierrors.ValidationError("template id is required")
	}

	template.UpdatedAt = time.Now().UTC()
	if err := s.templateRepo.Update(ctx, template); err != nil {
		return apierrors.InternalError("failed to update template").WithCause(err)
	}
	return nil
}

// GetTemplate retrieves a template by ID.
func (s *Service) GetTemplate(ctx context.Context, templateID string) (*models.VMTemplate, error) {
	if s.templateRepo == nil {
		return nil, apierrors.InternalError("template repository is not configured")
	}
	if templateID == "" {
		return nil, apierrors.ValidationError("template id is required")
	}

	template, err := s.templateRepo.Get(ctx, templateID)
	if err != nil {
		return nil, apierrors.NotFound("template").WithCause(err)
	}
	return template, nil
}

// ListTemplates returns templates.
func (s *Service) ListTemplates(ctx context.Context, filters map[string]string) ([]*models.VMTemplate, error) {
	if s.templateRepo == nil {
		return nil, apierrors.InternalError("template repository is not configured")
	}

	templates, err := s.templateRepo.List(ctx, filters)
	if err != nil {
		return nil, apierrors.InternalError("failed to list templates").WithCause(err)
	}
	return templates, nil
}

// DeleteTemplate removes a template.
func (s *Service) DeleteTemplate(ctx context.Context, templateID string) error {
	if s.templateRepo == nil {
		return apierrors.InternalError("template repository is not configured")
	}
	if templateID == "" {
		return apierrors.ValidationError("template id is required")
	}

	if err := s.templateRepo.Delete(ctx, templateID); err != nil {
		return apierrors.InternalError("failed to delete template").WithCause(err)
	}
	return nil
}

// CreateVMFromTemplate creates a VM using a template and overrides.
func (s *Service) CreateVMFromTemplate(ctx context.Context, templateID string, overrides *models.VirtualMachine) error {
	template, err := s.GetTemplate(ctx, templateID)
	if err != nil {
		return err
	}

	vm := &models.VirtualMachine{
		ID:            models.NewID(),
		Name:          template.Name,
		Namespace:     "default",
		Flavor:        template.Flavor,
		Image:         template.Image,
		NetworkConfig: template.NetworkConfig,
		StorageConfig: template.StorageConfig,
		GPURequests:   template.GPURequests,
		Metadata:      template.Metadata,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}

	if overrides != nil {
		if overrides.Name != "" {
			vm.Name = overrides.Name
		}
		if overrides.Namespace != "" {
			vm.Namespace = overrides.Namespace
		}
		if overrides.Flavor.CPU > 0 {
			vm.Flavor = overrides.Flavor
		}
		if overrides.Image.ID != "" {
			vm.Image = overrides.Image
		}
		if len(overrides.GPURequests) > 0 {
			vm.GPURequests = overrides.GPURequests
		}
		if overrides.Metadata != nil {
			vm.Metadata = overrides.Metadata
		}
	}

	return s.CreateVM(ctx, vm)
}
