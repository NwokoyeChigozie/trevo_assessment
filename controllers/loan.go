package controllers

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/gregoflash05/trove/models"
	"github.com/gregoflash05/trove/utils"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TakeLoan(response http.ResponseWriter, request *http.Request) {
	userID, err := TokenValid(request)
	if err != nil {
		utils.GetError(fmt.Errorf("token Invalid"), nttrep, response)
		return
	}

	var rBody models.TakeLoanRequest
	if e := utils.ParseJSONFromRequest(request, &rBody); e != nil {
		utils.GetError(e, nttrep, response)
		return
	}

	objID, _ := primitive.ObjectIDFromHex(userID)
	user, userA := models.User{}, models.User{}

	userMap, err := utils.GetMongoDBDoc(models.UserCollectionName, bson.M{"_id": objID})

	if err != nil {
		utils.GetError(errors.New("user not found"), nttrep, response)
		return
	}

	bsonBytes, err := bson.Marshal(userMap)

	if err != nil {
		utils.GetError(errors.New("user not found"), nttrep, response)
		return
	}

	err = bson.Unmarshal(bsonBytes, &userA)
	if err != nil {
		utils.GetError(errors.New("user not found"), nttrep, response)
		return
	}

	err = mapstructure.Decode(userMap, &user)
	if err != nil {
		utils.GetError(err, nttrep, response)
		return
	}
	user.PortfolioPositions = userA.PortfolioPositions

	totalPositionsValue := TotalPositionsValue(user.PortfolioPositions)

	res, _ := utils.GetMongoDBDoc(models.LoanCollection, bson.M{"user_id": userID})

	if res == nil {
		// create loan
		err = ValidateTakeLoan(rBody, totalPositionsValue)
		if err != nil {
			utils.GetError(err, nttrep, response)
			return
		}

		err = CreateLoan(userID, rBody)
		if err != nil {
			utils.GetError(err, nttrep, response)
			return
		}

	} else {
		err = ValidateTakeLoan(rBody, totalPositionsValue)
		if err != nil {
			utils.GetError(err, nttrep, response)
			return
		}

		loanD := models.Loan{}
		bsonBytes, _ := bson.Marshal(res)
		bson.Unmarshal(bsonBytes, &loanD)

		if loanD.Balance > 0 {
			utils.GetError(fmt.Errorf("you have outstanding loan, pay back to take more loan"), nttrep, response)
			return
		}

		// update loan
		err = UpdateLoan(loanD.ID, userID, rBody.Amount, rBody.Amount, rBody.Duration, int(time.Now().Unix()), int(time.Now().Unix()))
		if err != nil {
			utils.GetError(err, nttrep, response)
			return
		}
	}
	utils.GetSuccess("Loan taken", "", response)
}

func GetLoan(response http.ResponseWriter, request *http.Request) {
	userID, err := TokenValid(request)
	if err != nil {
		utils.GetError(fmt.Errorf("token Invalid"), nttrep, response)
		return
	}

	res, _ := utils.GetMongoDBDoc(models.LoanCollection, bson.M{"user_id": userID})

	if res != nil {
		subLoan, layoutUS := models.Loan{}, "January 2, 2006"
		bsonBytes, _ := bson.Marshal(res)
		bson.Unmarshal(bsonBytes, &subLoan)

		// calculating prorated payment
		submonths, monthSecs := int(math.Round(float64(subLoan.Duration/(60*60*24*30)))), 60*60*24*30
		installment := subLoan.TotalAmount / float64(submonths)
		timeFraction := int(math.Floor((float64(time.Now().Unix()) - float64(subLoan.TimeTaken)) / float64(monthSecs)))
		amountDue := math.Round(float64((timeFraction*int(installment))*100)) / 100
		amountPaid := subLoan.TotalAmount - subLoan.Balance
		actualDue := amountDue - amountPaid
		whenNextDue := subLoan.TimeTaken + ((timeFraction + 1) * monthSecs)
		amountNextDue := (math.Round(float64(((timeFraction+1)*int(installment))*100)) / 100) - amountPaid
		if actualDue < 0 {
			actualDue = 0
		}
		if amountNextDue < 0 {
			amountNextDue = 0
		}

		t := time.Unix(int64(whenNextDue), 0)
		strDate := t.Format(layoutUS)
		res["whenNextDue"], res["amountNextDue"] = strDate, amountNextDue
		res["amountDue"], res["actualDue"] = amountDue, actualDue

		utils.GetSuccess("loan retrieved successfully", res, response)
		return

	}

	utils.GetSuccess("loan retrieved successfully", res, response)

}

