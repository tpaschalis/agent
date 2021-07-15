// +build !ignore_autogenerated

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CRIStageSpec) DeepCopyInto(out *CRIStageSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CRIStageSpec.
func (in *CRIStageSpec) DeepCopy() *CRIStageSpec {
	if in == nil {
		return nil
	}
	out := new(CRIStageSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DockerStageSpec) DeepCopyInto(out *DockerStageSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DockerStageSpec.
func (in *DockerStageSpec) DeepCopy() *DockerStageSpec {
	if in == nil {
		return nil
	}
	out := new(DockerStageSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DropStageSpec) DeepCopyInto(out *DropStageSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DropStageSpec.
func (in *DropStageSpec) DeepCopy() *DropStageSpec {
	if in == nil {
		return nil
	}
	out := new(DropStageSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GrafanaAgent) DeepCopyInto(out *GrafanaAgent) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GrafanaAgent.
func (in *GrafanaAgent) DeepCopy() *GrafanaAgent {
	if in == nil {
		return nil
	}
	out := new(GrafanaAgent)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *GrafanaAgent) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GrafanaAgentList) DeepCopyInto(out *GrafanaAgentList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]*GrafanaAgent, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(GrafanaAgent)
				(*in).DeepCopyInto(*out)
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GrafanaAgentList.
func (in *GrafanaAgentList) DeepCopy() *GrafanaAgentList {
	if in == nil {
		return nil
	}
	out := new(GrafanaAgentList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *GrafanaAgentList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GrafanaAgentSpec) DeepCopyInto(out *GrafanaAgentSpec) {
	*out = *in
	if in.APIServerConfig != nil {
		in, out := &in.APIServerConfig, &out.APIServerConfig
		*out = new(v1.APIServerConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.PodMetadata != nil {
		in, out := &in.PodMetadata, &out.PodMetadata
		*out = new(v1.EmbeddedObjectMetadata)
		(*in).DeepCopyInto(*out)
	}
	if in.Image != nil {
		in, out := &in.Image, &out.Image
		*out = new(string)
		**out = **in
	}
	if in.ImagePullSecrets != nil {
		in, out := &in.ImagePullSecrets, &out.ImagePullSecrets
		*out = make([]corev1.LocalObjectReference, len(*in))
		copy(*out, *in)
	}
	if in.Storage != nil {
		in, out := &in.Storage, &out.Storage
		*out = new(v1.StorageSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Volumes != nil {
		in, out := &in.Volumes, &out.Volumes
		*out = make([]corev1.Volume, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.VolumeMounts != nil {
		in, out := &in.VolumeMounts, &out.VolumeMounts
		*out = make([]corev1.VolumeMount, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	in.Resources.DeepCopyInto(&out.Resources)
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Secrets != nil {
		in, out := &in.Secrets, &out.Secrets
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.ConfigMaps != nil {
		in, out := &in.ConfigMaps, &out.ConfigMaps
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Affinity != nil {
		in, out := &in.Affinity, &out.Affinity
		*out = new(corev1.Affinity)
		(*in).DeepCopyInto(*out)
	}
	if in.Tolerations != nil {
		in, out := &in.Tolerations, &out.Tolerations
		*out = make([]corev1.Toleration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.TopologySpreadConstraints != nil {
		in, out := &in.TopologySpreadConstraints, &out.TopologySpreadConstraints
		*out = make([]corev1.TopologySpreadConstraint, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.SecurityContext != nil {
		in, out := &in.SecurityContext, &out.SecurityContext
		*out = new(corev1.PodSecurityContext)
		(*in).DeepCopyInto(*out)
	}
	if in.Containers != nil {
		in, out := &in.Containers, &out.Containers
		*out = make([]corev1.Container, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.InitContainers != nil {
		in, out := &in.InitContainers, &out.InitContainers
		*out = make([]corev1.Container, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	in.Prometheus.DeepCopyInto(&out.Prometheus)
	in.Logs.DeepCopyInto(&out.Logs)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GrafanaAgentSpec.
func (in *GrafanaAgentSpec) DeepCopy() *GrafanaAgentSpec {
	if in == nil {
		return nil
	}
	out := new(GrafanaAgentSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JSONStageSpec) DeepCopyInto(out *JSONStageSpec) {
	*out = *in
	if in.Expressions != nil {
		in, out := &in.Expressions, &out.Expressions
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JSONStageSpec.
func (in *JSONStageSpec) DeepCopy() *JSONStageSpec {
	if in == nil {
		return nil
	}
	out := new(JSONStageSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LogsBackoffConfigSpec) DeepCopyInto(out *LogsBackoffConfigSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LogsBackoffConfigSpec.
func (in *LogsBackoffConfigSpec) DeepCopy() *LogsBackoffConfigSpec {
	if in == nil {
		return nil
	}
	out := new(LogsBackoffConfigSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LogsClientSpec) DeepCopyInto(out *LogsClientSpec) {
	*out = *in
	if in.BasicAuth != nil {
		in, out := &in.BasicAuth, &out.BasicAuth
		*out = new(v1.BasicAuth)
		(*in).DeepCopyInto(*out)
	}
	if in.TLSConfig != nil {
		in, out := &in.TLSConfig, &out.TLSConfig
		*out = new(v1.TLSConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.BackoffConfig != nil {
		in, out := &in.BackoffConfig, &out.BackoffConfig
		*out = new(LogsBackoffConfigSpec)
		**out = **in
	}
	if in.ExternalLabels != nil {
		in, out := &in.ExternalLabels, &out.ExternalLabels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LogsClientSpec.
func (in *LogsClientSpec) DeepCopy() *LogsClientSpec {
	if in == nil {
		return nil
	}
	out := new(LogsClientSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LogsInstance) DeepCopyInto(out *LogsInstance) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LogsInstance.
func (in *LogsInstance) DeepCopy() *LogsInstance {
	if in == nil {
		return nil
	}
	out := new(LogsInstance)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *LogsInstance) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LogsInstanceList) DeepCopyInto(out *LogsInstanceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]*LogsInstance, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(LogsInstance)
				(*in).DeepCopyInto(*out)
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LogsInstanceList.
func (in *LogsInstanceList) DeepCopy() *LogsInstanceList {
	if in == nil {
		return nil
	}
	out := new(LogsInstanceList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *LogsInstanceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LogsInstanceSpec) DeepCopyInto(out *LogsInstanceSpec) {
	*out = *in
	if in.Clients != nil {
		in, out := &in.Clients, &out.Clients
		*out = make([]LogsClientSpec, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.PodLogsSelector != nil {
		in, out := &in.PodLogsSelector, &out.PodLogsSelector
		*out = new(metav1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
	if in.PodLogsNamespaceSelector != nil {
		in, out := &in.PodLogsNamespaceSelector, &out.PodLogsNamespaceSelector
		*out = new(metav1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
	if in.AdditionalScrapeConfigs != nil {
		in, out := &in.AdditionalScrapeConfigs, &out.AdditionalScrapeConfigs
		*out = new(corev1.SecretKeySelector)
		(*in).DeepCopyInto(*out)
	}
	if in.TargetConfig != nil {
		in, out := &in.TargetConfig, &out.TargetConfig
		*out = new(LogsTargetConfigSpec)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LogsInstanceSpec.
func (in *LogsInstanceSpec) DeepCopy() *LogsInstanceSpec {
	if in == nil {
		return nil
	}
	out := new(LogsInstanceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LogsSubsystemSpec) DeepCopyInto(out *LogsSubsystemSpec) {
	*out = *in
	if in.Clients != nil {
		in, out := &in.Clients, &out.Clients
		*out = make([]LogsClientSpec, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.LogsExternalLabelName != nil {
		in, out := &in.LogsExternalLabelName, &out.LogsExternalLabelName
		*out = new(string)
		**out = **in
	}
	if in.InstanceSelector != nil {
		in, out := &in.InstanceSelector, &out.InstanceSelector
		*out = new(metav1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
	if in.InstanceNamespaceSelector != nil {
		in, out := &in.InstanceNamespaceSelector, &out.InstanceNamespaceSelector
		*out = new(metav1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LogsSubsystemSpec.
func (in *LogsSubsystemSpec) DeepCopy() *LogsSubsystemSpec {
	if in == nil {
		return nil
	}
	out := new(LogsSubsystemSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LogsTargetConfigSpec) DeepCopyInto(out *LogsTargetConfigSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LogsTargetConfigSpec.
func (in *LogsTargetConfigSpec) DeepCopy() *LogsTargetConfigSpec {
	if in == nil {
		return nil
	}
	out := new(LogsTargetConfigSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MatchStageSpec) DeepCopyInto(out *MatchStageSpec) {
	*out = *in
	if in.Stages != nil {
		in, out := &in.Stages, &out.Stages
		*out = make([]*PipelineStageSpec, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(PipelineStageSpec)
				(*in).DeepCopyInto(*out)
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MatchStageSpec.
func (in *MatchStageSpec) DeepCopy() *MatchStageSpec {
	if in == nil {
		return nil
	}
	out := new(MatchStageSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MetadataConfig) DeepCopyInto(out *MetadataConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MetadataConfig.
func (in *MetadataConfig) DeepCopy() *MetadataConfig {
	if in == nil {
		return nil
	}
	out := new(MetadataConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MetricsStageSpec) DeepCopyInto(out *MetricsStageSpec) {
	*out = *in
	if in.MatchAll != nil {
		in, out := &in.MatchAll, &out.MatchAll
		*out = new(bool)
		**out = **in
	}
	if in.CountEntryBytes != nil {
		in, out := &in.CountEntryBytes, &out.CountEntryBytes
		*out = new(bool)
		**out = **in
	}
	if in.Buckets != nil {
		in, out := &in.Buckets, &out.Buckets
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MetricsStageSpec.
func (in *MetricsStageSpec) DeepCopy() *MetricsStageSpec {
	if in == nil {
		return nil
	}
	out := new(MetricsStageSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MultilineStageSpec) DeepCopyInto(out *MultilineStageSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MultilineStageSpec.
func (in *MultilineStageSpec) DeepCopy() *MultilineStageSpec {
	if in == nil {
		return nil
	}
	out := new(MultilineStageSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OutputStageSpec) DeepCopyInto(out *OutputStageSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OutputStageSpec.
func (in *OutputStageSpec) DeepCopy() *OutputStageSpec {
	if in == nil {
		return nil
	}
	out := new(OutputStageSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PackStageSpec) DeepCopyInto(out *PackStageSpec) {
	*out = *in
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PackStageSpec.
func (in *PackStageSpec) DeepCopy() *PackStageSpec {
	if in == nil {
		return nil
	}
	out := new(PackStageSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PipelineStageSpec) DeepCopyInto(out *PipelineStageSpec) {
	*out = *in
	if in.CRI != nil {
		in, out := &in.CRI, &out.CRI
		*out = new(CRIStageSpec)
		**out = **in
	}
	if in.Docker != nil {
		in, out := &in.Docker, &out.Docker
		*out = new(DockerStageSpec)
		**out = **in
	}
	if in.Drop != nil {
		in, out := &in.Drop, &out.Drop
		*out = new(DropStageSpec)
		**out = **in
	}
	if in.JSON != nil {
		in, out := &in.JSON, &out.JSON
		*out = new(JSONStageSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.LabelAllow != nil {
		in, out := &in.LabelAllow, &out.LabelAllow
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.LabelDrop != nil {
		in, out := &in.LabelDrop, &out.LabelDrop
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Match != nil {
		in, out := &in.Match, &out.Match
		*out = new(MatchStageSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Metrics != nil {
		in, out := &in.Metrics, &out.Metrics
		*out = make(map[string]MetricsStageSpec, len(*in))
		for key, val := range *in {
			(*out)[key] = *val.DeepCopy()
		}
	}
	if in.Multiline != nil {
		in, out := &in.Multiline, &out.Multiline
		*out = new(MultilineStageSpec)
		**out = **in
	}
	if in.Output != nil {
		in, out := &in.Output, &out.Output
		*out = new(OutputStageSpec)
		**out = **in
	}
	if in.Pack != nil {
		in, out := &in.Pack, &out.Pack
		*out = new(PackStageSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Regex != nil {
		in, out := &in.Regex, &out.Regex
		*out = new(RegexStageSpec)
		**out = **in
	}
	if in.Replace != nil {
		in, out := &in.Replace, &out.Replace
		*out = new(ReplaceStageSpec)
		**out = **in
	}
	if in.Template != nil {
		in, out := &in.Template, &out.Template
		*out = new(TemplateStageSpec)
		**out = **in
	}
	if in.Tenant != nil {
		in, out := &in.Tenant, &out.Tenant
		*out = new(TenantStageSpec)
		**out = **in
	}
	if in.Timestamp != nil {
		in, out := &in.Timestamp, &out.Timestamp
		*out = new(TimestampStageSpec)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PipelineStageSpec.
func (in *PipelineStageSpec) DeepCopy() *PipelineStageSpec {
	if in == nil {
		return nil
	}
	out := new(PipelineStageSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PodLogs) DeepCopyInto(out *PodLogs) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PodLogs.
func (in *PodLogs) DeepCopy() *PodLogs {
	if in == nil {
		return nil
	}
	out := new(PodLogs)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PodLogs) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PodLogsList) DeepCopyInto(out *PodLogsList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]*PodLogs, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(PodLogs)
				(*in).DeepCopyInto(*out)
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PodLogsList.
func (in *PodLogsList) DeepCopy() *PodLogsList {
	if in == nil {
		return nil
	}
	out := new(PodLogsList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PodLogsList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PodLogsSpec) DeepCopyInto(out *PodLogsSpec) {
	*out = *in
	if in.PodTargetLabels != nil {
		in, out := &in.PodTargetLabels, &out.PodTargetLabels
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	in.Selector.DeepCopyInto(&out.Selector)
	in.NamespaceSelector.DeepCopyInto(&out.NamespaceSelector)
	if in.PipelineStages != nil {
		in, out := &in.PipelineStages, &out.PipelineStages
		*out = make([]*PipelineStageSpec, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(PipelineStageSpec)
				(*in).DeepCopyInto(*out)
			}
		}
	}
	if in.RelabelConfigs != nil {
		in, out := &in.RelabelConfigs, &out.RelabelConfigs
		*out = make([]*v1.RelabelConfig, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(v1.RelabelConfig)
				(*in).DeepCopyInto(*out)
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PodLogsSpec.
func (in *PodLogsSpec) DeepCopy() *PodLogsSpec {
	if in == nil {
		return nil
	}
	out := new(PodLogsSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PrometheusInstance) DeepCopyInto(out *PrometheusInstance) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PrometheusInstance.
func (in *PrometheusInstance) DeepCopy() *PrometheusInstance {
	if in == nil {
		return nil
	}
	out := new(PrometheusInstance)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PrometheusInstance) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PrometheusInstanceList) DeepCopyInto(out *PrometheusInstanceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]*PrometheusInstance, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(PrometheusInstance)
				(*in).DeepCopyInto(*out)
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PrometheusInstanceList.
func (in *PrometheusInstanceList) DeepCopy() *PrometheusInstanceList {
	if in == nil {
		return nil
	}
	out := new(PrometheusInstanceList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PrometheusInstanceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PrometheusInstanceSpec) DeepCopyInto(out *PrometheusInstanceSpec) {
	*out = *in
	if in.WriteStaleOnShutdown != nil {
		in, out := &in.WriteStaleOnShutdown, &out.WriteStaleOnShutdown
		*out = new(bool)
		**out = **in
	}
	if in.ServiceMonitorSelector != nil {
		in, out := &in.ServiceMonitorSelector, &out.ServiceMonitorSelector
		*out = new(metav1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
	if in.ServiceMonitorNamespaceSelector != nil {
		in, out := &in.ServiceMonitorNamespaceSelector, &out.ServiceMonitorNamespaceSelector
		*out = new(metav1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
	if in.PodMonitorSelector != nil {
		in, out := &in.PodMonitorSelector, &out.PodMonitorSelector
		*out = new(metav1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
	if in.PodMonitorNamespaceSelector != nil {
		in, out := &in.PodMonitorNamespaceSelector, &out.PodMonitorNamespaceSelector
		*out = new(metav1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
	if in.ProbeSelector != nil {
		in, out := &in.ProbeSelector, &out.ProbeSelector
		*out = new(metav1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
	if in.ProbeNamespaceSelector != nil {
		in, out := &in.ProbeNamespaceSelector, &out.ProbeNamespaceSelector
		*out = new(metav1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
	if in.RemoteWrite != nil {
		in, out := &in.RemoteWrite, &out.RemoteWrite
		*out = make([]RemoteWriteSpec, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.AdditionalScrapeConfigs != nil {
		in, out := &in.AdditionalScrapeConfigs, &out.AdditionalScrapeConfigs
		*out = new(corev1.SecretKeySelector)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PrometheusInstanceSpec.
func (in *PrometheusInstanceSpec) DeepCopy() *PrometheusInstanceSpec {
	if in == nil {
		return nil
	}
	out := new(PrometheusInstanceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PrometheusSubsystemSpec) DeepCopyInto(out *PrometheusSubsystemSpec) {
	*out = *in
	if in.RemoteWrite != nil {
		in, out := &in.RemoteWrite, &out.RemoteWrite
		*out = make([]RemoteWriteSpec, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Replicas != nil {
		in, out := &in.Replicas, &out.Replicas
		*out = new(int32)
		**out = **in
	}
	if in.Shards != nil {
		in, out := &in.Shards, &out.Shards
		*out = new(int32)
		**out = **in
	}
	if in.ReplicaExternalLabelName != nil {
		in, out := &in.ReplicaExternalLabelName, &out.ReplicaExternalLabelName
		*out = new(string)
		**out = **in
	}
	if in.PrometheusExternalLabelName != nil {
		in, out := &in.PrometheusExternalLabelName, &out.PrometheusExternalLabelName
		*out = new(string)
		**out = **in
	}
	if in.ExternalLabels != nil {
		in, out := &in.ExternalLabels, &out.ExternalLabels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	out.ArbitraryFSAccessThroughSMs = in.ArbitraryFSAccessThroughSMs
	if in.EnforcedSampleLimit != nil {
		in, out := &in.EnforcedSampleLimit, &out.EnforcedSampleLimit
		*out = new(uint64)
		**out = **in
	}
	if in.EnforcedTargetLimit != nil {
		in, out := &in.EnforcedTargetLimit, &out.EnforcedTargetLimit
		*out = new(uint64)
		**out = **in
	}
	if in.InstanceSelector != nil {
		in, out := &in.InstanceSelector, &out.InstanceSelector
		*out = new(metav1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
	if in.InstanceNamespaceSelector != nil {
		in, out := &in.InstanceNamespaceSelector, &out.InstanceNamespaceSelector
		*out = new(metav1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PrometheusSubsystemSpec.
func (in *PrometheusSubsystemSpec) DeepCopy() *PrometheusSubsystemSpec {
	if in == nil {
		return nil
	}
	out := new(PrometheusSubsystemSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *QueueConfig) DeepCopyInto(out *QueueConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new QueueConfig.
func (in *QueueConfig) DeepCopy() *QueueConfig {
	if in == nil {
		return nil
	}
	out := new(QueueConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RegexStageSpec) DeepCopyInto(out *RegexStageSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RegexStageSpec.
func (in *RegexStageSpec) DeepCopy() *RegexStageSpec {
	if in == nil {
		return nil
	}
	out := new(RegexStageSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RemoteWriteSpec) DeepCopyInto(out *RemoteWriteSpec) {
	*out = *in
	if in.Headers != nil {
		in, out := &in.Headers, &out.Headers
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.WriteRelabelConfigs != nil {
		in, out := &in.WriteRelabelConfigs, &out.WriteRelabelConfigs
		*out = make([]v1.RelabelConfig, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.BasicAuth != nil {
		in, out := &in.BasicAuth, &out.BasicAuth
		*out = new(v1.BasicAuth)
		(*in).DeepCopyInto(*out)
	}
	if in.SigV4 != nil {
		in, out := &in.SigV4, &out.SigV4
		*out = new(SigV4Config)
		(*in).DeepCopyInto(*out)
	}
	if in.TLSConfig != nil {
		in, out := &in.TLSConfig, &out.TLSConfig
		*out = new(v1.TLSConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.QueueConfig != nil {
		in, out := &in.QueueConfig, &out.QueueConfig
		*out = new(QueueConfig)
		**out = **in
	}
	if in.MetadataConfig != nil {
		in, out := &in.MetadataConfig, &out.MetadataConfig
		*out = new(MetadataConfig)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RemoteWriteSpec.
func (in *RemoteWriteSpec) DeepCopy() *RemoteWriteSpec {
	if in == nil {
		return nil
	}
	out := new(RemoteWriteSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ReplaceStageSpec) DeepCopyInto(out *ReplaceStageSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ReplaceStageSpec.
func (in *ReplaceStageSpec) DeepCopy() *ReplaceStageSpec {
	if in == nil {
		return nil
	}
	out := new(ReplaceStageSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SigV4Config) DeepCopyInto(out *SigV4Config) {
	*out = *in
	if in.AccessKey != nil {
		in, out := &in.AccessKey, &out.AccessKey
		*out = new(corev1.SecretKeySelector)
		(*in).DeepCopyInto(*out)
	}
	if in.SecretKey != nil {
		in, out := &in.SecretKey, &out.SecretKey
		*out = new(corev1.SecretKeySelector)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SigV4Config.
func (in *SigV4Config) DeepCopy() *SigV4Config {
	if in == nil {
		return nil
	}
	out := new(SigV4Config)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TemplateStageSpec) DeepCopyInto(out *TemplateStageSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TemplateStageSpec.
func (in *TemplateStageSpec) DeepCopy() *TemplateStageSpec {
	if in == nil {
		return nil
	}
	out := new(TemplateStageSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TenantStageSpec) DeepCopyInto(out *TenantStageSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TenantStageSpec.
func (in *TenantStageSpec) DeepCopy() *TenantStageSpec {
	if in == nil {
		return nil
	}
	out := new(TenantStageSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TimestampStageSpec) DeepCopyInto(out *TimestampStageSpec) {
	*out = *in
	if in.FallbackFormats != nil {
		in, out := &in.FallbackFormats, &out.FallbackFormats
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TimestampStageSpec.
func (in *TimestampStageSpec) DeepCopy() *TimestampStageSpec {
	if in == nil {
		return nil
	}
	out := new(TimestampStageSpec)
	in.DeepCopyInto(out)
	return out
}
