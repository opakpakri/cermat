package models

import (
	"encoding/json"
	"time"
)

type User struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Role         string    `json:"role"`
	BusinessName string    `json:"business_name"`
	BusinessMode string    `json:"business_mode"`
	Phone        *string   `json:"phone"`
	Address      *string   `json:"address"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type AdminUserRequest struct {
	Name         string  `json:"name"`
	Email        string  `json:"email"`
	Password     string  `json:"password"`
	Role         string  `json:"role"`
	BusinessName string  `json:"business_name"`
	BusinessMode string  `json:"business_mode"`
	Phone        *string `json:"phone"`
	Address      *string `json:"address"`
}

type Contact struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Name        string    `json:"name"`
	Phone       string    `json:"phone"`
	Email       *string   `json:"email"`
	CompanyName *string   `json:"company_name"`
	Address     *string   `json:"address"`
	LeadSource  string    `json:"lead_source"`
	Status      string    `json:"status"`
	Notes       *string   `json:"notes"`
	CreatedAt   time.Time `json:"created_at"`
}

type CatalogItem struct {
	ID        int             `json:"id"`
	UserID    int             `json:"user_id"`
	Code      *string         `json:"code"`
	Name      string          `json:"name"`
	Category  *string         `json:"category"`
	Price     float64         `json:"price"`
	Status    string          `json:"status"`
	Specs     json.RawMessage `json:"specs"`
	CreatedAt time.Time       `json:"created_at"`
}

type Deal struct {
	ID                int        `json:"id"`
	UserID            int        `json:"user_id"`
	DealCode          string     `json:"deal_code"`
	Title             string     `json:"title"`
	ContactID         int        `json:"contact_id"`
	ContactName       string     `json:"contact_name,omitempty"`
	ContactPhone      string     `json:"contact_phone,omitempty"`
	LeadSource        string     `json:"lead_source,omitempty"`
	CatalogItemID     *int       `json:"catalog_item_id"`
	CatalogItemName   *string    `json:"catalog_item_name,omitempty"`
	CatalogItemCode   *string    `json:"catalog_item_code,omitempty"`
	Stage             string     `json:"stage"`
	DealValue         float64    `json:"deal_value"`
	PaymentPlan       *string    `json:"payment_plan"`
	ExpectedCloseDate *time.Time `json:"expected_close_date"`
	Notes             *string    `json:"notes"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

type DealActivity struct {
	ID           int        `json:"id"`
	UserID       int        `json:"user_id"`
	DealID       *int       `json:"deal_id"`
	DealTitle    *string    `json:"deal_title,omitempty"`
	ContactID    int        `json:"contact_id"`
	ContactName  string     `json:"contact_name,omitempty"`
	ContactPhone string     `json:"contact_phone,omitempty"`
	ActivityType string     `json:"activity_type"`
	Title        string     `json:"title"`
	Description  *string    `json:"description"`
	ScheduledAt  *time.Time `json:"scheduled_at"`
	IsCompleted  bool       `json:"is_completed"`
	CreatedAt    time.Time  `json:"created_at"`
}

type Ticket struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	TicketCode  string    `json:"ticket_code"`
	ContactID   int       `json:"contact_id"`
	ContactName string    `json:"contact_name,omitempty"`
	DealID      *int      `json:"deal_id"`
	DealTitle   *string   `json:"deal_title,omitempty"`
	Subject     string    `json:"subject"`
	Description *string   `json:"description"`
	Priority    string    `json:"priority"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type DashboardSummary struct {
	NewLeadsThisMonth int            `json:"new_leads_this_month"`
	PipelineValue     float64        `json:"pipeline_value"`
	WonDealsCount     int            `json:"won_deals_count"`
	WonDealsValue     float64        `json:"won_deals_value"`
	FollowupsToday    int            `json:"followups_today"`
	RecentActivities  []DealActivity `json:"recent_activities"`
	RecentDeals       []Deal         `json:"recent_deals"`
}
