package quagga

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

// skipIfUnavailable skips the test when the endpoint doesn't exist on this
// platform (older os-frr versions lack RIP, staticd, and OSPFv3 redistribution).
func skipIfUnavailable(t *testing.T, err error) {
	t.Helper()
	if err != nil && strings.Contains(err.Error(), "404") {
		t.Skipf("endpoint not available on this os-frr version: %v", err)
	}
}

func newController(t *testing.T) Controller {
	t.Helper()
	return Controller{
		Api: api.NewClient(api.Options{
			Uri:           os.Getenv("OPNSENSE_URI"),
			APIKey:        os.Getenv("OPNSENSE_API_KEY"),
			APISecret:     os.Getenv("OPNSENSE_API_SECRET"),
			AllowInsecure: true,
			MaxBackoff:    30,
			MinBackoff:    1,
			MaxRetries:    4,
		}),
	}
}

// ---- General singleton ----

func TestGeneralSingleton(t *testing.T) {
	c := newController(t)
	ctx := context.Background()

	orig, err := c.GeneralGet(ctx)
	if err != nil {
		t.Fatalf("GeneralGet failed: %v", err)
	}

	update := &QuaggaGeneral{
		Enabled:      "0",
		Profile:      orig.General.Profile,
		EnableCarp:   orig.General.EnableCarp,
		EnableSyslog: orig.General.EnableSyslog,
		EnableSNMP:   orig.General.EnableSNMP,
		SyslogLevel:  orig.General.SyslogLevel,
		FWRules:      orig.General.FWRules,
	}
	if _, err := c.GeneralSet(ctx, update); err != nil {
		t.Fatalf("GeneralSet failed: %v", err)
	}

	got, err := c.GeneralGet(ctx)
	if err != nil {
		t.Fatalf("GeneralGet after set failed: %v", err)
	}
	if got.General.Enabled != "0" {
		t.Fatalf("expected enabled=0, got %s", got.General.Enabled)
	}

	// restore
	if _, err := c.GeneralSet(ctx, &orig.General); err != nil {
		t.Fatalf("GeneralSet restore failed: %v", err)
	}
	if _, err := c.ServiceReconfigure(ctx); err != nil {
		t.Fatalf("ServiceReconfigure failed: %v", err)
	}
}

// ---- BFD singleton ----

func TestBFDSingleton(t *testing.T) {
	c := newController(t)
	ctx := context.Background()

	orig, err := c.BFDGet(ctx)
	if err != nil {
		t.Fatalf("BFDGet failed: %v", err)
	}

	if _, err := c.BFDSet(ctx, &QuaggaBFD{Enabled: "1"}); err != nil {
		t.Fatalf("BFDSet failed: %v", err)
	}
	got, err := c.BFDGet(ctx)
	if err != nil {
		t.Fatalf("BFDGet after set failed: %v", err)
	}
	if got.BFD.Enabled != "1" {
		t.Fatalf("expected enabled=1, got %s", got.BFD.Enabled)
	}

	if _, err := c.BFDSet(ctx, &orig.BFD); err != nil {
		t.Fatalf("BFDSet restore failed: %v", err)
	}
	if _, err := c.ServiceReconfigure(ctx); err != nil {
		t.Fatalf("ServiceReconfigure failed: %v", err)
	}
}

// ---- BFD Neighbor ----

func TestBFDNeighbor(t *testing.T) {
	c := newController(t)
	ctx := context.Background()

	nb := &BFDNeighbor{
		Enabled:            "1",
		Description:        "test-bfd-neighbor",
		Address:            "192.0.2.1",
		MultiHop:           "0",
		LocalAddress:       "",
		DetectMultiplier:   "3",
		ReceiveInterval:    "300",
		TransmitInterval:   "300",
	}

	id, err := c.AddBFDNeighbor(ctx, nb)
	if err != nil {
		t.Fatalf("AddBFDNeighbor failed: %v", err)
	}
	t.Logf("Added BFD neighbor: %s", id)

	got, err := c.GetBFDNeighbor(ctx, id)
	if err != nil {
		t.Fatalf("GetBFDNeighbor failed: %v", err)
	}
	if got.Address != nb.Address {
		t.Fatalf("address mismatch: got %s, want %s", got.Address, nb.Address)
	}
	if got.Description != nb.Description {
		t.Fatalf("description mismatch: got %s, want %s", got.Description, nb.Description)
	}

	nb.Description = "test-bfd-neighbor-updated"
	if err := c.UpdateBFDNeighbor(ctx, id, nb); err != nil {
		t.Fatalf("UpdateBFDNeighbor failed: %v", err)
	}

	got, err = c.GetBFDNeighbor(ctx, id)
	if err != nil {
		t.Fatalf("GetBFDNeighbor after update failed: %v", err)
	}
	if got.Description != "test-bfd-neighbor-updated" {
		t.Fatalf("description after update mismatch: got %s", got.Description)
	}

	if err := c.DeleteBFDNeighbor(ctx, id); err != nil {
		t.Fatalf("DeleteBFDNeighbor failed: %v", err)
	}
}

