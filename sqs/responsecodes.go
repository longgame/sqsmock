package sqs

type ResponseCode int

const(
	ResponseOk ResponseCode = 200
	ResponseResouceNotFound = 404
	ResponseInternalServerError = 500
)
