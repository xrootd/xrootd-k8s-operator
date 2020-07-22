package utils

import (
	xrootdv1alpha1 "github.com/xrootd/xrootd-k8s-operator/pkg/apis/xrootd/v1alpha1"
)

func NewMemberStatus(ready []string, unready []string) xrootdv1alpha1.MemberStatus {
	size := len(ready) + len(unready)
	return xrootdv1alpha1.MemberStatus{
		Size:    size,
		Ready:   ready,
		Unready: unready,
	}
}
