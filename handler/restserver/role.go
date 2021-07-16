package restserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/n-creativesystem/api-rbac/proto"
	"google.golang.org/grpc"
)

type roleHandle interface {
	create(*gin.Context)
	findById(*gin.Context)
	findAll(*gin.Context)
	update(*gin.Context)
	delete(*gin.Context)
	getPermissions(*gin.Context)
	addPermissions(*gin.Context)
	deletePermissions(*gin.Context)
}

type roleHandler struct {
	con *grpc.ClientConn
}

func newRoleHandler(con *grpc.ClientConn) *roleHandler {
	return &roleHandler{con: con}
}

func (h *roleHandler) grpcClient() proto.RoleClient {
	return proto.NewRoleClient(h.con)
}

func (h *roleHandler) create(c *gin.Context) {
	client := h.grpcClient()
	var req proto.RoleEntities
	if err := c.BindJSON(&req); requestError(c, err, body) {
		return
	}
	res, err := client.Create(c.Request.Context(), &req)
	if responseError(c, err) {
		return
	}
	Render(c, http.StatusOK, res)
}

func (h *roleHandler) findById(c *gin.Context) {
	client := h.grpcClient()
	id := c.Param("id")
	res, err := client.FindById(c.Request.Context(), &proto.RoleKey{
		Id: id,
	})
	if responseError(c, err) {
		return
	}
	Render(c, http.StatusOK, res)
}

func (h *roleHandler) findAll(c *gin.Context) {
	client := h.grpcClient()
	res, err := client.FindAll(c.Request.Context(), &proto.Empty{})
	if responseError(c, err) {
		return
	}
	Render(c, http.StatusOK, res)
}

func (h *roleHandler) update(c *gin.Context) {
	id := c.Param("id")
	var req proto.RoleUpdateEntity
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

func (h *roleHandler) delete(c *gin.Context) {
	id := c.Param("id")
	client := h.grpcClient()
	if _, err := client.Delete(c.Request.Context(), &proto.RoleKey{Id: id}); responseError(c, err) {
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *roleHandler) getPermissions(c *gin.Context) {
	id := c.Param("id")
	req := &proto.RoleKey{
		Id: id,
	}
	client := h.grpcClient()
	res, err := client.GetPermissions(c.Request.Context(), req)
	if responseError(c, err) {
		return
	}
	Render(c, http.StatusOK, res)
}

func (h *roleHandler) addPermissions(c *gin.Context) {
	id := c.Param("id")
	var req proto.RoleReleationPermissions
	if err := c.BindJSON(&req); requestError(c, err, body) {
		return
	}
	req.Id = id
	client := h.grpcClient()
	res, err := client.AddPermissions(c.Request.Context(), &req)
	if responseError(c, err) {
		return
	}
	Render(c, http.StatusOK, res)
}

func (h *roleHandler) deletePermissions(c *gin.Context) {
	roleId := c.Param("id")
	permissionId := c.Param("permissionId")
	client := h.grpcClient()
	_, err := client.DeletePermission(c.Request.Context(), &proto.RoleReleationPermission{
		Id:           roleId,
		PermissionId: permissionId,
	})
	if responseError(c, err) {
		return
	}
	c.Status(http.StatusNoContent)
}
