package service

import (
	"errors"
	"jastip-jakarta/features/admin"
	ud "jastip-jakarta/features/user"
	"jastip-jakarta/utils/encrypts"
	"jastip-jakarta/utils/middlewares"
	"mime/multipart"
)

type adminService struct {
	adminData   admin.AdminDataInterface
	hashService encrypts.HashInterface
	userData    ud.UserDataInterface
}

// dependency injection
func New(repo admin.AdminDataInterface, hash encrypts.HashInterface, userData ud.UserDataInterface) admin.AdminServiceInterface {
	return &adminService{
		adminData:   repo,
		hashService: hash,
		userData:    userData,
	}
}

// CreateSuper implements admin.AdminServiceInterface.
func (u *adminService) CreateSuper(input admin.Admin) error {
	if input.Password != "" {
		hashedPass, errHash := u.hashService.HashPassword(input.Password)
		if errHash != nil {
			return errors.New("error hash password")
		}
		input.Password = hashedPass
	}

	if input.Role == "" {
		input.Role = "Super"
	}

	err := u.adminData.Insert(input)
	return err
}

// Create implements admin.AdminServiceInterface.
func (u *adminService) Create(adminIdLogin int, input admin.Admin) error {
	role, err := u.adminData.SelectById(adminIdLogin)
	if err != nil {
		return err
	}

	if role.Role != "Super" {
		return errors.New("anda tidak memiliki akses untuk menggunakan fitur ini")
	}

	if input.Name == "" {
		return errors.New("nama tidak boleh kosong")
	}
	if input.Email == "" {
		return errors.New("email tidak boleh kosong")
	}
	if input.Password == "" {
		return errors.New("password tidak boleh kosong")
	}
	if input.PhoneNumber == 0 {
		return errors.New("nomor Telephone tidak boleh kosong")
	}
	if input.Role == "" {
		return errors.New("role tidak boleh kosong")
	}

	if input.Password != "" {
		hashedPass, errHash := u.hashService.HashPassword(input.Password)
		if errHash != nil {
			return errors.New("error hash password")
		}
		input.Password = hashedPass
	}

	err = u.adminData.Insert(input)
	return err
}

// GetById implements admin.AdminServiceInterface.
func (u *adminService) GetById(adminIdLogin int) (*admin.Admin, error) {
	adminData, err := u.adminData.SelectById(adminIdLogin)
	if err != nil {
		return nil, err
	}
	return adminData, nil
}

// Login implements admin.AdminServiceInterface.
func (u *adminService) Login(phoneOrEmail string, password string) (data *admin.Admin, token string, err error) {
	// Validasi jika email atau password kosong
	if phoneOrEmail == "" {
		return nil, "", errors.New("email atau nomor telepon tidak boleh kosong")
	}
	if password == "" {
		return nil, "", errors.New("password tidak boleh kosong")
	}

	data, err = u.adminData.Login(phoneOrEmail, password)
	if err != nil {
		return nil, "", err
	}

	isValid := u.hashService.CheckPasswordHash(data.Password, password)
	if !isValid {
		return nil, "", errors.New("sandi salah")
	}

	token, errJwt := middlewares.CreateToken(int(data.ID))
	if errJwt != nil {
		return nil, "", errJwt
	}

	return data, token, err
}

// Update implements admin.AdminServiceInterface.
func (u *adminService) Update(adminIdLogin int, photo *multipart.FileHeader) error {
	if photo == nil {
		return errors.New("tidak ada foto yang di upload")
	}
	err := u.adminData.Update(adminIdLogin, photo)
	return err
}

// CreateRegionCode implements admin.AdminServiceInterface.
func (u *adminService) CreateRegionCode(adminIdLogin int, input admin.RegionCode) error {
	role, err := u.adminData.SelectById(adminIdLogin)
	if err != nil {
		return err
	}

	if role.Role != "Super" {
		return errors.New("anda tidak memiliki akses untuk menggunakan fitur ini")
	}

	err = u.adminData.InsertRegionCode(input)
	return err
}

// GetAllRegionCode implements admin.AdminServiceInterface.
func (u *adminService) GetAllRegionCode() ([]admin.RegionCode, error) {
	return u.adminData.SelectAllRegionCode()
}

// GettByIdRegion implements admin.AdminServiceInterface.
func (u *adminService) GettByIdRegion(IdRegion string) (*admin.RegionCode, error) {
	regionCode, err := u.adminData.SelectByIdRegion(IdRegion)
	if err != nil {
		return nil, err
	}
	return regionCode, nil
}

// CreateBatchDelivery implements admin.AdminServiceInterface.
func (u *adminService) CreateBatchDelivery(adminIdLogin int, input admin.DeliveryBatch) error {
	existingBatch, _ := u.adminData.SelectDeliveryBatch(input.ID)
	if existingBatch != nil {
		return errors.New("batch sudah ada")
	}

	err := u.adminData.InsertBatchDelivery(adminIdLogin, input)
	if err != nil {
		return err
	}
	return nil
}

