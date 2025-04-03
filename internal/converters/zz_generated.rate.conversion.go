// Code generated by github.com/jmattheis/goverter, DO NOT EDIT.
//go:build !ignore_autogenerated

package converters

import (
	v1alpha1 "github.com/redbackthomson/provider-metronome/apis/rate/v1alpha1"
	metronome "github.com/redbackthomson/provider-metronome/internal/clients/metronome"
)

// +k8s:deepcopy-gen=false
type RateConverterImpl struct{}

func (c *RateConverterImpl) FromRate(source *metronome.Rate) *v1alpha1.ObservedRate {
	var pV1alpha1ObservedRate *v1alpha1.ObservedRate
	if source != nil {
		var v1alpha1ObservedRate v1alpha1.ObservedRate
		v1alpha1ObservedRate.Entitled = (*source).Entitled
		if (*source).ProductCustomField != nil {
			v1alpha1ObservedRate.ProductCustomField = make(map[string]string, len((*source).ProductCustomField))
			for key, value := range (*source).ProductCustomField {
				v1alpha1ObservedRate.ProductCustomField[key] = value
			}
		}
		v1alpha1ObservedRate.ProductID = (*source).ProductID
		v1alpha1ObservedRate.ProductName = (*source).ProductName
		if (*source).ProductTags != nil {
			v1alpha1ObservedRate.ProductTags = make([]string, len((*source).ProductTags))
			for i := 0; i < len((*source).ProductTags); i++ {
				v1alpha1ObservedRate.ProductTags[i] = (*source).ProductTags[i]
			}
		}
		v1alpha1ObservedRate.Details = c.metronomeRateDetailsToV1alpha1RateDetails((*source).Details)
		v1alpha1ObservedRate.StartingAt = (*source).StartingAt
		v1alpha1ObservedRate.CommitRate = c.pMetronomeCommitRateToV1alpha1CommitRate((*source).CommitRate)
		v1alpha1ObservedRate.EndingBefore = (*source).EndingBefore
		if (*source).PricingGroupValues != nil {
			v1alpha1ObservedRate.PricingGroupValues = make(map[string]string, len((*source).PricingGroupValues))
			for key2, value2 := range (*source).PricingGroupValues {
				v1alpha1ObservedRate.PricingGroupValues[key2] = value2
			}
		}
		pV1alpha1ObservedRate = &v1alpha1ObservedRate
	}
	return pV1alpha1ObservedRate
}
func (c *RateConverterImpl) FromRateSpec(source *v1alpha1.RateParameters) *metronome.AddRateRequest {
	var pMetronomeAddRateRequest *metronome.AddRateRequest
	if source != nil {
		var metronomeAddRateRequest metronome.AddRateRequest
		metronomeAddRateRequest.Entitled = (*source).Entitled
		metronomeAddRateRequest.ProductID = (*source).ProductID
		metronomeAddRateRequest.RateCardID = (*source).RateCardID
		metronomeAddRateRequest.RateType = (*source).RateType
		metronomeAddRateRequest.StartingAt = (*source).StartingAt
		metronomeAddRateRequest.CommitRate = c.pV1alpha1CommitRateToPMetronomeCommitRate((*source).CommitRate)
		metronomeAddRateRequest.CreditTypeID = (*source).CreditTypeID
		metronomeAddRateRequest.EndingBefore = (*source).EndingBefore
		metronomeAddRateRequest.IsProrated = (*source).IsProrated
		metronomeAddRateRequest.Price = (*source).Price
		if (*source).PricingGroupValues != nil {
			metronomeAddRateRequest.PricingGroupValues = make(map[string]string, len((*source).PricingGroupValues))
			for key, value := range (*source).PricingGroupValues {
				metronomeAddRateRequest.PricingGroupValues[key] = value
			}
		}
		metronomeAddRateRequest.Quantity = (*source).Quantity
		if (*source).Tiers != nil {
			metronomeAddRateRequest.Tiers = make([]metronome.Tier, len((*source).Tiers))
			for i := 0; i < len((*source).Tiers); i++ {
				metronomeAddRateRequest.Tiers[i] = c.v1alpha1TierToMetronomeTier((*source).Tiers[i])
			}
		}
		metronomeAddRateRequest.UseListPrices = (*source).UseListPrices
		pMetronomeAddRateRequest = &metronomeAddRateRequest
	}
	return pMetronomeAddRateRequest
}
func (c *RateConverterImpl) FromRateToParameters(source *metronome.Rate) *v1alpha1.RateParameters {
	var pV1alpha1RateParameters *v1alpha1.RateParameters
	if source != nil {
		var v1alpha1RateParameters v1alpha1.RateParameters
		v1alpha1RateParameters.ProductID = (*source).ProductID
		v1alpha1RateParameters.StartingAt = (*source).StartingAt
		v1alpha1RateParameters.Entitled = (*source).Entitled
		v1alpha1RateParameters.RateType = (*source).Details.RateType
		v1alpha1RateParameters.Price = (*source).Details.Price
		if (*source).PricingGroupValues != nil {
			v1alpha1RateParameters.PricingGroupValues = make(map[string]string, len((*source).PricingGroupValues))
			for key, value := range (*source).PricingGroupValues {
				v1alpha1RateParameters.PricingGroupValues[key] = value
			}
		}
		v1alpha1RateParameters.CommitRate = c.pMetronomeCommitRateToPV1alpha1CommitRate((*source).CommitRate)
		v1alpha1RateParameters.CreditTypeID = (*source).Details.CreditType.ID
		v1alpha1RateParameters.EndingBefore = (*source).EndingBefore
		v1alpha1RateParameters.IsProrated = (*source).Details.IsProrated
		v1alpha1RateParameters.Quantity = (*source).Details.Quantity
		if (*source).Details.Tiers != nil {
			v1alpha1RateParameters.Tiers = make([]v1alpha1.Tier, len((*source).Details.Tiers))
			for i := 0; i < len((*source).Details.Tiers); i++ {
				v1alpha1RateParameters.Tiers[i] = c.metronomeTierToV1alpha1Tier((*source).Details.Tiers[i])
			}
		}
		v1alpha1RateParameters.UseListPrices = (*source).Details.UseListPrices
		pV1alpha1RateParameters = &v1alpha1RateParameters
	}
	return pV1alpha1RateParameters
}
func (c *RateConverterImpl) ToRate(source *v1alpha1.ObservedRate) *metronome.Rate {
	var pMetronomeRate *metronome.Rate
	if source != nil {
		var metronomeRate metronome.Rate
		metronomeRate.Entitled = (*source).Entitled
		if (*source).ProductCustomField != nil {
			metronomeRate.ProductCustomField = make(map[string]string, len((*source).ProductCustomField))
			for key, value := range (*source).ProductCustomField {
				metronomeRate.ProductCustomField[key] = value
			}
		}
		metronomeRate.ProductID = (*source).ProductID
		metronomeRate.ProductName = (*source).ProductName
		if (*source).ProductTags != nil {
			metronomeRate.ProductTags = make([]string, len((*source).ProductTags))
			for i := 0; i < len((*source).ProductTags); i++ {
				metronomeRate.ProductTags[i] = (*source).ProductTags[i]
			}
		}
		metronomeRate.Details = c.v1alpha1RateDetailsToMetronomeRateDetails((*source).Details)
		metronomeRate.StartingAt = (*source).StartingAt
		metronomeRate.CommitRate = c.v1alpha1CommitRateToPMetronomeCommitRate((*source).CommitRate)
		metronomeRate.EndingBefore = (*source).EndingBefore
		if (*source).PricingGroupValues != nil {
			metronomeRate.PricingGroupValues = make(map[string]string, len((*source).PricingGroupValues))
			for key2, value2 := range (*source).PricingGroupValues {
				metronomeRate.PricingGroupValues[key2] = value2
			}
		}
		pMetronomeRate = &metronomeRate
	}
	return pMetronomeRate
}
func (c *RateConverterImpl) ToRateSpec(source *metronome.AddRateRequest) *v1alpha1.RateParameters {
	var pV1alpha1RateParameters *v1alpha1.RateParameters
	if source != nil {
		var v1alpha1RateParameters v1alpha1.RateParameters
		v1alpha1RateParameters.RateCardID = (*source).RateCardID
		v1alpha1RateParameters.ProductID = (*source).ProductID
		v1alpha1RateParameters.StartingAt = (*source).StartingAt
		v1alpha1RateParameters.Entitled = (*source).Entitled
		v1alpha1RateParameters.RateType = (*source).RateType
		v1alpha1RateParameters.Price = (*source).Price
		if (*source).PricingGroupValues != nil {
			v1alpha1RateParameters.PricingGroupValues = make(map[string]string, len((*source).PricingGroupValues))
			for key, value := range (*source).PricingGroupValues {
				v1alpha1RateParameters.PricingGroupValues[key] = value
			}
		}
		v1alpha1RateParameters.CommitRate = c.pMetronomeCommitRateToPV1alpha1CommitRate((*source).CommitRate)
		v1alpha1RateParameters.CreditTypeID = (*source).CreditTypeID
		v1alpha1RateParameters.EndingBefore = (*source).EndingBefore
		v1alpha1RateParameters.IsProrated = (*source).IsProrated
		v1alpha1RateParameters.Quantity = (*source).Quantity
		if (*source).Tiers != nil {
			v1alpha1RateParameters.Tiers = make([]v1alpha1.Tier, len((*source).Tiers))
			for i := 0; i < len((*source).Tiers); i++ {
				v1alpha1RateParameters.Tiers[i] = c.metronomeTierToV1alpha1Tier((*source).Tiers[i])
			}
		}
		v1alpha1RateParameters.UseListPrices = (*source).UseListPrices
		pV1alpha1RateParameters = &v1alpha1RateParameters
	}
	return pV1alpha1RateParameters
}
func (c *RateConverterImpl) metronomeCreditTypeToV1alpha1CreditType(source metronome.CreditType) v1alpha1.CreditType {
	var v1alpha1CreditType v1alpha1.CreditType
	v1alpha1CreditType.ID = source.ID
	v1alpha1CreditType.Name = source.Name
	return v1alpha1CreditType
}
func (c *RateConverterImpl) metronomeRateDetailsToV1alpha1RateDetails(source metronome.RateDetails) v1alpha1.RateDetails {
	var v1alpha1RateDetails v1alpha1.RateDetails
	v1alpha1RateDetails.RateType = source.RateType
	v1alpha1RateDetails.CreditType = c.metronomeCreditTypeToV1alpha1CreditType(source.CreditType)
	v1alpha1RateDetails.IsProrated = source.IsProrated
	v1alpha1RateDetails.Price = source.Price
	if source.PricingGroupValues != nil {
		v1alpha1RateDetails.PricingGroupValues = make(map[string]string, len(source.PricingGroupValues))
		for key, value := range source.PricingGroupValues {
			v1alpha1RateDetails.PricingGroupValues[key] = value
		}
	}
	v1alpha1RateDetails.Quantity = source.Quantity
	if source.Tiers != nil {
		v1alpha1RateDetails.Tiers = make([]v1alpha1.Tier, len(source.Tiers))
		for i := 0; i < len(source.Tiers); i++ {
			v1alpha1RateDetails.Tiers[i] = c.metronomeTierToV1alpha1Tier(source.Tiers[i])
		}
	}
	v1alpha1RateDetails.UseListPrices = source.UseListPrices
	return v1alpha1RateDetails
}
func (c *RateConverterImpl) metronomeTierToV1alpha1Tier(source metronome.Tier) v1alpha1.Tier {
	var v1alpha1Tier v1alpha1.Tier
	v1alpha1Tier.Price = source.Price
	v1alpha1Tier.Size = source.Size
	return v1alpha1Tier
}
func (c *RateConverterImpl) pMetronomeCommitRateToPV1alpha1CommitRate(source *metronome.CommitRate) *v1alpha1.CommitRate {
	var pV1alpha1CommitRate *v1alpha1.CommitRate
	if source != nil {
		var v1alpha1CommitRate v1alpha1.CommitRate
		v1alpha1CommitRate.RateType = (*source).RateType
		v1alpha1CommitRate.Price = (*source).Price
		if (*source).Tiers != nil {
			v1alpha1CommitRate.Tiers = make([]v1alpha1.Tier, len((*source).Tiers))
			for i := 0; i < len((*source).Tiers); i++ {
				v1alpha1CommitRate.Tiers[i] = c.metronomeTierToV1alpha1Tier((*source).Tiers[i])
			}
		}
		pV1alpha1CommitRate = &v1alpha1CommitRate
	}
	return pV1alpha1CommitRate
}
func (c *RateConverterImpl) pMetronomeCommitRateToV1alpha1CommitRate(source *metronome.CommitRate) v1alpha1.CommitRate {
	var v1alpha1CommitRate v1alpha1.CommitRate
	if source != nil {
		var v1alpha1CommitRate2 v1alpha1.CommitRate
		v1alpha1CommitRate2.RateType = (*source).RateType
		v1alpha1CommitRate2.Price = (*source).Price
		if (*source).Tiers != nil {
			v1alpha1CommitRate2.Tiers = make([]v1alpha1.Tier, len((*source).Tiers))
			for i := 0; i < len((*source).Tiers); i++ {
				v1alpha1CommitRate2.Tiers[i] = c.metronomeTierToV1alpha1Tier((*source).Tiers[i])
			}
		}
		v1alpha1CommitRate = v1alpha1CommitRate2
	}
	return v1alpha1CommitRate
}
func (c *RateConverterImpl) pV1alpha1CommitRateToPMetronomeCommitRate(source *v1alpha1.CommitRate) *metronome.CommitRate {
	var pMetronomeCommitRate *metronome.CommitRate
	if source != nil {
		var metronomeCommitRate metronome.CommitRate
		metronomeCommitRate.RateType = (*source).RateType
		metronomeCommitRate.Price = (*source).Price
		if (*source).Tiers != nil {
			metronomeCommitRate.Tiers = make([]metronome.Tier, len((*source).Tiers))
			for i := 0; i < len((*source).Tiers); i++ {
				metronomeCommitRate.Tiers[i] = c.v1alpha1TierToMetronomeTier((*source).Tiers[i])
			}
		}
		pMetronomeCommitRate = &metronomeCommitRate
	}
	return pMetronomeCommitRate
}
func (c *RateConverterImpl) v1alpha1CommitRateToPMetronomeCommitRate(source v1alpha1.CommitRate) *metronome.CommitRate {
	var metronomeCommitRate metronome.CommitRate
	metronomeCommitRate.RateType = source.RateType
	metronomeCommitRate.Price = source.Price
	if source.Tiers != nil {
		metronomeCommitRate.Tiers = make([]metronome.Tier, len(source.Tiers))
		for i := 0; i < len(source.Tiers); i++ {
			metronomeCommitRate.Tiers[i] = c.v1alpha1TierToMetronomeTier(source.Tiers[i])
		}
	}
	return &metronomeCommitRate
}
func (c *RateConverterImpl) v1alpha1CreditTypeToMetronomeCreditType(source v1alpha1.CreditType) metronome.CreditType {
	var metronomeCreditType metronome.CreditType
	metronomeCreditType.ID = source.ID
	metronomeCreditType.Name = source.Name
	return metronomeCreditType
}
func (c *RateConverterImpl) v1alpha1RateDetailsToMetronomeRateDetails(source v1alpha1.RateDetails) metronome.RateDetails {
	var metronomeRateDetails metronome.RateDetails
	metronomeRateDetails.RateType = source.RateType
	metronomeRateDetails.CreditType = c.v1alpha1CreditTypeToMetronomeCreditType(source.CreditType)
	metronomeRateDetails.IsProrated = source.IsProrated
	metronomeRateDetails.Price = source.Price
	if source.PricingGroupValues != nil {
		metronomeRateDetails.PricingGroupValues = make(map[string]string, len(source.PricingGroupValues))
		for key, value := range source.PricingGroupValues {
			metronomeRateDetails.PricingGroupValues[key] = value
		}
	}
	metronomeRateDetails.Quantity = source.Quantity
	if source.Tiers != nil {
		metronomeRateDetails.Tiers = make([]metronome.Tier, len(source.Tiers))
		for i := 0; i < len(source.Tiers); i++ {
			metronomeRateDetails.Tiers[i] = c.v1alpha1TierToMetronomeTier(source.Tiers[i])
		}
	}
	metronomeRateDetails.UseListPrices = source.UseListPrices
	return metronomeRateDetails
}
func (c *RateConverterImpl) v1alpha1TierToMetronomeTier(source v1alpha1.Tier) metronome.Tier {
	var metronomeTier metronome.Tier
	metronomeTier.Price = source.Price
	metronomeTier.Size = source.Size
	return metronomeTier
}
