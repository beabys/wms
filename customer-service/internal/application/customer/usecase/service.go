package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/beabys/wms/customer-service/internal/application/customer/command"
	"github.com/beabys/wms/customer-service/internal/application/customer/repository"
	"github.com/beabys/wms/customer-service/internal/domain/customer/model"
)

// NewCustomerService creates a new CustomerService.
func NewCustomerService(customerRepo repository.CustomerRepository, inviteRepo repository.InviteLinkRepository) *CustomerService {
	return &CustomerService{
		customerRepo: customerRepo,
		inviteRepo:   inviteRepo,
	}
}

// CreateCustomer handles the create customer use case.
func (s *CustomerService) CreateCustomer(ctx context.Context, req *command.CreateCustomerCommand) (*model.Customer, error) {
	vat, err := model.NewVATNumber(req.VATNumber)
	if err != nil {
		return nil, fmt.Errorf("create customer: %w", err)
	}

	now := time.Now()
	customer := &model.Customer{
		ID:            model.NewID(),
		CompanyName:   req.CompanyName,
		VATNumber:     vat.String(),
		Address:       req.Address,
		Status:        model.CustomerStatusPending,
		RateCardID:    req.RateCardID,
		CreditBalance: 0,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := s.customerRepo.Save(ctx, customer); err != nil {
		return nil, fmt.Errorf("create customer: save: %w", err)
	}

	return customer, nil
}

// InviteCustomer handles the invite customer use case.
func (s *CustomerService) InviteCustomer(ctx context.Context, req *command.InviteCustomerCommand) (*model.InviteLink, error) {
	if req.Email == "" {
		return nil, fmt.Errorf("invite customer: email is required")
	}

	token, err := generateToken()
	if err != nil {
		return nil, fmt.Errorf("invite customer: generate token: %w", err)
	}

	now := time.Now()
	invite := &model.InviteLink{
		ID:        model.NewID(),
		Email:     req.Email,
		Token:     token,
		ExpiresAt: now.Add(7 * 24 * time.Hour), // 7 days
		Used:      false,
		CreatedAt: now,
	}

	if err := s.inviteRepo.Create(ctx, invite); err != nil {
		return nil, fmt.Errorf("invite customer: save: %w", err)
	}

	return invite, nil
}

// ApproveCustomer handles the approve customer use case.
func (s *CustomerService) ApproveCustomer(ctx context.Context, req *command.ApproveCustomerCommand) (*model.Customer, error) {
	customer, err := s.customerRepo.GetByID(ctx, req.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("approve customer: %w", err)
	}

	if err := customer.Approve(); err != nil {
		return nil, fmt.Errorf("approve customer: %w", err)
	}

	if err := s.customerRepo.Update(ctx, customer); err != nil {
		return nil, fmt.Errorf("approve customer: update: %w", err)
	}

	return customer, nil
}

// SuspendCustomer handles the suspend customer use case.
func (s *CustomerService) SuspendCustomer(ctx context.Context, req *command.SuspendCustomerCommand) (*model.Customer, error) {
	customer, err := s.customerRepo.GetByID(ctx, req.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("suspend customer: %w", err)
	}

	if err := customer.Suspend(); err != nil {
		return nil, fmt.Errorf("suspend customer: %w", err)
	}

	if err := s.customerRepo.Update(ctx, customer); err != nil {
		return nil, fmt.Errorf("suspend customer: update: %w", err)
	}

	return customer, nil
}

// GetCustomer handles the get customer query.
func (s *CustomerService) GetCustomer(ctx context.Context, req *command.GetCustomerQuery) (*model.Customer, error) {
	customer, err := s.customerRepo.GetByID(ctx, req.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("get customer: %w", err)
	}
	return customer, nil
}

// ListCustomers handles the list customers query.
func (s *CustomerService) ListCustomers(ctx context.Context, req *command.ListCustomersQuery) ([]*model.Customer, int32, error) {
	customers, total, err := s.customerRepo.List(ctx, req.Status, req.Page, req.PageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("list customers: %w", err)
	}
	return customers, total, nil
}

// generateToken creates a random hex token for invitation links.
func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
