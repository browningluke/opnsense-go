package interfaces

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestLaggMarshal(t *testing.T) {
	lagg := Lagg{
		Members:         api.SelectedMapList{"vtnet2", "vtnet1"},
		PrimaryMember:   api.SelectedMap("vtnet1"),
		Protocol:        api.SelectedMap("lacp"),
		LACPFastTimeout: "1",
		UseFlowID:       api.SelectedMap(""),
		HashLayers:      api.SelectedMapList{"l4", "l2", "l3"},
		LACPStrict:      api.SelectedMap("0"),
		MTU:             "9000",
		Description:     "uplink bundle",
	}

	got, err := json.Marshal(&lagg)
	if err != nil {
		t.Fatalf("marshal LAGG: %v", err)
	}

	want := `{"laggif":"","members":"vtnet1,vtnet2","primary_member":"vtnet1","proto":"lacp","lacp_fast_timeout":"1","use_flowid":"","lagghash":"l2,l3,l4","lacp_strict":"0","mtu":"9000","descr":"uplink bundle"}`
	if string(got) != want {
		t.Fatalf("unexpected LAGG request JSON\ngot:  %s\nwant: %s", got, want)
	}
}

func TestLaggClient(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/interfaces/lagg_settings/addItem", "/api/interfaces/lagg_settings/setItem/abc":
			var body map[string]json.RawMessage
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Errorf("decode request body: %v", err)
			}
			if _, ok := body["lagg"]; !ok {
				t.Errorf("request body missing lagg monad: %s", body)
			}
			_, _ = io.WriteString(w, `{"result":"saved","uuid":"abc"}`)
		case "/api/interfaces/lagg_settings/getItem/abc":
			_, _ = io.WriteString(w, `{"lagg":{"laggif":"lagg0","members":{"vtnet1":{"value":"vtnet1","selected":1}},"primary_member":{},"proto":{"lacp":{"value":"LACP","selected":1}},"lacp_fast_timeout":"0","use_flowid":{},"lagghash":{},"lacp_strict":{},"mtu":"","descr":"test"}}`)
		case "/api/interfaces/lagg_settings/delItem/abc":
			_, _ = io.WriteString(w, `{"result":"deleted"}`)
		case "/api/interfaces/lagg_settings/reconfigure":
			_, _ = io.WriteString(w, `{"status":"ok"}`)
		case "/api/interfaces/lagg_settings/search_item":
			_, _ = io.WriteString(w, `{"current":1,"rowCount":-1,"total":1,"rows":[{"uuid":"abc","laggif":"lagg0","members":"vtnet1","proto":"lacp","descr":"test"}]}`)
		default:
			http.NotFound(w, r)
		}
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	controller := Controller{Api: api.NewClient(api.Options{
		Uri:        server.URL,
		Logger:     log.New(io.Discard, "", 0),
		MaxRetries: 1,
	})}
	ctx := context.Background()
	lagg := &Lagg{Members: api.SelectedMapList{"vtnet1"}, Protocol: api.SelectedMap("lacp")}

	id, err := controller.AddLagg(ctx, lagg)
	if err != nil || id != "abc" {
		t.Fatalf("AddLagg() = %q, %v", id, err)
	}
	got, err := controller.GetLagg(ctx, id)
	if err != nil || got.Device != "lagg0" || !reflect.DeepEqual([]string(got.Members), []string{"vtnet1"}) {
		t.Fatalf("GetLagg() = %+v, %v", got, err)
	}
	if err := controller.UpdateLagg(ctx, id, lagg); err != nil {
		t.Fatalf("UpdateLagg(): %v", err)
	}
	search, err := controller.SearchLagg(ctx, "-1")
	if err != nil || search.Total != 1 || search.Rows[0].Id != id {
		t.Fatalf("SearchLagg() = %+v, %v", search, err)
	}
	if err := controller.DeleteLagg(ctx, id); err != nil {
		t.Fatalf("DeleteLagg(): %v", err)
	}
}

func TestLaggUnmarshal(t *testing.T) {
	response := []byte(`{
		"laggif":"lagg0",
		"members":{
			"vtnet1":{"value":"vtnet1","selected":1},
			"vtnet2":{"value":"vtnet2","selected":true},
			"vtnet3":{"value":"vtnet3","selected":0}
		},
		"primary_member":{
			"vtnet1":{"value":"vtnet1","selected":true},
			"vtnet2":{"value":"vtnet2","selected":false}
		},
		"proto":{"lacp":{"value":"LACP","selected":1}},
		"lacp_fast_timeout":"0",
		"use_flowid":{"":{"value":"Default","selected":1},"1":{"value":"Yes","selected":0},"0":{"value":"No","selected":0}},
		"lagghash":{"l2":{"value":"L2","selected":1},"l3":{"value":"L3","selected":1},"l4":{"value":"L4","selected":0}},
		"lacp_strict":{"":{"value":"Default","selected":1},"1":{"value":"Yes","selected":0},"0":{"value":"No","selected":0}},
		"mtu":"",
		"descr":"uplink bundle"
	}`)

	var got Lagg
	if err := json.Unmarshal(response, &got); err != nil {
		t.Fatalf("unmarshal LAGG: %v", err)
	}

	if got.Device != "lagg0" || got.Protocol.String() != "lacp" || got.PrimaryMember.String() != "vtnet1" {
		t.Fatalf("unexpected scalar fields: %+v", got)
	}
	if !reflect.DeepEqual([]string(got.Members), []string{"vtnet1", "vtnet2"}) {
		t.Fatalf("unexpected members: %#v", got.Members)
	}
	if !reflect.DeepEqual([]string(got.HashLayers), []string{"l2", "l3"}) {
		t.Fatalf("unexpected hash layers: %#v", got.HashLayers)
	}
	if got.UseFlowID.String() != "" || got.LACPStrict.String() != "" || got.MTU != "" {
		t.Fatalf("unexpected default-valued fields: %+v", got)
	}
}

func TestSearchLaggResponseUnmarshal(t *testing.T) {
	response := []byte(`{"current":1,"rowCount":-1,"total":1,"rows":[{"uuid":"abc","laggif":"lagg0","members":"vtnet1,vtnet2","proto":"lacp","descr":"uplink bundle"}]}`)

	var got SearchLaggResponse
	if err := json.Unmarshal(response, &got); err != nil {
		t.Fatalf("unmarshal LAGG search response: %v", err)
	}

	if got.Total != 1 || len(got.Rows) != 1 || got.Rows[0].Id != "abc" || got.Rows[0].Device != "lagg0" {
		t.Fatalf("unexpected LAGG search response: %+v", got)
	}
}
