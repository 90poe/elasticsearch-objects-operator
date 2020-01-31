// +build !ignore_autogenerated

// Code generated by operator-sdk. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ESAlias) DeepCopyInto(out *ESAlias) {
	*out = *in
	if in.Indices != nil {
		in, out := &in.Indices, &out.Indices
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Aliases != nil {
		in, out := &in.Aliases, &out.Aliases
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ESAlias.
func (in *ESAlias) DeepCopy() *ESAlias {
	if in == nil {
		return nil
	}
	out := new(ESAlias)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ESAliases) DeepCopyInto(out *ESAliases) {
	*out = *in
	in.Add.DeepCopyInto(&out.Add)
	in.Remove.DeepCopyInto(&out.Remove)
	in.RemoveIndex.DeepCopyInto(&out.RemoveIndex)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ESAliases.
func (in *ESAliases) DeepCopy() *ESAliases {
	if in == nil {
		return nil
	}
	out := new(ESAliases)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ESAnalyze) DeepCopyInto(out *ESAnalyze) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ESAnalyze.
func (in *ESAnalyze) DeepCopy() *ESAnalyze {
	if in == nil {
		return nil
	}
	out := new(ESAnalyze)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ESHighlights) DeepCopyInto(out *ESHighlights) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ESHighlights.
func (in *ESHighlights) DeepCopy() *ESHighlights {
	if in == nil {
		return nil
	}
	out := new(ESHighlights)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ESIdle) DeepCopyInto(out *ESIdle) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ESIdle.
func (in *ESIdle) DeepCopy() *ESIdle {
	if in == nil {
		return nil
	}
	out := new(ESIdle)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ESIndexBlocks) DeepCopyInto(out *ESIndexBlocks) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ESIndexBlocks.
func (in *ESIndexBlocks) DeepCopy() *ESIndexBlocks {
	if in == nil {
		return nil
	}
	out := new(ESIndexBlocks)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ESIndexRouting) DeepCopyInto(out *ESIndexRouting) {
	*out = *in
	out.Allocation = in.Allocation
	out.Rebalance = in.Rebalance
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ESIndexRouting.
func (in *ESIndexRouting) DeepCopy() *ESIndexRouting {
	if in == nil {
		return nil
	}
	out := new(ESIndexRouting)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ESIndexSettings) DeepCopyInto(out *ESIndexSettings) {
	*out = *in
	out.Shard = in.Shard
	out.SearchIdleAfter = in.SearchIdleAfter
	out.Blocks = in.Blocks
	out.Analyze = in.Analyze
	out.Highlight = in.Highlight
	out.Routing = in.Routing
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ESIndexSettings.
func (in *ESIndexSettings) DeepCopy() *ESIndexSettings {
	if in == nil {
		return nil
	}
	out := new(ESIndexSettings)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ESRoutingAllocationEnable) DeepCopyInto(out *ESRoutingAllocationEnable) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ESRoutingAllocationEnable.
func (in *ESRoutingAllocationEnable) DeepCopy() *ESRoutingAllocationEnable {
	if in == nil {
		return nil
	}
	out := new(ESRoutingAllocationEnable)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ESRoutingRebalanceEnable) DeepCopyInto(out *ESRoutingRebalanceEnable) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ESRoutingRebalanceEnable.
func (in *ESRoutingRebalanceEnable) DeepCopy() *ESRoutingRebalanceEnable {
	if in == nil {
		return nil
	}
	out := new(ESRoutingRebalanceEnable)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ESSearch) DeepCopyInto(out *ESSearch) {
	*out = *in
	out.Idle = in.Idle
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ESSearch.
func (in *ESSearch) DeepCopy() *ESSearch {
	if in == nil {
		return nil
	}
	out := new(ESSearch)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ESShard) DeepCopyInto(out *ESShard) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ESShard.
func (in *ESShard) DeepCopy() *ESShard {
	if in == nil {
		return nil
	}
	out := new(ESShard)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ElasticSearchIndex) DeepCopyInto(out *ElasticSearchIndex) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ElasticSearchIndex.
func (in *ElasticSearchIndex) DeepCopy() *ElasticSearchIndex {
	if in == nil {
		return nil
	}
	out := new(ElasticSearchIndex)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ElasticSearchIndex) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ElasticSearchIndexList) DeepCopyInto(out *ElasticSearchIndexList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ElasticSearchIndex, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ElasticSearchIndexList.
func (in *ElasticSearchIndexList) DeepCopy() *ElasticSearchIndexList {
	if in == nil {
		return nil
	}
	out := new(ElasticSearchIndexList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ElasticSearchIndexList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ElasticSearchIndexSpec) DeepCopyInto(out *ElasticSearchIndexSpec) {
	*out = *in
	out.Settings = in.Settings
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ElasticSearchIndexSpec.
func (in *ElasticSearchIndexSpec) DeepCopy() *ElasticSearchIndexSpec {
	if in == nil {
		return nil
	}
	out := new(ElasticSearchIndexSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ElasticSearchIndexStatus) DeepCopyInto(out *ElasticSearchIndexStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ElasticSearchIndexStatus.
func (in *ElasticSearchIndexStatus) DeepCopy() *ElasticSearchIndexStatus {
	if in == nil {
		return nil
	}
	out := new(ElasticSearchIndexStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ElasticSearchTemplate) DeepCopyInto(out *ElasticSearchTemplate) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ElasticSearchTemplate.
func (in *ElasticSearchTemplate) DeepCopy() *ElasticSearchTemplate {
	if in == nil {
		return nil
	}
	out := new(ElasticSearchTemplate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ElasticSearchTemplate) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ElasticSearchTemplateList) DeepCopyInto(out *ElasticSearchTemplateList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ElasticSearchTemplate, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ElasticSearchTemplateList.
func (in *ElasticSearchTemplateList) DeepCopy() *ElasticSearchTemplateList {
	if in == nil {
		return nil
	}
	out := new(ElasticSearchTemplateList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ElasticSearchTemplateList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ElasticSearchTemplateSpec) DeepCopyInto(out *ElasticSearchTemplateSpec) {
	*out = *in
	if in.IndexPatterns != nil {
		in, out := &in.IndexPatterns, &out.IndexPatterns
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	in.Aliases.DeepCopyInto(&out.Aliases)
	out.Settings = in.Settings
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ElasticSearchTemplateSpec.
func (in *ElasticSearchTemplateSpec) DeepCopy() *ElasticSearchTemplateSpec {
	if in == nil {
		return nil
	}
	out := new(ElasticSearchTemplateSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ElasticSearchTemplateStatus) DeepCopyInto(out *ElasticSearchTemplateStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ElasticSearchTemplateStatus.
func (in *ElasticSearchTemplateStatus) DeepCopy() *ElasticSearchTemplateStatus {
	if in == nil {
		return nil
	}
	out := new(ElasticSearchTemplateStatus)
	in.DeepCopyInto(out)
	return out
}