// ---- BGP singleton ----

func TestBGPSingleton(t *testing.T) {
	c := newController(t)
	ctx := context.Background()

	orig, err := c.BGPGet(ctx)
	if err != nil {
		t.Fatalf("BGPGet failed: %v", err)
	}

	update := &QuaggaBGP{
		Enabled:            "0",
		ASNumber:           "65001",
		Distance:           orig.BGP.Distance,
		RouterID:           orig.BGP.RouterID,
		Graceful:           orig.BGP.Graceful,
		NetworkImportCheck: orig.BGP.NetworkImportCheck,
		EnforceFirstAS:     orig.BGP.EnforceFirstAS,
		LogNeighborChanges: orig.BGP.LogNeighborChanges,
		MaximumPaths:       orig.BGP.MaximumPaths,
		MaximumPathsIBGP:   orig.BGP.MaximumPathsIBGP,
	}
	if _, err := c.BGPSet(ctx, update); err != nil {
		t.Fatalf("BGPSet failed: %v", err)
	}

	got, err := c.BGPGet(ctx)
	if err != nil {
		t.Fatalf("BGPGet after set failed: %v", err)
	}
	if got.BGP.ASNumber != "65001" {
		t.Fatalf("asnumber mismatch: got %s", got.BGP.ASNumber)
	}

	if _, err := c.BGPSet(ctx, &orig.BGP); err != nil {
		t.Fatalf("BGPSet restore failed: %v", err)
	}
	if _, err := c.ServiceReconfigure(ctx); err != nil {
		t.Fatalf("ServiceReconfigure failed: %v", err)
	}
}

// ---- BGP PeerGroup ----

func TestBGPPeerGroup(t *testing.T) {
	c := newController(t)
	ctx := context.Background()

	pg := &BGPPeerGroup{
		Enabled:        "1",
		Name:           "test-peergroup",
		RemoteASMode:   api.SelectedMap(""),
		RemoteAS:       "65002",
		Family:         api.SelectedMap("ipv4"),
		NextHopSelf:    "0",
		DefaultOriginate: "0",
	}

	id, err := c.AddBGPPeerGroup(ctx, pg)
	if err != nil {
		t.Fatalf("AddBGPPeerGroup failed: %v", err)
	}
	t.Logf("Added BGP PeerGroup: %s", id)

	got, err := c.GetBGPPeerGroup(ctx, id)
	if err != nil {
		t.Fatalf("GetBGPPeerGroup failed: %v", err)
	}
	if got.Name != pg.Name {
		t.Fatalf("name mismatch: got %s, want %s", got.Name, pg.Name)
	}
	if got.RemoteAS != pg.RemoteAS {
		t.Fatalf("remoteas mismatch: got %s, want %s", got.RemoteAS, pg.RemoteAS)
	}

	pg.RemoteAS = "65003"
	if err := c.UpdateBGPPeerGroup(ctx, id, pg); err != nil {
		t.Fatalf("UpdateBGPPeerGroup failed: %v", err)
	}

	got, err = c.GetBGPPeerGroup(ctx, id)
	if err != nil {
		t.Fatalf("GetBGPPeerGroup after update failed: %v", err)
	}
	if got.RemoteAS != "65003" {
		t.Fatalf("remoteas after update mismatch: got %s", got.RemoteAS)
	}

	if err := c.DeleteBGPPeerGroup(ctx, id); err != nil {
		t.Fatalf("DeleteBGPPeerGroup failed: %v", err)
	}
}

// ---- BGP Redistribution ----

