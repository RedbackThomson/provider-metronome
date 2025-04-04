//go:build !ignore_autogenerated

/*
Copyright 2025 RedbackThomson.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"github.com/crossplane/crossplane-runtime/apis/common/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CommitRate) DeepCopyInto(out *CommitRate) {
	*out = *in
	if in.Tiers != nil {
		in, out := &in.Tiers, &out.Tiers
		*out = make([]Tier, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CommitRate.
func (in *CommitRate) DeepCopy() *CommitRate {
	if in == nil {
		return nil
	}
	out := new(CommitRate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CreditType) DeepCopyInto(out *CreditType) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CreditType.
func (in *CreditType) DeepCopy() *CreditType {
	if in == nil {
		return nil
	}
	out := new(CreditType)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ObservedRate) DeepCopyInto(out *ObservedRate) {
	*out = *in
	if in.ProductCustomField != nil {
		in, out := &in.ProductCustomField, &out.ProductCustomField
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.ProductTags != nil {
		in, out := &in.ProductTags, &out.ProductTags
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	in.Details.DeepCopyInto(&out.Details)
	in.CommitRate.DeepCopyInto(&out.CommitRate)
	if in.PricingGroupValues != nil {
		in, out := &in.PricingGroupValues, &out.PricingGroupValues
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ObservedRate.
func (in *ObservedRate) DeepCopy() *ObservedRate {
	if in == nil {
		return nil
	}
	out := new(ObservedRate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Rate) DeepCopyInto(out *Rate) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Rate.
func (in *Rate) DeepCopy() *Rate {
	if in == nil {
		return nil
	}
	out := new(Rate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Rate) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RateDetails) DeepCopyInto(out *RateDetails) {
	*out = *in
	out.CreditType = in.CreditType
	if in.PricingGroupValues != nil {
		in, out := &in.PricingGroupValues, &out.PricingGroupValues
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Tiers != nil {
		in, out := &in.Tiers, &out.Tiers
		*out = make([]Tier, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RateDetails.
func (in *RateDetails) DeepCopy() *RateDetails {
	if in == nil {
		return nil
	}
	out := new(RateDetails)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RateList) DeepCopyInto(out *RateList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Rate, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RateList.
func (in *RateList) DeepCopy() *RateList {
	if in == nil {
		return nil
	}
	out := new(RateList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *RateList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RateParameters) DeepCopyInto(out *RateParameters) {
	*out = *in
	if in.RateCardRef != nil {
		in, out := &in.RateCardRef, &out.RateCardRef
		*out = new(v1.Reference)
		(*in).DeepCopyInto(*out)
	}
	if in.RateCardSelector != nil {
		in, out := &in.RateCardSelector, &out.RateCardSelector
		*out = new(v1.Selector)
		(*in).DeepCopyInto(*out)
	}
	if in.ProductRef != nil {
		in, out := &in.ProductRef, &out.ProductRef
		*out = new(v1.Reference)
		(*in).DeepCopyInto(*out)
	}
	if in.ProductSelector != nil {
		in, out := &in.ProductSelector, &out.ProductSelector
		*out = new(v1.Selector)
		(*in).DeepCopyInto(*out)
	}
	if in.PricingGroupValues != nil {
		in, out := &in.PricingGroupValues, &out.PricingGroupValues
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.CommitRate != nil {
		in, out := &in.CommitRate, &out.CommitRate
		*out = new(CommitRate)
		(*in).DeepCopyInto(*out)
	}
	if in.Tiers != nil {
		in, out := &in.Tiers, &out.Tiers
		*out = make([]Tier, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RateParameters.
func (in *RateParameters) DeepCopy() *RateParameters {
	if in == nil {
		return nil
	}
	out := new(RateParameters)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RateSpec) DeepCopyInto(out *RateSpec) {
	*out = *in
	in.ResourceSpec.DeepCopyInto(&out.ResourceSpec)
	in.ForProvider.DeepCopyInto(&out.ForProvider)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RateSpec.
func (in *RateSpec) DeepCopy() *RateSpec {
	if in == nil {
		return nil
	}
	out := new(RateSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RateStatus) DeepCopyInto(out *RateStatus) {
	*out = *in
	in.ResourceStatus.DeepCopyInto(&out.ResourceStatus)
	in.AtProvider.DeepCopyInto(&out.AtProvider)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RateStatus.
func (in *RateStatus) DeepCopy() *RateStatus {
	if in == nil {
		return nil
	}
	out := new(RateStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Tier) DeepCopyInto(out *Tier) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Tier.
func (in *Tier) DeepCopy() *Tier {
	if in == nil {
		return nil
	}
	out := new(Tier)
	in.DeepCopyInto(out)
	return out
}
