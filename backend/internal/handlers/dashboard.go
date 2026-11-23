package handlers

import (
	"cermat-api/internal/middleware"
	"cermat-api/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DashboardHandler struct {
	DB *pgxpool.Pool
}

func (h *DashboardHandler) GetSummary(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	ctx := c.Context()

	var summary models.DashboardSummary

	// 1. Prospek baru bulan ini
	err := h.DB.QueryRow(ctx, `
		SELECT COUNT(*) FROM contacts
		WHERE user_id = $1
		  AND date_trunc('month', created_at) = date_trunc('month', CURRENT_TIMESTAMP)
	`, userID).Scan(&summary.NewLeadsThisMonth)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memuat ringkasan kontak"})
	}

	// 2. Estimasi nilai pipeline aktif (selain WON dan LOST)
	err = h.DB.QueryRow(ctx, `
		SELECT COALESCE(SUM(deal_value), 0) FROM deals
		WHERE user_id = $1 AND stage NOT IN ('WON', 'LOST')
	`, userID).Scan(&summary.PipelineValue)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memuat nilai pipeline"})
	}

	// 3. Total Closing / Deal bulan ini (WON)
	err = h.DB.QueryRow(ctx, `
		SELECT COUNT(*), COALESCE(SUM(deal_value), 0) FROM deals
		WHERE user_id = $1 AND stage = 'WON'
		  AND date_trunc('month', updated_at) = date_trunc('month', CURRENT_TIMESTAMP)
	`, userID).Scan(&summary.WonDealsCount, &summary.WonDealsValue)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memuat data closing"})
	}

	// 4. Tugas follow-up hari ini (atau yang belum selesai sampai hari ini)
	err = h.DB.QueryRow(ctx, `
		SELECT COUNT(*) FROM deal_activities
		WHERE user_id = $1 AND is_completed = FALSE
		  AND scheduled_at IS NOT NULL
		  AND scheduled_at::date <= CURRENT_DATE
	`, userID).Scan(&summary.FollowupsToday)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memuat jadwal follow-up"})
	}

	// 5. Aktivitas & Pengingat Terbaru
	summary.RecentActivities = make([]models.DealActivity, 0)
	actRows, err := h.DB.Query(ctx, `
		SELECT a.id, a.user_id, a.deal_id, d.title, a.contact_id, c.name, c.phone,
		       a.activity_type, a.title, a.description, a.scheduled_at, a.is_completed, a.created_at
		FROM deal_activities a
		JOIN contacts c ON c.id = a.contact_id AND c.user_id = a.user_id
		LEFT JOIN deals d ON d.id = a.deal_id AND d.user_id = a.user_id
		WHERE a.user_id = $1
		ORDER BY a.is_completed ASC, COALESCE(a.scheduled_at, a.created_at) ASC
		LIMIT 6
	`, userID)
	if err == nil {
		defer actRows.Close()
		for actRows.Next() {
			var act models.DealActivity
			if scanErr := actRows.Scan(
				&act.ID, &act.UserID, &act.DealID, &act.DealTitle, &act.ContactID, &act.ContactName, &act.ContactPhone,
				&act.ActivityType, &act.Title, &act.Description, &act.ScheduledAt, &act.IsCompleted, &act.CreatedAt,
			); scanErr == nil {
				summary.RecentActivities = append(summary.RecentActivities, act)
			}
		}
	}

	// 6. Progres transaksi terkini
	summary.RecentDeals = make([]models.Deal, 0)
	dealRows, err := h.DB.Query(ctx, `
		SELECT d.id, d.user_id, d.deal_code, d.title, d.contact_id, c.name, c.phone, c.lead_source,
		       d.catalog_item_id, ci.name, ci.code, d.stage, d.deal_value, d.payment_plan,
		       d.expected_close_date, d.notes, d.created_at, d.updated_at
		FROM deals d
		JOIN contacts c ON c.id = d.contact_id AND c.user_id = d.user_id
		LEFT JOIN catalog_items ci ON ci.id = d.catalog_item_id AND ci.user_id = d.user_id
		WHERE d.user_id = $1
		ORDER BY d.updated_at DESC
		LIMIT 5
	`, userID)
	if err == nil {
		defer dealRows.Close()
		for dealRows.Next() {
			var dl models.Deal
			if scanErr := dealRows.Scan(
				&dl.ID, &dl.UserID, &dl.DealCode, &dl.Title, &dl.ContactID, &dl.ContactName, &dl.ContactPhone, &dl.LeadSource,
				&dl.CatalogItemID, &dl.CatalogItemName, &dl.CatalogItemCode, &dl.Stage, &dl.DealValue, &dl.PaymentPlan,
				&dl.ExpectedCloseDate, &dl.Notes, &dl.CreatedAt, &dl.UpdatedAt,
			); scanErr == nil {
				summary.RecentDeals = append(summary.RecentDeals, dl)
			}
		}
	}

	return c.JSON(summary)
}