func TestBGPRedistribution(t *testing.T) {
	c := newController(t)
	ctx := context.Background()

	rd := &BGPRedistribution{
		Enabled:      "1",
		Description:  "test-bgp-redistribution",
		Redistribute: api.SelectedMap("connected"),
	}

	id, err := c.AddBGPRedistribution(ctx, rd)
	if err != nil {
		t.Fatalf("AddBGPRedistribution failed: %v", err)
	}
	t.Logf("Added BGP Redistribution: %s", id)

	got, err := c.GetBGPRedistribution(ctx, id)
	if err != nil {
		t.Fatalf("GetBGPRedistribution failed: %v", err)
	}
	if got.Redistribute.String() != "connected" {
		t.Fatalf("redistribute mismatch: got %s", got.Redistribute.String())
	}

	rd.Description = "test-bgp-redistribution-updated"
	if err := c.UpdateBGPRedistribution(ctx, id, rd); err != nil {
		t.Fatalf("UpdateBGPRedistribution failed: %v", err)
	}

	got, err = c.GetBGPRedistribution(ctx, id)
	if err != nil {
		t.Fatalf("GetBGPRedistribution after update failed: %v", err)
	}
	if got.Description != "test-bgp-redistribution-updated" {
		t.Fatalf("description after update: got %s", got.Description)
	}

	if err := c.DeleteBGPRedistribution(ctx, id); err != nil {
		t.Fatalf("DeleteBGPRedistribution failed: %v", err)
	}
}

// ---- BGP ASPath ----

func TestBGPASPath(t *testing.T) {
	c := newController(t)
	ctx := context.Background()

	aspath := &BGPASPath{
		Enabled:     "1",
		Description: "test-aspath",
		Number:      "10",
		Action:      api.SelectedMap("permit"),
		AS:          "^65100$",
	}

	id, err := c.AddBGPASPath(ctx, aspath)
	if err != nil {
		t.Fatalf("AddBGPASPath failed: %v", err)
	}

	got, err := c.GetBGPASPath(ctx, id)
	if err != nil {
		t.Fatalf("GetBGPASPath failed: %v", err)
	}
	if got.AS != aspath.AS {
		t.Fatalf("as mismatch: got %s, want %s", got.AS, aspath.AS)
	}

	aspath.AS = "^65200$"
	if err := c.UpdateBGPASPath(ctx, id, aspath); err != nil {
		t.Fatalf("UpdateBGPASPath failed: %v", err)
	}

	got, err = c.GetBGPASPath(ctx, id)
	if err != nil {
		t.Fatalf("GetBGPASPath after update failed: %v", err)
	}
	if got.AS != "^65200$" {
		t.Fatalf("as after update: got %s", got.AS)
	}

	if err := c.DeleteBGPASPath(ctx, id); err != nil {
		t.Fatalf("DeleteBGPASPath failed: %v", err)
	}
}

// ---- BGP PrefixList ----

func TestBGPPrefixList(t *testing.T) {
	c := newController(t)
	ctx := context.Background()

	pl := &BGPPrefixList{
		Enabled:        "1",
		Description:    "test-prefixlist",
		Name:           "TEST_PL",
		IPVersion:      api.SelectedMap("IPv4"),
		SequenceNumber: "10",
		Action:         api.SelectedMap("permit"),
		Network:        "10.0.0.0/8",
	}

	id, err := c.AddBGPPrefixList(ctx, pl)
	if err != nil {
		t.Fatalf("AddBGPPrefixList failed: %v", err)
	}

	got, err := c.GetBGPPrefixList(ctx, id)
	if err != nil {
		t.Fatalf("GetBGPPrefixList failed: %v", err)
	}
	if got.Name != pl.Name {
		t.Fatalf("name mismatch: got %s, want %s", got.Name, pl.Name)
	}
	if got.Network != pl.Network {
		t.Fatalf("network mismatch: got %s, want %s", got.Network, pl.Network)
	}

	pl.Network = "192.168.0.0/16"
	if err := c.UpdateBGPPrefixList(ctx, id, pl); err != nil {
		t.Fatalf("UpdateBGPPrefixList failed: %v", err)
	}

	got, err = c.GetBGPPrefixList(ctx, id)
	if err != nil {
		t.Fatalf("GetBGPPrefixList after update failed: %v", err)
	}
	if got.Network != "192.168.0.0/16" {
		t.Fatalf("network after update: got %s", got.Network)
	}

	if err := c.DeleteBGPPrefixList(ctx, id); err != nil {
		t.Fatalf("DeleteBGPPrefixList failed: %v", err)
	}
}

// ---- BGP CommunityList ----

