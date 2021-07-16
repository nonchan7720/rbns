package restserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/n-creativesystem/api-rbac/proto"
	"google.golang.org/grpc"
)

type organizationHandle interface {
	create(*gin.Context)
	findById(*gin.Context)
	findAll(*gin.Context)
	update(*gin.Context)
	delete(*gin.Context)
}

type organizationHandler struct {
	con *grpc.ClientConn
}

func newOrganizationHandler(con *grpc.ClientConn) *organizationHandler {
	return &organizationHandler{con: con}
}

func (h *organizationHandler) grpcClient() proto.OrganizationClient {
	return proto.NewOrganizationClient(h.con)
}

func (h *organizationHandler) create(c *gin.Context) {
	var req proto.OrganizationEntity
	if err := c.Bind(&req); requestError(c, err, body) {
		return
	}
	client := h.grpcClient()
	res, err := client.Create(c.Request.Context(), &req)
	if responseError(c, err) {
		return
	}
	Render(c, http.StatusOK, res)
}

func (h *organizationHandler) findById(c *gin.Context) {
	id := c.Param("id")
	client := h.grpcClient()
	res, err := client.FindById(c.Request.Context(), &proto.OrganizationKey{Id: id})
	if responseError(c, err) {
		return
	}
	Render(c, http.StatusOK, res)
}

func (h *organizationHandler) findAll(c *gin.Context) {
	client := h.grpcClient()
	res, err := client.FindAll(c.Request.Context(), &proto.Empty{})
	if responseError(c, err) {
		return
	}
	Render(c, http.StatusOK, res)
}

func (h *organizationHandler) update(c *gin.Context) {
	id := c.Param("id")
	var req proto.OrganizationUpdateEntity
	if err := c.BindJSON(&req); requestError(c, err, body) {
		return
	}
	req.Id = id
	client := h.grpcClient()
	if _, err := client.Update(c.Request.Context(), &req); responseError(c, err) {
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *organizationHandler) delete(c *gin.Context) {
	id := c.Param("id")
	client := h.grpcClient()
	if _, err := client.Delete(c.Request.Context(), &proto.OrganizationKey{Id: id}); responseError(c, err) {
		return
	}
	c.Status(http.StatusNoContent)
}
