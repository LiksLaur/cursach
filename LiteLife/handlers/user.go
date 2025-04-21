package handlers

import (
	"litelife/database"
	"litelife/models"
	"net/http"
)

func UserIndexHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	username, ok := session.Values["username"].(string)
	if !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	rows, err := database.DB.Query(`
		SELECT id, name, apartment, repair_type, comment, created_at, is_approved
		FROM repair_requests
		WHERE username = $1
		ORDER BY created_at DESC
	`, username)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var requests []models.RepairRequest
	for rows.Next() {
		var req models.RepairRequest
		err := rows.Scan(&req.ID, &req.Name, &req.Apartment, &req.RepairType,
			&req.Comment, &req.CreatedAt, &req.IsApproved)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		requests = append(requests, req)
	}

	messages, err := LoadMessages()
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	data := models.PageData{
		Username: username,
		Requests: requests,
		Messages: messages,
	}
	templates.ExecuteTemplate(w, "userindex.html", data)
}
