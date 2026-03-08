package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// --- ListTemplates ---

func TestTemplateHandler_ListTemplates_Empty(t *testing.T) {
	env := newTestEnv(t)

	req := httptest.NewRequest("GET", "/api/admin/templates", nil)
	rec := httptest.NewRecorder()
	env.TemplateHandler.ListTemplates(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", rec.Code, http.StatusOK)
	}

	data, errMsg := decodeEnvelope(t, rec)
	if errMsg != "" {
		t.Fatalf("unexpected error: %s", errMsg)
	}

	var resp struct {
		Templates []json.RawMessage `json:"templates"`
	}
	json.Unmarshal(data, &resp)
	if len(resp.Templates) != 0 {
		t.Errorf("template count: got %d, want 0", len(resp.Templates))
	}
}

// --- CreateTemplate ---

func TestTemplateHandler_CreateTemplate_Success(t *testing.T) {
	env := newTestEnv(t)

	body := `{"name":"test-map","description":"A test map","map_size":11}`
	req := httptest.NewRequest("POST", "/api/admin/templates", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	env.TemplateHandler.CreateTemplate(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status: got %d, want %d. Body: %s", rec.Code, http.StatusCreated, rec.Body.String())
	}

	data, errMsg := decodeEnvelope(t, rec)
	if errMsg != "" {
		t.Fatalf("unexpected error: %s", errMsg)
	}

	var resp struct {
		Message  string         `json:"message"`
		Template map[string]any `json:"template"`
	}
	json.Unmarshal(data, &resp)
	if resp.Template["name"] != "test-map" {
		t.Errorf("name: got %v, want test-map", resp.Template["name"])
	}
}

func TestTemplateHandler_CreateTemplate_MissingName(t *testing.T) {
	env := newTestEnv(t)

	body := `{"description":"no name"}`
	req := httptest.NewRequest("POST", "/api/admin/templates", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	env.TemplateHandler.CreateTemplate(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestTemplateHandler_CreateTemplate_Duplicate(t *testing.T) {
	env := newTestEnv(t)

	body := `{"name":"dup-map"}`
	req := httptest.NewRequest("POST", "/api/admin/templates", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	env.TemplateHandler.CreateTemplate(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("first create failed: %d %s", rec.Code, rec.Body.String())
	}

	// Duplicate
	req2 := httptest.NewRequest("POST", "/api/admin/templates", strings.NewReader(body))
	req2.Header.Set("Content-Type", "application/json")
	rec2 := httptest.NewRecorder()
	env.TemplateHandler.CreateTemplate(rec2, req2)

	if rec2.Code != http.StatusConflict {
		t.Errorf("status: got %d, want %d", rec2.Code, http.StatusConflict)
	}
}

// --- GetTemplate ---

func TestTemplateHandler_GetTemplate_Success(t *testing.T) {
	env := newTestEnv(t)

	// Create first
	createBody := `{"name":"get-me","map_size":5}`
	createReq := httptest.NewRequest("POST", "/api/admin/templates", strings.NewReader(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()
	env.TemplateHandler.CreateTemplate(createRec, createReq)

	// Get
	req := httptest.NewRequest("GET", "/api/admin/templates/get-me", nil)
	req.SetPathValue("name", "get-me")
	rec := httptest.NewRecorder()
	env.TemplateHandler.GetTemplate(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d. Body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}
}

func TestTemplateHandler_GetTemplate_NotFound(t *testing.T) {
	env := newTestEnv(t)

	req := httptest.NewRequest("GET", "/api/admin/templates/nope", nil)
	req.SetPathValue("name", "nope")
	rec := httptest.NewRecorder()
	env.TemplateHandler.GetTemplate(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusNotFound)
	}
}

// --- DeleteTemplate ---

func TestTemplateHandler_DeleteTemplate_Success(t *testing.T) {
	env := newTestEnv(t)

	// Create
	createBody := `{"name":"del-me"}`
	createReq := httptest.NewRequest("POST", "/api/admin/templates", strings.NewReader(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()
	env.TemplateHandler.CreateTemplate(createRec, createReq)

	// Delete
	req := httptest.NewRequest("DELETE", "/api/admin/templates/del-me", nil)
	req.SetPathValue("name", "del-me")
	rec := httptest.NewRecorder()
	env.TemplateHandler.DeleteTemplate(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", rec.Code, http.StatusOK)
	}

	// Verify gone
	getReq := httptest.NewRequest("GET", "/api/admin/templates/del-me", nil)
	getReq.SetPathValue("name", "del-me")
	getRec := httptest.NewRecorder()
	env.TemplateHandler.GetTemplate(getRec, getReq)

	if getRec.Code != http.StatusNotFound {
		t.Errorf("expected 404 after delete, got %d", getRec.Code)
	}
}

func TestTemplateHandler_DeleteTemplate_NotFound(t *testing.T) {
	env := newTestEnv(t)

	req := httptest.NewRequest("DELETE", "/api/admin/templates/ghost", nil)
	req.SetPathValue("name", "ghost")
	rec := httptest.NewRecorder()
	env.TemplateHandler.DeleteTemplate(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusNotFound)
	}
}

// --- UpdateTerrain ---

func TestTemplateHandler_UpdateTerrain_Success(t *testing.T) {
	env := newTestEnv(t)

	// Create template
	createBody := `{"name":"terrain-test","map_size":5}`
	createReq := httptest.NewRequest("POST", "/api/admin/templates", strings.NewReader(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()
	env.TemplateHandler.CreateTemplate(createRec, createReq)

	// Update terrain
	body := `{"tiles":[{"x":0,"y":0,"terrain_type":"mountain"},{"x":1,"y":0,"terrain_type":"forest"}]}`
	req := httptest.NewRequest("PUT", "/api/admin/templates/terrain-test/terrain", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("name", "terrain-test")
	rec := httptest.NewRecorder()
	env.TemplateHandler.UpdateTerrain(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d. Body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}
}

func TestTemplateHandler_UpdateTerrain_InvalidType(t *testing.T) {
	env := newTestEnv(t)

	createBody := `{"name":"bad-terrain","map_size":5}`
	createReq := httptest.NewRequest("POST", "/api/admin/templates", strings.NewReader(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()
	env.TemplateHandler.CreateTemplate(createRec, createReq)

	body := `{"tiles":[{"x":0,"y":0,"terrain_type":"lava"}]}`
	req := httptest.NewRequest("PUT", "/api/admin/templates/bad-terrain/terrain", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("name", "bad-terrain")
	rec := httptest.NewRecorder()
	env.TemplateHandler.UpdateTerrain(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

// --- UpdateZones ---

func TestTemplateHandler_UpdateZones_Success(t *testing.T) {
	env := newTestEnv(t)

	createBody := `{"name":"zone-test","map_size":5}`
	createReq := httptest.NewRequest("POST", "/api/admin/templates", strings.NewReader(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()
	env.TemplateHandler.CreateTemplate(createRec, createReq)

	body := `{"tiles":[{"x":0,"y":0,"kingdom_zone":"veridor"},{"x":1,"y":0,"kingdom_zone":"sylvara"}]}`
	req := httptest.NewRequest("PUT", "/api/admin/templates/zone-test/zones", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("name", "zone-test")
	rec := httptest.NewRecorder()
	env.TemplateHandler.UpdateZones(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d. Body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}
}

func TestTemplateHandler_UpdateZones_InvalidZone(t *testing.T) {
	env := newTestEnv(t)

	createBody := `{"name":"bad-zone","map_size":5}`
	createReq := httptest.NewRequest("POST", "/api/admin/templates", strings.NewReader(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()
	env.TemplateHandler.CreateTemplate(createRec, createReq)

	body := `{"tiles":[{"x":0,"y":0,"kingdom_zone":"atlantis"}]}`
	req := httptest.NewRequest("PUT", "/api/admin/templates/bad-zone/zones", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("name", "bad-zone")
	rec := httptest.NewRecorder()
	env.TemplateHandler.UpdateZones(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

// --- ResizeTemplate ---

func TestTemplateHandler_ResizeTemplate_Success(t *testing.T) {
	env := newTestEnv(t)

	createBody := `{"name":"resize-test","map_size":5}`
	createReq := httptest.NewRequest("POST", "/api/admin/templates", strings.NewReader(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()
	env.TemplateHandler.CreateTemplate(createRec, createReq)

	body := `{"map_size":11}`
	req := httptest.NewRequest("PUT", "/api/admin/templates/resize-test/resize", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("name", "resize-test")
	rec := httptest.NewRecorder()
	env.TemplateHandler.ResizeTemplate(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d. Body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}
}

func TestTemplateHandler_ResizeTemplate_TooSmall(t *testing.T) {
	env := newTestEnv(t)

	createBody := `{"name":"tiny-resize","map_size":5}`
	createReq := httptest.NewRequest("POST", "/api/admin/templates", strings.NewReader(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()
	env.TemplateHandler.CreateTemplate(createRec, createReq)

	body := `{"map_size":1}`
	req := httptest.NewRequest("PUT", "/api/admin/templates/tiny-resize/resize", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("name", "tiny-resize")
	rec := httptest.NewRecorder()
	env.TemplateHandler.ResizeTemplate(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

// --- ApplyTemplate ---

func TestTemplateHandler_ApplyTemplate_Success(t *testing.T) {
	env := newTestEnv(t)

	// Create 51x51 template (must match default map size)
	createBody := `{"name":"apply-test","map_size":51}`
	createReq := httptest.NewRequest("POST", "/api/admin/templates", strings.NewReader(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()
	env.TemplateHandler.CreateTemplate(createRec, createReq)

	// Apply
	body := `{"confirm":true}`
	req := httptest.NewRequest("POST", "/api/admin/templates/apply-test/apply", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("name", "apply-test")
	req = req.WithContext(authCtx(1, "admin"))
	rec := httptest.NewRecorder()
	env.TemplateHandler.ApplyTemplate(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d. Body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}
}

func TestTemplateHandler_ApplyTemplate_NoConfirm(t *testing.T) {
	env := newTestEnv(t)

	createBody := `{"name":"no-confirm"}`
	createReq := httptest.NewRequest("POST", "/api/admin/templates", strings.NewReader(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()
	env.TemplateHandler.CreateTemplate(createRec, createReq)

	body := `{"confirm":false}`
	req := httptest.NewRequest("POST", "/api/admin/templates/no-confirm/apply", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("name", "no-confirm")
	rec := httptest.NewRecorder()
	env.TemplateHandler.ApplyTemplate(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestTemplateHandler_ApplyTemplate_NotFound(t *testing.T) {
	env := newTestEnv(t)

	body := `{"confirm":true}`
	req := httptest.NewRequest("POST", "/api/admin/templates/missing/apply", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("name", "missing")
	rec := httptest.NewRecorder()
	env.TemplateHandler.ApplyTemplate(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}

// --- ExportTemplate ---

func TestTemplateHandler_ExportTemplate_Success(t *testing.T) {
	env := newTestEnv(t)

	createBody := `{"name":"export-test","map_size":5}`
	createReq := httptest.NewRequest("POST", "/api/admin/templates", strings.NewReader(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()
	env.TemplateHandler.CreateTemplate(createRec, createReq)

	req := httptest.NewRequest("GET", "/api/admin/templates/export-test/export", nil)
	req.SetPathValue("name", "export-test")
	rec := httptest.NewRecorder()
	env.TemplateHandler.ExportTemplate(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", rec.Code, http.StatusOK)
	}

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("content-type: got %q, want application/json", ct)
	}

	cd := rec.Header().Get("Content-Disposition")
	if !strings.Contains(cd, "export-test.json") {
		t.Errorf("content-disposition: got %q, want to contain export-test.json", cd)
	}
}

// --- ListTemplates after create ---

func TestTemplateHandler_ListTemplates_AfterCreate(t *testing.T) {
	env := newTestEnv(t)

	// Create two templates
	for _, name := range []string{"list-a", "list-b"} {
		body := `{"name":"` + name + `"}`
		req := httptest.NewRequest("POST", "/api/admin/templates", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		env.TemplateHandler.CreateTemplate(rec, req)
	}

	req := httptest.NewRequest("GET", "/api/admin/templates", nil)
	rec := httptest.NewRecorder()
	env.TemplateHandler.ListTemplates(rec, req)

	data, _ := decodeEnvelope(t, rec)
	var resp struct {
		Templates []json.RawMessage `json:"templates"`
	}
	json.Unmarshal(data, &resp)
	if len(resp.Templates) != 2 {
		t.Errorf("template count: got %d, want 2", len(resp.Templates))
	}
}