// GetAllBatchDelivery implements admin.AdminServiceInterface.
func (u *adminService) GetAllBatchDelivery() ([]admin.DeliveryBatch, error) {
	responseBatch, err := u.adminData.SelectAllBatchDelivery()
	if err != nil {
		return nil, err
	}
	return responseBatch, nil
}

// GetDeliveryBatch implements admin.AdminServiceInterface.
func (u *adminService) GetDeliveryBatch(batchID string) (*admin.DeliveryBatch, error) {
	deliveryBatch, err := u.adminData.SelectDeliveryBatch(batchID)
	if err != nil {
		return nil, err
	}
	return deliveryBatch, nil
}

// GettAdminsByRole implements admin.AdminServiceInterface.
func (u *adminService) GetAdminsByRole(adminIdLogin int, role string) ([]admin.Admin, error) {
	adminCheck, err := u.GetById(adminIdLogin)
	if err != nil || adminCheck.Role != "Super" {
		return nil, errors.New("anda bukan admin super")
	}

	adminRes, err := u.adminData.SelectAdminsByRole(role)
	if err != nil {
		return nil, err
	}
	return adminRes, nil
}

// GettAllAdmins implements admin.AdminServiceInterface.
func (u *adminService) GetAllAdmins(adminIdLogin int) ([]admin.Admin, error) {
	adminCheck, err := u.GetById(adminIdLogin)
	if err != nil || adminCheck.Role != "Super" {
		return nil, errors.New("anda bukan admin super")
	}

	adminRes, err := u.adminData.SelectAllAdmins()
	if err != nil {
		return nil, err
	}
	return adminRes, nil
}

// SearchRegionCode implements admin.AdminServiceInterface.
func (u *adminService) SearchRegionCode(adminIdLogin int, code string) ([]admin.RegionCode, error) {
	adminCheck, err := u.GetById(adminIdLogin)
	if err != nil || adminCheck.Role != "Super" {
		return nil, errors.New("anda bukan admin super")
	}

	regionCodes, err := u.adminData.SearchRegionCode(code)
	if err != nil {
		return nil, err
	}
	return regionCodes, nil
}

// UpdateRegionCode implements admin.AdminServiceInterface.
func (u *adminService) UpdateRegionCode(adminIdLogin int, code string, updatedRegion admin.RegionCode) error {
	adminCheck, err := u.GetById(adminIdLogin)
	if err != nil || adminCheck.Role != "Super" {
		return errors.New("anda bukan admin super")
	}

	err = u.adminData.UpdateRegionCode(code, updatedRegion)
	if err != nil {
		return err
	}
	return nil
}

// SearchUser implements admin.AdminServiceInterface.
func (u *adminService) SearchUser(adminIdLogin int, query string) ([]ud.User, error) {
	adminCheck, err := u.GetById(adminIdLogin)
	if err != nil || adminCheck.Role != "Super" {
		return nil, errors.New("anda bukan admin super")
	}

	users, err := u.userData.SelectByNameOrEmail(query)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *adminService) UpdateUserByName(adminIdLogin int, name string, input ud.User) error {
	adminCheck, err := u.GetById(adminIdLogin)
	if err != nil || adminCheck.Role != "Super" {
		return errors.New("anda bukan admin super")
	}

	if input.Password != "" {
		hashedPass, errHash := u.hashService.HashPassword(input.Password)
		if errHash != nil {
			return errors.New("error hash password")
		}
		input.Password = hashedPass
	}

	errUpdate := u.userData.UpdateUserByName(name, input)
	if errUpdate != nil {
		return errors.New("error update user: " + errUpdate.Error())
	}

	return nil
}

// CreateUser implements admin.AdminServiceInterface.
func (u *adminService) CreateUser(adminIdLogin int, input ud.User) error {
	adminCheck, err := u.GetById(adminIdLogin)
	if err != nil || adminCheck.Role != "Super" {
		return errors.New("anda bukan admin super")
	}

	if input.Name == "" {
		return errors.New("nama tidak boleh kosong")
	}
	if input.Email == "" {
		return errors.New("email tidak boleh kosong")
	}
	if input.PhoneNumber == 0 {
		return errors.New("nomor Telephone tidak boleh kosong")
	}

	if input.Password != "" {
		hashedPass, errHash := u.hashService.HashPassword(input.Password)
		if errHash != nil {
			return errors.New("error hash password")
		}
		input.Password = hashedPass
	}

	err = u.userData.Insert(input)
	return err
}

// GetAllUser implements admin.AdminServiceInterface.
func (u *adminService) GetAllUser(adminIdLogin int) ([]ud.User, error) {
	adminCheck, err := u.GetById(adminIdLogin)
	if err != nil || adminCheck.Role != "Super" {
		return nil, errors.New("anda bukan admin super")
	}

	userResponse, err := u.userData.SelectAllUser()
	if err != nil {
		return nil, err
	}
	return userResponse, nil
}