func TestBGPCommunityList(t *testing.T) {
	c := newController(t)
	ctx := context.Background()

	cl := &BGPCommunityList{
		Enabled:        "1",
		Description:    "test-communitylist",
		Number:         "10",
		SequenceNumber: "10",
		Action:         api.SelectedMap("permit"),
		Community:      "65000:100",
	}

	id, err := c.AddBGPCommunityList(ctx, cl)
	if err != nil {
		t.Fatalf("AddBGPCommunityList failed: %v", err)
	}

	got, err := c.GetBGPCommunityList(ctx, id)
	if err != nil {
		t.Fatalf("GetBGPCommunityList failed: %v", err)
	}
	if got.Community != cl.Community {
		t.Fatalf("community mismatch: got %s, want %s", got.Community, cl.Community)
	}

	cl.Community = "65000:200"
	if err := c.UpdateBGPCommunityList(ctx, id, cl); err != nil {
		t.Fatalf("UpdateBGPCommunityList failed: %v", err)
	}

	got, err = c.GetBGPCommunityList(ctx, id)
	if err != nil {
		t.Fatalf("GetBGPCommunityList after update failed: %v", err)
	}
	if got.Community != "65000:200" {
		t.Fatalf("community after update: got %s", got.Community)
	}

	if err := c.DeleteBGPCommunityList(ctx, id); err != nil {
		t.Fatalf("DeleteBGPCommunityList failed: %v", err)
	}
}

// ---- BGP RouteMap ----

func TestBGPRouteMap(t *testing.T) {
	c := newController(t)
	ctx := context.Background()

	rm := &BGPRouteMap{
		Enabled:     "1",
		Description: "test-routemap",
		Name:        "TEST_RM",
		Action:      api.SelectedMap("permit"),
		RouteMapID:  "10",
		Set:         "local-preference 200",
	}

	id, err := c.AddBGPRouteMap(ctx, rm)
	if err != nil {
		t.Fatalf("AddBGPRouteMap failed: %v", err)
	}

	got, err := c.GetBGPRouteMap(ctx, id)
	if err != nil {
		t.Fatalf("GetBGPRouteMap failed: %v", err)
	}
	if got.Name != rm.Name {
		t.Fatalf("name mismatch: got %s, want %s", got.Name, rm.Name)
	}

	rm.Set = "local-preference 300"
	if err := c.UpdateBGPRouteMap(ctx, id, rm); err != nil {
		t.Fatalf("UpdateBGPRouteMap failed: %v", err)
	}

	got, err = c.GetBGPRouteMap(ctx, id)
	if err != nil {
		t.Fatalf("GetBGPRouteMap after update failed: %v", err)
	}
	if got.Set != "local-preference 300" {
		t.Fatalf("set after update: got %s", got.Set)
	}

	if err := c.DeleteBGPRouteMap(ctx, id); err != nil {
		t.Fatalf("DeleteBGPRouteMap failed: %v", err)
	}
}

// ---- BGP Neighbor ----

func TestBGPNeighbor(t *testing.T) {
	c := newController(t)
	ctx := context.Background()

	nb := &BGPNeighbor{
		Enabled:     "1",
		Description: "test-bgp-neighbor",
		PeerIP:      "192.0.2.2",
		RemoteASMode: api.SelectedMap(""),
		RemoteAS:    "65100",
	}

	id, err := c.AddBGPNeighbor(ctx, nb)
	if err != nil {
		t.Fatalf("AddBGPNeighbor failed: %v", err)
	}

	got, err := c.GetBGPNeighbor(ctx, id)
	if err != nil {
		t.Fatalf("GetBGPNeighbor failed: %v", err)
	}
	if got.PeerIP != nb.PeerIP {
		t.Fatalf("peer_ip mismatch: got %s, want %s", got.PeerIP, nb.PeerIP)
	}
	if got.RemoteAS != nb.RemoteAS {
		t.Fatalf("remoteas mismatch: got %s, want %s", got.RemoteAS, nb.RemoteAS)
	}

	nb.Description = "test-bgp-neighbor-updated"
	nb.RemoteAS = "65200"
	if err := c.UpdateBGPNeighbor(ctx, id, nb); err != nil {
		t.Fatalf("UpdateBGPNeighbor failed: %v", err)
	}

	got, err = c.GetBGPNeighbor(ctx, id)
	if err != nil {
		t.Fatalf("GetBGPNeighbor after update failed: %v", err)
	}
	if got.RemoteAS != "65200" {
		t.Fatalf("remoteas after update: got %s", got.RemoteAS)
	}

	if err := c.DeleteBGPNeighbor(ctx, id); err != nil {
		t.Fatalf("DeleteBGPNeighbor failed: %v", err)
	}
}

