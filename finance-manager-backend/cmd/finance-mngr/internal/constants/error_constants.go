package constants

//DB Errors
const EntityNotFoundError = "entity not found"
const UnexpectedSQLError = "an unexpected error occured during the database call"

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
