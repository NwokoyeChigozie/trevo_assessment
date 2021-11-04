package controllers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	validator "github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/gregoflash05/trove/models"
	"github.com/gregoflash05/trove/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var (
	validate               = validator.New()
	errEmailNotValid       = errors.New("email address is not valid")
	errHashingFailed       = errors.New("failed to hashed password")
	ErrUserNotFound        = errors.New("user not found, confirm and try again")
	ErrInvalidCredentials  = errors.New("invalid login credentials, confirm and try again")
	ErrAccountConfirmError = errors.New("your account is not verified, kindly check your email for verification code")
	ErrAccessExpired       = errors.New("error fetching user info, access token expired, kindly login again")
	ErrGeneratingToken     = errors.New("error generating token")
	ErrConfirmPassword     = errors.New("passwords do not match")
	nttrep                 = 201
	DefaultHashCode        = 14
)

// Method to hash password.
func GenerateHashPassword(password string) (string, error) {
	cost := 14
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)

	return string(bytes), err
}

func FetchUserByEmail(filter map[string]interface{}) (*models.User, error) {
	u := &models.User{}
	userCollection, err := utils.GetMongoDBCollection(os.Getenv("DB_NAME"), models.UserCollectionName)

	if err != nil {
		return u, err
	}

	result := userCollection.FindOne(context.TODO(), filter)
	err = result.Decode(&u)

	return u, err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func UserCreate(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")

	var user models.User
	err := utils.ParseJSONFromRequest(request, &user)

	if err != nil {
		utils.GetError(err, nttrep, response)
		return
	}

	userEmail := strings.ToLower(user.Email)
	if !utils.IsValidEmail(userEmail) {
		utils.GetError(errEmailNotValid, nttrep, response)
		return
	}

	// confirm if user_email exists
	result, _ := utils.GetMongoDBDoc(models.UserCollectionName, bson.M{"email": userEmail})
	if result != nil {
		utils.GetError(
			fmt.Errorf("user with email: %s already exists", userEmail),
			nttrep,
			response,
		)

		return
	}

	hashPassword, err := GenerateHashPassword(user.Password)
	if err != nil {
		utils.GetError(errHashingFailed, nttrep, response)
		return
	}

	user.Email = userEmail
	user.CreatedAt = time.Now()
	user.Password = hashPassword
	user.IsVerified = true
	user.PortfolioPositions = models.Positions
	detail, _ := utils.StructToMap(user)

	res, err := utils.CreateMongoDBDoc(models.UserCollectionName, detail)

	if err != nil {
		utils.GetError(err, nttrep, response)
		return
	}

	respse := map[string]interface{}{
		"user_id": res.InsertedID,
	}

	utils.GetSuccess("user created", respse, response)
}

func UserLogin(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")

	var creds models.AuthCredentials
	if err := utils.ParseJSONFromRequest(request, &creds); err != nil {
		utils.GetError(err, nttrep, response)
		return
	}

	if err := validate.Struct(creds); err != nil {
		utils.GetError(err, nttrep, response)
		return
	}

	vser, err := FetchUserByEmail(bson.M{"email": strings.ToLower(creds.Email)})
	if err != nil {
		utils.GetError(ErrUserNotFound, nttrep, response)
		return
	}
	// check if user is verified
	if !vser.IsVerified {
		utils.GetError(ErrAccountConfirmError, nttrep, response)
		return
	}

	// check password
	check := CheckPassword(creds.Password, vser.Password)
	if !check {
		utils.GetError(ErrInvalidCredentials, nttrep, response)
		return
	}

	token, err := CreateToken(vser.ID)
	if err != nil {
		utils.GetError(ErrGeneratingToken, nttrep, response)
		return
	}
	utils.GetSuccess("user created", token, response)
}

func VerifyTokenHandler(response http.ResponseWriter, request *http.Request) {
	_, err := TokenValid(request)
	if err != nil {
		utils.GetError(fmt.Errorf("token Invalid"), nttrep, response)
		return
	}

	utils.GetSuccess("token is valid", "", response)
}

func GetUser(response http.ResponseWriter, request *http.Request) {
	userID, err := TokenValid(request)
	if err != nil {
		utils.GetError(fmt.Errorf("token Invalid"), nttrep, response)
		return
	}

	objID, _ := primitive.ObjectIDFromHex(userID)

	res, err := utils.GetMongoDBDoc(models.UserCollectionName, bson.M{"_id": objID})

	if err != nil {
		utils.GetError(errors.New("user not found"), nttrep, response)
		return
	}

	DeleteMapProps(res, []string{"password"})
	utils.GetSuccess("user retrieved successfully", res, response)
}

func UserUpdate(response http.ResponseWriter, request *http.Request) {
	userID, err := TokenValid(request)
	if err != nil {
		utils.GetError(fmt.Errorf("token Invalid"), nttrep, response)
		return
	}

	var user models.UserUpdate
	if err = utils.ParseJSONFromRequest(request, &user); err != nil {
		utils.GetError(errors.New("bad update data"), nttrep, response)
		return
	}

	userMap, err := utils.StructToMap(user)
	if err != nil {
		utils.GetError(err, nttrep, response)
	}

	updateFields := make(map[string]interface{})

	for key, value := range userMap {
		if value != "" {
			updateFields[key] = value
		}
	}

	if len(updateFields) == 0 {
		utils.GetError(errors.New("empty/invalid user input data"), nttrep, response)

		return
	}

	_, err = utils.UpdateOneMongoDBDoc(models.UserCollectionName, userID, updateFields)

	if err != nil {
		utils.GetError(errors.New("user update failed"), nttrep, response)

		return
	}

	utils.GetSuccess("user successfully updated", nil, response)
}

func PasswordUpdate(response http.ResponseWriter, request *http.Request) {
	userID, err := TokenValid(request)
	if err != nil {
		utils.GetError(fmt.Errorf("token Invalid"), nttrep, response)
		return
	}

	var rBody models.PasswordUpdate
	if e := utils.ParseJSONFromRequest(request, &rBody); e != nil {
		utils.GetError(e, nttrep, response)
		return
	}

	if er := validate.Struct(rBody); er != nil {
		utils.GetError(er, nttrep, response)
		return
	}

	if rBody.Password != rBody.ConfirmPassword {
		utils.GetError(ErrConfirmPassword, nttrep, response)
		return
	}

	// update password & delete passwordreset object
	bytes, err := bcrypt.GenerateFromPassword([]byte(rBody.Password), DefaultHashCode)
	if err != nil {
		utils.GetError(err, nttrep, response)
		return
	}

	id, _ := primitive.ObjectIDFromHex(userID)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"password": string(bytes)}}

	if _, err := utils.GetCollection(models.UserCollectionName).UpdateOne(context.Background(), filter, update); err != nil {
		utils.GetError(err, nttrep, response)
		return
	}

	utils.GetSuccess("Password update successful", nil, response)
}

func CreateToken(userid string) (string, error) {
	var err error
	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userid
	atClaims["exp"] = time.Now().Add(time.Minute * 60 * 24).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}
	return token, nil
}

func VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func TokenValid(r *http.Request) (string, error) {
	token, err := VerifyToken(r)
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		return "", err
	}
	return fmt.Sprintf("%v", claims["user_id"]), nil
}
func DeleteMapProps(m map[string]interface{}, s []string) {
	for _, v := range s {
		delete(m, v)
	}
}