// ---- OSPF singleton ----

func TestOSPFSingleton(t *testing.T) {
	c := newController(t)
	ctx := context.Background()

	orig, err := c.OSPFGet(ctx)
	if err != nil {
		t.Fatalf("OSPFGet failed: %v", err)
	}

	update := &QuaggaOSPF{
		Enabled:             "0",
		CARPDemote:          orig.OSPF.CARPDemote,
		RouterID:            "1.2.3.4",
		CostReference:       orig.OSPF.CostReference,
		LogAdjacencyChanges: orig.OSPF.LogAdjacencyChanges,
		Originate:           orig.OSPF.Originate,
		OriginateAlways:     orig.OSPF.OriginateAlways,
		OriginateMetric:     orig.OSPF.OriginateMetric,
	}
	if _, err := c.OSPFSet(ctx, update); err != nil {
		t.Fatalf("OSPFSet failed: %v", err)
	}

	got, err := c.OSPFGet(ctx)
	if err != nil {
		t.Fatalf("OSPFGet after set failed: %v", err)
	}
	if got.OSPF.RouterID != "1.2.3.4" {
		t.Fatalf("routerid mismatch: got %s", got.OSPF.RouterID)
	}

	if _, err := c.OSPFSet(ctx, &orig.OSPF); err != nil {
		t.Fatalf("OSPFSet restore failed: %v", err)
	}
	if _, err := c.ServiceReconfigure(ctx); err != nil {
		t.Fatalf("ServiceReconfigure failed: %v", err)
	}
}

// ---- OSPF Area ----

func TestOSPFArea(t *testing.T) {
	c := newController(t)
	ctx := context.Background()

	area := &OSPFArea{
		Enabled: "1",
		AreaID:  "0.0.0.1",
		Type:    api.SelectedMap("stub"),
	}

	id, err := c.AddOSPFArea(ctx, area)
	if err != nil {
		t.Fatalf("AddOSPFArea failed: %v", err)
	}
	t.Logf("Added OSPF Area: %s", id)

	got, err := c.GetOSPFArea(ctx, id)
	if err != nil {
		t.Fatalf("GetOSPFArea failed: %v", err)
	}
	if got.AreaID != area.AreaID {
		t.Fatalf("area_id mismatch: got %s, want %s", got.AreaID, area.AreaID)
	}

	area.Type = api.SelectedMap("nssa")
	if err := c.UpdateOSPFArea(ctx, id, area); err != nil {
		t.Fatalf("UpdateOSPFArea failed: %v", err)
	}

	got, err = c.GetOSPFArea(ctx, id)
	if err != nil {
		t.Fatalf("GetOSPFArea after update failed: %v", err)
	}
	if got.Type.String() != "nssa" {
		t.Fatalf("type after update: got %s", got.Type.String())
	}

	if err := c.DeleteOSPFArea(ctx, id); err != nil {
		t.Fatalf("DeleteOSPFArea failed: %v", err)
	}
}

// ---- OSPF Network ----

func TestOSPFNetwork(t *testing.T) {
	c := newController(t)
	ctx := context.Background()

	net := &OSPFNetwork{
		Enabled: "1",
		IPAddr:  "10.0.0.0",
		Area:    "0.0.0.0",
		NetMask: "24",
	}

	id, err := c.AddOSPFNetwork(ctx, net)
	if err != nil {
		t.Fatalf("AddOSPFNetwork failed: %v", err)
	}
	t.Logf("Added OSPF Network: %s", id)

	got, err := c.GetOSPFNetwork(ctx, id)
	if err != nil {
		t.Fatalf("GetOSPFNetwork failed: %v", err)
	}
	if got.IPAddr != net.IPAddr {
		t.Fatalf("ipaddr mismatch: got %s, want %s", got.IPAddr, net.IPAddr)
	}

	net.NetMask = "16"
	if err := c.UpdateOSPFNetwork(ctx, id, net); err != nil {
		t.Fatalf("UpdateOSPFNetwork failed: %v", err)
	}

	got, err = c.GetOSPFNetwork(ctx, id)
	if err != nil {
		t.Fatalf("GetOSPFNetwork after update failed: %v", err)
	}
	if got.NetMask != "16" {
		t.Fatalf("netmask after update: got %s", got.NetMask)
	}

	if err := c.DeleteOSPFNetwork(ctx, id); err != nil {
		t.Fatalf("DeleteOSPFNetwork failed: %v", err)
	}
}

