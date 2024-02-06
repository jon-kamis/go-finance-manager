// Package constants contains application constants
package constants

//Generic Errors
const GenericForbiddenError = "forbidden"
const GenericServerError = "an unexpected error has occured"
const GenericNotFoundError = "not found"
const GenericBadRequestError = "bad request"

//DB Errors
const EntityNotFoundError = "entity not found"
const UnexpectedSQLError = "an unexpected error occured during the database call: %v"
const InsertMultStockDataError = "one or more errors occured when inserting stock data\n%v"
const FailedToDeleteEntityError = "failed to delete the requested entity: \n%v"
const FailedToUpdateEntityError = "failed to update the requested entity: \n%v"
const FailedToSaveEntityError = "failed to save the requested entity: \n%v"
const FailedToRetrieveEntityError = "failed to retrieve entity(s) from the database: \n%v"

const FailedToLoadUserError = "failed to load user\n%v"

//Authentication
const InvalidCredentialsError = "invalid credentials"

//Default Error Respones
const JSONDefaultErrorMessage = "an unexpected error occured\n%v"

//Validation Errors
const EntityDoesNotBelongToUserError = "requested entity does not belong to the requesting user: \n%v"
const ProcessUserIdError = "an error occured when attempting to process the id: \n%v"
const ProcessIdError = "an error occured when attempting to process the id: \n%v"
const UserRoleDoesNotBelongToUserError = "user role does not belong to the given user: \n%v"
const BillDoesNotBelongToUserError = "bill does not belong to the given user: \n%v"
const LoanDoesNotBelongToUserError = "loan does not belong to the given user: \n%v"
const FailedToReadUserIdFromAuthHeaderError = "failed to read the logged in user's ID from the auth header: \n%v"
const UserForbiddenToViewOtherUserDataError = "user is forbidden from viewing other user data: \n%v"
const FailedToParseIdError = "failed to parse id: \n%v"
const InvalidCreditCardError = "credit card is invalid: \n%v"
const UsernameOrEmailExistError = "username or email already exists: \n%v"

//Auth Errors
const InvalidAuthHeaderError = "authorization header is invalid\n%v"
const InvalidSigningMethodError = "unexpected signing method\n%v"
const ExpiredTokenError = "token is expired\n%v"
const InvalidIssuerError = "invalid issuer\n%v"

//External Calls
const UnexpectedExternalCallError = "unexpected error was returned when making external API call\n%v"
const FailedToParseJsonBodyError = "failed to unmarshal json payload:\n%v"
const UnexpectedResponseCodeError = "unexpected response code during remote call\n%v"

//Stock Errors
const StockOperationInvalidOperationError = "invalid stock operation"
const StockOperationInvalidDateError = "date is required and cannot be a future date"
const StockOperationInvalidAmountError = "amount must be greater than 0"
const StockOperationTickerRequiredError = "ticker is required"
const StockOperationAlreadyExistsError = "a stock operation already exists for the given time"
const StockOperationBelowZeroError = "stock operations cannot result in a quantity below 0"
