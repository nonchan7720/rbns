package restserver_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/n-creativesystem/api-rbac/handler/grpcserver"
	"github.com/n-creativesystem/api-rbac/handler/restserver"
	"github.com/n-creativesystem/api-rbac/infra/dao"
	"github.com/n-creativesystem/api-rbac/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis = bufconn.Listen(bufSize)

func bufDialer(ctx context.Context, addr string) (net.Conn, error) {
	return lis.Dial()
}

func newRequest(method, target string, body io.Reader) (*httptest.ResponseRecorder, *http.Request) {
	req := httptest.NewRequest(method, target, body)
	w := httptest.NewRecorder()
	return w, req
}

func obj2json(obj interface{}) io.Reader {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(obj); err != nil {
		panic(err)
	}
	return &buf
}

func json2obj(reader io.Reader, obj interface{}) {
	if err := json.NewDecoder(reader).Decode(obj); err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {
	_ = dao.New(dao.Migration, dao.MigrationBack)
	_ = m.Run()
}

func TestPermission(t *testing.T) {
	db := dao.New()
	grpcSrv := grpcserver.New(db)
	go func() {
		err := grpcSrv.Serve(lis)
		assert.NoError(t, err)
	}()
	con, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if !assert.NoError(t, err) {
		return
	}
	var req *http.Request
	var w *httptest.ResponseRecorder
	mpPermission := map[string]struct {
		name        string
		description string
	}{
		"1": {
			name:        "create:user",
			description: "cerate user permission",
		},
		"2": {
			name:        "read:user",
			description: "read user permission",
		},
	}
	restSrv := restserver.New(con)
	p := &proto.PermissionEntities{
		Permissions: []*proto.PermissionEntity{
			{
				Name:        mpPermission["1"].name,
				Description: mpPermission["1"].description,
			},
			{
				Name:        mpPermission["2"].name,
				Description: mpPermission["2"].description,
			},
		},
	}
	w, req = newRequest(http.MethodPost, "/api/v1/permissions", obj2json(p))
	restSrv.ServeHTTP(w, req)
	assert.Equal(t, w.Result().StatusCode, http.StatusOK)
	w, req = newRequest(http.MethodGet, "/api/v1/permissions/1", nil)
	restSrv.ServeHTTP(w, req)
	var res proto.PermissionEntities
	json2obj(w.Result().Body, &res)
	for _, permission := range res.GetPermissions() {
		id := permission.GetId()
		assert.Equal(t, mpPermission[id].name, permission.Name)
		assert.Equal(t, mpPermission[id].description, permission.GetDescription())
	}
	w, req = newRequest(http.MethodGet, "/api/v1/permissions", nil)
	restSrv.ServeHTTP(w, req)
	res = proto.PermissionEntities{}
	json2obj(w.Result().Body, &res)
	for _, permission := range res.GetPermissions() {
		id := permission.GetId()
		assert.Equal(t, mpPermission[id].name, permission.Name)
		assert.Equal(t, mpPermission[id].description, permission.GetDescription())
	}
}