// ---- OSPF PrefixList ----

func TestOSPFPrefixList(t *testing.T) {
	c := newController(t)
	ctx := context.Background()

	pl := &OSPFPrefixList{
		Enabled:        "1",
		Name:           "TEST_OSPF_PL",
		SequenceNumber: "10",
		Action:         api.SelectedMap("permit"),
		Network:        "10.0.0.0/8",
	}

	id, err := c.AddOSPFPrefixList(ctx, pl)
	if err != nil {
		t.Fatalf("AddOSPFPrefixList failed: %v", err)
	}

	got, err := c.GetOSPFPrefixList(ctx, id)
	if err != nil {
		t.Fatalf("GetOSPFPrefixList failed: %v", err)
	}
	if got.Name != pl.Name {
		t.Fatalf("name mismatch: got %s, want %s", got.Name, pl.Name)
	}

	pl.Network = "172.16.0.0/12"
	if err := c.UpdateOSPFPrefixList(ctx, id, pl); err != nil {
		t.Fatalf("UpdateOSPFPrefixList failed: %v", err)
	}

	got, err = c.GetOSPFPrefixList(ctx, id)
	if err != nil {
		t.Fatalf("GetOSPFPrefixList after update failed: %v", err)
	}
	if got.Network != "172.16.0.0/12" {
		t.Fatalf("network after update: got %s", got.Network)
	}

	if err := c.DeleteOSPFPrefixList(ctx, id); err != nil {
		t.Fatalf("DeleteOSPFPrefixList failed: %v", err)
	}
}

// ---- OSPF Redistribution ----

func TestOSPFRedistribution(t *testing.T) {
	c := newController(t)
	ctx := context.Background()

	rd := &OSPFRedistribution{
		Enabled:      "1",
		Description:  "test-ospf-redistribution",
		Redistribute: api.SelectedMap("connected"),
	}

	id, err := c.AddOSPFRedistribution(ctx, rd)
	if err != nil {
		t.Fatalf("AddOSPFRedistribution failed: %v", err)
	}

	got, err := c.GetOSPFRedistribution(ctx, id)
	if err != nil {
		t.Fatalf("GetOSPFRedistribution failed: %v", err)
	}
	if got.Redistribute.String() != "connected" {
		t.Fatalf("redistribute mismatch: got %s", got.Redistribute.String())
	}

	rd.Description = "test-ospf-redistribution-updated"
	if err := c.UpdateOSPFRedistribution(ctx, id, rd); err != nil {
		t.Fatalf("UpdateOSPFRedistribution failed: %v", err)
	}

	got, err = c.GetOSPFRedistribution(ctx, id)
	if err != nil {
		t.Fatalf("GetOSPFRedistribution after update failed: %v", err)
	}
	if got.Description != "test-ospf-redistribution-updated" {
		t.Fatalf("description after update: got %s", got.Description)
	}

	if err := c.DeleteOSPFRedistribution(ctx, id); err != nil {
		t.Fatalf("DeleteOSPFRedistribution failed: %v", err)
	}
}

// ---- OSPF RouteMap ----

func TestOSPFRouteMap(t *testing.T) {
	c := newController(t)
	ctx := context.Background()

	rm := &OSPFRouteMap{
		Enabled:    "1",
		Name:       "TEST_OSPF_RM",
		Action:     api.SelectedMap("permit"),
		RouteMapID: "10",
		Set:        "metric 100",
	}

	id, err := c.AddOSPFRouteMap(ctx, rm)
	if err != nil {
		t.Fatalf("AddOSPFRouteMap failed: %v", err)
	}

	got, err := c.GetOSPFRouteMap(ctx, id)
	if err != nil {
		t.Fatalf("GetOSPFRouteMap failed: %v", err)
	}
	if got.Name != rm.Name {
		t.Fatalf("name mismatch: got %s, want %s", got.Name, rm.Name)
	}

	rm.Set = "metric 200"
	if err := c.UpdateOSPFRouteMap(ctx, id, rm); err != nil {
		t.Fatalf("UpdateOSPFRouteMap failed: %v", err)
	}

	got, err = c.GetOSPFRouteMap(ctx, id)
	if err != nil {
		t.Fatalf("GetOSPFRouteMap after update failed: %v", err)
	}
	if got.Set != "metric 200" {
		t.Fatalf("set after update: got %s", got.Set)
	}

	if err := c.DeleteOSPFRouteMap(ctx, id); err != nil {
		t.Fatalf("DeleteOSPFRouteMap failed: %v", err)
	}
}

