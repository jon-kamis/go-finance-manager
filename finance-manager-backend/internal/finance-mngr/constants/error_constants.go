// Package constants contains application constants
package constants

//Generic Errors
const GenericForbiddenError = "forbidden"
const GenericServerError = "an unexpected error has occured"
const GenericNotFoundError = "not found"
const GenericBadRequestError = "bad request"

//DB Errors
const EntityNotFoundError = "entity not found"
const UnexpectedSQLError = "an unexpected error occured during the database call"
const InsertMultStockDataError = "one or more errors occured when inserting stock data"

const FailedToLoadUserError = "failed to load user"

//Default Error Respones
const JSONDefaultErrorMessage = "an unexpected error occured"

//Validation Errors
const UserRoleDoesNotBelongToUserError = "user role does not belong to the given user"
const BillDoesNotBelongToUserError = "bill does not belong to the given user"
const LoanDoesNotBelongToUserError = "loan does not belong to the given user"
const FailedToReadUserIdFromAuthHeaderError = "failed to read the logged in user's ID from the auth header"
const UserForbiddenToViewOtherUserDataError = "user is forbidden from viewing other user data"
const FailedToParseIdError = "failed to parse id"
const InvalidCreditCardError = "credit card is invalid"
const UsernameOrEmailExistError = "username or email already exists"

//Auth Errors
const InvalidAuthHeaderError = "authorization header is invalid"
const InvalidSigningMethodError = "unexpected signing method"
const ExpiredTokenError = "token is expired"
const InvalidIssuerError = "invalid issuer"

//External Calls
const UnexpectedExternalCallError = "unexpected error was returned when making external API call"
const FailedToParseJsonBodyError = "failed to unmarshal json payload"
const UnexpectedResponseCodeError = "unexpected response code during remote call"

//Stock Errors
const StockOperationInvalidOperationError = "invalid stock operation"
const StockOperationInvalidDateError = "date is required and cannot be a future date"
const StockOperationInvalidAmountError = "amount must be greater than 0"
const StockOperationTickerRequiredError = "ticker is required"
const StockOperationAlreadyExistsError = "a stock operation already exists for the given time"
const StockOperationBelowZeroError = "stock operations cannot result in a quantity below 0"