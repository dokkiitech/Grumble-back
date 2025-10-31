package api

import (
	"net/http"

	"github.com/dokkiitech/grumble-back/internal/controller"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// ServerImpl implements the OpenAPI ServerInterface by delegating to controllers
type ServerImpl struct {
	grumbleController  *controller.GrumbleController
	timelineController *controller.TimelineController
	authController     *controller.AuthController
	vibeController     *controller.VibeController
}

// NewServerImpl wires controllers into a ServerImpl
func NewServerImpl(
	gc *controller.GrumbleController,
	tc *controller.TimelineController,
	ac *controller.AuthController,
	vc *controller.VibeController,
) *ServerImpl {
	return &ServerImpl{
		grumbleController:  gc,
		timelineController: tc,
		authController:     ac,
		vibeController:     vc,
	}
}

// GetEvents implements ServerInterface (placeholder for future implementation)
func (s *ServerImpl) GetEvents(c *gin.Context, params GetEventsParams) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Events not yet implemented"})
}

// GetEvent implements ServerInterface (placeholder for future implementation)
func (s *ServerImpl) GetEvent(c *gin.Context, eventID int) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Events not yet implemented"})
}

// GetGrumbles implements ServerInterface
func (s *ServerImpl) GetGrumbles(c *gin.Context, params GetGrumblesParams) {
	s.timelineController.GetGrumbles(
		c,
		params.UserID,
		params.ToxicLevelMin,
		params.ToxicLevelMax,
		params.UnpurifiedOnly,
		params.Limit,
		params.Offset,
	)
}

// CreateGrumble implements ServerInterface
func (s *ServerImpl) CreateGrumble(c *gin.Context) {
	s.grumbleController.CreateGrumble(c)
}

// AddVibe implements ServerInterface
func (s *ServerImpl) AddVibe(c *gin.Context, grumbleID openapi_types.UUID) {
	s.vibeController.PostVibe(c, uuid.UUID(grumbleID))
}

// GetMyProfile implements ServerInterface
func (s *ServerImpl) GetMyProfile(c *gin.Context) {
	s.authController.GetMyProfile(c)
}
