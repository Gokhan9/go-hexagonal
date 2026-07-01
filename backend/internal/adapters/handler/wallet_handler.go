package handler

import (
	"encoding/json"
	"go-hexagonal/internal/adapters/handler/middleware"
	"go-hexagonal/internal/api/dto"
	"go-hexagonal/internal/core/ports"
	"net/http"

	"github.com/go-playground/validator/v10"
)

/*
*Şimdi HTTP isteklerini karşılayacak olan ana adaptörümüzü yazacağız.
 */

// dependency inject
type WalletHandler struct {
	service ports.WalletService
}

func NewWalletHandler(walletService ports.WalletService) *WalletHandler {
	return &WalletHandler{
		service: walletService,
	}
}

// Validator
var validate = validator.New()

// 1. POST /wallets
func (h *WalletHandler) Create(w http.ResponseWriter, r *http.Request) {

	var req dto.CreateWalletRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.WriteError(w, http.StatusBadRequest, "Invalid Request Body.") // 400 (Client/Kullanıcı Hatası)
		return
	}

	// Validasyon
	if err := validate.Struct(req); err != nil {
		h.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	wallet, err := h.service.CreateWallet(r.Context(), req.Owner, req.Currency)
	if err != nil {
		h.WriteError(w, http.StatusInternalServerError, err.Error()) // 500 (Sunucu Hatası)
		return
	}

	h.WriteJSON(w, http.StatusCreated, dto.ToDomainResponse(wallet))
}

// 2. GET /wallets/{id}
func (h *WalletHandler) GetByID(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id") // * "Path" parametre okuma.

	user, _ := middleware.GetUsernameFromContext(r.Context())

	if id == "" {
		h.WriteError(w, http.StatusBadRequest, "Wallet ID required.")
		return
	}

	wallet, err := h.service.GetWallet(r.Context(), user.UserID, id)
	if err != nil {
		h.WriteError(w, http.StatusNotFound, "No wallet information was found for this ID.") // 404 (İstenen kayıt/sayfa/api bulunamadığı durumlar.)
		return
	}

	h.WriteJSON(w, http.StatusOK, dto.ToDomainResponse(wallet))
}

// 3. POST /wallets/{id}/deposit
func (h *WalletHandler) Deposit(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")
	idempotencyKey := r.Header.Get("X-Idempotency-Key") //! X-Idempotency-Key

	var req dto.TransactionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.WriteError(w, http.StatusBadRequest, "Geçersiz bakiye bilgisi")
		return
	}

	// ? Tutar Doğrulaması
	if req.Amount <= 0 {
		h.WriteError(w, http.StatusBadRequest, "Yatırılacak tutar 0'dan büyük olmalıdır.")
		return
	}

	// ! idempotencyKey
	user, err := middleware.GetUsernameFromContext(r.Context())
	if err != nil {
		h.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}
	err = h.service.Deposit(r.Context(), idempotencyKey, id, user.UserID, req.TransactionID, req.ToCents())
	if err != nil {
		h.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.WriteJSON(w, http.StatusOK, map[string]string{"message": "Para yatırma işlemi başarılı."})
}

// 4. POST /wallets{id}/withdraw
func (h *WalletHandler) Withdraw(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")
	idempotencyKey := r.Header.Get("X-Idempotency-Key") //! X-Idempotency-Key

	var req dto.TransactionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.WriteError(w, http.StatusBadRequest, "Geçersiz bakiye bilgisi")
		return
	}

	if req.Amount <= 0 {
		h.WriteError(w, http.StatusBadRequest, "Çekeceğiniz tutar 0'dan büyük olmalıdır.")
		return
	}

	// ! idempotencyKey
	user, err := middleware.GetUsernameFromContext(r.Context())
	if err != nil {
		h.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}
	err = h.service.Withdraw(r.Context(), idempotencyKey, id, user.UserID, req.TransactionID, req.ToCents())
	if err != nil {
		h.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.WriteJSON(w, http.StatusOK, map[string]string{"message": "Para çekme işlemi başarılı."})
}

func (h *WalletHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")

	tns, err := h.service.GetTransactions(r.Context(), id)
	if err != nil {
		h.WriteError(w, http.StatusInternalServerError, "Geçmiş kayıtlar alınamadı..")
		return
	}

	h.WriteJSON(w, http.StatusOK, tns)
}

// Yardımcı JSON Metodları.
func (h *WalletHandler) WriteJSON(w http.ResponseWriter, status int, data interface{}) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// Errors Function
func (h *WalletHandler) WriteError(w http.ResponseWriter, status int, message string) {

	h.WriteJSON(w, status, map[string]string{"error": message})
}
