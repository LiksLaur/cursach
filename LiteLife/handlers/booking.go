package handlers

import (
	"litelife/database"
	"litelife/models"
	"net/http"
)

func RoomBookingHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	username, ok := session.Values["username"].(string)
	if !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	data := models.PageData{Username: username}
	templates.ExecuteTemplate(w, "roomBooking.html", data)
}

func SubmitBookingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/room-booking", http.StatusSeeOther)
		return
	}

	session, _ := store.Get(r, "session-name")
	username, ok := session.Values["username"].(string)
	if !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	name := r.FormValue("name")
	phone := r.FormValue("phone")
	roomNumber := r.FormValue("room")
	bookingDate := r.FormValue("booking_date")

	var count int
	err := database.DB.QueryRow(`
		SELECT COUNT(*) FROM room_bookings
		WHERE room_number = $1 AND booking_date = $2 AND is_approved = true
	`, roomNumber, bookingDate).Scan(&count)

	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	if count > 0 {
		data := models.PageData{
			Username: username,
			Error:    "Room is already booked for this date",
		}
		templates.ExecuteTemplate(w, "roomBooking.html", data)
		return
	}

	_, err = database.DB.Exec(`
		INSERT INTO room_bookings (name, phone, room_number, booking_date)
		VALUES ($1, $2, $3, $4)
	`, name, phone, roomNumber, bookingDate)

	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/user", http.StatusSeeOther)
}

func AdminRoomBookingsHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	username, ok := session.Values["username"].(string)
	if !ok || username != "admin" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	rows, err := database.DB.Query(`
		SELECT id, name, phone, room_number, booking_date, is_approved, created_at
		FROM room_bookings
		ORDER BY created_at DESC
	`)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var bookings []models.RoomBooking
	for rows.Next() {
		var booking models.RoomBooking
		err := rows.Scan(&booking.ID, &booking.Name, &booking.Phone,
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
		Bookings: bookings,
	}
	templates.ExecuteTemplate(w, "adminRoomBooking.html", data)
}

func ApproveBookingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/admin/bookings", http.StatusSeeOther)
		return
	}

	session, _ := store.Get(r, "session-name")
	username, ok := session.Values["username"].(string)
	if !ok || username != "admin" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	bookingID := r.FormValue("booking_id")
	_, err := database.DB.Exec("UPDATE room_bookings SET is_approved = true WHERE id = $1", bookingID)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/bookings", http.StatusSeeOther)
}

func RejectBookingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/admin/bookings", http.StatusSeeOther)
		return
	}

	session, _ := store.Get(r, "session-name")
	username, ok := session.Values["username"].(string)
	if !ok || username != "admin" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	bookingID := r.FormValue("booking_id")
	_, err := database.DB.Exec("DELETE FROM room_bookings WHERE id = $1", bookingID)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/bookings", http.StatusSeeOther)
}
