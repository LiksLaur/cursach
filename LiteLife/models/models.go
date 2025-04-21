package models

import (
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn    *websocket.Conn
	Send    chan []byte
	Cleanup func()
}

func (c *Client) ReadPump() {
	defer func() {
		if c.Cleanup != nil {
			c.Cleanup()
		}
		c.Conn.Close()
	}()

	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (c *Client) WritePump() {
	defer c.Conn.Close()

	for {
		message, ok := <-c.Send
		if !ok {
			c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		err := c.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			return
		}
	}
}

type User struct {
	ID       int
	Username string
	Password string
}

type PageData struct {
	Error    string
	Username string
	Requests []RepairRequest
	Messages []ChatMessage
	Bookings []RoomBooking
}

type RepairRequest struct {
	ID         int
	Name       string
	Username   string
	Apartment  string
	RepairType string
	Comment    string
	CreatedAt  time.Time
	IsApproved bool
}

type ChatMessage struct {
	ID        int
	Username  string
	Message   string
	CreatedAt time.Time
}

type RoomBooking struct {
	ID          int
	Name        string
	Phone       string
	RoomNumber  int
	BookingDate time.Time
	IsApproved  bool
	CreatedAt   time.Time
}
