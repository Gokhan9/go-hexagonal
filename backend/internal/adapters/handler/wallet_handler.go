package handler

import (
	"encoding/json"
	"errors"
	"go-hexagonal/internal/adapters/handler/middleware"
	"go-hexagonal/internal/api/dto"
	"go-hexagonal/internal/core/domain"
	"go-hexagonal/internal/core/ports"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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

// @Summary Cüzdan Oluştur
// @Tags Wallets
// @Accept json
// @Produce json
// @Param body body dto.CreateWalletRequest true "Cüzdan bilgileri"
// @Success 201 {object} map[string]string
// @Router /wallets [post]
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

// @Summary Cüzdan Detayını Getir
// @Tags Wallets
// @Produce json
// @Param id path string true "Cüzdan ID"
// @Success 200 {object} domain.Wallet
// @Failure 404 {object} map[string]string
// @Router /wallets{id} [get]
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

// @Summary Para yatır
// @Tags Wallets
// @Accept json
// @Produce json
// @Param id path string true "Cüzdan ID"
// @Param X-Idempotency-Key header string false "Mükerrer işlem koruması"
// @Param body body dto.TransactionRequest true "Yatırma detayları"
// @Success 200 {object} map[string]string
// @Router /wallets{id}/deposit [post]
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
	//user, err := middleware.GetUsernameFromContext(r.Context())
	//if err != nil {
	//	h.WriteError(w, http.StatusUnauthorized, err.Error())
	//	return
	// }

	userID := "Gökhan"
	transactionID := uuid.NewString()

	log.Println("Generated transaction:", transactionID)

	err := h.service.Deposit(r.Context(), idempotencyKey, id, userID, transactionID, req.ToCents())
	if err != nil {
		h.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.WriteJSON(w, http.StatusOK, map[string]string{"message": "Para yatırma işlemi başarılı."})
}

// @Summary Para çek
// @Tags Wallets
// @Accept json
// @Produce json
// @Param id path string true "Cüzdan ID"
// @Param X-Idempotency-Key header string false "Mükerrer işlem koruması"
// @Param body body dto.TransactionRequest true "Çekim detayları"
// @Success 200 {object} map[string]string
// @Router /wallets/{id}/withdraw [post]
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
	//user, err := middleware.GetUsernameFromContext(r.Context())
	//if err != nil {
	//	h.WriteError(w, http.StatusUnauthorized, err.Error())
	//	return
	//}

	userID := "Gökhan"
	transactionID := uuid.NewString()
	log.Println("Generated transaction:", transactionID)
	err := h.service.Withdraw(r.Context(), idempotencyKey, id, userID, transactionID, req.ToCents())
	if err != nil {
		h.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.WriteJSON(w, http.StatusOK, map[string]string{"message": "Para çekme işlemi başarılı."})
}

// @Summary İşlem geçmişini getir
// @Tags Wallets
// @Produce json
// @Param id path string true "Cüzdan ID"
// @Success 200 {array} domain.Transaction
// @Router /wallets/{id}/transactions [get]
func (h *WalletHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")

	tns, err := h.service.GetTransactions(r.Context(), id)
	if err != nil {
		h.WriteError(w, http.StatusInternalServerError, "Geçmiş kayıtlar alınamadı..")
		return
	}

	h.WriteJSON(w, http.StatusOK, tns)
}

// @Summary Transfer yap
// @Tags Wallets
// @Accept json
// @Produce json
// @Param id path string true "Gönderen Cüzdan ID"
// @Param X-Idempotency-Key header string false "Mükerrer işlem koruması"
// @Param body body dto.TransferRequest true "Transfer detayları"
// @Success 202 {object} map[string]string
// @Router /wallets/{id}/transfer [post]
func (h *WalletHandler) Transfer(w http.ResponseWriter, r *http.Request) {

	var req dto.TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.WriteError(w, http.StatusBadRequest, "Invalid request body..")
		return
	}

	// validasyon kontrolü
	if err := validate.Struct(req); err != nil {
		h.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	fromID := r.PathValue("id")

	err := h.service.Transfer(r.Context(), r.Header.Get("X-Idempotency-Key"), fromID, req.ToWalletID, req.OwnerID, req.Amount)
	if err != nil {
		h.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"message": "Transfer Successful"})
}

// @Summary Bakiye sorgula
// @Tags Wallets
// @Produce json
// @Param id path string true "Cüzdan ID"
// @Success 200 {object} map[string]int64
// @Router /wallets/{id}/balance [get]
func (h *WalletHandler) GetBalance(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")

	balance, err := h.service.GetBalance(r.Context(), id)
	if err != nil {
		// hataya göre ugun http status
		if errors.Is(err, domain.ErrorWalletNotFound) {
			h.WriteError(w, http.StatusNotFound, err.Error())
			return
		}

		h.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int64{"balance": balance})
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