func PayBackLoan(response http.ResponseWriter, request *http.Request) {
	userID, err := TokenValid(request)
	if err != nil {
		utils.GetError(fmt.Errorf("token Invalid"), nttrep, response)
		return
	}
	res, _ := utils.GetMongoDBDoc(models.LoanCollection, bson.M{"user_id": userID})
	loanD := models.Loan{}
	bsonBytes, _ := bson.Marshal(res)
	bson.Unmarshal(bsonBytes, &loanD)

	var rBody models.PayBackLoanRequest
	if e := utils.ParseJSONFromRequest(request, &rBody); e != nil {
		utils.GetError(e, nttrep, response)
		return
	}

	if rBody.Amount > loanD.Balance {
		utils.GetError(fmt.Errorf("payment is beyond your loan balance"), nttrep, response)
		return

	}

	err = UpdateLoan(loanD.ID, userID, loanD.TotalAmount, math.Round((loanD.Balance-rBody.Amount)*100)/100, loanD.Duration, loanD.TimeTaken, int(time.Now().Unix()))
	if err != nil {
		utils.GetError(err, nttrep, response)
		return
	}

	utils.GetSuccess("Payment successfully", res, response)
}

func CreateLoan(userID string, rBody models.TakeLoanRequest) error {
	newLoan := models.Loan{
		UserID:      userID,
		TotalAmount: rBody.Amount,
		Balance:     rBody.Amount,
		Duration:    rBody.Duration,
		TimeTaken:   int(time.Now().Unix()),
		LastPayment: int(time.Now().Unix()),
	}

	detail, _ := utils.StructToMap(newLoan)

	_, err := utils.CreateMongoDBDoc(models.LoanCollection, detail)

	if err != nil {
		return err
	}
	return nil
}

func UpdateLoan(loanId string, userID string, totalAmount, Balance float64, duration, Timetaken, lastPayment int) error {
	newLoan := models.Loan{
		UserID:      userID,
		TotalAmount: totalAmount,
		Balance:     Balance,
		Duration:    duration,
		TimeTaken:   Timetaken,
		LastPayment: lastPayment,
	}

	detail, _ := utils.StructToMap(newLoan)

	updateFields := make(map[string]interface{})

	for key, value := range detail {
		if value != "" {
			updateFields[key] = value
		}
	}

	_, err := utils.UpdateOneMongoDBDoc(models.LoanCollection, loanId, updateFields)

	if err != nil {
		return err
	}
	return nil
}

func ValidateTakeLoan(rbody models.TakeLoanRequest, totalPositionsValue float64) error {
	if rbody.Amount < 1 {
		return fmt.Errorf("please enter a valid amout")

	} else if rbody.Amount > (0.6 * totalPositionsValue) {
		return fmt.Errorf("you cannot loan more than 60 percent of total value")
	} else if (rbody.Duration < 15551999) || (rbody.Duration > 31104001) {
		return fmt.Errorf("duration must be withing 6 to 12 months")
	}
	return nil
}

func TotalPositionsValue(positions []models.PortfolioPosition) float64 {
	var total float64 = 0
	for _, positon := range positions {
		total += positon.EquityValue
	}
	return total
}
