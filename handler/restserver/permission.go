package restserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/n-creativesystem/rbns/proto"
	"google.golang.org/grpc"
)

type permissionHandle interface {
	create(*gin.Context)
	findById(*gin.Context)
	findAll(*gin.Context)
	update(*gin.Context)
	delete(*gin.Context)
}

type permissionHandler struct {
	con *grpc.ClientConn
}

func newPermissionHandler(con *grpc.ClientConn) *permissionHandler {
	return &permissionHandler{con: con}
}

func (h *permissionHandler) grpcClient() proto.PermissionClient {
	return proto.NewPermissionClient(h.con)
}

func (h *permissionHandler) create(c *gin.Context) {
	client := h.grpcClient()
	var req proto.PermissionEntities
	if err := c.BindJSON(&req); requestError(c, err, body) {
		return
	}
	res, err := client.Create(c.Request.Context(), &req)
	if responseError(c, err) {
		return
	}
	Render(c, http.StatusOK, res)
}

func (h *permissionHandler) findById(c *gin.Context) {
	id := c.Param("id")
	client := h.grpcClient()
	res, err := client.FindById(c.Request.Context(), &proto.PermissionKey{
		Id: id,
	})
	if responseError(c, err) {
		return
	}
	Render(c, http.StatusOK, res)
}

func (h *permissionHandler) findAll(c *gin.Context) {
	client := h.grpcClient()
	res, err := client.FindAll(c.Request.Context(), &proto.Empty{})
	if responseError(c, err) {
		return
	}
	Render(c, http.StatusOK, res)
}

func (h *permissionHandler) update(c *gin.Context) {
	id := c.Param("id")
	var req proto.PermissionEntity
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

func (h *permissionHandler) delete(c *gin.Context) {
	id := c.Param("id")
	client := h.grpcClient()
	if _, err := client.Delete(c.Request.Context(), &proto.PermissionKey{Id: id}); responseError(c, err) {
		return
	}
	c.Status(http.StatusNoContent)
}
