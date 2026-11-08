package ledger

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

var ErrInvalidManualAsset = errors.New("asset requires a name, type, supported currency, non-negative value, and valid date")

type ManualAsset struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	Name             string    `gorm:"size:120;not null" json:"name"`
	Kind             string    `gorm:"size:12;not null;index" json:"kind"`
	Currency         string    `gorm:"size:3;not null;index" json:"currency"`
	LatestValueMinor int64     `gorm:"->;column:latest_value_minor" json:"latest_value_minor"`
	ValuationDate    string    `gorm:"->;column:valuation_date" json:"valuation_date,omitempty"`
	ValuationSource  string    `gorm:"->;column:valuation_source" json:"valuation_source,omitempty"`
	HasValuation     bool      `gorm:"->;column:has_valuation" json:"has_valuation"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (ManualAsset) TableName() string {
	return "manual_assets"
}

type AssetValuation struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	AssetID    uint      `gorm:"not null;index:idx_asset_valuation_date,priority:1" json:"asset_id"`
	ValueMinor int64     `gorm:"not null" json:"value_minor"`
	Date       string    `gorm:"size:10;not null;index:idx_asset_valuation_date,priority:2" json:"date"`
	Source     string    `gorm:"size:100;not null" json:"source"`
	CreatedAt  time.Time `json:"created_at"`
}

func (AssetValuation) TableName() string {
	return "asset_valuations"
}

type CreateManualAssetRequest struct {
	Name       string `json:"name"`
	Kind       string `json:"kind"`
	Currency   string `json:"currency"`
	ValueMinor int64  `json:"value_minor"`
	Date       string `json:"date"`
	Source     string `json:"source"`
}

type UpdateManualAssetRequest struct {
	Name string `json:"name"`
	Kind string `json:"kind"`
}

func GetManualAssets() ([]ManualAsset, error) {
	return GetManualAssetsAt(time.Now().UTC().Format("2006-01-02"))
}

func GetManualAssetsAt(asOf string) ([]ManualAsset, error) {
	asOfDate, err := time.Parse("2006-01-02", asOf)
	if err != nil || asOfDate.Format("2006-01-02") != asOf {
		return nil, ErrInvalidManualAsset
	}
	assets := make([]ManualAsset, 0)
	err = db.Table("manual_assets").
		Select(`manual_assets.id, manual_assets.name, manual_assets.kind, manual_assets.currency,
			COALESCE(asset_valuations.value_minor, 0) AS latest_value_minor,
			COALESCE(asset_valuations.date, '') AS valuation_date,
			COALESCE(asset_valuations.source, '') AS valuation_source,
			CASE WHEN asset_valuations.id IS NULL THEN false ELSE true END AS has_valuation,
			manual_assets.created_at, manual_assets.updated_at`).
		Joins(`LEFT JOIN asset_valuations ON asset_valuations.id = (
			SELECT id FROM asset_valuations AS latest
			WHERE latest.asset_id = manual_assets.id AND latest.date <= ?
			ORDER BY latest.date DESC, latest.id DESC LIMIT 1
		)`, asOf).
		Where("manual_assets.created_at < ?", asOfDate.AddDate(0, 0, 1)).
		Order("manual_assets.name ASC, manual_assets.id ASC").
		Scan(&assets).Error
	if err != nil {
		return nil, err
	}
	return assets, nil
}

func GetManualAssetValuations(assetID uint) ([]AssetValuation, error) {
	var asset ManualAsset
	if err := db.First(&asset, assetID).Error; err != nil {
		return nil, err
	}
	valuations := make([]AssetValuation, 0)
	err := db.Where("asset_id = ?", assetID).Order("date DESC, id DESC").Find(&valuations).Error
	return valuations, err
}

func CreateManualAsset(request CreateManualAssetRequest) (*ManualAsset, error) {
	asset := ManualAsset{
		Name:     strings.TrimSpace(request.Name),
		Kind:     strings.ToLower(strings.TrimSpace(request.Kind)),
		Currency: strings.ToUpper(strings.TrimSpace(request.Currency)),
	}
	valuation := AssetValuation{
		ValueMinor: request.ValueMinor,
		Date:       request.Date,
		Source:     request.Source,
	}
	if err := validateManualAsset(&asset, &valuation); err != nil {
		return nil, err
	}
	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&asset).Error; err != nil {
			return err
		}
		valuation.AssetID = asset.ID
		return tx.Create(&valuation).Error
	}); err != nil {
		return nil, err
	}
	asset.LatestValueMinor = valuation.ValueMinor
	asset.ValuationDate = valuation.Date
	asset.ValuationSource = valuation.Source
	asset.HasValuation = true
	return &asset, nil
}

func UpdateManualAsset(assetID uint, request UpdateManualAssetRequest) error {
	name := strings.TrimSpace(request.Name)
	kind := strings.ToLower(strings.TrimSpace(request.Kind))
	if name == "" || len(name) > 120 || (kind != "asset" && kind != "liability") {
		return ErrInvalidManualAsset
	}
	var asset ManualAsset
	if err := db.First(&asset, assetID).Error; err != nil {
		return err
	}
	return db.Model(&asset).Select("Name", "Kind").Updates(&ManualAsset{Name: name, Kind: kind}).Error
}

func CreateAssetValuation(assetID uint, valuation *AssetValuation) error {
	valuation.Source = strings.TrimSpace(valuation.Source)
	if valuation.Source == "" {
		valuation.Source = "manual"
	}
	if valuation.ValueMinor < 0 || !validDate(valuation.Date) || len(valuation.Source) > 100 {
		return ErrInvalidManualAsset
	}
	valuation.ID = 0
	valuation.CreatedAt = time.Time{}
	return db.Transaction(func(tx *gorm.DB) error {
		var asset ManualAsset
		if err := tx.First(&asset, assetID).Error; err != nil {
			return err
		}
		valuation.AssetID = assetID
		return tx.Create(valuation).Error
	})
}

func DeleteManualAsset(assetID uint) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&AssetValuation{}, "asset_id = ?", assetID).Error; err != nil {
			return err
		}
		result := tx.Delete(&ManualAsset{}, assetID)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}

func validateManualAsset(asset *ManualAsset, valuation *AssetValuation) error {
	if asset.Name == "" || len(asset.Name) > 120 || (asset.Kind != "asset" && asset.Kind != "liability") ||
		!validCurrency(asset.Currency) || valuation.ValueMinor < 0 || !validDate(valuation.Date) {
		return ErrInvalidManualAsset
	}
	valuation.Source = strings.TrimSpace(valuation.Source)
	if valuation.Source == "" {
		valuation.Source = "manual"
	}
	if len(valuation.Source) > 100 {
		return fmt.Errorf("%w: valuation source must be 100 characters or fewer", ErrInvalidManualAsset)
	}
	return nil
}
