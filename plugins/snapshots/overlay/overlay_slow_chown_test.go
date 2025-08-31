//go:build linux

package overlay

import (
	"context"
	"testing"

	"github.com/containerd/containerd/v2/core/snapshots"
	"github.com/containerd/containerd/v2/internal/userns"
	"github.com/opencontainers/runtime-spec/specs-go"
)

func TestSnapshotterWithSlowChown(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	
	// Create a snapshotter with slow_chown enabled
	s, err := NewSnapshotter(tmpDir, WithSlowChown)
	if err != nil {
		t.Fatalf("Failed to create snapshotter: %v", err)
	}
	defer s.Close()

	// Create a snapshot with user namespace mapping
	ctx := context.Background()
	key := "test-snapshot"
	parent := ""

	// Create ID mappings for user namespace
	uidMaps := []specs.LinuxIDMapping{{ContainerID: 0, HostID: 1000, Size: 65536}}
	gidMaps := []specs.LinuxIDMapping{{ContainerID: 0, HostID: 1000, Size: 65536}}
	
	idMap := userns.IDMap{
		UidMap: uidMaps,
		GidMap: gidMaps,
	}
	uidmapLabel, gidmapLabel := idMap.Marshal()

	opts := []snapshots.Opt{
		snapshots.WithLabels(map[string]string{
			snapshots.LabelSnapshotUIDMapping: uidmapLabel,
			snapshots.LabelSnapshotGIDMapping: gidmapLabel,
		}),
	}

	// Prepare the snapshot
	mounts, err := s.Prepare(ctx, key, parent, opts...)
	if err != nil {
		t.Fatalf("Failed to prepare snapshot: %v", err)
	}

	// Verify that mounts were created
	if len(mounts) == 0 {
		t.Fatal("No mounts returned from Prepare")
	}

	// Clean up
	if err := s.Remove(ctx, key); err != nil {
		t.Fatalf("Failed to remove snapshot: %v", err)
	}
}

func TestSnapshotterWithoutSlowChown(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	
	// Create a snapshotter without slow_chown enabled
	s, err := NewSnapshotter(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create snapshotter: %v", err)
	}
	defer s.Close()

	// Create a snapshot with user namespace mapping
	ctx := context.Background()
	key := "test-snapshot"
	parent := ""

	// Create ID mappings for user namespace
	uidMaps := []specs.LinuxIDMapping{{ContainerID: 0, HostID: 1000, Size: 65536}}
	gidMaps := []specs.LinuxIDMapping{{ContainerID: 0, HostID: 1000, Size: 65536}}
	
	idMap := userns.IDMap{
		UidMap: uidMaps,
		GidMap: gidMaps,
	}
	uidmapLabel, gidmapLabel := idMap.Marshal()

	opts := []snapshots.Opt{
		snapshots.WithLabels(map[string]string{
			snapshots.LabelSnapshotUIDMapping: uidmapLabel,
			snapshots.LabelSnapshotGIDMapping: gidmapLabel,
		}),
	}

	// Prepare the snapshot
	mounts, err := s.Prepare(ctx, key, parent, opts...)
	if err != nil {
		t.Fatalf("Failed to prepare snapshot: %v", err)
	}

	// Verify that mounts were created
	if len(mounts) == 0 {
		t.Fatal("No mounts returned from Prepare")
	}

	// Clean up
	if err := s.Remove(ctx, key); err != nil {
		t.Fatalf("Failed to remove snapshot: %v", err)
	}
}