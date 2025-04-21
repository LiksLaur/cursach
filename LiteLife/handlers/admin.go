package handlers

import (
	"litelife/database"
	"litelife/models"
	"net/http"
)

func AdminIndexHandler(w http.ResponseWriter, r *http.Request) {
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

	messages, err := LoadMessages()
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	bookingsRows, err := database.DB.Query(`
		SELECT id, name, phone, room_number, booking_date, is_approved, created_at
		FROM room_bookings
		ORDER BY created_at DESC
	`)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	defer bookingsRows.Close()

	var bookings []models.RoomBooking
	for bookingsRows.Next() {
		var booking models.RoomBooking
		err := bookingsRows.Scan(&booking.ID, &booking.Name, &booking.Phone,
			&booking.RoomNumber, &booking.BookingDate, &booking.IsApproved,
			&booking.CreatedAt)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		bookings = append(bookings, booking)
	}

	data := models.PageData{
		Username: username,
		Requests: requests,
		Messages: messages,
		Bookings: bookings,
	}
	templates.ExecuteTemplate(w, "adminindex.html", data)
}
