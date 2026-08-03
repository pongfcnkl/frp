package nathole

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	clocktesting "k8s.io/utils/clock/testing"
)

func TestAnalyzerUsesClockForRecordTimestamps(t *testing.T) {
	require := require.New(t)

	start := time.Date(2026, time.May, 8, 12, 30, 0, 0, time.UTC)
	clk := clocktesting.NewFakeClock(start)
	analyzer := newAnalyzerWithClock(time.Hour, clk)
	clientFeature := &NatFeature{NatType: EasyNAT, Behavior: BehaviorNoChange}
	visitorFeature := &NatFeature{NatType: EasyNAT, Behavior: BehaviorNoChange}

	mode, index, _, _ := analyzer.GetRecommandBehaviors("key", clientFeature, visitorFeature)
	require.Equal(start, analyzer.records["key"].lastUpdateTime)

	updatedAt := start.Add(time.Minute)
	clk.SetTime(updatedAt)
	analyzer.ReportSuccess("key", mode, index)
	require.Equal(updatedAt, analyzer.records["key"].lastUpdateTime)

	clk.SetTime(start.Add(2 * time.Hour))
	count, total := analyzer.Clean()
	require.Equal(1, count)
	require.Equal(1, total)
	require.Empty(analyzer.records)
}

func TestHardNATBehaviorsUseLimitedRandomProbes(t *testing.T) {
	for _, tc := range []struct {
		name string
		mode int
	}{
		{name: "mode2", mode: DetectMode2},
		{name: "mode4", mode: DetectMode4},
	} {
		t.Run(tc.name, func(t *testing.T) {
			sender, receiver := getBehaviorByModeAndIndex(tc.mode, 0)

			if sender.PortsRandomNumber != defaultRandomPortProbes {
				t.Fatalf("unexpected random port probes: got %d, want %d", sender.PortsRandomNumber, defaultRandomPortProbes)
			}
			if receiver.ListenRandomPorts != defaultRandomListenPorts {
				t.Fatalf("unexpected random listen ports: got %d, want %d", receiver.ListenRandomPorts, defaultRandomListenPorts)
			}
		})
	}
}

func TestControllerAppliesNatHoleBehaviorOptions(t *testing.T) {
	controller, err := NewController(time.Hour, ControllerOptions{
		RandomPortProbes:  120,
		RandomListenPorts: 32,
	})
	require.NoError(t, err)

	behavior := RecommandBehavior{
		PortsRandomNumber: defaultRandomPortProbes,
		ListenRandomPorts: defaultRandomListenPorts,
	}
	controller.applyBehaviorOptions(&behavior)

	require.Equal(t, 120, behavior.PortsRandomNumber)
	require.Equal(t, 32, behavior.ListenRandomPorts)
}
