package service

import (
	"net/http"

	"flow-desk/dto"
	"flow-desk/entity"
	"flow-desk/repository"
	"flow-desk/utils"
)

type AuthService interface {
	Register(req dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(req dto.LoginRequest) (*dto.AuthResponse, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Register(req dto.RegisterRequest) (*dto.AuthResponse, error) {
	existingUser, _ := s.userRepo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, utils.NewBadRequestError("EMAIL_ALREADY_EXISTS", "Email sudah terdaftar", nil)
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, utils.NewInternalServerError("HASH_PASSWORD_ERROR", "Gagal memproses password", err)
	}

	user := entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
	}

	err = s.userRepo.Create(&user)
	if err != nil {
		return nil, utils.NewInternalServerError("CREATE_USER_ERROR", "Gagal menyimpan user", err)
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return nil, utils.NewInternalServerError("GENERATE_TOKEN_ERROR", "Gagal membuat token autentikasi", err)
	}

	return &dto.AuthResponse{
		Token: token,
		User: dto.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	}, nil
}

func (s *authService) Login(req dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, utils.NewAppError(http.StatusUnauthorized, "AUTH_INVALID_CREDENTIALS", "Email atau password salah", err)
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, utils.NewAppError(http.StatusUnauthorized, "AUTH_INVALID_CREDENTIALS", "Email atau password salah", nil)
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return nil, utils.NewInternalServerError("GENERATE_TOKEN_ERROR", "Gagal membuat token autentikasi", err)
	}

	return &dto.AuthResponse{
		Token: token,
		User: dto.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	}, nil
}