// ---- OSPFv3 singleton ----

func TestOSPF6Singleton(t *testing.T) {
	c := newController(t)
	ctx := context.Background()

	orig, err := c.OSPF6Get(ctx)
	if err != nil {
		t.Fatalf("OSPF6Get failed: %v", err)
	}

	update := &QuaggaOSPF6{
		Enabled:         "0",
		CARPDemote:      orig.OSPF6.CARPDemote,
		RouterID:        "4.3.2.1",
		Originate:       orig.OSPF6.Originate,
		OriginateAlways: orig.OSPF6.OriginateAlways,
		OriginateMetric: orig.OSPF6.OriginateMetric,
	}
	if _, err := c.OSPF6Set(ctx, update); err != nil {
		t.Fatalf("OSPF6Set failed: %v", err)
	}

	got, err := c.OSPF6Get(ctx)
	if err != nil {
		t.Fatalf("OSPF6Get after set failed: %v", err)
	}
	if got.OSPF6.RouterID != "4.3.2.1" {
		t.Fatalf("routerid mismatch: got %s", got.OSPF6.RouterID)
	}

	if _, err := c.OSPF6Set(ctx, &orig.OSPF6); err != nil {
		t.Fatalf("OSPF6Set restore failed: %v", err)
	}
	if _, err := c.ServiceReconfigure(ctx); err != nil {
		t.Fatalf("ServiceReconfigure failed: %v", err)
	}
}

// ---- OSPFv3 PrefixList ----

func TestOSPF6PrefixList(t *testing.T) {
	c := newController(t)
	ctx := context.Background()

	pl := &OSPF6PrefixList{
		Enabled:        "1",
		Name:           "TEST_OSPF6_PL",
		SequenceNumber: "10",
		Action:         api.SelectedMap("permit"),
		Network:        "2001:db8::/32",
	}

	id, err := c.AddOSPF6PrefixList(ctx, pl)
	if err != nil {
		t.Fatalf("AddOSPF6PrefixList failed: %v", err)
	}

	got, err := c.GetOSPF6PrefixList(ctx, id)
	if err != nil {
		t.Fatalf("GetOSPF6PrefixList failed: %v", err)
	}
	if got.Name != pl.Name {
		t.Fatalf("name mismatch: got %s, want %s", got.Name, pl.Name)
	}

	pl.Network = "2001:db8:1::/48"
	if err := c.UpdateOSPF6PrefixList(ctx, id, pl); err != nil {
		t.Fatalf("UpdateOSPF6PrefixList failed: %v", err)
	}

	got, err = c.GetOSPF6PrefixList(ctx, id)
	if err != nil {
		t.Fatalf("GetOSPF6PrefixList after update failed: %v", err)
	}
	if got.Network != "2001:db8:1::/48" {
		t.Fatalf("network after update: got %s", got.Network)
	}

	if err := c.DeleteOSPF6PrefixList(ctx, id); err != nil {
		t.Fatalf("DeleteOSPF6PrefixList failed: %v", err)
	}
}

// ---- OSPFv3 Redistribution ----

func TestOSPF6Redistribution(t *testing.T) {
	c := newController(t)
	ctx := context.Background()

	rd := &OSPF6Redistribution{
		Enabled:      "1",
		Description:  "test-ospf6-redistribution",
		Redistribute: api.SelectedMap("connected"),
	}

	id, err := c.AddOSPF6Redistribution(ctx, rd)
	skipIfUnavailable(t, err)
	if err != nil {
		t.Fatalf("AddOSPF6Redistribution failed: %v", err)
	}

	got, err := c.GetOSPF6Redistribution(ctx, id)
	if err != nil {
		t.Fatalf("GetOSPF6Redistribution failed: %v", err)
	}
	if got.Redistribute.String() != "connected" {
		t.Fatalf("redistribute mismatch: got %s", got.Redistribute.String())
	}

	rd.Description = "test-ospf6-redistribution-updated"
	if err := c.UpdateOSPF6Redistribution(ctx, id, rd); err != nil {
		t.Fatalf("UpdateOSPF6Redistribution failed: %v", err)
	}

	got, err = c.GetOSPF6Redistribution(ctx, id)
	if err != nil {
		t.Fatalf("GetOSPF6Redistribution after update failed: %v", err)
	}
	if got.Description != "test-ospf6-redistribution-updated" {
		t.Fatalf("description after update: got %s", got.Description)
	}

	if err := c.DeleteOSPF6Redistribution(ctx, id); err != nil {
		t.Fatalf("DeleteOSPF6Redistribution failed: %v", err)
	}
}

