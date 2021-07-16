package service

// func TestService(t *testing.T) {
// 	db := dao.New(dao.Migration, dao.MigrationBack, dao.Debug)
// 	var repo repository.Repository = infra.NewRepository(db)
// 	ctx := context.Background()
// 	rSrv := NewRoleService(repo)
// 	oSrv := NewOrganizationService(repo)
// 	uSrv := NewUserService(repo)

// 	cases := []tests.Case{
// 		{
// 			Name: "Role create",
// 			Fn: func(t *testing.T) {
// 				out, err := rSrv.Create(ctx, &proto.RoleEntities{
// 					Roles: []*proto.RoleEntity{
// 						{
// 							Name:        "admin",
// 							Description: "administrator",
// 						},
// 						{
// 							Name:        "user",
// 							Description: "user",
// 						},
// 						{
// 							Name:        "guest",
// 							Description: "guest",
// 						},
// 					},
// 				})
// 				expecteds := []struct {
// 					id                uint64
// 					name, description string
// 				}{
// 					{
// 						id:          1,
// 						name:        "admin",
// 						description: "administrator",
// 					},
// 					{
// 						id:          2,
// 						name:        "user",
// 						description: "user",
// 					},
// 					{
// 						id:          3,
// 						name:        "guest",
// 						description: "guest",
// 					},
// 				}
// 				assert.NoError(t, err)
// 				for idx, entity := range out.GetRoles() {
// 					expected := expecteds[idx]
// 					assert.Equal(t, expected.id, entity.GetId())
// 					assert.Equal(t, expected.name, entity.GetName())
// 					assert.Equal(t, expected.description, entity.GetDescription())
// 				}
// 			},
// 		},
// 		{
// 			Name: "Role find",
// 			Fn: func(t *testing.T) {
// 				out, err := rSrv.FindById(ctx, &proto.RoleKey{Id: 1})
// 				assert.NoError(t, err)
// 				assert.Equal(t, "admin", out.GetName())
// 			},
// 		},
// 		{
// 			Name: "Role findall",
// 			Fn: func(t *testing.T) {
// 				out, err := rSrv.FindAll(ctx, &proto.Empty{})
// 				expecteds := []struct {
// 					id                uint64
// 					name, description string
// 				}{
// 					{
// 						id:          1,
// 						name:        "admin",
// 						description: "administrator",
// 					},
// 					{
// 						id:          2,
// 						name:        "user",
// 						description: "user",
// 					},
// 					{
// 						id:          3,
// 						name:        "guest",
// 						description: "guest",
// 					},
// 				}
// 				assert.NoError(t, err)
// 				for idx, entity := range out.GetRoles() {
// 					expected := expecteds[idx]
// 					assert.Equal(t, expected.id, entity.GetId())
// 					assert.Equal(t, expected.name, entity.GetName())
// 					assert.Equal(t, expected.description, entity.GetDescription())
// 				}
// 			},
// 		},
// 		{
// 			Name: "Role add permissions",
// 			Fn: func(t *testing.T) {
// 				_, err := rSrv.AddPermissions(ctx, &proto.RoleReleationPermission{
// 					Id: 1,
// 					Permissions: []*proto.PermissionKey{
// 						{
// 							Id: 1,
// 						},
// 						{
// 							Id: 2,
// 						},
// 					},
// 				})
// 				assert.NoError(t, err)
// 			},
// 		},
// 		{
// 			Name: "Organization create",
// 			Fn: func(t *testing.T) {
// 				out, err := oSrv.Create(ctx, &proto.OrganizationEntity{
// 					Name:        "organization",
// 					Description: "organization test",
// 				})
// 				assert.NoError(t, err)
// 				assert.Equal(t, uint64(1), out.Id)
// 				assert.Equal(t, "organization", out.GetName())
// 				assert.Equal(t, "organization test", out.GetDescription())
// 				out, err = oSrv.Create(ctx, &proto.OrganizationEntity{
// 					Name:        "organization2",
// 					Description: "organization2 test",
// 				})
// 				assert.NoError(t, err)
// 				assert.Equal(t, uint64(2), out.Id)
// 				assert.Equal(t, "organization2", out.GetName())
// 				assert.Equal(t, "organization2 test", out.GetDescription())
// 			},
// 		},
// 		{
// 			Name: "Organization find",
// 			Fn: func(t *testing.T) {
// 				out, err := oSrv.FindById(ctx, &proto.OrganizationKey{Id: 1})
// 				assert.NoError(t, err)
// 				assert.Equal(t, uint64(1), out.Id)
// 				assert.Equal(t, "organization", out.GetName())
// 				assert.Equal(t, "organization test", out.GetDescription())
// 			},
// 		},
// 		{
// 			Name: "Organization findall",
// 			Fn: func(t *testing.T) {
// 				expecteds := []struct {
// 					id                uint64
// 					name, description string
// 				}{
// 					{
// 						id:          1,
// 						name:        "organization",
// 						description: "organization test",
// 					},
// 					{
// 						id:          2,
// 						name:        "organization2",
// 						description: "organization2 test",
// 					},
// 				}
// 				out, err := oSrv.FindAll(ctx, &proto.Empty{})
// 				assert.NoError(t, err)
// 				for i, entity := range out.Entities {
// 					expected := expecteds[i]
// 					assert.Equal(t, expected.id, entity.GetId())
// 					assert.Equal(t, expected.name, entity.GetName())
// 					assert.Equal(t, expected.description, entity.GetDescription())
// 				}
// 			},
// 		},
// 		{
// 			Name: "User create",
// 			Fn: func(t *testing.T) {
// 				users := []struct {
// 					key            string
// 					organizationID uint64
// 				}{
// 					{
// 						key:            "user1",
// 						organizationID: 1,
// 					},
// 					{
// 						key:            "user1",
// 						organizationID: 2,
// 					},
// 					{
// 						key:            "user2",
// 						organizationID: 1,
// 					},
// 					{
// 						key:            "user2",
// 						organizationID: 2,
// 					},
// 				}
// 				for _, user := range users {
// 					_, err := uSrv.Create(ctx, &proto.UserKey{
// 						Key:            user.key,
// 						OrganizationId: user.organizationID,
// 					})
// 					assert.NoError(t, err)
// 				}
// 			},
// 		},
// 		{
// 			Name: "User add role",
// 			Fn: func(t *testing.T) {
// 				_, err := uSrv.AddRole(ctx, &proto.UserRole{
// 					User: &proto.UserKey{Key: "user1", OrganizationId: 1},
// 					Roles: []*proto.RoleKey{
// 						{
// 							Id: 1,
// 						},
// 					},
// 				})
// 				assert.NoError(t, err)
// 			},
// 		},
// 		{
// 			Name: "User find",
// 			Fn: func(t *testing.T) {
// 				mpRoles := map[uint64]bool{
// 					1: true,
// 				}
// 				mpPermissions := map[uint64]bool{
// 					1: true, 2: true, 3: true, 4: true,
// 				}
// 				u, err := uSrv.FindByKey(ctx, &proto.UserKey{Key: "user1", OrganizationId: 1})
// 				assert.NoError(t, err)
// 				assert.Equal(t, "user1", u.GetKey())
// 				assert.Equal(t, uint64(1), u.GetOrganizationId())
// 				for _, role := range u.GetRoles() {
// 					assert.True(t, mpRoles[role.GetId()])
// 				}
// 				for _, permission := range u.GetPermissions() {
// 					assert.True(t, mpPermissions[permission.GetId()])
// 				}
// 			},
// 		},
// 		{
// 			Name: "User delete",
// 			Fn:   func(t *testing.T) {},
// 		},
// 	}

// 	for _, tt := range cases {
// 		t.Run(tt.Name, tt.Fn)
// 	}
// }
