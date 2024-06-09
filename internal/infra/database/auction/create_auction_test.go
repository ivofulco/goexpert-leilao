package auction_test

import (
	"context"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/infra/database/auction"
	"testing"
	"time"
)

func TestCreateAuction(t *testing.T) {
	mockRepo := auction.NewAuctionRepositoryMock()

	auctionEntity := &auction_entity.Auction{
		Id:          "test_auction",
		ProductName: "Test Product",
		Category:    "Test Category",
		Description: "Test Description",
		Condition:   auction_entity.New,
		Status:      auction_entity.Active,
		Timestamp:   time.Now(),
	}

	// Simulando a criação do leilão
	err := mockRepo.CreateAuction(context.Background(), auctionEntity)
	if err != nil {
		t.Fatalf("Failed to create auction: %v", err)
	}

	// Verificar se o leilão foi criado corretamente
	createdAuction, err := mockRepo.FindAuctionById(context.Background(), "test_auction")
	if err != nil {
		t.Fatalf("Failed to find auction: %v", err)
	}

	if createdAuction.ProductName != "Test Product" {
		t.Errorf("Expected product name to be 'Test Product', got %s", createdAuction.ProductName)
	}
}

func TestMonitorAndCloseExpiredAuctions_AlreadyCompleted(t *testing.T) {
	// Configuração inicial
	mockRepo := auction.NewAuctionRepositoryMock()

	// Criar leilão já encerrado
	now := time.Now()
	completedAuction := &auction_entity.Auction{
		Id:          "auction1",
		ProductName: "Completed Product",
		Status:      auction_entity.Completed,
		Timestamp:   now.Add(-time.Hour), // Expirado há mais de uma hora
	}

	mockRepo.SaveAuction(completedAuction)

	// Monitorar e fechar leilões expirados
	mockRepo.MonitorAndCloseExpiredAuctions(context.Background())

	// Verificar se o leilão já encerrado permanece encerrado
	completedAuctionResult, err := mockRepo.FindAuctionById(context.Background(), "auction1")
	if err != nil {
		t.Fatalf("Error finding completed auction: %v", err)
	}
	if completedAuctionResult.Status != auction_entity.Completed {
		t.Errorf("Expected completed auction status to remain 'Completed', got %v", completedAuctionResult.Status)
	}
}

func TestMonitorAndCloseExpiredAuctions_NotExpired(t *testing.T) {
	// Configuração inicial
	mockRepo := auction.NewAuctionRepositoryMock()

	// Criar leilão que não deve ser encerrado
	now := time.Now()
	notExpiredAuction := &auction_entity.Auction{
		Id:          "auction2",
		ProductName: "Not Expired Product",
		Status:      auction_entity.Active,
		Timestamp:   now.Add(-time.Second * 10), // Expirado há menos de 30 minutos
	}

	mockRepo.SaveAuction(notExpiredAuction)

	// Monitorar e fechar leilões expirados
	mockRepo.MonitorAndCloseExpiredAuctions(context.Background())

	// Verificar se o leilão não expirado ainda está ativo
	notExpiredAuctionResult, err := mockRepo.FindAuctionById(context.Background(), "auction2")
	if err != nil {
		t.Fatalf("Error finding not expired auction: %v", err)
	}
	if notExpiredAuctionResult.Status != auction_entity.Active {
		t.Errorf("Expected not expired auction status to be 'Active', got %v", notExpiredAuctionResult.Status)
	}
}

func TestMonitorAndCloseExpiredAuctions_MixedStatus(t *testing.T) {
	// Configuração inicial
	mockRepo := auction.NewAuctionRepositoryMock()

	// Criar alguns leilões simulados com status diferentes
	now := time.Now()
	expiredAuction := &auction_entity.Auction{
		Id:          "auction1",
		ProductName: "Expired Product",
		Status:      auction_entity.Active,
		Timestamp:   now.Add(-time.Minute * 40),
	}
	notExpiredAuction := &auction_entity.Auction{
		Id:          "auction2",
		ProductName: "Not Expired Product",
		Status:      auction_entity.Active,
		Timestamp:   now.Add(-time.Second * 10),
	}
	completedAuction := &auction_entity.Auction{
		Id:          "auction3",
		ProductName: "Completed Product",
		Status:      auction_entity.Completed,
		Timestamp:   now.Add(-time.Hour),
	}

	mockRepo.SaveAuction(expiredAuction)
	mockRepo.SaveAuction(notExpiredAuction)
	mockRepo.SaveAuction(completedAuction)

	// Monitorar e fechar leilões expirados
	mockRepo.MonitorAndCloseExpiredAuctions(context.Background())

	// Verificar se o leilão expirado foi fechado
	expiredAuctionResult, err := mockRepo.FindAuctionById(context.Background(), "auction1")
	if err != nil {
		t.Fatalf("Error finding expired auction: %v", err)
	}
	if expiredAuctionResult.Status != auction_entity.Completed {
		t.Errorf("Expected expired auction status to be 'Completed', got %v", expiredAuctionResult.Status)
	}

	// Verificar se o leilão não expirado ainda está ativo
	notExpiredAuctionResult, err := mockRepo.FindAuctionById(context.Background(), "auction2")
	if err != nil {
		t.Fatalf("Error finding not expired auction: %v", err)
	}
	if notExpiredAuctionResult.Status != auction_entity.Active {
		t.Errorf("Expected not expired auction status to be 'Active', got %v", notExpiredAuctionResult.Status)
	}

	// Verificar se o leilão já encerrado permanece encerrado
	completedAuctionResult, err := mockRepo.FindAuctionById(context.Background(), "auction3")
	if err != nil {
		t.Fatalf("Error finding completed auction: %v", err)
	}
	if completedAuctionResult.Status != auction_entity.Completed {
		t.Errorf("Expected completed auction status to remain 'Completed', got %v", completedAuctionResult.Status)
	}
}
