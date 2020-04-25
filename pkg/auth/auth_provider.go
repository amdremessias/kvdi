package auth

import (
	"github.com/tinyzimmer/kvdi/pkg/apis/kvdi/v1alpha1"
	"github.com/tinyzimmer/kvdi/pkg/auth/providers/local"
	"github.com/tinyzimmer/kvdi/pkg/auth/types"
)

func GetAuthProvider(cluster *v1alpha1.VDICluster) types.AuthProvider {
	return local.New()
}
