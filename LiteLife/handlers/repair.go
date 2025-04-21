package handlers

import (
	"net/http"
	"litelife/models"
	"litelife/database"
)

func BuildRequestHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	username, ok := session.Values["username"].(string)
	if !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	data := models.PageData{Username: username}
	templates.ExecuteTemplate(w, "buildRequest.html", data)
}

func SubmitRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/build-request", http.StatusSeeOther)
		return
	}

	session, _ := store.Get(r, "session-name")
	username, ok := session.Values["username"].(string)
	if !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	name := r.FormValue("name")
	apartment := r.FormValue("apartment")
	repairType := r.FormValue("repairType")
	comment := r.FormValue("comment")

	_, err := database.DB.Exec(`
		INSERT INTO repair_requests (name, username, apartment, repair_type, comment)
		VALUES ($1, $2, $3, $4, $5)
	`, name, username, apartment, repairType, comment)

	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/user", http.StatusSeeOther)
}

func AdminRequestsHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	username, ok := session.Values["username"].(string)
	if !ok || username != "admin" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	rows, err := database.DB.Query(`
		SELECT id, name, username, apartment, repair_type, comment, created_at, is_approved
		FROM repair_requests
		ORDER BY created_at DESC
	`)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var requests []models.RepairRequest
	for rows.Next() {
		var req models.RepairRequest
		err := rows.Scan(&req.ID, &req.Name, &req.Username, &req.Apartment,
			&req.RepairType, &req.Comment, &req.CreatedAt, &req.IsApproved)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		requests = append(requests, req)
	}

	data := models.PageData{
		Username: username,
		Requests: requests,
	}
	templates.ExecuteTemplate(w, "adminBuildRequests.html", data)
}

func ApproveRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/admin/requests", http.StatusSeeOther)
		return
	}

	session, _ := store.Get(r, "session-name")
	username, ok := session.Values["username"].(string)
	if !ok || username != "admin" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	requestID := r.FormValue("request_id")
	_, err := database.DB.Exec("UPDATE repair_requests SET is_approved = true WHERE id = $1", requestID)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/requests", http.StatusSeeOther)
} 