package db

import (
	"testing"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
)

// TestApplyHostPolicyLeavesDefaultWhenNoDatacenter: with no datacenter the
// cluster keeps gocql's default host-selection policy.
func TestApplyHostPolicyLeavesDefaultWhenNoDatacenter(t *testing.T) {
	cluster := gocql.NewCluster("127.0.0.1")
	before := cluster.PoolConfig.HostSelectionPolicy

	applyHostPolicy(cluster, "")

	if cluster.PoolConfig.HostSelectionPolicy != before {
		t.Fatal("applyHostPolicy changed the policy when no datacenter was given")
	}
}

// TestApplyHostPolicySetsPolicyForDatacenter: a datacenter installs a
// DC-aware host-selection policy.
func TestApplyHostPolicySetsPolicyForDatacenter(t *testing.T) {
	cluster := gocql.NewCluster("127.0.0.1")

	applyHostPolicy(cluster, "dc1")

	if cluster.PoolConfig.HostSelectionPolicy == nil {
		t.Fatal("applyHostPolicy did not set a host-selection policy")
	}
}

// TestApplyHostPolicyBuildsAFreshPolicyEachCall: gocql's token-aware policy may
// not be shared between sessions, so each call must install a new instance -
// this is what lets the protocol-version fallback loop and SetKeyspace recreate
// a session on the same cluster without panicking.
func TestApplyHostPolicyBuildsAFreshPolicyEachCall(t *testing.T) {
	cluster := gocql.NewCluster("127.0.0.1")

	applyHostPolicy(cluster, "dc1")
	first := cluster.PoolConfig.HostSelectionPolicy

	applyHostPolicy(cluster, "dc1")
	second := cluster.PoolConfig.HostSelectionPolicy

	if first == second {
		t.Fatal("applyHostPolicy reused the policy instance; sessions would share it")
	}
}