// ---- RIP singleton ----

func TestRIPSingleton(t *testing.T) {
	c := newController(t)
	ctx := context.Background()

	orig, err := c.RIPGet(ctx)
	skipIfUnavailable(t, err)
	if err != nil {
		t.Fatalf("RIPGet failed: %v", err)
	}

	update := &QuaggaRIP{
		Enabled:       "0",
		Version:       "2",
		DefaultMetric: orig.RIP.DefaultMetric,
	}
	if _, err := c.RIPSet(ctx, update); err != nil {
		t.Fatalf("RIPSet failed: %v", err)
	}

	got, err := c.RIPGet(ctx)
	if err != nil {
		t.Fatalf("RIPGet after set failed: %v", err)
	}
	if got.RIP.Version != "2" {
		t.Fatalf("version mismatch: got %s", got.RIP.Version)
	}

	if _, err := c.RIPSet(ctx, &orig.RIP); err != nil {
		t.Fatalf("RIPSet restore failed: %v", err)
	}
	if _, err := c.ServiceReconfigure(ctx); err != nil {
		t.Fatalf("ServiceReconfigure failed: %v", err)
	}
}

// ---- Static singleton ----

func TestStaticSingleton(t *testing.T) {
	c := newController(t)
	ctx := context.Background()

	orig, err := c.StaticGet(ctx)
	skipIfUnavailable(t, err)
	if err != nil {
		t.Fatalf("StaticGet failed: %v", err)
	}

	if _, err := c.StaticSet(ctx, &QuaggaStatic{Enabled: "1"}); err != nil {
		t.Fatalf("StaticSet failed: %v", err)
	}

	got, err := c.StaticGet(ctx)
	if err != nil {
		t.Fatalf("StaticGet after set failed: %v", err)
	}
	if got.Static.Enabled != "1" {
		t.Fatalf("enabled mismatch: got %s", got.Static.Enabled)
	}

	if _, err := c.StaticSet(ctx, &orig.Static); err != nil {
		t.Fatalf("StaticSet restore failed: %v", err)
	}
	if _, err := c.ServiceReconfigure(ctx); err != nil {
		t.Fatalf("ServiceReconfigure failed: %v", err)
	}
}

// ---- Static Route ----

func TestStaticRoute(t *testing.T) {
	c := newController(t)
	ctx := context.Background()

	route := &StaticRoute{
		Enabled: "1",
		Network: "10.100.0.0/24",
		Gateway: "192.168.1.1",
		BFD:     "0",
	}

	id, err := c.AddStaticRoute(ctx, route)
	skipIfUnavailable(t, err)
	if err != nil {
		t.Fatalf("AddStaticRoute failed: %v", err)
	}
	t.Logf("Added Static Route: %s", id)

	got, err := c.GetStaticRoute(ctx, id)
	if err != nil {
		t.Fatalf("GetStaticRoute failed: %v", err)
	}
	if got.Network != route.Network {
		t.Fatalf("network mismatch: got %s, want %s", got.Network, route.Network)
	}
	if got.Gateway != route.Gateway {
		t.Fatalf("gateway mismatch: got %s, want %s", got.Gateway, route.Gateway)
	}

	route.Gateway = "192.168.1.254"
	if err := c.UpdateStaticRoute(ctx, id, route); err != nil {
		t.Fatalf("UpdateStaticRoute failed: %v", err)
	}

	got, err = c.GetStaticRoute(ctx, id)
	if err != nil {
		t.Fatalf("GetStaticRoute after update failed: %v", err)
	}
	if got.Gateway != "192.168.1.254" {
		t.Fatalf("gateway after update: got %s", got.Gateway)
	}

	if err := c.DeleteStaticRoute(ctx, id); err != nil {
		t.Fatalf("DeleteStaticRoute failed: %v", err)
	}
}
