package handler

import "github.com/SawitProRecruitment/UserService/generated"

var (
	BadRequestErr = generated.BadRequest{Message: "Invalid input parameter(s)"}
	ForbiddenErr  = generated.Forbidden{Message: "Access to the requested resource is forbidden"}
	ConflictErr   = generated.Conflict{Message: "The request is conflict with the current state of the target resource"}
	InternalErr   = generated.Conflict{Message: "The server encountered an unexpected condition"}
)
