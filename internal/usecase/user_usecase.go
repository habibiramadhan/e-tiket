//internal/usecase/user_usecase.go

package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"
	
	"ticket-system/internal/domain/entity"
	"ticket-system/internal/domain/repository"
	"ticket-system/pkg/utils"
)

type RegisterRequest struct {
	Username       string `json:"username"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	RetypePassword string `json:"retype_password"`
	Role           string `json:"role"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	Role         string `json:"role"`
	Username     string `json:"username"`
}

type UserUsecase interface {
	Register(ctx context.Context, req RegisterRequest) (int, error)
	Login(ctx context.Context, req LoginRequest) (*LoginResponse, error)
	GetByID(ctx context.Context, id int) (*entity.User, error)
	VerifyEmail(ctx context.Context, token string) error
	ResendVerificationEmail(ctx context.Context, email string) error
	UpdateProfile(ctx context.Context, userID int, name, gender, address, phoneNumber string) error
}

type userUsecase struct {
	userRepo              repository.UserRepository
	userProfileRepo       repository.UserProfileRepository
	emailVerificationRepo repository.EmailVerificationRepository
	jwtSecret             string
	tokenExpiry           int
	smtpConfig            utils.SMTPConfig
	appURL                string
}

func (u *userUsecase) UpdateProfile(ctx context.Context, userID int, name, gender, address, phoneNumber string) error {
	existingProfile, err := u.userProfileRepo.FindByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if existingProfile == nil {
		return errors.New("profil pengguna tidak ditemukan")
	}

	profile := &entity.UserProfile{
		ID:          existingProfile.ID,
		UserID:      userID,
		Name:        name,
		Gender:      gender,
		Address:     address,
		PhoneNumber: phoneNumber,
		UpdatedAt:   time.Now(),
	}

	return u.userProfileRepo.Update(ctx, profile)
}

func NewUserUsecase(
	userRepo repository.UserRepository,
	userProfileRepo repository.UserProfileRepository,
	emailVerificationRepo repository.EmailVerificationRepository,
	jwtSecret string,
	tokenExpiry string,
	smtpConfig utils.SMTPConfig,
	appURL string,
) UserUsecase {
	expiry, _ := strconv.Atoi(tokenExpiry)
	if expiry == 0 {
		expiry = 24 // default 24 jam
	}
	
	return &userUsecase{
		userRepo:              userRepo,
		userProfileRepo:       userProfileRepo,
		emailVerificationRepo: emailVerificationRepo,
		jwtSecret:             jwtSecret,
		tokenExpiry:           expiry,
		smtpConfig:            smtpConfig,
		appURL:                appURL,
	}
}

func (u *userUsecase) Register(ctx context.Context, req RegisterRequest) (int, error) {
	if err := utils.ValidateUsername(req.Username); err != nil {
		return 0, err
	}
	
	if err := utils.ValidateEmail(req.Email); err != nil {
		return 0, err
	}
	
	if err := utils.ValidatePassword(req.Password); err != nil {
		return 0, err
	}

	if req.Password != req.RetypePassword {
		return 0, errors.New("password dan retype password tidak cocok")
	}
	
	existingUser, err := u.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		return 0, err
	}
	if existingUser != nil {
		return 0, errors.New("username sudah digunakan")
	}
	
	existingEmail, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return 0, err
	}
	if existingEmail != nil {
		return 0, errors.New("email sudah digunakan")
	}
	
	if req.Role != "user" && req.Role != "organizer" {
		req.Role = "user" // Default role
	}
	
	hashedPassword, err := utils.GeneratePassword(req.Password)
	if err != nil {
		return 0, err
	}
	
	user := &entity.User{
		Username:   req.Username,
		Email:      req.Email,
		Password:   hashedPassword,
		Role:       req.Role,
		IsVerified: false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	
	userID, err := u.userRepo.Create(ctx, user)
	if err != nil {
		return 0, err
	}
	
	token := utils.GenerateRandomString(64)
	expiredAt := time.Now().Add(24 * time.Hour)
	
	verification := &entity.EmailVerification{
		UserID:    userID,
		Token:     token,
		ExpiredAt: expiredAt,
		CreatedAt: time.Now(),
	}
	
	_, err = u.emailVerificationRepo.Create(ctx, verification)
	if err != nil {
		return 0, err
	}
	
	go u.sendVerificationEmail(user.Username, user.Email, token)
	
	return userID, nil
}

