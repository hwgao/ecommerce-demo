package service

import (
    "errors"
    "time"
    "golang.org/x/crypto/bcrypt"
    "github.com/golang-jwt/jwt/v4"
    "github.com/google/uuid"
    "ecommerce/services/user-service/internal/domain"
    "ecommerce/shared/pkg/cache"
    "ecommerce/shared/pkg/events"
)

type UserServiceImpl struct {
    userRepo    domain.UserRepository
    cache       cache.Cache
    eventBus    events.EventBus
    jwtSecret   string
}

func NewUserService(userRepo domain.UserRepository, cache cache.Cache, eventBus events.EventBus, jwtSecret string) *UserServiceImpl {
    return &UserServiceImpl{
        userRepo:  userRepo,
        cache:     cache,
        eventBus:  eventBus,
        jwtSecret: jwtSecret,
    }
}

func (s *UserServiceImpl) Register(email, password, firstName, lastName string) (*domain.User, error) {
    // Check if user already exists
    existingUser, _ := s.userRepo.GetByEmail(email)
    if existingUser != nil {
        return nil, errors.New("user already exists")
    }

    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }

    user := &domain.User{
        Email:     email,
        Password:  string(hashedPassword),
        FirstName: firstName,
        LastName:  lastName,
    }

    if err := s.userRepo.Create(user); err != nil {
        return nil, err
    }

    // Publish user registered event
    event := events.UserRegisteredEvent{
        UserID:    user.ID,
        Email:     user.Email,
        FirstName: user.FirstName,
        LastName:  user.LastName,
        Timestamp: time.Now(),
    }
    s.eventBus.Publish("user.registered", event)

    return user, nil
}

func (s *UserServiceImpl) Login(email, password string) (string, error) {
    user, err := s.userRepo.GetByEmail(email)
    if err != nil {
        return "", errors.New("invalid credentials")
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        return "", errors.New("invalid credentials")
    }

    // Generate JWT token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID,
        "email":   user.Email,
        "exp":     time.Now().Add(time.Hour * 24).Unix(),
    })

    tokenString, err := token.SignedString([]byte(s.jwtSecret))
    if err != nil {
        return "", err
    }

    // Cache user session
    s.cache.Set(user.ID.String(), user, time.Hour*24)

    return tokenString, nil
}

func (s *UserServiceImpl) GetProfile(userID uuid.UUID) (*domain.User, error) {
    // Try cache first
    if cachedUser, found := s.cache.Get(userID.String()); found {
        if user, ok := cachedUser.(*domain.User); ok {
            return user, nil
        }
    }

    user, err := s.userRepo.GetByID(userID)
    if err != nil {
        return nil, err
    }

    // Cache the user
    s.cache.Set(userID.String(), user, time.Hour*24)
    return user, nil
}

func (s *UserServiceImpl) UpdateProfile(userID uuid.UUID, firstName, lastName string) error {
    user, err := s.userRepo.GetByID(userID)
    if err != nil {
        return err
    }

    user.FirstName = firstName
    user.LastName = lastName

    if err := s.userRepo.Update(user); err != nil {
        return err
    }

    // Update cache
    s.cache.Set(userID.String(), user, time.Hour*24)
    return nil
}
