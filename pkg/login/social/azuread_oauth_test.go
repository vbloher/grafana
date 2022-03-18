package social

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/grafana/grafana/pkg/infra/log"
	"github.com/grafana/grafana/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

func TestSocialAzureAD_UserInfo(t *testing.T) {
	type fields struct {
		allowedGroups       []string
		autoAssignOrgRole   string
		roleAttributeStrict bool
	}
	type args struct {
		client *http.Client
	}

	tests := []struct {
		name                     string
		fields                   fields
		claims                   *azureClaims
		args                     args
		settingAutoAssignOrgRole string
		want                     *BasicUserInfo
		wantErr                  bool
	}{
		{
			name: "Email in email claim",
			claims: &azureClaims{
				Email:             "me@example.com",
				PreferredUsername: "",
				Roles:             []string{},
				Name:              "My Name",
				ID:                "1234",
			},
			fields: fields{
				autoAssignOrgRole: "Viewer",
			},
			want: &BasicUserInfo{
				Id:      "1234",
				Name:    "My Name",
				Email:   "me@example.com",
				Login:   "me@example.com",
				Company: "",
				OrgMemberships: map[int64]models.RoleType{
					1: models.ROLE_VIEWER,
				},
				Groups: []string{},
			},
		},
		{
			name: "No email",
			claims: &azureClaims{
				Email:             "",
				PreferredUsername: "",
				Roles:             []string{},
				Name:              "My Name",
				ID:                "1234",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "No id token",
			claims:  nil,
			want:    nil,
			wantErr: true,
		},
		{
			name: "Email in preferred_username claim",
			claims: &azureClaims{
				Email:             "",
				PreferredUsername: "me@example.com",
				Roles:             []string{},
				Name:              "My Name",
				ID:                "1234",
			},
			fields: fields{
				autoAssignOrgRole: "Viewer",
			},
			want: &BasicUserInfo{
				Id:      "1234",
				Name:    "My Name",
				Email:   "me@example.com",
				Login:   "me@example.com",
				Company: "",
				OrgMemberships: map[int64]models.RoleType{
					1: models.ROLE_VIEWER,
				},
				Groups: []string{},
			},
		},
		{
			name: "Admin role",
			claims: &azureClaims{
				Email:             "me@example.com",
				PreferredUsername: "",
				Roles:             []string{"Admin"},
				Name:              "My Name",
				ID:                "1234",
			},
			want: &BasicUserInfo{
				Id:      "1234",
				Name:    "My Name",
				Email:   "me@example.com",
				Login:   "me@example.com",
				Company: "",
				OrgMemberships: map[int64]models.RoleType{
					1: models.ROLE_ADMIN,
				},
				Groups: []string{},
			},
		},
		{
			name: "Lowercase Admin role",
			claims: &azureClaims{
				Email:             "me@example.com",
				PreferredUsername: "",
				Roles:             []string{"admin"},
				Name:              "My Name",
				ID:                "1234",
			},
			want: &BasicUserInfo{
				Id:      "1234",
				Name:    "My Name",
				Email:   "me@example.com",
				Login:   "me@example.com",
				Company: "",
				OrgMemberships: map[int64]models.RoleType{
					1: models.ROLE_ADMIN,
				},
				Groups: []string{},
			},
		},
		{
			name: "Only other roles",
			claims: &azureClaims{
				Email:             "me@example.com",
				PreferredUsername: "",
				Roles:             []string{"AppAdmin"},
				Name:              "My Name",
				ID:                "1234",
			},
			want: &BasicUserInfo{
				Id:      "1234",
				Name:    "My Name",
				Email:   "me@example.com",
				Login:   "me@example.com",
				Company: "",
				OrgMemberships: map[int64]models.RoleType{
					1: models.ROLE_VIEWER,
				},
				Groups: []string{},
			},
		},
		{
			name: "role from env variable",
			claims: &azureClaims{
				Email:             "me@example.com",
				PreferredUsername: "",
				Roles:             []string{},
				Name:              "My Name",
				ID:                "1234",
			},
			fields: fields{
				autoAssignOrgRole: "Editor",
			},
			want: &BasicUserInfo{
				Id:      "1234",
				Name:    "My Name",
				Email:   "me@example.com",
				Login:   "me@example.com",
				Company: "",
				Role:    "Editor",
				Groups:  []string{},
			},
		},
		{
			name: "Editor role",
			claims: &azureClaims{
				Email:             "me@example.com",
				PreferredUsername: "",
				Roles:             []string{"Editor"},
				Name:              "My Name",
				ID:                "1234",
			},
			want: &BasicUserInfo{
				Id:      "1234",
				Name:    "My Name",
				Email:   "me@example.com",
				Login:   "me@example.com",
				Company: "",
				OrgMemberships: map[int64]models.RoleType{
					1: models.ROLE_EDITOR,
				},
				Groups: []string{},
			},
		},
		{
			name: "Admin and Editor roles in claim",
			claims: &azureClaims{
				Email:             "me@example.com",
				PreferredUsername: "",
				Roles:             []string{"Admin", "Editor"},
				Name:              "My Name",
				ID:                "1234",
			},
			want: &BasicUserInfo{
				Id:      "1234",
				Name:    "My Name",
				Email:   "me@example.com",
				Login:   "me@example.com",
				Company: "",
				OrgMemberships: map[int64]models.RoleType{
					1: models.ROLE_ADMIN,
				},
				Groups: []string{},
			},
		},
		{
			name: "Error if user is not a member of allowed_groups",
			fields: fields{
				allowedGroups: []string{"dead-beef"},
			},
			claims: &azureClaims{
				Email:             "me@example.com",
				PreferredUsername: "",
				Roles:             []string{},
				Groups:            []string{"foo", "bar"},
				Name:              "My Name",
				ID:                "1234",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Error if user is a member of allowed_groups",
			fields: fields{
				allowedGroups:     []string{"foo", "bar"},
				autoAssignOrgRole: "Viewer",
			},
			claims: &azureClaims{
				Email:             "me@example.com",
				PreferredUsername: "",
				Roles:             []string{},
				Groups:            []string{"foo"},
				Name:              "My Name",
				ID:                "1234",
			},
			want: &BasicUserInfo{
				Id:      "1234",
				Name:    "My Name",
				Email:   "me@example.com",
				Login:   "me@example.com",
				Company: "",
				OrgMemberships: map[int64]models.RoleType{
					1: models.ROLE_VIEWER,
				},
				Groups: []string{"foo"},
			},
		},
		{
			name: "Fetch groups when ClaimsNames and ClaimsSources is set",
			fields: fields{
				SocialBase: newSocialBase("azuread", &oauth2.Config{}, &OAuthInfo{}),
			},
			claims: &azureClaims{
				ID:                "1",
				Name:              "test",
				PreferredUsername: "test",
				Email:             "test@test.com",
				Roles:             []string{"Viewer"},
				ClaimNames:        claimNames{Groups: "src1"},
				ClaimSources:      nil, // set by the test
			},
			settingAutoAssignOrgRole: "",
			want: &BasicUserInfo{
				Id:     "1",
				Name:   "test",
				Email:  "test@test.com",
				Login:  "test@test.com",
				Role:   "Viewer",
				Groups: []string{"from_server"},
			},
			wantErr: false,
		},
		{
			name: "Fetch empty role when strict attribute role is true and no match",
			fields: fields{
				roleAttributeStrict: true,
			},
			claims: &azureClaims{
				Email:             "me@example.com",
				PreferredUsername: "",
				Roles:             []string{"foo"},
				Groups:            []string{},
				Name:              "My Name",
				ID:                "1234",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SocialAzureAD{
				SocialBase: &SocialBase{
					log: log.New("test"),
				},
				allowedGroups:       tt.fields.allowedGroups,
				autoAssignOrgRole:   tt.fields.autoAssignOrgRole,
				roleAttributeStrict: tt.fields.roleAttributeStrict,
			}

			key := []byte("secret")
			sig, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.HS256, Key: key}, (&jose.SignerOptions{}).WithType("JWT"))
			require.NoError(t, err)

			cl := jwt.Claims{
				Subject:   "subject",
				Issuer:    "issuer",
				NotBefore: jwt.NewNumericDate(time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC)),
				Audience:  jwt.Audience{"leela", "fry"},
			}

			var raw string
			if tt.claims != nil {
				if tt.claims.ClaimNames.Groups != "" {
					server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
						tokenParts := strings.Split(request.Header.Get("Authorization"), " ")
						require.Len(t, tokenParts, 2)
						require.Equal(t, "fake_token", tokenParts[1])

						writer.WriteHeader(http.StatusOK)

						type response struct {
							Value []string
						}
						res := response{Value: []string{"from_server"}}
						require.NoError(t, json.NewEncoder(writer).Encode(&res))
					}))
					// need to set the fake servers url as endpoint to capture request
					tt.claims.ClaimSources = map[string]claimSource{
						tt.claims.ClaimNames.Groups: {Endpoint: server.URL},
					}
				}
				raw, err = jwt.Signed(sig).Claims(cl).Claims(tt.claims).CompactSerialize()
				require.NoError(t, err)
			} else {
				raw, err = jwt.Signed(sig).Claims(cl).CompactSerialize()
				require.NoError(t, err)
			}

			token := &oauth2.Token{
				AccessToken: "fake_token",
			}
			if tt.claims != nil {
				token = token.WithExtra(map[string]interface{}{"id_token": raw})
			}

			if tt.fields.SocialBase != nil {
				tt.args.client = s.Client(context.Background(), token)
			}

			got, err := s.UserInfo(tt.args.client, token)
			if tt.wantErr {
				require.Error(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