func (u *userUsecase) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	var user *entity.User
	var err error
	
	if utils.ValidateEmail(req.Username) == nil {
		user, err = u.userRepo.FindByEmail(ctx, req.Username)
	} else {
		user, err = u.userRepo.FindByUsername(ctx, req.Username)
	}
	
	if err != nil {
		return nil, err
	}
	
	if user == nil {
		return nil, errors.New("username atau password salah")
	}
	
	match, err := utils.VerifyPassword(req.Password, user.Password)
	if err != nil {
		return nil, err
	}
	
	if !match {
		return nil, errors.New("username atau password salah")
	}
	
	if !user.IsVerified {
		return nil, errors.New("email belum diverifikasi, silakan periksa email Anda")
	}
	
	token, err := utils.GenerateJWT(user.ID, user.Username, user.Email, user.Role, u.jwtSecret, u.tokenExpiry)
	if err != nil {
		return nil, err
	}
	
	refreshToken, err := utils.GenerateJWT(user.ID, user.Username, user.Email, user.Role, u.jwtSecret, u.tokenExpiry*2)
	if err != nil {
		return nil, err
	}
	
	response := &LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		Role:         user.Role,
		Username:     user.Username,
	}
	
	return response, nil
}

func (u *userUsecase) GetByID(ctx context.Context, id int) (*entity.User, error) {
	return u.userRepo.FindByID(ctx, id)
}

func (u *userUsecase) CreateProfile(ctx context.Context, userID int, name, gender, address, phoneNumber string) (int, error) {
	existingProfile, err := u.userProfileRepo.FindByUserID(ctx, userID)
	if err != nil {
		return 0, err
	}
	
	if existingProfile != nil {
		return 0, errors.New("user sudah memiliki profil")
	}
	
	profile := &entity.UserProfile{
		UserID:      userID,
		Name:        name,
		Gender:      gender,
		Address:     address,
		PhoneNumber: phoneNumber,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	
	profileID, err := u.userProfileRepo.Create(ctx, profile)
	if err != nil {
		return 0, err
	}
	
	return profileID, nil
}

func (u *userUsecase) VerifyEmail(ctx context.Context, token string) error {
	verification, err := u.emailVerificationRepo.FindByToken(ctx, token)
	if err != nil {
		return err
	}
	
	if verification == nil {
		return errors.New("token verifikasi tidak valid")
	}
	
	if time.Now().After(verification.ExpiredAt) {
		return errors.New("token verifikasi sudah kedaluwarsa")
	}
	
	if err := u.userRepo.UpdateVerificationStatus(ctx, verification.UserID, true); err != nil {
		return err
	}
	
	if err := u.userRepo.CreateDefaultProfile(ctx, verification.UserID); err != nil {
		return err
	}
	
	return u.emailVerificationRepo.Delete(ctx, verification.ID)
}

func (u *userUsecase) ResendVerificationEmail(ctx context.Context, email string) error {
	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return err
	}
	
	if user == nil {
		return errors.New("email tidak terdaftar")
	}
	
	if user.IsVerified {
		return errors.New("email sudah diverifikasi")
	}
	
	err = u.emailVerificationRepo.DeleteByUserID(ctx, user.ID)
	if err != nil {
		return err
	}
	
	token := utils.GenerateRandomString(64)
	expiredAt := time.Now().Add(24 * time.Hour)
	
	verification := &entity.EmailVerification{
		UserID:    user.ID,
		Token:     token,
		ExpiredAt: expiredAt,
		CreatedAt: time.Now(),
	}
	
	_, err = u.emailVerificationRepo.Create(ctx, verification)
	if err != nil {
		return err
	}
	
	go u.sendVerificationEmail(user.Username, user.Email, token)
	
	return nil
}

func (u *userUsecase) sendVerificationEmail(username, email, token string) {
	verificationLink := fmt.Sprintf("%s/api/verify-email?token=%s", u.appURL, token)
	templateData := map[string]interface{}{
		"Username":         username,
		"VerificationLink": verificationLink,
		"Year":             time.Now().Year(),
	}
	
	body, err := utils.ParseTemplate("templates/email/verification.html", templateData)
	if err != nil {
		log.Printf("Gagal parse template email: %v", err)
		return
	}
	
	emailData := utils.EmailData{
		To:      []string{email},
		Subject: "Verifikasi Email - Sistem Tiket Event",
		Body:    body,
	}
	
	if err := utils.SendEmail(u.smtpConfig, emailData); err != nil {
		log.Printf("Gagal mengirim email verifikasi: %v", err)
	} else {
		log.Printf("Email verifikasi berhasil dikirim ke: %s", email)
	}
}