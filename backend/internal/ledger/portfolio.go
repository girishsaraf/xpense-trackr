package ledger

import (
	"errors"
	"math/big"
	"regexp"
	"strings"
	"time"

	"gorm.io/gorm"
)

var (
	ErrInvalidPortfolioHolding = errors.New("invalid portfolio holding or valuation")
	holdingQuantityPattern     = regexp.MustCompile(`^[0-9]{1,20}(\.[0-9]{1,8})?$`)
	holdingTickerPattern       = regexp.MustCompile(`^[A-Za-z0-9.^_-]{0,16}$`)
	holdingAssetClasses        = map[string]struct{}{
		"cash": {}, "equity": {}, "fixed_income": {}, "real_estate": {},
		"commodity": {}, "crypto": {}, "other": {},
	}
)

type PortfolioHolding struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	Name             string    `gorm:"size:120;not null" json:"name"`
	Ticker           string    `gorm:"size:16;not null;default:''" json:"ticker"`
	AssetClass       string    `gorm:"size:20;not null;index" json:"asset_class"`
	Currency         string    `gorm:"size:3;not null;index" json:"currency"`
	Quantity         string    `gorm:"size:32;not null" json:"quantity"`
	CostBasisMinor   int64     `gorm:"not null" json:"cost_basis_minor"`
	LatestValueMinor int64     `gorm:"->;column:latest_value_minor" json:"latest_value_minor"`
	ValuationDate    string    `gorm:"->;column:valuation_date" json:"valuation_date,omitempty"`
	ValuationSource  string    `gorm:"->;column:valuation_source" json:"valuation_source,omitempty"`
	HasValuation     bool      `gorm:"->;column:has_valuation" json:"has_valuation"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (PortfolioHolding) TableName() string {
	return "portfolio_holdings"
}

type HoldingValuation struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	HoldingID  uint      `gorm:"not null;index:idx_holding_valuation_date,priority:1" json:"holding_id"`
	ValueMinor int64     `gorm:"not null" json:"value_minor"`
	Date       string    `gorm:"size:10;not null;index:idx_holding_valuation_date,priority:2" json:"date"`
	Source     string    `gorm:"size:100;not null" json:"source"`
	CreatedAt  time.Time `json:"created_at"`
}

func (HoldingValuation) TableName() string {
	return "holding_valuations"
}

type CreatePortfolioHoldingRequest struct {
	Name           string `json:"name"`
	Ticker         string `json:"ticker"`
	AssetClass     string `json:"asset_class"`
	Currency       string `json:"currency"`
	Quantity       string `json:"quantity"`
	CostBasisMinor int64  `json:"cost_basis_minor"`
	ValueMinor     int64  `json:"value_minor"`
	Date           string `json:"date"`
	Source         string `json:"source"`
}

type UpdatePortfolioHoldingRequest struct {
	Name           string `json:"name"`
	Ticker         string `json:"ticker"`
	AssetClass     string `json:"asset_class"`
	Quantity       string `json:"quantity"`
	CostBasisMinor int64  `json:"cost_basis_minor"`
}

func GetPortfolioHoldings() ([]PortfolioHolding, error) {
	return GetPortfolioHoldingsAt(time.Now().UTC().Format("2006-01-02"))
}

func GetPortfolioHoldingsAt(asOf string) ([]PortfolioHolding, error) {
	asOfDate, err := time.Parse("2006-01-02", asOf)
	if err != nil || asOfDate.Format("2006-01-02") != asOf {
		return nil, ErrInvalidPortfolioHolding
	}
	holdings := make([]PortfolioHolding, 0)
	err = db.Table("portfolio_holdings").
		Select(`portfolio_holdings.id, portfolio_holdings.name, portfolio_holdings.ticker,
			portfolio_holdings.asset_class, portfolio_holdings.currency, portfolio_holdings.quantity,
			portfolio_holdings.cost_basis_minor,
			COALESCE(holding_valuations.value_minor, 0) AS latest_value_minor,
			COALESCE(holding_valuations.date, '') AS valuation_date,
			COALESCE(holding_valuations.source, '') AS valuation_source,
			CASE WHEN holding_valuations.id IS NULL THEN false ELSE true END AS has_valuation,
			portfolio_holdings.created_at, portfolio_holdings.updated_at`).
		Joins(`LEFT JOIN holding_valuations ON holding_valuations.id = (
			SELECT id FROM holding_valuations AS latest
			WHERE latest.holding_id = portfolio_holdings.id AND latest.date <= ?
			ORDER BY latest.date DESC, latest.id DESC LIMIT 1
		)`, asOf).
		Where("portfolio_holdings.created_at < ?", asOfDate.AddDate(0, 0, 1)).
		Order("portfolio_holdings.name ASC, portfolio_holdings.id ASC").
		Scan(&holdings).Error
	return holdings, err
}

func GetHoldingValuations(holdingID uint) ([]HoldingValuation, error) {
	var holding PortfolioHolding
	if err := db.First(&holding, holdingID).Error; err != nil {
		return nil, err
	}
	valuations := make([]HoldingValuation, 0)
	err := db.Where("holding_id = ?", holdingID).Order("date DESC, id DESC").Find(&valuations).Error
	return valuations, err
}

func CreatePortfolioHolding(request CreatePortfolioHoldingRequest) (*PortfolioHolding, error) {
	holding := PortfolioHolding{
		Name:           strings.TrimSpace(request.Name),
		Ticker:         strings.ToUpper(strings.TrimSpace(request.Ticker)),
		AssetClass:     strings.ToLower(strings.TrimSpace(request.AssetClass)),
		Currency:       strings.ToUpper(strings.TrimSpace(request.Currency)),
		Quantity:       strings.TrimSpace(request.Quantity),
		CostBasisMinor: request.CostBasisMinor,
	}
	valuation := HoldingValuation{
		ValueMinor: request.ValueMinor,
		Date:       request.Date,
		Source:     request.Source,
	}
	if err := validatePortfolioHolding(&holding, &valuation); err != nil {
		return nil, err
	}
	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&holding).Error; err != nil {
			return err
		}
		valuation.HoldingID = holding.ID
		return tx.Create(&valuation).Error
	}); err != nil {
		return nil, err
	}
	holding.LatestValueMinor = valuation.ValueMinor
	holding.ValuationDate = valuation.Date
	holding.ValuationSource = valuation.Source
	holding.HasValuation = true
	return &holding, nil
}

func UpdatePortfolioHolding(holdingID uint, request UpdatePortfolioHoldingRequest) error {
	holding := PortfolioHolding{
		Name:           strings.TrimSpace(request.Name),
		Ticker:         strings.ToUpper(strings.TrimSpace(request.Ticker)),
		AssetClass:     strings.ToLower(strings.TrimSpace(request.AssetClass)),
		Quantity:       strings.TrimSpace(request.Quantity),
		CostBasisMinor: request.CostBasisMinor,
	}
	if !validHoldingMetadata(&holding) {
		return ErrInvalidPortfolioHolding
	}
	var existing PortfolioHolding
	if err := db.First(&existing, holdingID).Error; err != nil {
		return err
	}
	return db.Model(&existing).Select("Name", "Ticker", "AssetClass", "Quantity", "CostBasisMinor").Updates(&holding).Error
}

func CreateHoldingValuation(holdingID uint, valuation *HoldingValuation) error {
	valuation.Source = strings.TrimSpace(valuation.Source)
	if valuation.Source == "" {
		valuation.Source = "manual"
	}
	if valuation.ValueMinor < 0 || !validDate(valuation.Date) || len(valuation.Source) > 100 {
		return ErrInvalidPortfolioHolding
	}
	valuation.ID = 0
	valuation.CreatedAt = time.Time{}
	return db.Transaction(func(tx *gorm.DB) error {
		var holding PortfolioHolding
		if err := tx.First(&holding, holdingID).Error; err != nil {
			return err
		}
		valuation.HoldingID = holdingID
		return tx.Create(valuation).Error
	})
}

func DeletePortfolioHolding(holdingID uint) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&HoldingValuation{}, "holding_id = ?", holdingID).Error; err != nil {
			return err
		}
		result := tx.Delete(&PortfolioHolding{}, holdingID)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}

func validatePortfolioHolding(holding *PortfolioHolding, valuation *HoldingValuation) error {
	if !validHoldingMetadata(holding) || !validCurrency(holding.Currency) ||
		valuation.ValueMinor < 0 || !validDate(valuation.Date) {
		return ErrInvalidPortfolioHolding
	}
	valuation.Source = strings.TrimSpace(valuation.Source)
	if valuation.Source == "" {
		valuation.Source = "manual"
	}
	if len(valuation.Source) > 100 {
		return ErrInvalidPortfolioHolding
	}
	return nil
}

func validHoldingMetadata(holding *PortfolioHolding) bool {
	if holding.Name == "" || len(holding.Name) > 120 || !holdingTickerPattern.MatchString(holding.Ticker) ||
		holding.CostBasisMinor < 0 || !holdingQuantityPattern.MatchString(holding.Quantity) {
		return false
	}
	quantity, ok := new(big.Rat).SetString(holding.Quantity)
	if !ok || quantity.Sign() <= 0 {
		return false
	}
	_, classSupported := holdingAssetClasses[holding.AssetClass]
	return classSupported
}
