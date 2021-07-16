package restserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/n-creativesystem/api-rbac/proto"
	"google.golang.org/grpc"
)

type userHandle interface {
	create(*gin.Context)
	findByKey(*gin.Context)
	delete(*gin.Context)
	addRole(*gin.Context)
	deleteRole(*gin.Context)
}

type userHandler struct {
	con *grpc.ClientConn
}

func newUserHandler(con *grpc.ClientConn) *userHandler {
	return &userHandler{con: con}
}

func (h *userHandler) grpcClient() proto.UserClient {
	return proto.NewUserClient(h.con)
}

func (h *userHandler) create(c *gin.Context) {
	var user proto.UserEntity
	if err := c.BindJSON(&user); requestError(c, err, body) {
		return
	}
	client := h.grpcClient()
	if _, err := client.Create(c.Request.Context(), &user); responseError(c, err) {
		return
	}
	c.Status(http.StatusOK)
}

func (h *userHandler) findByKey(c *gin.Context) {
	id := c.Param("id")
	key := c.Param("key")
	client := h.grpcClient()
	res, err := client.FindByKey(c.Request.Context(), &proto.UserKey{Key: key, OrganizationId: id})
	if responseError(c, err) {
		return
	}
	Render(c, http.StatusOK, res)
}

func (h *userHandler) delete(c *gin.Context) {
	id := c.Param("id")
	key := c.Param("key")
	client := h.grpcClient()
	if _, err := client.Delete(c.Request.Context(), &proto.UserKey{Key: key, OrganizationId: id}); responseError(c, err) {
		return
	}
	c.Status(http.StatusOK)
}

func (h *userHandler) addRole(c *gin.Context) {
	id := c.Param("id")
	key := c.Param("key")
	var reqbody proto.UserRole
	if err := c.BindJSON(&reqbody); requestError(c, err, body) {
		return
	}
	reqbody.User = &proto.UserKey{
		Key:            key,
		OrganizationId: id,
	}
	client := h.grpcClient()
	if _, err := client.AddRole(c.Request.Context(), &reqbody); responseError(c, err) {
		return
	}
	c.Status(http.StatusCreated)
}

func (h *userHandler) deleteRole(c *gin.Context) {
	id := c.Param("id")
	key := c.Param("key")
	roleId := c.Param("roleId")
	client := h.grpcClient()
	_, err := client.DeleteRole(c.Request.Context(), &proto.UserDeleteRole{
		User: &proto.UserKey{
			OrganizationId: id,
			Key:            key,
		},
		Role: &proto.RoleKey{
			Id: roleId,
		},
	})
	if responseError(c, err) {
		return
	}
	c.Status(http.StatusNoContent)
}
