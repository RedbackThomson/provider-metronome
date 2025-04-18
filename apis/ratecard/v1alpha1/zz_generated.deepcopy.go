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
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CreditTypeConversion) DeepCopyInto(out *CreditTypeConversion) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CreditTypeConversion.
func (in *CreditTypeConversion) DeepCopy() *CreditTypeConversion {
	if in == nil {
		return nil
	}
	out := new(CreditTypeConversion)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FiatCreditType) DeepCopyInto(out *FiatCreditType) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FiatCreditType.
func (in *FiatCreditType) DeepCopy() *FiatCreditType {
	if in == nil {
		return nil
	}
	out := new(FiatCreditType)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ObservedRateCard) DeepCopyInto(out *ObservedRateCard) {
	*out = *in
	out.FiatCreditType = in.FiatCreditType
	if in.Aliases != nil {
		in, out := &in.Aliases, &out.Aliases
		*out = make([]RateCardAlias, len(*in))
		copy(*out, *in)
	}
	if in.CustomFields != nil {
		in, out := &in.CustomFields, &out.CustomFields
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ObservedRateCard.
func (in *ObservedRateCard) DeepCopy() *ObservedRateCard {
	if in == nil {
		return nil
	}
	out := new(ObservedRateCard)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RateCard) DeepCopyInto(out *RateCard) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RateCard.
func (in *RateCard) DeepCopy() *RateCard {
	if in == nil {
		return nil
	}
	out := new(RateCard)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *RateCard) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RateCardAlias) DeepCopyInto(out *RateCardAlias) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RateCardAlias.
func (in *RateCardAlias) DeepCopy() *RateCardAlias {
	if in == nil {
		return nil
	}
	out := new(RateCardAlias)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RateCardList) DeepCopyInto(out *RateCardList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]RateCard, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RateCardList.
func (in *RateCardList) DeepCopy() *RateCardList {
	if in == nil {
		return nil
	}
	out := new(RateCardList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *RateCardList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RateCardParameters) DeepCopyInto(out *RateCardParameters) {
	*out = *in
	if in.CreditTypeConversions != nil {
		in, out := &in.CreditTypeConversions, &out.CreditTypeConversions
		*out = make([]CreditTypeConversion, len(*in))
		copy(*out, *in)
	}
	if in.Aliases != nil {
		in, out := &in.Aliases, &out.Aliases
		*out = make([]RateCardAlias, len(*in))
		copy(*out, *in)
	}
	if in.CustomFields != nil {
		in, out := &in.CustomFields, &out.CustomFields
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RateCardParameters.
func (in *RateCardParameters) DeepCopy() *RateCardParameters {
	if in == nil {
		return nil
	}
	out := new(RateCardParameters)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RateCardSpec) DeepCopyInto(out *RateCardSpec) {
	*out = *in
	in.ResourceSpec.DeepCopyInto(&out.ResourceSpec)
	in.ForProvider.DeepCopyInto(&out.ForProvider)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RateCardSpec.
func (in *RateCardSpec) DeepCopy() *RateCardSpec {
	if in == nil {
		return nil
	}
	out := new(RateCardSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RateCardStatus) DeepCopyInto(out *RateCardStatus) {
	*out = *in
	in.ResourceStatus.DeepCopyInto(&out.ResourceStatus)
	in.AtProvider.DeepCopyInto(&out.AtProvider)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RateCardStatus.
func (in *RateCardStatus) DeepCopy() *RateCardStatus {
	if in == nil {
		return nil
	}
	out := new(RateCardStatus)
	in.DeepCopyInto(out)
	return out
}
