// Copyright 2025 Ant Group Co., Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kusciascheduling

import (
	"strings"

	v1 "k8s.io/api/core/v1"
)

// getJobIDFromPod extracts the job ID from pod annotations or labels
func getJobIDFromPod(pod *v1.Pod) string {
	// First check annotations
	if pod.Annotations != nil {
		if jobID, exists := pod.Annotations[JobIDAnnotationKey]; exists && jobID != "" {
			return jobID
		}
	}

	return ""
}

// isRequireSameNode checks if the KusciaJob requires scheduling to the same node
func isRequireSameNode(pod *v1.Pod) bool {
	// In a real implementation, we would need to fetch the KusciaJob object to check its ScheduleAffinity field.
	// For now, let's assume that if the pod has a job-id annotation, it requires same node scheduling.
	// This is a simplified implementation - in practice, you'd check the actual KusciaJob spec.

	// We're assuming that all pods with job-id annotations need same node scheduling for this example
	// A more complete implementation would check the KusciaJob's ScheduleAffinity field

	// Check for a specific annotation indicating RequireSameNode
	if pod.Annotations != nil {
		if requireSameNode, exists := pod.Annotations[JobAffinityAnnotation]; exists {
			return strings.ToLower(requireSameNode) == "requiresamenode"
		}
	}

	return false
}
